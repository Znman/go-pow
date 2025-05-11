<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
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

interface SearchResponse {
  results: Block[]
}

const blockchain = ref<BlockchainData | null>(null)
const searchResults = ref<Block[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const searchQuery = ref('')
const searchType = ref('block') // 'block' or 'transaction'
const currentTime = ref(new Date().toISOString())
const refreshInterval = ref<number | null>(null)

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

// Search functionality
const searchBlocks = async () => {
  if (!searchQuery.value) {
    searchResults.value = []
    return
  }

  loading.value = true
  error.value = null

  try {
    const response = await axios.get<SearchResponse>(
      `http://localhost:8000/search?q=${encodeURIComponent(searchQuery.value)}&type=${searchType.value}`
    )
    searchResults.value = response.data.results
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Search failed'
    console.error('Search error:', err)
  } finally {
    loading.value = false
  }
}

// Format timestamp
const formatTimestamp = (timestamp: number): string => {
  return new Date(timestamp * 1000).toLocaleString()
}

// Display blocks
const displayedBlocks = computed(() => {
  return searchResults.value.length > 0 ? searchResults.value : 
         blockchain.value ? blockchain.value.chain : []
})

// Stats
const stats = computed(() => {
  if (!blockchain.value) return null

  const chain = blockchain.value.chain
  return {
    totalBlocks: chain.length,
    totalTransactions: chain.reduce((sum, block) => sum + block.transactions.length, 0),
    lastBlockTime: chain[chain.length - 1].timestamp
  }
})

// Setup auto-refresh
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
  <div class="explorer-container">
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

    <!-- Stats Section -->
    <div v-if="stats" class="stats-section">
      <div class="stat-card">
        <h3>Total Blocks</h3>
        <p>{{ stats.totalBlocks }}</p>
      </div>
      <div class="stat-card">
        <h3>Total Transactions</h3>
        <p>{{ stats.totalTransactions }}</p>
      </div>
      <div class="stat-card">
        <h3>Last Block</h3>
        <p>{{ formatTimestamp(stats.lastBlockTime) }}</p>
      </div>
    </div>

    <!-- Search Section -->
    <div class="search-section">
      <div class="search-box">
        <select v-model="searchType" class="search-type" @change="searchBlocks">
          <option value="block">Block Number</option>
          <option value="transaction">Transaction</option>
        </select>
        <input
          v-model="searchQuery"
          @input="searchBlocks"
          type="text"
          :placeholder="searchType === 'block' ? 'Enter block number...' : 'Search transactions...'"
          class="search-input"
        />
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <!-- Results -->
    <div class="blockchain-data">
      <div class="blocks-container">
        <div v-for="block in displayedBlocks" 
             :key="block.index" 
             class="block-card">
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
            <h4>Transactions ({{ block.transactions.length }})</h4>
            <div v-if="block.transactions.length > 0" class="transactions-list">
              <div v-for="(tx, index) in block.transactions" 
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
.explorer-container {
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

.stats-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  text-align: center;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.stat-card h3 {
  margin: 0 0 10px;
  color: #666;
  font-size: 0.9em;
}

.stat-card p {
  margin: 0;
  font-size: 1.5em;
  font-weight: bold;
  color: #2c3e50;
}

.search-section {
  margin-bottom: 20px;
}

.search-box {
  display: flex;
  gap: 10px;
  background-color: white;
  padding: 15px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.search-type {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: white;
}

.search-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
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
}

.refresh-button:hover:not(:disabled) {
  background-color: #45a049;
}

.refresh-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.block-card {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.block-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.block-details {
  margin-bottom: 15px;
}

.detail-item {
  margin-bottom: 10px;
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

.hash {
  font-family: monospace;
  word-break: break-all;
  background-color: #f8f9fa;
  padding: 4px 8px;
  border-radius: 4px;
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
</style>