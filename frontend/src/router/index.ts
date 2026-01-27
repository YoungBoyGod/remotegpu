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
      component: () => import('@/layouts/MainLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/dashboard',
        },
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/DashboardView.vue'),
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

  if (requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
