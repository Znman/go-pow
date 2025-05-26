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
      <h1>基于工作量认证机制的共识系统</h1>
    </div>

    <div class="features-grid">
      <div class="feature-card" @click="navigateTo('/explorer')">
        <div class="icon">🔍</div>
        <h2>区块链浏览器</h2>
        <p>查看详细的区块链信息，包括区块内容、交易记录和挖矿证明。</p>
        <button class="action-button">打开区块链浏览器</button>
      </div>

      <div class="feature-card" @click="navigateTo('/blockchain')">
        <div class="icon">⛏️</div>
        <h2>挖矿</h2>
        <p>展示区块链挖矿过程，了解工作量证明机制。</p>
        <button class="action-button">开始挖矿</button>
      </div>

      <div class="feature-card" @click="navigateTo('/transactions')">
        <div class="icon">💸</div>
        <h2>交易/信息管理</h2>
        <p>创建新交易，体验消息上链过程。</p>
        <button class="action-button">交易/信息管理</button>
      </div>
    </div>

    <div class="info-section">
      <h3>关于本系统</h3>
      <p>该系统是一个基于工作量证明（PoW）共识机制的区块链演示项目，旨在展示区块链技术的基本概念和操作原理。通过一个互动界面，用户可以体验区块链的核心功能，包括：</p>
      <ul>
        <li>区块创建和验证过程</li>
        <li>交易发起和确认机制</li>
        <li>工作量证明挖矿操作</li>
        <li>区块链数据透明和不可篡改性</li>
      </ul>

      <!-- <div class="research-notes">
        <h4>研究细节：</h4>
        <p>本系统重点关注区块链技术的以下关键方面：</p>
        <ul>
          <li><strong></strong> Demonstrates a basic PoW implementation with adjustable difficulty
          </li>
          <li><strong>Data Structure:</strong> Implements linked blocks using cryptographic hashes</li>
          <li><strong>Transaction Management:</strong> Shows UTXO-style transaction handling</li>
          <li><strong>Network Transparency:</strong> Provides real-time visibility of blockchain state</li>
        </ul>
      </div> -->
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