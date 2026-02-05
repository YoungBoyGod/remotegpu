<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

interface MenuItem {
  id: string
  title: string
  icon: string
  path?: string
  children?: MenuItem[]
  badge?: number
}

const router = useRouter()
const route = useRoute()

// 菜单配置
const menuItems = ref<MenuItem[]>([
  {
    id: 'dashboard',
    title: '工作台首页',
    icon: '🏠',
    path: '/customer/dashboard'
  },
  {
    id: 'machines',
    title: '我的机器',
    icon: '💻',
    children: [
      { id: 'machine-list', title: '机器列表', path: '/customer/machines/list', icon: '📋' },
      { id: 'machine-add', title: '添加机器', path: '/customer/machines/enroll', icon: '➕' },
      { id: 'machine-enrollments', title: '添加进度', path: '/customer/machines/enrollments', icon: '⏳' },
      { id: 'quick-connect', title: '快速连接', path: '/customer/machines/connect', icon: '🔗' },
      { id: 'connect-history', title: '连接历史', path: '/customer/machines/history', icon: '🕐' }
    ]
  },
  {
    id: 'tasks',
    title: '任务管理',
    icon: '🚀',
    children: [
      { id: 'training-tasks', title: '训练任务', path: '/customer/tasks/training', icon: '🎯' },
      { id: 'inference-tasks', title: '推理任务', path: '/customer/tasks/inference', icon: '🔮' },
      { id: 'task-queue', title: '任务队列', path: '/customer/tasks/queue', icon: '⏳' },
      { id: 'task-history', title: '历史记录', path: '/customer/tasks/history', icon: '📜' }
    ]
  },
  {
    id: 'images',
    title: '镜像市场',
    icon: '🐳',
    children: [
      { id: 'image-market', title: '镜像市场', path: '/customer/images/market', icon: '🛒' },
      { id: 'my-images', title: '我的镜像', path: '/customer/images/my', icon: '📦' }
    ]
  },
  {
    id: 'datasets',
    title: '数据集管理',
    icon: '📁',
    children: [
      { id: 'dataset-list', title: '数据集列表', path: '/customer/datasets/list', icon: '📂' },
      { id: 'upload-dataset', title: '上传数据集', path: '/customer/datasets/upload', icon: '📤' },
      { id: 'mount-manage', title: '挂载管理', path: '/customer/datasets/mounts', icon: '🔗' }
    ]
  },
  {
    id: 'models',
    title: '模型管理',
    icon: '🤖',
    path: '/customer/models/list'
  },
  {
    id: 'monitoring',
    title: '监控与分析',
    icon: '📊',
    children: [
      { id: 'realtime-monitor', title: '实时监控', path: '/customer/monitoring/realtime', icon: '📈' },
      { id: 'usage-stats', title: '使用统计', path: '/customer/statistics/usage', icon: '📉' }
    ]
  },
  {
    id: 'files',
    title: '文件管理',
    icon: '📂',
    path: '/customer/files/browser'
  },
  {
    id: 'tools',
    title: '开发工具',
    icon: '🔧',
    children: [
      { id: 'jupyter', title: 'Jupyter', path: '/customer/tools/jupyter', icon: '📓' },
      { id: 'tensorboard', title: 'TensorBoard', path: '/customer/tools/tensorboard', icon: '📉' },
      { id: 'terminal', title: 'Terminal', path: '/customer/tools/terminal', icon: '💻' }
    ]
  },
  {
    id: 'settings',
    title: '设置',
    icon: '⚙️',
    children: [
      { id: 'profile', title: '个人设置', path: '/customer/settings/profile', icon: '👤' },
      { id: 'ssh-keys', title: 'SSH密钥', path: '/customer/settings/ssh-keys', icon: '🔑' },
      { id: 'notifications', title: '通知设置', path: '/customer/settings/notifications', icon: '🔔' }
    ]
  }
])

const expandedMenus = ref<string[]>(['machines'])

const toggleMenu = (menuId: string) => {
  const index = expandedMenus.value.indexOf(menuId)
  if (index > -1) {
    expandedMenus.value.splice(index, 1)
  } else {
    expandedMenus.value.push(menuId)
  }
}

const isExpanded = (menuId: string) => expandedMenus.value.includes(menuId)
const isActive = (path?: string) => path && route.path === path

const navigateTo = (path?: string) => {
  if (path) router.push(path)
}
</script>

<template>
  <div class="customer-sidebar">
    <!-- Logo -->
    <div class="sidebar-header">
      <div class="logo">
        <span class="logo-icon">🚀</span>
        <span class="logo-text">RemoteGPU</span>
      </div>
    </div>

    <!-- 菜单 -->
    <el-scrollbar class="sidebar-menu">
      <div v-for="item in menuItems" :key="item.id" class="menu-group">
        <!-- 无子菜单 -->
        <div
          v-if="!item.children"
          class="menu-item"
          :class="{ active: isActive(item.path) }"
          @click="navigateTo(item.path)"
        >
          <span class="menu-icon">{{ item.icon }}</span>
          <span class="menu-title">{{ item.title }}</span>
        </div>

        <!-- 有子菜单 -->
        <div v-else>
          <div
            class="menu-item parent"
            :class="{ expanded: isExpanded(item.id) }"
            @click="toggleMenu(item.id)"
          >
            <span class="menu-icon">{{ item.icon }}</span>
            <span class="menu-title">{{ item.title }}</span>
            <span class="arrow">{{ isExpanded(item.id) ? '▼' : '▶' }}</span>
          </div>

          <div v-show="isExpanded(item.id)" class="submenu">
            <div
              v-for="child in item.children"
              :key="child.id"
              class="submenu-item"
              :class="{ active: isActive(child.path) }"
              @click="navigateTo(child.path)"
            >
              <span class="submenu-icon">{{ child.icon }}</span>
              <span class="submenu-title">{{ child.title }}</span>
            </div>
          </div>
        </div>
      </div>
    </el-scrollbar>
  </div>
</template>

<style scoped>
.customer-sidebar {
  width: 220px;
  height: 100vh;
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  color: #e4e4e4;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.15);
}

.sidebar-header {
  padding: 20px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-icon {
  font-size: 28px;
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.sidebar-menu {
  flex: 1;
  padding: 12px 0;
}

.menu-group {
  margin-bottom: 4px;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;
}

.menu-item:hover {
  background: rgba(255, 255, 255, 0.08);
}

.menu-item.active {
  background: linear-gradient(90deg, rgba(102, 126, 234, 0.3) 0%, transparent 100%);
  border-left: 3px solid #667eea;
  padding-left: 13px;
}

.menu-icon {
  font-size: 20px;
  margin-right: 10px;
  width: 24px;
  text-align: center;
}

.menu-title {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
}

.arrow {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.5);
  transition: transform 0.2s;
}

.submenu {
  background: rgba(0, 0, 0, 0.2);
  animation: slideDown 0.2s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    max-height: 0;
  }
  to {
    opacity: 1;
    max-height: 500px;
  }
}

.submenu-item {
  display: flex;
  align-items: center;
  padding: 10px 16px 10px 40px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 13px;
}

.submenu-item:hover {
  background: rgba(255, 255, 255, 0.06);
}

.submenu-item.active {
  background: rgba(102, 126, 234, 0.2);
  color: #8b9cff;
}

.submenu-icon {
  font-size: 16px;
  margin-right: 8px;
}

.submenu-title {
  flex: 1;
}
</style>
