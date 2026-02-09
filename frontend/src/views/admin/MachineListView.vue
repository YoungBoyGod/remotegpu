<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getMachineList, deleteMachine, setMachineMaintenance, collectMachineSpec, allocateMachine, reclaimMachine, getCustomerList, batchImportMachines, batchSetMaintenance, batchAllocate, batchReclaim } from '@/api/admin'
import type { ImportMachineItem } from '@/api/admin'
import type { Machine } from '@/types/machine'
import type { Customer } from '@/types/customer'
import type { PageRequest } from '@/types/common'
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

// 密码显示状态
const passwordVisible = ref<Record<string | number, boolean>>({})

// 批量导入相关
const importDialogVisible = ref(false)
const importLoading = ref(false)
const importText = ref('')
const importPreview = ref<ImportMachineItem[]>([])
const importError = ref('')

// 批量操作相关
const selectedMachines = ref<Machine[]>([])
const batchAllocateDialogVisible = ref(false)
const batchAllocateLoading = ref(false)
const batchAllocateForm = ref({
  customer_id: null as number | null,
  duration_months: 1,
  remark: ''
})

const hasSelection = computed(() => selectedMachines.value.length > 0)
const isAllSelected = computed(() => machines.value.length > 0 && selectedMachines.value.length === machines.value.length)
const isIndeterminate = computed(() => selectedMachines.value.length > 0 && selectedMachines.value.length < machines.value.length)

const toggleSelectAll = (val: boolean) => {
  selectedMachines.value = val ? [...machines.value] : []
}

const handleSelectionChange = (selection: Machine[]) => {
  selectedMachines.value = selection
}

const toggleSelect = (machine: Machine) => {
  const idx = selectedMachines.value.findIndex(m => m.id === machine.id)
  if (idx >= 0) {
    selectedMachines.value.splice(idx, 1)
  } else {
    selectedMachines.value.push(machine)
  }
}

// 解析批量导入文本（CSV 格式）
const parseImportText = () => {
  importError.value = ''
  importPreview.value = []
  if (!importText.value.trim()) return

  const lines = importText.value.trim().split('\n').filter(l => l.trim())
  const items: ImportMachineItem[] = []

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    if (!line) continue
    const trimmed = line.trim()
    // 跳过表头行
    if (i === 0 && trimmed.includes('host_ip')) continue

    const cols = trimmed.split(',').map(c => c.trim())
    if (cols.length < 3) {
      importError.value = `第 ${i + 1} 行格式错误：需要至少 3 列（IP, 端口, 用户名）`
      return
    }

    items.push({
      host_ip: cols[0] ?? '',
      ssh_port: parseInt(cols[1] ?? '') || 22,
      ssh_username: cols[2] ?? 'root',
      ssh_password: cols[3] ?? '',
    })
  }

  importPreview.value = items
}

const handleImport = () => {
  importDialogVisible.value = true
  importText.value = ''
  importPreview.value = []
  importError.value = ''
}

const handleConfirmImport = async () => {
  if (importPreview.value.length === 0) {
    ElMessage.warning('请先输入并解析导入数据')
    return
  }

  try {
    importLoading.value = true
    const res = await batchImportMachines({ machines: importPreview.value })
    ElMessage.success(`成功导入 ${res.data.count} 台机器`)
    importDialogVisible.value = false
    loadMachines()
  } catch (error) {
    console.error('批量导入失败:', error)
  } finally {
    importLoading.value = false
  }
}

// 计算序号
const getRowIndex = (index: number) => {
  return (pageRequest.value.page - 1) * pageRequest.value.pageSize + index + 1
}

// 筛选条件
const filters = ref({
  status: '',
  region: '',
  gpu_model: '',
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
    gpu_model: '',
    keyword: ''
  }
  handleSearch()
}

const handleAdd = () => {
  router.push('/admin/machines/add')
}

