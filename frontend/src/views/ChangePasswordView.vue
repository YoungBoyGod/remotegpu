<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { changePassword } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const rules: FormRules = {
  oldPassword: [{ required: true, message: '请输入旧密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 32, message: '长度在 6 到 32 个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (value !== form.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
}

const handleSubmit = async (formEl: FormInstance | undefined) => {
  if (!formEl) return

  await formEl.validate(async (valid) => {
    if (!valid) return
    if (form.oldPassword === form.newPassword) {
      ElMessage.warning('新密码不能与旧密码相同')
      return
    }

    loading.value = true
    try {
      await changePassword({
        old_password: form.oldPassword,
        new_password: form.newPassword,
      })
      await authStore.fetchProfile()
      ElMessage.success('密码修改成功')
      router.push('/')
    } catch (error: any) {
      ElMessage.error(error?.msg || error.response?.data?.msg || '修改失败，请稍后重试')
    } finally {
      loading.value = false
    }
  })
}
</script>

<template>
  <div class="change-password-container">
    <div class="change-password-box">
      <div class="change-password-header">
        <h1>首次登录需修改密码</h1>
        <p>为了账号安全，请设置新的登录密码</p>
      </div>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        class="change-password-form"
        size="large"
      >
        <el-form-item prop="oldPassword">
          <el-input
            v-model="form.oldPassword"
            type="password"
            placeholder="旧密码"
            show-password
          />
        </el-form-item>
        <el-form-item prop="newPassword">
          <el-input
            v-model="form.newPassword"
            type="password"
            placeholder="新密码"
            show-password
          />
        </el-form-item>
        <el-form-item prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="确认新密码"
            show-password
            @keyup.enter="handleSubmit(formRef)"
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            class="submit-button"
            :loading="loading"
            @click="handleSubmit(formRef)"
          >
            保存新密码
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.change-password-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #ffe29f 0%, #ffa99f 100%);
}

.change-password-box {
  width: 420px;
  padding: 40px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 12px 36px rgba(0, 0, 0, 0.12);
}

.change-password-header {
  text-align: center;
  margin-bottom: 28px;
}

.change-password-header h1 {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 6px 0;
  color: #303133;
}

.change-password-header p {
  margin: 0;
  font-size: 13px;
  color: #909399;
}

.submit-button {
  width: 100%;
}
</style>
