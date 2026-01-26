/**
 * 问题单管理模块 - API 接口
 */
import request from '../common/request'
import type { Issue, IssueComment, IssueAttachment } from './types'
import type { PaginationParams, PaginationResponse, StatusResponse, IdResponse } from '../common/types'

// ==================== 问题管理 ====================

/**
 * 获取问题列表
 */
export function getIssueList(params?: PaginationParams & {
  status?: string
  type?: string
  priority?: string
  assignee_id?: number
  keyword?: string
}) {
  return request.get<PaginationResponse<Issue>>('/issues', { params })
}

/**
 * 创建问题
 */
export function createIssue(data: {
  title: string
  description: string
  type: string
  priority: string
  assignee_id?: number
  labels?: string[]
}) {
  return request.post<IdResponse>('/issues', data)
}

/**
 * 获取问题详情
 */
export function getIssueDetail(id: number) {
  return request.get<Issue>(`/issues/${id}`)
}

/**
 * 更新问题
 */
export function updateIssue(id: number, data: Partial<Issue>) {
  return request.put<StatusResponse>(`/issues/${id}`, data)
}

/**
 * 分配问题
 */
export function assignIssue(id: number, assigneeId: number) {
  return request.post<StatusResponse>(`/issues/${id}/assign`, { assignee_id: assigneeId })
}

/**
 * 关闭问题
 */
export function closeIssue(id: number, data: { resolution: string; comment?: string }) {
  return request.post<StatusResponse>(`/issues/${id}/close`, data)
}

// ==================== 评论管理 ====================

/**
 * 获取问题评论
 */
export function getIssueComments(issueId: number) {
  return request.get<IssueComment[]>(`/issues/${issueId}/comments`)
}

/**
 * 添加评论
 */
export function addIssueComment(issueId: number, content: string) {
  return request.post<IdResponse>(`/issues/${issueId}/comments`, { content })
}

// ==================== 附件管理 ====================

/**
 * 上传附件
 */
export function uploadIssueAttachment(issueId: number, file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<IdResponse & { file_url: string }>(`/issues/${issueId}/attachments`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

/**
 * 获取附件列表
 */
export function getIssueAttachments(issueId: number) {
  return request.get<IssueAttachment[]>(`/issues/${issueId}/attachments`)
}
