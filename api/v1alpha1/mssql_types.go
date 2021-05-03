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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MsSqlSpec defines the desired state of MsSql
type MsSqlSpec struct {
	Replicas int32 `json:"replicas"`
	Eula bool `json:"eula"`
	Hostname string `json:"hostname"`
	SqlAgentEnabled bool `json:"sqlagent"`
	SqlServerEdition string `json:"sqlserveredition"`
	SqlPassword string `json:"sqlpassword"`
}

// MsSqlStatus defines the observed state of MsSql
type MsSqlStatus struct {
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []MsSqlCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,6,rep,name=conditions"`
}

type MsSqlConditionType string

type MsSqlCondition struct {
	Type MsSqlConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=MsSqlConditionType"`
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty" protobuf:"bytes,6,opt,name=lastUpdateTime"`
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,7,opt,name=lastTransitionTime"`
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MsSql is the Schema for the mssqls API
type MsSql struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MsSqlSpec   `json:"spec,omitempty"`
	Status MsSqlStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MsSqlList contains a list of MsSql
type MsSqlList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MsSql `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MsSql{}, &MsSqlList{})
}
