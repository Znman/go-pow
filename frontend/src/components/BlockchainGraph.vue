<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale
} from 'chart.js'
import { Line } from 'vue-chartjs'
import 'chartjs-adapter-date-fns' // We'll need to install this

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale
)

interface Block {
  index: number
  timestamp: number
  transactions: any[]
  proof: number
  previousHash: string
  hash: string
}

interface BlockchainData {
  chain: Block[]
  currentTransactions: any[]
}

const props = defineProps<{
  blockchain: BlockchainData | null
}>()

const chartData = computed(() => {
  if (!props.blockchain?.chain) {
    return {
      datasets: [{
        label: '区块链数目',
        data: [],
        fill: false,
        borderColor: '#4CAF50',
        tension: 0.1
      }]
    }
  }

  // Sort blocks by timestamp
  const sortedBlocks = [...props.blockchain.chain].sort((a, b) => a.timestamp - b.timestamp)

  // Convert timestamp to JavaScript Date objects and create data points
  const dataPoints = sortedBlocks.map(block => ({
    x: block.timestamp * 1000, // Convert Unix timestamp to milliseconds
    y: block.index
  }))

  return {
    datasets: [{
      label: '区块链数目',
      data: dataPoints,
      fill: false,
      borderColor: '#4CAF50',
      tension: 0.1,
      pointRadius: 4,
      pointHoverRadius: 6,
      pointBackgroundColor: '#4CAF50'
    }]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: {
      type: 'time' as const,
      time: {
        unit: 'hour' as const,
        displayFormats: {
          hour: 'HH:mm'
        },
        tooltipFormat: 'yyyy-MM-dd HH:mm'
      },
      title: {
        display: true,
        text: '时间 (小时)',
        color: '#666',
        font: {
          size: 12,
          weight: 'bold'
        }
      }
    },
    y: {
      beginAtZero: true,
      title: {
        display: true,
        text: '区块数量',
        color: '#666',
        font: {
          size: 12,
          weight: 'bold'
        }
      },
      ticks: {
        stepSize: 1
      }
    }
  },
  plugins: {
    legend: {
      display: true,
      position: 'top' as const
    },
    title: {
      display: true,
      text: '区块链数目随时间变化图',
      font: {
        size: 16,
        weight: 'bold'
      }
    },
    tooltip: {
      callbacks: {
        label: (context: any) => {
          return `Block ${context.parsed.y}`
        },
        title: (tooltipItems: any[]) => {
          const date = new Date(tooltipItems[0].parsed.x)
          return date.toLocaleString()
        }
      }
    }
  },
  interaction: {
    intersect: false,
    mode: 'index'
  }
}
</script>

<template>
  <div class="chart-container">
    <Line :data="chartData" :options="chartOptions" style="height: 400px;" />
  </div>
</template>

<style scoped>
.chart-container {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

@media (max-width: 768px) {
  .chart-container {
    padding: 10px;
  }
}
</style>