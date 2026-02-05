<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Warning,
  InfoFilled,
  CircleCheck,
  Search
} from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { getAlertList, acknowledgeAlert } from '@/api/admin'

interface Alert {
  id: number
  rule_id: number
  host_id: string
  value: number
  message: string
  triggered_at: string
  acknowledged: boolean
  rule: {
    name: string
    severity: string
    metric_type: string
  }
}

const loading = ref(false)
const alerts = ref<Alert[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const selectedSeverity = ref<string>('')
const selectedAcknowledged = ref<string>('false') // Default show active

// 加载告警列表
const loadAlerts = async () => {
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      pageSize: pageSize.value,
    }
    if (selectedSeverity.value) params.severity = selectedSeverity.value
    if (selectedAcknowledged.value !== '') {
      params.acknowledged = selectedAcknowledged.value === 'true'
    }

    const res = await getAlertList(params)
    alerts.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error(error)
    ElMessage.error('加载告警列表失败')
  } finally {
    loading.value = false
  }
}

// 确认告警
const handleAcknowledge = async (alert: Alert) => {
  try {
    await ElMessageBox.confirm('确认已知晓此告警？', '确认告警', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await acknowledgeAlert(alert.id)
    ElMessage.success('告警已确认')
    loadAlerts()
  } catch {
    // 用户取消
  }
}

const handlePageChange = (val: number) => {
  page.value = val
  loadAlerts()
}

// 获取级别类型
const getSeverityType = (level: string) => {
  const map: Record<string, string> = {
    critical: 'danger',
    warning: 'warning',
    info: 'info'
  }
  return map[level] || 'info'
}

onMounted(() => {
  loadAlerts()
})
</script>

<template>
  <div class="alert-center">
    <PageHeader title="告警中心" />

    <!-- 筛选区域 -->
    <el-card class="filter-card">
      <div class="filter-container">
        <el-select v-model="selectedSeverity" placeholder="告警级别" clearable style="width: 150px" @change="loadAlerts">
          <el-option label="严重" value="critical" />
          <el-option label="警告" value="warning" />
          <el-option label="信息" value="info" />
        </el-select>

        <el-select v-model="selectedAcknowledged" placeholder="告警状态" style="width: 150px" @change="loadAlerts">
          <el-option label="未确认" value="false" />
          <el-option label="已确认" value="true" />
          <el-option label="全部" value="" />
        </el-select>

        <el-button type="primary" :icon="Search" @click="loadAlerts">刷新</el-button>
      </div>
    </el-card>

    <!-- 告警列表 -->
    <el-card class="alert-list-card">
      <div v-loading="loading" class="alert-list">
        <div
          v-for="alert in alerts"
          :key="alert.id"
          class="alert-item"
          :class="`alert-${alert.rule?.severity || 'info'}`"
        >
          <div class="alert-header">
            <div class="alert-title-section">
              <el-tag :type="getSeverityType(alert.rule?.severity)" size="small">
                {{ alert.rule?.severity || 'unknown' }}
              </el-tag>
              <h4 class="alert-title">{{ alert.rule?.name || 'Unknown Rule' }}</h4>
            </div>
            <el-tag v-if="alert.acknowledged" type="success" size="small">已确认</el-tag>
            <el-tag v-else type="danger" size="small">未确认</el-tag>
          </div>

          <p class="alert-description">{{ alert.message }} (Value: {{ alert.value }})</p>

          <div class="alert-meta">
            <span class="meta-item">主机: {{ alert.host_id }}</span>
            <span class="meta-item">触发时间: {{ new Date(alert.triggered_at).toLocaleString() }}</span>
          </div>

          <div class="alert-actions" v-if="!alert.acknowledged">
            <el-button
              size="small"
              type="warning"
              @click="handleAcknowledge(alert)"
            >
              确认
            </el-button>
          </div>
        </div>

        <el-empty
          v-if="!loading && alerts.length === 0"
          description="暂无告警信息"
        />
      </div>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.alert-center {
  padding: 24px;
}

.filter-card {
  margin-bottom: 24px;
}

.filter-container {
  display: flex;
  gap: 12px;
}

.alert-list {
  min-height: 400px;
}

.alert-item {
  padding: 20px;
  margin-bottom: 16px;
  border-radius: 8px;
  border-left: 4px solid #dcdfe6;
  background: #f8f9fa;
  transition: all 0.3s;
}

.alert-critical { border-left-color: #f56c6c; background: #fef0f0; }
.alert-warning { border-left-color: #e6a23c; background: #fdf6ec; }
.alert-info { border-left-color: #409eff; background: #ecf5ff; }

.alert-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.alert-title-section {
  display: flex;
  align-items: center;
  gap: 8px;
}

.alert-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.alert-description {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #606266;
}

.alert-meta {
  display: flex;
  gap: 16px;
  color: #909399;
  font-size: 13px;
  margin-bottom: 12px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
