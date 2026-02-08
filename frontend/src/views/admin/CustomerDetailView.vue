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
  reclaimMachine,
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

const handleReclaim = async (row: CustomerAllocation) => {
  try {
    await ElMessageBox.confirm(
      `确定要回收机器「${row.machine_name || row.machine_id}」吗？回收后客户将无法继续使用。`,
      '回收确认',
      { confirmButtonText: '确定回收', cancelButtonText: '取消', type: 'warning' }
    )
    await reclaimMachine(row.machine_id)
    ElMessage.success('回收成功')
    loadCustomerDetail()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.msg || '回收失败')
    }
  }
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

          <div v-if="allocations.length > 0" class="machine-items">
            <div v-for="row in allocations" :key="row.allocation_id" class="machine-row">
              <!-- 顶栏：名称 + 状态 + 时间 + 回收 -->
              <div class="row-header">
                <el-button link type="primary" class="row-name" @click="goToMachine(row.machine_id)">
                  {{ row.machine_name || row.machine_id }}
                </el-button>
                <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
                  {{ row.status === 'active' ? '使用中' : row.status || '-' }}
                </el-tag>
                <span class="row-time">{{ formatDateTime(row.allocated_at) }} ~ {{ formatDateTime(row.end_time) }}</span>
                <div class="row-actions">
                  <el-button type="danger" size="small" plain @click="handleReclaim(row)">回收</el-button>
                </div>
              </div>

              <!-- 两栏：对内 | 对外 -->
              <div class="row-body">
                <div class="body-col">
                  <div class="col-title">对内连接</div>
                  <div class="conn-line">
                    <span class="conn-tag">SSH</span>
                    <template v-if="row.ip_address || row.ssh_host">
                      <span class="conn-field"><span class="conn-label">IP:</span> {{ row.ip_address || row.ssh_host }}</span>
                      <span class="conn-field"><span class="conn-label">端口:</span> {{ row.ssh_port || 22 }}</span>
                      <span class="conn-field"><span class="conn-label">用户:</span> {{ row.ssh_username || 'root' }}</span>
                      <span class="conn-field" v-if="row.ssh_password">
                        <span class="conn-label">密码:</span>
                        <span class="password-text">{{ row.ssh_password }}</span>
                        <el-button link :icon="CopyDocument" size="small" @click="copyToClipboard(row.ssh_password!)" />
                      </span>
                    </template>
                    <span v-else class="text-muted">未配置</span>
                  </div>
                  <div class="conn-line">
                    <span class="conn-tag conn-tag-jupyter">Jupyter</span>
                    <template v-if="row.jupyter_url">
                      <a class="conn-link" :href="row.jupyter_url" target="_blank">{{ row.jupyter_url }}</a>
                      <el-button link :icon="CopyDocument" size="small" @click="copyToClipboard(row.jupyter_url)" />
                    </template>
                    <span v-else class="text-muted">未配置</span>
                  </div>
                  <div class="conn-line">
                    <span class="conn-tag conn-tag-vnc">VNC</span>
                    <template v-if="row.vnc_url">
                      <a class="conn-link" :href="row.vnc_url" target="_blank">{{ row.vnc_url }}</a>
                      <el-button link :icon="CopyDocument" size="small" @click="copyToClipboard(row.vnc_url)" />
                    </template>
                    <span v-else class="text-muted">未配置</span>
                  </div>
                </div>

                <div class="body-col">
                  <div class="col-title">对外连接</div>
                  <div class="conn-line">
                    <span class="conn-tag conn-tag-ext">SSH</span>
                    <template v-if="row.nginx_domain || row.external_ip">
                      <span class="conn-field">{{ row.nginx_domain || row.external_ip }}:{{ row.external_ssh_port || '-' }}</span>
                      <el-button link :icon="CopyDocument" size="small" @click="copyToClipboard(`ssh -p ${row.external_ssh_port || 22} ${row.ssh_username || 'root'}@${row.nginx_domain || row.external_ip}`)" />
                    </template>
                    <span v-else class="text-muted">未配置</span>
                  </div>
                  <div class="conn-line">
                    <span class="conn-tag conn-tag-ext">Jupyter</span>
                    <template v-if="row.external_jupyter_port && (row.nginx_domain || row.external_ip)">
                      <span class="conn-field">{{ row.nginx_domain || row.external_ip }}:{{ row.external_jupyter_port }}</span>
                      <el-button link :icon="CopyDocument" size="small" @click="copyToClipboard(`${row.nginx_domain || row.external_ip}:${row.external_jupyter_port}`)" />
                    </template>
                    <span v-else class="text-muted">未配置</span>
                  </div>
                  <div class="conn-line">
                    <span class="conn-tag conn-tag-ext">VNC</span>
                    <template v-if="row.external_vnc_port && (row.nginx_domain || row.external_ip)">
                      <span class="conn-field">{{ row.nginx_domain || row.external_ip }}:{{ row.external_vnc_port }}</span>
                      <el-button link :icon="CopyDocument" size="small" @click="copyToClipboard(`${row.nginx_domain || row.external_ip}:${row.external_vnc_port}`)" />
                    </template>
                    <span v-else class="text-muted">未配置</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
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

.machine-row {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #fff;
  margin-bottom: 12px;
  transition: box-shadow 0.2s;
}

.machine-row:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.row-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
}

.row-name {
  font-weight: 600;
  font-size: 15px;
}

.row-time {
  font-size: 12px;
  color: #909399;
}

.row-actions {
  margin-left: auto;
}

.row-body {
  display: flex;
  border-top: 1px solid #ebeef5;
}

.body-col {
  flex: 1;
  padding: 10px 20px;
}

.body-col + .body-col {
  border-left: 1px solid #ebeef5;
}

.col-title {
  font-size: 13px;
  font-weight: 600;
  color: #909399;
  margin-bottom: 8px;
}

.conn-line {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #606266;
  padding: 4px 0;
  border-top: 1px dashed #f0f0f0;
}

.conn-tag {
  display: inline-block;
  min-width: 52px;
  text-align: center;
  font-size: 12px;
  font-weight: 600;
  color: #409eff;
  background: #ecf5ff;
  padding: 2px 8px;
  border-radius: 3px;
  flex-shrink: 0;
}

.conn-tag-jupyter {
  color: #e6a23c;
  background: #fdf6ec;
}

.conn-tag-vnc {
  color: #67c23a;
  background: #f0f9eb;
}

.conn-tag-ext {
  color: #f56c6c;
  background: #fef0f0;
}

.conn-field {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

.conn-label {
  color: #909399;
  flex-shrink: 0;
}

.conn-link {
  color: #409eff;
  text-decoration: none;
  font-size: 13px;
}

.conn-link:hover {
  text-decoration: underline;
}

.password-text {
  font-family: monospace;
  font-size: 13px;
  color: #303133;
}

.text-muted {
  color: #c0c4cc;
}
</style>
