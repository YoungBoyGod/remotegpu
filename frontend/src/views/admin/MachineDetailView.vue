<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getMachineDetail, getMachineUsage } from '@/api/admin'
import type { Machine } from '@/types/machine'
import type { MachineUsage } from '@/api/admin'
import { ElMessage } from 'element-plus'
import { CopyDocument, View, Hide } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const machine = ref<Machine | null>(null)
const usage = ref<MachineUsage | null>(null)

const loadMachine = async () => {
  const id = route.params.id as string
  if (!id) {
    ElMessage.error('机器ID不存在')
    router.back()
    return
  }

  try {
    loading.value = true
    const response = await getMachineDetail(id)
    machine.value = response.data
  } catch (error) {
    console.error('加载机器详情失败:', error)
    ElMessage.error('加载机器详情失败')
  } finally {
    loading.value = false
  }
}

const loadUsage = async () => {
  const id = route.params.id as string
  if (!id) return
  try {
    const res = await getMachineUsage(id)
    usage.value = res.data
  } catch {
    // 使用情况加载失败不阻塞页面
  }
}

const usageColor = (percent: number) => {
  if (percent >= 90) return '#f56c6c'
  if (percent >= 70) return '#e6a23c'
  return '#67c23a'
}

const getStatusType = (status?: string) => {
  const map: Record<string, string> = {
    idle: 'success',
    allocated: 'warning',
    maintenance: 'info',
    offline: 'danger'
  }
  return map[status || ''] || 'info'
}

const getStatusText = (status?: string) => {
  const map: Record<string, string> = {
    idle: '空闲',
    allocated: '已分配',
    maintenance: '维护中',
    offline: '离线'
  }
  return map[status || ''] || status || '未知'
}

const formatDateTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

const formatSizeGB = (value?: number | null) => {
  if (!value) return '-'
  return `${value} GB`
}

const formatGPUSize = (mb?: number | null) => {
  if (!mb) return '-'
  return `${Math.round(mb / 1024)} GB`
}

const activeAllocation = computed(() => {
  return machine.value?.allocations?.find((item) => item.status === 'active') || null
})



const handleBack = () => {
  router.push('/admin/machines/list')
}

const showPassword = ref(false)
const showVncPassword = ref(false)
const connectionTab = ref('ssh')

const sshConnectHost = computed(() => {
  if (!machine.value) return ''
  return machine.value.ssh_host || machine.value.public_ip || machine.value.ip_address || ''
})

const sshCommand = computed(() => {
  if (!machine.value) return ''
  return machine.value.ssh_command || ''
})

const maskedPassword = computed(() => {
  const pwd = machine.value?.ssh_password
  if (!pwd) return ''
  if (pwd.length <= 6) return pwd.slice(0, 2) + '...'
  return pwd.slice(0, 6) + '...'
})

const maskedVncPassword = computed(() => {
  const pwd = machine.value?.vnc_password
  if (!pwd) return ''
  if (pwd.length <= 6) return pwd.slice(0, 2) + '...'
  return pwd.slice(0, 6) + '...'
})

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

onMounted(() => {
  loadMachine()
  loadUsage()
})
</script>

