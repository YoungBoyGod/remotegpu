<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCustomerList, deleteCustomer, toggleCustomerStatus } from '@/api/admin'
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

const handleDelete = async (customer: Customer) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除客户 "${customer.name}" 吗?`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await deleteCustomer(customer.id)
    ElMessage.success('删除成功')
    loadCustomers()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除客户失败:', error)
    }
  }
}

const handleToggleStatus = async (customer: Customer) => {
  try {
    const newStatus = customer.status !== 'active'
    await toggleCustomerStatus(customer.id, newStatus)
    ElMessage.success(newStatus ? '已启用' : '已禁用')
    loadCustomers()
  } catch (error) {
    console.error('切换状态失败:', error)
  }
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    active: 'success',
    inactive: 'info',
    suspended: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '正常',
    inactive: '未激活',
    suspended: '已停用'
  }
  return statusMap[status] || status
}

onMounted(() => {
  loadCustomers()
})
</script>

<template>
  <div class="customer-list">
    <div class="page-header">
      <h2 class="page-title">客户列表</h2>
    </div>

    <!-- 筛选栏 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters">
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="正常" value="active" />
            <el-option label="未激活" value="inactive" />
            <el-option label="已停用" value="suspended" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="filters.keyword" placeholder="客户名称/联系人" clearable style="width: 200px" />
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
      <el-table-column prop="name" label="客户名称" min-width="150" />
      <el-table-column prop="contactPerson" label="联系人" width="120" />
      <el-table-column prop="contactEmail" label="联系邮箱" min-width="180" />
      <el-table-column prop="contactPhone" label="联系电话" width="130" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="分配机器数" width="120">
        <template #default="{ row }">
          {{ row.allocatedMachines || 0 }}
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="创建时间" width="180" />
      <el-table-column label="操作" width="240" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="handleViewDetail(row)">
            详情
          </el-button>
          <el-button link type="warning" size="small" @click="handleToggleStatus(row)">
            {{ row.status === 'active' ? '禁用' : '启用' }}
          </el-button>
          <el-button link type="danger" size="small" @click="handleDelete(row)">
            删除
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
