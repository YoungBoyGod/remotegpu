<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getMachineList, deleteMachine, setMachineMaintenance } from '@/api/admin'
import type { Machine } from '@/types/machine'
import type { PageRequest } from '@/types/common'
import DataTable from '@/components/common/DataTable.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const machines = ref<Machine[]>([])
const total = ref(0)
const pageRequest = ref<PageRequest>({
  page: 1,
  pageSize: 10,
  filters: {}
})

// 筛选条件
const filters = ref({
  status: '',
  region: '',
  keyword: ''
})

const loadMachines = async () => {
  try {
    loading.value = true
    const response = await getMachineList({
      ...pageRequest.value,
      filters: filters.value
    })
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

const handleSearch = () => {
  pageRequest.value.page = 1
  loadMachines()
}

const handleReset = () => {
  filters.value = {
    status: '',
    region: '',
    keyword: ''
  }
  handleSearch()
}

const handleAdd = () => {
  router.push('/admin/machines/add')
}

const handleDelete = async (machine: Machine) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除机器 "${machine.name}" 吗?`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await deleteMachine(machine.id)
    ElMessage.success('删除成功')
    loadMachines()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除机器失败:', error)
    }
  }
}

const handleToggleMaintenance = async (machine: Machine) => {
  try {
    const newStatus = machine.status !== 'maintenance'
    await setMachineMaintenance(machine.id, newStatus)
    ElMessage.success(newStatus ? '已设置为维护状态' : '已取消维护状态')
    loadMachines()
  } catch (error) {
    console.error('设置维护状态失败:', error)
  }
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

onMounted(() => {
  loadMachines()
})
</script>

<template>
  <div class="machine-list">
    <div class="page-header">
      <h2 class="page-title">机器列表</h2>
      <el-button type="primary" @click="handleAdd">添加机器</el-button>
    </div>

    <!-- 筛选栏 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters">
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="在线" value="online" />
            <el-option label="离线" value="offline" />
            <el-option label="维护中" value="maintenance" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="filters.region" placeholder="请输入区域" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="filters.keyword" placeholder="机器名称/IP" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

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
      <el-table-column label="CPU/内存" min-width="150">
        <template #default="{ row }">
          {{ row.cpu }} / {{ row.memory }}GB
        </template>
      </el-table-column>
      <el-table-column label="分配状态" width="120">
        <template #default="{ row }">
          <el-tag v-if="row.allocatedTo" type="warning">已分配</el-tag>
          <el-tag v-else type="info">未分配</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="handleToggleMaintenance(row)">
            {{ row.status === 'maintenance' ? '取消维护' : '设为维护' }}
          </el-button>
          <el-button link type="danger" size="small" @click="handleDelete(row)">
            删除
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

.filter-card {
  margin-bottom: 20px;
}
</style>
