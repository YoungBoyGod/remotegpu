/**
 * 通知管理模块 - 类型定义
 */
import type { Timestamps } from '../common/types'

// 通知类型
export type NotificationType = 'system' | 'alert' | 'billing' | 'task' | 'message'

// 通知渠道
export type NotificationChannel = 'email' | 'sms' | 'webhook' | 'in_app'

// 通知信息
export interface Notification extends Timestamps {
  id: number
  customer_id: number
  type: NotificationType
  title: string
  content: string
  channel: NotificationChannel
  status: 'pending' | 'sent' | 'failed' | 'read'
  read_at?: string
  metadata?: Record<string, any>
}

// 通知设置
export interface NotificationSettings {
  customer_id: number
  email_enabled: boolean
  sms_enabled: boolean
  webhook_enabled: boolean
  in_app_enabled: boolean
  alert_types: string[]
}
