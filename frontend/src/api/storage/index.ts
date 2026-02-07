/**
 * 数据与存储模块 - API 接口
 */
import request from '@/utils/request'
import type {
  Dataset,
  DatasetVersion,
  Model,
  ModelVersion,
  FileInfo
} from './types'
import type { PaginationParams, PaginationResponse, StatusResponse, IdResponse, UploadResponse, DownloadUrlResponse } from '../common/types'

// ==================== 数据集管理 ====================

/**
 * 获取数据集列表
 */
export function getDatasetList(params?: PaginationParams & {
  visibility?: string
  tag?: string
  keyword?: string
}) {
  return request.get<PaginationResponse<Dataset>>('/datasets', { params })
}

/**
 * 创建数据集
 */
export function createDataset(data: {
  name: string
  description?: string
  visibility: string
  tags?: string[]
}) {
  return request.post<IdResponse & { storage_path: string }>('/datasets', data)
}

/**
 * 获取数据集详情
 */
export function getDatasetDetail(id: number) {
  return request.get<Dataset>(`/datasets/${id}`)
}

/**
 * 更新数据集
 */
export function updateDataset(id: number, data: Partial<Dataset>) {
  return request.put<StatusResponse>(`/datasets/${id}`, data)
}

/**
 * 删除数据集
 */
export function deleteDataset(id: number) {
  return request.delete<StatusResponse>(`/datasets/${id}`)
}

/**
 * 获取上传凭证
 */
export function getUploadUrl(datasetId: number, data: {
  file_name: string
  file_size: number
}) {
  return request.post<{ upload_url: string; expires_in: number }>(`/datasets/${datasetId}/upload-url`, data)
}

/**
 * 完成上传
 */
export function completeUpload(datasetId: number, data: {
  files: Array<{ file_name: string; file_size: number }>
}) {
  return request.post<StatusResponse>(`/datasets/${datasetId}/complete`, data)
}

/**
 * 浏览数据集文件
 */
export function browseDatasetFiles(datasetId: number, params?: { prefix?: string }) {
  return request.get<{ files: FileInfo[] }>(`/datasets/${datasetId}/files`, { params })
}

/**
 * 下载文件
 */
export function getDownloadUrl(datasetId: number, params: { file: string }) {
  return request.get<DownloadUrlResponse>(`/datasets/${datasetId}/download`, { params })
}

// ==================== 模型管理 ====================

/**
 * 获取模型列表
 */
export function getModelList(params?: PaginationParams & {
  framework?: string
  visibility?: string
  keyword?: string
}) {
  return request.get<PaginationResponse<Model>>('/models', { params })
}

/**
 * 创建模型
 */
export function createModel(data: {
  name: string
  framework: string
  description?: string
  visibility?: string
}) {
  return request.post<IdResponse>('/models', data)
}

/**
 * 获取模型详情
 */
export function getModelDetail(id: number) {
  return request.get<Model>(`/models/${id}`)
}

/**
 * 更新模型
 */
export function updateModel(id: number, data: Partial<Model>) {
  return request.put<StatusResponse>(`/models/${id}`, data)
}

/**
 * 删除模型
 */
export function deleteModel(id: number) {
  return request.delete<StatusResponse>(`/models/${id}`)
}

/**
 * 同步预训练模型
 */
export function syncPretrainedModel(name: string) {
  return request.post<IdResponse>(`/models/pretrained/${name}/sync`)
}
