package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevSpaceWorkspaceTemplate holds the DevSpaceWorkspaceTemplate information
// +k8s:openapi-gen=true
type DevSpaceWorkspaceTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevSpaceWorkspaceTemplateSpec   `json:"spec,omitempty"`
	Status DevSpaceWorkspaceTemplateStatus `json:"status,omitempty"`
}

func (a *DevSpaceWorkspaceTemplate) GetVersions() []VersionAccessor {
	var retVersions []VersionAccessor
	for _, v := range a.Spec.Versions {
		b := v
		retVersions = append(retVersions, &b)
	}

	return retVersions
}

func (a *DevSpaceWorkspaceTemplateVersion) GetVersion() string {
	return a.Version
}

func (a *DevSpaceWorkspaceTemplate) GetOwner() *UserOrTeam {
	return a.Spec.Owner
}

func (a *DevSpaceWorkspaceTemplate) SetOwner(userOrTeam *UserOrTeam) {
	a.Spec.Owner = userOrTeam
}

func (a *DevSpaceWorkspaceTemplate) GetAccess() []Access {
	return a.Spec.Access
}

func (a *DevSpaceWorkspaceTemplate) SetAccess(access []Access) {
	a.Spec.Access = access
}

// DevSpaceWorkspaceTemplateSpec holds the specification
type DevSpaceWorkspaceTemplateSpec struct {
	// DisplayName is the name that is shown in the UI
	// +optional
	DisplayName string `json:"displayName,omitempty"`

	// Description describes the virtual cluster template
	// +optional
	Description string `json:"description,omitempty"`

	// Owner holds the owner of this object
	// +optional
	Owner *UserOrTeam `json:"owner,omitempty"`

	// Parameters define additional app parameters that will set provider values
	// +optional
	Parameters []AppParameter `json:"parameters,omitempty"`

	// Template holds the DevSpace workspace template
	Template DevSpaceWorkspaceTemplateDefinition `json:"template,omitempty"`

	// Versions are different versions of the template that can be referenced as well
	// +optional
	Versions []DevSpaceWorkspaceTemplateVersion `json:"versions,omitempty"`

	// Access holds the access rights for users and teams
	// +optional
	Access []Access `json:"access,omitempty"`
}

type DevSpaceWorkspaceTemplateDefinition struct {
	// Provider holds the DevSpace provider configuration
	Provider DevSpaceWorkspaceProvider `json:"provider"`

	// SpaceTemplateRef is a reference to the space that should get created for this DevSpace.
	// If this is specified, the kubernetes provider will be selected automatically.
	// +optional
	SpaceTemplateRef *TemplateRef `json:"spaceTemplateRef,omitempty"`

	// SpaceTemplate is the inline template for a space that should get created for this DevSpace.
	// If this is specified, the kubernetes provider will be selected automatically.
	// +optional
	SpaceTemplate *SpaceTemplateDefinition `json:"spaceTemplate,omitempty"`

	// VirtualClusterTemplateRef is a reference to the virtual cluster that should get created for this DevSpace.
	// If this is specified, the kubernetes provider will be selected automatically.
	// +optional
	VirtualClusterTemplateRef *TemplateRef `json:"virtualClusterTemplateRef,omitempty"`

	// VirtualClusterTemplate is the inline template for a virtual cluster that should get created for this DevSpace.
	// If this is specified, the kubernetes provider will be selected automatically.
	// +optional
	VirtualClusterTemplate *VirtualClusterTemplateDefinition `json:"virtualClusterTemplate,omitempty"`

	// WorkspaceEnv are environment variables that should be available within the created workspace.
	// +optional
	WorkspaceEnv map[string]DevSpaceProviderOption `json:"workspaceEnv,omitempty"`

	// InitEnv are environment variables that should be available during the initialization phase of the created workspace.
	// +optional
	InitEnv map[string]DevSpaceProviderOption `json:"initEnv,omitempty"`

	// InstanceTemplate holds the workspace instance template
	// +optional
	InstanceTemplate DevSpaceWorkspaceInstanceTemplateDefinition `json:"instanceTemplate,omitempty"`

	// UseProjectGitCredentials specifies if the project git credentials should be used instead of local ones for this workspace
	// +optional
	UseProjectGitCredentials bool `json:"useProjectGitCredentials,omitempty"`

	// UseProjectSSHCredentials specifies if the project ssh credentials should be used instead of local ones for this workspace
	// +optional
	UseProjectSSHCredentials bool `json:"useProjectSSHCredentials,omitempty"`

	// GitCloneStrategy specifies how git based workspace are being cloned. Can be "" (full, default), treeless, blobless or shallow
	// +optional
	GitCloneStrategy GitCloneStrategy `json:"gitCloneStrategy,omitempty"`

	// CredentialForwarding specifies controls for how workspaces created by this template forward credentials into the workspace
	// +optional
	CredentialForwarding *CredentialForwarding `json:"credentialForwarding,omitempty"`

	// PreventWakeUpOnConnection is used to prevent workspace that uses sleep mode from waking up on incomming ssh connection.
	// +optional
	PreventWakeUpOnConnection bool `json:"preventWakeUpOnConnection,omitempty"`
}

