package v1

import (
	storagev1 "dev.khulnasoft.com/api/v4/pkg/apis/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevSpaceWorkspacePreset
// +k8s:openapi-gen=true
// +resource:path=devspaceworkspacepresets,rest=DevSpaceWorkspacePresetREST
type DevSpaceWorkspacePreset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevSpaceWorkspacePresetSpec   `json:"spec,omitempty"`
	Status DevSpaceWorkspacePresetStatus `json:"status,omitempty"`
}

// DevSpaceWorkspacePresetSpec holds the specification
type DevSpaceWorkspacePresetSpec struct {
	storagev1.DevSpaceWorkspacePresetSpec `json:",inline"`
}

// DevSpaceWorkspacePresetSource
// +k8s:openapi-gen=true
type DevSpaceWorkspacePresetSource struct {
	storagev1.DevSpaceWorkspacePresetSource `json:",inline"`
}

func (a *DevSpaceWorkspacePreset) GetOwner() *storagev1.UserOrTeam {
	return a.Spec.Owner
}

func (a *DevSpaceWorkspacePreset) SetOwner(userOrTeam *storagev1.UserOrTeam) {
	a.Spec.Owner = userOrTeam
}

func (a *DevSpaceWorkspacePreset) GetAccess() []storagev1.Access {
	return a.Spec.Access
}

func (a *DevSpaceWorkspacePreset) SetAccess(access []storagev1.Access) {
	a.Spec.Access = access
}

// DevSpaceWorkspacePresetStatus holds the status
type DevSpaceWorkspacePresetStatus struct{}
