/**
 * 镜像管理模块 - API 接口
 */
import request from '../common/request'
import type { Image, CustomImage, ImageBuild } from './types'
import type { PaginationParams, PaginationResponse, StatusResponse, IdResponse } from '../common/types'

// ==================== 官方镜像 ====================

/**
 * 获取官方镜像列表
 */
export function getOfficialImages(params?: { category?: string }) {
  return request.get<Image[]>('/images/official', { params })
}

/**
 * 获取镜像详情
 */
export function getImageDetail(name: string) {
  return request.get<Image>(`/images/${name}`)
}

// ==================== 自定义镜像 ====================

/**
 * 获取自定义镜像列表
 */
export function getCustomImages(params?: PaginationParams) {
  return request.get<PaginationResponse<CustomImage>>('/images/custom', { params })
}

/**
 * 创建自定义镜像
 */
export function createCustomImage(data: {
  name: string
  base_image: string
  dockerfile: string
  visibility?: string
}) {
  return request.post<IdResponse & { status: string }>('/images/custom', data)
}

/**
 * 获取自定义镜像详情
 */
export function getCustomImageDetail(id: number) {
  return request.get<CustomImage>(`/images/custom/${id}`)
}

/**
 * 删除自定义镜像
 */
export function deleteCustomImage(id: number) {
  return request.delete<StatusResponse>(`/images/custom/${id}`)
}

/**
 * 获取构建状态
 */
export function getBuildStatus(id: number) {
  return request.get<{ status: string; build_log?: string }>(`/images/custom/${id}/build-status`)
}

/**
 * 获取构建历史
 */
export function getBuildHistory(id: number) {
  return request.get<ImageBuild[]>(`/images/custom/${id}/builds`)
}
