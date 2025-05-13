import { ManagementV1DevSpaceWorkspacePreset } from "@loft-enterprise/client/gen/models/managementV1DevSpaceWorkspacePreset"

export function presetDisplayName(preset: ManagementV1DevSpaceWorkspacePreset | undefined) {
  return preset?.spec?.displayName ?? preset?.metadata?.name
}
