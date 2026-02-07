import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getUnreadCount, markAllNotificationsRead, markNotificationRead } from '@/api/customer'
import { useAuthStore } from './auth'

export const useNotificationStore = defineStore('notification', () => {
  const unreadCount = ref(0)
  const sseConnection = ref<EventSource | null>(null)
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null

  // 获取未读数
  const fetchUnreadCount = async () => {
    try {
      const res = await getUnreadCount()
      unreadCount.value = res.data.count || 0
    } catch {
      // 静默失败
    }
  }

  // 标记单条已读
  const markRead = async (id: number) => {
    await markNotificationRead(id)
    if (unreadCount.value > 0) unreadCount.value--
  }

  // 标记全部已读
  const markAllRead = async () => {
    await markAllNotificationsRead()
    unreadCount.value = 0
  }

  // 建立 SSE 连接
  const connectSSE = () => {
    const authStore = useAuthStore()
    if (!authStore.accessToken || sseConnection.value) return

    const baseURL = import.meta.env.VITE_API_BASE_URL || '/api/v1'
    const url = `${baseURL}/customer/notifications/sse?token=${encodeURIComponent(authStore.accessToken)}`

    const es = new EventSource(url)
    sseConnection.value = es

    es.addEventListener('connected', () => {
      fetchUnreadCount()
    })

    es.addEventListener('notification', () => {
      unreadCount.value++
    })

    es.addEventListener('task_status', () => {
      // 任务状态变更，可由页面自行监听
    })

    es.addEventListener('machine_status', () => {
      // 机器状态变更，可由页面自行监听
    })

    es.onerror = () => {
      sseConnection.value?.close()
      sseConnection.value = null
      // 5 秒后重连
      if (!reconnectTimer) {
        reconnectTimer = setTimeout(() => {
          reconnectTimer = null
          if (authStore.isAuthenticated) connectSSE()
        }, 5000)
      }
    }
  }

  // 断开 SSE 连接
  const disconnectSSE = () => {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    if (sseConnection.value) {
      sseConnection.value.close()
      sseConnection.value = null
    }
  }

  return {
    unreadCount,
    sseConnection,
    fetchUnreadCount,
    markRead,
    markAllRead,
    connectSSE,
    disconnectSSE,
  }
})
