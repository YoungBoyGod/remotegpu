import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { nextTick } from 'vue'
import ElementPlus from 'element-plus'
import MachineListView from '../views/admin/MachineListView.vue'

// Mock vue-router
const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush }),
  useRoute: () => ({ query: {} }),
}))

// Mock admin API
const mockGetMachineList = vi.fn()
const mockDeleteMachine = vi.fn()
const mockSetMachineMaintenance = vi.fn()
const mockCollectMachineSpec = vi.fn()
const mockAllocateMachine = vi.fn()
const mockReclaimMachine = vi.fn()
const mockGetCustomerList = vi.fn()
const mockBatchImportMachines = vi.fn()

vi.mock('@/api/admin', () => ({
  getMachineList: (...args: any[]) => mockGetMachineList(...args),
  deleteMachine: (...args: any[]) => mockDeleteMachine(...args),
  setMachineMaintenance: (...args: any[]) => mockSetMachineMaintenance(...args),
  collectMachineSpec: (...args: any[]) => mockCollectMachineSpec(...args),
  allocateMachine: (...args: any[]) => mockAllocateMachine(...args),
  reclaimMachine: (...args: any[]) => mockReclaimMachine(...args),
  getCustomerList: (...args: any[]) => mockGetCustomerList(...args),
  batchImportMachines: (...args: any[]) => mockBatchImportMachines(...args),
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

const sampleMachines = [
  {
    id: 'host-001',
    name: 'GPU-Node-1',
    region: 'cn-east',
    status: 'idle',
    device_status: 'online',
    allocation_status: 'idle',
    needs_collect: false,
    ip_address: '192.168.1.10',
    public_ip: '1.2.3.4',
    ssh_port: 22,
    ssh_host: '1.2.3.4',
    ssh_username: 'root',
    total_cpu: 64,
    total_memory_gb: 256,
    gpus: [{ id: 1, host_id: 'host-001', index: 0, uuid: 'gpu-0', name: 'A100', memory_total_mb: 81920 }],
  },
  {
    id: 'host-002',
    name: 'GPU-Node-2',
    region: 'cn-west',
    status: 'allocated',
    device_status: 'offline',
    allocation_status: 'allocated',
    needs_collect: true,
    ip_address: '192.168.1.11',
    total_cpu: 32,
    total_memory_gb: 128,
    gpus: [],
  },
]

describe('MachineListView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    setActivePinia(createPinia())
    mockGetMachineList.mockResolvedValue({
      data: { list: sampleMachines, total: 2 },
    })
  })

  const mountView = () => {
    return mount(MachineListView, {
      global: {
        plugins: [ElementPlus],
        stubs: {
          DataTable: {
            template: '<div class="data-table-stub"><slot /></div>',
            props: ['data', 'total', 'loading', 'currentPage', 'pageSize', 'showPagination'],
          },
          ElTableColumn: true,
          ElDialog: true,
        },
      },
    })
  }

  // ==================== 页面渲染测试 ====================

  it('渲染页面标题', () => {
    const wrapper = mountView()
    expect(wrapper.text()).toContain('机器列表')
  })

  it('渲染添加机器和批量导入按钮', () => {
    const wrapper = mountView()
    expect(wrapper.text()).toContain('添加机器')
    expect(wrapper.text()).toContain('批量导入')
  })

  it('渲染筛选栏', () => {
    const wrapper = mountView()
    expect(wrapper.text()).toContain('状态')
    expect(wrapper.text()).toContain('区域')
    expect(wrapper.text()).toContain('GPU型号')
    expect(wrapper.text()).toContain('关键词')
    expect(wrapper.text()).toContain('搜索')
    expect(wrapper.text()).toContain('重置')
  })

  it('挂载时调用 getMachineList 加载数据', async () => {
    mountView()
    await nextTick()
    expect(mockGetMachineList).toHaveBeenCalledTimes(1)
  })

  // ==================== 导航测试 ====================

  it('点击添加机器按钮跳转到添加页面', async () => {
    const wrapper = mountView()
    const addBtn = wrapper.findAll('button').find(b => b.text().includes('添加机器'))
    if (addBtn) {
      await addBtn.trigger('click')
      expect(mockPush).toHaveBeenCalledWith('/admin/machines/add')
    }
  })

  // ==================== API 调用测试 ====================

  it('加载失败时不抛出异常', async () => {
    mockGetMachineList.mockRejectedValue(new Error('network error'))
    expect(() => mountView()).not.toThrow()
  })

  it('getMachineList 传递正确的分页参数', async () => {
    mountView()
    await nextTick()
    const callArgs = mockGetMachineList.mock.calls[0]![0]
    expect(callArgs.page).toBe(1)
    expect(callArgs.pageSize).toBe(10)
  })
})
