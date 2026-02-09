<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Monitor, CircleCheck, Warning, Cpu } from '@element-plus/icons-vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { PieChart, BarChart, GaugeChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { getRealtimeMonitoring, getMachineList } from '@/api/admin'
import type { Machine } from '@/types/machine'

use([PieChart, BarChart, GaugeChart, TitleComponent, TooltipComponent, LegendComponent, GridComponent, CanvasRenderer])

interface MonitoredMachine extends Machine {
  cpuUsage?: number
  memoryUsage?: number
  gpuUsage?: number
}

const devices = ref<MonitoredMachine[]>([])
const loading = ref(false)
const timer = ref<number | null>(null)

const totalStats = ref({
  totalMachines: 0,
  onlineMachines: 0,
  idleMachines: 0,
  allocatedMachines: 0,
  offlineMachines: 0,
  maintenanceMachines: 0,
  avgGpuUtil: 0
})

const loadData = async () => {
  loading.value = true
  try {
    const snapshotRes = await getRealtimeMonitoring()
    if (snapshotRes.data) {
      totalStats.value = {
        totalMachines: snapshotRes.data.total_machines || 0,
        onlineMachines: snapshotRes.data.online_machines || 0,
        idleMachines: snapshotRes.data.idle_machines || 0,
        allocatedMachines: snapshotRes.data.allocated_machines || 0,
        offlineMachines: snapshotRes.data.offline_machines || 0,
        maintenanceMachines: snapshotRes.data.maintenance_machines || 0,
        avgGpuUtil: snapshotRes.data.avg_gpu_util || 0
      }
    }
    const listRes = await getMachineList({ page: 1, pageSize: 100 })
    devices.value = listRes.data.list.map((m: Machine) => ({
      ...m,
      cpuUsage: (m as any).cpu_usage ?? 0,
      memoryUsage: (m as any).memory_usage ?? 0,
      gpuUsage: (m as any).gpu_usage ?? 0,
    }))
  } catch (error) {
    console.error(error)
    ElMessage.error('加载监控数据失败')
  } finally {
    loading.value = false
  }
}

// 机器状态分布饼图
const statusPieOption = computed(() => ({
  title: { text: '机器状态分布', left: 'center', textStyle: { fontSize: 14, color: '#303133' } },
  tooltip: { trigger: 'item', formatter: '{b}: {c} 台 ({d}%)' },
  legend: { bottom: 0, itemWidth: 10, itemHeight: 10, textStyle: { fontSize: 12 } },
  series: [{
    type: 'pie',
    radius: ['40%', '65%'],
    center: ['50%', '45%'],
    avoidLabelOverlap: true,
    itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
    label: { show: true, formatter: '{b}\n{c}台' },
    data: [
      { value: totalStats.value.idleMachines, name: '空闲', itemStyle: { color: '#67C23A' } },
      { value: totalStats.value.allocatedMachines, name: '已分配', itemStyle: { color: '#409EFF' } },
      { value: totalStats.value.offlineMachines, name: '离线', itemStyle: { color: '#F56C6C' } },
      { value: totalStats.value.maintenanceMachines, name: '维护中', itemStyle: { color: '#E6A23C' } },
    ].filter(d => d.value > 0)
  }]
}))

// GPU 型号分布柱状图
const gpuModelData = computed(() => {
  const map: Record<string, number> = {}
  devices.value.forEach(d => {
    if (d.gpus && d.gpus.length > 0) {
      const name = d.gpus[0].name || '未知'
      map[name] = (map[name] || 0) + d.gpus.length
    }
  })
  const entries = Object.entries(map).sort((a, b) => b[1] - a[1])
  return { names: entries.map(e => e[0]), counts: entries.map(e => e[1]) }
})

