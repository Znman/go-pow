<script setup lang="ts">
import { useRouter } from 'vue-router'
import WelcomeItem from './WelcomeItem.vue'
import DocumentationIcon from './icons/IconDocumentation.vue'
import ToolingIcon from './icons/IconTooling.vue'
import EcosystemIcon from './icons/IconEcosystem.vue'
import CommunityIcon from './icons/IconCommunity.vue'
import SupportIcon from './icons/IconSupport.vue'

const router = useRouter()

const navigateTo = (route: string) => {
  router.push({ name: route })
}
</script>

<template>
  <!-- Navigation Section -->
  <div class="navigation-container">
    <h1>Welcome to Go-POW</h1>
    <div class="nav-buttons">
      <button @click="navigateTo('blockchain')" class="nav-button">
        <span class="button-icon">🔗</span>
        Blockchain
      </button>
      <button @click="navigateTo('transactions')" class="nav-button">
        <span class="button-icon">💸</span>
        Transactions
      </button>
      <!-- <button @click="navigateTo('mining')" class="nav-button">
        <span class="button-icon">⛏️</span>
        Mining
      </button> -->
    </div>
  </div>
</template>

<style scoped>
.navigation-container {
  text-align: center;
  padding: 20px;
}

.nav-buttons {
  display: flex;
  gap: 20px;
  justify-content: center;
  margin-top: 20px;
}

.nav-button {
  padding: 12px 24px;
  font-size: 16px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  display: flex;
  align-items: center;
}

.nav-button:hover {
  background-color: #45a049;
}

.button-icon {
  margin-right: 8px;
  font-size: 20px;
}
</style>