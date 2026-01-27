<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/common/PageHeader.vue'
import MetricCard from '@/components/dashboard/MetricCard.vue'

// 资源统计数据
const metrics = ref({
  hostResources: {
    total: 10,
    physical: 10,
    virtual: 0,
    cloud: 0
  },
  modelTotal: {
    total: 63,
    customer: 55,
    system: 3,
    organization: 5
  },
  organizationResources: {
    total: 0,
    users: 0
  },
  applications: {
    total: 0
  }
})

// 应用资源使用排名
const resourceUsage = ref([])
const loading = ref(false)

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取数据
    // 模拟数据
    await new Promise(resolve => setTimeout(resolve, 500))
  } catch (error) {
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="resource-platform">
    <PageHeader title="资源平台" />

    <!-- 统计指标卡片 -->
    <div class="metrics-container">
      <div class="metric-card-wrapper">
        <div class="metric-card-header">
          <h3>主机资源</h3>
          <span class="metric-total">{{ metrics.hostResources.total }}</span>
        </div>
        <div class="metric-details">
          <div class="detail-item">
            <span class="label">物理机:</span>
            <span class="value">{{ metrics.hostResources.physical }}</span>
          </div>
          <div class="detail-item">
            <span class="label">虚拟机:</span>
            <span class="value">{{ metrics.hostResources.virtual }}</span>
          </div>
          <div class="detail-item">
            <span class="label">云主机:</span>
            <span class="value">{{ metrics.hostResources.cloud }}</span>
          </div>
        </div>
      </div>

      <div class="metric-card-wrapper">
        <div class="metric-card-header">
          <h3>模型总数</h3>
          <span class="metric-total">{{ metrics.modelTotal.total }}</span>
        </div>
        <div class="metric-details">
          <div class="detail-item">
            <span class="label">客户模型:</span>
            <span class="value">{{ metrics.modelTotal.customer }}</span>
          </div>
          <div class="detail-item">
            <span class="label">系统模型:</span>
            <span class="value">{{ metrics.modelTotal.system }}</span>
          </div>
          <div class="detail-item">
            <span class="label">组织模型:</span>
            <span class="value">{{ metrics.modelTotal.organization }}</span>
          </div>
        </div>
      </div>

      <div class="metric-card-wrapper">
        <div class="metric-card-header">
          <h3>组织资源</h3>
          <span class="metric-total">{{ metrics.organizationResources.total }}</span>
        </div>
        <div class="metric-details">
          <div class="detail-item">
            <span class="label">用户:</span>
            <span class="value">{{ metrics.organizationResources.users }}</span>
          </div>
        </div>
      </div>

      <div class="metric-card-wrapper">
        <div class="metric-card-header">
          <h3>应用</h3>
          <span class="metric-total">{{ metrics.applications.total }}</span>
        </div>
      </div>
    </div>

    <!-- 数据表格和图表区域 -->
    <div class="content-container">
      <el-card class="table-card">
        <template #header>
          <div class="card-header">
            <span>应用资源使用排名</span>
            <el-link type="primary">更多</el-link>
          </div>
        </template>
        <el-table :data="resourceUsage" :loading="loading" style="width: 100%">
          <el-table-column prop="name" label="应用名称" />
          <el-table-column prop="business" label="所属业务" />
          <el-table-column prop="hosts" label="主机数" />
          <el-table-column prop="services" label="服务数" />
        </el-table>
        <el-empty v-if="!loading && resourceUsage.length === 0" description="暂无数据" />
      </el-card>

      <el-card class="chart-card">
        <template #header>
          <div class="card-header">
            <span>资源占比统计</span>
            <el-link type="primary">更多</el-link>
          </div>
        </template>
        <div class="chart-container">
          <div class="chart-placeholder">
            <p>图表区域</p>
            <p class="chart-total">总数量<br />12台</p>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<style scoped>
.resource-platform {
  padding: 24px;
}

.metrics-container {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.metric-card-wrapper {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.metric-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.metric-card-header h3 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.metric-total {
  font-size: 32px;
  font-weight: 600;
  color: #409EFF;
}

.metric-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.detail-item .label {
  color: #909399;
}

.detail-item .value {
  color: #303133;
  font-weight: 500;
}

.content-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.table-card,
.chart-card {
  min-height: 400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.chart-placeholder {
  text-align: center;
  color: #909399;
}

.chart-total {
  margin-top: 20px;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

@media (max-width: 1200px) {
  .metrics-container {
    grid-template-columns: repeat(2, 1fr);
  }

  .content-container {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .metrics-container {
    grid-template-columns: 1fr;
  }
}
</style>
