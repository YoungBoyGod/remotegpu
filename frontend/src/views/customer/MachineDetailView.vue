<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/common/PageHeader.vue'
import { getMyMachines, getMachineConnection, resetSSH } from '@/api/customer'
import type { Machine } from '@/types/machine'

const route = useRoute()
const router = useRouter()

const machineId = computed(() => String(route.params.id || ''))
const loading = ref(false)
const connectionLoading = ref(false)
const machine = ref<Machine | null>(null)
const connectionInfo = ref<any>(null)

const loadMachineDetail = async () => {
  loading.value = true
  try {
    const res = await getMyMachines({ page: 1, pageSize: 200 })
    const found = res.data.list.find(item => String(item.id) === machineId.value)
    if (!found) {
      ElMessage.error('未找到机器信息')
    }
    machine.value = found || null
  } catch (error) {
    console.error('加载机器详情失败:', error)
  } finally {
    loading.value = false
  }
}

const loadConnectionInfo = async () => {
  connectionLoading.value = true
  try {
    const res = await getMachineConnection(machineId.value)
    connectionInfo.value = res.data
  } catch (error) {
    console.error('加载连接信息失败:', error)
  } finally {
    connectionLoading.value = false
  }
}

const handleBack = () => {
  router.push('/customer/machines/list')
}

const handleResetSSH = async () => {
  try {
    await resetSSH(machineId.value)
    ElMessage.success('已触发 SSH 重置')
  } catch (error) {
    console.error('SSH 重置失败:', error)
  }
}

const statusTagType = (status?: string) => {
  const map: Record<string, any> = {
    online: 'success',
    offline: 'danger',
    maintenance: 'warning',
  }
  return map[status || ''] || 'info'
}

const copyText = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

const formatDate = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

onMounted(async () => {
  await loadMachineDetail()
  await loadConnectionInfo()
})
</script>

<template>
  <div class="machine-detail">
    <PageHeader title="机器详情">
      <template #actions>
        <el-button @click="handleBack">返回列表</el-button>
        <el-button type="warning" @click="handleResetSSH">重置SSH</el-button>
      </template>
    </PageHeader>

    <!-- 机器名称和状态 -->
    <div class="detail-banner" v-loading="loading">
      <div class="banner-left">
        <h3 class="machine-name">{{ machine?.name || machine?.hostname || '-' }}</h3>
        <el-tag :type="statusTagType(machine?.device_status)" size="large">
          {{ machine?.device_status === 'online' ? '在线' : '离线' }}
        </el-tag>
      </div>
      <div class="banner-meta">
        <span class="meta-item">ID: {{ machine?.id || '-' }}</span>
        <span class="meta-item">区域: {{ machine?.region || '-' }}</span>
        <span class="meta-item">分配时间: {{ formatDate(machine?.start_time) }} ~ {{ formatDate(machine?.end_time) }}</span>
      </div>
    </div>

    <!-- 资源概览 -->
    <div class="resource-cards">
      <div class="resource-card">
        <div class="resource-label">CPU</div>
        <div class="resource-value">{{ machine?.total_cpu || '-' }} <span class="resource-unit">核</span></div>
      </div>
      <div class="resource-card">
        <div class="resource-label">内存</div>
        <div class="resource-value">{{ machine?.total_memory_gb || '-' }} <span class="resource-unit">GB</span></div>
      </div>
      <div class="resource-card">
        <div class="resource-label">GPU</div>
        <div class="resource-value">{{ machine?.gpus?.length || 0 }}x <span class="resource-unit">{{ machine?.gpus?.[0]?.name || '-' }}</span></div>
      </div>
      <div class="resource-card">
        <div class="resource-label">GPU 显存</div>
        <div class="resource-value">{{ machine?.gpus?.[0]?.memory_total_mb ? Math.round(machine.gpus[0].memory_total_mb / 1024) : '-' }} <span class="resource-unit">GB</span></div>
      </div>
    </div>

    <!-- 连接信息卡片 -->
    <div class="connection-section" v-loading="connectionLoading">
      <h4 class="section-title">连接方式</h4>
      <div class="connection-cards">
        <!-- SSH -->
        <el-card class="conn-card" shadow="hover">
          <div class="conn-card-header">
            <span class="conn-tag conn-tag-ssh">SSH</span>
            <el-button v-if="connectionInfo?.ssh?.command" link size="small" @click="copyText(connectionInfo.ssh.command)">复制命令</el-button>
          </div>
          <template v-if="connectionInfo?.ssh">
            <div class="conn-field"><span class="conn-label">主机</span><span class="conn-value">{{ connectionInfo.ssh.host }}</span></div>
            <div class="conn-field"><span class="conn-label">端口</span><span class="conn-value">{{ connectionInfo.ssh.port }}</span></div>
            <div class="conn-field"><span class="conn-label">用户名</span><span class="conn-value">{{ connectionInfo.ssh.username }}</span></div>
            <div v-if="connectionInfo.ssh.command" class="conn-command">{{ connectionInfo.ssh.command }}</div>
          </template>
          <div v-else class="conn-empty">暂无 SSH 连接信息</div>
        </el-card>

        <!-- Jupyter -->
        <el-card class="conn-card" shadow="hover">
          <div class="conn-card-header">
            <span class="conn-tag conn-tag-jupyter">Jupyter</span>
          </div>
          <template v-if="connectionInfo?.jupyter?.url">
            <div class="conn-field"><span class="conn-label">地址</span><a class="conn-link" :href="connectionInfo.jupyter.url" target="_blank">{{ connectionInfo.jupyter.url }}</a></div>
            <div v-if="connectionInfo.jupyter.token" class="conn-field"><span class="conn-label">Token</span><span class="conn-value mono">{{ connectionInfo.jupyter.token }}</span></div>
          </template>
          <div v-else class="conn-empty">暂无 Jupyter 连接信息</div>
        </el-card>

        <!-- VNC -->
        <el-card class="conn-card" shadow="hover">
          <div class="conn-card-header">
            <span class="conn-tag conn-tag-vnc">VNC</span>
          </div>
          <template v-if="connectionInfo?.vnc?.url">
            <div class="conn-field"><span class="conn-label">地址</span><a class="conn-link" :href="connectionInfo.vnc.url" target="_blank">{{ connectionInfo.vnc.url }}</a></div>
            <div v-if="connectionInfo.vnc.password" class="conn-field"><span class="conn-label">密码</span><span class="conn-value mono">{{ connectionInfo.vnc.password }}</span></div>
          </template>
          <div v-else class="conn-empty">暂无 VNC 连接信息</div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<style scoped>
