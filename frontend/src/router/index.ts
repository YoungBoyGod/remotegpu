import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      redirect: (to) => {
        // 根据用户角色重定向到不同的首页
        const authStore = useAuthStore()
        if (authStore.user?.role === 'admin') {
          return '/admin/dashboard'
        }
        return '/portal/dashboard'
      },
    },
    // 管理员路由
    {
      path: '/admin',
      component: () => import('@/layouts/MainLayout.vue'),
      meta: { requiresAuth: true, requiresRole: 'admin' },
      children: [
        {
          path: '',
          redirect: '/admin/dashboard',
        },
        {
          path: 'dashboard',
          name: 'admin-dashboard',
          component: () => import('@/views/DashboardView.vue'),
        },
        {
          path: 'resource-center',
          name: 'admin-resource-center',
          component: () => import('@/views/ResourceCenterView.vue'),
        },
        {
          path: 'resource-platform',
          name: 'admin-resource-platform',
          component: () => import('@/views/ResourcePlatformView.vue'),
        },
        {
          path: 'resource-list',
          name: 'admin-resource-list',
          component: () => import('@/views/ResourceListView.vue'),
        },
        {
          path: 'monitoring-center',
          name: 'admin-monitoring-center',
          component: () => import('@/views/MonitoringCenterView.vue'),
        },
        {
          path: 'alert-center',
          name: 'admin-alert-center',
          component: () => import('@/views/AlertCenterView.vue'),
        },
        {
          path: 'customer-center',
          name: 'admin-customer-center',
          component: () => import('@/views/CustomerCenterView.vue'),
        },
        {
          path: 'customer-center/:id',
          name: 'admin-customer-dashboard',
          component: () => import('@/views/CustomerDashboardView.vue'),
        },
        {
          path: 'computing-market',
          name: 'admin-computing-market',
          component: () => import('@/views/ComputingMarketView.vue'),
        },
        {
          path: 'release-version',
          name: 'admin-release-version',
          component: () => import('@/views/ReleaseVersionView.vue'),
        },
        {
          path: 'document-center',
          name: 'admin-document-center',
          component: () => import('@/views/DocumentCenterView.vue'),
        },
        {
          path: 'download-center',
          name: 'admin-download-center',
          component: () => import('@/views/DownloadCenterView.vue'),
        },
      ],
    },
    // 客户门户路由
    {
      path: '/portal',
      component: () => import('@/layouts/MainLayout.vue'),
      meta: { requiresAuth: true, requiresRole: 'customer' },
      children: [
        {
          path: '',
          redirect: '/portal/dashboard',
        },
        {
          path: 'dashboard',
          name: 'portal-dashboard',
          component: () => import('@/views/DashboardView.vue'),
        },
        {
          path: 'environments',
          name: 'portal-environments',
          component: () => import('@/views/EnvironmentListView.vue'),
        },
        {
          path: 'environments/create',
          name: 'portal-environment-create',
          component: () => import('@/views/EnvironmentCreateView.vue'),
        },
        {
          path: 'environments/:id',
          name: 'portal-environment-detail',
          component: () => import('@/views/EnvironmentDetailView.vue'),
        },
        {
          path: 'hosts',
          name: 'portal-hosts',
          component: () => import('@/views/HostSelectionView.vue'),
        },
        {
          path: 'datasets',
          name: 'portal-datasets',
          component: () => import('@/views/DatasetListView.vue'),
        },
        {
          path: 'datasets/upload',
          name: 'portal-dataset-upload',
          component: () => import('@/views/DatasetUploadView.vue'),
        },
        {
          path: 'datasets/:id',
          name: 'portal-dataset-detail',
          component: () => import('@/views/DatasetDetailView.vue'),
        },
        {
          path: 'images',
          name: 'portal-images',
          component: () => import('@/views/ImageListView.vue'),
        },
        {
          path: 'images/build',
          name: 'portal-image-build',
          component: () => import('@/views/ImageBuildView.vue'),
        },
        {
          path: 'training',
          name: 'portal-training',
          component: () => import('@/views/TrainingListView.vue'),
        },
        {
          path: 'training/create',
          name: 'portal-training-create',
          component: () => import('@/views/TrainingCreateView.vue'),
        },
        {
          path: 'training/:id',
          name: 'portal-training-detail',
          component: () => import('@/views/TrainingDetailView.vue'),
        },
        {
          path: 'model-repository',
          name: 'portal-model-repository',
          component: () => import('@/views/ModelRepositoryView.vue'),
        },
        {
          path: 'settings',
          name: 'portal-settings',
          component: () => import('@/views/SettingsView.vue'),
        },
        {
          path: 'workspace',
          name: 'portal-workspace',
          component: () => import('@/views/WorkspaceSettingsView.vue'),
        },
        {
          path: 'computing-market',
          name: 'portal-computing-market',
          component: () => import('@/views/ComputingMarketView.vue'),
        },
        {
          path: 'release-version',
          name: 'portal-release-version',
          component: () => import('@/views/ReleaseVersionView.vue'),
        },
        {
          path: 'document-center',
          name: 'portal-document-center',
          component: () => import('@/views/DocumentCenterView.vue'),
        },
        {
          path: 'download-center',
          name: 'portal-download-center',
          component: () => import('@/views/DownloadCenterView.vue'),
        },
        {
          path: 'resource-center',
          name: 'resource-center',
          component: () => import('@/views/ResourceCenterView.vue'),
        },
        {
          path: 'resource-platform',
          name: 'resource-platform',
          component: () => import('@/views/ResourcePlatformView.vue'),
        },
        {
          path: 'resource-list',
          name: 'resource-list',
          component: () => import('@/views/ResourceListView.vue'),
        },
        {
          path: 'computing-market',
          name: 'computing-market',
          component: () => import('@/views/ComputingMarketView.vue'),
        },
        {
          path: 'release-version',
          name: 'release-version',
          component: () => import('@/views/ReleaseVersionView.vue'),
        },
        {
          path: 'document-center',
          name: 'document-center',
          component: () => import('@/views/DocumentCenterView.vue'),
        },
        {
          path: 'download-center',
          name: 'download-center',
          component: () => import('@/views/DownloadCenterView.vue'),
        },
        {
          path: 'monitoring-center',
          name: 'monitoring-center',
          component: () => import('@/views/MonitoringCenterView.vue'),
        },
        {
          path: 'alert-center',
          name: 'alert-center',
          component: () => import('@/views/AlertCenterView.vue'),
        },
        {
          path: 'customer-center',
          name: 'customer-center',
          component: () => import('@/views/CustomerCenterView.vue'),
        },
        {
          path: 'customer-center/:id',
          name: 'customer-dashboard',
          component: () => import('@/views/CustomerDashboardView.vue'),
        },
        {
          path: 'model-repository',
          name: 'model-repository',
          component: () => import('@/views/ModelRepositoryView.vue'),
        },
        {
          path: 'environments',
          name: 'environments',
          component: () => import('@/views/EnvironmentListView.vue'),
        },
        {
          path: 'environments/create',
          name: 'environment-create',
          component: () => import('@/views/EnvironmentCreateView.vue'),
        },
        {
          path: 'environments/:id',
          name: 'environment-detail',
          component: () => import('@/views/EnvironmentDetailView.vue'),
        },
        {
          path: 'hosts',
          name: 'hosts',
          component: () => import('@/views/HostSelectionView.vue'),
        },
        {
          path: 'datasets',
          name: 'datasets',
          component: () => import('@/views/DatasetListView.vue'),
        },
        {
          path: 'datasets/upload',
          name: 'dataset-upload',
          component: () => import('@/views/DatasetUploadView.vue'),
        },
        {
          path: 'datasets/:id',
          name: 'dataset-detail',
          component: () => import('@/views/DatasetDetailView.vue'),
        },
        {
          path: 'images',
          name: 'images',
          component: () => import('@/views/ImageListView.vue'),
        },
        {
          path: 'images/build',
          name: 'image-build',
          component: () => import('@/views/ImageBuildView.vue'),
        },
        {
          path: 'training',
          name: 'training',
          component: () => import('@/views/TrainingListView.vue'),
        },
        {
          path: 'training/create',
          name: 'training-create',
          component: () => import('@/views/TrainingCreateView.vue'),
        },
        {
          path: 'training/:id',
          name: 'training-detail',
          component: () => import('@/views/TrainingDetailView.vue'),
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/views/SettingsView.vue'),
        },
        {
          path: 'workspace',
          name: 'workspace',
          component: () => import('@/views/WorkspaceSettingsView.vue'),
        },
      ],
    },
  ],
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.meta.requiresAuth !== false
  const requiresRole = to.meta.requiresRole as 'admin' | 'customer' | undefined

  // 未登录用户访问需要认证的页面，重定向到登录页
  if (requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  }

  // 已登录用户访问登录页，根据角色重定向到对应的首页
  if (to.path === '/login' && authStore.isAuthenticated) {
    if (authStore.user?.role === 'admin') {
      next('/admin/dashboard')
    } else {
      next('/portal/dashboard')
    }
    return
  }

  // 检查角色权限
  if (requiresRole && authStore.user?.role !== requiresRole) {
    // 用户角色不匹配，重定向到对应角色的首页
    if (authStore.user?.role === 'admin') {
      next('/admin/dashboard')
    } else {
      next('/portal/dashboard')
    }
    return
  }

  next()
})

export default router
