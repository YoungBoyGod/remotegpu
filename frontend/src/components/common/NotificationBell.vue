<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Bell } from '@element-plus/icons-vue'
import { useNotificationStore } from '@/stores/notification'

const router = useRouter()
const notificationStore = useNotificationStore()

const handleClick = () => {
  router.push('/customer/notifications')
}

onMounted(() => {
  notificationStore.connectSSE()
  notificationStore.fetchUnreadCount()
})

onUnmounted(() => {
  notificationStore.disconnectSSE()
})
</script>

<template>
  <el-badge :value="notificationStore.unreadCount" :hidden="notificationStore.unreadCount === 0" :max="99">
    <el-button :icon="Bell" circle size="small" @click="handleClick" />
  </el-badge>
</template>
