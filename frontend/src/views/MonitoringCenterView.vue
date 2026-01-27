<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Monitor, CircleCheck, Warning, Cpu } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'

interface DeviceMetrics {
  id: number
  name: string
  status: 'online' | 'offline' | 'warning'
  cpuUsage: number
  memoryUsage: number
  gpuUsage: number
  diskUsage: number
  temperature: number
  uptime: string
}

const devices = ref<DeviceMetrics[]>([])
const loading = ref(false)

// 总体统计
const totalStats = ref({
  totalDevices: 0,
  onlineDevices: 0,
  warningDevices: 0,
  avgCpuUsage: 0,
  avgMemoryUsage: 0,
  avgGpuUsage: 0
})

// 加载设备监控数据
const loadDeviceMetrics = async () => {
  loading.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 500))
    devices.value = [
      {
        id: 1,
        name: '主机-001',
        status: 'online',
        cpuUsage: 45,
        memoryUsage: 68,
        gpuUsage: 82,
        diskUsage: 55,
        temperature: 65,
        uptime: '15天 8小时'
      },
      {
        id: 2,
        name: '主机-002',
        status: 'warning',
        cpuUsage: 88,
        memoryUsage: 92,
        gpuUsage: 95,
        diskUsage: 78,
        temperature: 82,
        uptime: '7天 12小时'
      },
      {
        id: 3,
        name: '主机-003',
        status: 'online',
        cpuUsage: 32,
        memoryUsage: 45,
        gpuUsage: 60,
        diskUsage: 42,
        temperature: 58,
        uptime: '22天 5小时'
      }
    ]

    // 计算总体统计
    totalStats.value = {
      totalDevices: devices.value.length,
      onlineDevices: devices.value.filter(d => d.status === 'online').length,
      warningDevices: devices.value.filter(d => d.status === 'warning').length,
      avgCpuUsage: Math.round(devices.value.reduce((sum, d) => sum + d.cpuUsage, 0) / devices.value.length),
      avgMemoryUsage: Math.round(devices.value.reduce((sum, d) => sum + d.memoryUsage, 0) / devices.value.length),
      avgGpuUsage: Math.round(devices.value.reduce((sum, d) => sum + d.gpuUsage, 0) / devices.value.length)
    }
  } catch (error) {
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
  const statusMap = {
    online: 'success',
    offline: 'info',
    warning: 'warning'
  }
  return statusMap[status as keyof typeof statusMap] || 'info'
}

// 获取状态文本
const getStatusText = (status: string) => {
  const statusMap = {
    online: '在线',
    offline: '离线',
    warning: '告警'
  }
  return statusMap[status as keyof typeof statusMap] || '未知'
}

onMounted(() => {
  loadDeviceMetrics()
  // 模拟实时更新
  setInterval(loadDeviceMetrics, 30000)
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
            <div class="stat-label">设备总数</div>
            <div class="stat-value">{{ totalStats.totalDevices }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #67C23A">
            <el-icon :size="32"><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">在线设备</div>
            <div class="stat-value">{{ totalStats.onlineDevices }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #E6A23C">
            <el-icon :size="32"><Warning /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">告警设备</div>
            <div class="stat-value">{{ totalStats.warningDevices }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #909399">
            <el-icon :size="32"><Cpu /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">平均CPU使用率</div>
            <div class="stat-value">{{ totalStats.avgCpuUsage }}%</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 设备监控列表 -->
    <el-card class="device-list-card">
      <template #header>
        <div class="card-header">
          <span>设备监控列表</span>
          <el-button type="primary" size="small" @click="loadDeviceMetrics">
            刷新数据
          </el-button>
        </div>
      </template>

      <el-table :data="devices" :loading="loading" stripe>
        <el-table-column prop="name" label="设备名称" width="150" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="CPU使用率" width="200">
          <template #default="{ row }">
            <div class="metric-cell">
              <el-progress
                :percentage="row.cpuUsage"
                :color="getUsageColor(row.cpuUsage)"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column label="内存使用率" width="200">
          <template #default="{ row }">
            <div class="metric-cell">
              <el-progress
                :percentage="row.memoryUsage"
                :color="getUsageColor(row.memoryUsage)"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column label="GPU使用率" width="200">
          <template #default="{ row }">
            <div class="metric-cell">
              <el-progress
                :percentage="row.gpuUsage"
                :color="getUsageColor(row.gpuUsage)"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column label="磁盘使用率" width="200">
          <template #default="{ row }">
            <div class="metric-cell">
              <el-progress
                :percentage="row.diskUsage"
                :color="getUsageColor(row.diskUsage)"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="temperature" label="温度(℃)" width="100" />
        <el-table-column prop="uptime" label="运行时间" width="150" />
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
  cursor: pointer;
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

@media (max-width: 1200px) {
  .stats-container {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-container {
    grid-template-columns: 1fr;
  }
}
</style>
