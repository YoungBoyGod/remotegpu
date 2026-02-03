import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

/**
 * 角色导航组合式函数
 * 根据用户角色自动添加路由前缀
 */
export function useRoleNavigation() {
  const router = useRouter()
  const authStore = useAuthStore()

  /**
   * 根据用户角色获取完整路径
   * @param path 基础路径，如 '/dashboard', '/environments'
   * @returns 带角色前缀的完整路径，如 '/admin/dashboard' 或 '/customer/dashboard'
   */
  const getRolePath = (path: string): string => {
    if (path.startsWith('/admin/') || path.startsWith('/customer/')) {
      return path
    }

    // 如果是登录、注册等公共页面，直接返回
    if (path === '/login' || path === '/register') {
      return path
    }

    // 根据用户角色添加前缀
    const role = authStore.user?.role
    if (role === 'admin') {
      return `/admin${path}`
    }
    return `/customer${path}`
  }

  /**
   * 角色感知的路由跳转
   * @param path 基础路径
   */
  const navigateTo = (path: string) => {
    router.push(getRolePath(path))
  }

  return {
    getRolePath,
    navigateTo,
  }
}
