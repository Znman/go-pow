import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'home',
    component: () => import('../components/TheWelcome.vue')
  },
  {
    path: '/blockchain',
    name: 'blockchain',
    component: () => import('../components/Blockchain.vue')
  },
  {
    path: '/transactions',
    name: 'transactions',
    component: () => import('../components/Transactions.vue')
  },
  {
    path: '/explorer',
    name: 'explorer',
    component: () => import('../components/Explorer.vue')
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router