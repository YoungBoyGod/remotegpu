/**
 * 需求单管理模块 - 类型定义
 */
import type { Timestamps, UUIDField } from '../common/types'

// 需求单
export interface Requirement extends Timestamps, UUIDField {
  id: number
  customer_id: number
  title: string
  description: string
  type: 'feature' | 'optimization' | 'refactor'
  priority: 'p0' | 'p1' | 'p2' | 'p3'
  status: 'draft' | 'reviewing' | 'approved' | 'rejected' | 'in_progress' | 'completed'
  owner_id: number
  owner_name?: string
  acceptance_criteria?: string
}

// 需求评审
export interface RequirementReview extends Timestamps {
  id: number
  requirement_id: number
  reviewer_id: number
  reviewer_name: string
  score: number
  comment: string
  status: 'approved' | 'rejected' | 'pending'
}

// Sprint
export interface Sprint extends Timestamps {
  id: number
  milestone_id?: number
  name: string
  start_date: string
  end_date: string
  capacity: number
  status: 'planning' | 'active' | 'completed'
}
