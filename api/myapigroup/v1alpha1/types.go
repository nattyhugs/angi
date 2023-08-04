package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Custom types and methods

// MyAppResourceSpec defines the desired state of MyAppResource
type MyAppResourceSpec struct {
	Image        ImageSpec    `json:"image"`
	Redis        RedisSpec    `json:"redis"`
	ReplicaCount int32        `json:"replicaCount"`
	Resources    ResourceSpec `json:"resources"`
	Ui           UiSpec       `json:"ui"`
}

type ImageSpec struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
}

type RedisSpec struct {
	Enabled bool `json:"enabled"`
}

type ResourceSpec struct {
	CpuRequest  string `json:"cpuRequest"`
	MemoryLimit string `json:"memoryLimit"`
}

type UiSpec struct {
	Color   string `json:"color"`
	Message string `json:"message"`
}

// MyAppResource is the Schema for the myappresources API
type MyAppResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MyAppResourceSpec `json:"spec,omitempty"`
}

// MyAppResourceList contains a list of MyAppResource
type MyAppResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyAppResource `json:"items"`
}

// DeepCopyObject implements runtime.Object interface for MyAppResource.
func (in *MyAppResource) DeepCopyObject() runtime.Object {
	out := in.DeepCopy()
	return out
}

// DeepCopyObject implements runtime.Object interface for MyAppResourceList.
func (in *MyAppResourceList) DeepCopyObject() runtime.Object {
	out := in.DeepCopy()
	return out
}

// DeepCopy returns a deep copy of the MyAppResource.
func (in *MyAppResource) DeepCopy() *MyAppResource {
	out := in.DeepCopy()
	return out
}

// DeepCopy returns a deep copy of the MyAppResourceList.
func (in *MyAppResourceList) DeepCopy() *MyAppResourceList {
	out := in.DeepCopy()
	return out
}
