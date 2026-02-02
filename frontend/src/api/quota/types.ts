/**
 * 资源配额模块 - 类型定义
 */

// 配额级别
export type QuotaLevel = 'free' | 'basic' | 'pro' | 'enterprise'

// 资源配额信息
export interface QuotaInfo {
  id: number
  customer_id: number
  workspace_id: number | null
  quota_level: QuotaLevel
  max_gpu: number
  max_cpu: number
  max_memory: number
  max_storage: number
  max_environments: number
  created_at: string
  updated_at: string
}

// 设置配额请求
export interface SetQuotaRequest {
  customer_id: number
  workspace_id?: number | null
  max_gpu: number
  max_cpu: number
  max_memory: number
  max_storage: number
  max_environments: number
  quota_level?: QuotaLevel
}

// 更新配额请求
export interface UpdateQuotaRequest {
  max_gpu: number
  max_cpu: number
  max_memory: number
  max_storage: number
  max_environments: number
  quota_level?: QuotaLevel
}

// 配额列表响应
export interface QuotaListResponse {
  items: QuotaInfo[]
  total: number
  page: number
  page_size: number
}

// 配额详情
export interface QuotaDetail {
  max_gpu: number
  max_cpu: number
  max_memory: number
  max_storage: number
  max_environments: number
}

// 已使用资源
export interface UsedResources {
  used_gpu: number
  used_cpu: number
  used_memory: number
  used_storage: number
  used_environments: number
}

// 可用资源
export interface AvailableResources {
  available_gpu: number
  available_cpu: number
  available_memory: number
  available_storage: number
  available_environments: number
}

// 使用百分比
export interface UsagePercentage {
  gpu: number
  cpu: number
  memory: number
  storage: number
  environments: number
}

// 配额使用情况响应
export interface QuotaUsageResponse {
  quota: QuotaDetail
  used: UsedResources
  available: AvailableResources
  usage_percentage: UsagePercentage
}
