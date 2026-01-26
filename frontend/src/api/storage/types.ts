/**
 * 数据与存储模块 - 类型定义
 */
import { Timestamps, UUIDField } from '../common/types'

// 可见性类型
export type Visibility = 'public' | 'workspace' | 'private'

// 数据集信息
export interface Dataset extends Timestamps, UUIDField {
  id: number
  customer_id: number
  workspace_id?: number
  name: string
  description?: string
  visibility: Visibility
  storage_path: string
  total_size: number
  file_count: number
  status: 'uploading' | 'ready' | 'error'
  tags?: string[]
}

// 数据集版本
export interface DatasetVersion extends Timestamps {
  id: number
  dataset_id: number
  version: string
  storage_path: string
  size: number
  is_default: boolean
  changelog?: string
}

// 模型信息
export interface Model extends Timestamps, UUIDField {
  id: number
  customer_id: number
  workspace_id?: number
  name: string
  description?: string
  framework: 'pytorch' | 'tensorflow' | 'onnx' | 'other'
  visibility: Visibility
  storage_path: string
  total_size: number
  status: 'uploading' | 'ready' | 'error'
  tags?: string[]
}

// 模型版本
export interface ModelVersion extends Timestamps {
  id: number
  model_id: number
  version: string
  storage_path: string
  size: number
  metrics?: Record<string, number>
  is_default: boolean
}

// 文件信息
export interface FileInfo {
  name: string
  path: string
  size: number
  is_directory: boolean
  modified_at: string
}
