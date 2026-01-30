<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { updateUserInfo } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'

const activeTab = ref('profile')
const authStore = useAuthStore()

const profileForm = reactive({
  username: '',
  email: '',
  nickname: '',
  avatar: '',
})

const loadProfile = async () => {
  await authStore.fetchProfile()
  if (authStore.user) {
    profileForm.username = authStore.user.username
    profileForm.email = authStore.user.email
    profileForm.nickname = authStore.user.nickname || ''
    profileForm.avatar = authStore.user.avatar || ''
  }
}

const updateProfile = async () => {
  try {
    await updateUserInfo({
      nickname: profileForm.nickname || undefined,
      avatar: profileForm.avatar || undefined,
    })
    await loadProfile()
    ElMessage.success('个人信息已更新')
  } catch (error) {
    ElMessage.error('更新失败')
  }
}

onMounted(() => {
  loadProfile()
})
</script>

<template>
  <div class="settings">
    <div class="page-header">
      <h1>个人设置</h1>
    </div>

    <el-tabs v-model="activeTab">
      <el-tab-pane label="个人信息" name="profile">
        <el-form :model="profileForm" label-width="120px" style="max-width: 600px">
          <el-form-item label="用户名">
            <el-input v-model="profileForm.username" disabled />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="profileForm.email" disabled />
          </el-form-item>
          <el-form-item label="昵称">
            <el-input v-model="profileForm.nickname" placeholder="请输入昵称" />
          </el-form-item>
          <el-form-item label="头像 URL">
            <el-input v-model="profileForm.avatar" placeholder="请输入头像地址" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="updateProfile">保存</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped>
.settings {
  padding: 24px;
}

.page-header h1 {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 24px 0;
}
</style>