<template>
  <div class="machine-detail" v-loading="loading">
    <div class="page-header">
      <div class="header-left">
        <el-button @click="handleBack" :icon="'ArrowLeft'">返回列表</el-button>
        <h2 class="page-title">{{ machine?.name || machine?.hostname || machine?.id || '机器详情' }}</h2>
        <el-tag v-if="machine" :type="(machine as any).device_status === 'online' ? 'success' : 'danger'">
          {{ (machine as any).device_status === 'online' ? '在线' : '离线' }}
        </el-tag>
        <el-tag v-if="machine" :type="getStatusType((machine as any).allocation_status)">
          {{ getStatusText((machine as any).allocation_status) }}
        </el-tag>
        <el-tag v-if="machine" :type="machine.needs_collect ? 'warning' : 'success'">
          {{ machine.needs_collect ? '待采集' : '已采集' }}
        </el-tag>
      </div>
    </div>

    <template v-if="machine">
      <!-- 基本信息 -->
      <el-card class="info-card">
        <template #header>
          <span>基本信息</span>
        </template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="ID">{{ machine.id }}</el-descriptions-item>
          <el-descriptions-item label="名称">{{ machine.name }}</el-descriptions-item>
          <el-descriptions-item label="主机名">{{ machine.hostname }}</el-descriptions-item>
          <el-descriptions-item label="区域">{{ machine.region }}</el-descriptions-item>
          <el-descriptions-item label="部署模式">{{ machine.deployment_mode || '-' }}</el-descriptions-item>
          <el-descriptions-item label="设备状态">
            <el-tag :type="(machine as any).device_status === 'online' ? 'success' : 'danger'">
              {{ (machine as any).device_status === 'online' ? '在线' : '离线' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="分配状态">
            <el-tag :type="getStatusType((machine as any).allocation_status)">
              {{ getStatusText((machine as any).allocation_status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="健康状态">
            <el-tag :type="machine.health_status === 'healthy' ? 'success' : 'danger'">
              {{ machine.health_status === 'healthy' ? '健康' : '异常' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="采集状态">
            <el-tag :type="machine.needs_collect ? 'warning' : 'success'">
              {{ machine.needs_collect ? '待采集' : '已采集' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 网络信息 -->
      <el-card class="info-card">
        <template #header>
          <span>网络信息</span>
        </template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="内网IP">{{ machine.ip_address || '-' }}</el-descriptions-item>
          <el-descriptions-item label="公网IP">{{ machine.public_ip || '-' }}</el-descriptions-item>
          <el-descriptions-item label="SSH端口">{{ machine.ssh_port || '-' }}</el-descriptions-item>
          <el-descriptions-item label="SSH用户">{{ machine.ssh_username || '-' }}</el-descriptions-item>
          <el-descriptions-item label="Agent端口">{{ machine.agent_port || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 连接信息 -->
      <el-card class="info-card">
        <template #header>
          <span>连接信息</span>
        </template>
        <el-tabs v-model="connectionTab">
          <el-tab-pane label="SSH" name="ssh">
            <div class="connection-info">
              <div class="info-item">
                <span class="info-label">连接主机：</span>
                <span class="info-value">{{ sshConnectHost || '-' }}</span>
                <el-button v-if="sshConnectHost" link :icon="CopyDocument" @click="copyToClipboard(sshConnectHost)" />
              </div>
              <div class="info-item">
                <span class="info-label">端口：</span>
                <span class="info-value">{{ machine.ssh_port || 22 }}</span>
                <el-button link :icon="CopyDocument" @click="copyToClipboard(String(machine.ssh_port || 22))" />
              </div>
              <div class="info-item">
                <span class="info-label">用户：</span>
                <span class="info-value">{{ machine.ssh_username || 'root' }}</span>
                <el-button link :icon="CopyDocument" @click="copyToClipboard(machine.ssh_username || 'root')" />
                <span class="info-label" style="margin-left: 24px">密码：</span>
                <span class="info-value">{{ showPassword ? (machine.ssh_password || '-') : (maskedPassword || '-') }}</span>
                <el-button v-if="machine.ssh_password" link :icon="showPassword ? Hide : View" @click="showPassword = !showPassword" />
                <el-button v-if="machine.ssh_password" link :icon="CopyDocument" @click="copyToClipboard(machine.ssh_password!)" />
              </div>
              <div class="info-item">
                <span class="info-label">连接命令：</span>
                <code class="info-value ssh-cmd">{{ sshCommand || '-' }}</code>
                <el-button v-if="sshCommand" link :icon="CopyDocument" @click="copyToClipboard(sshCommand)" />
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="Jupyter" name="jupyter">
            <div class="connection-info">
              <template v-if="machine.jupyter_url">
                <div class="info-item">
                  <span class="info-label">访问地址：</span>
                  <a class="info-value info-link" :href="machine.jupyter_url" target="_blank">{{ machine.jupyter_url }}</a>
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(machine.jupyter_url!)" />
                </div>
                <div class="info-item" v-if="machine.jupyter_token">
                  <span class="info-label">Token：</span>
                  <span class="info-value">{{ machine.jupyter_token }}</span>
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(machine.jupyter_token!)" />
                </div>
              </template>
              <div v-else class="info-empty">暂未配置 Jupyter 信息</div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="VNC" name="vnc">
            <div class="connection-info">
              <template v-if="machine.vnc_url">
                <div class="info-item">
                  <span class="info-label">访问地址：</span>
                  <a class="info-value info-link" :href="machine.vnc_url" target="_blank">{{ machine.vnc_url }}</a>
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(machine.vnc_url!)" />
                </div>
                <div class="info-item" v-if="machine.vnc_password">
                  <span class="info-label">密码：</span>
                  <span class="info-value">{{ showVncPassword ? machine.vnc_password : maskedVncPassword }}</span>
                  <el-button link :icon="showVncPassword ? Hide : View" @click="showVncPassword = !showVncPassword" />
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(machine.vnc_password!)" />
                </div>
              </template>
              <div v-else class="info-empty">暂未配置 VNC 信息</div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-card>

      <!-- 使用情况 -->
      <el-card class="info-card" v-if="usage">
        <template #header>
          <span>使用情况</span>
        </template>
        <div class="usage-grid">
          <div class="usage-item">
            <span class="usage-label">CPU</span>
            <el-progress
              type="dashboard"
              :percentage="Math.round(usage.cpu_usage)"
              :color="usageColor(usage.cpu_usage)"
              :width="90"
            />
          </div>
          <div class="usage-item">
            <span class="usage-label">内存</span>
            <el-progress
              type="dashboard"
              :percentage="Math.round(usage.memory_usage)"
              :color="usageColor(usage.memory_usage)"
              :width="90"
            />
          </div>
          <div class="usage-item">
            <span class="usage-label">磁盘</span>
            <el-progress
              type="dashboard"
              :percentage="Math.round(usage.disk_usage)"
              :color="usageColor(usage.disk_usage)"
              :width="90"
            />
          </div>
        </div>
        <div class="usage-meta" v-if="usage.collected_at">
          <span>采集时间：{{ formatDateTime(usage.collected_at) }}</span>
        </div>
      </el-card>

      <!-- 硬件配置 -->
      <el-card class="info-card">
        <template #header>
          <span>硬件配置</span>
        </template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="CPU">{{ machine.cpu_info || '-' }}</el-descriptions-item>
          <el-descriptions-item label="CPU核心数">{{ machine.total_cpu || '-' }}</el-descriptions-item>
          <el-descriptions-item label="内存">{{ formatSizeGB(machine.total_memory_gb) }}</el-descriptions-item>
          <el-descriptions-item label="磁盘">{{ formatSizeGB(machine.total_disk_gb) }}</el-descriptions-item>
          <el-descriptions-item label="操作系统">{{ machine.os_type || '-' }}</el-descriptions-item>
          <el-descriptions-item label="内核版本">{{ machine.os_version || '-' }}</el-descriptions-item>
          <el-descriptions-item label="GPU数量">{{ machine.gpus?.length || 0 }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- GPU 信息 -->
      <el-card class="info-card" v-if="machine.gpus && machine.gpus.length > 0">
        <template #header>
          <span>GPU 信息</span>
        </template>
        <el-table :data="machine.gpus" border>
          <el-table-column prop="name" label="型号" />
          <el-table-column label="显存">
            <template #default="{ row }">
              {{ formatGPUSize(row.memory_total_mb) }}
            </template>
          </el-table-column>
          <el-table-column prop="uuid" label="UUID" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" />
        </el-table>
      </el-card>

      <!-- 分配信息 -->
      <el-card class="info-card">
        <template #header>
          <span>分配信息</span>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="当前状态">
            <el-tag :type="activeAllocation ? 'warning' : 'info'">
              {{ activeAllocation ? '已分配' : '未分配' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="客户">
            {{
              activeAllocation?.customer?.company ||
              activeAllocation?.customer?.display_name ||
              activeAllocation?.customer?.username ||
              '-'
            }}
          </el-descriptions-item>
          <el-descriptions-item label="开始时间">
            {{ formatDateTime(activeAllocation?.start_time) }}
          </el-descriptions-item>
          <el-descriptions-item label="到期时间">
            {{ formatDateTime(activeAllocation?.end_time) }}
          </el-descriptions-item>
        </el-descriptions>

        <el-table
          v-if="machine.allocations && machine.allocations.length > 0"
          :data="machine.allocations"
          border
          style="margin-top: 16px"
        >
          <el-table-column prop="id" label="分配ID" width="180" show-overflow-tooltip />
          <el-table-column label="客户">
            <template #default="{ row }">
              {{ row.customer?.company || row.customer?.display_name || row.customer?.username || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" />
          <el-table-column label="开始时间">
            <template #default="{ row }">
              {{ formatDateTime(row.start_time) }}
            </template>
          </el-table-column>
          <el-table-column label="结束时间">
            <template #default="{ row }">
              {{ formatDateTime(row.end_time) }}
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <!-- 外映射配置 -->
      <el-card class="info-card" v-if="machine.external_ip || machine.nginx_domain || machine.external_ssh_port || machine.external_jupyter_port || machine.external_vnc_port">
        <template #header>
          <span>外映射配置</span>
        </template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="外部IP/域名">{{ machine.external_ip || '-' }}</el-descriptions-item>
          <el-descriptions-item label="Nginx域名">{{ machine.nginx_domain || '-' }}</el-descriptions-item>
          <el-descriptions-item label="Nginx配置路径">{{ machine.nginx_config_path || '-' }}</el-descriptions-item>
          <el-descriptions-item label="SSH映射端口">{{ machine.external_ssh_port || '-' }}</el-descriptions-item>
          <el-descriptions-item label="Jupyter映射端口">{{ machine.external_jupyter_port || '-' }}</el-descriptions-item>
          <el-descriptions-item label="VNC映射端口">{{ machine.external_vnc_port || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 时间信息 -->
      <el-card class="info-card">
        <template #header>
          <span>时间信息</span>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="创建时间">{{ formatDateTime(machine.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ formatDateTime(machine.updated_at) }}</el-descriptions-item>
          <el-descriptions-item label="最后心跳">{{ formatDateTime(machine.last_heartbeat) }}</el-descriptions-item>
        </el-descriptions>
      </el-card>
    </template>
  </div>
</template>

<style scoped>
.machine-detail {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.info-card {
  margin-bottom: 20px;
}

.ssh-cmd {
  font-size: 13px;
  color: #606266;
  background: #f5f7fa;
  padding: 4px 8px;
  border-radius: 4px;
}

.connection-info {
  display: flex;
  flex-direction: column;
  gap: 0;
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
  min-width: 70px;
}

.info-value {
  color: #303133;
  font-size: 14px;
  word-break: break-all;
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

.usage-grid {
  display: flex;
  justify-content: space-around;
  padding: 16px 0;
}

.usage-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.usage-label {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.usage-meta {
  text-align: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
  font-size: 12px;
  color: #909399;
}
</style>
