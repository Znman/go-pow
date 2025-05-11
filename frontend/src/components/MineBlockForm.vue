<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'

const emit = defineEmits(['blockMined'])

interface BlockData {
  author: string
  description: string
  metadata: string
}

const formData = ref<BlockData>({
  author: 'Znman',
  description: '',
  metadata: ''
})

const loading = ref(false)
const error = ref<string | null>(null)

const mineBlock = async () => {
  loading.value = true
  error.value = null

  try {
    const response = await axios.post('http://localhost:8000/mine', formData.value)
    emit('blockMined')
    clearForm()
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to mine block'
  } finally {
    loading.value = false
  }
}

const clearForm = () => {
  formData.value.description = ''
  formData.value.metadata = ''
}
</script>

<template>
  <div class="mine-block-form">
    <h3>Mine New Block</h3>
    <form @submit.prevent="mineBlock">
      <div class="form-group">
        <label for="author">Author</label>
        <input 
          id="author"
          v-model="formData.author"
          type="text"
          placeholder="Enter author name"
          required
        />
      </div>

      <div class="form-group">
        <label for="description">Description</label>
        <textarea
          id="description"
          v-model="formData.description"
          placeholder="Enter block description"
          rows="3"
        ></textarea>
      </div>

      <div class="form-group">
        <label for="metadata">Metadata</label>
        <input
          id="metadata"
          v-model="formData.metadata"
          type="text"
          placeholder="Enter additional metadata"
        />
      </div>

      <button type="submit" :disabled="loading" class="mine-button">
        <span class="button-icon">⛏️</span>
        {{ loading ? 'Mining...' : 'Mine Block' }}
      </button>
    </form>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>
  </div>
</template>

<style scoped>
.mine-block-form {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  color: #2c3e50;
  font-weight: bold;
}

input,
textarea {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

textarea {
  resize: vertical;
}

.mine-button {
  padding: 12px 24px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  transition: background-color 0.3s;
}

.mine-button:hover:not(:disabled) {
  background-color: #45a049;
}

.mine-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.error-message {
  margin-top: 15px;
  padding: 10px;
  background-color: #ffe6e6;
  color: #dc3545;
  border-radius: 4px;
}
</style>