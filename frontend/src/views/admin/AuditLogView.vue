<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Search, Refresh } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { getAuditLogs } from '@/api/admin'

interface AuditLog {
  id: number
  customer_id?: number
  username: string
  action: string
  resource_type: string
  resource_id: string
  ip_address: string
  method: string
  path: string
  status_code: number
  created_at: string
  detail: any
}

const loading = ref(false)
const logs = ref<AuditLog[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

// 筛选条件
const filters = ref({
  username: '',
  action: '',
  resource_type: '',
  time_range: null as [string, string] | null,
})

// 操作类型选项（与后端 middleware/audit.go parseAction 一致）
const actionOptions = [
  { label: '创建', value: 'create' },
  { label: '更新', value: 'update' },
  { label: '删除', value: 'delete' },
  { label: '分配', value: 'allocate' },
  { label: '回收', value: 'reclaim' },
  { label: '禁用', value: 'disable' },
  { label: '启用', value: 'enable' },
  { label: '停止', value: 'stop' },
  { label: '取消', value: 'cancel' },
  { label: '重试', value: 'retry' },
  { label: '采集', value: 'collect' },
  { label: '维护', value: 'maintenance' },
  { label: '确认', value: 'acknowledge' },
  { label: '同步', value: 'sync' },
  { label: '导入', value: 'import' },
]

// 资源类型选项（与后端 middleware/audit.go parseResource 一致）
const resourceTypeOptions = [
  { label: '机器', value: 'machine' },
  { label: '客户', value: 'customer' },
  { label: '分配记录', value: 'allocation' },
  { label: '任务', value: 'task' },
  { label: '数据集', value: 'dataset' },
  { label: '镜像', value: 'image' },
  { label: '文档', value: 'document' },
  { label: '告警', value: 'alert' },
  { label: 'SSH 密钥', value: 'ssh_key' },
  { label: '存储', value: 'storage' },
  { label: '系统配置', value: 'system_config' },
]

const loadLogs = async () => {
  loading.value = true
  try {
    const res = await getAuditLogs({
      page: page.value,
      pageSize: pageSize.value,
      username: filters.value.username || undefined,
      action: filters.value.action || undefined,
      resource_type: filters.value.resource_type || undefined,
      start_time: filters.value.time_range?.[0] || undefined,
      end_time: filters.value.time_range?.[1] || undefined,
    })
    logs.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('加载审计日志失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  page.value = 1
  loadLogs()
}

const handleReset = () => {
  filters.value = {
    username: '',
    action: '',
    resource_type: '',
    time_range: null,
  }
  page.value = 1
  loadLogs()
}

const handlePageChange = (val: number) => {
  page.value = val
  loadLogs()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  page.value = 1
  loadLogs()
}

const formatDate = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN')
}

// 操作类型标签颜色
const actionTagType = (action: string) => {
  const map: Record<string, string> = {
    create: 'success',
    update: 'warning',
    delete: 'danger',
    allocate: 'success',
    reclaim: 'warning',
    disable: 'danger',
    enable: 'success',
    stop: 'danger',
    cancel: 'info',
    retry: 'warning',
    collect: '',
    maintenance: 'info',
    acknowledge: 'success',
    sync: '',
    import: '',
  }
  return map[action] || 'info'
}

// 操作类型显示名称
const actionLabel = (action: string) => {
  const opt = actionOptions.find(o => o.value === action)
  return opt?.label || action
}

// 资源类型显示名称
const resourceTypeLabel = (type: string) => {
  const opt = resourceTypeOptions.find(o => o.value === type)
  return opt?.label || type
}

// HTTP 方法标签颜色
const methodTagType = (method: string) => {
  const map: Record<string, string> = {
    GET: 'info',
    POST: 'success',
    PUT: 'warning',
    PATCH: 'warning',
    DELETE: 'danger',
  }
  return map[method?.toUpperCase()] || 'info'
}

onMounted(() => {
  loadLogs()
})
</script>

<template>
  <div class="audit-log-view">
    <PageHeader title="审计日志" />

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters">
        <el-form-item label="用户名">
          <el-input v-model="filters.username" placeholder="请输入用户名" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="操作类型">
          <el-select v-model="filters.action" placeholder="全部" clearable style="width: 130px">
            <el-option v-for="opt in actionOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="资源类型">
          <el-select v-model="filters.resource_type" placeholder="全部" clearable style="width: 130px">
            <el-option v-for="opt in resourceTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="filters.time_range"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
            style="width: 340px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table :data="logs" v-loading="loading" stripe style="width: 100%">
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="detail-expand">
              <el-descriptions :column="2" border size="small">
                <el-descriptions-item label="请求方法">
                  <el-tag size="small" :type="methodTagType(row.method)">{{ row.method }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="请求路径">{{ row.path || '-' }}</el-descriptions-item>
                <el-descriptions-item label="资源ID">{{ row.resource_id || '-' }}</el-descriptions-item>
                <el-descriptions-item label="客户ID">{{ row.customer_id || '-' }}</el-descriptions-item>
              </el-descriptions>
              <div v-if="row.detail" class="detail-json">
                <span class="detail-json-label">详细信息：</span>
                <pre class="detail-pre">{{ JSON.stringify(row.detail, null, 2) }}</pre>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户" width="120" />
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="actionTagType(row.action)">
              {{ actionLabel(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="资源类型" width="110">
          <template #default="{ row }">
            {{ resourceTypeLabel(row.resource_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="resource_id" label="资源ID" min-width="140" show-overflow-tooltip />
        <el-table-column prop="ip_address" label="IP 地址" width="140" />
        <el-table-column label="状态码" width="90">
          <template #default="{ row }">
            <el-tag size="small" :type="row.status_code >= 200 && row.status_code < 300 ? 'success' : 'danger'">
              {{ row.status_code }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="时间" width="175">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </el-card>
  </div>
</template>

<style scoped>
.audit-log-view {
  padding: 24px;
  background: #f5f7fa;
  min-height: 100%;
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

.detail-expand {
  padding: 16px 24px;
}

.detail-json {
  margin-top: 12px;
}

.detail-json-label {
  font-size: 13px;
  font-weight: 600;
  color: #606266;
}

.detail-pre {
  margin: 8px 0 0 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 12px;
  background: #f5f7fa;
  padding: 12px;
  border-radius: 6px;
  color: #303133;
  max-height: 300px;
  overflow-y: auto;
}
</style>