const handleDelete = async (machine: Machine) => {
  if (machine.allocation_status === 'allocated') {
    ElMessage.warning('已分配的机器不能删除，请先回收')
    return
  }
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
    const newStatus = machine.allocation_status !== 'maintenance'
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

const getDeviceStatusType = (status?: string) => {
  return status === 'online' ? 'success' : 'danger'
}

const getDeviceStatusText = (status?: string) => {
  return status === 'online' ? '在线' : '离线'
}

const getAllocStatusType = (status?: string) => {
  const map: Record<string, string> = {
    idle: 'success',
    allocated: 'warning',
    maintenance: 'info'
  }
  return map[status || ''] || 'info'
}

const getAllocStatusText = (status?: string) => {
  const map: Record<string, string> = {
    idle: '空闲',
    allocated: '已分配',
    maintenance: '维护中'
  }
  return map[status || ''] || status || '-'
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
    // 只显示正常状态的客户
    customers.value = (response.data.list || []).filter((c: Customer) => c.status === 'active')
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

// 批量设为维护
const handleBatchMaintenance = async (maintenance: boolean) => {
  const ids = selectedMachines.value.map(m => String(m.id))
  const action = maintenance ? '设为维护' : '取消维护'
  try {
    await ElMessageBox.confirm(
      `确定要将选中的 ${ids.length} 台机器${action}吗？`,
      '批量操作确认',
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    const res = await batchSetMaintenance(ids, maintenance)
    ElMessage.success(`${action}成功，影响 ${res.data.affected} 台机器`)
    loadMachines()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.msg || `批量${action}失败`)
    }
  }
}

// 批量分配 - 打开弹窗
const handleBatchAllocate = async () => {
  batchAllocateForm.value = { customer_id: null, duration_months: 1, remark: '' }
  batchAllocateDialogVisible.value = true
  try {
    const response = await getCustomerList({ page: 1, pageSize: 100 })
    // 只显示正常状态的客户
    customers.value = (response.data.list || []).filter((c: Customer) => c.status === 'active')
  } catch (error) {
    console.error('加载客户列表失败:', error)
  }
}

// 确认批量分配
const handleConfirmBatchAllocate = async () => {
  if (!batchAllocateForm.value.customer_id) {
    ElMessage.warning('请选择客户')
    return
  }
  const ids = selectedMachines.value.map(m => String(m.id))
  try {
    batchAllocateLoading.value = true
    const res = await batchAllocate(
      ids,
      batchAllocateForm.value.customer_id,
      batchAllocateForm.value.duration_months,
      batchAllocateForm.value.remark
    )
    const successCount = res.data.success?.length ?? 0
    const failedCount = res.data.failed?.length ?? 0
    ElMessage.success(`批量分配完成：成功 ${successCount}，失败 ${failedCount}`)
    batchAllocateDialogVisible.value = false
    loadMachines()
  } catch (error: any) {
    ElMessage.error(error?.msg || '批量分配失败')
  } finally {
    batchAllocateLoading.value = false
  }
}

// 批量回收
const handleBatchReclaim = async () => {
  const ids = selectedMachines.value.map(m => String(m.id))
  try {
    await ElMessageBox.confirm(
      `确定要批量回收选中的 ${ids.length} 台机器吗？回收后客户将无法继续使用。`,
      '批量回收确认',
      { confirmButtonText: '确定回收', cancelButtonText: '取消', type: 'warning' }
    )
    const res = await batchReclaim(ids)
    const successCount = res.data.success?.length ?? 0
    const failedCount = res.data.failed?.length ?? 0
    ElMessage.success(`批量回收完成：成功 ${successCount}，失败 ${failedCount}`)
    loadMachines()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.msg || '批量回收失败')
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
      <div class="header-actions">
        <template v-if="hasSelection">
          <el-button type="warning" @click="handleBatchMaintenance(true)">批量维护</el-button>
          <el-button type="success" @click="handleBatchAllocate">批量分配</el-button>
          <el-button type="danger" @click="handleBatchReclaim">批量回收</el-button>
          <span class="selection-count">已选 {{ selectedMachines.length }} 台</span>
        </template>
        <el-button @click="handleImport">批量导入</el-button>
        <el-button type="primary" @click="handleAdd">添加机器</el-button>
      </div>
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
        <el-form-item label="GPU型号">
          <el-input v-model="filters.gpu_model" placeholder="如 A100/H100" clearable style="width: 150px" />
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

    <!-- 机器列表 -->
    <div v-loading="loading" class="machine-items">
      <!-- 列表头：全选 -->
      <div class="list-header" v-if="machines.length > 0">
        <el-checkbox
          :model-value="isAllSelected"
          :indeterminate="isIndeterminate"
          @change="toggleSelectAll"
        />
        <span class="select-all-label">
          {{ hasSelection ? `已选 ${selectedMachines.length} / ${machines.length} 台` : '全选' }}
        </span>
      </div>

      <div
        v-for="(machine, index) in machines"
        :key="machine.id"
        class="machine-row"
        :class="{ 'is-selected': selectedMachines.some(m => m.id === machine.id) }"
      >
        <!-- 顶栏：序号 + 勾选 + 名称 + 状态 + 操作 -->
        <div class="row-header">
          <span class="row-index">{{ getRowIndex(index) }}</span>
          <el-checkbox
            :model-value="selectedMachines.some(m => m.id === machine.id)"
            @change="toggleSelect(machine)"
          />
          <el-link type="primary" class="row-name" @click="handleViewDetail(machine)">
            {{ machine.name || machine.hostname || machine.id }}
          </el-link>
          <el-tag :type="getDeviceStatusType(machine.device_status)">{{ getDeviceStatusText(machine.device_status) }}</el-tag>
          <el-tag :type="getAllocStatusType(machine.allocation_status)">{{ getAllocStatusText(machine.allocation_status) }}</el-tag>
          <div class="row-actions">
            <el-button v-if="machine.allocation_status === 'idle'" type="success" size="small" plain @click="handleAllocate(machine)">分配</el-button>
            <el-button v-if="machine.allocation_status === 'allocated'" type="danger" size="small" plain @click="handleReclaim(machine)">回收</el-button>
            <el-button v-if="machine.needs_collect" type="warning" size="small" plain @click="handleCollectSpec(machine)">补采</el-button>
            <el-button v-if="machine.allocation_status !== 'allocated'" size="small" plain @click="handleToggleMaintenance(machine)">{{ machine.allocation_status === 'maintenance' ? '取消维护' : '维护' }}</el-button>
            <el-button type="danger" size="small" text :disabled="machine.allocation_status === 'allocated'" @click="handleDelete(machine)">删除</el-button>
          </div>
        </div>

        <!-- 三栏：配置 | 对内 | 对外 -->
        <div class="row-body">
          <div class="body-col">
            <div class="col-title">配置信息</div>
            <div class="spec-row">
              <span class="spec-label">地域</span>
              <span class="spec-value">{{ machine.region || '-' }}</span>
            </div>
            <div class="spec-row">
              <span class="spec-label">GPU</span>
              <span class="spec-value" v-if="machine.gpus && machine.gpus.length > 0">{{ machine.gpus.length }}x {{ machine.gpus[0]?.name }} ({{ Math.round((machine.gpus[0]?.memory_total_mb ?? 0) / 1024) }}GB)</span>
              <span class="spec-value text-muted" v-else>-</span>
            </div>
            <div class="spec-row">
              <span class="spec-label">CPU</span>
              <span class="spec-value">{{ machine.cpu_info || '-' }}{{ machine.total_cpu ? ' / ' + machine.total_cpu + ' 核' : '' }}</span>
            </div>
            <div class="spec-row">
              <span class="spec-label">内存</span>
              <span class="spec-value">{{ machine.total_memory_gb ? machine.total_memory_gb + ' GB' : '-' }}</span>
            </div>
            <div class="spec-row">
              <span class="spec-label">硬盘</span>
              <span class="spec-value">{{ machine.total_disk_gb ? machine.total_disk_gb + ' GB' : '-' }}</span>
            </div>
          </div>

          <div class="body-col">
            <div class="col-title">对内连接</div>
            <div class="conn-line">
              <span class="conn-tag">SSH</span>
              <template v-if="machine.ssh_host || machine.ip_address">
                <span class="conn-field"><span class="conn-label">IP:</span> {{ machine.ip_address || '-' }}</span>
                <span class="conn-field"><span class="conn-label">端口:</span> {{ machine.ssh_port || 22 }}</span>
                <span class="conn-field"><span class="conn-label">用户:</span> {{ machine.ssh_username || 'root' }}</span>
                <span class="conn-field" v-if="machine.ssh_password">
                  <span class="conn-label">密码:</span>
                  <span class="password-text">{{ machine.ssh_password }}</span>
                  <el-button link :icon="CopyDocument" size="small" @click="copyToClipboard(machine.ssh_password)" />
                </span>
              </template>
              <span v-else class="text-muted">未配置</span>
            </div>
            <div class="conn-line">
              <span class="conn-tag conn-tag-jupyter">Jupyter</span>
              <template v-if="machine.jupyter_url">
                <a class="conn-link" :href="machine.jupyter_url" target="_blank">{{ machine.jupyter_url }}</a>
                <el-button link :icon="CopyDocument" size="small" @click="copyToClipboard(machine.jupyter_url)" />
              </template>
              <span v-else class="text-muted">未配置</span>
            </div>
            <div class="conn-line">
              <span class="conn-tag conn-tag-vnc">VNC</span>
              <template v-if="machine.vnc_url">
                <a class="conn-link" :href="machine.vnc_url" target="_blank">{{ machine.vnc_url }}</a>
                <el-button link :icon="CopyDocument" size="small" @click="copyToClipboard(machine.vnc_url)" />
              </template>
              <span v-else class="text-muted">未配置</span>
            </div>
          </div>

          <div class="body-col">
            <div class="col-title">对外连接</div>
            <div class="conn-line">
              <span class="conn-tag conn-tag-ext">SSH</span>
              <template v-if="machine.nginx_domain || machine.external_ip">
                <span class="conn-field">{{ machine.nginx_domain || machine.external_ip }}:{{ machine.external_ssh_port || '-' }}</span>
                <el-button link :icon="CopyDocument" size="small" @click="copyToClipboard(`ssh -p ${machine.external_ssh_port || 22} ${machine.ssh_username || 'root'}@${machine.nginx_domain || machine.external_ip}`)" />
              </template>
              <span v-else class="text-muted">未配置</span>
            </div>
            <div class="conn-line">
              <span class="conn-tag conn-tag-ext">Jupyter</span>
              <template v-if="machine.external_jupyter_port && (machine.nginx_domain || machine.external_ip)">
                <span class="conn-field">{{ machine.nginx_domain || machine.external_ip }}:{{ machine.external_jupyter_port }}</span>
                <el-button link :icon="CopyDocument" size="small" @click="copyToClipboard(`${machine.nginx_domain || machine.external_ip}:${machine.external_jupyter_port}`)" />
              </template>
              <span v-else class="text-muted">未配置</span>
            </div>
            <div class="conn-line">
              <span class="conn-tag conn-tag-ext">VNC</span>
              <template v-if="machine.external_vnc_port && (machine.nginx_domain || machine.external_ip)">
                <span class="conn-field">{{ machine.nginx_domain || machine.external_ip }}:{{ machine.external_vnc_port }}</span>
                <el-button link :icon="CopyDocument" size="small" @click="copyToClipboard(`${machine.nginx_domain || machine.external_ip}:${machine.external_vnc_port}`)" />
              </template>
              <span v-else class="text-muted">未配置</span>
            </div>
          </div>
        </div>
      </div>

      <el-empty v-if="!loading && machines.length === 0" description="暂无机器数据" />
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="pageRequest.page"
        v-model:page-size="pageRequest.pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </div>

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

    <!-- 批量导入弹窗 -->
    <el-dialog
      v-model="importDialogVisible"
      title="批量导入机器"
      width="700px"
      :close-on-click-modal="false"
    >
      <div class="import-hint">
        请输入 CSV 格式数据，每行一台机器，字段顺序：<br />
        <code>IP地址, SSH端口, 用户名, 密码</code>
      </div>
      <el-input
        v-model="importText"
        type="textarea"
        :rows="8"
        placeholder="192.168.1.10, 22, root, password123"
      />
      <div class="import-actions">
        <el-button @click="parseImportText">解析预览</el-button>
        <span v-if="importError" class="import-error">{{ importError }}</span>
        <span v-else-if="importPreview.length > 0" class="import-count">
          共解析 {{ importPreview.length }} 条记录
        </span>
      </div>
      <el-table
        v-if="importPreview.length > 0"
        :data="importPreview"
        border
        max-height="250"
        style="margin-top: 12px"
      >
        <el-table-column prop="host_ip" label="IP地址" width="160" />
        <el-table-column prop="ssh_port" label="SSH端口" width="100" />
        <el-table-column prop="ssh_username" label="用户名" width="120" />
        <el-table-column prop="ssh_password" label="密码" width="140" />
      </el-table>
      <template #footer>
        <el-button @click="importDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="importLoading"
          :disabled="importPreview.length === 0"
          @click="handleConfirmImport"
        >
          确认导入
        </el-button>
      </template>
    </el-dialog>

    <!-- 批量分配弹窗 -->
    <el-dialog
      v-model="batchAllocateDialogVisible"
      title="批量分配机器"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="batchAllocateForm" label-width="100px">
        <el-form-item label="已选机器">
          <span>{{ selectedMachines.length }} 台</span>
        </el-form-item>
        <el-form-item label="选择客户" required>
          <el-select
            v-model="batchAllocateForm.customer_id"
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
            v-model="batchAllocateForm.duration_months"
            :min="1"
            :max="36"
          />
          <span style="margin-left: 8px">个月</span>
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="batchAllocateForm.remark"
            type="textarea"
            :rows="2"
            placeholder="可选备注"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="batchAllocateDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="batchAllocateLoading" @click="handleConfirmBatchAllocate">
          确认分配
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.machine-list {
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

.selection-count {
  font-size: 13px;
  color: #409eff;
  margin-left: 4px;
}

.header-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

/* 勾选框放大 */
.machine-items :deep(.el-checkbox__inner) {
  width: 18px;
  height: 18px;
}

.machine-items :deep(.el-checkbox__inner::after) {
  height: 9px;
  left: 6px;
  top: 2px;
  width: 4px;
}

/* 机器列表 */
.machine-items {
  min-height: 200px;
}

/* 列表头：全选栏 */
.list-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  margin-bottom: 8px;
  background: #f0f1f5;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
}

.select-all-label {
  font-size: 13px;
  color: #606266;
  user-select: none;
}

.machine-row {
  border: 1px solid #e5e6eb;
  border-radius: 10px;
  background: #fff;
  margin-bottom: 12px;
  transition: box-shadow 0.2s, border-color 0.2s;
}

.machine-row:hover {
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.07);
  border-color: #c9cdd4;
}

.machine-row.is-selected {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.12);
}

