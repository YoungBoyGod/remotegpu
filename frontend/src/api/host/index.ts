/**
 * Host Selection Module - API Functions
 */
import request from '@/utils/request'
import type { Host, HostQueryParams, HostListResponse } from './types'
import type { StatusResponse } from '../common/types'

/**
 * 获取可用主机列表（带筛选和分页）
 */
export function getAvailableHosts(params?: HostQueryParams) {
  return request.get<HostListResponse>('/admin/machines', { params })
}

/**
 * 获取主机详情
 */
export function getHostDetail(id: string) {
  return request.get<Host>(`/admin/machines/${id}`)
}

/**
 * 获取主机价格信息
 */
export function createHost(data: Host) {
  return request.post<Host>('/admin/machines', data)
}

export function updateHost(id: string, data: Host) {
  return request.put<StatusResponse>(`/admin/machines/${id}`, data)
}

export function deleteHost(id: string) {
  return request.delete<StatusResponse>(`/admin/machines/${id}`)
}

export function sendHeartbeat(id: string) {
  return request.post<StatusResponse>(`/admin/machines/${id}/heartbeat`)
}