// +enum
type GitCloneStrategy string

// WARN: Need to match https://dev.khulnasoft.com/devspace/pkg/git/clone.go
const (
	FullCloneStrategy     GitCloneStrategy = ""
	BloblessCloneStrategy GitCloneStrategy = "blobless"
	TreelessCloneStrategy GitCloneStrategy = "treeless"
	ShallowCloneStrategy  GitCloneStrategy = "shallow"
)

type CredentialForwarding struct {
	// Docker specifies controls for how workspaces created by this template forward docker credentials
	// +optional
	Docker *DockerCredentialForwarding `json:"docker,omitempty"`

	// Git specifies controls for how workspaces created by this template forward git credentials
	// +optional
	Git *GitCredentialForwarding `json:"git,omitempty"`
}

type DockerCredentialForwarding struct {
	// Disabled prevents all workspaces created by this template from forwarding credentials into the workspace
	// +optional
	Disabled bool `json:"disabled,omitempty"`
}

type GitCredentialForwarding struct {
	// Disabled prevents all workspaces created by this template from forwarding credentials into the workspace
	// +optional
	Disabled bool `json:"disabled,omitempty"`
}

type DevSpaceWorkspaceProvider struct {
	// Name is the name of the provider. This can also be an url.
	Name string `json:"name"`

	// Options are the provider option values
	// +optional
	Options map[string]DevSpaceProviderOption `json:"options,omitempty"`

	// Env are environment options to set when using the provider.
	// +optional
	Env map[string]DevSpaceProviderOption `json:"env,omitempty"`
}

type DevSpaceWorkspaceInstanceTemplateDefinition struct {
	// The virtual cluster instance metadata
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	TemplateMetadata `json:"metadata,omitempty"`
}

type DevSpaceProviderOption struct {
	// Value of this option.
	// +optional
	Value string `json:"value,omitempty"`

	// ValueFrom specifies a secret where this value should be taken from.
	// +optional
	ValueFrom *DevSpaceProviderOptionFrom `json:"valueFrom,omitempty"`
}

type DevSpaceProviderOptionFrom struct {
	// ProjectSecretRef is the project secret to use for this value.
	// +optional
	ProjectSecretRef *corev1.SecretKeySelector `json:"projectSecretRef,omitempty"`

	// SharedSecretRef is the shared secret to use for this value.
	// +optional
	SharedSecretRef *corev1.SecretKeySelector `json:"sharedSecretRef,omitempty"`
}

type DevSpaceProviderSource struct {
	// Github source for the provider
	Github string `json:"github,omitempty"`

	// File source for the provider
	File string `json:"file,omitempty"`

	// URL where the provider was downloaded from
	URL string `json:"url,omitempty"`
}

type DevSpaceWorkspaceTemplateVersion struct {
	// Template holds the DevSpace template
	// +optional
	Template DevSpaceWorkspaceTemplateDefinition `json:"template,omitempty"`

	// Parameters define additional app parameters that will set provider values
	// +optional
	Parameters []AppParameter `json:"parameters,omitempty"`

	// Version is the version. Needs to be in X.X.X format.
	// +optional
	Version string `json:"version,omitempty"`
}

// DevSpaceWorkspaceTemplateStatus holds the status
type DevSpaceWorkspaceTemplateStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevSpaceWorkspaceTemplateList contains a list of DevSpaceWorkspaceTemplate
type DevSpaceWorkspaceTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DevSpaceWorkspaceTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DevSpaceWorkspaceTemplate{}, &DevSpaceWorkspaceTemplateList{})
}
