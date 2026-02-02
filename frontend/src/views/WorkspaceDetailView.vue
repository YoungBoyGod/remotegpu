<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, Delete, Plus, User } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { getWorkspaceById, getMembers, addMember, removeMember } from '@/api/workspace'
import type { WorkspaceInfo, WorkspaceMemberInfo, AddMemberRequest } from '@/api/workspace/types'

const router = useRouter()
const route = useRoute()

const workspaceId = Number(route.params.id)
const workspace = ref<WorkspaceInfo | null>(null)
const members = ref<WorkspaceMemberInfo[]>([])
const loading = ref(false)
const membersLoading = ref(false)
const activeTab = ref('info')

// 添加成员对话框
const addMemberDialogVisible = ref(false)
const addMemberForm = ref<AddMemberRequest>({
  user_id: 0,
  role: 'member',
})

// 加载工作空间详情
const loadWorkspace = async () => {
  loading.value = true
  try {
    const response = await getWorkspaceById(workspaceId)
    workspace.value = response.data
  } catch (error: any) {
    ElMessage.error(error.message || '加载工作空间信息失败')
    router.back()
  } finally {
    loading.value = false
  }
}

// 加载成员列表
const loadMembers = async () => {
  membersLoading.value = true
  try {
    const response = await getMembers(workspaceId)
    members.value = response.data
  } catch (error: any) {
    ElMessage.error(error.message || '加载成员列表失败')
  } finally {
    membersLoading.value = false
  }
}

// 编辑工作空间
const handleEdit = () => {
  router.push(`/portal/workspaces/${workspaceId}/edit`)
}

// 删除工作空间
const handleDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除工作空间 "${workspace.value?.name}" 吗？此操作不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    // TODO: 调用删除API
    ElMessage.success('删除成功')
    router.push('/portal/workspaces')
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 显示添加成员对话框
const showAddMemberDialog = () => {
  addMemberForm.value = {
    user_id: 0,
    role: 'member',
  }
  addMemberDialogVisible.value = true
}

// 添加成员
const handleAddMember = async () => {
  try {
    await addMember(workspaceId, addMemberForm.value)
    ElMessage.success('添加成员成功')
    addMemberDialogVisible.value = false
    loadMembers()
    loadWorkspace() // 重新加载以更新成员数量
  } catch (error: any) {
    ElMessage.error(error.message || '添加成员失败')
  }
}

// 移除成员
const handleRemoveMember = async (member: WorkspaceMemberInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要移除成员 "${member.username}" 吗？`,
      '移除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await removeMember(workspaceId, member.customer_id)
    ElMessage.success('移除成员成功')
    loadMembers()
    loadWorkspace() // 重新加载以更新成员数量
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '移除成员失败')
    }
  }
}

// 格式化日期
const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

// 角色标签类型
const getRoleType = (role: string) => {
  const roleMap: Record<string, any> = {
    owner: 'danger',
    admin: 'warning',
    member: 'info',
  }
  return roleMap[role] || 'info'
}

// 角色显示文本
const getRoleText = (role: string) => {
  const roleMap: Record<string, string> = {
    owner: '所有者',
    admin: '管理员',
    member: '成员',
  }
  return roleMap[role] || role
}

onMounted(() => {
  loadWorkspace()
  loadMembers()
})
</script>

<template>
  <div class="workspace-detail">
    <PageHeader :title="workspace?.name || '工作空间详情'">
      <template #actions>
        <el-button :icon="Edit" @click="handleEdit">编辑</el-button>
        <el-button type="danger" :icon="Delete" @click="handleDelete">
          删除
        </el-button>
      </template>
    </PageHeader>

    <el-tabs v-model="activeTab">
      <!-- 基本信息 -->
      <el-tab-pane label="基本信息" name="info">
        <el-card v-loading="loading">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">
              {{ workspace?.id }}
            </el-descriptions-item>
            <el-descriptions-item label="名称">
              {{ workspace?.name }}
            </el-descriptions-item>
            <el-descriptions-item label="描述" :span="2">
              {{ workspace?.description || '无' }}
            </el-descriptions-item>
            <el-descriptions-item label="所有者ID">
              {{ workspace?.owner_id }}
            </el-descriptions-item>
            <el-descriptions-item label="成员数量">
              <el-tag type="info">
                <el-icon><User /></el-icon>
                {{ workspace?.member_count }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ workspace?.created_at ? formatDate(workspace.created_at) : '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="更新时间">
              {{ workspace?.updated_at ? formatDate(workspace.updated_at) : '-' }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-tab-pane>

      <!-- 成员管理 -->
      <el-tab-pane label="成员管理" name="members">
        <el-card v-loading="membersLoading">
          <div class="members-header">
            <el-button type="primary" :icon="Plus" @click="showAddMemberDialog">
              添加成员
            </el-button>
          </div>

          <el-table :data="members" style="width: 100%; margin-top: 16px">
            <el-table-column prop="customer_id" label="用户ID" width="100" />
            <el-table-column prop="username" label="用户名" width="150" />
            <el-table-column prop="email" label="邮箱" min-width="180" />
            <el-table-column label="角色" width="120">
              <template #default="{ row }">
                <el-tag :type="getRoleType(row.role)">
                  {{ getRoleText(row.role) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="加入时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.joined_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button
                  v-if="row.role !== 'owner'"
                  type="danger"
                  size="small"
                  @click="handleRemoveMember(row)"
                >
                  移除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 添加成员对话框 -->
    <el-dialog
      v-model="addMemberDialogVisible"
      title="添加成员"
      width="500px"
    >
      <el-form :model="addMemberForm" label-width="100px">
        <el-form-item label="用户ID" required>
          <el-input-number
            v-model="addMemberForm.user_id"
            :min="1"
            placeholder="请输入用户ID"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="角色" required>
          <el-select v-model="addMemberForm.role" style="width: 100%">
            <el-option label="成员" value="member" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addMemberDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAddMember">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.workspace-detail {
  padding: 24px;
}

.members-header {
  display: flex;
  justify-content: flex-end;
}
</style>
