<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getDashboardOverview, getGPUTrend, getRecentAllocations } from '@/api/admin'
import StatCard from '@/components/common/StatCard.vue'

// 后端 /admin/dashboard/stats 返回结构
interface DashboardStats {
  total_machines: number
  allocated_machines: number
  idle_machines: number
  offline_machines: number
  active_customers: number
}

// 后端 /admin/dashboard/gpu-trend 返回结构
interface GPUTrendPoint {
  time: string
  usage: number
}

// 后端 /admin/allocations/recent 返回结构
interface RecentAllocation {
  id: string
  customer_id: number
  host_id: string
  start_time: string
  end_time: string
  status: string
  created_at: string
  customer?: { username?: string; display_name?: string; company?: string }
  host?: { name?: string; ip_address?: string }
}

const statsLoading = ref(true)
const trendLoading = ref(true)
const allocationsLoading = ref(true)

const stats = ref<DashboardStats>({
  total_machines: 0,
  allocated_machines: 0,
  idle_machines: 0,
  offline_machines: 0,
  active_customers: 0,
})

const gpuTrend = ref<GPUTrendPoint[]>([])
const recentAllocationList = ref<RecentAllocation[]>([])

// 计算在线机器数 = 总数 - 离线数
const onlineMachines = () => stats.value.total_machines - stats.value.offline_machines

const loadStats = async () => {
  try {
    statsLoading.value = true
    const response = await getDashboardOverview()
    stats.value = response.data
  } catch (error) {
    ElMessage.error('加载统计数据失败')
    console.error('加载统计数据失败:', error)
  } finally {
    statsLoading.value = false
  }
}

const loadGPUTrend = async () => {
  try {
    trendLoading.value = true
    const response = await getGPUTrend()
    gpuTrend.value = response.data || []
  } catch (error) {
    ElMessage.error('加载GPU趋势数据失败')
    console.error('加载GPU趋势数据失败:', error)
  } finally {
    trendLoading.value = false
  }
}

const loadRecentAllocations = async () => {
  try {
    allocationsLoading.value = true
    const response = await getRecentAllocations()
    recentAllocationList.value = response.data || []
  } catch (error) {
    ElMessage.error('加载最近分配记录失败')
    console.error('加载最近分配记录失败:', error)
  } finally {
    allocationsLoading.value = false
  }
}

const loadAllData = () => {
  loadStats()
  loadGPUTrend()
  loadRecentAllocations()
}

// 格式化时间显示
const formatTime = (timeStr: string) => {
  if (!timeStr) return '-'
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

// 分配状态标签类型
const allocationStatusType = (status: string) => {
  const map: Record<string, string> = {
    active: 'success',
    expired: 'info',
    reclaimed: 'warning',
  }
  return (map[status] || 'info') as 'success' | 'info' | 'warning' | 'danger'
}

// 自动刷新
let refreshTimer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  loadAllData()
  // 每 60 秒自动刷新统计数据
  refreshTimer = setInterval(loadStats, 60000)
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
})
</script>

<template>
  <div class="admin-dashboard">
    <div class="page-header">
      <h2 class="page-title">管理后台首页</h2>
      <p class="page-description">欢迎回来,这是您的管理后台概览</p>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <StatCard
        title="总机器数"
        :value="stats.total_machines"
        icon="💻"
        color="primary"
        :loading="statsLoading"
      />
      <StatCard
        title="在线机器"
        :value="onlineMachines()"
        icon="✅"
        color="success"
        :loading="statsLoading"
      />
      <StatCard
        title="已分配机器"
        :value="stats.allocated_machines"
        icon="🔗"
        color="warning"
        :loading="statsLoading"
      />
      <StatCard
        title="活跃客户"
        :value="stats.active_customers"
        icon="👥"
        color="info"
        :loading="statsLoading"
      />
    </div>

    <!-- 图表和列表区域 -->
    <el-row :gutter="20" class="content-row">
      <!-- GPU 使用趋势 -->
      <el-col :span="12">
        <el-card class="content-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">GPU 使用趋势（24h）</span>
            </div>
          </template>
          <el-skeleton :loading="trendLoading" :rows="5" animated>
            <div v-if="gpuTrend.length > 0" class="gpu-trend-chart">
              <div class="trend-bars">
                <div
                  v-for="(point, index) in gpuTrend"
                  :key="index"
                  class="trend-bar-item"
                >
                  <div class="trend-bar-wrapper">
                    <div
                      class="trend-bar"
                      :style="{ height: point.usage + '%' }"
                      :title="point.usage + '%'"
                    />
                  </div>
                  <span class="trend-bar-label">{{ point.time }}</span>
                </div>
              </div>
            </div>
            <el-empty v-else description="暂无趋势数据" />
          </el-skeleton>
        </el-card>
      </el-col>

      <!-- 最近分配记录 -->
      <el-col :span="12">
        <el-card class="content-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">最近分配</span>
              <router-link to="/admin/allocations/list">
                <el-link type="primary" :underline="false">查看全部</el-link>
              </router-link>
            </div>
          </template>
          <el-skeleton :loading="allocationsLoading" :rows="5" animated>
            <div v-if="recentAllocationList.length > 0" class="allocation-list">
              <div
                v-for="alloc in recentAllocationList"
                :key="alloc.id"
                class="allocation-item"
              >
                <div class="allocation-info">
                  <span class="allocation-machine">{{ alloc.host?.name || alloc.host_id }}</span>
                  <span class="allocation-arrow">→</span>
                  <span class="allocation-customer">{{ alloc.customer?.display_name || alloc.customer?.username || '-' }}</span>
                </div>
                <div class="allocation-meta">
                  <el-tag :type="allocationStatusType(alloc.status)" size="small">
                    {{ alloc.status }}
                  </el-tag>
                  <span class="allocation-time">{{ formatTime(alloc.created_at) }}</span>
                </div>
              </div>
            </div>
            <el-empty v-else description="暂无分配记录" />
          </el-skeleton>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.admin-dashboard {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px 0;
}

.page-description {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.content-row {
  margin-bottom: 24px;
}

.content-card {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

/* GPU 趋势图 */
.gpu-trend-chart {
  padding: 8px 0;
}

.trend-bars {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  height: 160px;
}

.trend-bar-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.trend-bar-wrapper {
  flex: 1;
  width: 100%;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.trend-bar {
  width: 70%;
  min-height: 2px;
  background: linear-gradient(180deg, #409eff, #79bbff);
  border-radius: 3px 3px 0 0;
  transition: height 0.3s;
}

.trend-bar-label {
  font-size: 11px;
  color: #909399;
  margin-top: 6px;
  white-space: nowrap;
}

/* 最近分配记录 */
.allocation-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.allocation-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-bottom: 1px solid #ebeef5;
}

.allocation-item:last-child {
  border-bottom: none;
}

.allocation-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.allocation-machine {
  color: #303133;
  font-weight: 500;
}

.allocation-arrow {
  color: #c0c4cc;
}

.allocation-customer {
  color: #606266;
}

.allocation-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.allocation-time {
  font-size: 12px;
  color: #909399;
}
</style>
