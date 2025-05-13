import { Result, ResultError, Return, getErrorFromChildProcess } from "@/lib"
import {
  TImportWorkspaceConfig,
  TListProInstancesConfig,
  TPlatformHealthCheck,
  TProID,
  TProInstance,
  TPlatformVersionInfo,
  TPlatformUpdateCheck,
} from "@/types"
import { Command, isOk, serializeRawOptions, toFlagArg } from "../command"
import {
  DEVSPACE_COMMAND_DELETE,
  DEVSPACE_COMMAND_IMPORT_WORKSPACE,
  DEVSPACE_COMMAND_LIST,
  DEVSPACE_COMMAND_LOGIN,
  DEVSPACE_COMMAND_PRO,
  DEVSPACE_FLAG_ACCESS_KEY,
  DEVSPACE_FLAG_DEBUG,
  DEVSPACE_FLAG_FORCE_BROWSER,
  DEVSPACE_FLAG_HOST,
  DEVSPACE_FLAG_INSTANCE,
  DEVSPACE_FLAG_JSON_LOG_OUTPUT,
  DEVSPACE_FLAG_JSON_OUTPUT,
  DEVSPACE_FLAG_LOGIN,
  DEVSPACE_FLAG_PROJECT,
  DEVSPACE_FLAG_USE,
  DEVSPACE_FLAG_WORKSPACE_ID,
  DEVSPACE_FLAG_WORKSPACE_PROJECT,
  DEVSPACE_FLAG_WORKSPACE_UID,
} from "../constants"
import { TStreamEventListenerFn } from "../types"
import { ManagementV1DevSpaceWorkspaceInstance } from "@loft-enterprise/client/gen/models/managementV1DevSpaceWorkspaceInstance"
import { ManagementV1Project } from "@loft-enterprise/client/gen/models/managementV1Project"
import { ManagementV1Self } from "@loft-enterprise/client/gen/models/managementV1Self"
import { ManagementV1ProjectTemplates } from "@loft-enterprise/client/gen/models/managementV1ProjectTemplates"
import { ManagementV1ProjectClusters } from "@loft-enterprise/client/gen/models/managementV1ProjectClusters"

export class ProCommands {
  static DEBUG = false

  private static newCommand(args: string[]): Command {
    return new Command([...args, ...(ProCommands.DEBUG ? [DEVSPACE_FLAG_DEBUG] : [])])
  }

  static async Login(
    host: string,
    accessKey?: string,
    listener?: TStreamEventListenerFn
  ): Promise<ResultError> {
    const maybeAccessKeyFlag = accessKey ? [toFlagArg(DEVSPACE_FLAG_ACCESS_KEY, accessKey)] : []
    const useFlag = toFlagArg(DEVSPACE_FLAG_USE, "false")

    const cmd = ProCommands.newCommand([
      DEVSPACE_COMMAND_PRO,
      DEVSPACE_COMMAND_LOGIN,
      host,
      useFlag,
      DEVSPACE_FLAG_FORCE_BROWSER,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
      ...maybeAccessKeyFlag,
    ])
    if (listener) {
      return cmd.stream(listener)
    } else {
      const result = await cmd.run()
      if (result.err) {
        return result
      }

      if (!isOk(result.val)) {
        return getErrorFromChildProcess(result.val)
      }

      return Return.Ok()
    }
  }

