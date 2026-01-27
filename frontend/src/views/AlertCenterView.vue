<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Bell,
  Warning,
  InfoFilled,
  CircleCheck,
  Delete,
  View,
  Search
} from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'

interface Alert {
  id: number
  title: string
  description: string
  level: 'critical' | 'warning' | 'info'
  type: 'system' | 'resource' | 'security' | 'performance'
  status: 'active' | 'acknowledged' | 'resolved'
  source: string
  createdAt: string
  acknowledgedAt?: string
  resolvedAt?: string
}

const loading = ref(false)
const alerts = ref<Alert[]>([])
const searchKeyword = ref('')
const selectedLevel = ref<string>('all')
const selectedType = ref<string>('all')
const selectedStatus = ref<string>('all')

// 告警统计
const alertStats = computed(() => {
  return {
    total: alerts.value.length,
    critical: alerts.value.filter(a => a.level === 'critical').length,
    warning: alerts.value.filter(a => a.level === 'warning').length,
    info: alerts.value.filter(a => a.level === 'info').length,
    active: alerts.value.filter(a => a.status === 'active').length,
    acknowledged: alerts.value.filter(a => a.status === 'acknowledged').length,
    resolved: alerts.value.filter(a => a.status === 'resolved').length
  }
})

// 过滤后的告警列表
const filteredAlerts = computed(() => {
  let result = alerts.value

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(a =>
      a.title.toLowerCase().includes(keyword) ||
      a.description.toLowerCase().includes(keyword) ||
      a.source.toLowerCase().includes(keyword)
    )
  }

  // 级别筛选
  if (selectedLevel.value !== 'all') {
    result = result.filter(a => a.level === selectedLevel.value)
  }

  // 类型筛选
  if (selectedType.value !== 'all') {
    result = result.filter(a => a.type === selectedType.value)
  }

  // 状态筛选
  if (selectedStatus.value !== 'all') {
    result = result.filter(a => a.status === selectedStatus.value)
  }

  return result
})

// 加载告警列表
const loadAlerts = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 500))
    alerts.value = [
      {
        id: 1,
        title: 'GPU温度过高',
        description: '主机-001的GPU温度达到85℃，超过安全阈值',
        level: 'critical',
        type: 'resource',
        status: 'active',
        source: '主机-001',
        createdAt: '2026-01-27 14:30:00'
      },
      {
        id: 2,
        title: 'CPU使用率告警',
        description: '主机-002的CPU使用率持续超过90%',
        level: 'warning',
        type: 'performance',
        status: 'acknowledged',
        source: '主机-002',
        createdAt: '2026-01-27 13:15:00',
        acknowledgedAt: '2026-01-27 13:20:00'
      },
      {
        id: 3,
        title: '磁盘空间不足',
        description: '主机-003的磁盘使用率达到95%',
        level: 'warning',
        type: 'resource',
        status: 'resolved',
        source: '主机-003',
        createdAt: '2026-01-27 10:00:00',
        resolvedAt: '2026-01-27 12:00:00'
      },
      {
        id: 4,
        title: '系统更新可用',
        description: '检测到新的系统更新版本',
        level: 'info',
        type: 'system',
        status: 'active',
        source: '系统',
        createdAt: '2026-01-27 09:00:00'
      }
    ]
  } catch (error) {
    ElMessage.error('加载告警列表失败')
  } finally {
    loading.value = false
  }
}

// 获取级别类型
const getLevelType = (level: string) => {
  const levelMap = {
    critical: 'danger',
    warning: 'warning',
    info: 'info'
  }
  return levelMap[level as keyof typeof levelMap] || 'info'
}

// 获取级别文本
const getLevelText = (level: string) => {
  const levelMap = {
    critical: '严重',
    warning: '警告',
    info: '信息'
  }
  return levelMap[level as keyof typeof levelMap] || '未知'
}

// 获取状态类型
const getStatusType = (status: string) => {
  const statusMap = {
    active: 'danger',
    acknowledged: 'warning',
    resolved: 'success'
  }
  return statusMap[status as keyof typeof statusMap] || 'info'
}

// 获取状态文本
const getStatusText = (status: string) => {
  const statusMap = {
    active: '活跃',
    acknowledged: '已确认',
    resolved: '已解决'
  }
  return statusMap[status as keyof typeof statusMap] || '未知'
}

// 获取类型文本
const getTypeText = (type: string) => {
  const typeMap = {
    system: '系统',
    resource: '资源',
    security: '安全',
    performance: '性能'
  }
  return typeMap[type as keyof typeof typeMap] || '未知'
}

// 确认告警
const handleAcknowledge = async (alert: Alert) => {
  try {
    await ElMessageBox.confirm('确认已知晓此告警？', '确认告警', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })
    // TODO: 调用API
    alert.status = 'acknowledged'
    alert.acknowledgedAt = new Date().toLocaleString('zh-CN')
    ElMessage.success('告警已确认')
  } catch {
    // 用户取消
  }
}

// 解决告警
const handleResolve = async (alert: Alert) => {
  try {
    await ElMessageBox.confirm('确认此告警已解决？', '解决告警', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'success'
    })
    // TODO: 调用API
    alert.status = 'resolved'
    alert.resolvedAt = new Date().toLocaleString('zh-CN')
    ElMessage.success('告警已解决')
  } catch {
    // 用户取消
  }
}

