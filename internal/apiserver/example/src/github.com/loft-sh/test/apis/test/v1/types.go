package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterRole holds the cluster role information
// +k8s:openapi-gen=true
// +resource:path=clusterroles,rest=ClusterRoleREST
type ClusterRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterRoleSpec   `json:"spec,omitempty"`
	Status ClusterRoleStatus `json:"status,omitempty"`
}

// ClusterRoleSpec holds the user specification
type ClusterRoleSpec struct {
	// Rules holds all the PolicyRules for this ClusterRole
	// +optional
	Rules []string `json:"rules"`
}

// ClusterRoleStatus holds the user status
type ClusterRoleStatus struct {
}
