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
    </div>
  </div>
</template>

<style scoped>
.navigation-container {
  text-align: center;
  margin-bottom: 2rem;
  padding: 1rem;
  background-color: var(--color-background-soft);
  border-radius: 8px;
}

.navigation-container h1 {
  margin-bottom: 1.5rem;
  color: var(--color-heading);
}

.nav-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
}

.nav-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text);
  background-color: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.nav-button:hover {
  background-color: var(--color-background-mute);
  transform: translateY(-2px);
}

.button-icon {
  font-size: 1.2rem;
}

@media (max-width: 768px) {
  .nav-buttons {
    flex-direction: column;
  }

  .nav-button {
    width: 100%;
    justify-content: center;
  }
}
</style>