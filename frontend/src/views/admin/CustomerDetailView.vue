<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getCustomerDetail } from '@/api/admin'
import type { CustomerDetail } from '@/types/customer'
import StatCard from '@/components/common/StatCard.vue'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const customerDetail = ref<CustomerDetail | null>(null)

const loadCustomerDetail = async () => {
  try {
    loading.value = true
    const customerId = Number(route.params.id)
    const response = await getCustomerDetail(customerId)
    customerDetail.value = response.data
  } catch (error) {
    console.error('加载客户详情失败:', error)
  } finally {
    loading.value = false
  }
}

const handleBack = () => {
  router.back()
}

onMounted(() => {
  loadCustomerDetail()
})
</script>

<template>
  <div class="customer-detail">
    <div class="page-header">
      <div>
        <el-button @click="handleBack">返回</el-button>
        <h2 class="page-title">客户详情</h2>
      </div>
    </div>

    <el-skeleton :loading="loading" :rows="10" animated>
      <div v-if="customerDetail">
        <!-- 基本信息 -->
        <el-card class="info-card">
          <template #header>
            <span class="card-title">基本信息</span>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="客户名称">
              {{ customerDetail.name }}
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="customerDetail.status === 'active' ? 'success' : 'danger'">
                {{ customerDetail.status === 'active' ? '正常' : '已停用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="联系人">
              {{ customerDetail.contactPerson }}
            </el-descriptions-item>
            <el-descriptions-item label="联系邮箱">
              {{ customerDetail.contactEmail }}
            </el-descriptions-item>
            <el-descriptions-item label="联系电话">
              {{ customerDetail.contactPhone }}
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ customerDetail.createdAt }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- 使用统计 -->
        <div class="stats-grid">
          <StatCard
            title="分配机器数"
            :value="customerDetail.usageStats?.allocatedMachines || 0"
            icon="💻"
            color="primary"
          />
          <StatCard
            title="运行任务数"
            :value="customerDetail.usageStats?.runningTasks || 0"
            icon="🚀"
            color="success"
          />
          <StatCard
            title="总任务数"
            :value="customerDetail.usageStats?.totalTasks || 0"
            icon="📊"
            color="info"
          />
          <StatCard
            title="存储使用(GB)"
            :value="customerDetail.usageStats?.storageUsed || 0"
            icon="💾"
            color="warning"
          />
        </div>

        <!-- 分配的机器 -->
        <el-card class="machines-card">
          <template #header>
            <span class="card-title">分配的机器</span>
          </template>
          <el-table :data="customerDetail.allocatedMachines" stripe border>
            <el-table-column prop="machineName" label="机器名称" min-width="150" />
            <el-table-column prop="region" label="区域" width="120" />
            <el-table-column prop="allocatedAt" label="分配时间" width="180" />
            <el-table-column prop="expiresAt" label="到期时间" width="180" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'warning'">
                  {{ row.status === 'active' ? '使用中' : '即将到期' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <!-- 操作日志 -->
        <el-card class="logs-card">
          <template #header>
            <span class="card-title">操作日志</span>
          </template>
          <el-timeline>
            <el-timeline-item
              v-for="log in customerDetail.operationLogs"
              :key="log.id"
              :timestamp="log.timestamp"
              placement="top"
            >
              <div class="log-content">
                <div class="log-action">{{ log.action }}</div>
                <div class="log-operator">操作人: {{ log.operator }}</div>
              </div>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </div>
    </el-skeleton>
  </div>
</template>

<style scoped>
.customer-detail {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 8px 0 0 0;
}

.info-card,
.machines-card,
.logs-card {
  margin-bottom: 20px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.log-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.log-action {
  font-size: 14px;
  color: #303133;
}

.log-operator {
  font-size: 12px;
  color: #909399;
}
</style>
