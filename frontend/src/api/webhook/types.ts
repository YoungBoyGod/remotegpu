/**
 * Webhook 管理模块 - 类型定义
 */
import type { Timestamps } from '../common/types'

// Webhook 事件类型
export type WebhookEvent =
  | 'environment.created'
  | 'environment.started'
  | 'environment.stopped'
  | 'environment.deleted'
  | 'training.started'
  | 'training.completed'
  | 'training.failed'
  | 'alert.triggered'
  | 'billing.charged'

// Webhook 配置
export interface Webhook extends Timestamps {
  id: number
  customer_id: number
  name: string
  url: string
  secret?: string
  events: WebhookEvent[]
  enabled: boolean
  last_triggered_at?: string
}

// Webhook 日志
export interface WebhookLog extends Timestamps {
  id: number
  webhook_id: number
  event: string
  payload: Record<string, any>
  response_status?: number
  response_body?: string
  status: 'success' | 'failed'
  error_message?: string
}
