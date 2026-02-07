<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const goHome = () => {
  if (authStore.user?.role === 'admin') {
    router.push('/admin/dashboard')
  } else {
    router.push('/customer/dashboard')
  }
}
</script>

<template>
  <div class="forbidden">
    <h2 class="title">403 - 无权访问</h2>
    <p class="subtitle">你没有权限访问该页面，请联系管理员。</p>
    <el-button type="primary" @click="goHome">返回首页</el-button>
  </div>
</template>

<style scoped>
.forbidden {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 60vh;
  color: #606266;
}

.title {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 14px;
  margin-bottom: 16px;
}
</style>
