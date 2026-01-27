<script setup lang="ts">
import { ref, watch } from 'vue'
import type { HostFilterParams } from '@/api/host/types'
import { GPU_MODELS, GPU_COUNTS, REGIONS } from '@/config/hostConfig'

interface Props {
  modelValue: HostFilterParams
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: HostFilterParams]
  refresh: []
}>()

// 本地过滤器状态
const localFilters = ref<HostFilterParams>({ ...props.modelValue })

// 监听本地过滤器变化，发出更新事件
watch(localFilters, (newValue) => {
  emit('update:modelValue', { ...newValue })
}, { deep: true })

// 清除所有过滤器
const handleClearFilters = () => {
  localFilters.value = {
    region: '',
    gpu_count: '',
    gpu_model: '',
    keyword: ''
  }
}

// 刷新
const handleRefresh = () => {
  emit('refresh')
}
</script>

<template>
  <div class="host-filter-bar">
    <div class="filter-row">
      <el-input
        v-model="localFilters.keyword"
        placeholder="搜索主机名或IP地址"
        clearable
        class="search-input"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>

      <el-select
        v-model="localFilters.region"
        placeholder="选择地区"
        clearable
        class="filter-select"
      >
        <el-option
          v-for="region in REGIONS"
          :key="region.value"
          :label="region.label"
          :value="region.value"
        />
      </el-select>

      <el-select
        v-model="localFilters.gpu_count"
        placeholder="GPU数量"
        clearable
        class="filter-select"
      >
        <el-option
          v-for="count in GPU_COUNTS"
          :key="count.value"
          :label="count.label"
          :value="count.value"
        />
      </el-select>

      <el-select
        v-model="localFilters.gpu_model"
        placeholder="GPU型号"
        clearable
        class="filter-select"
      >
        <el-option
          v-for="model in GPU_MODELS"
          :key="model.value"
          :label="model.label"
          :value="model.value"
        />
      </el-select>

      <el-button @click="handleClearFilters">
        <el-icon><Delete /></el-icon>
        清除筛选
      </el-button>

      <el-button type="primary" @click="handleRefresh">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.host-filter-bar {
  background: white;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.filter-row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
}

.search-input {
  width: 280px;
}

.filter-select {
  width: 160px;
}

@media (max-width: 768px) {
  .search-input,
  .filter-select {
    width: 100%;
  }
}
</style>
