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
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
}

// 创建 axios 实例
const service: AxiosInstance = axios.create(config)

// 请求拦截器
service.interceptors.request.use(
  (config) => {
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

    // 如果返回的状态码不是 200，则认为是错误
    if (response.status !== 200) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }

    return res
  },
  (error: AxiosError | any) => {
    // 处理 Mock 数据
    if (error.isMock && error.response) {
      const mockResponse = error.response.data

      // Mock 数据错误处理（如登录失败）
      if (mockResponse.code !== 200) {
        ElMessage.error(mockResponse.message || '请求失败')
        return Promise.reject(mockResponse)
      }

      // 返回 Mock 数据
      return mockResponse
    }

    console.error('响应错误:', error)

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
          ElMessage.error((error.response.data as any)?.message || '请求失败')
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
