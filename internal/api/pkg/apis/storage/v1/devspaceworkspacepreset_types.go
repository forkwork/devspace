package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevSpaceWorkspacePreset
// +k8s:openapi-gen=true
type DevSpaceWorkspacePreset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevSpaceWorkspacePresetSpec   `json:"spec,omitempty"`
	Status DevSpaceWorkspacePresetStatus `json:"status,omitempty"`
}

func (a *DevSpaceWorkspacePreset) GetOwner() *UserOrTeam {
	return a.Spec.Owner
}

func (a *DevSpaceWorkspacePreset) SetOwner(userOrTeam *UserOrTeam) {
	a.Spec.Owner = userOrTeam
}

func (a *DevSpaceWorkspacePreset) GetAccess() []Access {
	return a.Spec.Access
}

func (a *DevSpaceWorkspacePreset) SetAccess(access []Access) {
	a.Spec.Access = access
}

type DevSpaceWorkspacePresetSpec struct {
	// DisplayName is the name that should be displayed in the UI
	// +optional
	DisplayName string `json:"displayName,omitempty"`

	// Source stores inline path of project source
	Source *DevSpaceWorkspacePresetSource `json:"source"`

	// InfrastructureRef stores reference to DevSpaceWorkspaceTemplate to use
	InfrastructureRef *TemplateRef `json:"infrastructureRef"`

	// EnvironmentRef stores reference to DevSpaceEnvironmentTemplate
	// +optional
	EnvironmentRef *EnvironmentRef `json:"environmentRef,omitempty"`

	// UseProjectGitCredentials specifies if the project git credentials should be used instead of local ones for this environment
	// +optional
	UseProjectGitCredentials bool `json:"useProjectGitCredentials,omitempty"`

	// Owner holds the owner of this object
	// +optional
	Owner *UserOrTeam `json:"owner,omitempty"`

	// Access to the DevSpace machine instance object itself
	// +optional
	Access []Access `json:"access,omitempty"`

	// Versions are different versions of the template that can be referenced as well
	// +optional
	Versions []DevSpaceWorkspacePresetVersion `json:"versions,omitempty"`
}

type DevSpaceWorkspacePresetSource struct {
	// Git stores path to git repo to use as workspace source
	// +optional
	Git string `json:"git,omitempty"`

	// Image stores container image to use as workspace source
	// +optional
	Image string `json:"image,omitempty"`
}

type DevSpaceWorkspacePresetVersion struct {
	// Version is the version. Needs to be in X.X.X format.
	// +optional
	Version string `json:"version,omitempty"`

	// Source stores inline path of project source
	// +optional
	Source *DevSpaceWorkspacePresetSource `json:"source,omitempty"`

	// InfrastructureRef stores reference to DevSpaceWorkspaceTemplate to use
	// +optional
	InfrastructureRef *TemplateRef `json:"infrastructureRef,omitempty"`

	// EnvironmentRef stores reference to DevSpaceEnvironmentTemplate
	// +optional
	EnvironmentRef *EnvironmentRef `json:"environmentRef,omitempty"`

	// UseProjectGitCredentials specifies if the project git credentials should be used instead of local ones for this environment
	// +optional
	UseProjectGitCredentials bool `json:"useProjectGitCredentials,omitempty"`
}

// DevSpaceWorkspacePresetStatus holds the status
type DevSpaceWorkspacePresetStatus struct {
}

type WorkspaceRef struct {
	// Name is the name of DevSpaceWorkspaceTemplate this references
	Name string `json:"name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// DevSpaceWorkspacePresetList contains a list of DevSpaceWorkspacePreset objects
type DevSpaceWorkspacePresetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DevSpaceWorkspacePreset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DevSpaceWorkspacePreset{}, &DevSpaceWorkspacePresetList{})
}
