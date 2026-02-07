<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getCustomerDetail } from '@/api/admin'
import type { Customer } from '@/types/customer'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const customerDetail = ref<Customer | null>(null)

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

const formatDateTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN')
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
            <el-descriptions-item label="用户名">
              {{ customerDetail.username }}
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="customerDetail.status === 'active' ? 'success' : 'danger'">
                {{ customerDetail.status === 'active' ? '正常' : '已停用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="公司代号">
              {{ customerDetail.company_code || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="公司名称">
              {{ customerDetail.company || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="联系人">
              {{ customerDetail.full_name || customerDetail.display_name || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="邮箱">
              {{ customerDetail.email }}
            </el-descriptions-item>
            <el-descriptions-item label="电话">
              {{ customerDetail.phone || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ formatDateTime(customerDetail.created_at) }}
            </el-descriptions-item>
          </el-descriptions>
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

.info-card {
  margin-bottom: 20px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

</style>
