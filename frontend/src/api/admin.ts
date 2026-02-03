import request from '@/utils/request'
import type { ApiResponse, PageRequest, PageResponse } from '@/types/common'
import type { Machine } from '@/types/machine'
import type { Customer, CustomerDetail, CustomerForm } from '@/types/customer'
import type { AllocationRecord, QuickAllocateForm, ExtendAllocationForm } from '@/types/allocation'

// ==================== 机器管理 ====================

/**
 * 获取机器列表
 */
export function getMachineList(params: PageRequest): Promise<ApiResponse<PageResponse<Machine>>> {
  return request.get('/api/admin/machines', { params })
}

/**
 * 获取机器详情
 */
export function getMachineDetail(id: number): Promise<ApiResponse<Machine>> {
  return request.get(`/api/admin/machines/${id}`)
}

/**
 * 添加机器
 */
export function addMachine(data: Partial<Machine>): Promise<ApiResponse<Machine>> {
  return request.post('/api/admin/machines', data)
}

/**
 * 更新机器信息
 */
export function updateMachine(id: number, data: Partial<Machine>): Promise<ApiResponse<Machine>> {
  return request.put(`/api/admin/machines/${id}`, data)
}

/**
 * 删除机器
 */
export function deleteMachine(id: number): Promise<ApiResponse<void>> {
  return request.delete(`/api/admin/machines/${id}`)
}

/**
 * 批量导入机器
 */
export function batchImportMachines(data: Partial<Machine>[]): Promise<ApiResponse<{ success: number; failed: number }>> {
  return request.post('/api/admin/machines/import', data)
}

/**
 * 设置机器维护状态
 */
export function setMachineMaintenance(id: number, maintenance: boolean): Promise<ApiResponse<void>> {
  return request.post(`/api/admin/machines/${id}/maintenance`, { maintenance })
}

// ==================== 客户管理 ====================

/**
 * 获取客户列表
 */
export function getCustomerList(params: PageRequest): Promise<ApiResponse<PageResponse<Customer>>> {
  return request.get('/api/admin/customers', { params })
}

/**
 * 获取客户详情
 */
export function getCustomerDetail(id: number): Promise<ApiResponse<CustomerDetail>> {
  return request.get(`/api/admin/customers/${id}`)
}

/**
 * 添加客户
 */
export function addCustomer(data: CustomerForm): Promise<ApiResponse<Customer>> {
  return request.post('/api/admin/customers', data)
}

/**
 * 更新客户信息
 */
export function updateCustomer(id: number, data: CustomerForm): Promise<ApiResponse<Customer>> {
  return request.put(`/api/admin/customers/${id}`, data)
}

/**
 * 删除客户
 */
export function deleteCustomer(id: number): Promise<ApiResponse<void>> {
  return request.delete(`/api/admin/customers/${id}`)
}

/**
 * 启用/禁用客户
 */
export function toggleCustomerStatus(id: number, enabled: boolean): Promise<ApiResponse<void>> {
  const action = enabled ? 'enable' : 'disable'
  return request.post(`/api/admin/customers/${id}/${action}`)
}

// ==================== 分配管理 ====================

/**
 * 获取分配记录列表
 */
export function getAllocationList(params: PageRequest): Promise<ApiResponse<PageResponse<AllocationRecord>>> {
  return request.get('/api/admin/allocations', { params })
}

/**
 * 获取分配记录详情
 */
export function getAllocationDetail(id: number): Promise<ApiResponse<AllocationRecord>> {
  return request.get(`/api/admin/allocations/${id}`)
}

/**
 * 快速分配机器
 */
export function quickAllocate(data: QuickAllocateForm): Promise<ApiResponse<AllocationRecord>> {
  return request.post('/api/admin/allocations', data)
}

/**
 * 延期分配
 */
export function extendAllocation(id: number, data: ExtendAllocationForm): Promise<ApiResponse<AllocationRecord>> {
  return request.post(`/api/admin/allocations/${id}/extend`, data)
}

/**
 * 回收机器
 */
export function reclaimMachine(id: number): Promise<ApiResponse<void>> {
  return request.post(`/api/admin/allocations/${id}/reclaim`)
}

/**
 * 获取可分配机器列表
 */
export function getAvailableMachines(): Promise<ApiResponse<Machine[]>> {
  return request.get('/api/admin/allocations/available-machines')
}

// ==================== 监控中心 ====================

/**
 * 获取实时监控数据
 */
export function getRealtimeMonitoring(params?: { machineIds?: number[] }): Promise<ApiResponse<any>> {
  return request.get('/api/admin/monitoring/realtime', { params })
}

/**
 * 获取告警列表
 */
export function getAlertList(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/admin/alerts', { params })
}

/**
 * 处理告警
 */
export function handleAlert(id: number, action: 'acknowledge' | 'resolve'): Promise<ApiResponse<void>> {
  const path = action === 'acknowledge' ? 'ack' : 'resolve'
  return request.post(`/api/admin/alerts/${id}/${path}`)
}

// ==================== 镜像管理 ====================

/**
 * 获取镜像列表
 */
export function getImageList(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/admin/images', { params })
}

/**
 * 上传镜像
 */
export function uploadImage(data: FormData): Promise<ApiResponse<any>> {
  return request.post('/api/admin/images/upload', data, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

/**
 * 删除镜像
 */
export function deleteImage(id: number): Promise<ApiResponse<void>> {
  return request.delete(`/api/admin/images/${id}`)
}

// ==================== 数据集管理 ====================

/**
 * 获取数据集列表
 */
export function getDatasetList(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/admin/datasets', { params })
}

/**
 * 上传数据集
 */
export function uploadDataset(data: FormData): Promise<ApiResponse<any>> {
  return request.post('/api/admin/datasets/upload', data, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

/**
 * 删除数据集
 */
export function deleteDataset(id: number): Promise<ApiResponse<void>> {
  return request.delete(`/api/admin/datasets/${id}`)
}

// ==================== 任务管理 ====================

/**
 * 获取任务列表
 */
export function getTaskList(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/admin/tasks', { params })
}

/**
 * 获取任务详情
 */
export function getTaskDetail(id: number): Promise<ApiResponse<any>> {
  return request.get(`/api/admin/tasks/${id}`)
}

/**
 * 终止任务
 */
export function terminateTask(id: number): Promise<ApiResponse<void>> {
  return request.post(`/api/admin/tasks/${id}/stop`)
}

// ==================== 数据统计 ====================

/**
 * 获取资源统计数据
 */
export function getResourceStats(params?: { startDate?: string; endDate?: string }): Promise<ApiResponse<any>> {
  return request.get('/api/admin/statistics/resources', { params })
}

/**
 * 获取客户统计数据
 */
export function getCustomerStats(params?: { startDate?: string; endDate?: string }): Promise<ApiResponse<any>> {
  return request.get('/api/admin/statistics/customers', { params })
}

/**
 * 获取Dashboard概览数据
 */
export function getDashboardOverview(): Promise<ApiResponse<any>> {
  return request.get('/api/admin/dashboard/stats')
}
