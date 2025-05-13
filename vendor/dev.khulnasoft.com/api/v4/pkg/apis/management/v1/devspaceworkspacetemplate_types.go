package v1

import (
	storagev1 "dev.khulnasoft.com/api/v4/pkg/apis/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevSpaceWorkspaceTemplate holds the information
// +k8s:openapi-gen=true
// +resource:path=devspaceworkspacetemplates,rest=DevSpaceWorkspaceTemplateREST
type DevSpaceWorkspaceTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevSpaceWorkspaceTemplateSpec   `json:"spec,omitempty"`
	Status DevSpaceWorkspaceTemplateStatus `json:"status,omitempty"`
}

// DevSpaceWorkspaceTemplateSpec holds the specification
type DevSpaceWorkspaceTemplateSpec struct {
	storagev1.DevSpaceWorkspaceTemplateSpec `json:",inline"`
}

// DevSpaceWorkspaceTemplateStatus holds the status
type DevSpaceWorkspaceTemplateStatus struct {
	storagev1.DevSpaceWorkspaceTemplateStatus `json:",inline"`
}

func (a *DevSpaceWorkspaceTemplate) GetVersions() []storagev1.VersionAccessor {
	var retVersions []storagev1.VersionAccessor
	for _, v := range a.Spec.Versions {
		b := v
		retVersions = append(retVersions, &b)
	}

	return retVersions
}

func (a *DevSpaceWorkspaceTemplate) GetOwner() *storagev1.UserOrTeam {
	return a.Spec.Owner
}

func (a *DevSpaceWorkspaceTemplate) SetOwner(userOrTeam *storagev1.UserOrTeam) {
	a.Spec.Owner = userOrTeam
}

func (a *DevSpaceWorkspaceTemplate) GetAccess() []storagev1.Access {
	return a.Spec.Access
}

func (a *DevSpaceWorkspaceTemplate) SetAccess(access []storagev1.Access) {
	a.Spec.Access = access
}
