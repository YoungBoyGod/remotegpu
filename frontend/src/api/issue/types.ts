/**
 * 问题单管理模块 - 类型定义
 */
import type { Timestamps, UUIDField } from '../common/types'

// 问题单
export interface Issue extends Timestamps, UUIDField {
  id: number
  customer_id: number
  title: string
  description: string
  type: 'bug' | 'feature' | 'task' | 'improvement'
  priority: 'low' | 'medium' | 'high' | 'urgent'
  status: 'open' | 'in_progress' | 'resolved' | 'closed' | 'rejected'
  assignee_id?: number
  assignee_name?: string
  labels?: string[]
}

// 问题评论
export interface IssueComment extends Timestamps {
  id: number
  issue_id: number
  customer_id: number
  customer_name: string
  content: string
}

// 问题附件
export interface IssueAttachment extends Timestamps {
  id: number
  issue_id: number
  file_name: string
  file_path: string
  file_size: number
}
