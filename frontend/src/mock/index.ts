/**
 * Mock 数据入口
 * 根据请求 URL 和方法返回对应的 mock 数据
 */
import { mockLogin, mockGetCurrentUser, mockLogout } from './auth'

export interface MockRequest {
  url: string
  method: string
  data?: any
}

/**
 * Mock 数据匹配器
 */
export const mockMatcher = (request: MockRequest) => {
  const { url, method, data } = request

  // 认证相关接口
  if (url === '/auth/login' && method === 'post') {
    return mockLogin(data)
  }

  if (url === '/users/me' && method === 'get') {
    return mockGetCurrentUser()
  }

  if (url === '/auth/logout' && method === 'post') {
    return mockLogout()
  }

  // 未匹配到 mock 数据，返回 null
  return null
}
