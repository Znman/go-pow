<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'

interface TransactionForm {
  sender: string
  recipient: string
  amount: number | ''
}

const form = ref<TransactionForm>({
  sender: '',
  recipient: '',
  amount: ''
})

const loading = ref(false)
const error = ref<string | null>(null)
const success = ref<string | null>(null)

const submitTransaction = async () => {
  if (!form.value.sender || !form.value.recipient || form.value.amount === '') {
    error.value = 'Please fill in all fields'
    return
  }

  loading.value = true
  error.value = null
  success.value = null

  try {
    const response = await axios.post('http://localhost:8000/transactions/new', {
      sender: form.value.sender,
      recipient: form.value.recipient,
      amount: Number(form.value.amount)
    })

    success.value = 'Transaction added successfully!'
    // Clear form after successful submission
    form.value = {
      sender: '',
      recipient: '',
      amount: ''
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to add transaction'
    console.error('Error adding transaction:', err)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="transactions-container">
    <h1>New Transaction</h1>

    <form @submit.prevent="submitTransaction" class="transaction-form">
      <div class="form-group">
        <label for="sender">Sender</label>
        <input
          id="sender"
          v-model="form.sender"
          type="text"
          required
          placeholder="Enter sender's name"
        />
      </div>

      <div class="form-group">
        <label for="recipient">Recipient</label>
        <input
          id="recipient"
          v-model="form.recipient"
          type="text"
          required
          placeholder="Enter recipient's name"
        />
      </div>

      <div class="form-group">
        <label for="amount">Amount</label>
        <input
          id="amount"
          v-model="form.amount"
          type="number"
          step="0.01"
          required
          placeholder="Enter amount"
        />
      </div>

      <button type="submit" :disabled="loading" class="submit-button">
        {{ loading ? 'Adding Transaction...' : 'Add Transaction' }}
      </button>
    </form>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div v-if="success" class="success-message">
      {{ success }}
    </div>
  </div>
</template>

<style scoped>
.transactions-container {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

.transaction-form {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  color: #2c3e50;
}

input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.submit-button {
  width: 100%;
  padding: 10px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.submit-button:hover:not(:disabled) {
  background-color: #45a049;
}

.submit-button:disabled {
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

.success-message {
  margin-top: 15px;
  padding: 10px;
  background-color: #e6ffe6;
  color: #28a745;
  border-radius: 4px;
}
</style>