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

// 格式化价格
const formattedPrice = computed(() => {
  return `¥${props.host.price_per_hour.toFixed(2)}/时`
})

// 价格颜色
const priceColor = computed(() => {
  const price = props.host.price_per_hour
  if (price < 2) return '#67C23A'
  if (price < 5) return '#409EFF'
  return '#E6A23C'
})

// 可用性状态映射
const availabilityStatusMap = {
  available: { text: '可用', type: 'success' },
  limited: { text: '库存紧张', type: 'warning' },
  unavailable: { text: '已售罄', type: 'danger' }
}

const availabilityStatus = computed(() => {
  return availabilityStatusMap[props.host.availability_status]
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
      <div class="gpu-info">
        <span class="gpu-model">{{ host.gpu_model }} / {{ host.gpu_memory }}GB</span>
      </div>
      <div class="header-right">
        <StatusTag :type="availabilityStatus.type" :text="availabilityStatus.text" />
        <span class="price" :style="{ color: priceColor }">{{ formattedPrice }}</span>
      </div>
    </div>

    <div class="card-body">
      <div class="spec-item">
        <el-icon><Location /></el-icon>
        <span>{{ host.region }}</span>
      </div>
      <div class="spec-item">
        <el-icon><Cpu /></el-icon>
        <span>{{ host.cpu_cores }}核</span>
      </div>
      <div class="spec-item">
        <el-icon><Memo /></el-icon>
        <span>{{ host.memory_total }}GB</span>
      </div>
      <div class="spec-item">
        <el-icon><Coin /></el-icon>
        <span>{{ host.disk_total }}GB</span>
      </div>
    </div>

    <div class="card-footer">
      <div class="cuda-version">
        <el-icon><Tools /></el-icon>
        <span>CUDA {{ host.cuda_version }}</span>
      </div>
      <el-button
        type="primary"
        size="small"
        :disabled="host.availability_status === 'unavailable'"
      >
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

.gpu-info {
  flex: 1;
}

.gpu-model {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.price {
  font-size: 20px;
  font-weight: 600;
}

.card-body {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.spec-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #606266;
}

.spec-item .el-icon {
  color: #909399;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #EBEEF5;
}

.cuda-version {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #909399;
}

@media (max-width: 768px) {
  .card-body {
    grid-template-columns: 1fr;
  }
}
</style>
