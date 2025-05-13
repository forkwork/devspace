package v1

import (
	clusterv1 "dev.khulnasoft.com/agentapi/v4/pkg/apis/loft/cluster/v1"
	agentstoragev1 "dev.khulnasoft.com/agentapi/v4/pkg/apis/loft/storage/v1"
	storagev1 "dev.khulnasoft.com/api/v4/pkg/apis/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +genclient:method=GetState,verb=get,subresource=state,result=dev.khulnasoft.com/api/v4/pkg/apis/management/v1.DevSpaceWorkspaceInstanceState
// +genclient:method=SetState,verb=create,subresource=state,input=dev.khulnasoft.com/api/v4/pkg/apis/management/v1.DevSpaceWorkspaceInstanceState,result=dev.khulnasoft.com/api/v4/pkg/apis/management/v1.DevSpaceWorkspaceInstanceState
// +genclient:method=Troubleshoot,verb=get,subresource=troubleshoot,result=dev.khulnasoft.com/api/v4/pkg/apis/management/v1.DevSpaceWorkspaceInstanceTroubleshoot
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevSpaceWorkspaceInstance holds the DevSpaceWorkspaceInstance information
// +k8s:openapi-gen=true
// +resource:path=devspaceworkspaceinstances,rest=DevSpaceWorkspaceInstanceREST
// +subresource:request=DevSpaceUpOptions,path=up,kind=DevSpaceUpOptions,rest=DevSpaceUpOptionsREST
// +subresource:request=DevSpaceDeleteOptions,path=delete,kind=DevSpaceDeleteOptions,rest=DevSpaceDeleteOptionsREST
// +subresource:request=DevSpaceSshOptions,path=ssh,kind=DevSpaceSshOptions,rest=DevSpaceSshOptionsREST
// +subresource:request=DevSpaceStopOptions,path=stop,kind=DevSpaceStopOptions,rest=DevSpaceStopOptionsREST
// +subresource:request=DevSpaceStatusOptions,path=getstatus,kind=DevSpaceStatusOptions,rest=DevSpaceStatusOptionsREST
// +subresource:request=DevSpaceWorkspaceInstanceState,path=state,kind=DevSpaceWorkspaceInstanceState,rest=DevSpaceWorkspaceInstanceStateREST
// +subresource:request=DevSpaceWorkspaceInstanceTroubleshoot,path=troubleshoot,kind=DevSpaceWorkspaceInstanceTroubleshoot,rest=DevSpaceWorkspaceInstanceTroubleshootREST
type DevSpaceWorkspaceInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevSpaceWorkspaceInstanceSpec   `json:"spec,omitempty"`
	Status DevSpaceWorkspaceInstanceStatus `json:"status,omitempty"`
}

// DevSpaceWorkspaceInstanceSpec holds the specification
type DevSpaceWorkspaceInstanceSpec struct {
	storagev1.DevSpaceWorkspaceInstanceSpec `json:",inline"`
}

// DevSpaceWorkspaceInstanceStatus holds the status
type DevSpaceWorkspaceInstanceStatus struct {
	storagev1.DevSpaceWorkspaceInstanceStatus `json:",inline"`

	// SleepModeConfig is the sleep mode config of the workspace. This will only be shown
	// in the front end.
	// +optional
	SleepModeConfig *clusterv1.SleepModeConfig `json:"sleepModeConfig,omitempty"`
}

func (a *DevSpaceWorkspaceInstance) GetConditions() agentstoragev1.Conditions {
	return a.Status.Conditions
}

func (a *DevSpaceWorkspaceInstance) SetConditions(conditions agentstoragev1.Conditions) {
	a.Status.Conditions = conditions
}

func (a *DevSpaceWorkspaceInstance) GetOwner() *storagev1.UserOrTeam {
	return a.Spec.Owner
}

func (a *DevSpaceWorkspaceInstance) SetOwner(userOrTeam *storagev1.UserOrTeam) {
	a.Spec.Owner = userOrTeam
}

func (a *DevSpaceWorkspaceInstance) GetAccess() []storagev1.Access {
	return a.Spec.Access
}

func (a *DevSpaceWorkspaceInstance) SetAccess(access []storagev1.Access) {
	a.Spec.Access = access
}
