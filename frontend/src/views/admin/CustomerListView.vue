<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCustomerList, disableCustomer, enableCustomer } from '@/api/admin'
import type { Customer } from '@/types/customer'
import type { PageRequest } from '@/types/common'
import DataTable from '@/components/common/DataTable.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const customers = ref<Customer[]>([])
const total = ref(0)
const pageRequest = ref<PageRequest>({
  page: 1,
  pageSize: 10,
  filters: {}
})

const filters = ref({
  status: '',
  keyword: ''
})

const loadCustomers = async () => {
  try {
    loading.value = true
    const response = await getCustomerList({
      ...pageRequest.value,
      filters: filters.value
    })
    customers.value = response.data.list
    total.value = response.data.total
  } catch (error) {
    console.error('加载客户列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  pageRequest.value.page = page
  loadCustomers()
}

const handleSizeChange = (size: number) => {
  pageRequest.value.pageSize = size
  loadCustomers()
}

const handleSearch = () => {
  pageRequest.value.page = 1
  loadCustomers()
}

const handleReset = () => {
  filters.value = {
    status: '',
    keyword: ''
  }
  handleSearch()
}

const handleViewDetail = (customer: Customer) => {
  router.push(`/admin/customers/${customer.id}`)
}

const handleDisable = async (customer: Customer) => {
  try {
    await ElMessageBox.confirm(
      `确定要禁用客户 "${getCustomerName(customer)}" 吗?`,
      '禁用确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await disableCustomer(customer.id)
    ElMessage.success('已禁用')
    loadCustomers()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.msg || '禁用失败')
    }
  }
}

const handleEnable = async (customer: Customer) => {
  try {
    await ElMessageBox.confirm(
      `确定要启用客户 "${getCustomerName(customer)}" 吗?`,
      '启用确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }
    )
    await enableCustomer(customer.id)
    ElMessage.success('已启用')
    loadCustomers()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.msg || '启用失败')
    }
  }
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    active: 'success',
    pending: 'info',
    disabled: 'warning',
    suspended: 'danger',
    deleted: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '正常',
    pending: '未激活',
    disabled: '已禁用',
    suspended: '已停用',
    deleted: '已删除'
  }
  return statusMap[status] || status
}

const formatDateTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

const maskPhone = (phone?: string | null) => {
  if (!phone) return '-'
  if (phone.length >= 7) {
    return phone.slice(0, 3) + '****' + phone.slice(-4)
  }
  return phone
}

const getCustomerName = (customer: Customer) => {
  return customer.company_code || customer.company || customer.display_name || customer.username || customer.name || '-'
}

const getContactName = (customer: Customer) => {
  return customer.full_name || customer.display_name || customer.contactPerson || customer.username || '-'
}

const getRowIndex = (index: number) => {
  return (pageRequest.value.page - 1) * pageRequest.value.pageSize + index + 1
}

onMounted(() => {
  loadCustomers()
})
</script>

<template>
  <div class="customer-list">
    <div class="page-header">
      <h2 class="page-title">客户列表</h2>
      <el-button type="primary" @click="router.push('/admin/customers/add')">添加客户</el-button>
    </div>

    <!-- 筛选栏 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters">
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="正常" value="active" />
            <el-option label="已禁用" value="disabled" />
            <el-option label="已删除" value="deleted" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="filters.keyword" placeholder="公司代号/用户名" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <DataTable
      :data="customers"
      :total="total"
      :loading="loading"
      :current-page="pageRequest.page"
      :page-size="pageRequest.pageSize"
      :show-pagination="true"
      @page-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <el-table-column type="index" label="序号" width="70" :index="getRowIndex" />
      <el-table-column label="公司代号" min-width="160">
        <template #default="{ row }">
          {{ row.company_code || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="联系人" width="140">
        <template #default="{ row }">
          {{ getContactName(row) }}
        </template>
      </el-table-column>
      <el-table-column prop="email" label="联系邮箱" min-width="200" />
      <el-table-column label="联系电话" width="140">
        <template #default="{ row }">
          {{ maskPhone(row.phone) }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDateTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="240" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="handleViewDetail(row)">
            详情
          </el-button>
          <el-button v-if="row.status === 'active'" link type="warning" size="small" @click="handleDisable(row)">
            禁用
          </el-button>
          <el-button v-if="row.status === 'disabled'" link type="success" size="small" @click="handleEnable(row)">
            启用
          </el-button>
        </template>
      </el-table-column>
    </DataTable>
  </div>
</template>

<style scoped>
.customer-list {
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
  margin: 0;
}

.filter-card {
  margin-bottom: 20px;
}
</style>
