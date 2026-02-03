<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { requestPasswordReset, confirmPasswordReset } from '@/api/auth'

const router = useRouter()
const loading = ref(false)
const sending = ref(false)

const formRef = ref<FormInstance>()

const formData = reactive({
  email: '',
  code: '',
  password: '',
  confirmPassword: ''
})

const rules: FormRules = {
  email: [{ required: true, message: '请输入企业邮箱或工号', trigger: 'blur' }],
  code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
  password: [{ required: true, message: '请输入新密码', trigger: 'blur' }],
  confirmPassword: [{ required: true, message: '请再次输入新密码', trigger: 'blur' }]
}

const sendCode = async () => {
  if (!formData.email) {
    ElMessage.warning('请先输入企业邮箱或工号')
    return
  }

  try {
    sending.value = true
    await requestPasswordReset({ email: formData.email })
    ElMessage.success('验证码已发送')
  } catch (error: any) {
    ElMessage.error(error?.message || '发送验证码失败')
  } finally {
    sending.value = false
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return
    if (formData.password !== formData.confirmPassword) {
      ElMessage.warning('两次输入的密码不一致')
      return
    }

    try {
      loading.value = true
      await confirmPasswordReset({
        email: formData.email,
        code: formData.code,
        password: formData.password
      })
      ElMessage.success('密码重置成功，请重新登录')
      router.push('/login')
    } catch (error: any) {
      ElMessage.error(error?.message || '密码重置失败')
    } finally {
      loading.value = false
    }
  })
}
</script>

<template>
  <div class="forgot-password">
    <el-card class="card">
      <h2 class="title">找回密码</h2>
      <el-form ref="formRef" :model="formData" :rules="rules" label-position="top">
        <el-form-item label="企业邮箱/工号" prop="email">
          <el-input v-model="formData.email" placeholder="请输入企业邮箱或工号" />
        </el-form-item>
        <el-form-item label="验证码" prop="code">
          <div class="code-row">
            <el-input v-model="formData.code" placeholder="请输入验证码" />
            <el-button :loading="sending" @click="sendCode">发送验证码</el-button>
          </div>
        </el-form-item>
        <el-form-item label="新密码" prop="password">
          <el-input v-model="formData.password" type="password" show-password placeholder="请输入新密码" />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="formData.confirmPassword" type="password" show-password placeholder="请再次输入新密码" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">提交</el-button>
          <el-button @click="router.push('/login')">返回登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.forgot-password {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: #f5f7fa;
}

.card {
  width: 420px;
}

.title {
  text-align: center;
  margin-bottom: 16px;
}

.code-row {
  display: flex;
  gap: 12px;
}

.code-row :deep(.el-input) {
  flex: 1;
}
</style>
