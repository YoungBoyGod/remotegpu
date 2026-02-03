<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getAllocationList, reclaimMachine } from '@/api/admin'
import type { AllocationRecord } from '@/types/allocation'
import type { PageRequest } from '@/types/common'
import DataTable from '@/components/common/DataTable.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

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

const handleQuickAllocate = () => {
  router.push('/admin/allocations/quick')
}

const handleReclaim = async (allocation: AllocationRecord) => {
  try {
    await ElMessageBox.confirm(
      `确定要回收机器 "${allocation.machineName}" 吗?`,
      '回收确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await reclaimMachine(allocation.id)
    ElMessage.success('回收成功')
    loadAllocations()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('回收机器失败:', error)
    }
  }
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    active: 'success',
    expiring: 'warning',
    expired: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '使用中',
    expiring: '即将到期',
    expired: '已到期'
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
      <el-button type="primary" @click="handleQuickAllocate">快速分配</el-button>
    </div>

    <!-- 筛选栏 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters">
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="使用中" value="active" />
            <el-option label="即将到期" value="expiring" />
            <el-option label="已到期" value="expired" />
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
      <el-table-column prop="customerName" label="客户名称" min-width="150" />
      <el-table-column prop="machineName" label="机器名称" min-width="150" />
      <el-table-column prop="region" label="区域" width="120" />
      <el-table-column prop="allocatedAt" label="分配时间" width="180" />
      <el-table-column prop="duration" label="分配时长(天)" width="120" />
      <el-table-column prop="expiresAt" label="到期时间" width="180" />
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button
            link
            type="danger"
            size="small"
            :disabled="row.status === 'expired'"
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

.filter-card {
  margin-bottom: 20px;
}
</style>
