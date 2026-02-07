import request from '@/utils/request'
import type { ApiResponse, PageRequest, PageResponse } from '@/types/common'
import type { Machine } from '@/types/machine'
import type { Customer, CustomerDetailResponse, CustomerForm } from '@/types/customer'
import type { AllocationRecord } from '@/types/allocation'
import type { Task, TaskLogResponse, TaskResultResponse } from '@/types/task'
import type { SystemConfig, UpdateSystemConfigsPayload } from '@/types/systemConfig'

// ==================== 机器管理 ====================

/**
 * 获取机器列表
 */
export function getMachineList(
  params: PageRequest & { status?: string; region?: string; gpu_model?: string; keyword?: string; filters?: Record<string, any> }
): Promise<ApiResponse<PageResponse<Machine>>> {
  const query: Record<string, any> = {
    page: params.page,
    page_size: params.pageSize,
  }
  const source = params.filters || {}
  const status = params.status ?? source.status
  const region = params.region ?? source.region
  const gpuModel = params.gpu_model ?? source.gpu_model
  const keyword = params.keyword ?? source.keyword

  if (status) query.status = status
  if (region) query.region = region
  if (gpuModel) query.gpu_model = gpuModel
  if (keyword) query.keyword = keyword

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
  ssh_host?: string
  ssh_port: number
  ssh_username?: string
  ssh_password?: string
  ssh_key?: string
  jupyter_url?: string
  jupyter_token?: string
  vnc_url?: string
  vnc_password?: string
  external_ip?: string
  external_ssh_port?: number
  external_jupyter_port?: number
  external_vnc_port?: number
  nginx_domain?: string
  nginx_config_path?: string
}

export function addMachine(data: CreateMachinePayload): Promise<ApiResponse<Machine>> {
  return request.post('/admin/machines', data)
}

/**
 * 批量导入机器条目
 */
export interface ImportMachineItem {
  host_ip: string
  ssh_port: number
  region: string
  gpu_model: string
  gpu_count: number
  cpu_cores: number
  ram_size: number
  disk_size: number
  price_hourly: number
}

/**
 * 批量导入机器
 */
