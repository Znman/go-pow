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

interface ConsensusProgress {
  attempt: number,
  proof: number,
  hash: string,
  found: boolean,
  blockIndex: number,
  message: string,
}

interface NodeState {
  id: number,
  eventSource: EventSource | null,
  progress: ConsensusProgress | null,
}



const consensusInProgress = ref(false)
const consensusProgress = ref<ConsensusProgress | null>(null)
let eventSource: EventSource | null = null

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8000';

const nodes = ref<NodeState[]>([
  { id: 0, eventSource: null, progress: null },
  { id: 1, eventSource: null, progress: null },
  { id: 2, eventSource: null, progress: null },
]);
const winnerNode = ref<number | null>(null);
const competitionActive = ref(false);
// const startConsensusStream = () => {
//   consensusInProgress.value = true
//   consensusProgress.value = null

//   eventSource = new EventSource('http://localhost:8000/mine/stream')
//   eventSource.onmessage = (e: MessageEvent) => {
//     const data = JSON.parse(e.data)
//     consensusProgress.value = data
//     if (data.found) {
//       setTimeout(() => {
//         consensusInProgress.value = false
//         eventSource?.close()
//         fetchChain()
//       }, 1200)
//     }
//   }
//   eventSource.onerror = () => {
//     consensusInProgress.value = false
//     eventSource?.close()
//   }
// }

onUnmounted(() => {
  if (eventSource) eventSource.close()
})
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

