import { exists, getErrorFromChildProcess, Result, ResultError, Return } from "../../lib"
import {
  TAddProviderConfig,
  TCheckProviderUpdateResult,
  TProviderID,
  TProviderOptions,
  TProviders,
  TProviderSource,
} from "../../types"
import { Command, isOk, serializeRawOptions, toFlagArg } from "../command"
import {
  DEVSPACE_COMMAND_ADD,
  DEVSPACE_COMMAND_DELETE,
  DEVSPACE_COMMAND_GET_PROVIDER_NAME,
  DEVSPACE_COMMAND_LIST,
  DEVSPACE_COMMAND_OPTIONS,
  DEVSPACE_COMMAND_PROVIDER,
  DEVSPACE_COMMAND_SET_OPTIONS,
  DEVSPACE_COMMAND_UPDATE,
  DEVSPACE_COMMAND_USE,
  DEVSPACE_FLAG_DEBUG,
  DEVSPACE_FLAG_DRY,
  DEVSPACE_FLAG_JSON_LOG_OUTPUT,
  DEVSPACE_FLAG_JSON_OUTPUT,
  DEVSPACE_FLAG_NAME,
  DEVSPACE_FLAG_RECONFIGURE,
  DEVSPACE_FLAG_SINGLE_MACHINE,
  DEVSPACE_FLAG_USE,
} from "../constants"
import { DEVSPACE_COMMAND_CHECK_PROVIDER_UPDATE, DEVSPACE_COMMAND_HELPER } from "./../constants"

export class ProviderCommands {
  static DEBUG = false

  private static newCommand(args: string[]): Command {
    return new Command([...args, ...(ProviderCommands.DEBUG ? [DEVSPACE_FLAG_DEBUG] : [])])
  }

  static async ListProviders(): Promise<Result<TProviders>> {
    const result = await new Command([
      DEVSPACE_COMMAND_PROVIDER,
      DEVSPACE_COMMAND_LIST,
      DEVSPACE_FLAG_JSON_OUTPUT,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    const rawProviders = JSON.parse(result.val.stdout) as TProviders
    for (const provider of Object.values(rawProviders)) {
      provider.isProxyProvider =
        provider.config?.exec?.proxy !== undefined || provider.config?.exec?.daemon !== undefined
    }

    return Return.Value(rawProviders)
  }

  static async GetProviderID(source: string) {
    const result = await new Command([
      DEVSPACE_COMMAND_HELPER,
      DEVSPACE_COMMAND_GET_PROVIDER_NAME,
      source,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Value(result.val.stdout)
  }

  static async AddProvider(
    rawProviderSource: string,
    config: TAddProviderConfig
  ): Promise<ResultError> {
    const maybeName = config.name
    const maybeNameFlag = exists(maybeName) ? [toFlagArg(DEVSPACE_FLAG_NAME, maybeName)] : []
    const useFlag = toFlagArg(DEVSPACE_FLAG_USE, "false")

    const result = await ProviderCommands.newCommand([
      DEVSPACE_COMMAND_PROVIDER,
      DEVSPACE_COMMAND_ADD,
      rawProviderSource,
      ...maybeNameFlag,
      useFlag,
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

  static async RemoveProvider(id: TProviderID) {
    const result = await ProviderCommands.newCommand([
      DEVSPACE_COMMAND_PROVIDER,
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

  static async UseProvider(
    id: TProviderID,
    rawOptions?: Record<string, unknown>,
    reuseMachine?: boolean
  ) {
    const optionsFlag = rawOptions ? serializeRawOptions(rawOptions) : []
    const maybeResuseMachineFlag = reuseMachine ? [DEVSPACE_FLAG_SINGLE_MACHINE] : []

    const result = await ProviderCommands.newCommand([
      DEVSPACE_COMMAND_PROVIDER,
      DEVSPACE_COMMAND_USE,
      id,
      ...optionsFlag,
      ...maybeResuseMachineFlag,
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

  static async SetProviderOptions(
    id: TProviderID,
    rawOptions: Record<string, unknown>,
    reuseMachine: boolean,
    dry?: boolean,
    reconfigure?: boolean
  ) {
    const optionsFlag = serializeRawOptions(rawOptions)
    const maybeResuseMachineFlag = reuseMachine ? [DEVSPACE_FLAG_SINGLE_MACHINE] : []
    const maybeDry = dry ? [DEVSPACE_FLAG_DRY] : []
    const maybeReconfigure = reconfigure ? [DEVSPACE_FLAG_RECONFIGURE] : []

    const result = await ProviderCommands.newCommand([
      DEVSPACE_COMMAND_PROVIDER,
      DEVSPACE_COMMAND_SET_OPTIONS,
      id,
      ...optionsFlag,
      ...maybeResuseMachineFlag,
      ...maybeDry,
      ...maybeReconfigure,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    } else if (dry) {
      return Return.Value(JSON.parse(result.val.stdout) as TProviderOptions)
    }

    return Return.Ok()
  }

  static async GetProviderOptions(id: TProviderID) {
    const result = await new Command([
      DEVSPACE_COMMAND_PROVIDER,
      DEVSPACE_COMMAND_OPTIONS,
      id,
      DEVSPACE_FLAG_JSON_OUTPUT,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Value(JSON.parse(result.val.stdout) as TProviderOptions)
  }

  static async CheckProviderUpdate(id: TProviderID) {
    const result = await new Command([
      DEVSPACE_COMMAND_HELPER,
      DEVSPACE_COMMAND_CHECK_PROVIDER_UPDATE,
      id,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Value(JSON.parse(result.val.stdout) as TCheckProviderUpdateResult)
  }

  static async UpdateProvider(id: TProviderID, source: TProviderSource) {
    const useFlag = toFlagArg(DEVSPACE_FLAG_USE, "false")

    const result = await new Command([
      DEVSPACE_COMMAND_PROVIDER,
      DEVSPACE_COMMAND_UPDATE,
      id,
      source.raw ?? source.github ?? source.url ?? source.file ?? "",
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
      useFlag,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Ok()
  }
}
