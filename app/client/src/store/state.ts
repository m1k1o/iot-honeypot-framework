import {
  NodeSpec,
  ServiceSpec,
  ProxySpec,
} from '@/api/index.ts'

export const state = {
  services: [] as ServiceSpec[],
  proxies: [] as ProxySpec[],
  nodes: [] as NodeSpec[],
}

export type State = typeof state
