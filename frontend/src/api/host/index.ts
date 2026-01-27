/**
 * Host Selection Module - API Functions
 */
import request from '../common/request'
import type { Host, HostQueryParams, HostPricing } from './types'
import type { PaginationResponse } from '../common/types'

/**
 * 获取可用主机列表（带筛选和分页）
 */
export function getAvailableHosts(params?: HostQueryParams) {
  return request.get<PaginationResponse<Host>>('/cmdb/servers/available', { params })
}

/**
 * 获取主机详情
 */
export function getHostDetail(id: number) {
  return request.get<Host>(`/cmdb/servers/${id}`)
}

/**
 * 获取主机价格信息
 */
export function getHostPricing(hostId: number) {
  return request.get<HostPricing>(`/cmdb/servers/${hostId}/pricing`)
}

/**
 * 获取可用地区列表
 */
export function getAvailableRegions() {
  return request.get<string[]>('/cmdb/regions')
}

/**
 * 获取可用GPU型号列表
 */
export function getAvailableGpuModels() {
  return request.get<string[]>('/cmdb/gpu-models')
}
