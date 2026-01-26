/**
 * 环境管理模块 - API 接口
 */
import request from '../common/request'
import type {
  Environment,
  EnvironmentAccess,
  PortMapping,
  DatasetUsage,
  CreateEnvironmentRequest
} from './types'
import type { PaginationParams, PaginationResponse, StatusResponse, IdResponse } from '../common/types'

// ==================== 环境管理 ====================

/**
 * 获取环境列表
 */
export function getEnvironmentList(params?: PaginationParams & {
  status?: string
  workspace_id?: number
  keyword?: string
}) {
  return request.get<PaginationResponse<Environment>>('/environments', { params })
}

/**
 * 创建环境
 */
export function createEnvironment(data: CreateEnvironmentRequest) {
  return request.post<IdResponse & { status: string }>('/environments', data)
}

/**
 * 获取环境详情
 */
export function getEnvironmentDetail(id: number) {
  return request.get<Environment>(`/environments/${id}`)
}

/**
 * 启动环境
 */
export function startEnvironment(id: number) {
  return request.post<StatusResponse>(`/environments/${id}/start`)
}

/**
 * 停止环境
 */
export function stopEnvironment(id: number) {
  return request.post<StatusResponse>(`/environments/${id}/stop`)
}

/**
 * 重启环境
 */
export function restartEnvironment(id: number) {
  return request.post<StatusResponse>(`/environments/${id}/restart`)
}

/**
 * 删除环境
 */
export function deleteEnvironment(id: number) {
  return request.delete<StatusResponse>(`/environments/${id}`)
}

// ==================== 访问管理 ====================

/**
 * 获取环境访问信息
 */
export function getEnvironmentAccess(id: number) {
  return request.get<EnvironmentAccess>(`/environments/${id}/access`)
}

/**
 * 重置访问密码
 */
export function resetAccessPassword(id: number) {
  return request.post<{ password: string }>(`/environments/${id}/reset-password`)
}

// ==================== 端口管理 ====================

/**
 * 获取端口映射列表
 */
export function getPortMappings(envId: number) {
  return request.get<PortMapping[]>(`/environments/${envId}/ports`)
}

/**
 * 添加端口映射
 */
export function addPortMapping(envId: number, data: {
  service_type: string
  internal_port: number
  protocol?: string
}) {
  return request.post<IdResponse>(`/environments/${envId}/ports`, data)
}

/**
 * 删除端口映射
 */
export function deletePortMapping(envId: number, portId: number) {
  return request.delete<StatusResponse>(`/environments/${envId}/ports/${portId}`)
}

// ==================== 数据集挂载 ====================

/**
 * 获取已挂载的数据集
 */
export function getMountedDatasets(envId: number) {
  return request.get<DatasetUsage[]>(`/environments/${envId}/datasets`)
}

/**
 * 挂载数据集
 */
export function mountDataset(envId: number, data: {
  dataset_id: number
  mount_path: string
}) {
  return request.post<StatusResponse>(`/environments/${envId}/mount-dataset`, data)
}

/**
 * 卸载数据集
 */
export function unmountDataset(envId: number, datasetId: number) {
  return request.delete<StatusResponse>(`/environments/${envId}/datasets/${datasetId}`)
}
