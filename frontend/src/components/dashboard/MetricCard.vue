<script setup lang="ts">
import { computed } from 'vue'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

interface Props {
  title: string
  value: number
  unit: string
  icon: string
  color: string
}

const props = defineProps<Props>()

// 动态获取图标组件
const IconComponent = computed(() => {
  return (ElementPlusIconsVue as any)[props.icon]
})

// 格式化数值
const formattedValue = computed(() => {
  if (props.value >= 1000) {
    return (props.value / 1000).toFixed(1) + 'K'
  }
  return props.value.toFixed(1)
})
</script>

<template>
  <div class="metric-card">
    <div class="card-icon" :style="{ backgroundColor: color }">
      <component :is="IconComponent" class="icon" />
    </div>
    <div class="card-content">
      <div class="card-title">{{ title }}</div>
      <div class="card-value">
        {{ formattedValue }}
        <span class="unit">{{ unit }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.metric-card {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 16px;
  transition: transform 0.2s, box-shadow 0.2s;
  cursor: pointer;
}

.metric-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.card-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.icon {
  width: 28px;
  height: 28px;
  color: white;
}

.card-content {
  flex: 1;
}

.card-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.card-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.unit {
  font-size: 14px;
  font-weight: 400;
  color: #909399;
  margin-left: 4px;
}
</style>
