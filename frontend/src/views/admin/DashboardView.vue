<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getDashboardOverview } from '@/api/admin'
import StatCard from '@/components/common/StatCard.vue'

interface DashboardData {
  totalMachines: number
  onlineMachines: number
  allocatedMachines: number
  totalCustomers: number
  machinesTrend: number
  customersTrend: number
  recentAlerts: any[]
  recentActivities: any[]
}

const loading = ref(true)
const dashboardData = ref<DashboardData>({
  totalMachines: 0,
  onlineMachines: 0,
  allocatedMachines: 0,
  totalCustomers: 0,
  machinesTrend: 0,
  customersTrend: 0,
  recentAlerts: [],
  recentActivities: []
})

const loadDashboardData = async () => {
  try {
    loading.value = true
    const response = await getDashboardOverview()
    dashboardData.value = response.data
  } catch (error) {
    console.error('加载Dashboard数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDashboardData()
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
        :value="dashboardData.totalMachines"
        icon="💻"
        color="primary"
        :trend="dashboardData.machinesTrend"
        :loading="loading"
      />
      <StatCard
        title="在线机器"
        :value="dashboardData.onlineMachines"
        icon="✅"
        color="success"
        :loading="loading"
      />
      <StatCard
        title="已分配机器"
        :value="dashboardData.allocatedMachines"
        icon="🔗"
        color="warning"
        :loading="loading"
      />
      <StatCard
        title="客户数量"
        :value="dashboardData.totalCustomers"
        icon="👥"
        color="info"
        :trend="dashboardData.customersTrend"
        :loading="loading"
      />
    </div>

    <!-- 图表和列表区域 -->
    <el-row :gutter="20" class="content-row">
      <!-- 最近告警 -->
      <el-col :span="12">
        <el-card class="content-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">最近告警</span>
              <el-link type="primary" :underline="false">查看全部</el-link>
            </div>
          </template>
          <el-skeleton :loading="loading" :rows="5" animated>
            <div v-if="dashboardData.recentAlerts.length > 0" class="alert-list">
              <div
                v-for="alert in dashboardData.recentAlerts"
                :key="alert.id"
                class="alert-item"
              >
                <el-tag :type="alert.level === 'critical' ? 'danger' : 'warning'" size="small">
                  {{ alert.level }}
                </el-tag>
                <span class="alert-message">{{ alert.message }}</span>
                <span class="alert-time">{{ alert.time }}</span>
              </div>
            </div>
            <el-empty v-else description="暂无告警" />
          </el-skeleton>
        </el-card>
      </el-col>

      <!-- 最近活动 -->
      <el-col :span="12">
        <el-card class="content-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">最近活动</span>
              <el-link type="primary" :underline="false">查看全部</el-link>
            </div>
          </template>
          <el-skeleton :loading="loading" :rows="5" animated>
            <div v-if="dashboardData.recentActivities.length > 0" class="activity-list">
              <div
                v-for="activity in dashboardData.recentActivities"
                :key="activity.id"
                class="activity-item"
              >
                <div class="activity-icon">{{ activity.icon }}</div>
                <div class="activity-content">
                  <div class="activity-title">{{ activity.title }}</div>
                  <div class="activity-time">{{ activity.time }}</div>
                </div>
              </div>
            </div>
            <el-empty v-else description="暂无活动" />
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

.alert-list,
.activity-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.alert-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.alert-message {
  flex: 1;
  font-size: 14px;
  color: #606266;
}

.alert-time {
  font-size: 12px;
  color: #909399;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  border-bottom: 1px solid #ebeef5;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  font-size: 24px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  border-radius: 50%;
}

.activity-content {
  flex: 1;
}

.activity-title {
  font-size: 14px;
  color: #303133;
  margin-bottom: 4px;
}

.activity-time {
  font-size: 12px;
  color: #909399;
}
</style>
