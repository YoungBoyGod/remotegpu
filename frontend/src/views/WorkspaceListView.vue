<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, User } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import FilterBar from '@/components/common/FilterBar.vue'
import { getWorkspaces, deleteWorkspace } from '@/api/workspace'
import type { WorkspaceInfo } from '@/api/workspace/types'

const router = useRouter()

const workspaces = ref<WorkspaceInfo[]>([])
const loading = ref(false)
const searchText = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 过滤后的工作空间列表
const filteredWorkspaces = computed(() => {
  let result = workspaces.value

  // 搜索过滤
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(ws =>
      ws.name.toLowerCase().includes(search) ||
      ws.description.toLowerCase().includes(search)
    )
  }

  return result
})

// 加载工作空间列表
const loadWorkspaces = async () => {
  loading.value = true
  try {
    const response = await getWorkspaces(currentPage.value, pageSize.value)
    workspaces.value = response.data.items
    total.value = response.data.total
  } catch (error: any) {
    ElMessage.error(error.message || '加载工作空间列表失败')
  } finally {
    loading.value = false
  }
}

// 创建工作空间
const handleCreate = () => {
  router.push('/portal/workspaces/create')
}

// 编辑工作空间
const handleEdit = (workspace: WorkspaceInfo) => {
  router.push(`/portal/workspaces/${workspace.id}/edit`)
}

// 查看详情
const handleDetail = (workspace: WorkspaceInfo) => {
  router.push(`/portal/workspaces/${workspace.id}`)
}

// 删除工作空间
const handleDelete = async (workspace: WorkspaceInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除工作空间 "${workspace.name}" 吗？此操作不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await deleteWorkspace(workspace.id)
    ElMessage.success('删除成功')
    loadWorkspaces()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 分页变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  loadWorkspaces()
}

// 格式化日期
const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

onMounted(() => {
  loadWorkspaces()
})
</script>

<template>
  <div class="workspace-list">
    <PageHeader title="工作空间管理">
      <template #actions>
        <el-button type="primary" :icon="Plus" @click="handleCreate">
          创建工作空间
        </el-button>
      </template>
    </PageHeader>

    <FilterBar
      v-model:search-value="searchText"
      search-placeholder="搜索工作空间名称或描述"
    />

    <el-table
      v-loading="loading"
      :data="filteredWorkspaces"
      style="width: 100%"
      stripe
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" min-width="150" />
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column label="成员数量" width="100" align="center">
        <template #default="{ row }">
          <el-tag type="info">
            <el-icon><User /></el-icon>
            {{ row.member_count }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="更新时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.updated_at) }}
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
.workspace-list {
  padding: 24px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
