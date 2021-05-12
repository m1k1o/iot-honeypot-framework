import Vue from 'vue'
import VueRouter, { RouteConfig } from 'vue-router'

import Dashboard from '../views/Dashboard.vue'
import Services from '../views/Services.vue'
import Proxies from '../views/Proxies.vue'
import Nodes from '../views/Nodes.vue'
import NodeJoin from '../views/NodeJoin.vue'

// Mocks
import Settings from '../views/mocks/Settings.vue'

Vue.use(VueRouter)

const routes: Array<RouteConfig> = [
  {
    path: '/',
    name: 'dashboard',
    meta: { title: 'Dashboard' },
    component: Dashboard
  },
  {
    path: '/services',
    name: 'services',
    meta: { title: 'Services' },
    component: Services
  },
  {
    path: '/proxies',
    name: 'proxies',
    meta: { title: 'Proxies' },
    component: Proxies
  },
  {
    path: '/nodes',
    name: 'nodes',
    meta: { title: 'Nodes' },
    component: Nodes
  },
  {
    path: '/nodes/join',
    name: 'node-join',
    meta: { title: 'Joien new Node' },
    component: NodeJoin
  },

  // --

  {
    path: '/settings',
    name: 'settings',
    meta: { title: 'Settings' },
    component: Settings
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
