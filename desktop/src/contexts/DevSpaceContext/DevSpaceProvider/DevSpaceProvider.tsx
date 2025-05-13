import { useQuery } from "@tanstack/react-query"
import { ReactNode, useMemo } from "react"
import { client } from "../../../client"
import { QueryKeys } from "../../../queryKeys"
import { REFETCH_PROVIDER_INTERVAL_MS } from "../constants"
import { usePollWorkspaces } from "../workspaces"
import { DevSpaceContext, TDevspaceContext } from "./DevSpaceContext"

export function DevSpaceProvider({ children }: Readonly<{ children?: ReactNode }>) {
  usePollWorkspaces()

  const providersQuery = useQuery({
    queryKey: QueryKeys.PROVIDERS,
    queryFn: async () => (await client.providers.listAll()).unwrap(),
    refetchInterval: REFETCH_PROVIDER_INTERVAL_MS,
    enabled: true,
  })

  const value = useMemo<TDevspaceContext>(
    () => ({
      providers: [
        providersQuery.data,
        { status: providersQuery.status, error: providersQuery.error },
      ],
    }),
    [providersQuery.data, providersQuery.status, providersQuery.error]
  )

  return <DevSpaceContext.Provider value={value}>{children}</DevSpaceContext.Provider>
}

export function ProviderProvider({ children }: Readonly<{ children?: ReactNode }>) {
  const providersQuery = useQuery({
    queryKey: QueryKeys.PROVIDERS,
    queryFn: async () => (await client.providers.listAll()).unwrap(),
    refetchInterval: REFETCH_PROVIDER_INTERVAL_MS,
    enabled: true,
  })

  const value = useMemo<TDevspaceContext>(
    () => ({
      providers: [
        providersQuery.data,
        { status: providersQuery.status, error: providersQuery.error },
      ],
    }),
    [providersQuery.data, providersQuery.status, providersQuery.error]
  )

  return <DevSpaceContext.Provider value={value}>{children}</DevSpaceContext.Provider>
}
