import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { nextTick } from 'vue'
import ElementPlus from 'element-plus'
import CustomerListView from '../views/admin/CustomerListView.vue'

// Mock vue-router
const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush }),
  useRoute: () => ({ query: {} }),
}))

// Mock admin API
const mockGetCustomerList = vi.fn()
const mockDisableCustomer = vi.fn()
const mockEnableCustomer = vi.fn()

vi.mock('@/api/admin', () => ({
  getCustomerList: (...args: any[]) => mockGetCustomerList(...args),
  disableCustomer: (...args: any[]) => mockDisableCustomer(...args),
  enableCustomer: (...args: any[]) => mockEnableCustomer(...args),
}))

// Mock ElMessage / ElMessageBox
vi.mock('element-plus', async (importOriginal) => {
  const mod = await importOriginal<typeof import('element-plus')>()
  return {
    ...mod,
    ElMessage: { success: vi.fn(), error: vi.fn(), warning: vi.fn() },
    ElMessageBox: { confirm: vi.fn().mockResolvedValue('confirm') },
  }
})

const sampleCustomers = [
  {
    id: 1,
    username: 'customer1',
    company_code: 'ACME',
    company: 'ACME Corp',
    display_name: 'Alice',
    full_name: 'Alice Wang',
    email: 'alice@acme.com',
    phone: '13800001111',
    status: 'active',
    created_at: '2026-01-15T10:00:00Z',
  },
  {
    id: 2,
    username: 'customer2',
    company_code: 'BETA',
    company: 'Beta Inc',
    display_name: 'Bob',
    full_name: 'Bob Li',
    email: 'bob@beta.com',
    phone: '13800002222',
    status: 'suspended',
    created_at: '2026-01-20T12:00:00Z',
  },
  {
    id: 3,
    username: 'customer3',
    company_code: 'GAMMA',
    company: 'Gamma Ltd',
    email: 'charlie@gamma.com',
    status: 'disabled',
    created_at: '2026-02-01T08:00:00Z',
  },
]

describe('CustomerListView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    setActivePinia(createPinia())
    mockGetCustomerList.mockResolvedValue({
      data: { list: sampleCustomers, total: 3 },
    })
  })

  const mountView = () => {
    return mount(CustomerListView, {
      global: {
        plugins: [ElementPlus],
        stubs: {
          DataTable: {
            template: '<div class="data-table-stub"><slot /></div>',
            props: ['data', 'total', 'loading', 'currentPage', 'pageSize', 'showPagination'],
          },
          ElTableColumn: true,
        },
      },
    })
  }

  // ==================== 页面渲染测试 ====================

  it('渲染页面标题', () => {
    const wrapper = mountView()
    expect(wrapper.text()).toContain('客户列表')
  })

  it('渲染添加客户按钮', () => {
    const wrapper = mountView()
    expect(wrapper.text()).toContain('添加客户')
  })

  it('渲染筛选栏', () => {
    const wrapper = mountView()
    expect(wrapper.text()).toContain('状态')
    expect(wrapper.text()).toContain('关键词')
    expect(wrapper.text()).toContain('搜索')
    expect(wrapper.text()).toContain('重置')
  })

  it('挂载时调用 getCustomerList 加载数据', async () => {
    mountView()
    await nextTick()
    expect(mockGetCustomerList).toHaveBeenCalledTimes(1)
  })

  // ==================== 导航测试 ====================

  it('点击添加客户按钮跳转到添加页面', async () => {
    const wrapper = mountView()
    const addBtn = wrapper.findAll('button').find(b => b.text().includes('添加客户'))
    if (addBtn) {
      await addBtn.trigger('click')
      expect(mockPush).toHaveBeenCalledWith('/admin/customers/add')
    }
  })

  // ==================== API 调用测试 ====================

  it('加载失败时不抛出异常', async () => {
    mockGetCustomerList.mockRejectedValue(new Error('network error'))
    expect(() => mountView()).not.toThrow()
  })

  it('getCustomerList 传递正确的分页参数', async () => {
    mountView()
    await nextTick()
    const callArgs = mockGetCustomerList.mock.calls[0]![0]
    expect(callArgs.page).toBe(1)
    expect(callArgs.pageSize).toBe(10)
  })

  // ==================== 筛选状态选项测试 ====================

  it('筛选栏包含状态下拉框和关键词输入框', () => {
    const wrapper = mountView()
    const filterCard = wrapper.find('.filter-card')
    expect(filterCard.exists()).toBe(true)
    // 验证有 el-select 组件（状态筛选）
    expect(filterCard.findComponent({ name: 'ElSelect' }).exists()).toBe(true)
    // 验证有 el-input 组件（关键词输入）
    expect(filterCard.findComponent({ name: 'ElInput' }).exists()).toBe(true)
  })
})
