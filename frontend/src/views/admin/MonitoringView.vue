<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Monitor, CircleCheck, Warning, Cpu } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { getRealtimeMonitoring, getMachineList } from '@/api/admin'
import type { Machine } from '@/types/machine'

// 扩展 Machine 类型以包含（模拟的）监控数据
interface MonitoredMachine extends Machine {
  cpuUsage?: number
  memoryUsage?: number
  gpuUsage?: number
  diskUsage?: number
  temperature?: number
  uptime?: string
}

const devices = ref<MonitoredMachine[]>([])
const loading = ref(false)
const timer = ref<number | null>(null)

// 总体统计
const totalStats = ref({
  totalMachines: 0,
  onlineMachines: 0,
  idleMachines: 0,
  allocatedMachines: 0,
  offlineMachines: 0,
  avgGpuUtil: 0
})

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    // 1. 获取聚合快照
    const snapshotRes = await getRealtimeMonitoring()
    if (snapshotRes.data) {
      totalStats.value = {
        totalMachines: snapshotRes.data.total_machines || 0,
        onlineMachines: snapshotRes.data.online_machines || 0,
        idleMachines: snapshotRes.data.idle_machines || 0,
        allocatedMachines: snapshotRes.data.allocated_machines || 0,
        offlineMachines: snapshotRes.data.offline_machines || 0,
        avgGpuUtil: snapshotRes.data.avg_gpu_util || 0
      }
    }

    // 2. 获取机器列表 (模拟详细监控数据)
    // 注意：实际生产中应该有一个专门的 /admin/monitoring/devices 接口返回详细指标
    const listRes = await getMachineList({ page: 1, pageSize: 100 })
    devices.value = listRes.data.list.map((m: Machine) => ({
      ...m,
      // Mock metrics until backend supports them
      cpuUsage: Math.floor(Math.random() * 60) + 10,
      memoryUsage: Math.floor(Math.random() * 70) + 20,
      gpuUsage: Math.floor(Math.random() * 80) + 5,
      diskUsage: 45,
      temperature: 60 + Math.floor(Math.random() * 20),
      uptime: 'Running'
    }))

  } catch (error) {
    console.error(error)
    ElMessage.error('加载监控数据失败')
  } finally {
    loading.value = false
  }
}

// 获取使用率颜色
const getUsageColor = (usage: number) => {
  if (usage >= 90) return 'danger'
  if (usage >= 70) return 'warning'
  return 'success'
}

// 获取状态类型
const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    idle: 'success',
    allocated: 'primary',
    maintenance: 'warning',
    offline: 'danger'
  }
  return statusMap[status] || 'info'
}

onMounted(() => {
  loadData()
  // 30秒刷新一次
  timer.value = setInterval(loadData, 30000) as unknown as number
})

onUnmounted(() => {
  if (timer.value) clearInterval(timer.value)
})
</script>

<template>
  <div class="monitoring-center">
    <PageHeader title="监控中心" />

    <!-- 总体统计卡片 -->
    <div class="stats-container">
      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #409EFF">
            <el-icon :size="32"><Monitor /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">总机器数</div>
            <div class="stat-value">{{ totalStats.totalMachines }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #67C23A">
            <el-icon :size="32"><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">在线(空闲+分配)</div>
            <div class="stat-value">{{ totalStats.onlineMachines }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #E6A23C">
            <el-icon :size="32"><Warning /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">离线机器</div>
            <div class="stat-value">{{ totalStats.offlineMachines }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #909399">
            <el-icon :size="32"><Cpu /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">平均GPU利用率</div>
            <div class="stat-value">{{ totalStats.avgGpuUtil }}%</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 设备监控列表 -->
    <el-card class="device-list-card">
      <template #header>
        <div class="card-header">
          <span>设备实时状态</span>
          <el-button type="primary" size="small" @click="loadData">
            刷新
          </el-button>
        </div>
      </template>

      <el-table :data="devices" v-loading="loading" stripe>
        <el-table-column prop="ip_address" label="IP地址" width="150" />
        <el-table-column prop="region" label="区域" width="100" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="CPU" width="180">
          <template #default="{ row }">
            <div class="metric-cell">
              <el-progress
                :percentage="row.cpuUsage"
                :color="getUsageColor(row.cpuUsage)"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column label="内存" width="180">
          <template #default="{ row }">
            <div class="metric-cell">
              <el-progress
                :percentage="row.memoryUsage"
                :color="getUsageColor(row.memoryUsage)"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column label="GPU" width="180">
          <template #default="{ row }">
            <div class="metric-cell">
              <el-progress
                :percentage="row.gpuUsage"
                :color="getUsageColor(row.gpuUsage)"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="gpu_model" label="GPU型号" />
      </el-table>
    </el-card>
  </div>
</template>

<style scoped>
.monitoring-center {
  padding: 24px;
}

.stats-container {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.device-list-card {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.metric-cell {
  padding: 4px 0;
}
</style>
