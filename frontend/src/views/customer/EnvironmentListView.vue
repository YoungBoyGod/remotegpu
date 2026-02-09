<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import { useRoleNavigation } from '@/composables/useRoleNavigation'
import PageHeader from '@/components/common/PageHeader.vue'
import FilterBar from '@/components/common/FilterBar.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import type { Environment } from '@/api/environment/types'
import {
  getEnvironmentList,
  startEnvironment as startEnv,
  stopEnvironment as stopEnv,
  deleteEnvironment as deleteEnv,
  getEnvironmentAccessInfo,
} from '@/api/environment'

const { navigateTo } = useRoleNavigation()

const environments = ref<Environment[]>([])
const accessInfoMap = ref<Record<string, any>>({})
const loading = ref(false)
const searchText = ref('')
const statusFilter = ref('')

// 分页相关
const currentPage = ref(1)
const pageSize = ref(5)

const formatDateTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

const formatMemoryToGB = (value: number) => {
  if (!value) return 0
  const gb = value / 1024
  return gb < 1 ? 1 : Math.round(gb)
}

const formatDuration = (startAt?: string | null) => {
  if (!startAt) return '-'
  const start = new Date(startAt)
  if (Number.isNaN(start.getTime())) return '-'
  const diffMs = Date.now() - start.getTime()
  if (diffMs < 0) return '-'
  const diffMinutes = Math.floor(diffMs / 60000)
  const hours = Math.floor(diffMinutes / 60)
  const minutes = diffMinutes % 60
  if (hours > 0) {
    return `${hours}小时${minutes}分`
  }
  return `${minutes}分`
}

const displayEnvironments = computed(() => {
  return environments.value.map(env => ({
    id: env.id,
    name: env.name,
    status: env.status,
    image: env.image,
    gpu: env.gpu,
    cpu: env.cpu,
    memory: formatMemoryToGB(env.memory),
    runningTime: env.status === 'running' ? formatDuration(env.started_at) : '-',
    createdAt: formatDateTime(env.created_at),
    sshPort: env.ssh_port,
    rdpPort: env.rdp_port,
    jupyterPort: env.jupyter_port,
    accessInfo: accessInfoMap.value[env.id],
  }))
})

const statusTextMap: Record<string, string> = {
  running: '运行中',
  stopped: '已停止',
  creating: '创建中',
  deleting: '删除中',
  error: '错误',
}

// 过滤后的环境列表
const filteredEnvironments = computed(() => {
  let result = displayEnvironments.value

  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(env =>
      env.name.toLowerCase().includes(search) ||
      env.image.toLowerCase().includes(search)
    )
  }

  if (statusFilter.value) {
    result = result.filter(env => env.status === statusFilter.value)
  }

  return result
})

// 分页后的环境列表
const paginatedEnvironments = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredEnvironments.value.slice(start, end)
})

// 总数
const total = computed(() => filteredEnvironments.value.length)

// 获取环境访问信息
const loadAccessInfo = async (envId: string) => {
  try {
    const response = await getEnvironmentAccessInfo(envId)
    if (response.data && response.data.access_info) {
      accessInfoMap.value[envId] = response.data.access_info
    }
  } catch (error) {
    console.error(`获取环境 ${envId} 访问信息失败:`, error)
  }
}

// 加载环境列表
const loadEnvironments = async () => {
  loading.value = true
  try {
    const response = await getEnvironmentList()
    environments.value = response.data

    // 为运行中的环境加载访问信息
    const runningEnvs = environments.value.filter(env => env.status === 'running')
    await Promise.all(runningEnvs.map(env => loadAccessInfo(env.id)))
  } catch (error) {
    ElMessage.error('加载环境列表失败')
  } finally {
    loading.value = false
  }
}

// 启动环境
const startEnvironment = async (id: string) => {
  try {
    await startEnv(id)
    ElMessage.success('环境启动中...')
    await loadEnvironments()
  } catch (error) {
    ElMessage.error('启动失败')
  }
}

// 停止环境
const stopEnvironment = async (id: string) => {
  try {
    await stopEnv(id)
    ElMessage.success('环境已停止')
    await loadEnvironments()
  } catch (error) {
    ElMessage.error('停止失败')
  }
}

