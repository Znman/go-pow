<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import axios from 'axios'
import BlockchainGraph from './BlockchainGraph.vue'
import BlockchainStats from './BlockchainStats.vue'
import ConsensusSection from './ConsensusSection.vue'

interface Transaction {
  sender: string
  recipient: string
  amount: number
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

// State Management
const blockchain = ref<BlockchainData | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const miningInProgress = ref(false)
const currentTime = ref(new Date().toISOString())
const selectedBlock = ref<Block | null>(null)
const showingConsensus = ref(false)
const refreshInterval = ref<number | null>(null)

// Constants for consensus visualization
const MINING_DIFFICULTY = 4
const TARGET_PREFIX = '0'.repeat(MINING_DIFFICULTY)

// Fetch blockchain data
const fetchChain = async () => {
  loading.value = true
  error.value = null

  try {
    const response = await axios.get<BlockchainData>('http://localhost:8000/chain')
    blockchain.value = response.data
    currentTime.value = new Date().toISOString()

    // Update selected block if it exists in the new chain data
    if (selectedBlock.value && blockchain.value) {
      const updatedBlock = blockchain.value.chain.find(
        block => block.index === selectedBlock.value?.index
      )
      if (updatedBlock) {
        selectedBlock.value = updatedBlock
      }
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to load blockchain'
    console.error('Error loading blockchain:', err)
  } finally {
    loading.value = false
  }
}

// Consensus visualization controls
const startConsensusVisualization = (block: Block) => {
  if (miningInProgress.value) return

  // Create a deep copy of the block to prevent reference issues
  selectedBlock.value = JSON.parse(JSON.stringify(block))
  showingConsensus.value = true

  // Pause auto-refresh while showing consensus
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value)
    refreshInterval.value = null
  }
}

const stopConsensusVisualization = () => {
  showingConsensus.value = false
  selectedBlock.value = null

  // Restart auto-refresh if it was previously running
  if (!refreshInterval.value) {
    refreshInterval.value = window.setInterval(fetchChain, 30000)
  }
}

// Mine new block
const mineBlock = async () => {
  if (miningInProgress.value) return

  miningInProgress.value = true
  error.value = null
  showingConsensus.value = true

  try {
    // Create temporary block for visualization
    if (blockchain.value && blockchain.value.chain.length > 0) {
      const lastBlock = blockchain.value.chain[blockchain.value.chain.length - 1]
      selectedBlock.value = {
        index: lastBlock.index + 1,
        timestamp: Math.floor(Date.now() / 1000),
        transactions: blockchain.value.currentTransactions,
        proof: 0,
        previousHash: lastBlock.hash,
        hash: ''
      }
    }

    const startTime = Date.now()
    await axios.post('http://localhost:8000/mine')
    const miningTime = Date.now() - startTime
    console.log(`Block mined successfully in ${miningTime}ms`)

    await fetchChain()
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Mining failed'
    console.error('Error mining block:', err)
  } finally {
    miningInProgress.value = false
    // Don't automatically hide consensus visualization after mining
  }
}

// Format timestamp
const formatTimestamp = (timestamp: number): string => {
  return new Date(timestamp * 1000).toLocaleString()
}

// Format hash for display
const formatHash = (hash: string): string => {
  return hash.length > 16
    ? `${hash.substring(0, 8)}...${hash.substring(hash.length - 8)}`
    : hash
}

// Watch for mining state changes
watch(miningInProgress, (newValue) => {
  if (!newValue && !selectedBlock.value) {
    showingConsensus.value = false
  }
})