/* 顶栏 */
.row-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
}

.row-index {
  font-size: 13px;
  color: #86909c;
  min-width: 20px;
  text-align: center;
  flex-shrink: 0;
}

.row-name {
  font-weight: 600;
  font-size: 15px;
  flex-shrink: 0;
}

.row-actions {
  margin-left: auto;
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

/* 三栏布局 */
.row-body {
  display: flex;
  border-top: 1px solid #f2f3f5;
}

.body-col {
  flex: 1;
  padding: 12px 20px;
}

.body-col + .body-col {
  border-left: 1px solid #f2f3f5;
}

.col-title {
  font-size: 12px;
  font-weight: 600;
  color: #86909c;
  margin-bottom: 10px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.spec-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
  padding: 3px 0;
  font-size: 13px;
}

.spec-label {
  color: #86909c;
  font-size: 12px;
  min-width: 32px;
  flex-shrink: 0;
}

.spec-value {
  color: #1d2129;
}

/* 连接信息行 */
.conn-line {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #4e5969;
  padding: 5px 0;
  border-top: 1px dashed #f2f3f5;
}

.conn-line:first-child {
  border-top: 1px solid #f2f3f5;
}

.conn-tag {
  display: inline-block;
  min-width: 52px;
  text-align: center;
  font-size: 11px;
  font-weight: 600;
  color: #409eff;
  background: #e8f3ff;
  padding: 2px 8px;
  border-radius: 4px;
  flex-shrink: 0;
}

.conn-tag-jupyter {
  color: #e6a23c;
  background: #fff7e6;
}

.conn-tag-vnc {
  color: #67c23a;
  background: #e8ffea;
}

.conn-field {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

.conn-label {
  color: #86909c;
  flex-shrink: 0;
}

.conn-empty {
  color: #c9cdd4;
  font-size: 13px;
  border-top: 1px solid #f2f3f5;
  padding: 4px 0;
}

.ssh-cmd {
  font-size: 12px;
  color: #606266;
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 3px;
  display: inline-block;
  vertical-align: middle;
}

.conn-link {
  color: #409eff;
  text-decoration: none;
  font-size: 13px;
}

.conn-link:hover {
  text-decoration: underline;
}

.conn-tag-ext {
  color: #f56c6c;
  background: #ffece8;
}

.password-text {
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 13px;
  color: #1d2129;
}

.text-muted {
  color: #c9cdd4;
}

/* 分页 */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

/* 导入弹窗 */
.import-hint {
  margin-bottom: 12px;
  font-size: 13px;
  color: #606266;
  line-height: 1.8;
}

.import-hint code {
  font-size: 12px;
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 3px;
  color: #409eff;
}

.import-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 12px;
}

.import-error {
  color: #f56c6c;
  font-size: 13px;
}

.import-count {
  color: #67c23a;
  font-size: 13px;
}
</style>