// 删除环境
const deleteEnvironment = async (id: string) => {
  try {
    await ElMessageBox.confirm('确定要删除这个环境吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await deleteEnv(id)
    ElMessage.success('环境已删除')
    await loadEnvironments()
  } catch (error) {
    // 用户取消
  }
}

onMounted(() => {
  loadEnvironments()
})
</script>

<template>
  <div class="environment-list">
    <PageHeader title="开发环境">
      <template #actions>
        <el-button type="primary" :icon="Plus" @click="navigateTo('/environments/create')">
          创建环境
        </el-button>
      </template>
    </PageHeader>

    <FilterBar
      v-model:search-value="searchText"
      search-placeholder="搜索环境名称"
    >
      <template #filters>
        <el-select v-model="statusFilter" placeholder="状态筛选" style="width: 150px" clearable>
          <el-option label="运行中" value="running" />
          <el-option label="已停止" value="stopped" />
          <el-option label="错误" value="error" />
        </el-select>
      </template>
      <template #actions>
        <el-button :icon="Refresh" @click="loadEnvironments">刷新</el-button>
      </template>
    </FilterBar>

    <!-- 卡片网格布局 -->
    <div v-loading="loading" class="environment-grid">
      <el-card v-for="env in paginatedEnvironments" :key="env.id" class="environment-card">
        <!-- 卡片头部 -->
        <template #header>
          <div class="card-header">
            <el-link type="primary" @click="navigateTo(`/environments/${env.id}`)" class="env-name">
              {{ env.name }}
            </el-link>
            <StatusTag :status="statusTextMap[env.status] || env.status" />
          </div>
        </template>

        <!-- 主机配置信息 -->
        <div class="config-section">
          <div class="config-item">
            <span class="config-label">镜像:</span>
            <span class="config-value">{{ env.image }}</span>
          </div>
          <div class="config-item">
            <span class="config-label">GPU:</span>
            <span class="config-value">{{ env.gpu }}</span>
          </div>
          <div class="config-item">
            <span class="config-label">CPU:</span>
            <span class="config-value">{{ env.cpu }}核</span>
          </div>
          <div class="config-item">
            <span class="config-label">内存:</span>
            <span class="config-value">{{ env.memory }}GB</span>
          </div>
          <div class="config-item">
            <span class="config-label">运行时长:</span>
            <span class="config-value">{{ env.runningTime }}</span>
          </div>
          <div class="config-item">
            <span class="config-label">创建时间:</span>
            <span class="config-value">{{ env.createdAt }}</span>
          </div>
        </div>

        <!-- 连接说明 -->
        <div v-if="env.status === 'running' && env.accessInfo" class="connection-section">
          <div class="connection-title">连接方式:</div>
          <div class="connection-methods">
            <!-- SSH 连接 -->
            <div v-if="env.accessInfo.ssh" class="connection-item">
              <el-tag type="success" size="small">SSH</el-tag>
              <div class="connection-details">
                <div>主机: {{ env.accessInfo.ssh.host }}:{{ env.accessInfo.ssh.port }}</div>
                <div>用户: {{ env.accessInfo.ssh.username }}</div>
                <div>密码: {{ env.accessInfo.ssh.password }}</div>
                <div class="connection-command">{{ env.accessInfo.ssh.command }}</div>
              </div>
            </div>

            <!-- RDP 连接 -->
            <div v-if="env.accessInfo.rdp" class="connection-item">
              <el-tag type="primary" size="small">RDP</el-tag>
              <div class="connection-details">
                <div>主机: {{ env.accessInfo.rdp.host }}:{{ env.accessInfo.rdp.port }}</div>
                <div>用户: {{ env.accessInfo.rdp.username }}</div>
                <div>密码: {{ env.accessInfo.rdp.password }}</div>
                <div class="connection-command">{{ env.accessInfo.rdp.command }}</div>
              </div>
            </div>

            <!-- Jupyter 连接 -->
            <div v-if="env.accessInfo.jupyter" class="connection-item">
              <el-tag type="warning" size="small">Jupyter</el-tag>
              <div class="connection-details">
                <div>URL: {{ env.accessInfo.jupyter.url }}</div>
                <div>Token: {{ env.accessInfo.jupyter.token }}</div>
              </div>
            </div>

            <!-- VNC 连接 -->
            <div v-if="env.accessInfo.vnc" class="connection-item">
              <el-tag type="info" size="small">VNC</el-tag>
              <div class="connection-details">
                <div>主机: {{ env.accessInfo.vnc.host }}:{{ env.accessInfo.vnc.port }}</div>
                <div>密码: {{ env.accessInfo.vnc.password }}</div>
                <div class="connection-command">{{ env.accessInfo.vnc.url }}</div>
              </div>
            </div>
          </div>
        </div>
        <div v-else-if="env.status === 'running' && !env.accessInfo" class="connection-section">
          <div class="connection-title">正在加载连接信息...</div>
        </div>

        <!-- 操作按钮 -->
        <div class="card-actions">
          <el-button
            v-if="env.status === 'stopped'"
            type="success"
            size="small"
            @click="startEnvironment(env.id)"
          >
            启动
          </el-button>
          <el-button
            v-if="env.status === 'running'"
            type="warning"
            size="small"
            @click="stopEnvironment(env.id)"
          >
            停止
          </el-button>
          <el-button type="danger" size="small" @click="deleteEnvironment(env.id)">
            删除
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- 分页 -->
    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="currentPage = $event"
      />
    </div>
  </div>
</template>

<style scoped>
.environment-list {
  padding: 24px;
  background: #f5f7fa;
  min-height: 100%;
}

.environment-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.environment-card {
  border-radius: 8px;
  border: none;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  transition: all 0.2s;
}

.environment-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.env-name {
  font-size: 15px;
  font-weight: 600;
}

.config-section {
  margin-bottom: 16px;
}

.config-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.config-item:last-child {
  border-bottom: none;
}

.config-label {
  color: #86909c;
  font-size: 13px;
}

.config-value {
  color: #1d2129;
  font-size: 13px;
  font-weight: 500;
}

.connection-section {
  background: #f5f7fa;
  padding: 12px;
  border-radius: 4px;
  margin-bottom: 16px;
}

.connection-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.connection-methods {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.connection-item {
  background: white;
  padding: 10px;
  border-radius: 4px;
  border: 1px solid #e4e7ed;
}

.connection-details {
  margin-top: 8px;
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
}

.connection-details > div {
  margin-bottom: 4px;
}

.connection-command {
  margin-top: 8px;
  padding: 8px;
  background: #f5f7fa;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #409eff;
  word-break: break-all;
}

.card-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}
</style>
