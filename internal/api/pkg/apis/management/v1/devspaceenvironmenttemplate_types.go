package v1

import (
	storagev1 "dev.khulnasoft.com/api/v4/pkg/apis/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevSpaceEnvironmentTemplate holds the DevSpaceEnvironmentTemplate information
// +k8s:openapi-gen=true
// +resource:path=devspaceenvironmenttemplates,rest=DevSpaceEnvironmentTemplateREST
type DevSpaceEnvironmentTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevSpaceEnvironmentTemplateSpec   `json:"spec,omitempty"`
	Status DevSpaceEnvironmentTemplateStatus `json:"status,omitempty"`
}

// DevSpaceEnvironmentTemplateSpec holds the specification
type DevSpaceEnvironmentTemplateSpec struct {
	storagev1.DevSpaceEnvironmentTemplateSpec `json:",inline"`
}

// DevSpaceEnvironmentTemplateStatus holds the status
type DevSpaceEnvironmentTemplateStatus struct{}

func (a *DevSpaceEnvironmentTemplate) GetVersions() []storagev1.VersionAccessor {
	var retVersions []storagev1.VersionAccessor
	for _, v := range a.Spec.Versions {
		b := v
		retVersions = append(retVersions, &b)
	}

	return retVersions
}

func (a *DevSpaceEnvironmentTemplate) GetOwner() *storagev1.UserOrTeam {
	return a.Spec.Owner
}

func (a *DevSpaceEnvironmentTemplate) SetOwner(userOrTeam *storagev1.UserOrTeam) {
	a.Spec.Owner = userOrTeam
}

func (a *DevSpaceEnvironmentTemplate) GetAccess() []storagev1.Access {
	return a.Spec.Access
}

func (a *DevSpaceEnvironmentTemplate) SetAccess(access []storagev1.Access) {
	a.Spec.Access = access
}
