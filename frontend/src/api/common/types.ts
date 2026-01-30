/**
 * 通用类型定义
 */

// 分页请求参数
export interface PaginationParams {
  page: number
  page_size: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

// 分页响应数据
export interface PaginationResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 通用响应结构
export interface ApiResponse<T = any> {
  code: number
  msg: string
  data: T
}

// 通用列表响应
export interface ListResponse<T> {
  items: T[]
  total: number
}

// 通用 ID 响应
export interface IdResponse {
  id: number | string
}

// 通用状态响应
export interface StatusResponse {
  status?: string
  message?: string
}

// 时间戳字段
export interface Timestamps {
  created_at: string
  updated_at: string
  deleted_at?: string
}

// UUID 字段
export interface UUIDField {
  uuid: string
}

// 操作人字段
export interface OperatorFields {
  created_by?: number
  updated_by?: number
  operator_name?: string
}

// 通用查询参数
export interface CommonQueryParams {
  keyword?: string
  status?: string
  start_date?: string
  end_date?: string
}

// 批量操作参数
export interface BatchOperationParams {
  ids: (number | string)[]
  action: string
}

// 文件上传响应
export interface UploadResponse {
  file_id: string
  file_name: string
  file_url: string
  file_size: number
}

// 下载 URL 响应
export interface DownloadUrlResponse {
  download_url: string
  expires_in: number
}
