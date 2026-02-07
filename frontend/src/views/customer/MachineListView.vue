<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { CopyDocument, View, Hide } from '@element-plus/icons-vue'
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

// 连接信息对话框
const connectionDialogVisible = ref(false)
const connectionLoading = ref(false)
const connectionTab = ref('ssh')
const connectionData = ref<any>(null)
const connectionMachine = ref<Machine | null>(null)
const showPassword = ref(false)
const showVncPassword = ref(false)

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

const getDeviceStatusType = (status?: string) => {
  return status === 'online' ? 'success' : 'danger'
}

const getDeviceStatusText = (status?: string) => {
  return status === 'online' ? '在线' : '离线'
}

const formatDate = (dateStr: string | undefined) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const formatGPUSize = (mb?: number | null) => {
  if (!mb) return '-'
  return `${Math.round(mb / 1024)} GB`
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

// 打开连接信息对话框
const handleConnect = async (machine: Machine) => {
  connectionMachine.value = machine
  connectionData.value = null
  connectionTab.value = 'ssh'
  showPassword.value = false
  showVncPassword.value = false
  connectionDialogVisible.value = true

  try {
    connectionLoading.value = true
    const response = await getMachineConnection(machine.id)
    connectionData.value = response.data
  } catch (error: any) {
    ElMessage.error(error?.msg || error?.message || '获取连接信息失败')
  } finally {
    connectionLoading.value = false
  }
}

// 密码遮蔽
const maskPassword = (pwd?: string) => {
  if (!pwd) return ''
  if (pwd.length <= 6) return pwd.slice(0, 2) + '...'
  return pwd.slice(0, 6) + '...'
}

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

// 从连接数据或机器数据中获取 SSH 信息
const sshHost = computed(() => {
  const conn = connectionData.value
  if (conn?.ssh?.host) return conn.ssh.host
  const m = connectionMachine.value
  return m?.ssh_host || m?.public_ip || m?.ip_address || ''
})

const sshPort = computed(() => {
  const conn = connectionData.value
  if (conn?.ssh?.port) return conn.ssh.port
  return connectionMachine.value?.ssh_port || 22
})

const sshUsername = computed(() => {
  const conn = connectionData.value
  if (conn?.ssh?.username) return conn.ssh.username
  return connectionMachine.value?.ssh_username || 'root'
})

const sshPassword = computed(() => {
  const conn = connectionData.value
  if (conn?.ssh?.password) return conn.ssh.password
  return connectionMachine.value?.ssh_password || ''
})

const sshCommand = computed(() => {
  const conn = connectionData.value
  if (conn?.ssh?.command) return conn.ssh.command
  if (connectionMachine.value?.ssh_command) return connectionMachine.value.ssh_command
  if (!sshHost.value) return ''
  return `ssh ${sshUsername.value}@${sshHost.value} -p ${sshPort.value}`
})

const jupyterUrl = computed(() => {
  const conn = connectionData.value
  return conn?.jupyter?.url || connectionMachine.value?.jupyter_url || ''
})

const jupyterToken = computed(() => {
  const conn = connectionData.value
  return conn?.jupyter?.token || connectionMachine.value?.jupyter_token || ''
})

const vncUrl = computed(() => {
  const conn = connectionData.value
  return conn?.vnc?.url || connectionMachine.value?.vnc_url || ''
})

const vncPassword = computed(() => {
  const conn = connectionData.value
  return conn?.vnc?.password || connectionMachine.value?.vnc_password || ''
})

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
      :total="total"
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
            {{ row.hostname || row.name || row.id }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column prop="ip_address" label="IP 地址" width="140" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="getDeviceStatusType(row.device_status)">
            {{ getDeviceStatusText(row.device_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="GPU" min-width="160">
        <template #default="{ row }">
          <template v-if="row.gpus && row.gpus.length > 0">
            {{ row.gpus[0].name }} x{{ row.gpus.length }}
            <span class="gpu-mem">{{ formatGPUSize(row.gpus[0].memory_total_mb) }}</span>
          </template>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="CPU / 内存" width="130">
        <template #default="{ row }">
          {{ row.total_cpu || '-' }} 核 / {{ row.total_memory_gb ? row.total_memory_gb + ' GB' : '-' }}
        </template>
      </el-table-column>
      <el-table-column label="到期时间" width="170">
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

    <!-- 连接信息对话框 -->
    <el-dialog
      v-model="connectionDialogVisible"
      :title="'连接信息 - ' + (connectionMachine?.hostname || connectionMachine?.name || '')"
      width="560px"
      destroy-on-close
    >
      <div v-loading="connectionLoading">
        <el-tabs v-model="connectionTab">
          <!-- SSH -->
          <el-tab-pane label="SSH" name="ssh">
            <div class="connection-info">
              <div class="info-item">
                <span class="info-label">主机：</span>
                <span class="info-value">{{ sshHost || '-' }}</span>
                <el-button v-if="sshHost" link :icon="CopyDocument" @click="copyToClipboard(sshHost)" />
              </div>
              <div class="info-item">
                <span class="info-label">端口：</span>
                <span class="info-value">{{ sshPort }}</span>
                <el-button link :icon="CopyDocument" @click="copyToClipboard(String(sshPort))" />
              </div>
              <div class="info-item">
                <span class="info-label">用户：</span>
                <span class="info-value">{{ sshUsername }}</span>
                <el-button link :icon="CopyDocument" @click="copyToClipboard(sshUsername)" />
              </div>
              <div class="info-item" v-if="sshPassword">
                <span class="info-label">密码：</span>
                <span class="info-value">{{ showPassword ? sshPassword : maskPassword(sshPassword) }}</span>
                <el-button link :icon="showPassword ? Hide : View" @click="showPassword = !showPassword" />
                <el-button link :icon="CopyDocument" @click="copyToClipboard(sshPassword)" />
              </div>
              <div class="info-item" v-if="sshCommand">
                <span class="info-label">命令：</span>
                <code class="info-value ssh-cmd">{{ sshCommand }}</code>
                <el-button link :icon="CopyDocument" @click="copyToClipboard(sshCommand)" />
              </div>
            </div>
          </el-tab-pane>

          <!-- Jupyter -->
          <el-tab-pane label="Jupyter" name="jupyter">
            <div class="connection-info">
              <template v-if="jupyterUrl">
                <div class="info-item">
                  <span class="info-label">地址：</span>
                  <a class="info-value info-link" :href="jupyterUrl" target="_blank">{{ jupyterUrl }}</a>
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(jupyterUrl)" />
                </div>
                <div class="info-item" v-if="jupyterToken">
                  <span class="info-label">Token：</span>
                  <span class="info-value">{{ jupyterToken }}</span>
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(jupyterToken)" />
                </div>
              </template>
              <div v-else class="info-empty">暂未配置 Jupyter</div>
            </div>
          </el-tab-pane>

          <!-- VNC -->
          <el-tab-pane label="VNC" name="vnc">
            <div class="connection-info">
              <template v-if="vncUrl">
                <div class="info-item">
                  <span class="info-label">地址：</span>
                  <a class="info-value info-link" :href="vncUrl" target="_blank">{{ vncUrl }}</a>
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(vncUrl)" />
                </div>
                <div class="info-item" v-if="vncPassword">
                  <span class="info-label">密码：</span>
                  <span class="info-value">{{ showVncPassword ? vncPassword : maskPassword(vncPassword) }}</span>
                  <el-button link :icon="showVncPassword ? Hide : View" @click="showVncPassword = !showVncPassword" />
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(vncPassword)" />
                </div>
              </template>
              <div v-else class="info-empty">暂未配置 VNC</div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
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

.page-actions {
  display: flex;
  gap: 12px;
}

.filter-card {
  margin-bottom: 16px;
}

.gpu-mem {
  color: #909399;
  font-size: 12px;
  margin-left: 4px;
}

.connection-info {
  display: flex;
  flex-direction: column;
  background: #f9fafb;
  border-radius: 8px;
  padding: 4px 0;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  color: #909399;
  font-size: 14px;
  white-space: nowrap;
  min-width: 50px;
}

.info-value {
  color: #303133;
  font-size: 14px;
  word-break: break-all;
}

.ssh-cmd {
  font-size: 13px;
  color: #606266;
  background: #f5f7fa;
  padding: 4px 8px;
  border-radius: 4px;
}

.info-link {
  color: #409eff;
  text-decoration: none;
}

.info-link:hover {
  text-decoration: underline;
}

.info-empty {
  padding: 24px 16px;
  text-align: center;
  color: #909399;
  font-size: 14px;
}
</style>
