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

    <el-row :gutter="16">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>基础信息</span>
          </template>
          <el-descriptions :column="1" border v-loading="loading">
            <el-descriptions-item label="机器ID">{{ machine?.id || '-' }}</el-descriptions-item>
            <el-descriptions-item label="机器名称">{{ machine?.name || machine?.hostname || '-' }}</el-descriptions-item>
            <el-descriptions-item label="主机名">{{ machine?.hostname || '-' }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="statusTagType(machine?.device_status)">
                {{ machine?.device_status === 'online' ? '在线' : '离线' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="内网IP">{{ machine?.ip_address || '-' }}</el-descriptions-item>
            <el-descriptions-item label="公网IP">{{ machine?.public_ip || '-' }}</el-descriptions-item>
            <el-descriptions-item label="开始时间">{{ formatDate(machine?.start_time) }}</el-descriptions-item>
            <el-descriptions-item label="到期时间">{{ formatDate(machine?.end_time) }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>资源信息</span>
          </template>
          <el-descriptions :column="1" border v-loading="loading">
            <el-descriptions-item label="CPU(核)">{{ machine?.total_cpu || '-' }}</el-descriptions-item>
            <el-descriptions-item label="内存(GB)">{{ machine?.total_memory_gb || '-' }}</el-descriptions-item>
            <el-descriptions-item label="GPU型号">{{ machine?.gpus?.[0]?.name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="GPU数量">{{ machine?.gpus?.length || '-' }}</el-descriptions-item>
            <el-descriptions-item label="GPU显存(GB)">{{ machine?.gpus?.[0]?.memory_total_mb ? Math.round(machine.gpus[0].memory_total_mb / 1024) : '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="connection-row">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>连接信息</span>
          </template>
          <el-descriptions :column="2" border v-loading="connectionLoading">
            <el-descriptions-item label="SSH 主机">{{ connectionInfo?.ssh?.host || '-' }}</el-descriptions-item>
            <el-descriptions-item label="SSH 端口">{{ connectionInfo?.ssh?.port || '-' }}</el-descriptions-item>
            <el-descriptions-item label="SSH 用户名">{{ connectionInfo?.ssh?.username || '-' }}</el-descriptions-item>
            <el-descriptions-item label="VNC 地址">{{ connectionInfo?.vnc?.url || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.machine-detail {
  padding: 24px;
}

.connection-row {
  margin-top: 16px;
}
</style>
