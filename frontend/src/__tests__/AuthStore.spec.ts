import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {}
  return {
    getItem: vi.fn((key: string) => store[key] ?? null),
    setItem: vi.fn((key: string, value: string) => { store[key] = value }),
    removeItem: vi.fn((key: string) => { delete store[key] }),
    clear: vi.fn(() => { store = {} }),
  }
})()
Object.defineProperty(globalThis, 'localStorage', { value: localStorageMock })

// Mock API
const mockLoginApi = vi.fn()
const mockLogoutApi = vi.fn()
const mockRefreshTokenApi = vi.fn()
const mockGetCurrentUser = vi.fn()

vi.mock('@/api/auth', () => ({
  login: (...args: any[]) => mockLoginApi(...args),
  logout: (...args: any[]) => mockLogoutApi(...args),
  refreshToken: (...args: any[]) => mockRefreshTokenApi(...args),
  getCurrentUser: (...args: any[]) => mockGetCurrentUser(...args),
}))

// 需要在 mock 之后动态导入
const { useAuthStore } = await import('@/stores/auth')

describe('useAuthStore', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.clear()
    setActivePinia(createPinia())
  })

  // ==================== 初始状态测试 ====================

  it('初始状态未认证', () => {
    const store = useAuthStore()
    expect(store.isAuthenticated).toBe(false)
    expect(store.user).toBeNull()
    expect(store.accessToken).toBeNull()
  })

  it('isAdmin 在无用户时为 false', () => {
    const store = useAuthStore()
    expect(store.isAdmin).toBe(false)
  })

  it('isCustomer 在无用户时为 false', () => {
    const store = useAuthStore()
    expect(store.isCustomer).toBe(false)
  })

  // ==================== 登录测试 ====================

  it('登录成功后保存 token 到 localStorage', async () => {
    mockLoginApi.mockResolvedValue({
      data: {
        accessToken: 'access-123',
        refreshToken: 'refresh-456',
        mustChangePassword: false,
      },
    })
    mockGetCurrentUser.mockResolvedValue({
      data: { id: 1, username: 'admin', role: 'admin', email: 'a@b.com' },
    })

    const store = useAuthStore()
    await store.login({ username: 'admin', password: 'pass' })

    expect(store.accessToken).toBe('access-123')
    expect(store.refreshToken).toBe('refresh-456')
    expect(store.isAuthenticated).toBe(true)
    expect(localStorageMock.setItem).toHaveBeenCalledWith('accessToken', 'access-123')
    expect(localStorageMock.setItem).toHaveBeenCalledWith('refreshToken', 'refresh-456')
  })

  it('登录成功后获取用户信息', async () => {
    mockLoginApi.mockResolvedValue({
      data: { accessToken: 'tk', refreshToken: 'rt', mustChangePassword: false },
    })
    mockGetCurrentUser.mockResolvedValue({
      data: { id: 1, username: 'admin', role: 'admin', email: 'a@b.com' },
    })

    const store = useAuthStore()
    await store.login({ username: 'admin', password: 'pass' })

    expect(store.user).not.toBeNull()
    expect(store.user?.username).toBe('admin')
    expect(store.isAdmin).toBe(true)
  })

  it('登录成功但 fetchProfile 失败且需改密时设置最小用户信息', async () => {
    mockLoginApi.mockResolvedValue({
      data: { accessToken: 'tk', refreshToken: 'rt', mustChangePassword: true },
    })
    mockGetCurrentUser.mockRejectedValue(new Error('forbidden'))

    const store = useAuthStore()
    await store.login({ username: 'new', password: 'temp' })

    expect(store.user).not.toBeNull()
    expect(store.user?.must_change_password).toBe(true)
  })

  it('登录 API 失败时抛出错误', async () => {
    mockLoginApi.mockRejectedValue({ code: 2003, msg: '密码错误' })

    const store = useAuthStore()
    await expect(store.login({ username: 'admin', password: 'wrong' })).rejects.toEqual({
      code: 2003,
      msg: '密码错误',
    })
  })

  // ==================== 登出测试 ====================

  it('登出后清除所有状态和 localStorage', async () => {
    mockLogoutApi.mockResolvedValue({})

    const store = useAuthStore()
    store.accessToken = 'old-token'
    store.refreshToken = 'old-refresh'
    store.user = { id: 1, username: 'admin', role: 'admin', email: '' }

    await store.logout()

    expect(store.accessToken).toBeNull()
    expect(store.refreshToken).toBeNull()
    expect(store.user).toBeNull()
    expect(store.isAuthenticated).toBe(false)
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('accessToken')
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('refreshToken')
  })

  it('登出 API 失败时仍然清除本地状态', async () => {
    mockLogoutApi.mockRejectedValue(new Error('network'))

    const store = useAuthStore()
    store.accessToken = 'token'

    await store.logout()

    expect(store.accessToken).toBeNull()
    expect(store.isAuthenticated).toBe(false)
  })

  // ==================== 角色判断测试 ====================

  it('admin 角色 isAdmin 为 true', () => {
    const store = useAuthStore()
    store.user = { id: 1, username: 'admin', role: 'admin', email: '' }
    expect(store.isAdmin).toBe(true)
    expect(store.isCustomer).toBe(false)
  })

  it('customer 角色 isCustomer 为 true', () => {
    const store = useAuthStore()
    store.user = { id: 2, username: 'user', role: 'customer', email: '' }
    expect(store.isCustomer).toBe(true)
    expect(store.isAdmin).toBe(false)
  })

  it('customer_owner 角色 isCustomer 为 true', () => {
    const store = useAuthStore()
    store.user = { id: 3, username: 'owner', role: 'customer_owner', email: '' }
    expect(store.isCustomer).toBe(true)
  })

  // ==================== Token 刷新测试 ====================

  it('刷新 token 成功后更新 localStorage', async () => {
    mockRefreshTokenApi.mockResolvedValue({
      data: { accessToken: 'new-access', refreshToken: 'new-refresh' },
    })

    const store = useAuthStore()
    store.refreshToken = 'old-refresh'

    await store.refreshAccessToken()

    expect(store.accessToken).toBe('new-access')
    expect(store.refreshToken).toBe('new-refresh')
    expect(localStorageMock.setItem).toHaveBeenCalledWith('accessToken', 'new-access')
  })

  it('无 refreshToken 时刷新抛出错误', async () => {
    const store = useAuthStore()
    store.refreshToken = null

    await expect(store.refreshAccessToken()).rejects.toThrow('没有刷新令牌')
  })
})
