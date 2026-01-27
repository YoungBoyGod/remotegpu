<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useRoleNavigation } from '@/composables/useRoleNavigation'
import {
  ArrowLeft,
  User,
  Phone,
  Message,
  TrendCharts,
  FolderOpened,
  Picture,
  Warning
} from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'

interface Customer {
  id: number
  name: string
  company: string
  contact: string
  email: string
  phone: string
  status: 'active' | 'inactive' | 'maintenance'
  createdAt: string
  resources: {
    cpuQuota: number
    memoryQuota: number
    gpuQuota: number
    storageQuota: number
    cpuUsed: number
    memoryUsed: number
    gpuUsed: number
    storageUsed: number
  }
  issueCount: number
}

interface TrainingTask {
  id: number
  name: string
  status: 'running' | 'completed' | 'failed' | 'pending'
  gpuUsed: number
  startTime: string
  duration: string
}

interface Dataset {
  id: number
  name: string
  size: number
  createdAt: string
}

interface Image {
  id: number
  name: string
  size: number
  createdAt: string
}

interface Issue {
  id: number
  title: string
  status: 'pending' | 'processing' | 'resolved' | 'closed'
  priority: 'low' | 'medium' | 'high' | 'urgent'
  createdAt: string
}

interface Machine {
  id: number
  hostname: string
  ip: string
  gpuModel: string
  gpuCount: number
  cpuCores: number
  memory: number
  status: 'running' | 'stopped' | 'maintenance'
  startTime: string
}

interface OperationLog {
  id: number
  type: 'web_login' | 'ssh_login' | 'web_access'
  sourceIp: string
  targetIp?: string
  description: string
  timestamp: string
}

const route = useRoute()
const router = useRouter()
const { navigateTo } = useRoleNavigation()
const loading = ref(false)
const customer = ref<Customer | null>(null)
const trainingTasks = ref<TrainingTask[]>([])
const datasets = ref<Dataset[]>([])
const images = ref<Image[]>([])
const issues = ref<Issue[]>([])
const machines = ref<Machine[]>([])
const operationLogs = ref<OperationLog[]>([])

const customerId = computed(() => Number(route.params.id))

// 计算资源使用百分比
const resourceUsagePercent = computed(() => {
  if (!customer.value) return { cpu: 0, memory: 0, gpu: 0, storage: 0 }
  const res = customer.value.resources
  return {
    cpu: res.cpuQuota > 0 ? Math.round((res.cpuUsed / res.cpuQuota) * 100) : 0,
    memory: res.memoryQuota > 0 ? Math.round((res.memoryUsed / res.memoryQuota) * 100) : 0,
    gpu: res.gpuQuota > 0 ? Math.round((res.gpuUsed / res.gpuQuota) * 100) : 0,
    storage: res.storageQuota > 0 ? Math.round((res.storageUsed / res.storageQuota) * 100) : 0
  }
})

// 获取进度条颜色
const getProgressColor = (percent: number) => {
  if (percent >= 90) return '#F56C6C'
  if (percent >= 70) return '#E6A23C'
  return '#67C23A'
}

// 获取状态类型
const getStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    active: 'success',
    inactive: 'info',
    maintenance: 'warning',
    running: 'primary',
    completed: 'success',
    failed: 'danger',
    pending: 'info',
    processing: 'warning',
    resolved: 'success',
    closed: 'info'
  }
  return statusMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '活跃',
    inactive: '停用',
    maintenance: '维护中',
    running: '运行中',
    completed: '已完成',
    failed: '失败',
    pending: '待处理',
    processing: '处理中',
    resolved: '已解决',
    closed: '已关闭'
  }
  return statusMap[status] || status
}

// 加载客户信息
const loadCustomer = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 500))

    // Mock数据
    customer.value = {
      id: customerId.value,
      name: '张三',
      company: '科技有限公司',
      contact: '张三',
      email: 'zhangsan@example.com',
      phone: '13800138000',
      status: 'active',
      createdAt: '2026-01-15',
      resources: {
        cpuQuota: 100,
        memoryQuota: 256,
        gpuQuota: 8,
        storageQuota: 1000,
        cpuUsed: 45,
        memoryUsed: 128,
        gpuUsed: 3,
        storageUsed: 450
      },
      issueCount: 2
    }
  } catch (error) {
    ElMessage.error('加载客户信息失败')
  } finally {
    loading.value = false
  }
}

