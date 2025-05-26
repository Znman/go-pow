<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'

interface Transaction {
  sender: string
  recipient: string
  amount: number
  message: string
}

interface TransactionResponse {
  message: string
  transaction: Transaction
}

const transactionType = ref<'transfer' | 'message'>('transfer')
const form = ref({
  sender: '',
  recipient: '',
  amount: 0,
  message: ''
})
const loading = ref(false)
const error = ref<string | null>(null)
const success = ref<string | null>(null)

const submitTransaction = async () => {
  loading.value = true
  error.value = null
  success.value = null

  try {
    // Validate form based on transaction type
    if (transactionType.value === 'transfer') {
      if (!form.value.sender || !form.value.recipient || !form.value.amount) {
        error.value = 'Please fill in all required fields'
        return
      }
    } else {
      if (!form.value.sender || !form.value.message) {
        error.value = 'Please provide sender and message'
        return
      }
      // For message-only transactions, set default values
      form.value.recipient = 'BLOCKCHAIN'
      form.value.amount = 0
    }

    const response = await axios.post<TransactionResponse>('http://localhost:8000/transactions/new', {
      sender: form.value.sender,
      recipient: form.value.recipient,
      amount: form.value.amount,
      message: form.value.message
    })

    success.value = response.data.message
    // Reset form
    form.value = {
      sender: '',
      recipient: '',
      amount: 0,
      message: ''
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to submit transaction'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="transactions-container">
    <h2>交易/信息上链</h2>

    <div class="transaction-type-selector">
      <button :class="{ active: transactionType === 'transfer' }" @click="transactionType = 'transfer'">
        💰 新交易
      </button>
      <button :class="{ active: transactionType === 'message' }" @click="transactionType = 'message'">
        ✉️ 新信息
      </button>
    </div>

    <form @submit.prevent="submitTransaction" class="transaction-form">
      <div class="form-group">
        <label for="sender">发送方</label>
        <input id="sender" v-model="form.sender" placeholder="发送方的名称或地址" required />
      </div>

      <template v-if="transactionType === 'transfer'">
        <div class="form-group">
          <label for="recipient">接受方</label>
          <input id="recipient" v-model="form.recipient" placeholder="接受方的名称或地址" required />
        </div>

        <div class="form-group">
          <label for="amount">交易金额</label>
          <input id="amount" v-model.number="form.amount" type="number" step="0.01" min="0" placeholder="交易金额" required />
        </div>
      </template>

      <div class="form-group">
        <label for="message">信息 (可选)</label>
        <textarea id="message" v-model="form.message" :placeholder="transactionType === 'message' ? '输入信息' : '在你的交易中附加信息'"
          :required="transactionType === 'message'" rows="3"></textarea>
      </div>

      <button type="submit" :disabled="loading" class="submit-button">
        {{ loading ? '上链中...' : transactionType === 'transfer' ? '上传交易' : '上传信息' }}
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

.transaction-type-selector {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.transaction-type-selector button {
  flex: 1;
  padding: 10px;
  border: 2px solid #e9ecef;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.transaction-type-selector button.active {
  border-color: #4CAF50;
  background: #4CAF50;
  color: white;
}

.transaction-form {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  color: #666;
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

.submit-button {
  width: 100%;
  padding: 12px;
  background: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: background 0.3s ease;
}

.submit-button:disabled {
  background: #cccccc;
  cursor: not-allowed;
}

.error-message {
  margin-top: 15px;
  padding: 10px;
  background: #ffe6e6;
  color: #dc3545;
  border-radius: 4px;
}

.success-message {
  margin-top: 15px;
  padding: 10px;
  background: #e6ffe6;
  color: #28a745;
  border-radius: 4px;
}
</style>