const gpuBarOption = computed(() => ({
  title: { text: 'GPU 型号分布', left: 'center', textStyle: { fontSize: 14, color: '#303133' } },
  tooltip: { trigger: 'axis', formatter: '{b}: {c} 张' },
  grid: { left: 16, right: 16, bottom: 8, top: 40, containLabel: true },
  xAxis: { type: 'category', data: gpuModelData.value.names, axisLabel: { fontSize: 11, rotate: gpuModelData.value.names.length > 4 ? 20 : 0 } },
  yAxis: { type: 'value', minInterval: 1, axisLabel: { fontSize: 11 } },
  series: [{
    type: 'bar',
    data: gpuModelData.value.counts,
    barMaxWidth: 40,
    itemStyle: { color: '#409EFF', borderRadius: [4, 4, 0, 0] },
    label: { show: true, position: 'top', fontSize: 12 }
  }]
}))

// 区域分布饼图
const regionData = computed(() => {
  const map: Record<string, number> = {}
  devices.value.forEach(d => {
    const r = d.region || '未知'
    map[r] = (map[r] || 0) + 1
  })
  return Object.entries(map).map(([name, value]) => ({ name, value }))
})

const regionPieOption = computed(() => ({
  title: { text: '区域分布', left: 'center', textStyle: { fontSize: 14, color: '#303133' } },
  tooltip: { trigger: 'item', formatter: '{b}: {c} 台 ({d}%)' },
  legend: { bottom: 0, itemWidth: 10, itemHeight: 10, textStyle: { fontSize: 12 } },
  color: ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399', '#B37FEB', '#36CFC9'],
  series: [{
    type: 'pie',
    radius: ['40%', '65%'],
    center: ['50%', '45%'],
    itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
    label: { show: true, formatter: '{b}\n{c}台' },
    data: regionData.value
  }]
}))

// GPU 利用率仪表盘
const gpuGaugeOption = computed(() => ({
  series: [{
    type: 'gauge',
    startAngle: 210,
    endAngle: -30,
    radius: '90%',
    progress: { show: true, width: 14, roundCap: true },
    axisLine: { lineStyle: { width: 14, color: [[0.3, '#67C23A'], [0.7, '#E6A23C'], [1, '#F56C6C']] } },
    axisTick: { show: false },
    splitLine: { show: false },
    axisLabel: { show: false },
    pointer: { show: false },
    title: { offsetCenter: [0, '70%'], fontSize: 13, color: '#909399' },
    detail: { valueAnimation: true, fontSize: 28, fontWeight: 600, offsetCenter: [0, '35%'], formatter: '{value}%', color: '#303133' },
    data: [{ value: totalStats.value.avgGpuUtil, name: '平均GPU利用率' }]
  }]
}))

const getDeviceStatusType = (status: string) => status === 'online' ? 'success' : 'danger'
const getAllocationStatusType = (status: string) => {
  const m: Record<string, string> = { idle: 'success', allocated: 'primary', maintenance: 'warning' }
  return m[status] || 'info'
}
const getUsageColor = (usage: number) => {
  if (usage >= 90) return '#F56C6C'
  if (usage >= 70) return '#E6A23C'
  return '#67C23A'
}

onMounted(() => {
  loadData()
  timer.value = setInterval(loadData, 30000) as unknown as number
})
onUnmounted(() => {
  if (timer.value) clearInterval(timer.value)
})
</script>

