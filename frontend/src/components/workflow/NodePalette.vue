<script setup lang="ts">
import { NODE_TYPE_CONFIGS, type NodeType } from '../../types/workflow'

const emit = defineEmits<{
  (e: 'dragstart', type: NodeType, event: DragEvent): void
}>()

const nodeTypes = Object.entries(NODE_TYPE_CONFIGS) as [NodeType, typeof NODE_TYPE_CONFIGS[NodeType]][]

function onDragStart(type: NodeType, event: DragEvent) {
  if (event.dataTransfer) {
    event.dataTransfer.setData('application/vueflow-nodetype', type)
    event.dataTransfer.effectAllowed = 'move'
  }
  emit('dragstart', type, event)
}
</script>

<template>
  <div class="node-palette">
    <h3 class="palette-title">Node Palette</h3>
    <p class="palette-hint">Drag nodes to canvas</p>
    
    <div class="palette-nodes">
      <div
        v-for="[type, config] in nodeTypes"
        :key="type"
        class="palette-node"
        :style="{ borderColor: config.color }"
        draggable="true"
        @dragstart="onDragStart(type, $event)"
      >
        <div class="node-icon" :style="{ backgroundColor: config.color }">
          {{ config.icon }}
        </div>
        <div class="node-info">
          <div class="node-name">{{ config.label }}</div>
          <div class="node-type">{{ type }}</div>
        </div>
      </div>
    </div>
    
    <div class="palette-footer">
      <div class="legend">
        <div class="legend-item">
          <span class="legend-dot success"></span>
          <span>Success</span>
        </div>
        <div class="legend-item">
          <span class="legend-dot failure"></span>
          <span>Failure</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.node-palette {
  width: 200px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  padding: 16px;
  display: flex;
  flex-direction: column;
  height: fit-content;
}

.palette-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 4px 0;
}

.palette-hint {
  font-size: 12px;
  color: #6b7280;
  margin: 0 0 16px 0;
}

.palette-nodes {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.palette-node {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border: 2px solid;
  border-radius: 8px;
  cursor: grab;
  transition: all 0.2s;
  background: white;
}

.palette-node:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.palette-node:active {
  cursor: grabbing;
}

.node-icon {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
}

.node-info {
  flex: 1;
  min-width: 0;
}

.node-name {
  font-size: 12px;
  font-weight: 600;
  color: #1f2937;
}

.node-type {
  font-size: 10px;
  color: #9ca3af;
  font-family: monospace;
}

.palette-footer {
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px solid #e5e7eb;
}

.legend {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11px;
  color: #6b7280;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.legend-dot.success {
  background: #10b981;
}

.legend-dot.failure {
  background: #ef4444;
}
</style>
