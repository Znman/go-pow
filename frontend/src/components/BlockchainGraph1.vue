<template>
  <div class="chart-container">
    <Line
      v-if="blocks.length > 0"
      :data="chartData"
      :options="chartOptions"
    />
    <div v-else class="no-data">
      No blocks available to display
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  type ChartData,
  type ChartOptions
} from 'chart.js'
import { Line } from 'vue-chartjs'

// Register ChartJS components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
)

interface Block {
  index: number
  timestamp: number
  transactions: any[]
  proof: number
  previousHash: string
  hash: string
}

interface Props {
  blocks: Block[]
}

const props = defineProps<Props>()

const chartData = computed<ChartData<'line'>>(() => {
  const sortedBlocks = [...props.blocks].sort((a, b) => a.timestamp - b.timestamp)
  
  return {
    labels: sortedBlocks.map(block => {
      const date = new Date(block.timestamp * 1000)
      return date.toLocaleString('en-US', {
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }),
    datasets: [
      {
        label: 'Block Height',
        backgroundColor: 'rgba(75, 192, 192, 0.2)',
        borderColor: 'rgba(75, 192, 192, 1)',
        borderWidth: 2,
        data: sortedBlocks.map(block => block.index),
        fill: true,
        tension: 0.4
      },
      {
        label: 'Transactions',
        backgroundColor: 'rgba(255, 99, 132, 0.2)',
        borderColor: 'rgba(255, 99, 132, 1)',
        borderWidth: 2,
        data: sortedBlocks.map(block => block.transactions?.length || 0),
        fill: true,
        tension: 0.4,
        yAxisID: 'transactions'
      }
    ]
  }
})

const chartOptions = computed<ChartOptions<'line'>>(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    mode: 'index' as const,
    intersect: false
  },
  scales: {
    x: {
      ticks: {
        maxRotation: 45,
        minRotation: 45
      },
      title: {
        display: true,
        text: 'Time'
      }
    },
    y: {
      type: 'linear' as const,
      display: true,
      position: 'left' as const,
      beginAtZero: true,
      title: {
        display: true,
        text: 'Block Height'
      }
    },
    transactions: {
      type: 'linear' as const,
      display: true,
      position: 'right' as const,
      beginAtZero: true,
      grid: {
        drawOnChartArea: false
      },
      title: {
        display: true,
        text: 'Transactions'
      }
    }
  },
  plugins: {
    legend: {
      position: 'top' as const
    },
    title: {
      display: true,
      text: 'Blockchain Growth',
      font: {
        size: 16
      }
    },
    tooltip: {
      callbacks: {
        title: (items) => {
          if (!items.length) return ''
          const idx = items[0].dataIndex
          const blocks = props.blocks
          const block = blocks[idx]
          return `Block #${block.index} - ${new Date(block.timestamp * 1000).toLocaleString()}`
        },
        label: (context) => {
          const label = context.dataset.label || ''
          const value = context.parsed.y
          if (label === 'Block Height') {
            return `Height: ${value}`
          }
          return `Transactions: ${value}`
        }
      }
    }
  }
}))
</script>

<style scoped>
.chart-container {
  position: relative;
  height: 400px;
  width: 100%;
  margin: 20px 0;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.no-data {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #666;
  font-size: 1.1em;
}
</style>