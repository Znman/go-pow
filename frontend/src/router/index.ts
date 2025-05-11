import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

// Importing components with proper TypeScript support
const TheWelcome = () => import('../components/TheWelcome.vue')
const Blockchain = () => import('../components/Blockchain.vue')
const Transactions = () => import('../components/Transactions.vue')
const Mining = () => import('../components/Mining.vue')

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'home',
    component: TheWelcome
  },
  {
    path: '/blockchain',
    name: 'blockchain',
    component: Blockchain
  },
  {
    path: '/transactions',
    name: 'transactions',
    component: Transactions
  },
  {
    path: '/mining',
    name: 'mining',
    component: Mining
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router