// Lifecycle hooks
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
          <button @click="mineBlock" :disabled="miningInProgress" class="mine-button"
            :class="{ 'mining': miningInProgress }">
            <span class="button-icon">⛏️</span>
            {{ miningInProgress ? 'Mining in Progress...' : 'Mine New Block' }}
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

    <!-- Consensus Section -->
    <div v-if="showingConsensus" class="consensus-container">
      <div class="consensus-header">
        <h3>
          Consensus Process
          {{ selectedBlock ? `for Block #${selectedBlock.index}` : '' }}
        </h3>
        <button v-if="!miningInProgress" class="close-button" @click="stopConsensusVisualization"
          title="Close consensus view">
          ✕
        </button>
      </div>
      <ConsensusSection v-if="showingConsensus" :block="selectedBlock" :isActive="showingConsensus" :difficulty="4"
        :targetPrefix="'0000'" @close="stopConsensusVisualization" />
    </div>

    <!-- Stats Section -->
    <BlockchainStats :blockchain="blockchain" />

    <!-- Graph Section -->
    <BlockchainGraph :blockchain="blockchain" />

    <!-- Blockchain Data -->
    <div class="blockchain-data" v-if="blockchain">
      <div class="blocks-container">
        <div v-for="block in [...blockchain.chain].reverse()" :key="block.index" class="block-card"
          :class="{ 'selected': selectedBlock?.index === block.index }">
          <div class="block-header">
            <h3>Block #{{ block.index }}</h3>
            <div class="block-actions">
              <button v-if="!miningInProgress && (!showingConsensus || selectedBlock?.index !== block.index)"
                @click="startConsensusVisualization(block)" class="consensus-button">
                <span class="button-icon">🔍</span>
                View Consensus
              </button>
              <button v-else-if="showingConsensus && selectedBlock?.index === block.index"
                @click="stopConsensusVisualization" class="consensus-button active">
                <span class="button-icon">✕</span>
                Hide Consensus
              </button>
              <span class="timestamp">{{ formatTimestamp(block.timestamp) }}</span>
            </div>
          </div>

          <div class="block-details">
            <div class="detail-item">
              <span class="label">Hash:</span>
              <span class="value hash" :title="block.hash">
                {{ formatHash(block.hash) }}
              </span>
            </div>
            <div class="detail-item">
              <span class="label">Previous Hash:</span>
              <span class="value hash" :title="block.previousHash">
                {{ formatHash(block.previousHash) }}
              </span>
            </div>
            <div class="detail-item">
              <span class="label">Proof:</span>
              <span class="value">{{ block.proof }}</span>
            </div>
          </div>

          <div class="transactions">
            <h4>
              Transactions
              <span v-if="block.transactions.length" class="transaction-count">
                ({{ block.transactions.length }})
              </span>
            </h4>
            <div v-if="block.transactions && block.transactions.length > 0" class="transactions-list">
              <div v-for="(tx, index) in block.transactions" :key="index" class="transaction">
                <div class="transaction-main">
                  <span class="sender" :title="tx.sender">
                    {{ formatHash(tx.sender) }}
                  </span>
                  <span class="arrow">→</span>
                  <span class="recipient" :title="tx.recipient">
                    {{ formatHash(tx.recipient) }}
                  </span>
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
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-section h1 {
  margin: 0;
  color: #2c3e50;
  font-size: 1.8em;
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
.refresh-button,
.consensus-button {
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

.mine-button:hover {
  background-color: #45a049;
}

.mine-button.mining {
  animation: pulse 1.5s infinite;
}

.refresh-button {
  background-color: #2196F3;
  color: white;
}

.refresh-button:hover {
  background-color: #1e88e5;
}

.consensus-button {
  background-color: #6c757d;
  color: white;
  padding: 6px 12px;
  font-size: 0.9em;
}

.consensus-button:hover {
  background-color: #5a6268;
}

.consensus-button.active {
  background-color: #dc3545;
}

.consensus-button.active:hover {
  background-color: #c82333;
}

.consensus-container {
  background: white;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.consensus-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.consensus-header h3 {
  margin: 0;
  color: #2c3e50;
}

.close-button {
  background: none;
  border: none;
  color: #666;
  font-size: 20px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.close-button:hover {
  background: #f8f9fa;
  color: #dc3545;
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
  transition: all 0.3s ease;
}

.block-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.block-card.selected {
  border: 2px solid #4CAF50;
}

.block-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.block-header h3 {
  margin: 0;
  color: #2c3e50;
  font-size: 1.2em;
}

.block-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.block-details {
  margin-bottom: 15px;
}

.detail-item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.label {
  color: #666;
  font-weight: bold;
  margin-right: 8px;
  min-width: 120px;
}

.hash {
  font-family: monospace;
  word-break: break-all;
  background-color: #f8f9fa;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.9em;
  cursor: help;
}

.timestamp {
  color: #666;
  font-size: 0.9em;
}

.transactions {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #eee;
}

.transactions h4 {
  margin: 0 0 10px 0;
  color: #2c3e50;
}

.transaction-count {
  color: #666;
  font-size: 0.9em;
  font-weight: normal;
}

.transactions-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.transaction {
  background-color: #f8f9fa;
  border-radius: 4px;
  padding: 12px;
}

.transaction-main {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.sender,
.recipient {
  font-family: monospace;
  padding: 2px 6px;
  background-color: #e9ecef;
  border-radius: 3px;
  font-size: 0.9em;
  cursor: help;
}

.arrow {
  color: #666;
}

.amount {
  margin-left: auto;
  font-weight: bold;
  color: #4CAF50;
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

@keyframes pulse {
  0% {
    opacity: 1;
  }

  50% {
    opacity: 0.7;
  }

  100% {
    opacity: 1;
  }
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    gap: 16px;
  }

  .actions {
    width: 100%;
    justify-content: stretch;
  }

  .mine-button,
  .refresh-button {
    flex: 1;
  }

  .block-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }

  .block-actions {
    width: 100%;
    justify-content: space-between;
  }

  .transaction-main {
    flex-direction: column;
    align-items: flex-start;
  }

  .amount {
    margin-left: 0;
    margin-top: 8px;
  }
}
</style>