function startCompetition() {
  winnerNode.value = null;
  competitionActive.value = true;
  nodes.value.forEach((node, idx) => {
    // Clean up any previous EventSource
    if (node.eventSource) {
      node.eventSource.close();
    }
    node.progress = null;
    // Each node opens its own SSE connection (simulate three independent miners)
    const es = new EventSource(`${apiUrl}/mine/stream`);
    node.eventSource = es;

    es.onmessage = (event: MessageEvent) => {
      const progress: ConsensusProgress = JSON.parse(event.data);
      node.progress = progress;
      // If this node finds the block, declare as winner
      if (progress.found && winnerNode.value === null) {
        winnerNode.value = idx;
        competitionActive.value = false;
        // Stop the mining for the other nodes
        nodes.value.forEach((otherNode, otherIdx) => {
          if (otherNode.eventSource && otherIdx !== idx) {
            otherNode.eventSource.close();
          }
        });
      }
    };

    es.onerror = () => {
      es.close();
      node.eventSource = null;
      node.progress = node.progress || {
        attempt: 0,
        proof: 0,
        hash: '',
        found: false,
        blockIndex: 0,
        message: 'Connection error',
      };
    };
  });
}
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
          <h1>区块链挖矿模拟</h1>
          <p class="subtitle">最后一次更新: {{ new Date(currentTime).toLocaleString() }}</p>
        </div>
        <div class="actions">
          <button @click="mineBlock" :disabled="miningInProgress" class="mine-button">
            <span class="button-icon">⛏️</span>
            {{ miningInProgress ? '挖矿种...' : '挖取新的区块' }}
          </button>
          <button @click="fetchChain" :disabled="loading" class="refresh-button">
            <span class="button-icon">🔄</span>
            {{ loading ? '同步中...' : '同步区块链' }}
          </button>
        </div>
      </div>
      <div class="consensus-section">
        <h2 class="section-title">共识演示</h2>
        <div class="competition-controls">
          <button class="mine-btn" :disabled="competitionActive" @click="startCompetition">
            开始共识
          </button>
          <span v-if="winnerNode !== null" class="winner-msg">
            🏆 区块 {{ winnerNode + 1 }} 成功! 新区块将加入到链上.
          </span>
        </div>
        <div class="competition-nodes">
          <div v-for="(node, idx) in nodes" :key="node.id" :class="['node-card', { winner: winnerNode === idx }]">
            <div class="node-header">
              <span class="node-title">节点 {{ idx + 1 }}</span>
              <span v-if="winnerNode === idx" class="node-crown">🏆</span>
            </div>
            <div class="node-body">
              <template v-if="node.progress">
                <div class="node-row">
                  <span>尝试次数:</span>
                  <span class="node-num">{{ node.progress.attempt }}</span>
                </div>
                <div class="node-row">
                  <span>验证次数:</span>
                  <span class="node-num">{{ node.progress.proof }}</span>
                </div>
                <div class="node-row">
                  <span>哈希值:</span>
                  <span class="node-hash">{{ node.progress.hash }}</span>
                </div>
                <div class="node-row">
                  <span>状态:</span>
                  <span :class="['node-status', { found: node.progress.found }]">
                    {{ node.progress.message }}
                  </span>
                </div>
              </template>
              <template v-else>
                <div class="node-row node-waiting">等待中...</div>
              </template>
            </div>
          </div>
        </div>
      </div>
      <!-- <div class="actions">
        <button @click="startConsensusStream" :disabled="consensusInProgress" class="mine-button">
          <span class="button-icon">⛏️</span>
          {{ consensusInProgress ? 'Mining (Consensus in Progress)...' : 'Mine With Consensus Stream' }}
        </button>
      </div> -->
      <!-- Consensus Progress Section -->
      <!-- <div v-if="consensusInProgress" class="consensus-progress">
        <h3>Consensus Progress (Proof of Work)</h3>
        <div v-if="consensusProgress">
          <p><b>Block:</b> #{{ consensusProgress.blockIndex }}</p>
          <p><b>Attempt:</b> {{ consensusProgress.attempt }}</p>
          <p><b>Proof:</b> {{ consensusProgress.proof }}</p>
          <p><b>Hash:</b> {{ consensusProgress.hash }}</p>
          <p :class="{ found: consensusProgress.found }">{{ consensusProgress.message }}</p>
        </div>
        <div v-else>
          <em>Starting consensus process...</em>
        </div>
      </div> -->
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
            <h3>区块 #{{ block.index }}</h3>
            <span class="timestamp">{{ formatTimestamp(block.timestamp) }}</span>
          </div>
          <div class="block-details">
            <div class="detail-item">
              <span class="label">哈希值:</span>
              <span class="value hash">{{ block.hash }}</span>
            </div>
            <div class="detail-item">
              <span class="label">上一个区块的哈希值:</span>
              <span class="value hash">{{ block.previousHash }}</span>
            </div>
            <div class="detail-item">
              <span class="label">验证次数:</span>
              <span class="value">{{ block.proof }}</span>
            </div>
          </div>
          <div class="transactions">
            <h4>交易</h4>
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
              本区块无交易
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

.consensus-section {
  margin: 2em 0 2em 0;
  padding: 1.5em;
  background: #f8faff;
  border-radius: 16px;
  border: 1px solid #bde0fe;
  box-shadow: 0 2px 12px #bde0fe33;
}

.section-title {
  font-size: 1.4em;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 0.5em;
}

.section-desc {
  color: #64748b;
  margin-bottom: 1em;
}

.competition-controls {
  display: flex;
  align-items: center;
  gap: 1.5em;
  margin-bottom: 1.5em;
}

.mine-btn {
  padding: 0.6em 1.2em;
  background: linear-gradient(90deg, #38bdf8 40%, #60a5fa 100%);
  color: #fff;
  border: none;
  border-radius: 6px;
  font-weight: bold;
  font-size: 1em;
  cursor: pointer;
  box-shadow: 0 2px 8px #38bdf833;
  transition: background 0.25s;
}

.mine-btn:disabled {
  background: #cbd5e1;
  cursor: not-allowed;
  opacity: 0.7;
}

.winner-msg {
  font-size: 1.07em;
  color: #059669;
  font-weight: 600;
  margin-left: 1em;
  background: #d1fae5;
  border-radius: 6px;
  padding: 0.25em 1em;
}

.competition-nodes {}

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