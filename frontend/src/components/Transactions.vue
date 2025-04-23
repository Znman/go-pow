<template>
  <div>
    <h1>New Transaction</h1>
    <form @submit.prevent="submitTransaction">
      <input v-model="sender" placeholder="Sender" required />
      <input v-model="recipient" placeholder="Recipient" required />
      <input v-model="amount" type="number" placeholder="Amount" required />
      <button type="submit">Submit</button>
    </form>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return { sender: '', recipient: '', amount: 0 };
  },
  methods: {
    async submitTransaction() {
      await axios.post('http://localhost:8000/transactions/new', {
        sender: this.sender,
        recipient: this.recipient,
        amount: this.amount,
      });
      alert('Transaction submitted!');
    },
  },
};
</script>