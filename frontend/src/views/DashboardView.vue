<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import MetricCard from '@/components/dashboard/MetricCard.vue'
import ResourceChart from '@/components/dashboard/ResourceChart.vue'
import RecentEnvironments from '@/components/dashboard/RecentEnvironments.vue'

// 加载状态
const loading = ref(false)

// 概览数据
const overview = ref({
  runningEnvs: 0,
  totalGpuHours: 0,
  monthlyBilling: 0,
  storageUsed: 0,
})

// 加载仪表盘数据
const loadDashboardData = async () => {
  loading.value = true
  try {
    // TODO: 调用 API 获取数据
    // const response = await getDashboardOverview()
    // overview.value = response.data

    // 模拟数据
    overview.value = {
      runningEnvs: 3,
      totalGpuHours: 156.5,
      monthlyBilling: 1280.50,
      storageUsed: 45.2,
    }
  } catch (error) {
    ElMessage.error('加载仪表盘数据失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDashboardData()
})
</script>

<template>
  <div class="dashboard-container">
    <div class="dashboard-header">
      <h1>概览</h1>
      <p class="subtitle">欢迎回来，查看您的资源使用情况</p>
    </div>

    <!-- 数据卡片 -->
    <el-row :gutter="20" class="metric-cards">
      <el-col :xs="24" :sm="12" :lg="6">
        <MetricCard
          title="运行中的环境"
          :value="overview.runningEnvs"
          unit="个"
          icon="Monitor"
          color="#409EFF"
        />
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <MetricCard
          title="GPU 使用时长"
          :value="overview.totalGpuHours"
          unit="小时"
          icon="Timer"
          color="#67C23A"
        />
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <MetricCard
          title="本月费用"
          :value="overview.monthlyBilling"
          unit="元"
          icon="Money"
          color="#E6A23C"
        />
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <MetricCard
          title="存储使用"
          :value="overview.storageUsed"
          unit="GB"
          icon="FolderOpened"
          color="#F56C6C"
        />
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="chart-section">
      <el-col :xs="24" :lg="16">
        <ResourceChart />
      </el-col>
      <el-col :xs="24" :lg="8">
        <div class="quick-actions">
          <h3>快速操作</h3>
          <el-button type="primary" size="large" class="action-btn">
            创建环境
          </el-button>
          <el-button type="success" size="large" class="action-btn">
            上传数据集
          </el-button>
        </div>
      </el-col>
    </el-row>

    <!-- 最近活动 -->
    <el-row :gutter="20" class="recent-section">
      <el-col :xs="24">
        <RecentEnvironments />
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.dashboard-container {
  padding: 24px;
}

.dashboard-header {
  margin-bottom: 24px;
}

.dashboard-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px 0;
}

.subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.metric-cards {
  margin-bottom: 24px;
}

.chart-section {
  margin-bottom: 24px;
}

.quick-actions {
  background: white;
  padding: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  height: 100%;
}

.quick-actions h3 {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 20px 0;
}

.action-btn {
  width: 100%;
  margin-bottom: 12px;
}

.recent-section {
  margin-bottom: 24px;
}
</style>
