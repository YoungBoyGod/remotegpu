/**
 * 通知管理模块 - API 接口
 */
import request from '../common/request'
import type { Notification, NotificationSettings } from './types'
import type { PaginationParams, PaginationResponse, StatusResponse } from '../common/types'

// ==================== 通知管理 ====================

/**
 * 获取通知列表
 */
export function getNotificationList(params?: PaginationParams & {
  type?: string
  status?: string
  start_date?: string
  end_date?: string
}) {
  return request.get<PaginationResponse<Notification>>('/notifications', { params })
}

/**
 * 获取未读通知数量
 */
export function getUnreadCount() {
  return request.get<{ count: number }>('/notifications/unread-count')
}

/**
 * 标记为已读
 */
export function markAsRead(id: number) {
  return request.post<StatusResponse>(`/notifications/${id}/read`)
}

/**
 * 批量标记为已读
 */
export function markAllAsRead() {
  return request.post<StatusResponse>('/notifications/read-all')
}

/**
 * 删除通知
 */
export function deleteNotification(id: number) {
  return request.delete<StatusResponse>(`/notifications/${id}`)
}

// ==================== 通知设置 ====================

/**
 * 获取通知设置
 */
export function getNotificationSettings() {
  return request.get<NotificationSettings>('/notifications/settings')
}

/**
 * 更新通知设置
 */
export function updateNotificationSettings(data: Partial<NotificationSettings>) {
  return request.put<StatusResponse>('/notifications/settings', data)
}