<template>
  <div class="monitoring-center">
    <div class="page-header">
      <h2 class="page-title">监控中心</h2>
      <el-button type="primary" size="small" @click="loadData" :loading="loading">刷新数据</el-button>
    </div>

    <!-- 顶部统计卡片 -->
    <div class="stats-row">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #409EFF, #66b1ff)">
            <el-icon :size="28"><Monitor /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">总机器数</div>
            <div class="stat-value">{{ totalStats.totalMachines }}</div>
          </div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #67C23A, #85ce61)">
            <el-icon :size="28"><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">在线机器</div>
            <div class="stat-value">{{ totalStats.onlineMachines }}</div>
          </div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #F56C6C, #f89898)">
            <el-icon :size="28"><Warning /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">离线机器</div>
            <div class="stat-value">{{ totalStats.offlineMachines }}</div>
          </div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon" style="background: linear-gradient(135deg, #E6A23C, #ebb563)">
            <el-icon :size="28"><Cpu /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">平均GPU利用率</div>
            <div class="stat-value">{{ totalStats.avgGpuUtil }}<span class="stat-unit">%</span></div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 图表区域 -->
    <div class="charts-row">
      <el-card class="chart-card" shadow="hover">
        <v-chart :option="statusPieOption" autoresize style="height: 280px" />
      </el-card>
      <el-card class="chart-card" shadow="hover">
        <v-chart :option="gpuBarOption" autoresize style="height: 280px" />
      </el-card>
      <el-card class="chart-card chart-card-sm" shadow="hover">
        <v-chart :option="gpuGaugeOption" autoresize style="height: 280px" />
      </el-card>
    </div>

    <div class="charts-row">
      <el-card class="chart-card" shadow="hover">
        <v-chart :option="regionPieOption" autoresize style="height: 280px" />
      </el-card>
    </div>

    <!-- 设备实时状态表格 -->
    <el-card class="device-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">设备实时状态</span>
          <span class="card-subtitle">共 {{ devices.length }} 台设备</span>
        </div>
      </template>
      <el-table :data="devices" v-loading="loading" stripe style="width: 100%">
        <el-table-column label="机器" min-width="140">
          <template #default="{ row }">
            <span class="device-name">{{ row.name || row.hostname || row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="ip_address" label="IP" width="140" />
        <el-table-column prop="region" label="区域" width="90" />
        <el-table-column label="状态" width="160">
          <template #default="{ row }">
            <el-tag :type="getDeviceStatusType(row.device_status)" size="small">
              {{ row.device_status === 'online' ? '在线' : '离线' }}
            </el-tag>
            <el-tag :type="getAllocationStatusType(row.allocation_status)" size="small" style="margin-left: 4px">
              {{ { idle: '空闲', allocated: '已分配', maintenance: '维护' }[row.allocation_status as string] || row.allocation_status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="CPU" width="160">
          <template #default="{ row }">
            <el-progress :percentage="row.cpuUsage || 0" :color="getUsageColor(row.cpuUsage || 0)" :stroke-width="10" />
          </template>
        </el-table-column>
        <el-table-column label="内存" width="160">
          <template #default="{ row }">
            <el-progress :percentage="row.memoryUsage || 0" :color="getUsageColor(row.memoryUsage || 0)" :stroke-width="10" />
          </template>
        </el-table-column>
        <el-table-column label="GPU" width="160">
          <template #default="{ row }">
            <el-progress :percentage="row.gpuUsage || 0" :color="getUsageColor(row.gpuUsage || 0)" :stroke-width="10" />
          </template>
        </el-table-column>
        <el-table-column label="GPU型号" min-width="140">
          <template #default="{ row }">
            {{ row.gpus?.[0]?.name || '-' }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<style scoped>
.monitoring-center {
  padding: 24px;
  background: #f5f7fa;
  min-height: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  font-size: 22px;
  font-weight: 700;
  color: #1d2129;
  margin: 0;
}

/* 统计卡片 */
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card :deep(.el-card__body) {
  padding: 16px 20px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 14px;
}

.stat-icon {
  width: 52px;
  height: 52px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
}

.stat-unit {
  font-size: 14px;
  font-weight: 400;
  color: #909399;
  margin-left: 2px;
}

/* 图表区域 */
.charts-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.chart-card :deep(.el-card__body) {
  padding: 12px;
}

/* 设备表格 */
.device-card {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.card-subtitle {
  font-size: 13px;
  color: #909399;
}

.device-name {
  font-weight: 500;
  color: #303133;
}
</style>
