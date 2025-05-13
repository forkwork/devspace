package v1

import (
	agentstoragev1 "dev.khulnasoft.com/agentapi/v4/pkg/apis/loft/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	DevSpaceWorkspaceConditions = []agentstoragev1.ConditionType{
		InstanceScheduled,
		InstanceTemplateResolved,
	}

	// DevSpaceWorkspaceIDLabel holds the actual workspace id of the devspace workspace
	DevSpaceWorkspaceIDLabel = "loft.sh/workspace-id"

	// DevSpaceWorkspaceUIDLabel holds the actual workspace uid of the devspace workspace
	DevSpaceWorkspaceUIDLabel = "loft.sh/workspace-uid"

	// DevSpaceKubernetesProviderWorkspaceUIDLabel holds the actual workspace uid of the devspace workspace on resources
	// created by the DevSpace Kubernetes provider.
	DevSpaceKubernetesProviderWorkspaceUIDLabel = "dev.khulnasoft.com/workspace-uid"

	// DevSpaceWorkspacePictureAnnotation holds the workspace picture url of the devspace workspace
	DevSpaceWorkspacePictureAnnotation = "loft.sh/workspace-picture"

	// DevSpaceWorkspaceSourceAnnotation holds the workspace source of the devspace workspace
	DevSpaceWorkspaceSourceAnnotation = "loft.sh/workspace-source"

	// DevSpaceWorkspaceRunnerNetworkPeerAnnotation holds the workspace runner network peer name of the devspace workspace
	DevSpaceWorkspaceRunnerEndpointAnnotation = "loft.sh/runner-endpoint"
)

var (
	DevSpaceFlagsUp     = "DEVSPACE_FLAGS_UP"
	DevSpaceFlagsDelete = "DEVSPACE_FLAGS_DELETE"
	DevSpaceFlagsStatus = "DEVSPACE_FLAGS_STATUS"
	DevSpaceFlagsSsh    = "DEVSPACE_FLAGS_SSH"
	DevSpaceFlagsStop   = "DEVSPACE_FLAGS_STOP"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevSpaceWorkspaceInstance
// +k8s:openapi-gen=true
type DevSpaceWorkspaceInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevSpaceWorkspaceInstanceSpec   `json:"spec,omitempty"`
	Status DevSpaceWorkspaceInstanceStatus `json:"status,omitempty"`
}

func (a *DevSpaceWorkspaceInstance) GetConditions() agentstoragev1.Conditions {
	return a.Status.Conditions
}

func (a *DevSpaceWorkspaceInstance) SetConditions(conditions agentstoragev1.Conditions) {
	a.Status.Conditions = conditions
}

func (a *DevSpaceWorkspaceInstance) GetOwner() *UserOrTeam {
	return a.Spec.Owner
}

func (a *DevSpaceWorkspaceInstance) SetOwner(userOrTeam *UserOrTeam) {
	a.Spec.Owner = userOrTeam
}

func (a *DevSpaceWorkspaceInstance) GetAccess() []Access {
	return a.Spec.Access
}

func (a *DevSpaceWorkspaceInstance) SetAccess(access []Access) {
	a.Spec.Access = access
}

type DevSpaceWorkspaceInstanceSpec struct {
	// DisplayName is the name that should be displayed in the UI
	// +optional
	DisplayName string `json:"displayName,omitempty"`

	// Description describes a DevSpace machine instance
	// +optional
	Description string `json:"description,omitempty"`

	// Owner holds the owner of this object
	// +optional
	Owner *UserOrTeam `json:"owner,omitempty"`

	// PresetRef holds the DevSpaceWorkspacePreset template reference
	// +optional
	PresetRef *PresetRef `json:"presetRef,omitempty"`

	// TemplateRef holds the DevSpace machine template reference
	// +optional
	TemplateRef *TemplateRef `json:"templateRef,omitempty"`

	// EnvironmentRef is the reference to DevSpaceEnvironmentTemplate that should be used
	// +optional
	EnvironmentRef *EnvironmentRef `json:"environmentRef,omitempty"`

	// Template is the inline template to use for DevSpace machine creation. This is mutually
	// exclusive with templateRef.
	// +optional
	Template *DevSpaceWorkspaceTemplateDefinition `json:"template,omitempty"`

	// RunnerRef is the reference to the connected runner holding
	// this workspace
	// +optional
	RunnerRef RunnerRef `json:"runnerRef,omitempty"`

	// Parameters are values to pass to the template.
	// The values should be encoded as YAML string where each parameter is represented as a top-level field key.
	// +optional
	Parameters string `json:"parameters,omitempty"`

	// Access to the DevSpace machine instance object itself
	// +optional
	Access []Access `json:"access,omitempty"`

	// PreventWakeUpOnConnection is used to prevent workspace that uses sleep mode from waking up on incomming ssh connection.
	// +optional
	PreventWakeUpOnConnection bool `json:"preventWakeUpOnConnection,omitempty"`
}