.machine-detail {
  padding: 24px;
  background: #f5f7fa;
  min-height: 100%;
}

.detail-banner {
  background: #fff;
  border-radius: 8px;
  padding: 20px 24px;
  margin-bottom: 16px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.banner-left {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.machine-name {
  font-size: 20px;
  font-weight: 700;
  color: #1d2129;
  margin: 0;
}

.banner-meta {
  display: flex;
  gap: 24px;
  font-size: 13px;
  color: #86909c;
}

.resource-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.resource-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px 20px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.resource-label {
  font-size: 13px;
  color: #86909c;
  margin-bottom: 8px;
}

.resource-value {
  font-size: 22px;
  font-weight: 700;
  color: #1d2129;
}

.resource-unit {
  font-size: 13px;
  font-weight: 400;
  color: #86909c;
}

.connection-section {
  margin-top: 4px;
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: #1d2129;
  margin: 0 0 12px 0;
}

.connection-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.conn-card {
  border-radius: 8px;
  border: none;
}

.conn-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.conn-tag {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

.conn-tag-ssh { color: #409eff; background: #ecf5ff; }
.conn-tag-jupyter { color: #e6a23c; background: #fdf6ec; }
.conn-tag-vnc { color: #67c23a; background: #f0f9eb; }

.conn-field {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  font-size: 13px;
  border-bottom: 1px solid #f2f3f5;
}

.conn-field:last-of-type {
  border-bottom: none;
}

.conn-label {
  color: #86909c;
  flex-shrink: 0;
}

.conn-value {
  color: #1d2129;
  text-align: right;
  word-break: break-all;
}

.conn-value.mono {
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.conn-link {
  color: #409eff;
  text-decoration: none;
  font-size: 13px;
  word-break: break-all;
}

.conn-link:hover {
  text-decoration: underline;
}

.conn-command {
  margin-top: 8px;
  padding: 8px 10px;
  background: #f7f8fa;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #4e5969;
  word-break: break-all;
}

.conn-empty {
  color: #c0c4cc;
  font-size: 13px;
  text-align: center;
  padding: 16px 0;
}

@media (max-width: 1000px) {
  .resource-cards { grid-template-columns: repeat(2, 1fr); }
  .connection-cards { grid-template-columns: 1fr; }
}
</style>
