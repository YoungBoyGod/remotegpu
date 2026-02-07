<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// 表单数据
const loginForm = reactive({
  username: '',
  password: '',
  remember: false,
})

// 表单引用
const loginFormRef = ref<FormInstance>()

// 加载状态
const loading = ref(false)

// 表单验证规则
const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名或邮箱', trigger: 'blur' },
    { min: 3, max: 50, message: '长度在 3 到 50 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 32, message: '长度在 6 到 32 个字符', trigger: 'blur' },
  ],
}

// 提交登录
const handleLogin = async (formEl: FormInstance | undefined) => {
  if (!formEl) return

  await formEl.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const response = await authStore.login({
          username: loginForm.username,
          password: loginForm.password,
        })
        ElMessage.success('登录成功')
        // 首次登录强制改密
        if (response.data.mustChangePassword) {
          router.push('/change-password')
        } else {
          const redirect = route.query.redirect as string
          router.push(redirect || '/')
        }
      } catch (error: any) {
        const code = error?.code ?? error?.response?.data?.code
        const msg = error?.msg || error?.message || error?.response?.data?.msg
        if (code === 2003) {
          ElMessage.error('密码错误，请重试')
        } else if (code === 2006) {
          ElMessage.error('账号已禁用，请联系管理员')
        } else {
          ElMessage.error(msg || '登录失败，请检查用户名和密码')
        }
      } finally {
        loading.value = false
      }
    }
  })
}

</script>

<template>
  <div class="login-container">
    <div class="login-box">
      <!-- Logo 和标题 -->
      <div class="login-header">
        <h1 class="login-title">RemoteGPU</h1>
        <p class="login-subtitle">企业级 GPU 云平台</p>
      </div>

      <!-- 登录表单 -->
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="rules"
        class="login-form"
        size="large"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名或邮箱"
            :prefix-icon="User"
            clearable
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            :prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin(loginFormRef)"
          />
        </el-form-item>

        <el-form-item>
          <div class="login-options">
            <el-checkbox v-model="loginForm.remember">记住我</el-checkbox>
            <el-link type="primary" underline="never" @click="router.push('/forgot-password')">
              忘记密码？
            </el-link>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            class="login-button"
            :loading="loading"
            @click="handleLogin(loginFormRef)"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>

      <div class="login-footer">
        <span>请使用公司账号登录</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-box {
  width: 420px;
  padding: 40px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-title {
  font-size: 32px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px 0;
}

.login-subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.login-form {
  margin-top: 24px;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.login-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  color: #606266;
}
</style>
