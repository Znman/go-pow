import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

// Importing components with proper TypeScript support
const TheWelcome = () => import('../components/TheWelcome.vue')
const HelloWorld = () => import('../components/HelloWorld.vue')
const Blockchain = () => import('../components/Blockchain.vue')
const Transactions = () => import('../components/Transactions.vue')

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'home',
    component: TheWelcome
  },
  {
    path: '/hello',
    name: 'hello',
    component: HelloWorld
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
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router