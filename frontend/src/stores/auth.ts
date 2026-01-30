import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login, getCurrentUser } from '@/api/auth'
import type { LoginRequest, User, BackendUserInfo } from '@/api/auth/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<User | null>(null)
  const isAuthenticated = ref<boolean>(!!token.value)

  const mapUser = (backendUser: BackendUserInfo): User => ({
    id: backendUser.id,
    username: backendUser.username,
    email: backendUser.email,
    nickname: backendUser.nickname,
    avatar: backendUser.avatar,
    role: backendUser.role === 'admin' ? 'admin' : 'customer',
    status: backendUser.status,
  })

  // 登录
  const loginUser = async (credentials: LoginRequest) => {
    const response = await login(credentials)
    token.value = response.data.token
    user.value = mapUser(response.data.user)
    isAuthenticated.value = true
    localStorage.setItem('token', response.data.token)
    return response
  }

  // 登出
  const logoutUser = async () => {
    token.value = null
    user.value = null
    isAuthenticated.value = false
    localStorage.removeItem('token')
  }

  // 获取用户信息
  const fetchProfile = async () => {
    const response = await getCurrentUser()
    user.value = mapUser(response.data)
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
