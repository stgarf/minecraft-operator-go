package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MinecraftSpec defines the desired state of Minecraft
// +k8s:openapi-gen=true
type MinecraftSpec struct {
	// Version is the version of the minecraft deployment
	Version string `json:"size"`
}

// MinecraftStatus defines the observed state of Minecraft
// +k8s:openapi-gen=true
type MinecraftStatus struct {
	// Nodes are the names of the minecraft pods
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Minecraft is the Schema for the minecrafts API
// +k8s:openapi-gen=true
type Minecraft struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MinecraftSpec   `json:"spec,omitempty"`
	Status MinecraftStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MinecraftList contains a list of Minecraft
type MinecraftList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Minecraft `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Minecraft{}, &MinecraftList{})
}
