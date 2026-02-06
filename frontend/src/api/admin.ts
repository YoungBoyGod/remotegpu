import request from '@/utils/request'
import type { ApiResponse, PageRequest, PageResponse } from '@/types/common'
import type { Machine } from '@/types/machine'
import type { Customer, CustomerDetail, CustomerForm } from '@/types/customer'
import type { AllocationRecord, QuickAllocateForm } from '@/types/allocation'
import type { Task, TaskLogResponse, TaskResultResponse } from '@/types/task'
import type { SystemConfig, UpdateSystemConfigsPayload } from '@/types/systemConfig'

// ==================== 机器管理 ====================

/**
 * 获取机器列表
 */
export function getMachineList(
  params: PageRequest & { status?: string; region?: string; gpu_model?: string; filters?: Record<string, any> }
): Promise<ApiResponse<PageResponse<Machine>>> {
  const query: Record<string, any> = {
    page: params.page,
    page_size: params.pageSize,
  }
  const source = params.filters || {}
  const status = params.status ?? source.status
  const region = params.region ?? source.region
  const gpuModel = params.gpu_model ?? source.gpu_model

  if (status) query.status = status
  if (region) query.region = region
  if (gpuModel) query.gpu_model = gpuModel

  return request.get('/admin/machines', { params: query })
}

/**
 * 获取机器详情
 */
export function getMachineDetail(id: string): Promise<ApiResponse<Machine>> {
  return request.get(`/admin/machines/${id}`)
}

/**
 * 添加机器
 */
export interface CreateMachinePayload {
  name: string
  hostname?: string
  region?: string
  ip_address: string
  public_ip?: string
  ssh_port: number
  ssh_username?: string
  ssh_password?: string
  ssh_key?: string
}

export function addMachine(data: CreateMachinePayload): Promise<ApiResponse<Machine>> {
  return request.post('/admin/machines', data)
}

/**
 * 批量导入机器
 */
export function batchImportMachines(data: { machines: Partial<Machine>[] }): Promise<ApiResponse<{ message: string; count: number }>> {
  return request.post('/admin/machines/import', data)
}

/**
 * 分配机器
 */
export function allocateMachine(id: string, data: { customer_id: number; duration_months: number; remark?: string }): Promise<ApiResponse<AllocationRecord>> {
  return request.post(`/admin/machines/${id}/allocate`, { ...data, host_id: id })
}

/**
 * 分配机器（扩展版，支持时间段、对接人、通知方式）
 */
export interface AssignMachinePayload {
  customer_id: number
  duration_months: number
  start_time?: string
  end_time?: string
  contact_person?: string
  notify_methods?: string[]
  remark?: string
}

export function assignMachine(id: string, data: AssignMachinePayload): Promise<ApiResponse<AllocationRecord>> {
  return request.post(`/admin/machines/${id}/allocate`, { ...data, host_id: id })
}

/**
 * 回收机器
 */
export function reclaimMachine(id: string, data?: { reason?: string; force?: boolean }): Promise<ApiResponse<void>> {
  return request.post(`/admin/machines/${id}/reclaim`, data)
}

/**
 * 删除机器
 */
export function deleteMachine(id: string): Promise<ApiResponse<void>> {
  return request.delete(`/admin/machines/${id}`)
}

/**
 * 设置机器维护状态
 */
export function setMachineMaintenance(id: string, maintenance: boolean): Promise<ApiResponse<void>> {
  return request.post(`/admin/machines/${id}/maintenance`, { maintenance })
}

/**
 * 触发机器硬件补采
 */
export function collectMachineSpec(id: string): Promise<ApiResponse<Machine>> {
  return request.post(`/admin/machines/${id}/collect`)
}

// ==================== 客户管理 ====================

/**
 * 获取客户列表
 */
export function getCustomerList(params: PageRequest): Promise<ApiResponse<PageResponse<Customer>>> {
  return request.get('/admin/customers', { params })
}

/**
 * 添加客户
 */
export function addCustomer(data: CustomerForm): Promise<ApiResponse<Customer>> {
  return request.post('/admin/customers', data)
}

/**
 * 禁用客户
 */
