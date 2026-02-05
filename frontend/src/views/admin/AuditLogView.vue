<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { getAuditLogs } from '@/api/admin'

interface AuditLog {
  id: number
  username: string
  action: string
  resource_type: string
  resource_id: string
  ip_address: string
  status_code: number
  created_at: string
  detail: any
}

const loading = ref(false)
const logs = ref<AuditLog[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const searchUser = ref('')
const searchAction = ref('')
const searchResource = ref('')

const loadLogs = async () => {
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      pageSize: pageSize.value
    }
    if (searchUser.value) params.username = searchUser.value
    if (searchAction.value) params.action = searchAction.value
    if (searchResource.value) params.resource_type = searchResource.value

    const res = await getAuditLogs(params)
    logs.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (val: number) => {
  page.value = val
  loadLogs()
}

onMounted(() => {
  loadLogs()
})
</script>

<template>
  <div class="audit-log-view">
    <PageHeader title="审计日志" />

    <el-card class="filter-card">
      <div class="filter-container">
        <el-input v-model="searchUser" placeholder="用户名" style="width: 150px" clearable @clear="loadLogs" />
        <el-input v-model="searchAction" placeholder="操作类型" style="width: 150px" clearable @clear="loadLogs" />
        <el-input v-model="searchResource" placeholder="资源类型" style="width: 150px" clearable @clear="loadLogs" />
        <el-button type="primary" :icon="Search" @click="loadLogs">查询</el-button>
      </div>
    </el-card>

    <div class="table-container">
      <el-table :data="logs" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户" width="120" />
        <el-table-column prop="action" label="操作" width="150" />
        <el-table-column prop="resource_type" label="资源类型" width="120" />
        <el-table-column prop="resource_id" label="资源ID" width="150" show-overflow-tooltip />
        <el-table-column prop="ip_address" label="IP地址" width="140" />
        <el-table-column prop="status_code" label="状态码" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status_code >= 200 && row.status_code < 300 ? 'success' : 'danger'">
              {{ row.status_code }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="detail-expand">
              <p><strong>Method:</strong> {{ row.method }} {{ row.path }}</p>
              <p><strong>Detail:</strong></p>
              <pre>{{ JSON.stringify(row.detail, null, 2) }}</pre>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.audit-log-view {
  padding: 24px;
}
.filter-card {
  margin-bottom: 20px;
}
.filter-container {
  display: flex;
  gap: 10px;
}
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
.detail-expand {
  padding: 10px 20px;
  background-color: #f8f9fa;
}
</style>
