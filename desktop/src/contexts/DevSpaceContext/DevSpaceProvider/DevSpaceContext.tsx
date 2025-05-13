import { createContext } from "react"
import { TProviders, TQueryResult } from "../../../types"

export type TDevspaceContext = Readonly<{
  providers: TQueryResult<TProviders>
}>
export const DevSpaceContext = createContext<TDevspaceContext | null>(null)