// 加载训练任务
const loadTrainingTasks = async () => {
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 300))
    trainingTasks.value = [
      {
        id: 1,
        name: 'ResNet训练任务',
        status: 'running',
        gpuUsed: 2,
        startTime: '2026-01-27 10:00:00',
        duration: '2小时30分'
      },
      {
        id: 2,
        name: 'BERT模型训练',
        status: 'completed',
        gpuUsed: 1,
        startTime: '2026-01-26 14:00:00',
        duration: '5小时15分'
      },
      {
        id: 3,
        name: 'YOLO目标检测',
        status: 'pending',
        gpuUsed: 0,
        startTime: '2026-01-27 15:00:00',
        duration: '-'
      }
    ]
  } catch (error) {
    ElMessage.error('加载训练任务失败')
  }
}

// 加载数据集
const loadDatasets = async () => {
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 300))
    datasets.value = [
      { id: 1, name: 'ImageNet数据集', size: 150.5, createdAt: '2026-01-20' },
      { id: 2, name: 'COCO数据集', size: 89.3, createdAt: '2026-01-18' },
      { id: 3, name: '自定义数据集', size: 25.8, createdAt: '2026-01-15' }
    ]
  } catch (error) {
    ElMessage.error('加载数据集失败')
  }
}

// 加载镜像
const loadImages = async () => {
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 300))
    images.value = [
      { id: 1, name: 'pytorch:2.0-cuda11.8', size: 5.2, createdAt: '2026-01-22' },
      { id: 2, name: 'tensorflow:2.13-gpu', size: 4.8, createdAt: '2026-01-20' },
      { id: 3, name: 'custom-ml-env:v1.0', size: 3.5, createdAt: '2026-01-18' }
    ]
  } catch (error) {
    ElMessage.error('加载镜像失败')
  }
}

// 加载问题列表
const loadIssues = async () => {
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 300))
    issues.value = [
      {
        id: 1,
        title: 'GPU资源不足',
        status: 'pending',
        priority: 'high',
        createdAt: '2026-01-27 10:00:00'
      },
      {
        id: 2,
        title: '网络连接问题',
        status: 'processing',
        priority: 'urgent',
        createdAt: '2026-01-26 15:30:00'
      }
    ]
  } catch (error) {
    ElMessage.error('加载问题列表失败')
  }
}

// 加载机器列表
const loadMachines = async () => {
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 300))
    machines.value = [
      {
        id: 1,
        hostname: 'gpu-server-01',
        ip: '192.168.1.101',
        gpuModel: 'RTX 4090',
        gpuCount: 2,
        cpuCores: 32,
        memory: 128,
        status: 'running',
        startTime: '2026-01-20 09:00:00'
      },
      {
        id: 2,
        hostname: 'gpu-server-02',
        ip: '192.168.1.102',
        gpuModel: 'A800',
        gpuCount: 4,
        cpuCores: 64,
        memory: 256,
        status: 'running',
        startTime: '2026-01-22 14:30:00'
      },
      {
        id: 3,
        hostname: 'gpu-server-03',
        ip: '192.168.1.103',
        gpuModel: 'H800',
        gpuCount: 8,
        cpuCores: 128,
        memory: 512,
        status: 'maintenance',
        startTime: '2026-01-15 10:00:00'
      }
    ]
  } catch (error) {
    ElMessage.error('加载机器列表失败')
  }
}

