<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getDashboardOverview } from '@/api/customer'
import StatCard from '@/components/common/StatCard.vue'

const router = useRouter()

interface DashboardData {
  myMachines: number
  runningTasks: number
  totalTasks: number
  storageUsed: number
  recentActivities: any[]
}

const loading = ref(true)
const dashboardData = ref<DashboardData>({
  myMachines: 0,
  runningTasks: 0,
  totalTasks: 0,
  storageUsed: 0,
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

const quickActions = [
  { title: '我的机器', icon: '💻', path: '/customer/machines/list' },
  { title: '创建任务', icon: '🚀', path: '/customer/tasks/training' },
  { title: '镜像市场', icon: '🐳', path: '/customer/images/market' },
  { title: '上传数据集', icon: '📤', path: '/customer/datasets/upload' }
]

const handleQuickAction = (path: string) => {
  router.push(path)
}

onMounted(() => {
  loadDashboardData()
})
</script>

<template>
  <div class="customer-dashboard">
    <div class="page-header">
      <h2 class="page-title">工作台首页</h2>
      <p class="page-description">欢迎回来,开始您的AI训练之旅</p>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <StatCard
        title="我的机器"
        :value="dashboardData.myMachines"
        icon="💻"
        color="primary"
        :loading="loading"
      />
      <StatCard
        title="运行中任务"
        :value="dashboardData.runningTasks"
        icon="🚀"
        color="success"
        :loading="loading"
      />
      <StatCard
        title="总任务数"
        :value="dashboardData.totalTasks"
        icon="📊"
        color="info"
        :loading="loading"
      />
      <StatCard
        title="存储使用(GB)"
        :value="dashboardData.storageUsed"
        icon="💾"
        color="warning"
        :loading="loading"
      />
    </div>

    <!-- 快捷操作 -->
    <el-card class="quick-actions-card">
      <template #header>
        <span class="card-title">快捷操作</span>
      </template>
      <div class="quick-actions">
        <div
          v-for="action in quickActions"
          :key="action.path"
          class="action-item"
          @click="handleQuickAction(action.path)"
        >
          <div class="action-icon">{{ action.icon }}</div>
          <div class="action-title">{{ action.title }}</div>
        </div>
      </div>
    </el-card>

    <!-- 最近活动 -->
    <el-card class="activities-card">
      <template #header>
        <span class="card-title">最近活动</span>
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
  </div>
</template>

<style scoped>
.customer-dashboard {
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

.quick-actions-card,
.activities-card {
  margin-bottom: 24px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 16px;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px;
  background: #f5f7fa;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.action-item:hover {
  background: #e6f7ff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.action-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.action-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
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
