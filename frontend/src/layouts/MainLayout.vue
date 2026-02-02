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
  {
    label: '资源',
    icon: Grid,
    children: [
      { path: '/resource-center', icon: Grid, label: '资源中心' },
      { path: '/resource-platform', icon: DataBoard, label: '资源平台' },
      { path: '/resource-list', icon: Management, label: '主机信息' },
      { path: '/computing-market', icon: Shop, label: '算力市场' },
    ]
  },
  {
    label: '监控',
    icon: DataAnalysis,
    children: [
      { path: '/monitoring-center', icon: DataAnalysis, label: '监控中心' },
      { path: '/alert-center', icon: Bell, label: '告警中心' },
    ]
  },
  {
    label: '客户',
    icon: User,
    children: [
      { path: '/customer-center', icon: User, label: '客户中心' },
      { path: '/quotas', icon: TrendCharts, label: '配额管理' },
    ]
  },
  {
    label: '发布',
    icon: Upload,
    children: [
      { path: '/release-version', icon: Upload, label: '发布版本' },
      { path: '/document-center', icon: Document, label: '文档中心' },
      { path: '/download-center', icon: Download, label: '下载中心' },
    ]
  },
  { path: '/settings', icon: Setting, label: '设置' },
]

// 客户菜单项
const customerMenuItems = [
  { path: '/environments', icon: Grid, label: '环境列表' },
  { path: '/datasets', icon: FolderOpened, label: '数据集' },
  { path: '/images', icon: Picture, label: '镜像' },
  { path: '/model-repository', icon: Box, label: '模型仓库' },
  { path: '/computing-market', icon: Shop, label: '算力市场' },
  { path: '/document-center', icon: Document, label: '文档中心' },
  { path: '/settings', icon: Setting, label: '设置' },
]

// 根据用户角色返回对应的菜单项
const menuItems = computed(() => {
  const items = authStore.user?.role === 'admin' ? adminMenuItems : customerMenuItems

  // 递归处理菜单项，为路径添加角色前缀
  const processItems = (items: any[]) => {
    return items.map(item => {
      if (item.children) {
        // 如果有子菜单，递归处理
        return {
          ...item,
          children: processItems(item.children)
        }
      } else if (item.path) {
        // 如果有路径，添加角色前缀
        return {
          ...item,
          path: getRolePath(item.path)
        }
      }
      return item
    })
  }

  return processItems(items)
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
        <template v-for="(item, index) in menuItems" :key="item.path || item.label">
          <!-- 分组菜单 -->
          <el-sub-menu v-if="item.children" :index="String(index)">
            <template #title>
              <el-icon><component :is="item.icon" /></el-icon>
              <span>{{ item.label }}</span>
            </template>
            <el-menu-item
              v-for="child in item.children"
              :key="child.path"
              :index="child.path"
            >
              <el-icon><component :is="child.icon" /></el-icon>
              <template #title>{{ child.label }}</template>
            </el-menu-item>
          </el-sub-menu>

          <!-- 普通菜单项 -->
          <el-menu-item v-else :index="item.path">
            <el-icon><component :is="item.icon" /></el-icon>
            <template #title>{{ item.label }}</template>
          </el-menu-item>
        </template>
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
