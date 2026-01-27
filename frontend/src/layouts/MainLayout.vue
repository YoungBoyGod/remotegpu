<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useRoleNavigation } from '@/composables/useRoleNavigation'
import {
  Menu as IconMenu,
  Monitor,
  FolderOpened,
  Picture,
  TrendCharts,
  Setting,
  User,
  SwitchButton,
  DataBoard,
  Management,
  Shop,
  Upload,
  Document,
  DataAnalysis,
  Grid,
  Bell,
  Box,
  Download,
} from '@element-plus/icons-vue'

const router = useRouter()
const authStore = useAuthStore()
const { getRolePath, navigateTo } = useRoleNavigation()
const collapsed = ref(false)

// 管理员菜单项
const adminMenuItems = [
  { path: '/dashboard', icon: Monitor, label: '概览' },
  { path: '/resource-center', icon: Grid, label: '资源中心' },
  { path: '/resource-platform', icon: DataBoard, label: '资源平台' },
  { path: '/resource-list', icon: Management, label: '资源管理' },
  { path: '/monitoring-center', icon: DataAnalysis, label: '监控中心' },
  { path: '/alert-center', icon: Bell, label: '告警中心' },
  { path: '/customer-center', icon: User, label: '客户中心' },
  { path: '/computing-market', icon: Shop, label: '算力市场' },
  { path: '/release-version', icon: Upload, label: '发布版本' },
  { path: '/document-center', icon: Document, label: '文档中心' },
  { path: '/download-center', icon: Download, label: '下载中心' },
]

// 客户菜单项
const customerMenuItems = [
  { path: '/dashboard', icon: Monitor, label: '概览' },
  { path: '/environments', icon: IconMenu, label: '环境部署' },
  { path: '/datasets', icon: FolderOpened, label: '数据集' },
  { path: '/images', icon: Picture, label: '镜像' },
  { path: '/model-repository', icon: Box, label: '模型仓库' },
  { path: '/training', icon: TrendCharts, label: '训练任务' },
  { path: '/computing-market', icon: Shop, label: '算力市场' },
  { path: '/release-version', icon: Upload, label: '发布版本' },
  { path: '/document-center', icon: Document, label: '文档中心' },
  { path: '/download-center', icon: Download, label: '下载中心' },
  { path: '/settings', icon: Setting, label: '平台设置' },
]

// 根据用户角色返回对应的菜单项
const menuItems = computed(() => {
  const items = authStore.user?.role === 'admin' ? adminMenuItems : customerMenuItems
  return items.map(item => ({
    ...item,
    path: getRolePath(item.path)
  }))
})

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>

<template>
  <el-container class="main-layout">
    <el-aside :width="collapsed ? '64px' : '200px'" class="sidebar">
      <div class="logo">
        <h2 v-if="!collapsed">RemoteGPU</h2>
        <span v-else>RG</span>
      </div>
      <el-menu
        :default-active="$route.path"
        :collapse="collapsed"
        router
      >
        <el-menu-item
          v-for="item in menuItems"
          :key="item.path"
          :index="item.path"
        >
          <el-icon><component :is="item.icon" /></el-icon>
          <template #title>{{ item.label }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <el-icon class="collapse-icon" @click="collapsed = !collapsed">
          <IconMenu />
        </el-icon>
        <div class="header-right">
          <el-dropdown>
            <el-icon class="user-icon"><User /></el-icon>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="navigateTo('/settings')">
                  个人设置
                </el-dropdown-item>
                <el-dropdown-item @click="navigateTo('/workspace')">
                  工作空间
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.main-layout {
  height: 100vh;
}

.sidebar {
  background: #fbfbfc;
  transition: width 0.3s;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
  font-weight: 600;
}

.logo h2 {
  margin: 0;
}

.el-menu {
  border-right: none;
  background: #fbfcfd;
}

.header {
  background: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #f0f0f0;
  padding: 0 20px;
}

.collapse-icon {
  font-size: 20px;
  cursor: pointer;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.user-icon {
  font-size: 20px;
  cursor: pointer;
}

.main-content {
  background: white;
  overflow-y: auto;
}
</style>
