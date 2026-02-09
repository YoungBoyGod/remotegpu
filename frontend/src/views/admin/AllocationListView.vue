<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAllocationList, reclaimMachine } from '@/api/admin'
import type { AllocationRecord } from '@/types/allocation'
import type { PageRequest } from '@/types/common'
import DataTable from '@/components/common/DataTable.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'

const router = useRouter()

const loading = ref(false)
const allocations = ref<AllocationRecord[]>([])
const total = ref(0)
const pageRequest = ref<PageRequest>({
  page: 1,
  pageSize: 10,
  filters: {}
})

const filters = ref({
  status: '',
  keyword: ''
})

const loadAllocations = async () => {
  try {
    loading.value = true
    const response = await getAllocationList({
      ...pageRequest.value,
      filters: filters.value
    })
    allocations.value = response.data.list
    total.value = response.data.total
  } catch (error) {
    console.error('加载分配记录失败:', error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  pageRequest.value.page = page
  loadAllocations()
}

const handleSizeChange = (size: number) => {
  pageRequest.value.pageSize = size
  loadAllocations()
}

const handleSearch = () => {
  pageRequest.value.page = 1
  loadAllocations()
}

const handleReset = () => {
  filters.value = {
    status: '',
    keyword: ''
  }
  handleSearch()
}

const handleReclaim = async (allocation: AllocationRecord) => {
  const machineName = allocation.host?.name || allocation.host_id
  try {
    await ElMessageBox.confirm(
      `确定要回收机器 "${machineName}" 吗？回收后客户将无法继续使用。`,
      '回收确认',
      {
        confirmButtonText: '确定回收',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await reclaimMachine(allocation.host_id)
    ElMessage.success('回收成功')
    loadAllocations()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.msg || '回收机器失败')
    }
  }
}

const handleViewMachine = (allocation: AllocationRecord) => {
  router.push(`/admin/machines/${allocation.host_id}`)
}

const formatDateTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    active: 'success',
    pending: 'warning',
    expired: 'danger',
    reclaimed: 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '使用中',
    pending: '待生效',
    expired: '已到期',
    reclaimed: '已回收'
  }
  return statusMap[status] || status
}

onMounted(() => {
  loadAllocations()
})
</script>

<template>
  <div class="allocation-list">
    <div class="page-header">
      <h2 class="page-title">分配记录</h2>
    </div>

    <!-- 筛选栏 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters">
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="使用中" value="active" />
            <el-option label="待生效" value="pending" />
            <el-option label="已到期" value="expired" />
            <el-option label="已回收" value="reclaimed" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="filters.keyword" placeholder="客户/机器名称" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <DataTable
      :data="allocations"
      :total="total"
      :loading="loading"
      :current-page="pageRequest.page"
      :page-size="pageRequest.pageSize"
      :show-pagination="true"
      @page-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <el-table-column label="客户" min-width="150">
        <template #default="{ row }">
          {{ row.customer?.company || row.customer?.display_name || row.customer?.username || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="机器" min-width="150">
        <template #default="{ row }">
          <el-link type="primary" @click="handleViewMachine(row)">
            {{ row.host?.name || row.host_id }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column label="区域" width="120">
        <template #default="{ row }">
          {{ row.host?.region || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="开始时间" width="180">
        <template #default="{ row }">
          {{ formatDateTime(row.start_time) }}
        </template>
      </el-table-column>
      <el-table-column label="到期时间" width="180">
        <template #default="{ row }">
          {{ formatDateTime(row.end_time) }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="handleViewMachine(row)">
            查看机器
          </el-button>
          <el-button
            link
            type="danger"
            size="small"
            :disabled="row.status !== 'active'"
            @click="handleReclaim(row)"
          >
            回收
          </el-button>
        </template>
      </el-table-column>
    </DataTable>
  </div>
</template>

<style scoped>
.allocation-list {
  padding: 24px;
  background: #f5f7fa;
  min-height: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  font-size: 22px;
  font-weight: 700;
  color: #1d2129;
  margin: 0;
}

.filter-card {
  margin-bottom: 16px;
  border-radius: 8px;
  border: none;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}
</style>
