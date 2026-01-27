<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/common/PageHeader.vue'
import HostSelector from '@/components/host/HostSelector.vue'
import type { Host } from '@/api/host/types'

const router = useRouter()

// 处理主机选择
const handleHostSelected = (host: Host) => {
  ElMessage.success(`已选择主机: ${host.hostname}`)
}
</script>

<template>
  <div class="host-selection-view">
    <PageHeader title="云主机" />

    <HostSelector @select="handleHostSelected" />
  </div>
</template>

<style scoped>
.host-selection-view {
  padding: 24px;
}
</style>
