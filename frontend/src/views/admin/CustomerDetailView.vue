<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Refresh } from '@element-plus/icons-vue'
import {
  getCustomerDetail,
  updateCustomer,
  disableCustomer,
  enableCustomer,
} from '@/api/admin'
import type { Customer, CustomerAllocation } from '@/types/customer'
import { CopyDocument } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const customer = ref<Customer | null>(null)
const allocations = ref<CustomerAllocation[]>([])

const customerId = () => Number(route.params.id)

const loadCustomerDetail = async () => {
  loading.value = true
  try {
    const res = await getCustomerDetail(customerId())
    customer.value = res.data.customer
    allocations.value = res.data.allocations || []
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

const statusTagType = (status: string) => {
  switch (status) {
    case 'active': return 'success'
    case 'suspended': return 'warning'
    case 'disabled': return 'danger'
    case 'deleted': return 'info'
    default: return 'info'
  }
}

const statusLabel = (status: string) => {
  switch (status) {
    case 'active': return '正常'
    case 'suspended': return '已暂停'
    case 'disabled': return '已停用'
    case 'deleted': return '已删除'
    case 'pending': return '待激活'
    default: return status
  }
}

const accountTypeLabel = (type?: string) => {
  switch (type) {
    case 'individual': return '个人'
    case 'enterprise': return '企业'
    default: return type || '-'
  }
}

const roleLabel = (role?: string) => {
  switch (role) {
    case 'admin': return '管理员'
    case 'customer_owner': return '客户（所有者）'
    case 'customer_member': return '客户（成员）'
    default: return role || '-'
  }
}

// 编辑客户
const editDialogVisible = ref(false)
const editLoading = ref(false)
const editForm = ref({
  display_name: '',
  full_name: '',
  company_code: '',
  company: '',
  email: '',
  phone: '',
})

const openEditDialog = () => {
  if (!customer.value) return
  editForm.value = {
    display_name: customer.value.display_name || '',
    full_name: customer.value.full_name || '',
    company_code: customer.value.company_code || '',
    company: customer.value.company || '',
    email: customer.value.email || '',
    phone: customer.value.phone || '',
  }
  editDialogVisible.value = true
}

const handleEditSubmit = async () => {
  editLoading.value = true
  try {
    await updateCustomer(customerId(), editForm.value)
    ElMessage.success('客户信息已更新')
    editDialogVisible.value = false
    loadCustomerDetail()
  } catch (error) {
    console.error('更新客户信息失败:', error)
  } finally {
    editLoading.value = false
  }
}

// 启用/停用客户
const handleToggleStatus = async () => {
  if (!customer.value) return
  const isActive = customer.value.status === 'active'
  const action = isActive ? '停用' : '启用'
  try {
    await ElMessageBox.confirm(
      `确定${action}客户「${customer.value.username || customer.value.display_name}」吗？`,
      `确认${action}`,
      { type: 'warning' }
    )
    if (isActive) {
      await disableCustomer(customerId())
    } else {
      await enableCustomer(customerId())
    }
    ElMessage.success(`客户已${action}`)
    loadCustomerDetail()
  } catch {
    // 取消
  }
}

const handleBack = () => {
  router.push('/admin/customers/list')
}

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

const goToMachine = (machineId: string) => {
  router.push(`/admin/machines/${machineId}`)
}

onMounted(() => {
  loadCustomerDetail()
})
</script>

<template>
  <div class="customer-detail">
    <PageHeader
      :title="customer ? `客户详情 - ${customer.username || customer.display_name || ''}` : '客户详情'"
      subtitle="查看和管理客户信息"
    >
      <template #actions>
        <el-button :icon="ArrowLeft" @click="handleBack">返回列表</el-button>
        <el-button :icon="Refresh" @click="loadCustomerDetail">刷新</el-button>
      </template>
    </PageHeader>

    <el-skeleton :loading="loading" :rows="10" animated>
      <div v-if="customer">
        <!-- 基本信息 -->
        <el-card class="info-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">基本信息</span>
              <el-button type="primary" link @click="openEditDialog">编辑</el-button>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="用户名">
              {{ customer.username || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="statusTagType(customer.status)">
                {{ statusLabel(customer.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="显示名称">
              {{ customer.display_name || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="联系人">
              {{ customer.full_name || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="公司代号">
              {{ customer.company_code || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="公司名称">
              {{ customer.company || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="邮箱">
              {{ customer.email || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="电话">
              {{ customer.phone || '-' }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- 账户信息 -->
        <el-card class="info-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">账户信息</span>
              <el-button
                :type="customer.status === 'active' ? 'danger' : 'success'"
                link
                @click="handleToggleStatus"
              >
                {{ customer.status === 'active' ? '停用账户' : '启用账户' }}
              </el-button>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="角色">
              {{ roleLabel((customer as any).role) }}
            </el-descriptions-item>
            <el-descriptions-item label="账户类型">
              {{ accountTypeLabel((customer as any).account_type) }}
            </el-descriptions-item>
            <el-descriptions-item label="账户余额">
              {{ (customer as any).balance != null ? `¥ ${Number((customer as any).balance).toFixed(2)}` : '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="货币">
              {{ (customer as any).currency || 'CNY' }}
            </el-descriptions-item>
            <el-descriptions-item label="邮箱验证">
              <el-tag :type="(customer as any).email_verified ? 'success' : 'info'" size="small">
                {{ (customer as any).email_verified ? '已验证' : '未验证' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="手机验证">
              <el-tag :type="(customer as any).phone_verified ? 'success' : 'info'" size="small">
                {{ (customer as any).phone_verified ? '已验证' : '未验证' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ formatDateTime(customer.created_at) }}
            </el-descriptions-item>
            <el-descriptions-item label="最后登录">
              {{ formatDateTime(customer.last_login_at) }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- 分配的机器 -->
        <el-card class="info-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">分配的机器</span>
              <el-tag size="small">{{ allocations.length }} 台</el-tag>
            </div>
          </template>
          <el-table v-if="allocations.length > 0" :data="allocations" border>
            <el-table-column label="机器名称" min-width="140">
              <template #default="{ row }">
                <el-button link type="primary" @click="goToMachine(row.machine_id)">
                  {{ row.machine_name || row.machine_id }}
                </el-button>
              </template>
            </el-table-column>
            <el-table-column label="分配时间" min-width="160">
              <template #default="{ row }">
                {{ formatDateTime(row.allocated_at) }}
              </template>
            </el-table-column>
            <el-table-column label="到期时间" min-width="160">
              <template #default="{ row }">
                {{ formatDateTime(row.end_time) }}
              </template>
            </el-table-column>
            <el-table-column label="SSH" min-width="180">
              <template #default="{ row }">
                <template v-if="row.ssh_host">
                  <span>{{ row.ssh_host }}:{{ row.ssh_port || 22 }}</span>
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(`${row.ssh_host}:${row.ssh_port || 22}`)" />
                </template>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column label="Jupyter" min-width="120">
              <template #default="{ row }">
                <template v-if="row.jupyter_url">
                  <a class="info-link" :href="row.jupyter_url" target="_blank">访问</a>
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(row.jupyter_url)" />
                </template>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column label="VNC" min-width="120">
              <template #default="{ row }">
                <template v-if="row.vnc_url">
                  <a class="info-link" :href="row.vnc_url" target="_blank">访问</a>
                  <el-button link :icon="CopyDocument" @click="copyToClipboard(row.vnc_url)" />
                </template>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
                  {{ row.status === 'active' ? '使用中' : row.status || '-' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
          <div v-else class="empty-tip">暂无分配的机器</div>
        </el-card>

        <!-- 编辑对话框 -->
        <el-dialog v-model="editDialogVisible" title="编辑客户信息" width="560px" :close-on-click-modal="false">
          <el-form :model="editForm" label-width="100px">
            <el-form-item label="显示名称">
              <el-input v-model="editForm.display_name" placeholder="请输入显示名称" />
            </el-form-item>
            <el-form-item label="联系人">
              <el-input v-model="editForm.full_name" placeholder="请输入联系人姓名" />
            </el-form-item>
            <el-form-item label="公司代号">
              <el-input v-model="editForm.company_code" placeholder="请输入公司代号" />
            </el-form-item>
            <el-form-item label="公司名称">
              <el-input v-model="editForm.company" placeholder="请输入公司名称" />
            </el-form-item>
            <el-form-item label="邮箱">
              <el-input v-model="editForm.email" placeholder="请输入邮箱" />
            </el-form-item>
            <el-form-item label="电话">
              <el-input v-model="editForm.phone" placeholder="请输入电话" />
            </el-form-item>
          </el-form>
          <template #footer>
            <el-button @click="editDialogVisible = false">取消</el-button>
            <el-button type="primary" :loading="editLoading" @click="handleEditSubmit">保存</el-button>
          </template>
        </el-dialog>
      </div>
    </el-skeleton>
  </div>
</template>

<style scoped>
.customer-detail {
  padding: 24px;
}

.info-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.info-link {
  color: #409eff;
  text-decoration: none;
}

.info-link:hover {
  text-decoration: underline;
}

.empty-tip {
  padding: 24px;
  text-align: center;
  color: #909399;
  font-size: 14px;
}
</style>
