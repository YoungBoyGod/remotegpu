import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// 布局组件
const AdminLayout = () => import('@/components/layout/AdminLayout.vue')
const CustomerLayout = () => import('@/components/layout/CustomerLayout.vue')

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
      path: '/forgot-password',
      name: 'forgot-password',
      component: () => import('@/views/ForgotPasswordView.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      redirect: (to) => {
        const authStore = useAuthStore()
        if (!authStore.isAuthenticated) {
          return '/login'
        }
        if (authStore.user?.role === 'admin') {
          return '/admin/dashboard'
        }
        return '/customer/dashboard'
      },
    },
    // 管理员路由
    {
      path: '/admin',
      component: AdminLayout,
      meta: { requiresAuth: true, requiresRole: 'admin' },
      children: [
        {
          path: '',
          redirect: '/admin/dashboard',
        },
        {
          path: 'dashboard',
          name: 'admin-dashboard',
          component: () => import('@/views/admin/DashboardView.vue'),
          meta: { title: '管理后台首页' },
        },
        // 机器管理
        {
          path: 'machines/list',
          name: 'admin-machines-list',
          component: () => import('@/views/admin/MachineListView.vue'),
          meta: { title: '机器列表' },
        },
        {
          path: 'machines/add',
          name: 'admin-machines-add',
          component: () => import('@/views/admin/MachineAddView.vue'),
          meta: { title: '添加机器' },
        },
        {
          path: 'machines/:id',
          name: 'admin-machines-detail',
          component: () => import('@/views/admin/MachineDetailView.vue'),
          meta: { title: '机器详情' },
        },
        // 客户管理
        {
          path: 'customers/list',
          name: 'admin-customers-list',
          component: () => import('@/views/admin/CustomerListView.vue'),
          meta: { title: '客户列表' },
        },
        {
          path: 'customers/:id',
          name: 'admin-customers-detail',
          component: () => import('@/views/admin/CustomerDetailView.vue'),
          meta: { title: '客户详情' },
        },
        // 分配管理
        {
          path: 'allocations/list',
          name: 'admin-allocations-list',
          component: () => import('@/views/admin/AllocationListView.vue'),
          meta: { title: '分配记录' },
        },
        {
          path: 'allocations/quick',
          name: 'admin-allocations-quick',
          component: () => import('@/views/admin/QuickAllocateView.vue'),
          meta: { title: '快速分配' },
        },
        // 运维管理
        {
          path: 'images',
          name: 'admin-images',
          component: () => import('@/views/admin/ImageListView.vue'),
          meta: { title: '镜像管理' },
        },
        {
          path: 'monitoring',
          name: 'admin-monitoring',
          component: () => import('@/views/admin/MonitoringView.vue'),
          meta: { title: '监控中心' },
        },
        {
          path: 'alerts',
          name: 'admin-alerts',
          component: () => import('@/views/admin/AlertListView.vue'),
          meta: { title: '告警中心' },
        },
        {
          path: 'audit',
          name: 'admin-audit',
          component: () => import('@/views/admin/AuditLogView.vue'),
          meta: { title: '审计日志' },
        },
        {
          path: ':pathMatch(.*)*',
          name: 'admin-coming-soon',
          component: () => import('@/views/ComingSoonView.vue'),
          meta: { title: '功能开发中' },
        },
      ],
    },
    // 客户路由
    {
      path: '/customer',
      component: CustomerLayout,
      meta: { requiresAuth: true, requiresRole: ['customer_owner', 'customer_member'] },
      children: [
        {
          path: '',
          redirect: '/customer/dashboard',
        },
        {
          path: 'dashboard',
          name: 'customer-dashboard',
          component: () => import('@/views/customer/DashboardView.vue'),
          meta: { title: '工作台首页' },
        },
        // 我的机器
        {
          path: 'machines/list',
          name: 'customer-machines-list',
          component: () => import('@/views/customer/MachineListView.vue'),
          meta: { title: '机器列表' },
        },
        {
          path: 'machines/enroll',
          name: 'customer-machines-enroll',
          component: () => import('@/views/customer/MachineEnrollView.vue'),
          meta: { title: '添加机器' },
        },
        {
          path: 'machines/enrollments',
          name: 'customer-machines-enrollments',
          component: () => import('@/views/customer/MachineEnrollmentListView.vue'),
          meta: { title: '添加进度' },
        },
        // 资源管理
        {
          path: 'ssh-keys',
          name: 'customer-ssh-keys',
          component: () => import('@/views/customer/SSHKeyView.vue'),
          meta: { title: 'SSH 密钥' },
        },
        {
          path: 'tasks',
          name: 'customer-tasks',
          component: () => import('@/views/customer/TaskListView.vue'),
          meta: { title: '任务管理' },
        },
        {
          path: 'datasets',
          name: 'customer-datasets',
          component: () => import('@/views/customer/DatasetListView.vue'),
          meta: { title: '数据集管理' },
        },
        {
          path: ':pathMatch(.*)*',
          name: 'customer-coming-soon',
          component: () => import('@/views/ComingSoonView.vue'),
          meta: { title: '功能开发中' },
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFoundView.vue'),
      meta: { requiresAuth: false },
    },
  ],
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.meta.requiresAuth !== false
  const requiresRole = to.meta.requiresRole as string | string[] | undefined

  // 未登录用户访问需要认证的页面，重定向到登录页
  if (requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  }

  if (requiresAuth && authStore.isAuthenticated && !authStore.user) {
    try {
      await authStore.fetchProfile()
    } catch (error) {
      await authStore.logout()
      next('/login')
      return
    }
  }

  // 已登录用户访问登录页，根据角色重定向到对应的首页
  if (to.path === '/login' && authStore.isAuthenticated) {
    if (authStore.user?.role === 'admin') {
      next('/admin/dashboard')
    } else {
      next('/customer/dashboard')
    }
    return
  }

  // 检查角色权限
  if (requiresRole) {
    const role = authStore.user?.role
    const allowedRoles = Array.isArray(requiresRole) ? requiresRole : [requiresRole]
    if (!role || !allowedRoles.includes(role)) {
      if (role === 'admin') {
        next('/admin/dashboard')
      } else {
        next('/customer/dashboard')
      }
      return
    }
  }

  next()
})

export default router
