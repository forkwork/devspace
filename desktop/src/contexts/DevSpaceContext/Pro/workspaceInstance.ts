import { TIDE, TIdentifiable, TWorkspaceSource } from "@/types"
import { ManagementV1DevSpaceWorkspaceInstance } from "@loft-enterprise/client/gen/models/managementV1DevSpaceWorkspaceInstance"
import { Labels, deepCopy } from "@/lib"
import { Resources } from "@loft-enterprise/client"
import { ManagementV1DevSpaceWorkspaceInstanceStatus } from "@loft-enterprise/client/gen/models/managementV1DevSpaceWorkspaceInstanceStatus"

export class ProWorkspaceInstance
  extends ManagementV1DevSpaceWorkspaceInstance
  implements TIdentifiable
{
  public readonly status: ProWorkspaceInstanceStatus | undefined

  public get id(): string {
    const maybeID = this.metadata?.labels?.[Labels.WorkspaceID]
    if (!maybeID) {
      // If we don't have an ID we should ignore the instance.
      // Throwing an error for now to see how often this happens
      throw new Error(`No Workspace ID label present on instance ${this.metadata?.name}`)
    }

    return maybeID
  }

  constructor(instance: ManagementV1DevSpaceWorkspaceInstance) {
    super()

    this.apiVersion = `${Resources.ManagementV1DevSpaceWorkspaceInstance.group}/${Resources.ManagementV1DevSpaceWorkspaceInstance.version}`
    this.kind = Resources.ManagementV1DevSpaceWorkspaceInstance.kind
    this.metadata = deepCopy(instance.metadata)
    this.spec = deepCopy(instance.spec)
    this.status = deepCopy(instance.status) as ProWorkspaceInstanceStatus
  }
}

class ProWorkspaceInstanceStatus extends ManagementV1DevSpaceWorkspaceInstanceStatus {
  "source"?: TWorkspaceSource
  "ide"?: TIDE
  "metrics"?: ProWorkspaceMetricsSummary

  constructor() {
    super()
  }
}

class ProWorkspaceMetricsSummary {
  "latencyMs"?: number
  "connectionType"?: "direct" | "DERP"
  "derpRegion"?: string
}
