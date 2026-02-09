<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, User, View } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import FilterBar from '@/components/common/FilterBar.vue'
import {
  getWorkspaces,
  createWorkspace,
  updateWorkspace,
  deleteWorkspace,
} from '@/api/workspace'
import type {
  WorkspaceInfo,
  CreateWorkspaceRequest,
  UpdateWorkspaceRequest,
} from '@/api/workspace/types'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()

// 列表数据
const workspaces = ref<WorkspaceInfo[]>([])
const loading = ref(false)
const searchText = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 创建/编辑对话框
const dialogVisible = ref(false)
const dialogTitle = ref('创建工作空间')
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()
const submitting = ref(false)
const formData = ref<CreateWorkspaceRequest>({
  name: '',
  description: '',
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入工作空间名称', trigger: 'blur' },
    { min: 2, max: 100, message: '名称长度在 2 到 100 个字符', trigger: 'blur' },
  ],
  description: [
    { max: 500, message: '描述不能超过 500 个字符', trigger: 'blur' },
  ],
}

// 过滤后的列表
const filteredWorkspaces = computed(() => {
  if (!searchText.value) return workspaces.value
  const search = searchText.value.toLowerCase()
  return workspaces.value.filter(
    (ws) =>
      ws.name.toLowerCase().includes(search) ||
      ws.description.toLowerCase().includes(search)
  )
})

// 加载列表
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

// 打开创建对话框
const handleCreate = () => {
  dialogTitle.value = '创建工作空间'
  editingId.value = null
  formData.value = { name: '', description: '' }
  dialogVisible.value = true
}

// 打开编辑对话框
const handleEdit = (ws: WorkspaceInfo) => {
  dialogTitle.value = '编辑工作空间'
  editingId.value = ws.id
  formData.value = { name: ws.name, description: ws.description }
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    submitting.value = true

    if (editingId.value) {
      const updateData: UpdateWorkspaceRequest = {
        name: formData.value.name,
        description: formData.value.description,
      }
      await updateWorkspace(editingId.value, updateData)
      ElMessage.success('更新成功')
    } else {
      await createWorkspace(formData.value)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    loadWorkspaces()
  } catch (error: any) {
    if (error.message) {
      ElMessage.error(error.message)
    }
  } finally {
    submitting.value = false
  }
}

// 查看详情
const handleDetail = (ws: WorkspaceInfo) => {
  router.push(`/customer/workspaces/${ws.id}`)
}

// 删除
const handleDelete = async (ws: WorkspaceInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除工作空间 "${ws.name}" 吗？此操作不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    await deleteWorkspace(ws.id)
    ElMessage.success('删除成功')
    loadWorkspaces()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 分页
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
  <div class="workspace-view">
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
      <el-table-column prop="name" label="名称" min-width="150">
        <template #default="{ row }">
          <el-link type="primary" @click="handleDetail(row)">
            {{ row.name }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column
        prop="description"
        label="描述"
        min-width="200"
        show-overflow-tooltip
      />
      <el-table-column label="成员" width="100" align="center">
        <template #default="{ row }">
          <el-tag type="info" size="small">
            <el-icon style="vertical-align: middle"><User /></el-icon>
            {{ row.member_count }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button size="small" :icon="View" @click="handleDetail(row)">
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

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="120px"
      >
        <el-form-item label="工作空间名称" prop="name">
          <el-input
            v-model="formData.name"
            placeholder="请输入工作空间名称"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="4"
            placeholder="请输入工作空间描述（可选）"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ editingId ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.workspace-view {
  padding: 24px;
  background: #f5f7fa;
  min-height: 100%;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

:deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;
}

:deep(.el-table th) {
  background: #f7f8fa !important;
  color: #4e5969;
  font-weight: 600;
  font-size: 13px;
}

:deep(.el-table td) {
  font-size: 13px;
  color: #1d2129;
}
</style>
