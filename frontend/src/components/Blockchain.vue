<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'

interface Transaction {
  Sender: string
  Recipient: string
  Amount: number
}

interface Block {
  Index: number
  TimeStamp: number
  Transactions: Transaction[] | null
  Proof: number
  PreviousHash: string
}

interface BlockchainData {
  Chain: Block[]
  CurrentTransactions: Transaction[]
}

const blockchain = ref<BlockchainData | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const currentTime = ref(new Date().toISOString())

const fetchChain = async () => {
  loading.value = true
  error.value = null

  try {
    const response = await axios.get<BlockchainData>('http://localhost:8000/chain')
    blockchain.value = response.data
    currentTime.value = new Date().toISOString()
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to load blockchain'
    console.error('Error loading blockchain:', err)
  } finally {
    loading.value = false
  }
}

const formatTimestamp = (timestamp: number): string => {
  return new Date(timestamp * 1000).toLocaleString()
}

onMounted(() => {
  fetchChain()
})
</script>

<template>
  <div class="blockchain-container">
    <header class="main-header">
      <div class="header-content">
        <div class="title-section">
          <h1>Blockchain Explorer</h1>
          <p class="subtitle">Last updated: {{ new Date(currentTime).toLocaleString() }}</p>
        </div>
        <div class="actions">
          <button @click="fetchChain" :disabled="loading" class="refresh-button">
            <span class="button-icon" v-if="loading">↻</span>
            <span v-else>⟳</span>
            {{ loading ? 'Syncing...' : 'Sync Chain' }}
          </button>
        </div>
      </div>
    </header>
  </div>
</template>