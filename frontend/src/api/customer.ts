import request from '@/utils/request'
import type { ApiResponse, PageRequest, PageResponse } from '@/types/common'
import type { Machine } from '@/types/machine'
import type { Task, TaskLogResponse, TaskResultResponse } from '@/types/task'

// ==================== 我的机器 ====================

/**
 * 获取我的机器列表
 */
export function getMyMachines(params?: PageRequest): Promise<ApiResponse<PageResponse<Machine>>> {
  return request.get('/customer/machines', { params })
}

/**
 * 获取机器连接信息
 */
export function getMachineConnection(id: string | number): Promise<ApiResponse<any>> {
  return request.get(`/customer/machines/${id}/connection`)
}

/**
 * 重置SSH连接
 */
export function resetSSH(id: string | number): Promise<ApiResponse<void>> {
  return request.post(`/customer/machines/${id}/ssh-reset`)
}

// ==================== 用户添加机器 ====================

export interface CreateMachineEnrollmentPayload {
  name?: string
  hostname?: string
  region: string
  ip_address?: string
  ssh_port: number
  ssh_username: string
  ssh_password?: string
  ssh_key?: string
}

export interface MachineEnrollment {
  id: number
  customer_id: number
  name?: string
  hostname?: string
  region: string
  address: string
  ssh_port: number
  ssh_username: string
  status: string
  error_message?: string
  host_id?: string
  created_at: string
  updated_at: string
}

export function createMachineEnrollment(data: CreateMachineEnrollmentPayload): Promise<ApiResponse<MachineEnrollment>> {
  return request.post('/customer/machines', data)
}

export function getMachineEnrollments(params?: PageRequest): Promise<ApiResponse<PageResponse<MachineEnrollment>>> {
  return request.get('/customer/machines/enrollments', { params })
}

export function getMachineEnrollment(id: number): Promise<ApiResponse<MachineEnrollment>> {
  return request.get(`/customer/machines/enrollments/${id}`)
}

// ==================== 任务管理 ====================

/**
 * 获取任务列表
 */
export function getTasks(params: PageRequest & { status?: string; keyword?: string; host_id?: string }): Promise<ApiResponse<PageResponse<Task>>> {
  return request.get('/customer/tasks', { params })
}

/**
 * 创建训练任务
 */
export interface CreateTrainingTaskPayload {
  name: string
  host_id: string
  command: string
  image_id?: number
  env_vars?: Record<string, string>
}

export function createTrainingTask(data: CreateTrainingTaskPayload): Promise<ApiResponse<Task>> {
  return request.post('/customer/tasks/training', data)
}

/**
 * 停止任务
 */
export function stopTask(id: string): Promise<ApiResponse<void>> {
  return request.post(`/customer/tasks/${id}/stop`)
}

/**
 * 获取任务详情
 */
export function getTaskDetail(id: string): Promise<ApiResponse<Task>> {
  return request.get(`/customer/tasks/${id}`)
}

/**
 * 取消任务
 */
export function cancelTask(id: string): Promise<ApiResponse<void>> {
  return request.post(`/customer/tasks/${id}/cancel`)
}

/**
 * 重试任务
 */
export function retryTask(id: string): Promise<ApiResponse<void>> {
  return request.post(`/customer/tasks/${id}/retry`)
}

/**
 * 获取任务日志
 */
export function getTaskLogs(id: string): Promise<ApiResponse<TaskLogResponse>> {
  return request.get(`/customer/tasks/${id}/logs`)
}

/**
 * 获取任务结果元信息
 */
export function getTaskResult(id: string): Promise<ApiResponse<TaskResultResponse>> {
  return request.get(`/customer/tasks/${id}/result`)
}

/**
 * 下载任务结果
 */
export function downloadTaskResult(id: string): Promise<Blob> {
  return request.get(`/customer/tasks/${id}/result`, { responseType: 'blob' })
}

// ==================== 数据集管理 ====================

/**
 * 获取数据集列表
 */
export function getDatasetList(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/customer/datasets', { params })
}

/**
 * 初始化分片上传
 */
export function initMultipartUpload(data: { filename: string; size: number }): Promise<ApiResponse<any>> {
  return request.post('/customer/datasets/init-multipart', data)
}

/**
 * 挂载数据集
 */
export function mountDataset(id: number, data: { machineId: number; mountPath: string }): Promise<ApiResponse<any>> {
  return request.post(`/customer/datasets/${id}/mount`, data)
}

// ==================== SSH 密钥管理 ====================

/**
 * 获取SSH密钥列表
 */
export function getSshKeys(): Promise<ApiResponse<any[]>> {
  return request.get('/customer/keys')
}

/**
 * 添加SSH密钥
 */
export function addSshKey(data: { name: string; publicKey: string }): Promise<ApiResponse<any>> {
  return request.post('/customer/keys', data)
}

/**
 * 删除SSH密钥
 */
export function deleteSshKey(id: number): Promise<ApiResponse<void>> {
  return request.delete(`/customer/keys/${id}`)
}
