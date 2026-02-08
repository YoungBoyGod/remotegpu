<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/common/PageHeader.vue'
import { useAuthStore } from '@/stores/auth'
import { changePassword, updateProfile } from '@/api/auth'
import type { UpdateProfileRequest } from '@/api/auth'

const authStore = useAuthStore()

const profileLoading = ref(false)
const profileForm = ref({
  username: '',
  email: '',
  display_name: '',
  full_name: '',
  phone: '',
  company: '',
  role: '',
})
const profileSubmitting = ref(false)

const passwordForm = ref({
  old_password: '',
  new_password: '',
  confirm_password: '',
})
const passwordSubmitting = ref(false)

const loadProfile = async () => {
  profileLoading.value = true
  try {
    await authStore.fetchProfile()
    const user = authStore.user
    if (user) {
      profileForm.value = {
        username: user.username || '',
        email: user.email || '',
        display_name: (user as any).display_name || '',
        full_name: (user as any).full_name || '',
        phone: (user as any).phone || '',
        company: (user as any).company || '',
        role: user.role || '',
      }
    }
  } catch {
    ElMessage.error('加载个人信息失败')
  } finally {
    profileLoading.value = false
  }
}

const handleUpdateProfile = async () => {
  profileSubmitting.value = true
  try {
    const data: UpdateProfileRequest = {
      display_name: profileForm.value.display_name || undefined,
      full_name: profileForm.value.full_name || undefined,
      phone: profileForm.value.phone || undefined,
      company: profileForm.value.company || undefined,
    }
    await updateProfile(data)
    await authStore.fetchProfile()
    ElMessage.success('个人信息已更新')
  } catch {
    ElMessage.error('更新失败')
  } finally {
    profileSubmitting.value = false
  }
}

const handleChangePassword = async () => {
  if (!passwordForm.value.old_password || !passwordForm.value.new_password) {
    ElMessage.warning('请填写完整密码信息')
    return
  }
  if (passwordForm.value.new_password.length < 6) {
    ElMessage.warning('新密码长度不能少于 6 位')
    return
  }
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    ElMessage.warning('两次输入的新密码不一致')
    return
  }
  passwordSubmitting.value = true
  try {
    await changePassword({
      old_password: passwordForm.value.old_password,
      new_password: passwordForm.value.new_password,
    })
    ElMessage.success('密码修改成功')
    passwordForm.value = { old_password: '', new_password: '', confirm_password: '' }
  } catch {
    ElMessage.error('密码修改失败，请检查原密码是否正确')
  } finally {
    passwordSubmitting.value = false
  }
}

const roleLabel = (role: string) => {
  const map: Record<string, string> = {
    admin: '管理员',
    customer_owner: '团队拥有者',
    customer_member: '团队成员',
    customer: '客户',
  }
  return map[role] || role
}

onMounted(() => {
  loadProfile()
})
</script>

<template>
  <div class="profile-view">
    <PageHeader title="个人信息" />

    <el-row :gutter="24">
      <!-- 个人资料 -->
      <el-col :span="14">
        <el-card v-loading="profileLoading">
          <template #header>
            <span class="card-title">基本资料</span>
          </template>
          <el-form label-width="100px" :model="profileForm">
            <el-form-item label="用户名">
              <el-input v-model="profileForm.username" disabled />
            </el-form-item>
            <el-form-item label="邮箱">
              <el-input v-model="profileForm.email" disabled />
            </el-form-item>
            <el-form-item label="角色">
              <el-tag>{{ roleLabel(profileForm.role) }}</el-tag>
            </el-form-item>
            <el-form-item label="显示名称">
              <el-input v-model="profileForm.display_name" placeholder="请输入显示名称" />
            </el-form-item>
            <el-form-item label="姓名">
              <el-input v-model="profileForm.full_name" placeholder="请输入姓名" />
            </el-form-item>
            <el-form-item label="手机号">
              <el-input v-model="profileForm.phone" placeholder="请输入手机号" />
            </el-form-item>
            <el-form-item label="公司">
              <el-input v-model="profileForm.company" placeholder="请输入公司名称" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="profileSubmitting" @click="handleUpdateProfile">
                保存修改
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 修改密码 -->
      <el-col :span="10">
        <el-card>
          <template #header>
            <span class="card-title">修改密码</span>
          </template>
          <el-form label-width="100px" :model="passwordForm">
            <el-form-item label="当前密码">
              <el-input
                v-model="passwordForm.old_password"
                type="password"
                show-password
                placeholder="请输入当前密码"
              />
            </el-form-item>
            <el-form-item label="新密码">
              <el-input
                v-model="passwordForm.new_password"
                type="password"
                show-password
                placeholder="请输入新密码（至少6位）"
              />
            </el-form-item>
            <el-form-item label="确认密码">
              <el-input
                v-model="passwordForm.confirm_password"
                type="password"
                show-password
                placeholder="请再次输入新密码"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="passwordSubmitting" @click="handleChangePassword">
                修改密码
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.profile-view {
  padding: 24px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}
</style>
