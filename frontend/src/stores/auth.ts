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

  // 登录
  const login = async (credentials: LoginRequest) => {
    const response = await loginApi(credentials)
    const { accessToken: newAccessToken, refreshToken: newRefreshToken } = response.data

    accessToken.value = newAccessToken
    refreshToken.value = newRefreshToken
    localStorage.setItem('accessToken', newAccessToken)
    localStorage.setItem('refreshToken', newRefreshToken)

    await fetchProfile()

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
    login,
    logout,
    refreshAccessToken,
    fetchProfile,
  }
})
