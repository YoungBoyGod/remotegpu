<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Edit,
  Delete,
  Setting,
  ChatDotRound,
  Search,
  Refresh,
  View
} from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { useRouter } from 'vue-router'

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

interface Issue {
  id: number
  customerId: number
  customerName: string
  title: string
  description: string
  status: 'pending' | 'processing' | 'resolved' | 'closed'
  priority: 'low' | 'medium' | 'high' | 'urgent'
  createdAt: string
  updatedAt: string
  feedback?: string
}

const router = useRouter()
const loading = ref(false)
const customers = ref<Customer[]>([])
const issues = ref<Issue[]>([])
const searchKeyword = ref('')
const selectedStatus = ref<string>('all')

// 对话框状态
const customerDialogVisible = ref(false)
const resourceDialogVisible = ref(false)
const issueDialogVisible = ref(false)
const isEditMode = ref(false)
const currentCustomer = ref<Customer | null>(null)

// 客户表单
const customerForm = ref({
  name: '',
  company: '',
  contact: '',
  email: '',
  phone: '',
  status: 'active' as 'active' | 'inactive' | 'maintenance'
})

// 资源配置表单
const resourceForm = ref({
  cpuQuota: 0,
  memoryQuota: 0,
  gpuQuota: 0,
  storageQuota: 0
})

// 问题反馈表单
const issueForm = ref({
  feedback: ''
})

// 当前选中的问题
const currentIssue = ref<Issue | null>(null)

// 过滤后的客户列表
const filteredCustomers = computed(() => {
  let result = customers.value

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(c =>
      c.name.toLowerCase().includes(keyword) ||
      c.company.toLowerCase().includes(keyword) ||
      c.contact.toLowerCase().includes(keyword) ||
      c.email.toLowerCase().includes(keyword)
    )
  }

  // 状态筛选
  if (selectedStatus.value !== 'all') {
    result = result.filter(c => c.status === selectedStatus.value)
  }

  return result
})

// 客户统计
const customerStats = computed(() => {
  return {
    total: customers.value.length,
    active: customers.value.filter(c => c.status === 'active').length,
    inactive: customers.value.filter(c => c.status === 'inactive').length,
    maintenance: customers.value.filter(c => c.status === 'maintenance').length,
    totalIssues: issues.value.length,
    pendingIssues: issues.value.filter(i => i.status === 'pending').length
  }
})

// 加载客户列表
const loadCustomers = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 500))
    customers.value = [
      {
        id: 1,
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
      },
      {
        id: 2,
        name: '李四',
        company: '数据科技公司',
        contact: '李四',
        email: 'lisi@example.com',
        phone: '13900139000',
        status: 'active',
        createdAt: '2026-01-10',
        resources: {
          cpuQuota: 200,
          memoryQuota: 512,
          gpuQuota: 16,
          storageQuota: 2000,
          cpuUsed: 120,
          memoryUsed: 300,
          gpuUsed: 8,
          storageUsed: 1200
        },
        issueCount: 0
      }
    ]
  } catch (error) {
    ElMessage.error('加载客户列表失败')
  } finally {
    loading.value = false
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
        customerId: 1,
        customerName: '张三',
        title: 'GPU资源不足',
        description: '当前分配的GPU资源无法满足训练需求',
        status: 'pending',
        priority: 'high',
        createdAt: '2026-01-27 10:00:00',
        updatedAt: '2026-01-27 10:00:00'
      },
      {
        id: 2,
        customerId: 1,
        customerName: '张三',
        title: '网络连接问题',
        description: '无法访问训练环境',
        status: 'processing',
        priority: 'urgent',
        createdAt: '2026-01-26 15:30:00',
        updatedAt: '2026-01-27 09:00:00'
      }
    ]
  } catch (error) {
    ElMessage.error('加载问题列表失败')
  }
}

// 获取状态类型
const getStatusType = (status: string) => {
  const statusMap = {
    active: 'success',
    inactive: 'info',
    maintenance: 'warning'
  }
  return statusMap[status as keyof typeof statusMap] || 'info'
}

// 获取状态文本
const getStatusText = (status: string) => {
  const statusMap = {
    active: '活跃',
    inactive: '停用',
    maintenance: '维护中'
  }
  return statusMap[status as keyof typeof statusMap] || '未知'
}

// 查看客户详情
const handleViewCustomer = (customer: Customer) => {
  router.push(`/customer-center/${customer.id}`)
}

// 打开新增客户对话框
const handleAddCustomer = () => {
  isEditMode.value = false
  customerForm.value = {
    name: '',
    company: '',
    contact: '',
    email: '',
    phone: '',
    status: 'active'
  }
  customerDialogVisible.value = true
}

// 打开编辑客户对话框
const handleEditCustomer = (customer: Customer) => {
  isEditMode.value = true
  currentCustomer.value = customer
  customerForm.value = {
    name: customer.name,
    company: customer.company,
    contact: customer.contact,
    email: customer.email,
    phone: customer.phone,
    status: customer.status
  }
  customerDialogVisible.value = true
}

