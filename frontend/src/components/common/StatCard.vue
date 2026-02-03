<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  title: string
  value: string | number
  icon?: string
  trend?: number // 增长趋势百分比
  color?: 'primary' | 'success' | 'warning' | 'danger' | 'info'
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  color: 'primary',
  loading: false
})

const colorClass = computed(() => {
  const colorMap = {
    primary: 'stat-card-primary',
    success: 'stat-card-success',
    warning: 'stat-card-warning',
    danger: 'stat-card-danger',
    info: 'stat-card-info'
  }
  return colorMap[props.color]
})

const trendClass = computed(() => {
  if (!props.trend) return ''
  return props.trend > 0 ? 'trend-up' : 'trend-down'
})

const trendIcon = computed(() => {
  if (!props.trend) return ''
  return props.trend > 0 ? '↑' : '↓'
})
</script>

<template>
  <div class="stat-card" :class="colorClass">
    <el-skeleton :loading="loading" animated>
      <template #template>
        <div class="stat-card-content">
          <el-skeleton-item variant="text" style="width: 60%" />
          <el-skeleton-item variant="h1" style="width: 80%; margin-top: 16px" />
        </div>
      </template>

      <template #default>
        <div class="stat-card-content">
          <div class="stat-card-header">
            <span class="stat-card-title">{{ title }}</span>
            <span v-if="icon" class="stat-card-icon">{{ icon }}</span>
          </div>

          <div class="stat-card-body">
            <div class="stat-card-value">{{ value }}</div>
            <div v-if="trend !== undefined" class="stat-card-trend" :class="trendClass">
              <span class="trend-icon">{{ trendIcon }}</span>
              <span class="trend-value">{{ Math.abs(trend) }}%</span>
            </div>
          </div>
        </div>
      </template>
    </el-skeleton>
  </div>
</template>

<style scoped>
.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
  border-left: 4px solid;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.stat-card-primary {
  border-left-color: #409eff;
}

.stat-card-success {
  border-left-color: #67c23a;
}

.stat-card-warning {
  border-left-color: #e6a23c;
}

.stat-card-danger {
  border-left-color: #f56c6c;
}

.stat-card-info {
  border-left-color: #909399;
}

.stat-card-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stat-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-card-title {
  font-size: 14px;
  color: #909399;
  font-weight: 500;
}

.stat-card-icon {
  font-size: 24px;
  opacity: 0.6;
}

.stat-card-body {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
}

.stat-card-value {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
}

.stat-card-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 4px;
}

.trend-up {
  color: #67c23a;
  background: #f0f9ff;
}

.trend-down {
  color: #f56c6c;
  background: #fef0f0;
}

.trend-icon {
  font-size: 14px;
}
</style>
