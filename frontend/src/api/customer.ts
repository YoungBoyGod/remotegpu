import request from '@/utils/request'
import type { ApiResponse, PageRequest, PageResponse } from '@/types/common'
import type { Machine } from '@/types/machine'

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
export function getMachineConnection(id: number): Promise<ApiResponse<any>> {
  return request.get(`/customer/machines/${id}/connection`)
}

/**
 * 重置SSH连接
 */
export function resetSSH(id: number): Promise<ApiResponse<void>> {
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
export function getTasks(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/customer/tasks', { params })
}

/**
 * 创建训练任务
 */
export function createTrainingTask(data: any): Promise<ApiResponse<any>> {
  return request.post('/customer/tasks/training', data)
}

/**
 * 停止任务
 */
export function stopTask(id: number): Promise<ApiResponse<void>> {
  return request.post(`/customer/tasks/${id}/stop`)
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
