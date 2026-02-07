<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getMachineList, deleteMachine, setMachineMaintenance, collectMachineSpec, allocateMachine, reclaimMachine, getCustomerList, batchImportMachines } from '@/api/admin'
import type { ImportMachineItem } from '@/api/admin'
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

// 批量导入相关
const importDialogVisible = ref(false)
const importLoading = ref(false)
const importText = ref('')
const importPreview = ref<ImportMachineItem[]>([])
const importError = ref('')

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
    if (cols.length < 9) {
      importError.value = `第 ${i + 1} 行格式错误：需要至少 9 列`
      return
    }

    items.push({
      host_ip: cols[0] ?? '',
      ssh_port: parseInt(cols[1] ?? '') || 22,
      region: cols[2] ?? '',
      gpu_model: cols[3] ?? '',
      gpu_count: parseInt(cols[4] ?? '') || 0,
      cpu_cores: parseInt(cols[5] ?? '') || 0,
      ram_size: parseInt(cols[6] ?? '') || 0,
      disk_size: parseInt(cols[7] ?? '') || 0,
      price_hourly: parseFloat(cols[8] ?? '') || 0,
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
    const newStatus = (machine as any).allocation_status !== 'maintenance'
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
      <div class="header-actions">
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
      <el-table-column type="expand">
        <template #default="{ row }">
          <div class="expand-connection">
            <div class="conn-section">
              <span class="conn-title">SSH</span>
              <template v-if="getSSHCommand(row)">
                <div class="conn-row">
                  <span class="conn-label">主机：</span>
                  <span>{{ row.ssh_host || row.public_ip || row.ip_address || '-' }}</span>
                </div>
                <div class="conn-row">
                  <span class="conn-label">端口：</span>
                  <span>{{ row.ssh_port || 22 }}</span>
                </div>
                <div class="conn-row">
                  <span class="conn-label">用户：</span>
                  <span>{{ row.ssh_username || 'root' }}</span>
                </div>
                <div class="conn-row">
                  <span class="conn-label">命令：</span>
                  <code class="ssh-cmd">{{ getSSHCommand(row) }}</code>
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(getSSHCommand(row))" />
                </div>
              </template>
              <span v-else class="conn-empty">未配置</span>
            </div>
            <div class="conn-section">
              <span class="conn-title">Jupyter</span>
              <template v-if="row.jupyter_url">
                <div class="conn-row">
                  <a class="conn-link" :href="row.jupyter_url" target="_blank">{{ row.jupyter_url }}</a>
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(row.jupyter_url)" />
                </div>
              </template>
              <span v-else class="conn-empty">未配置</span>
            </div>
            <div class="conn-section">
              <span class="conn-title">VNC</span>
              <template v-if="row.vnc_url">
                <div class="conn-row">
                  <a class="conn-link" :href="row.vnc_url" target="_blank">{{ row.vnc_url }}</a>
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(row.vnc_url)" />
                </div>
              </template>
              <span v-else class="conn-empty">未配置</span>
            </div>
          </div>
        </template>
      </el-table-column>
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
      <el-table-column label="设备状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getDeviceStatusType(row.device_status)">
            {{ getDeviceStatusText(row.device_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="分配状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getAllocStatusType(row.allocation_status)">
            {{ getAllocStatusText(row.allocation_status) }}
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
      <el-table-column label="操作" width="280" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.allocation_status !== 'allocated'"
            link
            type="success"
            size="small"
            @click="handleAllocate(row)"
          >
            分配
          </el-button>
          <el-button
            v-if="row.allocation_status === 'allocated'"
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
            {{ row.allocation_status === 'maintenance' ? '取消维护' : '设为维护' }}
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

    <!-- 批量导入弹窗 -->
    <el-dialog
      v-model="importDialogVisible"
      title="批量导入机器"
      width="700px"
      :close-on-click-modal="false"
    >
      <div class="import-hint">
        请输入 CSV 格式数据，每行一台机器，字段顺序：<br />
        <code>host_ip, ssh_port, region, gpu_model, gpu_count, cpu_cores, ram_size, disk_size, price_hourly</code>
      </div>
      <el-input
        v-model="importText"
        type="textarea"
        :rows="8"
        placeholder="192.168.1.10, 22, cn-east, A100, 8, 64, 512, 2000, 50.0"
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
        <el-table-column prop="host_ip" label="IP" width="140" />
        <el-table-column prop="ssh_port" label="SSH端口" width="80" />
        <el-table-column prop="region" label="区域" width="90" />
        <el-table-column prop="gpu_model" label="GPU型号" width="90" />
        <el-table-column prop="gpu_count" label="GPU数" width="70" />
        <el-table-column prop="cpu_cores" label="CPU核" width="70" />
        <el-table-column prop="ram_size" label="内存GB" width="80" />
        <el-table-column prop="disk_size" label="磁盘GB" width="80" />
        <el-table-column prop="price_hourly" label="时价" width="70" />
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

.expand-connection {
  display: flex;
  gap: 32px;
  padding: 12px 24px;
  background: #f9fafb;
}

.conn-section {
  flex: 1;
}

.conn-title {
  display: block;
  font-weight: 600;
  font-size: 13px;
  color: #303133;
  margin-bottom: 8px;
}

.conn-row {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #606266;
  margin-bottom: 4px;
}

.conn-label {
  color: #909399;
  min-width: 40px;
}

.conn-link {
  color: #409eff;
  text-decoration: none;
}

.conn-link:hover {
  text-decoration: underline;
}

.conn-empty {
  font-size: 13px;
  color: #c0c4cc;
}

.header-actions {
  display: flex;
  gap: 8px;
}

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
