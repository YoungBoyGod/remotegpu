/**
 * 制品管理模块 - API 接口
 */
import request from '@/utils/request'
import type { Artifact, ArtifactVersion, Repository } from './types'
import type { PaginationParams, PaginationResponse, StatusResponse, IdResponse } from '../common/types'

// ==================== 制品管理 ====================

/**
 * 获取制品列表
 */
export function getArtifactList(params?: PaginationParams & {
  type?: string
  keyword?: string
}) {
  return request.get<PaginationResponse<Artifact>>('/artifacts', { params })
}

/**
 * 上传制品
 */
export function uploadArtifact(data: FormData) {
  return request.post<IdResponse & { download_url: string }>('/artifacts', data, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

/**
 * 获取制品详情
 */
export function getArtifactDetail(id: number) {
  return request.get<Artifact>(`/artifacts/${id}`)
}

/**
 * 下载制品
 */
export function downloadArtifact(id: number) {
  return request.get(`/artifacts/${id}/download`, { responseType: 'blob' })
}

/**
 * 删除制品
 */
export function deleteArtifact(id: number) {
  return request.delete<StatusResponse>(`/artifacts/${id}`)
}

/**
 * 获取制品版本列表
 */
export function getArtifactVersions(artifactId: number) {
  return request.get<ArtifactVersion[]>(`/artifacts/${artifactId}/versions`)
}

// ==================== 仓库管理 ====================

/**
 * 获取仓库列表
 */
export function getRepositoryList() {
  return request.get<Repository[]>('/repositories')
}

/**
 * 创建仓库配置
 */
export function createRepository(data: {
  name: string
  type: string
  url: string
  credentials?: any
}) {
  return request.post<IdResponse>('/repositories', data)
}

/**
 * 更新仓库配置
 */
export function updateRepository(id: number, data: Partial<Repository>) {
  return request.put<StatusResponse>(`/repositories/${id}`, data)
}

/**
 * 删除仓库配置
 */
export function deleteRepository(id: number) {
  return request.delete<StatusResponse>(`/repositories/${id}`)
}
