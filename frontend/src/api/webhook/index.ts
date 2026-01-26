/**
 * Webhook 管理模块 - API 接口
 */
import request from '../common/request'
import type { Webhook, WebhookLog } from './types'
import type { PaginationParams, PaginationResponse, StatusResponse, IdResponse } from '../common/types'

// ==================== Webhook 管理 ====================

/**
 * 获取 Webhook 列表
 */
export function getWebhookList(params?: PaginationParams) {
  return request.get<PaginationResponse<Webhook>>('/webhooks', { params })
}

/**
 * 创建 Webhook
 */
export function createWebhook(data: {
  name: string
  url: string
  secret?: string
  events: string[]
}) {
  return request.post<IdResponse>('/webhooks', data)
}

/**
 * 获取 Webhook 详情
 */
export function getWebhookDetail(id: number) {
  return request.get<Webhook>(`/webhooks/${id}`)
}

/**
 * 更新 Webhook
 */
export function updateWebhook(id: number, data: Partial<Webhook>) {
  return request.put<StatusResponse>(`/webhooks/${id}`, data)
}

/**
 * 删除 Webhook
 */
export function deleteWebhook(id: number) {
  return request.delete<StatusResponse>(`/webhooks/${id}`)
}

/**
 * 测试 Webhook
 */
export function testWebhook(id: number) {
  return request.post<StatusResponse>(`/webhooks/${id}/test`)
}

// ==================== Webhook 日志 ====================

/**
 * 获取 Webhook 日志
 */
export function getWebhookLogs(webhookId: number, params?: PaginationParams) {
  return request.get<PaginationResponse<WebhookLog>>(`/webhooks/${webhookId}/logs`, { params })
}
