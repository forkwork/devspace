package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevSpaceWorkspaceEnvironmentSource
// +k8s:openapi-gen=true
type DevSpaceEnvironmentTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevSpaceEnvironmentTemplateSpec   `json:"spec,omitempty"`
	Status DevSpaceEnvironmentTemplateStatus `json:"status,omitempty"`
}

// DevSpaceEnvironmentTemplateStatus holds the status
type DevSpaceEnvironmentTemplateStatus struct {
}

func (a *DevSpaceEnvironmentTemplate) GetVersions() []VersionAccessor {
	var retVersions []VersionAccessor
	for _, v := range a.Spec.Versions {
		b := v
		retVersions = append(retVersions, &b)
	}

	return retVersions
}

func (a *DevSpaceEnvironmentTemplateVersion) GetVersion() string {
	return a.Version
}

func (a *DevSpaceEnvironmentTemplate) GetOwner() *UserOrTeam {
	return a.Spec.Owner
}

func (a *DevSpaceEnvironmentTemplate) SetOwner(userOrTeam *UserOrTeam) {
	a.Spec.Owner = userOrTeam
}

func (a *DevSpaceEnvironmentTemplate) GetAccess() []Access {
	return a.Spec.Access
}

func (a *DevSpaceEnvironmentTemplate) SetAccess(access []Access) {
	a.Spec.Access = access
}

type DevSpaceEnvironmentTemplateSpec struct {
	// DisplayName is the name that should be displayed in the UI
	// +optional
	DisplayName string `json:"displayName,omitempty"`

	// Description describes the environment template
	// +optional
	Description string `json:"description,omitempty"`

	// Owner holds the owner of this object
	// +optional
	Owner *UserOrTeam `json:"owner,omitempty"`

	// Access to the DevSpace machine instance object itself
	// +optional
	Access []Access `json:"access,omitempty"`

	// Template is the inline template to use for DevSpace environments
	// +optional
	Template *DevSpaceEnvironmentTemplateDefinition `json:"template,omitempty"`

	// Versions are different versions of the template that can be referenced as well
	// +optional
	Versions []DevSpaceEnvironmentTemplateVersion `json:"versions,omitempty"`
}

type DevSpaceEnvironmentTemplateDefinition struct {
	// Git holds configuration for git environment spec source
	// +optional
	Git *GitEnvironmentTemplate `json:"git,omitempty"`

	// Inline holds an inline devcontainer.json definition
	// +optional
	Inline string `json:"inline,omitempty"`
}

// GitEnvironmentTemplate stores configuration of Git environment template source
type GitEnvironmentTemplate struct {
	// Repository stores repository URL for Git environment spec source
	Repository string `json:"repository"`

	// Revision stores revision to checkout in repository
	// +optional
	Revision string `json:"revision,omitempty"`

	// SubPath stores subpath within Repositor where environment spec is
	// +optional
	SubPath string `json:"subpath,omitempty"`

	// UseProjectGitCredentials specifies if the project git credentials should be used instead of local ones for this environment
	// +optional
	UseProjectGitCredentials bool `json:"useProjectGitCredentials,omitempty"`
}

type DevSpaceEnvironmentTemplateVersion struct {
	// Template holds the environment template definition
	// +optional
	Template DevSpaceEnvironmentTemplateDefinition `json:"template,omitempty"`

	// Version is the version. Needs to be in X.X.X format.
	// +optional
	Version string `json:"version,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevSpaceEnvironmentTemplateList contains a list of DevSpaceEnvironmentTemplate objects
type DevSpaceEnvironmentTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DevSpaceEnvironmentTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DevSpaceEnvironmentTemplate{}, &DevSpaceEnvironmentTemplateList{})
}