// 加载操作记录
const loadOperationLogs = async () => {
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 300))
    operationLogs.value = [
      {
        id: 1,
        type: 'web_login',
        sourceIp: '120.78.45.123',
        description: '使用Web方式登录系统',
        timestamp: '2026-01-27 10:30:15'
      },
      {
        id: 2,
        type: 'ssh_login',
        sourceIp: '120.78.45.123',
        targetIp: '192.168.1.101',
        description: '使用SSH方式登录了 192.168.1.101',
        timestamp: '2026-01-27 10:32:45'
      },
      {
        id: 3,
        type: 'web_access',
        sourceIp: '120.78.45.123',
        description: '访问了训练任务管理页面',
        timestamp: '2026-01-27 10:35:20'
      },
      {
        id: 4,
        type: 'ssh_login',
        sourceIp: '120.78.45.123',
        targetIp: '192.168.1.102',
        description: '使用SSH方式登录了 192.168.1.102',
        timestamp: '2026-01-27 11:15:30'
      },
      {
        id: 5,
        type: 'web_login',
        sourceIp: '120.78.45.123',
        description: '使用Web方式登录系统',
        timestamp: '2026-01-27 14:20:10'
      }
    ]
  } catch (error) {
    ElMessage.error('加载操作记录失败')
  }
}

// 返回客户列表
const handleBack = () => {
  navigateTo('/customer-center')
}

// 格式化文件大小
const formatSize = (size: number) => {
  if (size < 1024) return `${size.toFixed(2)} GB`
  return `${(size / 1024).toFixed(2)} TB`
}

onMounted(() => {
  loadCustomer()
  loadTrainingTasks()
  loadDatasets()
  loadImages()
  loadIssues()
  loadMachines()
  loadOperationLogs()
})
</script>