  static async ListProInstances(
    config?: TListProInstancesConfig
  ): Promise<Result<readonly TProInstance[]>> {
    const maybeLoginFlag = config?.authenticate ? [DEVSPACE_FLAG_LOGIN] : []
    const result = await ProCommands.newCommand([
      DEVSPACE_COMMAND_PRO,
      DEVSPACE_COMMAND_LIST,
      DEVSPACE_FLAG_JSON_OUTPUT,
      ...maybeLoginFlag,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    const instances = JSON.parse(result.val.stdout) as readonly TProInstance[]

    return Return.Value(instances)
  }

  static async RemoveProInstance(id: TProID) {
    const result = await ProCommands.newCommand([
      DEVSPACE_COMMAND_PRO,
      DEVSPACE_COMMAND_DELETE,
      id,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Ok()
  }

  static async ImportWorkspace(config: TImportWorkspaceConfig): Promise<ResultError> {
    const optionsFlag = config.options ? serializeRawOptions(config.options) : []
    const result = await new Command([
      DEVSPACE_COMMAND_PRO,
      DEVSPACE_COMMAND_IMPORT_WORKSPACE,
      config.devSpaceProHost,
      DEVSPACE_FLAG_WORKSPACE_ID,
      config.workspaceID,
      DEVSPACE_FLAG_WORKSPACE_UID,
      config.workspaceUID,
      DEVSPACE_FLAG_WORKSPACE_PROJECT,
      config.project,
      ...optionsFlag,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Ok()
  }

  static WatchWorkspaces(id: TProID, projectName: string) {
    const hostFlag = toFlagArg(DEVSPACE_FLAG_HOST, id)
    const projectFlag = toFlagArg(DEVSPACE_FLAG_PROJECT, projectName)
    const args = [DEVSPACE_COMMAND_PRO, "watch-workspaces", hostFlag, projectFlag]

    return ProCommands.newCommand(args)
  }

  static async ListProjects(id: TProID) {
    const hostFlag = toFlagArg(DEVSPACE_FLAG_HOST, id)
    const args = [DEVSPACE_COMMAND_PRO, "list-projects", hostFlag]

    const result = await ProCommands.newCommand(args).run()
    if (result.err) {
      return result
    }
    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Value(JSON.parse(result.val.stdout) as readonly ManagementV1Project[])
  }

  static async GetSelf(id: TProID) {
    const hostFlag = toFlagArg(DEVSPACE_FLAG_HOST, id)
    const args = [DEVSPACE_COMMAND_PRO, "self", hostFlag]

    const result = await ProCommands.newCommand(args).run()
    if (result.err) {
      return result
    }
    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Value(JSON.parse(result.val.stdout) as ManagementV1Self)
  }

  static async ListTemplates(id: TProID, projectName: string) {
    const hostFlag = toFlagArg(DEVSPACE_FLAG_HOST, id)
    const projectFlag = toFlagArg(DEVSPACE_FLAG_PROJECT, projectName)
    const args = [DEVSPACE_COMMAND_PRO, "list-templates", hostFlag, projectFlag]

    const result = await ProCommands.newCommand(args).run()
    if (result.err) {
      return result
    }
    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Value(JSON.parse(result.val.stdout) as ManagementV1ProjectTemplates)
  }

  static async ListClusters(id: TProID, projectName: string) {
    const hostFlag = toFlagArg(DEVSPACE_FLAG_HOST, id)
    const projectFlag = toFlagArg(DEVSPACE_FLAG_PROJECT, projectName)
    const args = [DEVSPACE_COMMAND_PRO, "list-clusters", hostFlag, projectFlag]

    const result = await ProCommands.newCommand(args).run()
    if (result.err) {
      return result
    }
    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Value(JSON.parse(result.val.stdout) as ManagementV1ProjectClusters)
  }

  static async CreateWorkspace(id: TProID, instance: ManagementV1DevSpaceWorkspaceInstance) {
    const hostFlag = toFlagArg(DEVSPACE_FLAG_HOST, id)
    const instanceFlag = toFlagArg(DEVSPACE_FLAG_INSTANCE, JSON.stringify(instance))
    const args = [DEVSPACE_COMMAND_PRO, "create-workspace", hostFlag, instanceFlag]

    const result = await ProCommands.newCommand(args).run()
    if (result.err) {
      return result
    }
    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Value(JSON.parse(result.val.stdout) as ManagementV1DevSpaceWorkspaceInstance)
  }

  static async UpdateWorkspace(id: TProID, instance: ManagementV1DevSpaceWorkspaceInstance) {
    const hostFlag = toFlagArg(DEVSPACE_FLAG_HOST, id)
    const instanceFlag = toFlagArg(DEVSPACE_FLAG_INSTANCE, JSON.stringify(instance))
    const args = [DEVSPACE_COMMAND_PRO, "update-workspace", hostFlag, instanceFlag]

    const result = await ProCommands.newCommand(args).run()
    if (result.err) {
      return result
    }
    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Value(JSON.parse(result.val.stdout) as ManagementV1DevSpaceWorkspaceInstance)
  }

  static async CheckHealth(id: TProID) {
    const hostFlag = toFlagArg(DEVSPACE_FLAG_HOST, id)
    const args = [DEVSPACE_COMMAND_PRO, "check-health", hostFlag]

    const result = await ProCommands.newCommand(args).run()
    if (result.err) {
      return result
    }
    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Value(JSON.parse(result.val.stdout) as TPlatformHealthCheck)
  }

  static async GetVersion(id: TProID) {
    const hostFlag = toFlagArg(DEVSPACE_FLAG_HOST, id)
    const args = [DEVSPACE_COMMAND_PRO, "version", hostFlag]

    const result = await ProCommands.newCommand(args).run()
    if (result.err) {
      return result
    }
    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Value(JSON.parse(result.val.stdout) as TPlatformVersionInfo)
  }

  static async CheckUpdate(id: TProID) {
    const hostFlag = toFlagArg(DEVSPACE_FLAG_HOST, id)
    const args = [DEVSPACE_COMMAND_PRO, "check-update", hostFlag]

    const result = await ProCommands.newCommand(args).run()
    if (result.err) {
      return result
    }
    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Value(JSON.parse(result.val.stdout) as TPlatformUpdateCheck)
  }

  static async Update(id: TProID, version: string) {
    const hostFlag = toFlagArg(DEVSPACE_FLAG_HOST, id)
    const args = [DEVSPACE_COMMAND_PRO, "update-provider", version, hostFlag]

    const result = await ProCommands.newCommand(args).run()
    if (result.err) {
      return result
    }
    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Value(JSON.parse(result.val.stdout) as TPlatformUpdateCheck)
  }
}
