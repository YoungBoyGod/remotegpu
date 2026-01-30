/**
 * Host Selection Module - API Functions
 */
import request from '../common/request'
import type { Host, HostQueryParams, HostListResponse } from './types'
import type { StatusResponse } from '../common/types'

/**
 * 获取可用主机列表（带筛选和分页）
 */
export function getAvailableHosts(params?: HostQueryParams) {
  return request.get<HostListResponse>('/admin/hosts', { params })
}

/**
 * 获取主机详情
 */
export function getHostDetail(id: string) {
  return request.get<Host>(`/admin/hosts/${id}`)
}

/**
 * 获取主机价格信息
 */
export function createHost(data: Host) {
  return request.post<Host>('/admin/hosts', data)
}

export function updateHost(id: string, data: Host) {
  return request.put<StatusResponse>(`/admin/hosts/${id}`, data)
}

export function deleteHost(id: string) {
  return request.delete<StatusResponse>(`/admin/hosts/${id}`)
}

export function sendHeartbeat(id: string) {
  return request.post<StatusResponse>(`/admin/hosts/${id}/heartbeat`)
}
