<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, Setting } from '@element-plus/icons-vue'
import { useRoute } from 'vue-router'
import PageHeader from '@/components/common/PageHeader.vue'
import ConfigurableTable from '@/components/common/ConfigurableTable.vue'
import { customerTaskColumns } from '@/config/tableColumns'
import type { TableColumnConfig } from '@/config/tableColumns'
import type { Task, TaskLogResponse, TaskResultResponse } from '@/types/task'
import type { Machine } from '@/types/machine'
import {
  getTasks,
  getTaskDetail,
  createTrainingTask,
  stopTask,
  cancelTask,
  retryTask,
  getTaskLogs,
  getTaskResult,
  downloadTaskResult,
  getMyMachines,
} from '@/api/customer'

const tasks = ref<Task[]>([])
const total = ref(0)
const loading = ref(false)
const pageRequest = ref({
  page: 1,
  pageSize: 20,
})

const filters = ref({
  keyword: '',
  status: '',
  host_id: '',
})

const route = useRoute()

const viewMode = computed(() => {
  const path = route.path
  if (path.includes('/customer/tasks/training')) return 'training'
  if (path.includes('/customer/tasks/inference')) return 'inference'
  if (path.includes('/customer/tasks/queue')) return 'queue'
  if (path.includes('/customer/tasks/history')) return 'history'
  return 'all'
})

const viewTitle = computed(() => {
  switch (viewMode.value) {
    case 'training':
      return '训练任务'
    case 'inference':
      return '推理任务'
    case 'queue':
      return '任务队列'
    case 'history':
      return '历史记录'
    default:
      return '任务管理'
  }
})

const applyRouteDefaults = (mode: string) => {
  filters.value.keyword = ''
  filters.value.host_id = ''
  switch (mode) {
    case 'training':
      filters.value.status = ''
      break
    case 'inference':
      filters.value.status = ''
      break
    case 'queue':
      filters.value.status = 'queued'
      break
    case 'history':
      filters.value.status = 'completed'
      break
    default:
      filters.value.status = ''
  }
}

const columnStorageKey = 'customer-task-columns'
const columns = ref<TableColumnConfig[]>([])

const visibleColumnKeys = computed({
  get: () => columns.value.filter(col => !col.hidden).map(col => col.prop),
  set: (keys: string[]) => {
    columns.value = columns.value.map((col) => {
      if (col.prop === 'actions') {
        return { ...col, hidden: false }
      }
      return { ...col, hidden: !keys.includes(col.prop) }
    })
    persistColumnSettings()
  },
})

const initColumns = () => {
  const defaults = customerTaskColumns.map(col => ({ ...col }))
  const raw = localStorage.getItem(columnStorageKey)
  if (!raw) {
    columns.value = defaults
    return
  }
  try {
    const saved = JSON.parse(raw) as Record<string, boolean>
    columns.value = defaults.map(col => ({
      ...col,
      hidden: saved[col.prop] ?? col.hidden,
    }))
  } catch {
    columns.value = defaults
  }
}

const persistColumnSettings = () => {
  const payload: Record<string, boolean> = {}
  columns.value.forEach((col) => {
    payload[col.prop] = !!col.hidden
  })
  localStorage.setItem(columnStorageKey, JSON.stringify(payload))
}

const resetColumnSettings = () => {
  columns.value = customerTaskColumns.map(col => ({ ...col }))
  persistColumnSettings()
}

const statusOptions = [
  { label: '待处理', value: 'pending' },
  { label: '排队中', value: 'queued' },
  { label: '已分配', value: 'assigned' },
  { label: '运行中', value: 'running' },
  { label: '已完成', value: 'completed' },
  { label: '失败', value: 'failed' },
  { label: '已取消', value: 'cancelled' },
  { label: '已停止', value: 'stopped' },
]

const statusTagType = (status: string) => {
  switch (status) {
    case 'running':
      return 'success'
    case 'completed':
      return 'success'
    case 'failed':
      return 'danger'
    case 'cancelled':
    case 'stopped':
      return 'info'
    case 'pending':
    case 'queued':
    case 'assigned':
      return 'warning'
    default:
      return 'info'
  }
}

const statusLabel = (status: string) => {
  const option = statusOptions.find(item => item.value === status)
  return option?.label || status
}

const formatDate = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

