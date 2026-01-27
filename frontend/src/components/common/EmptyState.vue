<script setup lang="ts">
interface Props {
  icon?: any
  title?: string
  description?: string
  showAction?: boolean
  actionText?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: '暂无数据',
  description: '',
  showAction: false,
  actionText: '创建',
})

const emit = defineEmits<{
  action: []
}>()
</script>

<template>
  <div class="empty-state">
    <el-icon v-if="icon" class="empty-icon" :size="64">
      <component :is="icon" />
    </el-icon>
    <h3 class="empty-title">{{ title }}</h3>
    <p v-if="description" class="empty-description">{{ description }}</p>
    <el-button v-if="showAction" type="primary" @click="emit('action')">
      {{ actionText }}
    </el-button>
    <slot />
  </div>
</template>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.empty-icon {
  color: #c0c4cc;
  margin-bottom: 16px;
}

.empty-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin: 0 0 8px 0;
}

.empty-description {
  font-size: 14px;
  color: #909399;
  margin: 0 0 20px 0;
}
</style>
