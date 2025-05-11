<script setup lang="ts">
import { computed } from 'vue'
import type { BlockchainData } from './interfaces'

const props = defineProps<{
  blockchain: BlockchainData | null
}>()

const stats = computed(() => {
  if (!props.blockchain) {
    return {
      totalBlocks: 0,
      totalTransactions: 0,
      averageBlockTime: 0
    }
  }

  const chain = props.blockchain.chain
  const totalBlocks = chain.length
  let totalTransactions = 0
  let totalBlockTime = 0

  for (let i = 1; i < chain.length; i++) {
    totalTransactions += chain[i].transactions.length
    totalBlockTime += chain[i].timestamp - chain[i-1].timestamp
  }

  const averageBlockTime = chain.length > 1 ? totalBlockTime / (chain.length - 1) : 0

  return {
    totalBlocks,
    totalTransactions,
    averageBlockTime: Math.round(averageBlockTime)
  }
})
</script>

<template>
  <div class="stats-container">
    <div class="stat-card">
      <h3>Total Blocks</h3>
      <p>{{ stats.totalBlocks }}</p>
    </div>
    <div class="stat-card">
      <h3>Total Transactions</h3>
      <p>{{ stats.totalTransactions }}</p>
    </div>
    <div class="stat-card">
      <h3>Average Block Time</h3>
      <p>{{ stats.averageBlockTime }}s</p>
    </div>
  </div>
</template>

<style scoped>
.stats-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.stat-card h3 {
  margin: 0;
  color: #666;
  font-size: 0.9em;
  text-transform: uppercase;
}

.stat-card p {
  margin: 10px 0 0;
  font-size: 1.8em;
  font-weight: bold;
  color: #4CAF50;
}
</style>