<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

interface MiningProgress {
  attempts: number
  hashRate: number
  currentHash: string
}

const props = defineProps<{
  isActive: boolean
  block?: {
    miningStats?: MiningProgress
  }
}>()

const emit = defineEmits(['update:progress'])

const miningStats = ref<MiningProgress>({
  attempts: 0,
  hashRate: 0,
  currentHash: ''
})

let pollInterval: number | null = null

const updateMiningStats = () => {
  if (props.block?.miningStats) {
    miningStats.value = props.block.miningStats
    emit('update:progress', props.block.miningStats)
  }
}

onMounted(() => {
  if (props.isActive) {
    pollInterval = window.setInterval(updateMiningStats, 100)
  }
})

onUnmounted(() => {
  if (pollInterval !== null) {
    clearInterval(pollInterval)
  }
})

// Watch for block updates
watch(() => props.block?.miningStats, (newStats) => {
  if (newStats) {
    miningStats.value = newStats
  }
}, { deep: true })
</script>

<template>
  <div class="consensus-container">
    <div class="mining-header">
      <h3>Mining Progress</h3>
      <div class="mining-indicator" :class="{ active: isActive }"></div>
    </div>

    <div class="mining-stats">
      <div class="stat-box">
        <div class="stat-label">Hash Attempts</div>
        <div class="stat-value">{{ miningStats.attempts.toLocaleString() }}</div>
      </div>
      
      <div class="stat-box">
        <div class="stat-label">Hash Rate</div>
        <div class="stat-value">
          {{ Math.floor(miningStats.hashRate).toLocaleString() }} H/s
        </div>
      </div>
    </div>

    <div class="current-hash">
      <div class="hash-label">Current Hash:</div>
      <div class="hash-value" :title="miningStats.currentHash">
        <template v-if="miningStats.currentHash">
          <span 
            v-for="(char, index) in miningStats.currentHash.slice(0, 4)" 
            :key="index"
            :class="{ 'target-char': char === '0' }"
          >{{ char }}</span>
          <span class="hash-separator">...</span>
          <span>{{ miningStats.currentHash.slice(-4) }}</span>
        </template>
        <template v-else>
          Starting...
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
.consensus-container {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  margin: 20px 0;
}

.mining-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.mining-header h3 {
  margin: 0;
  color: #2c3e50;
}

.mining-indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background-color: #dc3545;
  transition: background-color 0.3s ease;
}

.mining-indicator.active {
  background-color: #28a745;
  animation: pulse 1s infinite;
}

.mining-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.stat-box {
  background: white;
  padding: 16px;
  border-radius: 6px;
  text-align: center;
}

.stat-label {
  color: #666;
  font-size: 0.9em;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 1.2em;
  font-weight: bold;
  font-family: monospace;
}

.current-hash {
  background: white;
  padding: 16px;
  border-radius: 6px;
}

.hash-label {
  color: #666;
  font-size: 0.9em;
  margin-bottom: 8px;
}

.hash-value {
  font-family: monospace;
  background: #f8f9fa;
  padding: 8px;
  border-radius: 4px;
  word-break: break-all;
  letter-spacing: 1px;
}

.target-char {
  color: #28a745;
  font-weight: bold;
}

.hash-separator {
  color: #666;
  margin: 0 4px;
}

@keyframes pulse {
  0% { opacity: 1; }
  50% { opacity: 0.6; }
  100% { opacity: 1; }
}

</style>