export function batchImportMachines(data: { machines: ImportMachineItem[] }): Promise<ApiResponse<{ message: string; count: number }>> {
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

/**
 * 获取机器使用情况
 */
export interface MachineUsageGPU {
  index: number
  name: string
  memory_total_mb: number
  status: string
}

export interface MachineUsage {
  host_id: string
  status: string
  device_status: string
  allocation_status: string
  collected_at: string
  cpu_usage: number
  memory_usage: number
  disk_usage: number
  gpu_usage: MachineUsageGPU[]
}

export function getMachineUsage(id: string): Promise<ApiResponse<MachineUsage>> {
  return request.get(`/admin/machines/${id}/usage`)
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
 * 获取客户详情
 */
export function getCustomerDetail(id: number): Promise<ApiResponse<CustomerDetailResponse>> {
  return request.get(`/admin/customers/${id}`)
}

/**
 * 更新客户信息
 */
export interface UpdateCustomerPayload {
  email?: string
  display_name?: string
  full_name?: string
  company_code?: string
  company?: string
  phone?: string
  role?: string
}

export function updateCustomer(id: number, data: UpdateCustomerPayload): Promise<ApiResponse<void>> {
  return request.put(`/admin/customers/${id}`, data)
}

/**
 * 禁用客户
 */
export function disableCustomer(id: number): Promise<ApiResponse<void>> {
  return request.post(`/admin/customers/${id}/disable`)
}

/**
 * 启用客户
 */
export function enableCustomer(id: number): Promise<ApiResponse<void>> {
  return request.post(`/admin/customers/${id}/enable`)
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
 * 获取分配记录列表
 */
export function getAllocationList(
  params: PageRequest & { status?: string; keyword?: string; filters?: Record<string, any> }
): Promise<ApiResponse<PageResponse<AllocationRecord>>> {
  const query: Record<string, any> = {
    page: params.page,
    page_size: params.pageSize,
  }
  const source = params.filters || {}
  const status = params.status ?? source.status
  const keyword = params.keyword ?? source.keyword
  if (status) query.status = status
  if (keyword) query.keyword = keyword
  return request.get('/admin/allocations', { params: query })
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

/**
 * 批量确认告警
 */
export function batchAcknowledgeAlerts(ids: number[]): Promise<ApiResponse<void>> {
  return request.post('/admin/alerts/batch-acknowledge', { ids })
}

// ==================== 告警规则管理 ====================

export interface AlertRule {
  id: number
  name: string
  metric_type: string
  condition: string
  threshold: number
  severity: string
  duration_seconds: number
  enabled: boolean
  description?: string
  created_at: string
  updated_at: string
}

export interface AlertRuleForm {
  name: string
  metric_type: string
  condition: string
  threshold: number
  severity: string
  duration_seconds: number
  enabled: boolean
  description?: string
}

/**
 * 获取告警规则列表
 */
export function getAlertRuleList(params?: PageRequest): Promise<ApiResponse<PageResponse<AlertRule>>> {
  return request.get('/admin/alert-rules', { params })
}

/**
 * 创建告警规则
 */
export function createAlertRule(data: AlertRuleForm): Promise<ApiResponse<AlertRule>> {
  return request.post('/admin/alert-rules', data)
}

/**
 * 更新告警规则
 */
export function updateAlertRule(id: number, data: AlertRuleForm): Promise<ApiResponse<AlertRule>> {
  return request.put(`/admin/alert-rules/${id}`, data)
}

/**
 * 删除告警规则
 */
export function deleteAlertRule(id: number): Promise<ApiResponse<void>> {
  return request.delete(`/admin/alert-rules/${id}`)
}

/**
 * 启用/禁用告警规则
 */
export function toggleAlertRule(id: number, enabled: boolean): Promise<ApiResponse<void>> {
  return request.post(`/admin/alert-rules/${id}/toggle`, { enabled })
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

// ==================== 文档管理 ====================

export interface DocumentItem {
  id: number
  title: string
  category: string
  file_name: string
  file_path: string
  file_size: number
  content_type: string
  uploaded_by: number
  created_at: string
  updated_at: string
  uploader?: { id: number; username?: string; display_name?: string }
}

/**
 * 获取文档列表
 */
export function getDocumentList(params: PageRequest & { category?: string; keyword?: string }): Promise<ApiResponse<PageResponse<DocumentItem>>> {
  return request.get('/admin/documents', { params })
}

/**
 * 获取文档详情
 */
export function getDocumentDetail(id: number): Promise<ApiResponse<DocumentItem>> {
  return request.get(`/admin/documents/${id}`)
}

/**
 * 上传文档（multipart/form-data）
 */
export function uploadDocument(data: FormData): Promise<ApiResponse<DocumentItem>> {
  return request.post('/admin/documents', data, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

/**
 * 更新文档信息
 */
export function updateDocument(id: number, data: { title?: string; category?: string }): Promise<ApiResponse<void>> {
  return request.put(`/admin/documents/${id}`, data)
}

/**
 * 删除文档
 */
export function deleteDocument(id: number): Promise<ApiResponse<void>> {
  return request.delete(`/admin/documents/${id}`)
}

/**
 * 获取文档分类列表
 */
export function getDocumentCategories(): Promise<ApiResponse<string[]>> {
  return request.get('/admin/documents/categories')
}

/**
 * 获取文档下载链接
 */
export function getDocumentDownloadUrl(id: number): Promise<ApiResponse<{ url: string }>> {
  return request.get(`/admin/documents/${id}/download`)
}

// ==================== 审计日志 ====================

/**
 * 获取审计日志
 */
export function getAuditLogs(params: PageRequest & { username?: string; action?: string; resource_type?: string; start_time?: string; end_time?: string }): Promise<ApiResponse<PageResponse<any>>> {
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

// ==================== Agent 管理 ====================

export interface AgentInfo {
  agent_id: string
  machine_id: string
  machine_name: string
  ip_address: string
  status: string
  last_heartbeat: string | null
  agent_port: number
  region: string
  gpu_count: number
  gpu_models: string[]
}

export interface AgentListResponse {
  total: number
  online: number
  offline: number
  agents: AgentInfo[]
}

/**
 * 获取 Agent 列表
 */
export function getAgentList(): Promise<ApiResponse<AgentListResponse>> {
  return request.get('/admin/agents')
}

// ==================== 系统配置 ====================

/**
 * 获取所有系统配置
 */
export function getSystemConfigs(group?: string): Promise<ApiResponse<SystemConfig[]>> {
  const params: Record<string, any> = {}
  if (group) params.group = group
  return request.get('/admin/settings/configs', { params })
}

/**
 * 获取配置分组列表
 */
export function getConfigGroups(): Promise<ApiResponse<string[]>> {
  return request.get('/admin/settings/configs/groups')
}

/**
 * 批量更新系统配置
 */
export function updateSystemConfigs(data: UpdateSystemConfigsPayload): Promise<ApiResponse<void>> {
  return request.put('/admin/settings/configs', data)
}

// ==================== 存储管理 ====================

export interface StorageBackend {
  name: string
  type: string
  is_default: boolean
}

export interface StorageStats {
  backend_name: string
  file_count: number
  total_size: number
}

/**
 * 获取存储后端列表
 */
export function getStorageBackends(): Promise<ApiResponse<{ backends: StorageBackend[] }>> {
  return request.get('/admin/storage/backends')
}

/**
 * 获取存储统计
 */
export function getStorageStats(backend?: string): Promise<ApiResponse<StorageStats>> {
  const params: Record<string, any> = {}
  if (backend) params.backend = backend
  return request.get('/admin/storage/stats', { params })
}

export interface StorageFileInfo {
  name: string
  size: number
  content_type: string
  last_modified: string
  etag?: string
  is_dir: boolean
}

/**
 * 获取存储文件列表
 */
export function getStorageFiles(backend?: string, prefix?: string): Promise<ApiResponse<{ files: StorageFileInfo[]; total: number }>> {
  const params: Record<string, any> = {}
  if (backend) params.backend = backend
  if (prefix) params.prefix = prefix
  return request.get('/admin/storage/files', { params })
}

/**
 * 删除存储文件
 */
export function deleteStorageFile(data: { backend?: string; path: string }): Promise<ApiResponse<{ message: string }>> {
  return request.post('/admin/storage/files/delete', data)
}

/**
 * 获取存储文件下载链接
 */
export function getStorageDownloadUrl(backend: string, path: string): Promise<ApiResponse<{ url: string }>> {
  return request.get('/admin/storage/files/download-url', { params: { backend, path } })
}
