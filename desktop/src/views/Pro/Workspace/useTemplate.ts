import { useTemplates } from "@/contexts"
import { TParameterWithValue, getParametersWithValues } from "@/lib"
import { ManagementV1DevSpaceWorkspaceInstance } from "@loft-enterprise/client/gen/models/managementV1DevSpaceWorkspaceInstance"
import { ManagementV1DevSpaceWorkspaceTemplate } from "@loft-enterprise/client/gen/models/managementV1DevSpaceWorkspaceTemplate"
import { useMemo } from "react"

export function useTemplate(instance: ManagementV1DevSpaceWorkspaceInstance | undefined) {
  const { data: templates } = useTemplates()

  return useMemo<{
    parameters: readonly TParameterWithValue[]
    template: ManagementV1DevSpaceWorkspaceTemplate | undefined
  }>(() => {
    // find template for workspace
    const currentTemplate = templates?.workspace.find(
      (template) => instance?.spec?.templateRef?.name === template.metadata?.name
    )
    const empty = { parameters: [], template: undefined }
    if (!currentTemplate || !instance) {
      return empty
    }

    const parameters = getParametersWithValues(instance, currentTemplate)
    if (!parameters) {
      return { parameters: [], template: currentTemplate }
    }

    return { parameters, template: currentTemplate }
  }, [instance, templates])
}
