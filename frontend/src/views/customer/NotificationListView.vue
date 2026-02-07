<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Check, Bell } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { getNotifications, type Notification } from '@/api/customer'
import { useNotificationStore } from '@/stores/notification'

const notificationStore = useNotificationStore()
const loading = ref(false)
const notifications = ref<Notification[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const onlyUnread = ref(false)

const loadNotifications = async () => {
  loading.value = true
  try {
    const res = await getNotifications({
      page: page.value,
      pageSize: pageSize.value,
      unread: onlyUnread.value || undefined,
    })
    notifications.value = res.data.list || []
    total.value = res.data.total || 0
  } catch {
    // 静默
  } finally {
    loading.value = false
  }
}

const handleMarkRead = async (item: Notification) => {
  if (item.is_read) return
  try {
    await notificationStore.markRead(item.id)
    item.is_read = true
    ElMessage.success('已标记为已读')
  } catch {
    ElMessage.error('操作失败')
  }
}

const handleMarkAllRead = async () => {
  try {
    await notificationStore.markAllRead()
    notifications.value.forEach(n => { n.is_read = true })
    ElMessage.success('已全部标记为已读')
  } catch {
    ElMessage.error('操作失败')
  }
}

const handleFilterChange = () => {
  page.value = 1
  loadNotifications()
}

const handlePageChange = (val: number) => {
  page.value = val
  loadNotifications()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  page.value = 1
  loadNotifications()
}

const formatDate = (value?: string) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN')
}

const levelTagType = (level: string) => {
  const map: Record<string, string> = {
    info: 'info',
    warning: 'warning',
    error: 'danger',
  }
  return map[level] || 'info'
}

const typeLabel = (type: string) => {
  const map: Record<string, string> = {
    task: '任务',
    alert: '告警',
    machine: '机器',
    system: '系统',
  }
  return map[type] || type
}

onMounted(() => {
  loadNotifications()
})
</script>

<template>
  <div class="notification-list-view">
    <PageHeader title="消息通知">
      <template #actions>
        <el-button :icon="Check" @click="handleMarkAllRead">全部已读</el-button>
      </template>
    </PageHeader>

    <el-card class="filter-card">
      <el-checkbox v-model="onlyUnread" @change="handleFilterChange">仅显示未读</el-checkbox>
    </el-card>

    <el-card>
      <el-table :data="notifications" v-loading="loading" style="width: 100%">
        <el-table-column label="" width="40">
          <template #default="{ row }">
            <el-icon v-if="!row.is_read" color="#409eff" :size="8"><Bell /></el-icon>
          </template>
        </el-table-column>
        <el-table-column label="类型" width="90">
          <template #default="{ row }">
            <el-tag size="small" :type="levelTagType(row.level)">
              {{ typeLabel(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <span :class="{ 'unread-title': !row.is_read }">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="内容" min-width="250" show-overflow-tooltip />
        <el-table-column label="时间" width="175">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button
              v-if="!row.is_read"
              link
              type="primary"
              size="small"
              @click="handleMarkRead(row)"
            >
              标记已读
            </el-button>
            <span v-else class="read-label">已读</span>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </el-card>
  </div>
</template>

<style scoped>
.notification-list-view {
  padding: 24px;
}

.filter-card {
  margin-bottom: 16px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.unread-title {
  font-weight: 600;
  color: #303133;
}

.read-label {
  font-size: 13px;
  color: #c0c4cc;
}
</style>
