/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	databasev1alpha1 "github.com/opdev/mssql-server-operator-go/api/v1alpha1"
)

// MsSqlReconciler reconciles a MsSql object
type MsSqlReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=database.microsoft.com,resources=mssqls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=database.microsoft.com,resources=mssqls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=database.microsoft.com,resources=mssqls/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MsSql object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *MsSqlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logv := r.Log.WithValues("mssql", req.NamespacedName)

	mssql := &databasev1alpha1.MsSql{}

	if err := r.Get(ctx, req.NamespacedName, mssql); err != nil {
		if errors.IsNotFound(err) {
			logv.Info("MsSql resource not found. Ignoring...")
			return ctrl.Result{}, nil
		}

		logv.Error(err, "Failed to get MsSql")
		return ctrl.Result{}, err
	}

	mssqlStatefulSet := &appsv1.StatefulSet{}
	err := r.Get(ctx, types.NamespacedName{Name: mssql.Name, Namespace: mssql.Namespace}, mssqlStatefulSet)
	if err != nil && errors.IsNotFound(err) {
		ss := r.statefulSetForMsSql(mssql)
		if err = r.Create(ctx, ss); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		logv.Error(err, "Failed to get StatefulSet")
		return ctrl.Result{}, err
	}

	_replicas := mssql.Spec.Replicas
	if *mssqlStatefulSet.Spec.Replicas != _replicas {
		mssqlStatefulSet.Spec.Replicas = &_replicas
		err = r.Update(ctx, mssqlStatefulSet); if err != nil {
			logv.Error(err, "Failed to update StatefulSet", "StatefulSet.Namespace", mssqlStatefulSet.Namespace, "StatefulSet.Name", mssqlStatefulSet.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	_mssqlStatefulSetPods := &corev1.PodList{}
	_opts := []client.ListOption{
		client.InNamespace(mssql.Namespace),
		client.MatchingLabels{"app": "mssql", "mssql": mssql.Name},
	}

	if err = r.List(ctx, _mssqlStatefulSetPods, _opts...); err != nil {
		logv.Error(err, "Failed to list pods", "MsSql.Namespace", mssql.Namespace, "MsSql.Name", mssql.Name)
		return ctrl.Result{}, err
	}

	//_mssqlStatefulSetPodNames := getPodNames(_mssqlStatefulSetPods.Items)
	_ = getPodNames(_mssqlStatefulSetPods.Items)

	return ctrl.Result{}, nil
}

// StatefulSetForMsSql creates a statefulset for mssql custom resource when one is not available
func (r *MsSqlReconciler) statefulSetForMsSql(mssql *databasev1alpha1.MsSql) *appsv1.StatefulSet {
	_labels := make(map[string]string)
	_replicas := mssql.Spec.Replicas
	var _mssqlFsGroup int64 = 10001
	var _containerPort int32 = 1433

	// mssqlserver pod specification
	_mssqlPodSpec := corev1.PodSpec{
		Volumes: []corev1.Volume{
			{
				Name: "mssql-config-volume",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{Name: "mssql"},
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Name: "mssqlserver",
				Command: []string{"/bin/bash", "-c",
					"cp /var/opt/config/mssql.conf /var/opt/mssql/mssql.conf && /opt/mssql/bin/sqlserver"},
				Image:           "mcr.microsoft.com/mssql/server:2019-latest",
				ImagePullPolicy: "IfNotPresent",
				Ports: []corev1.ContainerPort{
					{
						ContainerPort: _containerPort,
					},
				},
				Env: []corev1.EnvVar{
					{Name: "MSSQL_PID", Value: "Developer"},
					{Name: "ACCEPT_EULA", Value: "Y"},
					{Name: "MSSQL_AGENT_ENABLED", Value: "false"},
					{
						Name: "SA_PASSWORD",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: "mssql"},
								Key:                  "SA_PASSWORD",
							},
						},
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{Name: mssql.Name, MountPath: "/var/opt/mssql"},
					{Name: "mssql-config-volume", MountPath: "/var/opt/config"},
				},
			},
		},
		SecurityContext: &corev1.PodSecurityContext{
			FSGroup: &_mssqlFsGroup,
		},
	}

	// storage for mssqlserver pod
	_pvsize := resource.MustParse("8Gi")
	_pvclaim := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: mssql.Name,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					"storage": _pvsize,
				},
			},
		},
	}

	// statefulset for mssqlserver operand
	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mssql.Name,
			Namespace: mssql.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &_replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: _labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{},
					Labels:      map[string]string{},
				},
				Spec: _mssqlPodSpec,
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				_pvclaim,
			},
		},
	}

	return ss
}

func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// SetupWithManager sets up the controller with the Manager.
func (r *MsSqlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1alpha1.MsSql{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}