type PresetRef struct {
	// Name is the name of DevSpaceWorkspacePreset
	Name string `json:"name"`

	// Version holds the preset version to use. Version is expected to
	// be in semantic versioning format. Alternatively, you can also exchange
	// major, minor or patch with an 'x' to tell Loft to automatically select
	// the latest major, minor or patch version.
	// +optional
	Version string `json:"version,omitempty"`
}

type RunnerRef struct {
	// Runner is the connected runner the workspace will be created in
	// +optional
	Runner string `json:"runner,omitempty"`
}

type EnvironmentRef struct {
	// Name is the name of DevSpaceEnvironmentTemplate this references
	Name string `json:"name"`

	// Version is the version of DevSpaceEnvironmentTemplate this references
	// +optional
	Version string `json:"version,omitempty"`
}

type DevSpaceWorkspaceInstanceStatus struct {
	// LastWorkspaceStatus is the last workspace status reported by the runner.
	// +optional
	LastWorkspaceStatus WorkspaceStatus `json:"lastWorkspaceStatus,omitempty"`

	// Phase describes the current phase the DevSpace machine instance is in
	// +optional
	Phase InstancePhase `json:"phase,omitempty"`

	// Reason describes the reason in machine-readable form why the cluster is in the current
	// phase
	// +optional
	Reason string `json:"reason,omitempty"`

	// Message describes the reason in human-readable form why the DevSpace machine is in the current
	// phase
	// +optional
	Message string `json:"message,omitempty"`

	// Conditions holds several conditions the DevSpace machine might be in
	// +optional
	Conditions agentstoragev1.Conditions `json:"conditions,omitempty"`

	// Instance is the template rendered with all the parameters
	// +optional
	Instance *DevSpaceWorkspaceTemplateDefinition `json:"instance,omitempty"`

	// IgnoreReconciliation ignores reconciliation for this object
	// +optional
	IgnoreReconciliation bool `json:"ignoreReconciliation,omitempty"`

	// ClusterRef holds the runners cluster if the workspace is scheduled
	// on kubernetes based runner
	ClusterRef *ClusterRef `json:"clusterRef,omitempty"`
}

type WorkspaceStatusResult struct {
	ID       string `json:"id,omitempty"`
	Context  string `json:"context,omitempty"`
	Provider string `json:"provider,omitempty"`
	State    string `json:"state,omitempty"`
}

var AllowedWorkspaceStatus = []WorkspaceStatus{
	WorkspaceStatusNotFound,
	WorkspaceStatusStopped,
	WorkspaceStatusBusy,
	WorkspaceStatusRunning,
}

type WorkspaceStatus string

var (
	WorkspaceStatusNotFound WorkspaceStatus = "NotFound"
	WorkspaceStatusStopped  WorkspaceStatus = "Stopped"
	WorkspaceStatusBusy     WorkspaceStatus = "Busy"
	WorkspaceStatusRunning  WorkspaceStatus = "Running"
)

type DevSpaceCommandStopOptions struct{}

type DevSpaceCommandDeleteOptions struct {
	IgnoreNotFound bool   `json:"ignoreNotFound,omitempty"`
	Force          bool   `json:"force,omitempty"`
	GracePeriod    string `json:"gracePeriod,omitempty"`
}

type DevSpaceCommandStatusOptions struct {
	ContainerStatus bool `json:"containerStatus,omitempty"`
}

type DevSpaceCommandUpOptions struct {
	// up options
	ID                   string   `json:"id,omitempty"`
	Source               string   `json:"source,omitempty"`
	IDE                  string   `json:"ide,omitempty"`
	IDEOptions           []string `json:"ideOptions,omitempty"`
	PrebuildRepositories []string `json:"prebuildRepositories,omitempty"`
	DevContainerPath     string   `json:"devContainerPath,omitempty"`
	WorkspaceEnv         []string `json:"workspaceEnv,omitempty"`
	Recreate             bool     `json:"recreate,omitempty"`
	Proxy                bool     `json:"proxy,omitempty"`
	DisableDaemon        bool     `json:"disableDaemon,omitempty"`
	DaemonInterval       string   `json:"daemonInterval,omitempty"`

	// build options
	Repository string   `json:"repository,omitempty"`
	SkipPush   bool     `json:"skipPush,omitempty"`
	Platform   []string `json:"platform,omitempty"`

	// TESTING
	ForceBuild            bool `json:"forceBuild,omitempty"`
	ForceInternalBuildKit bool `json:"forceInternalBuildKit,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevSpaceWorkspaceInstanceList contains a list of DevSpaceWorkspaceInstance objects
type DevSpaceWorkspaceInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DevSpaceWorkspaceInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DevSpaceWorkspaceInstance{}, &DevSpaceWorkspaceInstanceList{})
}
