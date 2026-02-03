/**
 * 资源调度模块 - 类型定义
 */
import type { Timestamps } from '../common/types'

// 调度策略
export interface SchedulerPolicy extends Timestamps {
  id: number
  name: string
  strategy: 'least_used' | 'round_robin' | 'random' | 'custom'
  enabled: boolean
  priority: number
  config?: Record<string, any>
}

// 调度历史
export interface ScheduleHistory extends Timestamps {
  id: number
  env_id: number
  customer_id: number
  selected_host_id: number
  strategy_used: string
  status: 'success' | 'failed'
  reason?: string
}

// 端口池
export interface PortPool extends Timestamps {
  id: number
  name: string
  port_range_start: number
  port_range_end: number
  used_ports: number[]
  available_count: number
}
