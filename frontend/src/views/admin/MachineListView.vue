<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getMachineList, deleteMachine, setMachineMaintenance, collectMachineSpec, allocateMachine, reclaimMachine, getCustomerList } from '@/api/admin'
import type { Machine } from '@/types/machine'
import type { Customer } from '@/types/customer'
import type { PageRequest } from '@/types/common'
import DataTable from '@/components/common/DataTable.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { CopyDocument } from '@element-plus/icons-vue'

const router = useRouter()

const loading = ref(false)
const machines = ref<Machine[]>([])
const total = ref(0)
const pageRequest = ref<PageRequest>({
  page: 1,
  pageSize: 10,
  filters: {}
})

// 分配弹窗相关
const allocateDialogVisible = ref(false)
const allocateLoading = ref(false)
const currentMachine = ref<Machine | null>(null)
const customers = ref<Customer[]>([])
const allocateForm = ref({
  customer_id: null as number | null,
  duration_months: 1,
  remark: ''
})

// 计算序号
const getRowIndex = (index: number) => {
  return (pageRequest.value.page - 1) * pageRequest.value.pageSize + index + 1
}

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
    await deleteMachine(String(machine.id))
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
    await setMachineMaintenance(String(machine.id), newStatus)
    ElMessage.success(newStatus ? '已设置为维护状态' : '已取消维护状态')
    loadMachines()
  } catch (error) {
    console.error('设置维护状态失败:', error)
  }
}

const handleCollectSpec = async (machine: Machine) => {
  try {
    await ElMessageBox.confirm(
      `确认对机器 "${machine.name}" 进行硬件补采?`,
      '补采确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await collectMachineSpec(String(machine.id))
    ElMessage.success('补采任务已触发')
    loadMachines()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('补采硬件失败:', error)
    }
  }
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    idle: 'success',
    allocated: 'primary',
    offline: 'danger',
    maintenance: 'warning'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    idle: '空闲',
    allocated: '已分配',
    offline: '离线',
    maintenance: '维护中'
  }
  return statusMap[status] || status
}

const getCollectType = (needsCollect?: boolean) => (needsCollect ? 'warning' : 'success')
const getCollectText = (needsCollect?: boolean) => (needsCollect ? '待采集' : '已采集')

const getSSHCommand = (machine: Machine) => {
  if (machine.ssh_command) return machine.ssh_command
  const host = machine.ssh_host || machine.public_ip || machine.ip_address || ''
  const user = machine.ssh_username || 'root'
  const port = machine.ssh_port || 22
  if (!host) return ''
  return `ssh -p ${port} ${user}@${host}`
}

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

// 跳转到机器详情
const handleViewDetail = (machine: Machine) => {
  router.push(`/admin/machines/${machine.id}`)
}

// 打开分配弹窗
const handleAllocate = async (machine: Machine) => {
  currentMachine.value = machine
  allocateForm.value = {
    customer_id: null,
    duration_months: 1,
    remark: ''
  }
  allocateDialogVisible.value = true

  // 加载客户列表
  try {
    const response = await getCustomerList({ page: 1, pageSize: 100 })
    customers.value = response.data.list || []
  } catch (error) {
    console.error('加载客户列表失败:', error)
  }
}

// 确认分配
const handleConfirmAllocate = async () => {
  if (!currentMachine.value || !allocateForm.value.customer_id) {
    ElMessage.warning('请选择客户')
    return
  }

  try {
    allocateLoading.value = true
    await allocateMachine(String(currentMachine.value.id), {
      customer_id: allocateForm.value.customer_id,
      duration_months: allocateForm.value.duration_months,
      remark: allocateForm.value.remark
    })
    ElMessage.success('分配成功')
    allocateDialogVisible.value = false
    loadMachines()
  } catch (error) {
    console.error('分配失败:', error)
  } finally {
    allocateLoading.value = false
  }
}

