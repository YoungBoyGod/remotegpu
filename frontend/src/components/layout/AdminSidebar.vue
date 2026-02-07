<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  House,
  Monitor,
  User,
  Connection,
  DataAnalysis,
  Box,
  List,
  Setting,
  ArrowRight,
  ArrowDown
} from '@element-plus/icons-vue'

interface MenuItem {
  id: string
  title: string
  icon?: any
  path?: string
  children?: MenuItem[]
  badge?: number
}

const router = useRouter()
const route = useRoute()

// 菜单配置
const menuItems: MenuItem[] = [
  {
    id: 'dashboard',
    title: '管理后台首页',
    icon: House,
    path: '/admin/dashboard'
  },
  {
    id: 'machines',
    title: '机器管理',
    icon: Monitor,
    children: [
      { id: 'machine-list', title: '机器列表', path: '/admin/machines/list' },
      { id: 'add-machine', title: '添加机器', path: '/admin/machines/add' },
      { id: 'batch-import', title: '批量导入', path: '/admin/machines/import' }
    ]
  },
  {
    id: 'customers',
    title: '客户管理',
    icon: User,
    children: [
      { id: 'customer-list', title: '客户列表', path: '/admin/customers/list' }
    ]
  },
  {
    id: 'allocations',
    title: '分配管理',
    icon: Connection,
    children: [
      { id: 'allocation-list', title: '分配记录', path: '/admin/allocations/list' },
      { id: 'machine-allocate', title: '机器分配', path: '/admin/allocations/assign' },
    ]
  },
  {
    id: 'monitoring',
    title: '监控中心',
    icon: DataAnalysis,
    children: [
      { id: 'realtime-monitor', title: '实时监控', path: '/admin/monitoring' },
      { id: 'alerts', title: '告警管理', path: '/admin/alerts' }
    ]
  },
  {
    id: 'images',
    title: '镜像管理',
    icon: Box,
    path: '/admin/images'
  },
  {
    id: 'tasks',
    title: '任务管理',
    icon: List,
    path: '/admin/tasks/list'
  },
  {
    id: 'settings',
    title: '系统设置',
    icon: Setting,
    children: [
      { id: 'platform-config', title: '平台配置', path: '/admin/settings/platform' },
      { id: 'audit-log', title: '审计日志', path: '/admin/audit' }
    ]
  }
]

// 展开的菜单项
const expandedMenus = ref<string[]>(['machines'])

// 切换菜单展开/收起
const toggleMenu = (menuId: string) => {
  const index = expandedMenus.value.indexOf(menuId)
  if (index > -1) {
    expandedMenus.value.splice(index, 1)
  } else {
    expandedMenus.value.push(menuId)
  }
}

// 判断菜单是否展开
const isExpanded = (menuId: string) => {
  return expandedMenus.value.includes(menuId)
}

// 判断菜单是否激活
const isActive = (path?: string) => {
  if (!path) return false
  return route.path === path
}

// 导航到指定路径
const navigateTo = (path?: string) => {
  if (path) {
    router.push(path)
  }
}
</script>

<template>
  <div class="admin-sidebar">
    <!-- Logo区域 -->
    <div class="sidebar-header">
      <div class="logo">
        <span class="logo-icon">🚀</span>
        <span class="logo-text">RemoteGPU</span>
      </div>
      <div class="admin-badge">管理后台</div>
    </div>

    <!-- 菜单列表 -->
    <el-scrollbar class="sidebar-menu">
      <div v-for="item in menuItems" :key="item.id" class="menu-item-wrapper">
        <!-- 一级菜单（无子菜单） -->
        <div
          v-if="!item.children"
          class="menu-item"
          :class="{ active: isActive(item.path) }"
          @click="navigateTo(item.path)"
        >
          <el-icon class="menu-icon"><component :is="item.icon" /></el-icon>
          <span class="menu-title">{{ item.title }}</span>
          <el-badge v-if="item.badge" :value="item.badge" class="menu-badge" />
        </div>

        <!-- 有子菜单的一级菜单 -->
        <div v-else>
          <div
            class="menu-item with-children"
            :class="{ expanded: isExpanded(item.id) }"
            @click="toggleMenu(item.id)"
          >
            <el-icon class="menu-icon"><component :is="item.icon" /></el-icon>
            <span class="menu-title">{{ item.title }}</span>
            <el-badge v-if="item.badge" :value="item.badge" class="menu-badge" />
            <el-icon class="expand-icon">
              <ArrowRight v-if="!isExpanded(item.id)" />
              <ArrowDown v-else />
            </el-icon>
          </div>

          <!-- 二级菜单 -->
          <transition name="submenu">
            <div v-show="isExpanded(item.id)" class="submenu">
              <div
                v-for="child in item.children"
                :key="child.id"
                class="submenu-item"
                :class="{ active: isActive(child.path) }"
                @click="navigateTo(child.path)"
              >
                <span class="submenu-title">{{ child.title }}</span>
                <el-badge v-if="child.badge" :value="child.badge" class="submenu-badge" />
              </div>
            </div>
          </transition>
        </div>
      </div>
    </el-scrollbar>
  </div>
</template>

<style scoped>
.admin-sidebar {
  width: 240px;
  height: 100vh;
  background: #001529;
  display: flex;
  flex-direction: column;
  color: rgba(255, 255, 255, 0.85);
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  font-size: 18px;
  font-weight: 600;
}

.logo-icon {
  font-size: 24px;
}

.admin-badge {
  display: inline-block;
  padding: 2px 8px;
  background: rgba(24, 144, 255, 0.2);
  border: 1px solid #1890ff;
  border-radius: 4px;
  font-size: 12px;
  color: #40a9ff;
}

.sidebar-menu {
  flex: 1;
  padding: 8px 0;
}

.menu-item-wrapper {
  margin-bottom: 4px;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
}

.menu-item:hover {
  background: rgba(255, 255, 255, 0.08);
}

.menu-item.active {
  background: #1890ff;
  color: #fff;
}

.menu-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: #fff;
}

.menu-icon {
  font-size: 18px;
  margin-right: 12px;
}

.menu-title {
  flex: 1;
  font-size: 14px;
}

.menu-badge {
  margin-left: 8px;
}

.expand-icon {
  font-size: 14px;
  transition: transform 0.3s;
}

.submenu {
  background: rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

.submenu-item {
  padding: 10px 16px 10px 48px;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 13px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.submenu-item:hover {
  background: rgba(255, 255, 255, 0.08);
}

.submenu-item.active {
  background: rgba(24, 144, 255, 0.3);
  color: #40a9ff;
}

.submenu-title {
  flex: 1;
}

.submenu-badge {
  margin-left: 8px;
}

/* 动画效果 */
.submenu-enter-active,
.submenu-leave-active {
  transition: all 0.3s ease;
}

.submenu-enter-from,
.submenu-leave-to {
  max-height: 0;
  opacity: 0;
}

.submenu-enter-to,
.submenu-leave-from {
  max-height: 500px;
  opacity: 1;
}
</style>
