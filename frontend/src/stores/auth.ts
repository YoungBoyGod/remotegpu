import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login, logout, getCurrentUser } from '@/api/auth'
import type { LoginRequest, User } from '@/api/auth/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<User | null>(null)
  const isAuthenticated = ref<boolean>(!!token.value)

  // 登录
  const loginUser = async (credentials: LoginRequest) => {
    const response = await login(credentials)
    token.value = response.data.token
    user.value = response.data.user
    isAuthenticated.value = true
    localStorage.setItem('token', response.data.token)
    return response
  }

  // 登出
  const logoutUser = async () => {
    try {
      await logout()
    } finally {
      token.value = null
      user.value = null
      isAuthenticated.value = false
      localStorage.removeItem('token')
    }
  }

  // 获取用户信息
  const fetchProfile = async () => {
    const response = await getCurrentUser()
    user.value = response.data
    return response
  }

  return {
    token,
    user,
    isAuthenticated,
    login: loginUser,
    logout: logoutUser,
    fetchProfile,
  }
})
