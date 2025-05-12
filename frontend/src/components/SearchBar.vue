<script setup lang="ts">
import { ref, watch } from 'vue'
import { debounce } from 'lodash'

interface SearchCriteria {
  type: 'hash' | 'block' | 'transaction' | 'message' | 'timerange'
  value: string | number | { start: string; end: string }
}

const searchType = ref<SearchCriteria['type']>('block')
const searchValue = ref('')
const startDate = ref('')
const endDate = ref('')

const emit = defineEmits<{
  (e: 'search', criteria: SearchCriteria): void
  (e: 'clear'): void
}>()

const handleSearch = debounce(() => {
  if (!searchValue.value && searchType.value !== 'timerange') {
    emit('clear')
    return
  }

  const criteria: SearchCriteria = {
    type: searchType.value,
    value: searchType.value === 'timerange' 
      ? { start: startDate.value, end: endDate.value }
      : searchType.value === 'block' 
        ? parseInt(searchValue.value) || 0
        : searchValue.value
  }
  
  emit('search', criteria)
}, 300)

watch([searchType], () => {
  searchValue.value = ''
  startDate.value = ''
  endDate.value = ''
})

const clearSearch = () => {
  searchValue.value = ''
  startDate.value = ''
  endDate.value = ''
  emit('clear')
}
</script>

<template>
  <div class="search-container">
    <div class="search-header">
      <h3>Search Blockchain</h3>
      <button @click="clearSearch" class="clear-button">
        Clear Search
      </button>
    </div>
    
    <div class="search-controls">
      <div class="search-type">
        <select v-model="searchType">
          <option value="block">Block Number</option>
          <option value="hash">Block Hash</option>
          <option value="transaction">Transaction</option>
          <option value="message">Message</option>
          <option value="timerange">Time Range</option>
        </select>
      </div>

      <template v-if="searchType === 'timerange'">
        <div class="time-range-inputs">
          <input
            type="datetime-local"
            v-model="startDate"
            placeholder="Start Date"
            @input="handleSearch"
          />
          <span>to</span>
          <input
            type="datetime-local"
            v-model="endDate"
            placeholder="End Date"
            @input="handleSearch"
          />
        </div>
      </template>
      <template v-else>
        <div class="search-input">
          <input
            v-model="searchValue"
            :type="searchType === 'block' ? 'number' : 'text'"
            :placeholder="'Search by ' + searchType"
            @input="handleSearch"
          />
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.search-container {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.search-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.search-header h3 {
  margin: 0;
  color: #2c3e50;
}

.clear-button {
  padding: 6px 12px;
  background: #dc3545;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.search-controls {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.search-type {
  min-width: 150px;
}

.search-type select {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.search-input {
  flex: 1;
  min-width: 200px;
}

.search-input input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.time-range-inputs {
  display: flex;
  gap: 8px;
  align-items: center;
  flex: 1;
}

.time-range-inputs input {
  flex: 1;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.time-range-inputs span {
  color: #666;
}

@media (max-width: 768px) {
  .search-controls {
    flex-direction: column;
  }

  .search-type,
  .search-input,
  .time-range-inputs {
    width: 100%;
  }
}
</style>