/**
 * CMDB 设备管理模块 - API 接口
 */
import request from '@/utils/request'
import type { Asset, Server, GPU, ChangeLog } from './types'
import type { PaginationParams, PaginationResponse, StatusResponse } from '../common/types'

// ==================== 资产管理 ====================

/**
 * 获取资产列表
 */
export function getAssetList(params?: PaginationParams & {
  type?: string
  operational_status?: string
  keyword?: string
}) {
  return request.get<PaginationResponse<Asset>>('/cmdb/assets', { params })
}

/**
 * 获取资产详情
 */
export function getAssetDetail(id: number) {
  return request.get<Asset>(`/cmdb/assets/${id}`)
}

/**
 * 更新资产状态
 */
export function updateAssetStatus(id: number, data: {
  operational_status: string
  reason?: string
  operator: string
}) {
  return request.put<StatusResponse>(`/cmdb/assets/${id}/status`, data)
}

/**
 * 获取资产变更历史
 */
export function getAssetChanges(id: number, params?: PaginationParams) {
  return request.get<PaginationResponse<ChangeLog>>(`/cmdb/assets/${id}/changes`, { params })
}

// ==================== 服务器管理 ====================

/**
 * 获取服务器列表
 */
export function getServerList(params?: PaginationParams & {
  os_type?: string
  status?: string
  keyword?: string
}) {
  return request.get<PaginationResponse<Server>>('/cmdb/servers', { params })
}

/**
 * 获取服务器详情
 */
export function getServerDetail(id: number) {
  return request.get<Server>(`/cmdb/servers/${id}`)
}

/**
 * 查询可用服务器
 */
export function getAvailableServers(params: {
  cpu?: number
  memory?: number
  gpu?: number
  os_type?: string
}) {
  return request.get<Server[]>('/cmdb/servers/available', { params })
}

/**
 * 分配服务器资源
 */
export function allocateServerResource(id: number, data: {
  env_id: number
  cpu: number
  memory: number
  gpu?: number
}) {
  return request.post<StatusResponse>(`/cmdb/servers/${id}/allocate`, data)
}

/**
 * 释放服务器资源
 */
export function releaseServerResource(id: number, data: { env_id: number }) {
  return request.post<StatusResponse>(`/cmdb/servers/${id}/release`, data)
}

// ==================== GPU 管理 ====================

/**
 * 获取 GPU 列表
 */
export function getGPUList(params?: PaginationParams & {
  server_id?: number
  status?: string
  model?: string
}) {
  return request.get<PaginationResponse<GPU>>('/cmdb/gpus', { params })
}

/**
 * 获取 GPU 详情
 */
export function getGPUDetail(id: number) {
  return request.get<GPU>(`/cmdb/gpus/${id}`)
}

/**
 * 获取服务器的 GPU 列表
 */
export function getServerGPUs(serverId: number) {
  return request.get<GPU[]>(`/cmdb/servers/${serverId}/gpus`)
}
