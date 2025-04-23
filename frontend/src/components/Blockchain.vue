<template>
  <div>
    <h1>Blockchain</h1>
    <button @click="fetchChain">Load Blockchain</button>
    <ul>
      <li v-for="block in chain" :key="block.Index">
        Block {{ block.Index }}: {{ block.PreviousHash }}
      </li>
    </ul>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return { chain: [] };
  },
  methods: {
    async fetchChain() {
      const response = await axios.get('http://localhost:8000/chain');
      this.chain = response.data.Chain;
    },
  },
};
</script>