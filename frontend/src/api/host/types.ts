/**
 * Host Management Module - Type Definitions
 */
import type { PaginationParams } from '../common/types'

export interface Host {
  id: string
  name: string
  hostname?: string
  ip_address: string
  public_ip?: string
  os_type: string
  os_version?: string
  arch?: string
  deployment_mode?: string
  k8s_node_name?: string
  status: string
  health_status: string
  total_cpu: number
  total_memory: number
  total_disk?: number
  total_gpu: number
  used_cpu: number
  used_memory: number
  used_disk?: number
  used_gpu: number
  ssh_port: number
  winrm_port?: number | null
  agent_port: number
  labels?: Record<string, string>
  tags?: string[]
  last_heartbeat?: string | null
  registered_at?: string
  created_at: string
  updated_at: string
}

export interface HostListResponse {
  list: Host[]
  total: number
}

export interface HostFilterParams {
  keyword?: string
  status?: string
  os_type?: string
}

export interface HostQueryParams extends PaginationParams, HostFilterParams {}
