<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import axios from 'axios'

interface Transaction {
  Sender: string
  Recipient: string
  Amount: number
}

interface Block {
  index: number
  timestamp: number
  transactions: Transaction[]
  proof: number
  previousHash: string
  hash: string
}

interface BlockchainData {
  chain: Block[]
  currentTransactions: Transaction[]
}

const blockchain = ref<BlockchainData | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const miningInProgress = ref(false)

// Fetch blockchain data
const fetchChain = async () => {
  loading.value = true
  error.value = null

  try {
    const response = await axios.get<BlockchainData>('http://localhost:8000/chain')
    blockchain.value = response.data
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to load blockchain'
    console.error('Error loading blockchain:', err)
  } finally {
    loading.value = false
  }
}

// Mine new block
const mineBlock = async () => {
  if (miningInProgress.value) return

  miningInProgress.value = true
  error.value = null

  try {
    await axios.post('http://localhost:8000/mine')
    await fetchChain() // Refresh the chain after mining
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Mining failed'
    console.error('Error mining block:', err)
  } finally {
    miningInProgress.value = false
  }
}

// Format timestamp
const formatTimestamp = (timestamp: number): string => {
  return new Date(timestamp * 1000).toLocaleString()
}

// Auto refresh setup
const refreshInterval = ref<number | null>(null)

onMounted(() => {
  fetchChain()
  refreshInterval.value = window.setInterval(fetchChain, 30000)
})

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
          <h1>Blockchain</h1>
        </div>
        <div class="actions">
          <button @click="mineBlock" :disabled="miningInProgress" class="mine-button">
            <span class="button-icon">⛏️</span>
            {{ miningInProgress ? 'Mining...' : 'Mine New Block' }}
          </button>
          <button @click="fetchChain" :disabled="loading" class="refresh-button">
            <span class="button-icon">🔄</span>
            Refresh
          </button>
        </div>
      </div>
    </header>

    <!-- Error Message -->
    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <!-- Blockchain Data -->
    <div class="blockchain-data" v-if="blockchain">
      <div class="blocks-container">
        <div v-for="block in [...blockchain.chain].reverse()" :key="block.index" class="block-card">
          <div class="block-header">
            <h3>Block #{{ block.index }}</h3>
            <span class="timestamp">{{ formatTimestamp(block.timestamp) }}</span>
          </div>
          <div class="block-details">
            <div class="detail-item">
              <span class="label">Hash:</span>
              <span class="value hash">{{ block.hash }}</span>
            </div>
            <div class="detail-item">
              <span class="label">Previous Hash:</span>
              <span class="value hash">{{ block.previousHash }}</span>
            </div>
            <div class="detail-item">
              <span class="label">Proof:</span>
              <span class="value">{{ block.proof }}</span>
            </div>
          </div>
          <div class="transactions">
            <h4>Transactions</h4>
            <div v-if="block.transactions && block.transactions.length > 0" class="transactions-list">
              <div v-for="(tx, index) in block.transactions" :key="index" class="transaction">
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

.actions {
  display: flex;
  gap: 10px;
}

.mine-button,
.refresh-button {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  transition: background-color 0.3s;
}

.mine-button {
  background-color: #4CAF50;
  color: white;
}

.mine-button:hover:not(:disabled) {
  background-color: #45a049;
}

.refresh-button {
  background-color: #2196F3;
  color: white;
}

.refresh-button:hover:not(:disabled) {
  background-color: #1e88e5;
}

.button-icon {
  font-size: 16px;
}

.block-card {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
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

.block-details {
  margin-bottom: 15px;
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
}

.transactions {
  border-top: 1px solid #eee;
  padding-top: 15px;
}

.transaction {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  background-color: #f8f9fa;
  border-radius: 4px;
  margin-bottom: 8px;
}

.arrow {
  color: #666;
}

.amount {
  margin-left: auto;
  font-weight: bold;
  color: #4CAF50;
}

.error-message {
  background-color: #ffe6e6;
  color: #dc3545;
  padding: 15px;
  border-radius: 4px;
  margin-bottom: 20px;
}

.no-transactions {
  color: #666;
  text-align: center;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
}

button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}
</style>