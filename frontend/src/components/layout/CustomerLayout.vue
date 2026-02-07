<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import CustomerSidebar from './CustomerSidebar.vue'
import NotificationBell from '@/components/common/NotificationBell.vue'
import { UserFilled, ArrowDown } from '@element-plus/icons-vue'
import type { UserInfo } from '@/types/common'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const userInfo = computed<UserInfo>(() => {
  return (
    authStore.user || {
      id: 0,
      username: 'customer',
      email: '',
      role: 'customer_member'
    }
  )
})

// 退出登录
const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="customer-layout">
    <!-- 侧边栏 -->
    <CustomerSidebar />

    <!-- 主内容区 -->
    <div class="main-container">
      <!-- 顶部导航栏 -->
      <div class="top-navbar">
        <div class="navbar-left">
          <span class="page-title">{{ $route.meta.title || '工作台' }}</span>
        </div>
        <div class="navbar-right">
          <!-- 通知铃铛 -->
          <NotificationBell />

          <!-- 租户信息 -->
          <div class="tenant-info">
            <span class="tenant-label">企业：</span>
            <span class="tenant-name">{{ userInfo.tenantName }}</span>
          </div>

          <!-- 用户信息 -->
          <el-dropdown @command="handleLogout">
            <div class="user-info">
              <el-avatar :size="32" :src="userInfo.avatar">
                <el-icon><UserFilled /></el-icon>
              </el-avatar>
              <span class="username">{{ userInfo.username }}</span>
              <el-icon class="arrow-icon"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>

      <!-- 页面内容 -->
      <div class="page-content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<style scoped>
.customer-layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #f5f7fa;
}

.top-navbar {
  height: 60px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.navbar-left {
  display: flex;
  align-items: center;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.tenant-info {
  display: flex;
  align-items: center;
  padding: 4px 12px;
  background: #f0f9ff;
  border-radius: 4px;
  border: 1px solid #d1e7ff;
}

.tenant-label {
  font-size: 12px;
  color: #909399;
  margin-right: 4px;
}

.tenant-name {
  font-size: 13px;
  font-weight: 600;
  color: #409eff;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background 0.3s;
}

.user-info:hover {
  background: #f5f7fa;
}

.username {
  font-size: 14px;
  color: #606266;
}

.arrow-icon {
  font-size: 12px;
  color: #909399;
}

.page-content {
  flex: 1;
  overflow: auto;
}
</style>
