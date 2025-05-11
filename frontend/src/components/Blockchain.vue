<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
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
  Hash: string
}

interface BlockchainData {
  Chain: Block[]
  CurrentTransactions: Transaction[]
}

const blockchain = ref<BlockchainData | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const currentTime = ref(new Date().toISOString())
const refreshInterval = ref<number | null>(null)

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

// Auto refresh every 30 seconds
onMounted(() => {
  fetchChain()
  refreshInterval.value = window.setInterval(fetchChain, 30000)
})

// Clean up interval on component unmount
onUnmounted(() => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value)
  }
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
            {{ loading ? 'Loading...' : 'Refresh' }}
          </button>
        </div>
      </div>
    </header>

    <!-- Error Message -->
    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <!-- Loading Message -->
    <div v-if="loading && !blockchain" class="loading-message">
      Loading blockchain data...
    </div>

    <!-- Blockchain Data -->
    <div class="blockchain-data" v-if="blockchain">
      <div class="blocks-container">
        <div v-for="block in [...blockchain.Chain].reverse()" 
             :key="block.Index" 
             class="block-card">
          <div class="block-header">
            <h3>Block #{{ block.Index }}</h3>
            <span class="timestamp">{{ formatTimestamp(block.TimeStamp) }}</span>
          </div>
          <div class="block-details">
            <div class="detail-item">
              <span class="label">Hash:</span>
              <span class="value hash">{{ block.Hash }}</span>
            </div>
            <div class="detail-item">
              <span class="label">Previous Hash:</span>
              <span class="value hash">{{ block.PreviousHash }}</span>
            </div>
            <div class="detail-item">
              <span class="label">Proof:</span>
              <span class="value">{{ block.Proof }}</span>
            </div>
          </div>
          <div class="transactions">
            <h4>Transactions ({{ block.Transactions ? block.Transactions.length : 0 }})</h4>
            <div v-if="block.Transactions && block.Transactions.length > 0" 
                 class="transactions-list">
              <div v-for="(tx, index) in block.Transactions" 
                   :key="index" 
                   class="transaction">
                <span class="sender">{{ tx.Sender }}</span>
                <span class="arrow">→</span>
                <span class="recipient">{{ tx.Recipient }}</span>
                <span class="amount">{{ tx.Amount }} coins</span>
              </div>
            </div>
            <div v-else class="no-transactions">
              No transactions in this block
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.blockchain-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.main-header {
  background-color: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-section h1 {
  margin: 0;
  color: #2c3e50;
}

.subtitle {
  margin: 5px 0 0;
  color: #666;
}

.refresh-button {
  padding: 10px 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: background-color 0.3s;
}

.refresh-button:hover:not(:disabled) {
  background-color: #45a049;
}

.refresh-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.error-message {
  background-color: #ffe6e6;
  color: #dc3545;
  padding: 15px;
  border-radius: 4px;
  margin-bottom: 20px;
}

.loading-message {
  text-align: center;
  padding: 20px;
  color: #666;
}

.blocks-container {
  display: grid;
  gap: 20px;
}

.block-card {
  background-color: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.block-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.block-header h3 {
  margin: 0;
  color: #2c3e50;
}

.timestamp {
  color: #666;
  font-size: 0.9em;
}

.detail-item {
  margin-bottom: 10px;
}

.label {
  color: #666;
  font-weight: bold;
  margin-right: 8px;
}

.hash {
  font-family: monospace;
  word-break: break-all;
  background-color: #f8f9fa;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.9em;
}

.transactions {
  margin-top: 15px;
}

.transactions h4 {
  margin: 0 0 10px;
  color: #2c3e50;
}

.transactions-list {
  display: grid;
  gap: 10px;
}

.transaction {
  background-color: #f8f9fa;
  padding: 10px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.arrow {
  color: #666;
}

.amount {
  margin-left: auto;
  font-weight: bold;
  color: #4CAF50;
}

.no-transactions {
  color: #666;
  text-align: center;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
}
</style>