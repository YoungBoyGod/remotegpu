<script setup lang="ts">
interface Props {
  title: string
  extra?: string
}

defineProps<Props>()
</script>

<template>
  <el-card class="detail-card" shadow="never">
    <template #header>
      <div class="card-header">
        <span class="card-title">{{ title }}</span>
        <span v-if="extra" class="card-extra">{{ extra }}</span>
        <slot name="extra" />
      </div>
    </template>
    <slot />
  </el-card>
</template>

<style scoped>
.detail-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.card-extra {
  font-size: 14px;
  color: #909399;
}
</style>