export function disableCustomer(id: number): Promise<ApiResponse<void>> {
  return request.post(`/admin/customers/${id}/disable`)
}

// ==================== 仪表盘 & 监控 ====================

/**
 * 获取Dashboard概览数据
 */
export function getDashboardOverview(): Promise<ApiResponse<any>> {
  return request.get('/admin/dashboard/stats')
}

/**
 * 获取GPU趋势
 */
export function getGPUTrend(): Promise<ApiResponse<any>> {
  return request.get('/admin/dashboard/gpu-trend')
}

/**
 * 获取最近分配记录
 */
export function getRecentAllocations(): Promise<ApiResponse<any>> {
  return request.get('/admin/allocations/recent')
}

/**
 * 获取实时监控数据
 */
export function getRealtimeMonitoring(): Promise<ApiResponse<any>> {
  return request.get('/admin/monitoring/realtime')
}

/**
 * 获取告警列表
 */
export function getAlertList(params: PageRequest & { severity?: string; acknowledged?: boolean }): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/admin/alerts', { params })
}

/**
 * 确认告警
 */
export function acknowledgeAlert(id: number): Promise<ApiResponse<void>> {
  return request.post(`/admin/alerts/${id}/acknowledge`)
}

// ==================== 镜像管理 ====================

/**
 * 获取镜像列表
 */
export function getImageList(params: PageRequest & { category?: string; framework?: string; status?: string }): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/admin/images', { params })
}

/**
 * 同步镜像
 */
export function syncImages(): Promise<ApiResponse<{ message: string }>> {
  return request.post('/admin/images/sync')
}

// ==================== 审计日志 ====================

/**
 * 获取审计日志
 */
export function getAuditLogs(params: PageRequest & { username?: string; action?: string; resource_type?: string }): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/admin/audit/logs', { params })
}

// ==================== 任务管理 ====================

export interface CreateTaskPayload {
  name: string
  type?: string
  customer_id: number
  host_id: string
  image_id?: number
  command: string
  env_vars?: Record<string, string>
}

/**
 * 获取任务列表
 */
export function getTaskList(
  params: PageRequest & { status?: string; customer_id?: number; host_id?: string; keyword?: string }
): Promise<ApiResponse<PageResponse<Task>>> {
  return request.get('/admin/tasks', { params })
}

/**
 * 获取任务详情
 */
export function getTaskDetail(id: string): Promise<ApiResponse<Task>> {
  return request.get(`/admin/tasks/${id}`)
}

/**
 * 创建任务
 */
export function createTask(data: CreateTaskPayload): Promise<ApiResponse<Task>> {
  return request.post('/admin/tasks', data)
}

/**
 * 停止任务
 */
export function stopTask(id: string): Promise<ApiResponse<void>> {
  return request.post(`/admin/tasks/${id}/stop`)
}

/**
 * 取消任务
 */
export function cancelTask(id: string): Promise<ApiResponse<void>> {
  return request.post(`/admin/tasks/${id}/cancel`)
}

/**
 * 重试任务
 */
export function retryTask(id: string): Promise<ApiResponse<void>> {
  return request.post(`/admin/tasks/${id}/retry`)
}

/**
 * 获取任务日志
 */
export function getTaskLogs(id: string): Promise<ApiResponse<TaskLogResponse>> {
  return request.get(`/admin/tasks/${id}/logs`)
}

/**
 * 获取任务结果元信息
 */
export function getTaskResult(id: string): Promise<ApiResponse<TaskResultResponse>> {
  return request.get(`/admin/tasks/${id}/result`)
}

/**
 * 下载任务结果
 */
export function downloadTaskResult(id: string): Promise<Blob> {
  return request.get(`/admin/tasks/${id}/result`, { responseType: 'blob' })
}

// ==================== 系统配置 ====================

/**
 * 获取所有系统配置
 */
export function getSystemConfigs(): Promise<ApiResponse<SystemConfig[]>> {
  return request.get('/admin/settings/configs')
}

/**
 * 批量更新系统配置
 */
export function updateSystemConfigs(data: UpdateSystemConfigsPayload): Promise<ApiResponse<void>> {
  return request.put('/admin/settings/configs', data)
}
