import request from '@/utils/request'
import type { ApiResponse } from '@/types/common'

// Proxy 节点接口
export interface ProxyNode {
  id: string
  name: string
  host: string
  api_port: number
  http_port: number
  range_start: number
  range_end: number
  version: string
  status: string
  active_mappings: number
  used_ports: number
  last_heartbeat: string | null
  created_at: string
  updated_at: string
}

// Proxy 端口映射接口
export interface ProxyMapping {
  id: number
  env_id: string
  service_type: string
  external_port: number
  internal_port: number
  proxy_id: string
  target_host: string
  target_port: number
  protocol: string
  status: string
}

// 获取 Proxy 节点列表
export function getProxyNodes(): Promise<ApiResponse<ProxyNode[]>> {
  return request.get('/admin/proxy/nodes')
}

// 获取 Proxy 节点详情
export function getProxyNodeDetail(id: string): Promise<ApiResponse<ProxyNode>> {
  return request.get(`/admin/proxy/nodes/${id}`)
}

// 删除 Proxy 节点
export function deleteProxyNode(id: string): Promise<ApiResponse<void>> {
  return request.delete(`/admin/proxy/nodes/${id}`)
}

// 获取所有端口映射
export function getProxyMappings(): Promise<ApiResponse<ProxyMapping[]>> {
  return request.get('/admin/proxy/mappings')
}
