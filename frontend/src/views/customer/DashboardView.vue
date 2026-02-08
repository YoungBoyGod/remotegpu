<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getDashboardOverview } from '@/api/customer'
import type { DashboardOverview } from '@/api/customer'
import type { Task } from '@/types/task'
import StatCard from '@/components/common/StatCard.vue'

const router = useRouter()

const loading = ref(true)
const dashboardData = ref<DashboardOverview>({
  myMachines: 0,
  runningTasks: 0,
  totalTasks: 0,
  datasetCount: 0,
  recentTasks: [],
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
  { title: '工作空间', icon: '📂', path: '/customer/workspaces' },
  { title: '开发环境', icon: '🖥️', path: '/customer/environments' },
  { title: '数据集', icon: '📦', path: '/customer/datasets' },
  { title: 'SSH 密钥', icon: '🔑', path: '/customer/ssh-keys' },
]

const handleQuickAction = (path: string) => {
  router.push(path)
}

// 任务状态标签
const taskStatusType = (status: string) => {
  const map: Record<string, string> = {
    running: 'success',
    pending: 'warning',
    completed: '',
    failed: 'danger',
    cancelled: 'info',
    stopped: 'info',
  }
  return (map[status] || 'info') as '' | 'success' | 'warning' | 'danger' | 'info'
}

const taskStatusLabel = (status: string) => {
  const map: Record<string, string> = {
    running: '运行中',
    pending: '等待中',
    completed: '已完成',
    failed: '失败',
    cancelled: '已取消',
    stopped: '已停止',
  }
  return map[status] || status
}

const formatDate = (value?: string) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN')
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
        title="数据集数量"
        :value="dashboardData.datasetCount"
        icon="📦"
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

    <!-- 最近任务 -->
    <el-card class="activities-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">最近任务</span>
          <router-link to="/customer/tasks">
            <el-link type="primary" :underline="false">查看全部</el-link>
          </router-link>
        </div>
      </template>
      <el-skeleton :loading="loading" :rows="5" animated>
        <el-table v-if="dashboardData.recentTasks.length > 0" :data="dashboardData.recentTasks" stripe>
          <el-table-column prop="name" label="任务名称" min-width="160" show-overflow-tooltip />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="taskStatusType(row.status)" size="small">
                {{ taskStatusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="机器" width="140" show-overflow-tooltip>
            <template #default="{ row }">
              {{ row.host?.name || row.host_id || '-' }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="175">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
        </el-table>
        <el-empty v-else description="暂无任务记录" />
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

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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
</style>
