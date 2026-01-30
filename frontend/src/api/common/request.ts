/**
 * 通用请求配置
 * 基于 axios 封装的请求工具
 */
import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
  type AxiosError,
} from 'axios'
import { ElMessage } from 'element-plus'
import { mockMatcher } from '@/mock'

// 请求配置
const config: AxiosRequestConfig = {
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
}

const enableDebugLog = import.meta.env.VITE_DEBUG_LOG === 'true'

// 创建 axios 实例
const service: AxiosInstance = axios.create(config)

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    if (enableDebugLog) {
      const metadata = (config as any).metadata || {}
      metadata.startTime = Date.now()
      ;(config as any).metadata = metadata
      console.debug('[API Request]', {
        method: config.method,
        url: config.url,
        params: config.params,
        data: config.data,
      })
    }

    // Mock 数据拦截
    const useMock = import.meta.env.VITE_USE_MOCK === 'true'
    if (useMock) {
      const mockData = mockMatcher({
        url: config.url || '',
        method: config.method || 'get',
        data: config.data,
      })

      if (mockData) {
        // 返回 mock 数据，取消真实请求
        return Promise.reject({
          config,
          response: {
            data: mockData,
            status: mockData.code || 200,
            statusText: 'OK',
            headers: {},
            config,
          },
          isMock: true,
        })
      }
    }

    // 从 localStorage 获取 token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error: AxiosError) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data

    if (enableDebugLog) {
      const metadata = (response.config as any)?.metadata
      const duration = metadata?.startTime ? Date.now() - metadata.startTime : undefined
      console.debug('[API Response]', {
        url: response.config.url,
        code: res?.code,
        duration,
        data: res?.data,
      })
    }

    if (res && typeof res.code !== 'undefined' && res.code !== 0) {
      const message = res.msg || res.message || '请求失败'

      if (res.code === 401) {
        ElMessage.error('未授权，请重新登录')
        localStorage.removeItem('token')
        window.location.href = '/login'
      } else if (res.code === 403) {
        ElMessage.error('拒绝访问')
      } else {
        ElMessage.error(message)
      }

      return Promise.reject(res)
    }

    return res
  },
  (error: AxiosError | any) => {
    // 处理 Mock 数据
    if (error.isMock && error.response) {
      const mockResponse = error.response.data

      // Mock 数据错误处理（如登录失败）
      if (mockResponse.code !== 200) {
        ElMessage.error(mockResponse.msg || mockResponse.message || '请求失败')
        return Promise.reject(mockResponse)
      }

      // 返回 Mock 数据
      return mockResponse
    }

    console.error('响应错误:', error)
    if (enableDebugLog) {
      console.debug('[API Error]', {
        url: error.config?.url,
        message: error.message,
        response: error.response?.data,
      })
    }

    if (error.response) {
      const status = error.response.status

      switch (status) {
        case 401:
          ElMessage.error('未授权，请重新登录')
          // 清除 token 并跳转到登录页
          localStorage.removeItem('token')
          window.location.href = '/login'
          break
        case 403:
          ElMessage.error('拒绝访问')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 429:
          ElMessage.error('请求过于频繁，请稍后再试')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error((error.response.data as any)?.msg || (error.response.data as any)?.message || '请求失败')
      }
    } else if (error.request) {
      ElMessage.error('网络错误，请检查网络连接')
    } else {
      ElMessage.error('请求配置错误')
    }

    return Promise.reject(error)
  }
)

export default service
