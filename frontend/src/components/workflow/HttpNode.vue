<script setup lang="ts">
import { Handle, Position } from '@vue-flow/core'
import type { WorkflowNodeData, HttpNodeConfig } from '../../types/workflow'

defineProps<{
  id: string
  data: WorkflowNodeData
  selected?: boolean
}>()
</script>

<template>
  <div
    class="workflow-node http-node"
    :class="{ selected }"
  >
    <Handle
      type="target"
      :position="Position.Left"
      class="handle-input"
    />
    
    <div class="node-header">
      <span class="node-icon">🌐</span>
      <span class="node-title">{{ data.label || 'HTTP Request' }}</span>
    </div>
    
    <div class="node-content">
      <div class="node-badge" :class="(data.config as HttpNodeConfig).method.toLowerCase()">
        {{ (data.config as HttpNodeConfig).method }}
      </div>
      <div class="node-url" :title="(data.config as HttpNodeConfig).url">
        {{ (data.config as HttpNodeConfig).url?.slice(0, 20) }}{{ ((data.config as HttpNodeConfig).url?.length || 0) > 20 ? '...' : '' }}
      </div>
    </div>
    
    <Handle
      id="success"
      type="source"
      :position="Position.Right"
      class="handle-success"
      :style="{ top: '35%' }"
    />
    <Handle
      id="failure"
      type="source"
      :position="Position.Right"
      class="handle-failure"
      :style="{ top: '65%' }"
    />
    
    <div class="handle-labels">
      <span class="handle-label success" style="top: 35%">✓</span>
      <span class="handle-label failure" style="top: 65%">✗</span>
    </div>
  </div>
</template>

<style scoped>
.workflow-node {
  background: white;
  border: 2px solid #3b82f6;
  border-radius: 8px;
  padding: 0;
  min-width: 160px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: relative;
}

.workflow-node.selected {
  border-color: #7c3aed;
  box-shadow: 0 0 0 2px rgba(124, 58, 237, 0.3);
}

.node-header {
  background: #3b82f6;
  color: white;
  padding: 8px 12px;
  border-radius: 6px 6px 0 0;
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 13px;
}

.node-icon {
  font-size: 16px;
}

.node-content {
  padding: 10px 12px;
}

.node-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  display: inline-block;
  margin-bottom: 6px;
}

.node-badge.get {
  background: #dbeafe;
  color: #1d4ed8;
}

.node-badge.post {
  background: #d1fae5;
  color: #065f46;
}

.node-badge.put {
  background: #fef3c7;
  color: #92400e;
}

.node-badge.delete {
  background: #fee2e2;
  color: #991b1b;
}

.node-badge.patch {
  background: #e0e7ff;
  color: #3730a3;
}

.node-url {
  font-size: 11px;
  color: #6b7280;
  font-family: monospace;
  word-break: break-all;
}

.handle-input {
  width: 12px !important;
  height: 12px !important;
  background: #6b7280 !important;
  border: 2px solid white !important;
}

.handle-success {
  width: 12px !important;
  height: 12px !important;
  background: #10b981 !important;
  border: 2px solid white !important;
}

.handle-failure {
  width: 12px !important;
  height: 12px !important;
  background: #ef4444 !important;
  border: 2px solid white !important;
}

.handle-labels {
  position: absolute;
  right: 16px;
  top: 0;
  bottom: 0;
  pointer-events: none;
}

.handle-label {
  position: absolute;
  font-size: 10px;
  transform: translateY(-50%);
}

.handle-label.success {
  color: #10b981;
}

.handle-label.failure {
  color: #ef4444;
}
</style>
