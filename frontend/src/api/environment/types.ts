/**
 * 环境管理模块 - 类型定义
 */
import { Timestamps, UUIDField } from '../common/types'

// 环境状态
export type EnvironmentStatus = 'creating' | 'running' | 'stopped' | 'stopping' | 'starting' | 'failed' | 'deleting'

// 环境信息
export interface Environment extends Timestamps, UUIDField {
  id: number
  customer_id: number
  workspace_id?: number
  name: string
  description?: string
  image: string
  os_type: 'linux' | 'windows'
  status: EnvironmentStatus
  host_id: number
  host_name?: string
  cpu: number
  memory: number
  gpu: number
  storage: number
}

// 环境访问信息
export interface EnvironmentAccess {
  ssh_host?: string
  ssh_port?: number
  ssh_username?: string
  ssh_password?: string
  rdp_host?: string
  rdp_port?: number
  rdp_username?: string
  rdp_password?: string
  jupyter_url?: string
  jupyter_token?: string
}

// 端口映射
export interface PortMapping extends Timestamps {
  id: number
  env_id: number
  service_type: 'ssh' | 'rdp' | 'jupyter' | 'custom'
  internal_port: number
  external_port: number
  protocol: 'tcp' | 'udp'
  status: 'active' | 'inactive'
}

// 数据集使用记录
export interface DatasetUsage extends Timestamps {
  id: number
  dataset_id: number
  dataset_name: string
  env_id: number
  mount_path: string
  mounted_at: string
}

// 创建环境请求
export interface CreateEnvironmentRequest {
  name: string
  description?: string
  image: string
  resources: {
    cpu: number
    memory: number
    gpu?: number
    storage?: number
  }
  datasets?: number[]
  env_vars?: Record<string, string>
  startup_script?: string
}
