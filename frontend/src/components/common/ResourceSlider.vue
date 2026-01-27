<script setup lang="ts">
interface Props {
  label: string
  modelValue: number
  min?: number
  max?: number
  step?: number
  unit?: string
  showInput?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  min: 0,
  max: 100,
  step: 1,
  unit: '',
  showInput: true,
})

const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

const handleChange = (value: number | null) => {
  if (value !== null) {
    emit('update:modelValue', value)
  }
}
</script>

<template>
  <div class="resource-slider">
    <div class="slider-header">
      <span class="slider-label">{{ label }}</span>
      <span class="slider-value">{{ modelValue }} {{ unit }}</span>
    </div>
    <div class="slider-content">
      <el-slider
        :model-value="modelValue"
        :min="min"
        :max="max"
        :step="step"
        :show-input="showInput"
        @update:model-value="handleChange"
      />
    </div>
  </div>
</template>

<style scoped>
.resource-slider {
  margin-bottom: 24px;
}

.slider-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.slider-label {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.slider-value {
  font-size: 14px;
  color: #606266;
}

.slider-content {
  padding: 0 12px;
}
</style>
