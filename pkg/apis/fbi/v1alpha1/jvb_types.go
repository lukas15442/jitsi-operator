package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JVBSpec defines the desired state of JVB
type JVBSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Size          int32             `json:"size"`
	Name          string            `json:"name" protobuf:"bytes,1,opt,name=name"`
	JVBDeployment appsv1.Deployment `json:"jvbDeployment"`
	Services      []corev1.Service  `json:"services"`
}

// JVBStatus defines the observed state of JVB
type JVBStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JVB is the Schema for the jvbs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=jvbs,scope=Namespaced
type JVB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JVBSpec   `json:"spec,omitempty"`
	Status JVBStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JVBList contains a list of JVB
type JVBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JVB `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JVB{}, &JVBList{})
}
