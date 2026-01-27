<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'

const activeTab = ref('profile')

const profileForm = reactive({
  username: 'admin',
  email: 'admin@example.com',
  avatar: '',
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const updateProfile = async () => {
  ElMessage.success('个人信息已更新')
}

const changePassword = async () => {
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }
  ElMessage.success('密码修改成功')
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
}
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
            <el-input v-model="profileForm.username" />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="profileForm.email" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="updateProfile">保存</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="安全设置" name="security">
        <el-form :model="passwordForm" label-width="120px" style="max-width: 600px">
          <el-form-item label="当前密码">
            <el-input v-model="passwordForm.oldPassword" type="password" />
          </el-form-item>
          <el-form-item label="新密码">
            <el-input v-model="passwordForm.newPassword" type="password" />
          </el-form-item>
          <el-form-item label="确认密码">
            <el-input v-model="passwordForm.confirmPassword" type="password" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="changePassword">修改密码</el-button>
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
