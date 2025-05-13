import { exists, Result, Return } from "../../lib"
import {
  TWorkspace,
  TWorkspaceID,
  TWorkspaceStartConfig,
  TWorkspaceStatusResult,
  TWorkspaceWithoutStatus,
} from "../../types"
import { Command, isOk, serializeRawOptions, toFlagArg } from "../command"
import {
  DEVSPACE_COMMAND_DELETE,
  DEVSPACE_COMMAND_GET_WORKSPACE_CONFIG,
  DEVSPACE_COMMAND_GET_WORKSPACE_NAME,
  DEVSPACE_COMMAND_GET_WORKSPACE_UID,
  DEVSPACE_COMMAND_HELPER,
  DEVSPACE_COMMAND_LIST,
  DEVSPACE_COMMAND_STATUS,
  DEVSPACE_COMMAND_STOP,
  DEVSPACE_COMMAND_UP,
  DEVSPACE_COMMAND_TROUBLESHOOT,
  DEVSPACE_FLAG_DEBUG,
  DEVSPACE_FLAG_DEVCONTAINER_PATH,
  DEVSPACE_FLAG_FORCE,
  DEVSPACE_FLAG_ID,
  DEVSPACE_FLAG_IDE,
  DEVSPACE_FLAG_JSON_LOG_OUTPUT,
  DEVSPACE_FLAG_JSON_OUTPUT,
  DEVSPACE_FLAG_PREBUILD_REPOSITORY,
  DEVSPACE_FLAG_PROVIDER,
  DEVSPACE_FLAG_PROVIDER_OPTION,
  DEVSPACE_FLAG_RECREATE,
  DEVSPACE_FLAG_RESET,
  DEVSPACE_FLAG_SKIP_PRO,
  DEVSPACE_FLAG_SOURCE,
  DEVSPACE_FLAG_TIMEOUT,
  WORKSPACE_COMMAND_ADDITIONAL_FLAGS_KEY,
} from "../constants"

type TRawWorkspaces = readonly (Omit<TWorkspace, "status" | "id"> &
  Readonly<{ id: string | null }>)[]

export class WorkspaceCommands {
  static DEBUG = false
  static ADDITIONAL_FLAGS = new Map<string, string>()

  private static newCommand(args: string[]): Command {
    const extraFlags = []
    if (WorkspaceCommands.DEBUG) {
      extraFlags.push(DEVSPACE_FLAG_DEBUG)
    }

    return new Command([...args, ...extraFlags])
  }

  static async ListWorkspaces(skipPro: boolean): Promise<Result<TWorkspaceWithoutStatus[]>> {
    const maybeSkipProFlag = skipPro ? [DEVSPACE_FLAG_SKIP_PRO] : []

    const result = await new Command([
      DEVSPACE_COMMAND_LIST,
      DEVSPACE_FLAG_JSON_OUTPUT,
      ...maybeSkipProFlag,
    ]).run()
    if (result.err) {
      return result
    }

    const rawWorkspaces = JSON.parse(result.val.stdout) as TRawWorkspaces

    return Return.Value(
      rawWorkspaces.filter((workspace): workspace is TWorkspaceWithoutStatus =>
        exists(workspace.id)
      )
    )
  }

  static async FetchWorkspaceStatus(
    id: string
  ): Promise<Result<Pick<TWorkspace, "id" | "status">>> {
    const result = await new Command([DEVSPACE_COMMAND_STATUS, id, DEVSPACE_FLAG_JSON_OUTPUT]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return Return.Failed(`Failed to get status for workspace ${id}: ${result.val.stderr}`)
    }

    const { state } = JSON.parse(result.val.stdout) as TWorkspaceStatusResult

    return Return.Value({ id, status: state })
  }

