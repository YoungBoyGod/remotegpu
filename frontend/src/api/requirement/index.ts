/**
 * 需求单管理模块 - API 接口
 */
import request from '../common/request'
import type { Requirement, RequirementReview, Sprint } from './types'
import type { PaginationParams, PaginationResponse, StatusResponse, IdResponse } from '../common/types'

// ==================== 需求管理 ====================

/**
 * 获取需求列表
 */
export function getRequirementList(params?: PaginationParams & {
  status?: string
  type?: string
  priority?: string
  milestone_id?: number
  keyword?: string
}) {
  return request.get<PaginationResponse<Requirement>>('/requirements', { params })
}

/**
 * 创建需求
 */
export function createRequirement(data: {
  title: string
  description: string
  type: string
  priority: string
  acceptance_criteria?: string
}) {
  return request.post<IdResponse>('/requirements', data)
}

/**
 * 获取需求详情
 */
export function getRequirementDetail(id: number) {
  return request.get<Requirement>(`/requirements/${id}`)
}

/**
 * 更新需求
 */
export function updateRequirement(id: number, data: Partial<Requirement>) {
  return request.put<StatusResponse>(`/requirements/${id}`, data)
}

/**
 * 提交评审
 */
export function submitReview(id: number) {
  return request.post<StatusResponse>(`/requirements/${id}/submit-review`)
}

/**
 * 评审需求
 */
export function reviewRequirement(id: number, data: {
  score: number
  comment: string
  approved: boolean
}) {
  return request.post<StatusResponse>(`/requirements/${id}/review`, data)
}

// ==================== Sprint 管理 ====================

/**
 * 获取 Sprint 列表
 */
export function getSprintList(params?: PaginationParams) {
  return request.get<PaginationResponse<Sprint>>('/sprints', { params })
}

/**
 * 创建 Sprint
 */
export function createSprint(data: {
  name: string
  milestone_id?: number
  start_date: string
  end_date: string
  capacity: number
}) {
  return request.post<IdResponse>('/sprints', data)
}

/**
 * 添加需求到 Sprint
 */
export function addRequirementToSprint(sprintId: number, data: {
  requirement_id: number
  story_points: number
}) {
  return request.post<StatusResponse>(`/sprints/${sprintId}/items`, data)
}
