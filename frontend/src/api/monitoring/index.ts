/**
 * 监控告警模块 - API 接口
 */
import request from '../common/request'
import type { Metric, SystemMetrics, GPUMetrics, AlertRule, Alert } from './types'
import type { PaginationParams, PaginationResponse, StatusResponse, IdResponse } from '../common/types'

// ==================== 资源监控 ====================

/**
 * 获取主机监控数据
 */
export function getHostMetrics(hostId: number, params: {
  metric_name: string
  start_time: string
  end_time: string
}) {
  return request.get<{ metrics: Metric[] }>(`/monitoring/hosts/${hostId}/metrics`, { params })
}

/**
 * 获取 GPU 监控数据
 */
export function getGPUMetrics(gpuId: number, params: {
  start_time: string
  end_time: string
}) {
  return request.get<{ metrics: GPUMetrics[] }>(`/monitoring/gpus/${gpuId}/metrics`, { params })
}

/**
 * 获取环境监控数据
 */
export function getEnvironmentMetrics(envId: number, params: {
  start_time: string
  end_time: string
}) {
  return request.get<{ metrics: SystemMetrics[] }>(`/monitoring/environments/${envId}/metrics`, { params })
}

// ==================== 告警管理 ====================

/**
 * 获取告警规则列表
 */
export function getAlertRules(params?: PaginationParams) {
  return request.get<PaginationResponse<AlertRule>>('/monitoring/alert-rules', { params })
}

/**
 * 创建告警规则
 */
export function createAlertRule(data: {
  name: string
  metric: string
  threshold: number
  comparison: string
  severity: string
  notification_channels?: string[]
}) {
  return request.post<IdResponse>('/monitoring/alert-rules', data)
}

/**
 * 更新告警规则
 */
export function updateAlertRule(id: number, data: Partial<AlertRule>) {
  return request.put<StatusResponse>(`/monitoring/alert-rules/${id}`, data)
}

/**
 * 删除告警规则
 */
export function deleteAlertRule(id: number) {
  return request.delete<StatusResponse>(`/monitoring/alert-rules/${id}`)
}

/**
 * 获取告警历史
 */
export function getAlerts(params?: PaginationParams & {
  status?: string
  severity?: string
  start_date?: string
  end_date?: string
}) {
  return request.get<PaginationResponse<Alert>>('/monitoring/alerts', { params })
}

/**
 * 确认告警
 */
export function acknowledgeAlert(id: number) {
  return request.post<StatusResponse>(`/monitoring/alerts/${id}/acknowledge`)
}
