import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

// Token 刷新并发控制
let isRefreshing = false
let refreshSubscribers: Array<(token: string) => void> = []

const onRefreshed = (token: string) => {
  refreshSubscribers.forEach((cb) => cb(token))
  refreshSubscribers = []
}

const addRefreshSubscriber = (cb: (token: string) => void) => {
  refreshSubscribers.push(cb)
}

/**
 * 创建axios实例
 */
const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json;charset=UTF-8',
  },
})

/**
 * 请求拦截器
 */
service.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()

    // 添加访问令牌
    if (authStore.accessToken) {
      config.headers.Authorization = `Bearer ${authStore.accessToken}`
    }

    // 添加租户ID (如果存在)
    if (authStore.user?.tenantId) {
      config.headers['X-Tenant-Id'] = authStore.user.tenantId.toString()
    }

    return config
  },
  (error) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

/**
 * 响应拦截器
 */
service.interceptors.response.use(
  (response: AxiosResponse) => {
    const { data } = response

    // 如果返回的是Blob类型(文件下载),直接返回
    if (response.config.responseType === 'blob') {
      return data
    }

    // 检查业务状态码
    if (data.code !== undefined && data.code !== 0 && data.code !== 200) {
      // 传递完整的 data 对象，让调用方决定如何处理错误
      return Promise.reject(data)
    }

    return data
  },
  async (error) => {
    const authStore = useAuthStore()

    if (error.response) {
      const { status, data } = error.response

      switch (status) {
        case 401:
          // 未授权,尝试刷新令牌
          if (authStore.refreshToken && !error.config._retry) {
            error.config._retry = true

            if (isRefreshing) {
              // 已有刷新请求在进行，排队等待
              return new Promise((resolve) => {
                addRefreshSubscriber((token: string) => {
                  error.config.headers.Authorization = `Bearer ${token}`
                  resolve(service(error.config))
                })
              })
            }

            isRefreshing = true
            try {
              await authStore.refreshAccessToken()
              const newToken = authStore.accessToken!
              onRefreshed(newToken)
              // 重试原请求
              return service(error.config)
            } catch (refreshError) {
              refreshSubscribers = []
              await authStore.logout()
              router.push('/login')
              ElMessage.error('登录已过期,请重新登录')
              return Promise.reject(refreshError)
            } finally {
              isRefreshing = false
            }
          } else {
            await authStore.logout()
            router.push('/login')
            ElMessage.error('登录已过期,请重新登录')
          }
          break

        case 403:
          ElMessage.error('没有权限访问该资源')
          break

        case 404:
          ElMessage.error('请求的资源不存在')
          break

        case 500:
          ElMessage.error(data?.message || '服务器内部错误')
          break

        case 502:
          ElMessage.error('网关错误')
          break

        case 503:
          ElMessage.error('服务暂时不可用')
          break

        default:
          ElMessage.error(data?.message || `请求失败 (${status})`)
      }
    } else if (error.request) {
      // 请求已发送但没有收到响应
      ElMessage.error('网络连接失败,请检查网络设置')
    } else {
      // 请求配置出错
      ElMessage.error('请求配置错误')
    }

    return Promise.reject(error)
  }
)

export default service