// 删除告警
const handleDelete = async (alert: Alert) => {
  try {
    await ElMessageBox.confirm('确认删除此告警？', '删除告警', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'error'
    })
    // TODO: 调用API
    const index = alerts.value.findIndex(a => a.id === alert.id)
    if (index > -1) {
      alerts.value.splice(index, 1)
    }
    ElMessage.success('告警已删除')
  } catch {
    // 用户取消
  }
}

// 清除筛选
const clearFilters = () => {
  searchKeyword.value = ''
  selectedLevel.value = 'all'
  selectedType.value = 'all'
  selectedStatus.value = 'all'
}

onMounted(() => {
  loadAlerts()
})
</script>

<template>
  <div class="alert-center">
    <PageHeader title="告警中心" />

    <!-- 统计卡片 -->
    <div class="stats-container">
      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #F56C6C">
            <el-icon :size="32"><Warning /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">严重告警</div>
            <div class="stat-value">{{ alertStats.critical }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #E6A23C">
            <el-icon :size="32"><Warning /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">警告</div>
            <div class="stat-value">{{ alertStats.warning }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #409EFF">
            <el-icon :size="32"><InfoFilled /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">信息</div>
            <div class="stat-value">{{ alertStats.info }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #67C23A">
            <el-icon :size="32"><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">已解决</div>
            <div class="stat-value">{{ alertStats.resolved }}</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 筛选区域 -->
    <el-card class="filter-card">
      <div class="filter-container">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索告警标题、描述或来源"
          :prefix-icon="Search"
          clearable
          style="width: 300px"
        />

        <el-select v-model="selectedLevel" placeholder="告警级别" clearable style="width: 150px">
          <el-option label="全部级别" value="all" />
          <el-option label="严重" value="critical" />
          <el-option label="警告" value="warning" />
          <el-option label="信息" value="info" />
        </el-select>

        <el-select v-model="selectedType" placeholder="告警类型" clearable style="width: 150px">
          <el-option label="全部类型" value="all" />
          <el-option label="系统" value="system" />
          <el-option label="资源" value="resource" />
          <el-option label="安全" value="security" />
          <el-option label="性能" value="performance" />
        </el-select>

        <el-select v-model="selectedStatus" placeholder="告警状态" clearable style="width: 150px">
          <el-option label="全部状态" value="all" />
          <el-option label="活跃" value="active" />
          <el-option label="已确认" value="acknowledged" />
          <el-option label="已解决" value="resolved" />
        </el-select>

        <el-button @click="clearFilters">清除筛选</el-button>
        <el-button type="primary" @click="loadAlerts">刷新</el-button>
      </div>
    </el-card>

    <!-- 告警列表 -->
    <el-card class="alert-list-card">
      <template #header>
        <div class="card-header">
          <span>告警列表 ({{ filteredAlerts.length }})</span>
        </div>
      </template>

      <div v-loading="loading" class="alert-list">
        <div
          v-for="alert in filteredAlerts"
          :key="alert.id"
          class="alert-item"
          :class="`alert-${alert.level}`"
        >
          <div class="alert-header">
            <div class="alert-title-section">
              <el-tag :type="getLevelType(alert.level)" size="small">
                {{ getLevelText(alert.level) }}
              </el-tag>
              <el-tag type="info" size="small">
                {{ getTypeText(alert.type) }}
              </el-tag>
              <h4 class="alert-title">{{ alert.title }}</h4>
            </div>
            <el-tag :type="getStatusType(alert.status)" size="small">
              {{ getStatusText(alert.status) }}
            </el-tag>
          </div>

          <p class="alert-description">{{ alert.description }}</p>

          <div class="alert-meta">
            <span class="meta-item">来源: {{ alert.source }}</span>
            <span class="meta-item">创建时间: {{ alert.createdAt }}</span>
            <span v-if="alert.acknowledgedAt" class="meta-item">
              确认时间: {{ alert.acknowledgedAt }}
            </span>
            <span v-if="alert.resolvedAt" class="meta-item">
              解决时间: {{ alert.resolvedAt }}
            </span>
          </div>

          <div class="alert-actions">
            <el-button
              v-if="alert.status === 'active'"
              size="small"
              type="warning"
              @click="handleAcknowledge(alert)"
            >
              确认
            </el-button>
            <el-button
              v-if="alert.status !== 'resolved'"
              size="small"
              type="success"
              @click="handleResolve(alert)"
            >
              解决
            </el-button>
            <el-button
              size="small"
              type="danger"
              :icon="Delete"
              @click="handleDelete(alert)"
            >
              删除
            </el-button>
          </div>
        </div>

        <el-empty
          v-if="!loading && filteredAlerts.length === 0"
          description="暂无告警信息"
        />
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.alert-center {
  padding: 24px;
}

.stats-container {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  cursor: pointer;
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.filter-card {
  margin-bottom: 24px;
}

.filter-container {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.alert-list-card {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
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

.alert-item:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.alert-critical {
  border-left-color: #f56c6c;
  background: #fef0f0;
}

.alert-warning {
  border-left-color: #e6a23c;
  background: #fdf6ec;
}

.alert-info {
  border-left-color: #409eff;
  background: #ecf5ff;
}

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
  flex: 1;
}

.alert-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.alert-description {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}

.alert-meta {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  margin-bottom: 12px;
  padding-top: 12px;
  border-top: 1px solid #e4e7ed;
}

.meta-item {
  font-size: 13px;
  color: #909399;
}

.alert-actions {
  display: flex;
  gap: 8px;
}

@media (max-width: 1200px) {
  .stats-container {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-container {
    grid-template-columns: 1fr;
  }

  .filter-container {
    flex-direction: column;
  }

  .filter-container > * {
    width: 100% !important;
  }
}
</style>
