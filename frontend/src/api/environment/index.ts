/**
 * 环境管理模块 - API 接口
 */
import request from '@/utils/request'
import type {
  Environment,
  CreateEnvironmentRequest
} from './types'
import type { StatusResponse } from '../common/types'

// ==================== 环境管理 ====================

/**
 * 获取环境列表
 */
export function getEnvironmentList(params?: {
  workspace_id?: number
}) {
  return request.get<Environment[]>('/environments', { params })
}

/**
 * 创建环境
 */
export function createEnvironment(data: CreateEnvironmentRequest) {
  return request.post<Environment>('/environments', data)
}

/**
 * 获取环境详情
 */
export function getEnvironmentDetail(id: string) {
  return request.get<Environment>(`/environments/${id}`)
}

/**
 * 启动环境
 */
export function startEnvironment(id: string) {
  return request.post<StatusResponse>(`/environments/${id}/start`)
}

/**
 * 停止环境
 */
export function stopEnvironment(id: string) {
  return request.post<StatusResponse>(`/environments/${id}/stop`)
}

/**
 * 删除环境
 */
export function deleteEnvironment(id: string) {
  return request.delete<StatusResponse>(`/environments/${id}`)
}

/**
 * 获取环境访问信息
 */
export function getEnvironmentAccessInfo(id: string) {
  return request.get<any>(`/environments/${id}/access`)
}
