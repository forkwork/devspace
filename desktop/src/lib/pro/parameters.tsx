import { ManagementV1DevSpaceWorkspaceInstance } from "@loft-enterprise/client/gen/models/managementV1DevSpaceWorkspaceInstance"
import { ManagementV1DevSpaceWorkspaceTemplate } from "@loft-enterprise/client/gen/models/managementV1DevSpaceWorkspaceTemplate"
import { StorageV1AppParameter } from "@loft-enterprise/client/gen/models/storageV1AppParameter"
import { StorageV1DevSpaceWorkspaceTemplateVersion } from "@loft-enterprise/client/gen/models/storageV1DevSpaceWorkspaceTemplateVersion"
import { compareVersions } from "compare-versions"
import jsyaml from "js-yaml"

export type TParameterWithValue = StorageV1AppParameter & { value?: string | number | boolean }

export function getParametersWithValues(
  instance: ManagementV1DevSpaceWorkspaceInstance,
  template: ManagementV1DevSpaceWorkspaceTemplate
): readonly TParameterWithValue[] | undefined {
  let rawParameters: StorageV1AppParameter[] | undefined = template.spec?.parameters
  if (instance.spec?.templateRef?.version) {
    // find versioned parameters
    rawParameters = template.spec?.versions?.find(
      (version) => version.version === instance.spec?.templateRef?.version
    )?.parameters
  } else if (template.spec?.versions && template.spec.versions.length > 0) {
    // fall back to latest version
    rawParameters = template.spec.versions[0]?.parameters
  }

  if (!instance.spec?.parameters || !rawParameters) {
    return undefined
  }

  try {
    const out = jsyaml.load(instance.spec.parameters) as Record<string, string | number | boolean>

    return rawParameters.map((param) => {
      const path = param.variable
      if (path) {
        return { ...param, value: out[path] }
      }

      return param
    })
  } catch {
    return undefined
  }
}

export function getParameters(
  template: ManagementV1DevSpaceWorkspaceTemplate | undefined,
  selectedVersion: string | undefined
): readonly StorageV1AppParameter[] | undefined {
  if (!template?.spec) {
    return undefined
  }

  if (selectedVersion) {
    return template.spec.versions?.find((version) => version.version === selectedVersion)
      ?.parameters
  }

  if (template.spec.versions && template.spec.versions.length > 0) {
    const latestVersion = findLatestVersion(template.spec.versions)
    if (latestVersion) {
      return template.spec.versions.find((version) => version.version === latestVersion.version)
        ?.parameters
    }
  }

  return template.spec.parameters
}

function findLatestVersion(
  versions: readonly StorageV1DevSpaceWorkspaceTemplateVersion[]
): StorageV1DevSpaceWorkspaceTemplateVersion | undefined {
  return versions.slice().sort(sortByVersionDesc)[0]
}

export function sortByVersionDesc(
  a: StorageV1DevSpaceWorkspaceTemplateVersion,
  b: StorageV1DevSpaceWorkspaceTemplateVersion
): number {
  return compareVersions(stripVersionPrefix(b.version ?? ""), stripVersionPrefix(a.version ?? ""))
}

function stripVersionPrefix(version: string): string {
  if (version.startsWith("v")) {
    return version.substring(1, version.length)
  }

  return version
}
