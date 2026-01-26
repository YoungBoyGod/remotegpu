/**
 * 制品管理模块 - 类型定义
 */
import { Timestamps } from '../common/types'

// 制品类型
export type ArtifactType = 'python' | 'npm' | 'maven' | 'docker' | 'generic'

// 制品信息
export interface Artifact extends Timestamps {
  id: number
  name: string
  type: ArtifactType
  version: string
  storage_path: string
  size: number
  checksum?: string
  visibility: 'public' | 'workspace' | 'private'
}

// 制品版本
export interface ArtifactVersion extends Timestamps {
  id: number
  artifact_id: number
  version: string
  release_date: string
  changelog?: string
  download_count: number
}

// 仓库配置
export interface Repository extends Timestamps {
  id: number
  name: string
  type: ArtifactType
  url: string
  enabled: boolean
  credentials?: {
    username?: string
    password?: string
  }
}
