<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import FilterBar from '@/components/common/FilterBar.vue'
import { getQuotas, deleteQuota } from '@/api/quota'
import type { QuotaInfo } from '@/api/quota/types'

const router = useRouter()

const quotas = ref<QuotaInfo[]>([])
const loading = ref(false)
const searchText = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 过滤后的配额列表
const filteredQuotas = computed(() => {
  let result = quotas.value

  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(quota =>
      quota.customer_id.toString().includes(search) ||
      quota.quota_level.toLowerCase().includes(search)
    )
  }

  return result
})

// 加载配额列表
const loadQuotas = async () => {
  loading.value = true
  try {
    const response = await getQuotas(currentPage.value, pageSize.value)
    quotas.value = response.data.items
    total.value = response.data.total
  } catch (error: any) {
    ElMessage.error(error.message || '加载配额列表失败')
  } finally {
    loading.value = false
  }
}

// 创建配额
const handleCreate = () => {
  router.push('/admin/quotas/create')
}

// 编辑配额
const handleEdit = (quota: QuotaInfo) => {
  router.push(`/admin/quotas/${quota.id}/edit`)
}

// 查看详情
const handleDetail = (quota: QuotaInfo) => {
  router.push(`/admin/quotas/${quota.id}`)
}

// 删除配额
const handleDelete = async (quota: QuotaInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除配额 ID ${quota.id} 吗？此操作不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await deleteQuota(quota.id)
    ElMessage.success('删除成功')
    loadQuotas()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 分页变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  loadQuotas()
}

// 格式化日期
const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

// 配额级别标签类型
const getQuotaLevelType = (level: string) => {
  const levelMap: Record<string, any> = {
    free: 'info',
    basic: 'success',
    pro: 'warning',
    enterprise: 'danger',
  }
  return levelMap[level] || 'info'
}

// 配额级别显示文本
const getQuotaLevelText = (level: string) => {
  const levelMap: Record<string, string> = {
    free: '免费版',
    basic: '基础版',
    pro: '专业版',
    enterprise: '企业版',
  }
  return levelMap[level] || level
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

onMounted(() => {
  loadQuotas()
})
</script>

<template>
  <div class="quota-management">
    <PageHeader title="资源配额管理">
      <template #actions>
        <el-button type="primary" :icon="Plus" @click="handleCreate">
          设置配额
        </el-button>
      </template>
    </PageHeader>

    <FilterBar
      v-model:search-value="searchText"
      search-placeholder="搜索客户ID或配额级别"
    />

    <el-table
      v-loading="loading"
      :data="filteredQuotas"
      style="width: 100%"
      stripe
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="customer_id" label="客户ID" width="100" />
      <el-table-column label="工作空间ID" width="120">
        <template #default="{ row }">
          {{ row.workspace_id || '用户级' }}
        </template>
      </el-table-column>
      <el-table-column label="配额级别" width="120">
        <template #default="{ row }">
          <el-tag :type="getQuotaLevelType(row.quota_level)">
            {{ getQuotaLevelText(row.quota_level) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="GPU" width="80" align="center">
        <template #default="{ row }">
          {{ row.max_gpu }}
        </template>
      </el-table-column>
      <el-table-column label="CPU" width="80" align="center">
        <template #default="{ row }">
          {{ row.max_cpu }}
        </template>
      </el-table-column>
      <el-table-column label="内存" width="120" align="center">
        <template #default="{ row }">
          {{ formatMemory(row.max_memory) }}
        </template>
      </el-table-column>
      <el-table-column label="存储" width="120" align="center">
        <template #default="{ row }">
          {{ formatStorage(row.max_storage) }}
        </template>
      </el-table-column>
      <el-table-column label="环境数" width="100" align="center">
        <template #default="{ row }">
          {{ row.max_environments }}
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="handleDetail(row)">
            详情
          </el-button>
          <el-button size="small" :icon="Edit" @click="handleEdit(row)">
            编辑
          </el-button>
          <el-button
            size="small"
            type="danger"
            :icon="Delete"
            @click="handleDelete(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<style scoped>
.quota-management {
  padding: 24px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
