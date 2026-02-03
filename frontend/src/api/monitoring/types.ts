/**
 * 监控告警模块 - 类型定义
 */
import type { Timestamps } from '../common/types'

// 监控指标
export interface Metric {
  timestamp: string
  value: number
}

// 系统指标
export interface SystemMetrics {
  cpu_usage: number
  memory_usage: number
  disk_usage: number
  network_in: number
  network_out: number
}

// GPU 指标
export interface GPUMetrics {
  gpu_usage: number
  memory_usage: number
  temperature: number
  power_usage: number
}

// 告警规则
export interface AlertRule extends Timestamps {
  id: number
  name: string
  metric: string
  threshold: number
  comparison: 'gt' | 'lt' | 'eq' | 'gte' | 'lte'
  severity: 'info' | 'warning' | 'error' | 'critical'
  enabled: boolean
  notification_channels: string[]
}

// 告警历史
export interface Alert extends Timestamps {
  id: number
  rule_id: number
  rule_name: string
  resource_type: string
  resource_id: number
  severity: string
  message: string
  triggered_at: string
  resolved_at?: string
  status: 'firing' | 'resolved' | 'acknowledged'
}
