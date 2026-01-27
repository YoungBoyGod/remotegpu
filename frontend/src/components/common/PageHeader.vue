<script setup lang="ts">
interface Props {
  title: string
  subtitle?: string
}

defineProps<Props>()
</script>

<template>
  <div class="page-header">
    <div class="header-left">
      <h1>{{ title }}</h1>
      <p v-if="subtitle" class="subtitle">{{ subtitle }}</p>
    </div>
    <div class="header-right">
      <slot name="actions" />
    </div>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-left h1 {
  font-size: 24px;
  font-weight: 600;
  margin: 0;
}

.subtitle {
  font-size: 14px;
  color: #909399;
  margin: 4px 0 0 0;
}

.header-right {
  display: flex;
  gap: 12px;
}
</style>
