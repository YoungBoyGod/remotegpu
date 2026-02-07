/**
 * 资源调度模块 - API 接口
 */
import request from '@/utils/request'
import type { SchedulerPolicy, ScheduleHistory, PortPool } from './types'
import type { PaginationParams, PaginationResponse, StatusResponse } from '../common/types'

// ==================== 调度历史 ====================

/**
 * 获取调度历史
 */
export function getScheduleHistory(params?: PaginationParams & {
  env_id?: number
  customer_id?: number
  start_date?: string
  end_date?: string
}) {
  return request.get<PaginationResponse<ScheduleHistory>>('/scheduler/history', { params })
}

// ==================== 端口管理 ====================

/**
 * 获取端口使用情况
 */
export function getPortUsage() {
  return request.get<{ total: number; used: number; available: number }>('/scheduler/ports/usage')
}

/**
 * 获取端口池列表
 */
export function getPortPools() {
  return request.get<PortPool[]>('/scheduler/ports/pools')
}

// ==================== 调度策略 ====================

/**
 * 获取调度策略列表
 */
export function getSchedulerPolicies() {
  return request.get<SchedulerPolicy[]>('/scheduler/policies')
}

/**
 * 更新调度策略
 */
export function updateSchedulerPolicy(id: number, data: {
  enabled?: boolean
  priority?: number
  config?: Record<string, any>
}) {
  return request.put<StatusResponse>(`/scheduler/policies/${id}`, data)
}
