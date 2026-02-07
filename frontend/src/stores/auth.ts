import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, logout as logoutApi, refreshToken as refreshTokenApi, getCurrentUser } from '@/api/auth'
import type { LoginRequest } from '@/api/auth'
import type { UserInfo } from '@/types/common'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const accessToken = ref<string | null>(localStorage.getItem('accessToken'))
  const refreshToken = ref<string | null>(localStorage.getItem('refreshToken'))
  const user = ref<UserInfo | null>(null)

  // 计算属性
  const isAuthenticated = computed(() => !!accessToken.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isCustomer = computed(() =>
    !!user.value?.role && ['customer', 'customer_owner', 'customer_member'].includes(user.value.role)
  )

  // 登录
  const login = async (credentials: LoginRequest) => {
    const response = await loginApi(credentials)
    const { accessToken: newAccessToken, refreshToken: newRefreshToken, mustChangePassword } = response.data

    accessToken.value = newAccessToken
    refreshToken.value = newRefreshToken
    localStorage.setItem('accessToken', newAccessToken)
    localStorage.setItem('refreshToken', newRefreshToken)

    // 先尝试获取完整用户信息
    try {
      await fetchProfile()
    } catch {
      // fetchProfile 失败时，用登录响应中的 mustChangePassword 构建最小用户信息
      if (mustChangePassword) {
        user.value = { must_change_password: true } as UserInfo
      }
    }

    return response
  }

  // 登出
  const logout = async () => {
    try {
      await logoutApi()
    } catch (error) {
      console.error('登出API调用失败:', error)
    } finally {
      accessToken.value = null
      refreshToken.value = null
      user.value = null
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
    }
  }

  // 刷新访问令牌
  const refreshAccessToken = async () => {
    if (!refreshToken.value) {
      throw new Error('没有刷新令牌')
    }

    const response = await refreshTokenApi(refreshToken.value)
    const { accessToken: newAccessToken, refreshToken: newRefreshToken } = response.data

    accessToken.value = newAccessToken
    localStorage.setItem('accessToken', newAccessToken)

    if (newRefreshToken) {
      refreshToken.value = newRefreshToken
      localStorage.setItem('refreshToken', newRefreshToken)
    }

    return response
  }

  // 获取用户信息
  const fetchProfile = async () => {
    const response = await getCurrentUser()
    user.value = response.data
    return response
  }

  return {
    accessToken,
    refreshToken,
    user,
    isAuthenticated,
    isAdmin,
    isCustomer,
    login,
    logout,
    refreshAccessToken,
    fetchProfile,
  }
})
