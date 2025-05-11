<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'

interface Transaction {
  Sender: string
  Recipient: string
  Amount: number
}

interface Block {
  index: number
  timestamp: number
  transactions: Transaction[] | null
  proof: number
  previousHash: string
  hash: string
}

interface MiningResponse {
  message: string
  block: Block
}

const isMining = ref(false)
const lastMinedBlock = ref<Block | null>(null)
const miningDuration = ref(0)
const error = ref<string | null>(null)
const blockchainStatus = ref({
  totalBlocks: 0,
  pendingTransactions: 0
})

// Fetch blockchain status
const fetchBlockchainStatus = async () => {
  try {
    const response = await axios.get('http://localhost:8000/chain')
    if (response.data && response.data.chain) {
      blockchainStatus.value.totalBlocks = response.data.chain.length
      blockchainStatus.value.pendingTransactions = 
        response.data.currentTransactions ? response.data.currentTransactions.length : 0
    }
  } catch (err) {
    console.error('Error fetching blockchain status:', err)
  }
}

// Start mining process
const startMining = async () => {
  isMining.value = true
  error.value = null
  const startTime = Date.now()

  try {
    const response = await axios.post<MiningResponse>('http://localhost:8000/mine')
    lastMinedBlock.value = response.data.block
    miningDuration.value = (Date.now() - startTime) / 1000 // Convert to seconds
    await fetchBlockchainStatus() // Update blockchain status after mining
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Mining failed'
    console.error('Error during mining:', err)
  } finally {
    isMining.value = false
  }
}

// Format timestamp to local date and time
const formatTimestamp = (timestamp: number): string => {
  if (!timestamp) return 'N/A'
  return new Date(timestamp * 1000).toLocaleString()
}

// Format hash with ellipsis
const formatHash = (hash: string | undefined): string => {
  if (!hash) return 'N/A'
  return hash.length > 16 
    ? `${hash.substring(0, 8)}...${hash.substring(hash.length - 8)}`
    : hash
}

// Format number with commas
const formatNumber = (num: number): string => {
  return num.toLocaleString()
}

// Check if block has valid data
const isValidBlock = (block: Block | null): boolean => {
  return block !== null && 
         typeof block.index !== 'undefined' && 
         typeof block.timestamp !== 'undefined'
}

onMounted(async () => {
  await fetchBlockchainStatus()
})
</script>

<template>
  <div class="mining-container">
    <h2>Blockchain Mining Dashboard</h2>
    
    <!-- Blockchain Status -->
    <div class="status-panel">
      <div class="status-item">
        <h3>Blockchain Status</h3>
        <p>Total Blocks: {{ formatNumber(blockchainStatus.totalBlocks) }}</p>
        <p>Pending Transactions: {{ blockchainStatus.pendingTransactions }}</p>
      </div>
    </div>

    <!-- Mining Control -->
    <div class="mining-control">
      <button 
        @click="startMining" 
        :disabled="isMining"
        class="mine-button"
      >
        <span class="button-icon">⛏️</span>
        {{ isMining ? 'Mining in Progress...' : 'Start Mining' }}
      </button>
    </div>

    <!-- Mining Status -->
    <div v-if="isMining" class="mining-status">
      <div class="spinner"></div>
      <p>Mining new block... Please wait</p>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="error-message">
      <span class="error-icon">⚠️</span>
      {{ error }}
    </div>

    <!-- Last Mined Block -->
    <div v-if="lastMinedBlock && isValidBlock(lastMinedBlock)" class="last-mined-block">
      <h3>Last Mined Block</h3>
      <div class="block-info">
        <div class="info-grid">
          <div class="info-item">
            <label>Block #:</label>
            <span>{{ formatNumber(lastMinedBlock.index) }}</span>
          </div>
          <div class="info-item">
            <label>Timestamp:</label>
            <span>{{ formatTimestamp(lastMinedBlock.timestamp) }}</span>
          </div>
          <div class="info-item">
            <label>Mining Duration:</label>
            <span>{{ miningDuration.toFixed(2) }} seconds</span>
          </div>
          <div class="info-item">
            <label>Proof:</label>
            <span>{{ formatNumber(lastMinedBlock.proof) }}</span>
          </div>
          <div class="info-item full-width">
            <label>Block Hash:</label>
            <span class="hash">{{ lastMinedBlock.hash }}</span>
          </div>
          <div class="info-item full-width">
            <label>Previous Hash:</label>
            <span class="hash">{{ lastMinedBlock.previousHash }}</span>
          </div>
        </div>
        
        <!-- Transactions -->
        <div class="transactions-section">
          <h4>Transactions in this block</h4>
          <div v-if="lastMinedBlock.transactions && lastMinedBlock.transactions.length > 0" 
               class="transactions">
            <div v-for="(tx, index) in lastMinedBlock.transactions" 
                 :key="index" 
                 class="transaction-item">
              <div class="transaction-details">
                <div class="address">
                  <span class="label">From:</span>
                  <span class="value">{{ tx.Sender }}</span>
                </div>
                <div class="arrow">→</div>
                <div class="address">
                  <span class="label">To:</span>
                  <span class="value">{{ tx.Recipient }}</span>
                </div>
                <div class="amount">
                  <span class="value">{{ tx.Amount }} coins</span>
                </div>
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
</template>

<style scoped>
.mining-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.status-panel {
  background-color: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.status-item h3 {
  margin-top: 0;
  color: #2c3e50;
}

.mining-control {
  margin: 20px 0;
  text-align: center;
}

.mine-button {
  padding: 12px 24px;
  font-size: 16px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.mine-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.mine-button:hover:not(:disabled) {
  background-color: #45a049;
  transform: translateY(-1px);
}

.mining-status {
  text-align: center;
  margin: 20px 0;
}

.spinner {
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin: 0 auto 10px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-message {
  color: #dc3545;
  background-color: #ffe6e6;
  padding: 12px;
  border-radius: 4px;
  margin: 10px 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.last-mined-block {
  background-color: #ffffff;
  border-radius: 8px;
  padding: 20px;
  margin-top: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.info-item label {
  font-weight: bold;
  color: #666;
  font-size: 0.9em;
}

.info-item .hash {
  font-family: monospace;
  word-break: break-all;
  background-color: #f8f9fa;
  padding: 8px;
  border-radius: 4px;
}

.transactions-section {
  margin-top: 20px;
  border-top: 1px solid #eee;
  padding-top: 20px;
}

.transaction-item {
  background-color: #f8f9fa;
  border-radius: 4px;
  padding: 12px;
  margin-bottom: 8px;
}

.transaction-details {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.address {
  display: flex;
  gap: 4px;
  align-items: center;
}

.address .label {
  color: #666;
  font-size: 0.9em;
}

.arrow {
  color: #666;
  font-weight: bold;
}

.amount {
  margin-left: auto;
  font-weight: bold;
  color: #4CAF50;
}

.no-transactions {
  color: #666;
  text-align: center;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 4px;
}
</style>