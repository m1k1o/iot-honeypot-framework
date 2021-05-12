import Vue from 'vue'
import Vuex, { ActionContext } from 'vuex';

import {
  Configuration,
  ServiceId,
  ServiceSpec,
  ServicesApi,
  ProxyId,
  ProxySpec,
  ProxiesApi,
  NodeSpec,
  NodesApi,
  NodeJoin,
} from '@/api/index.ts'

import { state, State } from './state'

Vue.use(Vuex)

const apiConfig = new Configuration({ basePath: (location.protocol + '//' + location.host + location.pathname).replace(/\/+$/, '') + '/api' })

export default new Vuex.Store({
  state,
  mutations: {
    SERVICES_SET(state: State, services: ServiceSpec[]) {
      Vue.set(state, 'services', services)
    },
    SERVICES_ADD(state: State, service: ServiceSpec) {
      Vue.set(state, 'services', [...state.services, service])
    },
    SERVICES_DEL(state: State, serviceId: string) {
      const services = state.services.filter(({ id }) => id != serviceId)
      Vue.set(state, 'services', services)
    },

    PROXIES_SET(state: State, proxies: ProxySpec[]) {
      Vue.set(state, 'proxies', proxies)
    },
    PROXIES_ADD(state: State, proxy: ProxySpec) {
      Vue.set(state, 'proxies', [...state.proxies, proxy])
    },
    PROXIES_DEL(state: State, proxyId: string) {
      const proxies = state.proxies.filter(({ id }) => id != proxyId)
      Vue.set(state, 'proxies', proxies)
    },

    NODES_SET(state: State, nodes: NodeSpec[]) {
      Vue.set(state, 'nodes', nodes)
    },
    NODES_DEL(state: State, nodeId: string) {
      const nodes = state.nodes.filter(({ id }) => id != nodeId)
      Vue.set(state, 'nodes', nodes)
    },
  },
  getters: {
    services(state: State) {
      return state.services;
    },
    proxies(state: State) {
      return state.proxies;
    },
    nodes(state: State) {
      return state.nodes;
    },
  },
  actions: {
    async SERVICES_LOAD({ commit }: ActionContext<State, State>) {
      const api = new ServicesApi(apiConfig)
      const services: ServiceSpec[] = await api.servicesList()
      commit('SERVICES_SET', services);
    },
    async SERVICE_CREATE({ commit }: ActionContext<State, State>, service: ServiceSpec) {
      const api = new ServicesApi(apiConfig)
      const { id }: ServiceId = await api.serviceCreate(service)
      commit('SERVICES_ADD', { id, ...service, status: { running: 0, desired: service.replicas } });
    },
    async SERVICE_REMOVE({ commit }: ActionContext<State, State>, serviceId: string) {
      const api = new ServicesApi(apiConfig)
      await api.serviceRemove(serviceId)
      commit('SERVICES_DEL', serviceId);
    },

    async PROXIES_LOAD({ commit }: ActionContext<State, State>) {
      const api = new ProxiesApi(apiConfig)
      const proxies: ProxySpec[] = await api.proxiesList()
      commit('PROXIES_SET', proxies);
    },
    async PROXY_CREATE({ commit }: ActionContext<State, State>, proxy: ProxySpec) {
      const api = new ProxiesApi(apiConfig)
      const { id }: ProxyId = await api.proxyCreate(proxy)
      commit('PROXIES_ADD', { id, ...proxy });
    },
    async PROXY_REMOVE({ commit }: ActionContext<State, State>, proxyId: string) {
      const api = new ProxiesApi(apiConfig)
      await api.proxyRemove(proxyId)
      commit('PROXIES_DEL', proxyId);
    },

    async NODES_LOAD({ commit }: ActionContext<State, State>) {
      const api = new NodesApi(apiConfig)
      const nodes: ProxySpec[] = await api.nodesList()
      commit('NODES_SET', nodes);
    },
    async NODE_JOIN(): Promise<NodeJoin> {
      const api = new NodesApi(apiConfig)
      return await api.nodeJoin()
    },
    async NODE_REMOVE({ commit }: ActionContext<State, State>, nodeId: string) {
      const api = new NodesApi(apiConfig)
      await api.nodeRemove(nodeId)
      commit('NODES_DEL', nodeId);
    },
  },
  modules: {
  }
})