const buildParams = () => {
  return {
    page: pageRequest.value.page,
    pageSize: pageRequest.value.pageSize,
    keyword: filters.value.keyword || undefined,
    status: filters.value.status || undefined,
    host_id: filters.value.host_id || undefined,
  }
}

const loadTasks = async () => {
  loading.value = true
  try {
    const res = await getTasks(buildParams())
    tasks.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('加载任务列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pageRequest.value.page = 1
  loadTasks()
}

const handleReset = () => {
  filters.value = {
    keyword: '',
    status: '',
    host_id: '',
  }
  pageRequest.value.page = 1
  loadTasks()
}

const handlePageChange = (page: number) => {
  pageRequest.value.page = page
  loadTasks()
}

const handleSizeChange = (size: number) => {
  pageRequest.value.pageSize = size
  pageRequest.value.page = 1
  loadTasks()
}

// 创建任务
const createDialogVisible = ref(false)
const createLoading = ref(false)
const createFormRef = ref<FormInstance>()
const createForm = ref({
  name: '',
  host_id: '',
  image_id: '',
  command: '',
  env_vars: '',
})

const createRules: FormRules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  host_id: [{ required: true, message: '请选择机器', trigger: 'change' }],
  command: [{ required: true, message: '请输入执行命令', trigger: 'blur' }],
}

const machines = ref<Machine[]>([])

const loadMachines = async () => {
  try {
    const res = await getMyMachines({ page: 1, pageSize: 200 })
    machines.value = res.data.list || []
  } catch (error) {
    console.error('加载机器列表失败:', error)
  }
}

const openCreateDialog = async () => {
  createDialogVisible.value = true
  await loadMachines()
}

const handleCreate = async () => {
  if (!createFormRef.value) return
  const valid = await createFormRef.value.validate()
  if (!valid) return
  createLoading.value = true
  try {
    const payload = {
      name: createForm.value.name,
      host_id: createForm.value.host_id,
      command: createForm.value.command,
      image_id: createForm.value.image_id ? Number(createForm.value.image_id) : undefined,
      env_vars: createForm.value.env_vars
        ? (JSON.parse(createForm.value.env_vars) as Record<string, string>)
        : undefined,
    }
    await createTrainingTask(payload)
    ElMessage.success('任务创建成功')
    createDialogVisible.value = false
    createFormRef.value.resetFields()
    loadTasks()
  } catch (error) {
    if (error instanceof SyntaxError) {
      ElMessage.error('环境变量 JSON 格式不正确')
    } else {
      console.error('创建任务失败:', error)
    }
  } finally {
    createLoading.value = false
  }
}

// 详情
const detailVisible = ref(false)
const detailLoading = ref(false)
const detailTask = ref<Task | null>(null)

const handleViewDetail = async (task: Task) => {
  detailVisible.value = true
  detailTask.value = task
  detailLoading.value = true
  try {
    const res = await getTaskDetail(task.id)
    detailTask.value = res.data
  } catch (error) {
    console.error('获取任务详情失败:', error)
  } finally {
    detailLoading.value = false
  }
}

// 日志
const logDialogVisible = ref(false)
const logLoading = ref(false)
const logTab = ref('stdout')
const logData = ref<TaskLogResponse>({})

const handleViewLogs = async (task: Task) => {
  logDialogVisible.value = true
  logLoading.value = true
  try {
    const res = await getTaskLogs(task.id)
    logData.value = res.data || {}
    if (logData.value.stderr) {
      logTab.value = 'stderr'
    } else if (logData.value.logs) {
      logTab.value = 'combined'
    } else {
      logTab.value = 'stdout'
    }
  } catch (error) {
    console.error('获取任务日志失败:', error)
  } finally {
    logLoading.value = false
  }
}

// 结果下载
const handleDownloadResult = async (task: Task) => {
  try {
    const res = await getTaskResult(task.id)
    const data = res.data as TaskResultResponse
    if (data?.presigned_url || data?.url) {
      window.open(data.presigned_url || data.url, '_blank')
      return
    }
    const blob = await downloadTaskResult(task.id)
    triggerDownload(blob, data?.filename || `task-${task.id}-result`)
  } catch (error) {
    try {
      const blob = await downloadTaskResult(task.id)
      triggerDownload(blob, `task-${task.id}-result`)
    } catch (downloadError) {
      console.error('下载结果失败:', downloadError)
    }
  }
}

const triggerDownload = (blob: Blob, filename: string) => {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(url)
}

// 操作
const handleStop = async (task: Task) => {
  try {
    await ElMessageBox.confirm('确定停止该任务吗？', '提示', { type: 'warning' })
    await stopTask(task.id)
    ElMessage.success('任务已停止')
    loadTasks()
  } catch {
    //
  }
}

const handleCancel = async (task: Task) => {
  try {
    await ElMessageBox.confirm('确定取消该任务吗？', '提示', { type: 'warning' })
    await cancelTask(task.id)
    ElMessage.success('任务已取消')
    loadTasks()
  } catch {
    //
  }
}

const handleRetry = async (task: Task) => {
  try {
    await ElMessageBox.confirm('确定重试该任务吗？', '提示', { type: 'warning' })
    await retryTask(task.id)
    ElMessage.success('任务已重试')
    loadTasks()
  } catch {
    //
  }
}

const canStop = (status: string) => ['running'].includes(status)
const canCancel = (status: string) => ['pending', 'queued', 'assigned'].includes(status)
const canRetry = (status: string) => ['failed', 'cancelled', 'stopped', 'preempted'].includes(status)

const formatEnvVars = (env: Task['env_vars']) => {
  if (!env) return '-'
  if (typeof env === 'string') {
    try {
      return JSON.stringify(JSON.parse(env), null, 2)
    } catch {
      return env
    }
  }
  return JSON.stringify(env, null, 2)
}

// 轮询：当列表中存在运行中/待处理的任务时，自动刷新
const hasActiveTasks = computed(() =>
  tasks.value.some(t => ['pending', 'queued', 'assigned', 'running'].includes(t.status))
)

let pollTimer: ReturnType<typeof setInterval> | null = null

const startPolling = () => {
  stopPolling()
  pollTimer = setInterval(() => {
    if (hasActiveTasks.value) {
      loadTasks()
    }
  }, 5000)
}

const stopPolling = () => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

onMounted(() => {
  initColumns()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})

watch(viewMode, (mode) => {
  applyRouteDefaults(mode)
  pageRequest.value.page = 1
  loadTasks()
}, { immediate: true })
</script>

<template>
  <div class="task-list-view">
    <PageHeader :title="viewTitle" subtitle="管理我的任务执行">
      <template #actions>
        <el-popover placement="bottom-end" width="220" trigger="click">
          <div class="column-settings">
            <div class="column-settings__header">
              <span>列显示</span>
              <el-button link type="primary" size="small" @click="resetColumnSettings">重置</el-button>
            </div>
            <el-checkbox-group v-model="visibleColumnKeys">
              <el-checkbox
                v-for="column in columns"
                :key="column.prop"
                :label="column.prop"
                :disabled="column.prop === 'actions'"
              >
                {{ column.label }}
              </el-checkbox>
            </el-checkbox-group>
          </div>
          <template #reference>
            <el-button :icon="Setting">列配置</el-button>
          </template>
        </el-popover>
        <el-button :icon="Refresh" @click="loadTasks">刷新</el-button>
        <el-button type="primary" :icon="Plus" @click="openCreateDialog">创建任务</el-button>
      </template>
    </PageHeader>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters">
        <el-form-item label="关键词">
          <el-input v-model="filters.keyword" placeholder="任务ID/名称/命令" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 140px">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="机器">
          <el-input v-model="filters.host_id" placeholder="机器ID" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <ConfigurableTable :columns="columns" :data="tasks" :loading="loading">
        <template #name="{ row }">
          <el-link type="primary" @click="handleViewDetail(row)">
            {{ row.name || row.id }}
          </el-link>
        </template>
        <template #status="{ row }">
          <el-tag :type="statusTagType(row.status)">
            {{ statusLabel(row.status) }}
          </el-tag>
        </template>
        <template #host="{ row }">
          <span>{{ row.host_id || '-' }}</span>
        </template>
        <template #image="{ row }">
          <span>{{ row.image?.name || row.image_id || '-' }}</span>
        </template>
        <template #command="{ row }">
          <el-tooltip v-if="row.command" :content="row.command" placement="top">
            <span class="truncate">{{ row.command }}</span>
          </el-tooltip>
          <span v-else>-</span>
        </template>
        <template #created_at="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
        <template #started_at="{ row }">
          {{ formatDate(row.started_at) }}
        </template>
        <template #ended_at="{ row }">
          {{ formatDate(row.ended_at) }}
        </template>
        <template #actions="{ row }">
          <el-button link type="primary" size="small" @click="handleViewDetail(row)">详情</el-button>
          <el-button link type="primary" size="small" @click="handleViewLogs(row)">日志</el-button>
          <el-button link type="primary" size="small" @click="handleDownloadResult(row)">结果下载</el-button>
          <el-dropdown>
            <el-button link type="primary" size="small">更多</el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item :disabled="!canStop(row.status)" @click="handleStop(row)">停止</el-dropdown-item>
                <el-dropdown-item :disabled="!canCancel(row.status)" @click="handleCancel(row)">取消</el-dropdown-item>
                <el-dropdown-item :disabled="!canRetry(row.status)" @click="handleRetry(row)">重试</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </ConfigurableTable>

      <el-pagination
        class="pagination"
        :current-page="pageRequest.page"
        :page-size="pageRequest.pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </el-card>

    <el-dialog v-model="createDialogVisible" title="创建任务" width="680px" :close-on-click-modal="false">
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="110px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="机器" prop="host_id">
          <el-select v-model="createForm.host_id" placeholder="请选择机器" filterable>
            <el-option
              v-for="machine in machines"
              :key="machine.id"
              :label="machine.name || machine.hostname || machine.id"
              :value="machine.id.toString()"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="镜像ID">
          <el-input v-model="createForm.image_id" placeholder="可选" />
        </el-form-item>
        <el-form-item label="执行命令" prop="command">
          <el-input v-model="createForm.command" type="textarea" :rows="3" placeholder="如：python train.py" />
        </el-form-item>
        <el-form-item label="环境变量">
          <el-input
            v-model="createForm.env_vars"
            type="textarea"
            :rows="3"
            placeholder='JSON格式，如：{"CUDA_VISIBLE_DEVICES":"0"}'
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="createLoading" @click="handleCreate">创建</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="detailVisible" title="任务详情" size="480px">
      <div v-loading="detailLoading" class="detail-content">
        <el-descriptions v-if="detailTask" :column="1" border>
          <el-descriptions-item label="任务ID">{{ detailTask.id }}</el-descriptions-item>
          <el-descriptions-item label="任务名称">{{ detailTask.name }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="statusTagType(detailTask.status)">
              {{ statusLabel(detailTask.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="机器ID">{{ detailTask.host_id || '-' }}</el-descriptions-item>
          <el-descriptions-item label="镜像">
            {{ detailTask.image?.name || detailTask.image_id || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="命令">
            <span class="command-text">{{ detailTask.command || '-' }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="环境变量">
            <pre class="json-pre">{{ formatEnvVars(detailTask.env_vars) }}</pre>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(detailTask.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="开始时间">{{ formatDate(detailTask.started_at) }}</el-descriptions-item>
          <el-descriptions-item label="结束时间">{{ formatDate(detailTask.ended_at) }}</el-descriptions-item>
          <el-descriptions-item label="退出码">{{ detailTask.exit_code ?? '-' }}</el-descriptions-item>
          <el-descriptions-item label="错误信息">{{ detailTask.error_msg || '-' }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-drawer>

    <el-dialog v-model="logDialogVisible" title="任务日志" width="820px">
      <div v-loading="logLoading" class="log-content">
        <el-tabs v-model="logTab">
          <el-tab-pane label="标准输出" name="stdout">
            <pre class="log-pre">{{ logData.stdout || '暂无输出' }}</pre>
          </el-tab-pane>
          <el-tab-pane label="标准错误" name="stderr">
            <pre class="log-pre">{{ logData.stderr || '暂无输出' }}</pre>
          </el-tab-pane>
          <el-tab-pane label="合并日志" name="combined">
            <pre class="log-pre">{{ logData.logs || '暂无输出' }}</pre>
          </el-tab-pane>
        </el-tabs>
      </div>
      <template #footer>
        <el-button @click="logDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.task-list-view {
  padding: 24px;
}

.filter-card {
  margin-bottom: 16px;
}

.table-card {
  padding-bottom: 16px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.truncate {
  display: inline-block;
  max-width: 220px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.column-settings {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.column-settings__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.detail-content {
  padding: 8px 0;
}

.command-text {
  word-break: break-all;
}

.json-pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 12px;
  background: #f5f7fa;
  padding: 8px;
  border-radius: 4px;
}

.log-content {
  min-height: 240px;
}

.log-pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 12px;
  background: #0f172a;
  color: #e2e8f0;
  padding: 12px;
  border-radius: 6px;
  min-height: 200px;
}
</style>
