/**
 * 资源配额模块 - API 接口
 */
import request from '@/utils/request'
import type {
  QuotaInfo,
  SetQuotaRequest,
  UpdateQuotaRequest,
  QuotaListResponse,
  QuotaUsageResponse
} from './types'
import type { StatusResponse } from '../common/types'

// ==================== 管理员接口 ====================

/**
 * 设置资源配额
 */
export function setQuota(data: SetQuotaRequest) {
  return request.post<QuotaInfo>('/admin/quotas', data)
}

/**
 * 获取配额列表
 */
export function getQuotas(page: number = 1, pageSize: number = 10) {
  return request.get<QuotaListResponse>('/admin/quotas', {
    params: { page, page_size: pageSize }
  })
}

/**
 * 获取配额详情
 */
export function getQuotaById(id: number) {
  return request.get<QuotaInfo>(`/admin/quotas/${id}`)
}

/**
 * 更新资源配额
 */
export function updateQuota(id: number, data: UpdateQuotaRequest) {
  return request.put<QuotaInfo>(`/admin/quotas/${id}`, data)
}

/**
 * 删除资源配额
 */
export function deleteQuota(id: number) {
  return request.delete<StatusResponse>(`/admin/quotas/${id}`)
}

// ==================== 用户接口 ====================

/**
 * 获取配额使用情况
 */
export function getQuotaUsage(customerId: number, workspaceId?: number) {
  const params: Record<string, any> = {
    customer_id: customerId
  }
  if (workspaceId !== undefined) {
    params.workspace_id = workspaceId
  }
  return request.get<QuotaUsageResponse>('/quotas/usage', { params })
}
