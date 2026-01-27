<script setup lang="ts">
import { Search } from '@element-plus/icons-vue'

interface Props {
  searchPlaceholder?: string
  searchValue?: string
}

const props = withDefaults(defineProps<Props>(), {
  searchPlaceholder: '搜索',
  searchValue: '',
})

const emit = defineEmits<{
  'update:searchValue': [value: string]
  search: [value: string]
}>()

const handleSearch = (value: string) => {
  emit('update:searchValue', value)
  emit('search', value)
}
</script>

<template>
  <div class="filter-bar">
    <el-input
      :model-value="searchValue"
      :placeholder="searchPlaceholder"
      :prefix-icon="Search"
      style="width: 300px"
      clearable
      @update:model-value="handleSearch"
    />
    <slot name="filters" />
    <slot name="actions" />
  </div>
</template>

<style scoped>
.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  align-items: center;
}
</style>
