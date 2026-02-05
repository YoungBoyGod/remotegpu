<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { getMachineEnrollments } from '@/api/customer'
import type { MachineEnrollment } from '@/api/customer'
import type { PageRequest } from '@/types/common'
import DataTable from '@/components/common/DataTable.vue'

const router = useRouter()
const loading = ref(false)
const enrollments = ref<MachineEnrollment[]>([])
const total = ref(0)
const pageRequest = ref<PageRequest>({
  page: 1,
  pageSize: 10
})
const refreshTimer = ref<number | null>(null)

const statusMap: Record<string, { text: string; type: string }> = {
  pending: { text: '采集中', type: 'warning' },
  success: { text: '已完成', type: 'success' },
  failed: { text: '失败', type: 'danger' }
}

const loadEnrollments = async (silent = false) => {
  try {
    if (!silent) {
      loading.value = true
    }
    const response = await getMachineEnrollments(pageRequest.value)
    enrollments.value = response.data.list
    total.value = response.data.total
  } catch (error) {
    console.error('加载添加任务失败:', error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  pageRequest.value.page = page
  loadEnrollments()
}

const handleSizeChange = (size: number) => {
  pageRequest.value.pageSize = size
  loadEnrollments()
}

const handleRefresh = () => {
  loadEnrollments()
}

const navigateToAdd = () => {
  router.push('/customer/machines/enroll')
}

const navigateToList = () => {
  router.push('/customer/machines/list')
}

const getStatusInfo = (status: string) => {
  return statusMap[status] || { text: status, type: 'info' }
}

onMounted(() => {
  loadEnrollments()
  refreshTimer.value = window.setInterval(() => loadEnrollments(true), 5000)
})

onUnmounted(() => {
  if (refreshTimer.value) {
    window.clearInterval(refreshTimer.value)
  }
})
</script>

<template>
  <div class="machine-enrollment-list">
    <div class="page-header">
      <h2 class="page-title">添加进度</h2>
      <div class="page-actions">
        <el-button @click="handleRefresh">刷新</el-button>
        <el-button type="primary" @click="navigateToAdd">添加机器</el-button>
      </div>
    </div>

    <el-alert
      type="info"
      show-icon
      :closable="false"
      title="页面每 5 秒自动刷新一次采集状态。"
      class="tip"
    />

    <DataTable
      :data="enrollments"
      :total="total"
      :loading="loading"
      :current-page="pageRequest.page"
      :page-size="pageRequest.pageSize"
      :show-pagination="true"
      @page-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <el-table-column prop="name" label="机器名称" min-width="140" />
      <el-table-column prop="address" label="连接地址" min-width="160" />
      <el-table-column prop="region" label="区域" width="120" />
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="getStatusInfo(row.status).type">
            {{ getStatusInfo(row.status).text }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="host_id" label="机器ID" min-width="140" />
      <el-table-column prop="error_message" label="失败原因" min-width="200" />
      <el-table-column prop="created_at" label="提交时间" width="180" />
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status === 'success'" link type="primary" size="small" @click="navigateToList">
            查看机器
          </el-button>
          <el-button v-else link type="primary" size="small" @click="navigateToAdd">
            再次添加
          </el-button>
        </template>
      </el-table-column>
    </DataTable>
  </div>
</template>

<style scoped>
.machine-enrollment-list {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.page-actions {
  display: flex;
  gap: 12px;
}

.tip {
  margin-bottom: 16px;
}
</style>
