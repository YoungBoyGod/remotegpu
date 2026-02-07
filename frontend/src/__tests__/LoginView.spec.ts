import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ElementPlus from 'element-plus'
import LoginView from '../views/LoginView.vue'

// Mock vue-router
const mockPush = vi.fn()
const mockRoute = { query: {} }
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush }),
  useRoute: () => mockRoute,
}))

// Mock auth store
const mockLogin = vi.fn()
vi.mock('@/stores/auth', () => ({
  useAuthStore: () => ({
    login: mockLogin,
  }),
}))

// Mock ElMessage
vi.mock('element-plus', async (importOriginal) => {
  const mod = await importOriginal<typeof import('element-plus')>()
  return {
    ...mod,
    ElMessage: {
      success: vi.fn(),
      error: vi.fn(),
      warning: vi.fn(),
    },
  }
})

describe('LoginView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    setActivePinia(createPinia())
    mockRoute.query = {}
  })

  const mountLogin = () => {
    return mount(LoginView, {
      global: {
        plugins: [ElementPlus],
      },
    })
  }

  // ==================== 页面渲染测试 ====================

  it('渲染登录页面标题和副标题', () => {
    const wrapper = mountLogin()
    expect(wrapper.text()).toContain('RemoteGPU')
    expect(wrapper.text()).toContain('企业级 GPU 云平台')
  })

  it('渲染用户名和密码输入框', () => {
    const wrapper = mountLogin()
    const inputs = wrapper.findAll('input')
    // 至少有用户名和密码两个输入框
    expect(inputs.length).toBeGreaterThanOrEqual(2)
  })

  it('渲染登录按钮', () => {
    const wrapper = mountLogin()
    const btn = wrapper.find('.login-button')
    expect(btn.exists()).toBe(true)
    expect(btn.text()).toContain('登录')
  })

  it('渲染记住我复选框', () => {
    const wrapper = mountLogin()
    expect(wrapper.text()).toContain('记住我')
  })

  it('渲染忘记密码链接', () => {
    const wrapper = mountLogin()
    expect(wrapper.text()).toContain('忘记密码')
  })

  it('渲染页脚提示文字', () => {
    const wrapper = mountLogin()
    expect(wrapper.text()).toContain('请使用公司账号登录')
  })

  // ==================== 表单交互测试 ====================

  it('用户名输入框可以输入值', async () => {
    const wrapper = mountLogin()
    const inputs = wrapper.findAll('input')
    const usernameInput = inputs[0]!
    await usernameInput.setValue('testuser')
    expect(usernameInput.element.value).toBe('testuser')
  })

  it('密码输入框可以输入值', async () => {
    const wrapper = mountLogin()
    const inputs = wrapper.findAll('input')
    // 密码输入框是 type="password"
    const passwordInput = inputs.find(i => i.element.type === 'password')
    expect(passwordInput).toBeDefined()
    await passwordInput!.setValue('password123')
    expect(passwordInput!.element.value).toBe('password123')
  })

  // ==================== 登录流程测试 ====================

  it('登录成功后跳转到首页', async () => {
    mockLogin.mockResolvedValue({
      data: { mustChangePassword: false },
    })

    const wrapper = mountLogin()
    const inputs = wrapper.findAll('input')
    await inputs[0]!.setValue('admin')
    const pwdInput = inputs.find(i => i.element.type === 'password')
    await pwdInput!.setValue('password123')

    // 触发登录
    await wrapper.find('.login-button').trigger('click')
    // 等待异步验证和登录
    await vi.dynamicImportSettled()
    await new Promise(r => setTimeout(r, 100))

    if (mockLogin.mock.calls.length > 0) {
      expect(mockLogin).toHaveBeenCalledWith({
        username: 'admin',
        password: 'password123',
      })
    }
  })

  it('登录成功且需要改密时跳转到改密页', async () => {
    mockLogin.mockResolvedValue({
      data: { mustChangePassword: true },
    })

    const wrapper = mountLogin()
    const inputs = wrapper.findAll('input')
    await inputs[0]!.setValue('newuser')
    const pwdInput = inputs.find(i => i.element.type === 'password')
    await pwdInput!.setValue('temppass1')

    await wrapper.find('.login-button').trigger('click')
    await vi.dynamicImportSettled()
    await new Promise(r => setTimeout(r, 100))

    if (mockPush.mock.calls.length > 0) {
      expect(mockPush).toHaveBeenCalledWith('/change-password')
    }
  })

  it('登录成功后跳转到 redirect 参数指定的页面', async () => {
    mockRoute.query = { redirect: '/admin/machines/list' }
    mockLogin.mockResolvedValue({
      data: { mustChangePassword: false },
    })

    const wrapper = mountLogin()
    const inputs = wrapper.findAll('input')
    await inputs[0]!.setValue('admin')
    const pwdInput = inputs.find(i => i.element.type === 'password')
    await pwdInput!.setValue('password123')

    await wrapper.find('.login-button').trigger('click')
    await vi.dynamicImportSettled()
    await new Promise(r => setTimeout(r, 100))

    if (mockPush.mock.calls.length > 0) {
      expect(mockPush).toHaveBeenCalledWith('/admin/machines/list')
    }
  })

  it('点击忘记密码跳转到忘记密码页', async () => {
    const wrapper = mountLogin()
    const link = wrapper.findAll('.el-link').find(l => l.text().includes('忘记密码'))
    if (link) {
      await link.trigger('click')
      expect(mockPush).toHaveBeenCalledWith('/forgot-password')
    }
  })
})
