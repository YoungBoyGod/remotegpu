/**
 * CMDB 设备管理模块 - 类型定义
 */
import type { Timestamps, UUIDField } from '../common/types'

// 资产状态
export type OperationalStatus = 'available' | 'maintenance' | 'faulty' | 'retired' | 'reserved'
export type UsageStatus = 'idle' | 'partial' | 'full' | 'overcommit'
export type HealthStatus = 'healthy' | 'degraded' | 'unhealthy'

// 资产信息
export interface Asset extends Timestamps, UUIDField {
  id: number
  asset_number: string
  name: string
  type: 'server' | 'gpu' | 'storage' | 'network'
  operational_status: OperationalStatus
  usage_status: UsageStatus
  health_status: HealthStatus
  location?: string
  owner?: string
  tags?: string[]
}

// 服务器信息
export interface Server extends Asset {
  hostname: string
  ip_address: string
  os_type: 'linux' | 'windows'
  os_version: string
  cpu_cores: number
  cpu_model: string
  memory_total: number
  memory_used: number
  disk_total: number
  disk_used: number
  gpu_count: number
  deployment_mode: 'traditional' | 'kubernetes'
}

// GPU 信息
export interface GPU extends Timestamps {
  id: number
  asset_id: number
  server_id: number
  gpu_index: number
  uuid: string
  model: string
  memory_total: number
  memory_used: number
  status: 'available' | 'allocated' | 'faulty'
  temperature?: number
  power_usage?: number
  utilization?: number
}

// 变更记录
export interface ChangeLog extends Timestamps {
  id: number
  asset_id: number
  change_type: 'status_change' | 'config_change' | 'resource_change'
  old_value: string
  new_value: string
  reason?: string
  operator: string
  operator_id: number
}