  static async GetWorkspaceID(source: string) {
    const result = await new Command([
      DEVSPACE_COMMAND_HELPER,
      DEVSPACE_COMMAND_GET_WORKSPACE_NAME,
      source,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return Return.Failed(`Failed to get ID for workspace source ${source}: ${result.val.stderr}`)
    }

    return Return.Value(result.val.stdout)
  }

  static async GetWorkspaceUID() {
    const result = await new Command([
      DEVSPACE_COMMAND_HELPER,
      DEVSPACE_COMMAND_GET_WORKSPACE_UID,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return Return.Failed(`Failed to get UID: ${result.val.stderr}`)
    }

    return Return.Value(result.val.stdout)
  }

  static GetStatusLogs(id: string) {
    return new Command([DEVSPACE_COMMAND_STATUS, id, DEVSPACE_FLAG_JSON_LOG_OUTPUT])
  }

  static StartWorkspace(id: TWorkspaceID, config: TWorkspaceStartConfig) {
    const maybeSource = config.sourceConfig?.source
    const maybeIDFlag = exists(maybeSource) ? [toFlagArg(DEVSPACE_FLAG_ID, id)] : []

    const maybeSourceType = config.sourceConfig?.type
    const maybeSourceFlag =
      exists(maybeSourceType) && exists(maybeSource)
        ? [toFlagArg(DEVSPACE_FLAG_SOURCE, `${maybeSourceType}:${maybeSource}`)]
        : []
    const identifier = exists(maybeSource) && exists(maybeIDFlag) ? maybeSource : id

    const maybeIdeName = config.ideConfig?.name
    const maybeIDEFlag = exists(maybeIdeName) ? [toFlagArg(DEVSPACE_FLAG_IDE, maybeIdeName)] : []

    const maybeProviderID = config.providerConfig?.providerID
    const maybeProviderFlag = exists(maybeProviderID)
      ? [toFlagArg(DEVSPACE_FLAG_PROVIDER, maybeProviderID)]
      : []
    const maybeProviderOptions = config.providerConfig?.options
    const maybeProviderOptionsFlag = exists(maybeProviderOptions)
      ? serializeRawOptions(maybeProviderOptions, DEVSPACE_FLAG_PROVIDER_OPTION)
      : []

    const maybePrebuildRepositories = config.prebuildRepositories?.length
      ? [toFlagArg(DEVSPACE_FLAG_PREBUILD_REPOSITORY, config.prebuildRepositories.join(","))]
      : []

    const maybeDevcontainerPath = config.devcontainerPath
      ? [toFlagArg(DEVSPACE_FLAG_DEVCONTAINER_PATH, config.devcontainerPath)]
      : []

    const additionalFlags = []
    if (WorkspaceCommands.ADDITIONAL_FLAGS.size > 0) {
      for (const [key, value] of WorkspaceCommands.ADDITIONAL_FLAGS.entries()) {
        if (key === WORKSPACE_COMMAND_ADDITIONAL_FLAGS_KEY) {
          additionalFlags.push(value)
          continue
        }

        additionalFlags.push(toFlagArg(key, value))
      }
    }

    return WorkspaceCommands.newCommand([
      DEVSPACE_COMMAND_UP,
      identifier,
      ...maybeIDFlag,
      ...maybeSourceFlag,
      ...maybeIDEFlag,
      ...maybeProviderFlag,
      ...maybePrebuildRepositories,
      ...maybeDevcontainerPath,
      ...additionalFlags,
      ...maybeProviderOptionsFlag,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
    ])
  }

  static StopWorkspace(id: TWorkspaceID) {
    return WorkspaceCommands.newCommand([DEVSPACE_COMMAND_STOP, id, DEVSPACE_FLAG_JSON_LOG_OUTPUT])
  }

  static RebuildWorkspace(id: TWorkspaceID) {
    return WorkspaceCommands.newCommand([
      DEVSPACE_COMMAND_UP,
      id,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
      DEVSPACE_FLAG_RECREATE,
    ])
  }

  static ResetWorkspace(id: TWorkspaceID) {
    return WorkspaceCommands.newCommand([
      DEVSPACE_COMMAND_UP,
      id,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
      DEVSPACE_FLAG_RESET,
    ])
  }

  static TroubleshootWorkspace(id: TWorkspaceID) {
    return WorkspaceCommands.newCommand([DEVSPACE_COMMAND_TROUBLESHOOT, id])
  }

  static RemoveWorkspace(id: TWorkspaceID, force?: boolean) {
    const args = [DEVSPACE_COMMAND_DELETE, id, DEVSPACE_FLAG_JSON_LOG_OUTPUT]
    if (force) {
      args.push(DEVSPACE_FLAG_FORCE)
    }

    return WorkspaceCommands.newCommand(args)
  }

  static GetDevcontainerConfig(rawSource: string) {
    return new Command([
      DEVSPACE_COMMAND_HELPER,
      DEVSPACE_COMMAND_GET_WORKSPACE_CONFIG,
      rawSource,
      DEVSPACE_FLAG_TIMEOUT,
      "10s",
    ])
  }
}
