<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getMyMachines, getMachineConnection } from '@/api/customer'
import type { Machine } from '@/types/machine'
import DataTable from '@/components/common/DataTable.vue'

const loading = ref(false)
const machines = ref<Machine[]>([])
const total = ref(0)
const router = useRouter()
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

const filteredMachines = computed(() => {
  const kw = searchKeyword.value.trim().toLowerCase()
  if (!kw) return machines.value
  return machines.value.filter((m) =>
    (m.hostname || '').toLowerCase().includes(kw) ||
    (m.ip_address || '').toLowerCase().includes(kw)
  )
})

const loadMachines = async () => {
  try {
    loading.value = true
    const response = await getMyMachines({
      page: currentPage.value,
      pageSize: pageSize.value,
    })
    machines.value = response.data?.list || []
    total.value = response.data?.total || 0
  } catch (error: any) {
    ElMessage.error(error?.msg || error?.message || '加载机器列表失败')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  loadMachines()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadMachines()
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    available: 'success',
    allocated: 'primary',
    maintenance: 'warning',
    online: 'success',
    offline: 'danger',
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    available: '可用',
    allocated: '已分配',
    maintenance: '维护中',
    online: '在线',
    offline: '离线',
  }
  return statusMap[status] || status
}

const formatDate = (dateStr: string | undefined) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const navigateToEnroll = () => {
  router.push('/customer/machines/enroll')
}

const navigateToEnrollments = () => {
  router.push('/customer/machines/enrollments')
}

const handleViewDetail = (machine: Machine) => {
  router.push(`/customer/machines/${machine.id}`)
}

const handleConnect = async (machine: Machine) => {
  try {
    const response = await getMachineConnection(machine.id)
    const conn = response.data
    if (conn?.ssh) {
      const cmd = `ssh ${conn.ssh.username}@${conn.ssh.host} -p ${conn.ssh.port}`
      await navigator.clipboard.writeText(cmd)
      ElMessage.success('SSH 连接命令已复制到剪贴板')
    } else {
      ElMessage.warning('暂无连接信息')
    }
  } catch (error: any) {
    ElMessage.error(error?.msg || error?.message || '获取连接信息失败')
  }
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

    <!-- 搜索栏 -->
    <el-card class="filter-card" shadow="never">
      <el-row :gutter="16" align="middle">
        <el-col :span="8">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索主机名或 IP"
            clearable
            @clear="loadMachines"
            @keyup.enter="loadMachines"
          />
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="loadMachines">搜索</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 数据表格 -->
    <DataTable
      :data="filteredMachines"
      :total="filteredMachines.length"
      :loading="loading"
      :current-page="currentPage"
      :page-size="pageSize"
      :show-pagination="true"
      @page-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <el-table-column label="机器名称" min-width="150">
        <template #default="{ row }">
          <el-link type="primary" @click="handleViewDetail(row)">
            {{ row.hostname || row.id }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column prop="ip_address" label="IP 地址" width="140" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="CPU" width="80">
        <template #default="{ row }">
          {{ row.total_cpu || '-' }} 核
        </template>
      </el-table-column>
      <el-table-column label="内存" width="100">
        <template #default="{ row }">
          {{ row.total_memory_gb ? row.total_memory_gb + ' GB' : '-' }}
        </template>
      </el-table-column>
      <el-table-column label="分配时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.start_time) }}
        </template>
      </el-table-column>
      <el-table-column label="到期时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.end_time) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="handleConnect(row)">
            连接
          </el-button>
          <el-button link type="primary" size="small" @click="handleViewDetail(row)">
            详情
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

.filter-card {
  margin-bottom: 16px;
}
</style>