<template>
  <div v-loading="loading" class="customer-dashboard">
    <!-- 页面头部 -->
    <div class="page-header">
      <el-button :icon="ArrowLeft" @click="handleBack">返回</el-button>
      <h2 class="page-title">客户详情 - {{ customer?.name }}</h2>
    </div>

    <div v-if="customer" class="dashboard-content">
      <!-- 客户基本信息 -->
      <el-card class="info-card">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
            <el-tag :type="getStatusType(customer.status)">
              {{ getStatusText(customer.status) }}
            </el-tag>
          </div>
        </template>
        <div class="info-grid">
          <div class="info-item">
            <el-icon class="info-icon"><User /></el-icon>
            <div class="info-content">
              <div class="info-label">客户姓名</div>
              <div class="info-value">{{ customer.name }}</div>
            </div>
          </div>
          <div class="info-item">
            <el-icon class="info-icon"><User /></el-icon>
            <div class="info-content">
              <div class="info-label">公司名称</div>
              <div class="info-value">{{ customer.company }}</div>
            </div>
          </div>
          <div class="info-item">
            <el-icon class="info-icon"><User /></el-icon>
            <div class="info-content">
              <div class="info-label">联系人</div>
              <div class="info-value">{{ customer.contact }}</div>
            </div>
          </div>
          <div class="info-item">
            <el-icon class="info-icon"><Phone /></el-icon>
            <div class="info-content">
              <div class="info-label">电话</div>
              <div class="info-value">{{ customer.phone }}</div>
            </div>
          </div>
          <div class="info-item">
            <el-icon class="info-icon"><Message /></el-icon>
            <div class="info-content">
              <div class="info-label">邮箱</div>
              <div class="info-value">{{ customer.email }}</div>
            </div>
          </div>
          <div class="info-item">
            <el-icon class="info-icon"><User /></el-icon>
            <div class="info-content">
              <div class="info-label">创建时间</div>
              <div class="info-value">{{ customer.createdAt }}</div>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 使用的机器列表 -->
      <el-card class="machine-card">
        <template #header>
          <div class="card-header">
            <span>使用的机器 ({{ machines.length }})</span>
          </div>
        </template>
        <el-table :data="machines" stripe>
          <el-table-column prop="hostname" label="主机名" width="150" />
          <el-table-column prop="ip" label="IP地址" width="140" />
          <el-table-column prop="gpuModel" label="GPU型号" width="120" />
          <el-table-column prop="gpuCount" label="GPU数量" width="100" />
          <el-table-column prop="cpuCores" label="CPU核心" width="100" />
          <el-table-column label="内存" width="100">
            <template #default="{ row }">
              {{ row.memory }} GB
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="startTime" label="启动时间" width="180" />
        </el-table>
      </el-card>

      <!-- 操作记录 -->
      <el-card class="operation-card">
        <template #header>
          <div class="card-header">
            <span>操作记录 ({{ operationLogs.length }})</span>
          </div>
        </template>
        <el-table :data="operationLogs" stripe>
          <el-table-column label="操作类型" width="120">
            <template #default="{ row }">
              <el-tag v-if="row.type === 'web_login'" type="success" size="small">Web登录</el-tag>
              <el-tag v-else-if="row.type === 'ssh_login'" type="warning" size="small">SSH登录</el-tag>
              <el-tag v-else type="info" size="small">Web访问</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="sourceIp" label="来源IP" width="140" />
          <el-table-column prop="targetIp" label="目标IP" width="140">
            <template #default="{ row }">
              {{ row.targetIp || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="description" label="操作描述" min-width="300" />
          <el-table-column prop="timestamp" label="操作时间" width="180" />
        </el-table>
      </el-card>

      <!-- 训练任务 -->
      <el-card class="task-card">
        <template #header>
          <div class="card-header">
            <span><el-icon><TrendCharts /></el-icon> 训练任务 ({{ trainingTasks.length }})</span>
          </div>
        </template>
        <el-table :data="trainingTasks" stripe>
          <el-table-column prop="name" label="任务名称" min-width="200" />
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="gpuUsed" label="GPU使用" width="100" />
          <el-table-column prop="startTime" label="开始时间" width="180" />
          <el-table-column prop="duration" label="运行时长" width="120" />
        </el-table>
      </el-card>

      <!-- 数据集和镜像 -->
      <div class="two-column-layout">
        <el-card class="dataset-card">
          <template #header>
            <div class="card-header">
              <span><el-icon><FolderOpened /></el-icon> 数据集 ({{ datasets.length }})</span>
            </div>
          </template>
          <el-table :data="datasets" stripe>
            <el-table-column prop="name" label="数据集名称" min-width="150" />
            <el-table-column label="大小" width="100">
              <template #default="{ row }">
                {{ formatSize(row.size) }}
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="创建时间" width="120" />
          </el-table>
        </el-card>

        <el-card class="image-card">
          <template #header>
            <div class="card-header">
              <span><el-icon><Picture /></el-icon> 镜像 ({{ images.length }})</span>
            </div>
          </template>
          <el-table :data="images" stripe>
            <el-table-column prop="name" label="镜像名称" min-width="150" />
            <el-table-column label="大小" width="100">
              <template #default="{ row }">
                {{ formatSize(row.size) }}
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="创建时间" width="120" />
          </el-table>
        </el-card>
      </div>

      <!-- 问题列表 -->
      <el-card class="issue-card">
        <template #header>
          <div class="card-header">
            <span><el-icon><Warning /></el-icon> 问题列表 ({{ issues.length }})</span>
          </div>
        </template>
        <el-table :data="issues" stripe>
          <el-table-column prop="title" label="问题标题" min-width="200" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="priority" label="优先级" width="100" />
          <el-table-column prop="createdAt" label="创建时间" width="180" />
        </el-table>
      </el-card>
    </div>
  </div>
</template>

<style scoped>
.customer-dashboard {
  padding: 24px;
  min-height: 100vh;
  background: #f5f7fa;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.page-title {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.dashboard-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 16px;
}

.card-header span {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 客户信息卡片 */
.info-card {
  margin-bottom: 0;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.info-icon {
  font-size: 24px;
  color: #409EFF;
}

.info-content {
  flex: 1;
}

.info-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 4px;
}

.info-value {
  font-size: 15px;
  color: #303133;
  font-weight: 500;
}

/* 机器列表和操作记录卡片 */
.machine-card,
.operation-card {
  margin-bottom: 0;
}

/* 两列布局 */
.two-column-layout {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .info-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .info-grid {
    grid-template-columns: 1fr;
  }

  .two-column-layout {
    grid-template-columns: 1fr;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>


