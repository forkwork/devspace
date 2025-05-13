import { Command, isOk } from "../command"
import {
  DEVSPACE_COMMAND_IDE,
  DEVSPACE_COMMAND_LIST,
  DEVSPACE_COMMAND_USE,
  DEVSPACE_FLAG_DEBUG,
  DEVSPACE_FLAG_JSON_LOG_OUTPUT,
  DEVSPACE_FLAG_JSON_OUTPUT,
} from "../constants"
import { getErrorFromChildProcess, Result, ResultError, Return } from "@/lib"
import { TIDEs } from "@/types"

export class IDECommands {
  static DEBUG = false

  private static newCommand(args: string[]): Command {
    return new Command([...args, ...(IDECommands.DEBUG ? [DEVSPACE_FLAG_DEBUG] : [])])
  }

  static async UseIDE(ide: string): Promise<ResultError> {
    const result = await IDECommands.newCommand([
      DEVSPACE_COMMAND_IDE,
      DEVSPACE_COMMAND_USE,
      ide,
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

  static async ListIDEs(): Promise<Result<TIDEs>> {
    const result = await IDECommands.newCommand([
      DEVSPACE_COMMAND_IDE,
      DEVSPACE_COMMAND_LIST,
      DEVSPACE_FLAG_JSON_OUTPUT,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    const ides = JSON.parse(result.val.stdout) as TIDEs

    return Return.Value(ides)
  }
}
