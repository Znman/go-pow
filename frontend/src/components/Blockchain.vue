<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import axios from 'axios'
import BlockchainGraph from './BlockchainGraph.vue'


interface Transaction {
  sender: string    // Changed from Sender to match backend JSON
  recipient: string // Changed from Recipient to match backend JSON
  amount: number    // Changed from Amount to match backend JSON
  message: string
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
const currentTime = ref(new Date().toISOString())

// Fetch blockchain data
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
          <h1>Blockchain Explorer</h1>
          <p class="subtitle">Last updated: {{ new Date(currentTime).toLocaleString() }}</p>
        </div>
        <div class="actions">
          <button @click="mineBlock" :disabled="miningInProgress" class="mine-button">
            <span class="button-icon">⛏️</span>
            {{ miningInProgress ? 'Mining...' : 'Mine New Block' }}
          </button>
          <button @click="fetchChain" :disabled="loading" class="refresh-button">
            <span class="button-icon">🔄</span>
            {{ loading ? 'Syncing...' : 'Sync Chain' }}
          </button>
        </div>
      </div>
    </header>

    <!-- Error Message -->
    <div v-if="error" class="error-message">
      {{ error }}
    </div>
    <!-- Stats Section -->
    <BlockchainStats :blockchain="blockchain" />

    <!-- Graph Section -->
    <BlockchainGraph :blockchain="blockchain" />

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
                <div class="transaction-main">
                  <span class="sender">{{ tx.sender }}</span>
                  <span class="arrow">→</span>
                  <span class="recipient">{{ tx.recipient }}</span>
                  <span class="amount">{{ tx.amount }} coins</span>
                </div>
                <div v-if="tx.message" class="transaction-message">
                  "{{ tx.message }}"
                </div>
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
  font-size: 0.9em;
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
  transition: all 0.3s ease;
}

.mine-button {
  background-color: #4CAF50;
  color: white;
}

.refresh-button {
  background-color: #2196F3;
  color: white;
}

button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.block-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s ease;
}

.block-card:hover {
  transform: translateY(-2px);
}

.block-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
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
  font-size: 0.9em;
}

.transactions {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #eee;
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

.sender,
.recipient {
  font-family: monospace;
  padding: 2px 6px;
  background-color: #e9ecef;
  border-radius: 3px;
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

.error-message {
  background-color: #ffe6e6;
  color: #dc3545;
  padding: 15px;
  border-radius: 4px;
  margin-bottom: 20px;
}

.button-icon {
  font-size: 16px;
}

.transaction {
  flex-direction: column;
  padding: 15px;
}

.transaction-main {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
}

.transaction-message {
  margin-top: 8px;
  padding: 8px;
  background: #fff;
  border-radius: 4px;
  color: #666;
  font-style: italic;
  border-left: 3px solid #4CAF50;
}

.blockchain-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}
</style>