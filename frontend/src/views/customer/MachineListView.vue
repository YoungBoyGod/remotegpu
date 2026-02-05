<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getMyMachines, getMachineMonitoring } from '@/api/customer'
import type { Machine } from '@/types/machine'
import type { PageRequest } from '@/types/common'
import DataTable from '@/components/common/DataTable.vue'

const loading = ref(false)
const machines = ref<Machine[]>([])
const total = ref(0)
const router = useRouter()
const pageRequest = ref<PageRequest>({
  page: 1,
  pageSize: 10
})

const loadMachines = async () => {
  try {
    loading.value = true
    const response = await getMyMachines(pageRequest.value)
    machines.value = response.data.list
    total.value = response.data.total
  } catch (error) {
    console.error('加载机器列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  pageRequest.value.page = page
  loadMachines()
}

const handleSizeChange = (size: number) => {
  pageRequest.value.pageSize = size
  loadMachines()
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    online: 'success',
    offline: 'danger',
    maintenance: 'warning'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    online: '在线',
    offline: '离线',
    maintenance: '维护中'
  }
  return statusMap[status] || status
}

const navigateToEnroll = () => {
  router.push('/customer/machines/enroll')
}

const navigateToEnrollments = () => {
  router.push('/customer/machines/enrollments')
}

onMounted(() => {
  loadMachines()
})
</script>

<template>
  <div class="machine-list">
    <div class="page-header">
      <h2 class="page-title">我的机器</h2>
      <div class="page-actions">
        <el-button @click="navigateToEnrollments">添加进度</el-button>
        <el-button type="primary" @click="navigateToEnroll">添加机器</el-button>
      </div>
    </div>

    <!-- 数据表格 -->
    <DataTable
      :data="machines"
      :total="total"
      :loading="loading"
      :current-page="pageRequest.page"
      :page-size="pageRequest.pageSize"
      :show-pagination="true"
      @page-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <el-table-column prop="name" label="机器名称" min-width="150" />
      <el-table-column prop="region" label="区域" width="120" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="GPU信息" min-width="200">
        <template #default="{ row }">
          {{ row.gpuCount }}x {{ row.gpuModel }} ({{ row.gpuMemory }}GB)
        </template>
      </el-table-column>
      <el-table-column label="GPU使用率" width="120">
        <template #default="{ row }">
          <el-progress :percentage="row.gpuUsage || 0" :color="row.gpuUsage > 80 ? '#f56c6c' : '#67c23a'" />
        </template>
      </el-table-column>
      <el-table-column label="内存使用率" width="120">
        <template #default="{ row }">
          <el-progress :percentage="row.memoryUsage || 0" :color="row.memoryUsage > 80 ? '#f56c6c' : '#67c23a'" />
        </template>
      </el-table-column>
      <el-table-column label="到期时间" width="180">
        <template #default="{ row }">
          {{ row.allocatedTo?.expiresAt || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" :disabled="row.status !== 'online'">
            连接
          </el-button>
          <el-button link type="primary" size="small">
            监控
          </el-button>
        </template>
      </el-table-column>
    </DataTable>
  </div>
</template>

<style scoped>
.machine-list {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
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
</style>
