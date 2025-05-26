<template>
  <section class="consensus-visualizer">
    <h2>Consensus Process Visualization</h2>
    <p>
      Three nodes are competing to mine the next block. The first node to find a valid proof wins!
    </p>
    <button @click="startCompetition" :disabled="competitionActive">Start Competition</button>
    <div class="nodes">
      <div v-for="(node, idx) in nodes" :key="node.id" :class="['node-panel', { winner: winnerNode === node.id }]">
        <h3>Node {{ idx + 1 }}</h3>
        <div v-if="node.progress">
          <p>Attempt: {{ node.progress.attempt }}</p>
          <p>Proof: {{ node.progress.proof }}</p>
          <p>Hash: <span class="hash">{{ node.progress.hash }}</span></p>
          <p>Status: <span :class="{ found: node.progress.found }">
              {{ node.progress.message }}
            </span></p>
        </div>
        <div v-else>
          <p>Waiting for mining to start...</p>
        </div>
      </div>
    </div>
    <div v-if="winnerNode !== null" class="winner-banner">
      🏆 Node {{ winnerNode + 1 }} wins! Block added to chain.
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue';

interface ConsensusProgress {
  attempt: number,
  proof: number,
  hash: string,
  found: boolean,
  blockIndex: number,
  message: string,
}

interface NodeState {
  id: number,
  eventSource: EventSource | null,
  progress: ConsensusProgress | null,
}

const nodes = ref<NodeState[]>([
  { id: 0, eventSource: null, progress: null },
  { id: 1, eventSource: null, progress: null },
  { id: 2, eventSource: null, progress: null },
]);
const winnerNode = ref<number | null>(null);
const competitionActive = ref(false);

function startCompetition() {
  winnerNode.value = null;
  competitionActive.value = true;
  nodes.value.forEach((node, idx) => {
    // Clean up any previous EventSource
    if (node.eventSource) {
      node.eventSource.close();
    }
    node.progress = null;
    // Each node opens its own SSE connection (simulate three independent miners)
    const es = new EventSource(`${'http://localhost:8000'}/mine/stream`);
    node.eventSource = es;

    es.onmessage = (event: MessageEvent) => {
      const progress: ConsensusProgress = JSON.parse(event.data);
      node.progress = progress;
      // If this node finds the block, declare as winner
      if (progress.found && winnerNode.value === null) {
        winnerNode.value = idx;
        competitionActive.value = false;
        // Stop the mining for the other nodes
        nodes.value.forEach((otherNode, otherIdx) => {
          if (otherNode.eventSource && otherIdx !== idx) {
            otherNode.eventSource.close();
          }
        });
      }
    };

    es.onerror = () => {
      es.close();
      node.eventSource = null;
      node.progress = node.progress || { attempt: 0, proof: 0, hash: '', found: false, blockIndex: 0, message: 'Connection error' };
    };
  });
}
</script>

<style scoped>
.consensus-visualizer {
  border: 1px solid #d0d0d0;
  border-radius: 8px;
  padding: 1.5em;
  margin: 2em auto;
  max-width: 900px;
  background: #f8f9fb;
}

.nodes {
  display: flex;
  justify-content: space-around;
  margin-top: 2em;
  gap: 1em;
}

.node-panel {
  flex: 1 1 0;
  background: #fff;
  border: 2px solid #b9b9b9;
  border-radius: 8px;
  padding: 1em;
  text-align: center;
  box-shadow: 0 2px 8px #0001;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.node-panel.winner {
  border-color: #41c300;
  background: #e9ffe2;
  box-shadow: 0 4px 16px #41c30044;
}

.hash {
  font-family: monospace;
  font-size: smaller;
  word-break: break-all;
}

.found {
  color: #41c300;
  font-weight: bold;
}

.winner-banner {
  margin-top: 2em;
  font-size: 1.4em;
  color: #41c300;
  text-align: center;
  font-weight: bold;
  background: #eaffd8;
  border-radius: 6px;
  padding: 0.6em 0;
}

button[disabled] {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>