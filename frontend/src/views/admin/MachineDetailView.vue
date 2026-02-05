<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getMachineDetail } from '@/api/admin'
import type { Machine } from '@/types/machine'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const machine = ref<Machine | null>(null)

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

onMounted(() => {
  loadMachine()
})
</script>

<template>
  <div class="machine-detail" v-loading="loading">
    <div class="page-header">
      <div class="header-left">
        <el-button @click="handleBack" :icon="'ArrowLeft'">返回列表</el-button>
        <h2 class="page-title">{{ machine?.name || machine?.hostname || machine?.id || '机器详情' }}</h2>
        <el-tag v-if="machine" :type="getStatusType(machine.status)">
          {{ getStatusText(machine.status) }}
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
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(machine.status)">
              {{ getStatusText(machine.status) }}
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
</style>