// 保存客户
const handleSaveCustomer = async () => {
  try {
    // TODO: 调用API保存数据
    await new Promise(resolve => setTimeout(resolve, 300))

    if (isEditMode.value && currentCustomer.value) {
      // 更新客户
      Object.assign(currentCustomer.value, customerForm.value)
      ElMessage.success('客户信息已更新')
    } else {
      // 新增客户
      const newCustomer: Customer = {
        id: customers.value.length + 1,
        ...customerForm.value,
        createdAt: new Date().toLocaleDateString('zh-CN'),
        resources: {
          cpuQuota: 0,
          memoryQuota: 0,
          gpuQuota: 0,
          storageQuota: 0,
          cpuUsed: 0,
          memoryUsed: 0,
          gpuUsed: 0,
          storageUsed: 0
        },
        issueCount: 0
      }
      customers.value.push(newCustomer)
      ElMessage.success('客户已添加')
    }

    customerDialogVisible.value = false
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

// 删除客户
const handleDeleteCustomer = async (customer: Customer) => {
  try {
    await ElMessageBox.confirm('确认删除此客户？', '删除客户', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'error'
    })

    // TODO: 调用API删除
    const index = customers.value.findIndex(c => c.id === customer.id)
    if (index > -1) {
      customers.value.splice(index, 1)
    }
    ElMessage.success('客户已删除')
  } catch {
    // 用户取消
  }
}

// 打开资源配置对话框
const handleConfigResource = (customer: Customer) => {
  currentCustomer.value = customer
  resourceForm.value = {
    cpuQuota: customer.resources.cpuQuota,
    memoryQuota: customer.resources.memoryQuota,
    gpuQuota: customer.resources.gpuQuota,
    storageQuota: customer.resources.storageQuota
  }
  resourceDialogVisible.value = true
}

// 保存资源配置
const handleSaveResource = async () => {
  try {
    // TODO: 调用API保存数据
    await new Promise(resolve => setTimeout(resolve, 300))

    if (currentCustomer.value) {
      currentCustomer.value.resources.cpuQuota = resourceForm.value.cpuQuota
      currentCustomer.value.resources.memoryQuota = resourceForm.value.memoryQuota
      currentCustomer.value.resources.gpuQuota = resourceForm.value.gpuQuota
      currentCustomer.value.resources.storageQuota = resourceForm.value.storageQuota
    }

    ElMessage.success('资源配置已更新')
    resourceDialogVisible.value = false
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

// 获取问题状态类型
const getIssueStatusType = (status: string) => {
  const statusMap = {
    pending: 'danger',
    processing: 'warning',
    resolved: 'success',
    closed: 'info'
  }
  return statusMap[status as keyof typeof statusMap] || 'info'
}

// 获取问题优先级类型
const getPriorityType = (priority: string) => {
  const priorityMap = {
    low: 'info',
    medium: 'warning',
    high: 'danger',
    urgent: 'danger'
  }
  return priorityMap[priority as keyof typeof priorityMap] || 'info'
}

// 打开问题反馈对话框
const handleFeedbackIssue = (issue: Issue) => {
  currentIssue.value = issue
  issueForm.value.feedback = issue.feedback || ''
  issueDialogVisible.value = true
}

// 提交问题反馈
const handleSubmitFeedback = async () => {
  try {
    // TODO: 调用API提交反馈
    await new Promise(resolve => setTimeout(resolve, 300))

    if (currentIssue.value) {
      currentIssue.value.feedback = issueForm.value.feedback
      currentIssue.value.status = 'resolved'
      currentIssue.value.updatedAt = new Date().toLocaleString('zh-CN')
    }

    ElMessage.success('反馈已提交')
    issueDialogVisible.value = false
  } catch (error) {
    ElMessage.error('提交失败')
  }
}

onMounted(() => {
  loadCustomers()
  loadIssues()
})
</script>

<template>
  <div class="customer-center">
    <PageHeader title="客户中心" />

    <!-- 统计卡片 -->
    <div class="stats-container">
      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #409EFF">
            <el-icon :size="32"><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">客户总数</div>
            <div class="stat-value">{{ customerStats.total }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #67C23A">
            <el-icon :size="32"><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">活跃客户</div>
            <div class="stat-value">{{ customerStats.active }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #E6A23C">
            <el-icon :size="32"><Warning /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">维护中</div>
            <div class="stat-value">{{ customerStats.maintenance }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #F56C6C">
            <el-icon :size="32"><ChatDotRound /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">待处理问题</div>
            <div class="stat-value">{{ customerStats.pendingIssues }}</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 客户列表 -->
    <el-card class="customer-list-card">
      <template #header>
        <div class="card-header">
          <span>客户列表</span>
          <div class="header-actions">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索客户"
              :prefix-icon="Search"
              clearable
              style="width: 250px; margin-right: 12px"
            />
            <el-select
              v-model="selectedStatus"
              placeholder="状态"
              clearable
              style="width: 120px; margin-right: 12px"
            >
              <el-option label="全部" value="all" />
              <el-option label="活跃" value="active" />
              <el-option label="停用" value="inactive" />
              <el-option label="维护中" value="maintenance" />
            </el-select>
            <el-button :icon="Refresh" @click="loadCustomers">刷新</el-button>
            <el-button type="primary" :icon="Plus" @click="handleAddCustomer">
              新增客户
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="filteredCustomers" :loading="loading" stripe>
        <el-table-column prop="name" label="客户姓名" width="120" />
        <el-table-column prop="company" label="公司名称" width="180" />
        <el-table-column prop="contact" label="联系人" width="100" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="phone" label="电话" width="130" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="资源使用" width="200">
          <template #default="{ row }">
            <div class="resource-usage">
              <div class="usage-item">
                GPU: {{ row.resources.gpuUsed }}/{{ row.resources.gpuQuota }}
              </div>
              <div class="usage-item">
                CPU: {{ row.resources.cpuUsed }}/{{ row.resources.cpuQuota }}
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="issueCount" label="问题数" width="80" />
        <el-table-column prop="createdAt" label="创建时间" width="120" />
        <el-table-column label="操作" width="360" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" :icon="View" @click="handleViewCustomer(row)">
              查看
            </el-button>
            <el-button size="small" :icon="Edit" @click="handleEditCustomer(row)">
              编辑
            </el-button>
            <el-button size="small" :icon="Setting" @click="handleConfigResource(row)">
              配置资源
            </el-button>
            <el-button
              size="small"
              type="danger"
              :icon="Delete"
              @click="handleDeleteCustomer(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 客户问题列表 -->
    <el-card class="issue-list-card">
      <template #header>
        <div class="card-header">
          <span>客户问题 ({{ issues.length }})</span>
        </div>
      </template>

      <el-table :data="issues" stripe>
        <el-table-column prop="customerName" label="客户" width="120" />
        <el-table-column prop="title" label="问题标题" width="200" />
        <el-table-column prop="description" label="问题描述" min-width="250" />
        <el-table-column label="优先级" width="100">
          <template #default="{ row }">
            <el-tag :type="getPriorityType(row.priority)" size="small">
              {{ row.priority }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getIssueStatusType(row.status)" size="small">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="160" />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button
              size="small"
              type="primary"
              :icon="ChatDotRound"
              @click="handleFeedbackIssue(row)"
            >
              反馈
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 客户编辑对话框 -->
    <el-dialog
      v-model="customerDialogVisible"
      :title="isEditMode ? '编辑客户' : '新增客户'"
      width="600px"
    >
      <el-form :model="customerForm" label-width="100px">
        <el-form-item label="客户姓名">
          <el-input v-model="customerForm.name" placeholder="请输入客户姓名" />
        </el-form-item>
        <el-form-item label="公司名称">
          <el-input v-model="customerForm.company" placeholder="请输入公司名称" />
        </el-form-item>
        <el-form-item label="联系人">
          <el-input v-model="customerForm.contact" placeholder="请输入联系人" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="customerForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="customerForm.phone" placeholder="请输入电话" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="customerForm.status" placeholder="请选择状态">
            <el-option label="活跃" value="active" />
            <el-option label="停用" value="inactive" />
            <el-option label="维护中" value="maintenance" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="customerDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveCustomer">保存</el-button>
      </template>
    </el-dialog>

    <!-- 资源配置对话框 -->
    <el-dialog v-model="resourceDialogVisible" title="配置资源" width="600px">
      <el-form :model="resourceForm" label-width="120px">
        <el-form-item label="CPU配额(核)">
          <el-input-number v-model="resourceForm.cpuQuota" :min="0" :step="1" />
        </el-form-item>
        <el-form-item label="内存配额(GB)">
          <el-input-number v-model="resourceForm.memoryQuota" :min="0" :step="1" />
        </el-form-item>
        <el-form-item label="GPU配额(卡)">
          <el-input-number v-model="resourceForm.gpuQuota" :min="0" :step="1" />
        </el-form-item>
        <el-form-item label="存储配额(GB)">
          <el-input-number v-model="resourceForm.storageQuota" :min="0" :step="10" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resourceDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveResource">保存</el-button>
      </template>
    </el-dialog>

    <!-- 问题反馈对话框 -->
    <el-dialog v-model="issueDialogVisible" title="问题反馈" width="600px">
      <el-form :model="issueForm" label-width="100px">
        <el-form-item label="问题标题">
          <el-input :value="currentIssue?.title" disabled />
        </el-form-item>
        <el-form-item label="问题描述">
          <el-input :value="currentIssue?.description" type="textarea" :rows="3" disabled />
        </el-form-item>
        <el-form-item label="反馈内容">
          <el-input
            v-model="issueForm.feedback"
            type="textarea"
            :rows="5"
            placeholder="请输入反馈内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="issueDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitFeedback">提交反馈</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.customer-center {
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

.customer-list-card,
.issue-list-card {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.header-actions {
  display: flex;
  align-items: center;
}

.resource-usage {
  font-size: 13px;
  color: #606266;
}

.usage-item {
  margin-bottom: 4px;
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

  .header-actions {
    flex-wrap: wrap;
    gap: 8px;
  }
}
</style>
