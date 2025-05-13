package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +subresource-request
type ProjectTemplates struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// DefaultVirtualClusterTemplate is the default template for the project
	DefaultVirtualClusterTemplate string `json:"defaultVirtualClusterTemplate,omitempty"`

	// VirtualClusterTemplates holds all the allowed virtual cluster templates
	VirtualClusterTemplates []VirtualClusterTemplate `json:"virtualClusterTemplates,omitempty"`

	// DefaultSpaceTemplate
	DefaultSpaceTemplate string `json:"defaultSpaceTemplate,omitempty"`

	// SpaceTemplates holds all the allowed space templates
	SpaceTemplates []SpaceTemplate `json:"spaceTemplates,omitempty"`

	// DefaultDevSpaceWorkspaceTemplate
	DefaultDevSpaceWorkspaceTemplate string `json:"defaultDevSpaceWorkspaceTemplate,omitempty"`

	// DevSpaceWorkspaceTemplates holds all the allowed space templates
	DevSpaceWorkspaceTemplates []DevSpaceWorkspaceTemplate `json:"devSpaceWorkspaceTemplates,omitempty"`

	// DevSpaceEnvironmentTemplates holds all the allowed environment templates
	DevSpaceEnvironmentTemplates []DevSpaceEnvironmentTemplate `json:"devSpaceEnvironmentTemplates,omitempty"`

	// DevSpaceWorkspacePresets holds all the allowed workspace presets
	DevSpaceWorkspacePresets []DevSpaceWorkspacePreset `json:"devSpaceWorkspacePresets,omitempty"`

	// DefaultDevSpaceEnvironmentTemplate
	DefaultDevSpaceEnvironmentTemplate string `json:"defaultDevSpaceEnvironmentTemplate,omitempty"`
}
