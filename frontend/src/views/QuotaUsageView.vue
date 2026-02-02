<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Warning } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { getQuotaUsage } from '@/api/quota'
import type { QuotaUsageResponse } from '@/api/quota/types'

const loading = ref(false)
const usageData = ref<QuotaUsageResponse | null>(null)
const customerId = ref<number>(1) // TODO: 从用户信息中获取
const workspaceId = ref<number | undefined>(undefined)

// 加载配额使用情况
const loadUsage = async () => {
  loading.value = true
  try {
    const response = await getQuotaUsage(customerId.value, workspaceId.value)
    usageData.value = response.data
  } catch (error: any) {
    ElMessage.error(error.message || '加载配额使用情况失败')
  } finally {
    loading.value = false
  }
}

// 刷新数据
const handleRefresh = () => {
  loadUsage()
}

// 获取进度条状态
const getProgressStatus = (percentage: number) => {
  if (percentage >= 90) return 'exception'
  if (percentage >= 80) return 'warning'
  return 'success'
}

// 格式化内存大小(MB)
const formatMemory = (mb: number) => {
  if (mb >= 1024) {
    return `${(mb / 1024).toFixed(1)} GB`
  }
  return `${mb} MB`
}

// 格式化存储大小(GB)
const formatStorage = (gb: number) => {
  if (gb >= 1024) {
    return `${(gb / 1024).toFixed(1)} TB`
  }
  return `${gb} GB`
}

// 是否显示告警
const showWarning = computed(() => {
  if (!usageData.value) return false
  const { usage_percentage } = usageData.value
  return (
    usage_percentage.gpu >= 80 ||
    usage_percentage.cpu >= 80 ||
    usage_percentage.memory >= 80 ||
    usage_percentage.storage >= 80 ||
    usage_percentage.environments >= 80
  )
})

onMounted(() => {
  loadUsage()
})
</script>

<template>
  <div class="quota-usage">
    <PageHeader title="配额使用情况">
      <template #actions>
        <el-button :icon="Refresh" @click="handleRefresh" :loading="loading">
          刷新
        </el-button>
      </template>
    </PageHeader>

    <el-alert
      v-if="showWarning"
      title="资源使用告警"
      type="warning"
      :icon="Warning"
      show-icon
      :closable="false"
      style="margin-bottom: 20px"
    >
      <template #default>
        部分资源使用率已超过80%,请注意及时释放资源或申请扩容
      </template>
    </el-alert>

    <el-card v-loading="loading">
      <template v-if="usageData">
        <!-- GPU使用情况 -->
        <div class="resource-item">
          <div class="resource-header">
            <span class="resource-name">GPU</span>
            <span class="resource-value">
              {{ usageData.used.used_gpu }} / {{ usageData.quota.max_gpu }} 个
              (可用: {{ usageData.available.available_gpu }} 个)
            </span>
          </div>
          <el-progress
            :percentage="usageData.usage_percentage.gpu"
            :status="getProgressStatus(usageData.usage_percentage.gpu)"
            :stroke-width="20"
          />
        </div>

        <!-- CPU使用情况 -->
        <div class="resource-item">
          <div class="resource-header">
            <span class="resource-name">CPU</span>
            <span class="resource-value">
              {{ usageData.used.used_cpu }} / {{ usageData.quota.max_cpu }} 核
              (可用: {{ usageData.available.available_cpu }} 核)
            </span>
          </div>
          <el-progress
            :percentage="usageData.usage_percentage.cpu"
            :status="getProgressStatus(usageData.usage_percentage.cpu)"
            :stroke-width="20"
          />
        </div>

        <!-- 内存使用情况 -->
        <div class="resource-item">
          <div class="resource-header">
            <span class="resource-name">内存</span>
            <span class="resource-value">
              {{ formatMemory(usageData.used.used_memory) }} /
              {{ formatMemory(usageData.quota.max_memory) }}
              (可用: {{ formatMemory(usageData.available.available_memory) }})
            </span>
          </div>
          <el-progress
            :percentage="usageData.usage_percentage.memory"
            :status="getProgressStatus(usageData.usage_percentage.memory)"
            :stroke-width="20"
          />
        </div>

        <!-- 存储使用情况 -->
        <div class="resource-item">
          <div class="resource-header">
            <span class="resource-name">存储</span>
            <span class="resource-value">
              {{ formatStorage(usageData.used.used_storage) }} /
              {{ formatStorage(usageData.quota.max_storage) }}
              (可用: {{ formatStorage(usageData.available.available_storage) }})
            </span>
          </div>
          <el-progress
            :percentage="usageData.usage_percentage.storage"
            :status="getProgressStatus(usageData.usage_percentage.storage)"
            :stroke-width="20"
          />
        </div>

        <!-- 环境数量使用情况 -->
        <div class="resource-item">
          <div class="resource-header">
            <span class="resource-name">环境数量</span>
            <span class="resource-value">
              {{ usageData.used.used_environments }} /
              {{ usageData.quota.max_environments }} 个
              (可用: {{ usageData.available.available_environments }} 个)
            </span>
          </div>
          <el-progress
            :percentage="usageData.usage_percentage.environments"
            :status="getProgressStatus(usageData.usage_percentage.environments)"
            :stroke-width="20"
          />
        </div>
      </template>

      <el-empty v-else description="暂无配额数据" />
    </el-card>
  </div>
</template>

<style scoped>
.quota-usage {
  padding: 24px;
}

.resource-item {
  margin-bottom: 32px;
}

.resource-item:last-child {
  margin-bottom: 0;
}

.resource-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.resource-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.resource-value {
  font-size: 14px;
  color: #606266;
}
</style>
