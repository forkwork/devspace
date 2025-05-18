import { useProContext } from "@/contexts"
import { useQuery, UseQueryResult } from "@tanstack/react-query"
import { QueryKeys } from "@/queryKeys"
import { ManagementV1DevSpaceWorkspaceTemplate } from "@loft-enterprise/client/gen/models/managementV1DevSpaceWorkspaceTemplate"
import { ManagementV1DevSpaceEnvironmentTemplate } from "@loft-enterprise/client/gen/models/managementV1DevSpaceEnvironmentTemplate"
import { ManagementV1DevSpaceWorkspacePreset } from "@loft-enterprise/client/gen/models/managementV1DevSpaceWorkspacePreset"

type TTemplates = Readonly<{
  default: ManagementV1DevSpaceWorkspaceTemplate | undefined
  workspace: readonly ManagementV1DevSpaceWorkspaceTemplate[]
  environment: readonly ManagementV1DevSpaceEnvironmentTemplate[]
  presets: readonly ManagementV1DevSpaceWorkspacePreset[]
}>
export function useTemplates(): UseQueryResult<TTemplates> {
  const { host, currentProject, client } = useProContext()
  const query = useQuery<TTemplates>({
    queryKey: QueryKeys.proWorkspaceTemplates(host, currentProject?.metadata!.name!),
    queryFn: async () => {
      const projectTemplates = (
        await client.getProjectTemplates(currentProject?.metadata!.name!)
      ).unwrap()

      // try to find default template in list
      let defaultTemplate: ManagementV1DevSpaceWorkspaceTemplate | undefined = undefined
      if (projectTemplates?.defaultDevSpaceWorkspaceTemplate) {
        defaultTemplate = projectTemplates.devSpaceWorkspaceTemplates?.find(
          (template) =>
            template.metadata?.name === projectTemplates.defaultDevSpaceWorkspaceTemplate
        )
      }

      return {
        default: defaultTemplate,
        workspace: projectTemplates?.devSpaceWorkspaceTemplates ?? [],
        environment: projectTemplates?.devSpaceEnvironmentTemplates ?? [],
        presets: projectTemplates?.devSpaceWorkspacePresets ?? [],
      }
    },
    enabled: !!currentProject,
  })

  return query
}
