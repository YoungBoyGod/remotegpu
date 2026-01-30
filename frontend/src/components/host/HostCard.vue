<script setup lang="ts">
import { computed } from 'vue'
import type { Host } from '@/api/host/types'
import StatusTag from '@/components/common/StatusTag.vue'

interface Props {
  host: Host
  selected?: boolean
  selectionMode?: 'single' | 'multiple'
}

const props = withDefaults(defineProps<Props>(), {
  selected: false,
  selectionMode: 'single'
})

const emit = defineEmits<{
  select: [host: Host]
}>()

const formatMemory = (value?: number | null) => {
  if (!value && value !== 0) return '-'
  const gb = value / 1024
  if (gb >= 1) {
    return `${gb.toFixed(1)} GB`
  }
  return `${value} MB`
}

const formatStorage = (value?: number | null) => {
  if (!value && value !== 0) return '-'
  const gb = value / 1024
  if (gb >= 1) {
    return `${gb.toFixed(1)} GB`
  }
  return `${value} MB`
}

const hostTitle = computed(() => {
  return props.host.name || props.host.hostname || props.host.id
})

// 处理选择
const handleSelect = () => {
  emit('select', props.host)
}
</script>

<template>
  <el-card
    class="host-card"
    :class="{ 'is-selected': selected }"
    shadow="hover"
    @click="handleSelect"
  >
    <div class="card-header">
      <div class="host-title">
        <span class="host-name">{{ hostTitle }}</span>
        <span class="host-ip">{{ host.ip_address }}</span>
      </div>
      <div class="header-right">
        <StatusTag :status="host.status" />
      </div>
    </div>

    <div class="card-body">
      <div class="spec-item">
        <span class="spec-label">系统</span>
        <span>{{ host.os_type }} {{ host.os_version || '' }}</span>
      </div>
      <div class="spec-item">
        <span class="spec-label">CPU</span>
        <span>{{ host.used_cpu }} / {{ host.total_cpu }} 核</span>
      </div>
      <div class="spec-item">
        <span class="spec-label">内存</span>
        <span>{{ formatMemory(host.used_memory) }} / {{ formatMemory(host.total_memory) }}</span>
      </div>
      <div class="spec-item">
        <span class="spec-label">GPU</span>
        <span>{{ host.used_gpu }} / {{ host.total_gpu }} 张</span>
      </div>
      <div class="spec-item">
        <span class="spec-label">磁盘</span>
        <span>{{ formatStorage(host.used_disk) }} / {{ formatStorage(host.total_disk) }}</span>
      </div>
    </div>

    <div class="card-footer">
      <div class="health-status">健康状态：{{ host.health_status }}</div>
      <el-button type="primary" size="small" :disabled="host.status !== 'active'">
        选择此主机
      </el-button>
    </div>
  </el-card>
</template>

<style scoped>
.host-card {
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
}

.host-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.host-card.is-selected {
  border-color: #409EFF;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.3);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #EBEEF5;
}


.host-title {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.host-name {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.host-ip {
  font-size: 13px;
  color: #909399;
}

.header-right {
  display: flex;
  align-items: center;
}

.card-body {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.spec-item {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  font-size: 14px;
  color: #606266;
}

.spec-label {
  color: #909399;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #EBEEF5;
}

.health-status {
  font-size: 13px;
  color: #909399;
}

@media (max-width: 768px) {
  .card-body {
    grid-template-columns: 1fr;
  }
}
</style>
