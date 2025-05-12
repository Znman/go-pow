<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const navigateTo = (path: string) => {
  router.push(path)
}
</script>

<template>
  <div class="welcome-container">
    <div class="welcome-header">
      <h1>Blockchain PoW Explorer</h1>
      <p class="subtitle">A Proof of Work Blockchain Demonstration System</p>
    </div>

    <div class="features-grid">
      <div class="feature-card" @click="navigateTo('/explorer')">
        <div class="icon">🔍</div>
        <h2>Blockchain Explorer</h2>
        <p>View detailed blockchain information, including block contents, transaction records, and mining proofs.</p>
        <button class="action-button">Open Explorer</button>
      </div>

      <div class="feature-card" @click="navigateTo('/blockchain')">
        <div class="icon">⛏️</div>
        <h2>Block Mining</h2>
        <p>Experience the blockchain mining process and understand how Proof of Work mechanisms operate.</p>
        <button class="action-button">Start Mining</button>
      </div>

      <div class="feature-card" @click="navigateTo('/transactions')">
        <div class="icon">💸</div>
        <h2>Transaction Management</h2>
        <p>Create new transactions and experience the complete blockchain transaction process.</p>
        <button class="action-button">New Transaction</button>
      </div>
    </div>

    <div class="info-section">
      <h3>About This System</h3>
      <p>This system is a blockchain demonstration project based on Proof of Work (PoW) consensus mechanism, designed to
        illustrate the fundamental concepts and operational principles of blockchain technology. Through an interactive
        interface, users can experience core blockchain functionalities, including:</p>
      <ul>
        <li>Block creation and validation processes</li>
        <li>Transaction initiation and confirmation mechanisms</li>
        <li>Proof of Work mining operations</li>
        <li>Blockchain data transparency and immutability</li>
      </ul>

      <div class="research-notes">
        <h4>Research Implementation Notes</h4>
        <p>This implementation focuses on the following key aspects of blockchain technology:</p>
        <ul>
          <li><strong>Consensus Mechanism:</strong> Demonstrates a basic PoW implementation with adjustable difficulty
          </li>
          <li><strong>Data Structure:</strong> Implements linked blocks using cryptographic hashes</li>
          <li><strong>Transaction Management:</strong> Shows UTXO-style transaction handling</li>
          <li><strong>Network Transparency:</strong> Provides real-time visibility of blockchain state</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
.welcome-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.welcome-header {
  text-align: center;
  margin-bottom: 60px;
}

.welcome-header h1 {
  font-size: 2.5em;
  color: #2c3e50;
  margin-bottom: 16px;
}

.subtitle {
  font-size: 1.2em;
  color: #666;
  margin-top: 0;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 30px;
  margin-bottom: 60px;
}

.feature-card {
  background: white;
  border-radius: 12px;
  padding: 30px;
  text-align: center;
  transition: transform 0.3s, box-shadow 0.3s;
  cursor: pointer;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.1);
}

.icon {
  font-size: 3em;
  margin-bottom: 20px;
}

.feature-card h2 {
  color: #2c3e50;
  margin-bottom: 15px;
  font-size: 1.5em;
}

.feature-card p {
  color: #666;
  margin-bottom: 20px;
  line-height: 1.6;
}

.action-button {
  background-color: #2196F3;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 6px;
  font-size: 1em;
  cursor: pointer;
  transition: background-color 0.3s;
}

.action-button:hover {
  background-color: #1976D2;
}

.info-section {
  background: #f8f9fa;
  padding: 40px;
  border-radius: 12px;
  margin-top: 40px;
}

.info-section h3 {
  color: #2c3e50;
  margin-bottom: 20px;
  font-size: 1.8em;
}

.info-section h4 {
  color: #2c3e50;
  margin: 30px 0 15px;
  font-size: 1.4em;
}

.info-section p {
  color: #666;
  line-height: 1.6;
  margin-bottom: 20px;
}

.info-section ul {
  color: #666;
  padding-left: 20px;
  margin-bottom: 20px;
}

.info-section li {
  margin-bottom: 10px;
  line-height: 1.6;
}

.research-notes {
  margin-top: 30px;
  padding-top: 30px;
  border-top: 1px solid #e0e0e0;
}

.research-notes strong {
  color: #2c3e50;
}

@media (max-width: 768px) {
  .welcome-container {
    padding: 20px;
  }

  .welcome-header h1 {
    font-size: 2em;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }

  .info-section {
    padding: 20px;
  }
}
</style>