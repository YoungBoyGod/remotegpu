<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  status: string
  type?: 'success' | 'info' | 'warning' | 'danger' | 'default'
}

const props = defineProps<Props>()

// 自动根据状态文本确定标签类型
const tagType = computed(() => {
  if (props.type) return props.type

  const status = props.status.toLowerCase()

  // 成功状态
  if (['running', 'active', '运行中', '活跃'].includes(status)) {
    return 'info'
  }

  // 停止状态
  if (['stopped', 'inactive', '已停止', '未激活'].includes(status)) {
    return 'default'
  }

  // 完成状态
  if (['completed', 'success', 'ready', '已完成', '成功', '就绪'].includes(status)) {
    return 'success'
  }

  // 失败状态
  if (['failed', 'error', '失败', '错误'].includes(status)) {
    return 'danger'
  }

  // 处理中状态
  if (['pending', 'processing', 'building', '等待中', '处理中', '构建中'].includes(status)) {
    return 'warning'
  }

  return 'default'
})
</script>

<template>
  <el-tag :type="tagType" size="small">
    {{ status }}
  </el-tag>
</template>
