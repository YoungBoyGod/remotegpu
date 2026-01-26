/**
 * 镜像管理模块 - 类型定义
 */
import { Timestamps } from '../common/types'

// 镜像信息
export interface Image extends Timestamps {
  id: number
  name: string
  description?: string
  category: 'base' | 'pytorch' | 'tensorflow' | 'custom'
  is_official: boolean
  size: number
  tags: string[]
  pull_count: number
}

// 自定义镜像
export interface CustomImage extends Timestamps {
  id: number
  uuid: string
  customer_id: number
  name: string
  base_image: string
  dockerfile: string
  status: 'building' | 'ready' | 'failed'
  size?: number
  visibility: 'public' | 'workspace' | 'private'
}

// 构建历史
export interface ImageBuild extends Timestamps {
  id: number
  image_id: number
  build_number: number
  status: 'pending' | 'building' | 'success' | 'failed'
  build_log?: string
  duration?: number
}
