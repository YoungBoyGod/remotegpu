<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const activeTab = ref('basic')

const workspaceForm = reactive({
  name: '我的工作空间',
  description: '',
  type: 'personal',
})

const members = ref([
  { id: 1, username: 'user1', role: 'admin', joinedAt: '2026-01-20' },
  { id: 2, username: 'user2', role: 'member', joinedAt: '2026-01-22' },
])

const updateWorkspace = async () => {
  ElMessage.success('工作空间信息已更新')
}
</script>

<template>
  <div class="workspace-settings">
    <div class="page-header">
      <h1>工作空间设置</h1>
    </div>

    <el-tabs v-model="activeTab">
      <el-tab-pane label="基本信息" name="basic">
        <el-form :model="workspaceForm" label-width="120px" style="max-width: 600px">
          <el-form-item label="工作空间名称">
            <el-input v-model="workspaceForm.name" />
          </el-form-item>
          <el-form-item label="描述">
            <el-input v-model="workspaceForm.description" type="textarea" :rows="3" />
          </el-form-item>
          <el-form-item label="类型">
            <el-radio-group v-model="workspaceForm.type">
              <el-radio value="personal">个人</el-radio>
              <el-radio value="team">团队</el-radio>
              <el-radio value="enterprise">企业</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="updateWorkspace">保存</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="成员管理" name="members">
        <div class="members-section">
          <el-button type="primary" style="margin-bottom: 16px">邀请成员</el-button>
          <el-table :data="members" style="width: 100%">
            <el-table-column prop="username" label="用户名" />
            <el-table-column prop="role" label="角色" />
            <el-table-column prop="joinedAt" label="加入时间" />
            <el-table-column label="操作" width="150">
              <template #default>
                <el-button type="danger" size="small">移除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <el-tab-pane label="配额管理" name="quota">
        <el-descriptions title="资源配额" :column="2" border>
          <el-descriptions-item label="CPU">32 核</el-descriptions-item>
          <el-descriptions-item label="内存">128 GB</el-descriptions-item>
          <el-descriptions-item label="GPU">8 个</el-descriptions-item>
          <el-descriptions-item label="存储">500 GB</el-descriptions-item>
        </el-descriptions>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped>
.workspace-settings {
  padding: 24px;
}

.page-header h1 {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 24px 0;
}

.members-section {
  max-width: 800px;
}
</style>