// 回收机器
const handleReclaim = async (machine: Machine) => {
  try {
    await ElMessageBox.confirm(
      `确定要回收机器 "${machine.name || machine.id}" 吗？回收后客户将无法继续使用。`,
      '回收确认',
      {
        confirmButtonText: '确定回收',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await reclaimMachine(String(machine.id))
    ElMessage.success('回收成功')
    loadMachines()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.msg || '回收失败')
    }
  }
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
            <el-option label="空闲" value="idle" />
            <el-option label="已分配" value="allocated" />
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
      <el-table-column label="序号" width="70" align="center">
        <template #default="{ $index }">
          {{ getRowIndex($index) }}
        </template>
      </el-table-column>
      <el-table-column prop="id" label="ID" width="150" show-overflow-tooltip />
      <el-table-column label="机器名称" min-width="150">
        <template #default="{ row }">
          <el-link type="primary" @click="handleViewDetail(row)">
            {{ row.name || row.hostname || row.id }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column prop="region" label="区域" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="采集状态" width="120">
        <template #default="{ row }">
          <el-tag :type="getCollectType(row.needs_collect)">
            {{ getCollectText(row.needs_collect) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="GPU信息" min-width="200">
        <template #default="{ row }">
          <template v-if="row.gpus && row.gpus.length > 0">
            {{ row.gpus.length }}x {{ row.gpus[0].name }} ({{ Math.round(row.gpus[0].memory_total_mb / 1024) }}GB)
          </template>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="CPU/内存" min-width="150">
        <template #default="{ row }">
          {{ row.total_cpu || '-' }} 核 / {{ row.total_memory_gb ? row.total_memory_gb + ' GB' : '-' }}
        </template>
      </el-table-column>
      <el-table-column label="SSH 连接" min-width="220">
        <template #default="{ row }">
          <template v-if="getSSHCommand(row)">
            <el-tooltip :content="getSSHCommand(row)" placement="top">
              <code class="ssh-cmd">{{ getSSHCommand(row) }}</code>
            </el-tooltip>
            <el-button
              link
              type="primary"
              size="small"
              :icon="CopyDocument"
              @click.stop="copyToClipboard(getSSHCommand(row))"
            />
          </template>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="分配状态" width="120">
        <template #default="{ row }">
          <el-tag v-if="row.status === 'allocated'" type="warning">已分配</el-tag>
          <el-tag v-else type="info">未分配</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="280" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status !== 'allocated'"
            link
            type="success"
            size="small"
            @click="handleAllocate(row)"
          >
            分配
          </el-button>
          <el-button
            v-if="row.status === 'allocated'"
            link
            type="danger"
            size="small"
            @click="handleReclaim(row)"
          >
            回收
          </el-button>
          <el-button
            v-if="row.needs_collect"
            link
            type="warning"
            size="small"
            @click="handleCollectSpec(row)"
          >
            补采
          </el-button>
          <el-button link type="primary" size="small" @click="handleToggleMaintenance(row)">
            {{ row.status === 'maintenance' ? '取消维护' : '设为维护' }}
          </el-button>
          <el-button link type="danger" size="small" @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </DataTable>

    <!-- 分配弹窗 -->
    <el-dialog
      v-model="allocateDialogVisible"
      title="分配机器"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="allocateForm" label-width="100px">
        <el-form-item label="机器">
          <span>{{ currentMachine?.name || currentMachine?.id }}</span>
        </el-form-item>
        <el-form-item label="选择客户" required>
          <el-select
            v-model="allocateForm.customer_id"
            placeholder="请选择客户"
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="customer in customers"
              :key="customer.id"
              :label="customer.company || customer.username"
              :value="customer.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="租用时长">
          <el-input-number
            v-model="allocateForm.duration_months"
            :min="1"
            :max="36"
          />
          <span style="margin-left: 8px">个月</span>
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="allocateForm.remark"
            type="textarea"
            :rows="2"
            placeholder="可选备注"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="allocateDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="allocateLoading" @click="handleConfirmAllocate">
          确认分配
        </el-button>
      </template>
    </el-dialog>
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

.ssh-cmd {
  font-size: 12px;
  color: #606266;
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 3px;
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
  vertical-align: middle;
}
</style>
