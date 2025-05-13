import { ProWorkspaceInstance } from "@/contexts"
import { TWorkspaceResult } from "@/contexts/DevSpaceContext/workspaces/useWorkspace"
import { ManagementV1DevSpaceWorkspaceTemplate } from "@loft-enterprise/client/gen/models/managementV1DevSpaceWorkspaceTemplate"

export type TTabProps = Readonly<{
  host: string
  workspace: TWorkspaceResult<ProWorkspaceInstance>
  instance: ProWorkspaceInstance
  template: ManagementV1DevSpaceWorkspaceTemplate | undefined
}>
