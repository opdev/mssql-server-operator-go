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
	"k8s.io/apimachinery/pkg/api/errors"
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
	if err := r.Get(ctx,
		types.NamespacedName{Name: mssql.Name, Namespace: mssql.Namespace},
		mssqlStatefulSet)
	err != nil && errors.IsNotFound(err) {

	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MsSqlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1alpha1.MsSql{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}
