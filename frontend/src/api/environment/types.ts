/**
 * 环境管理模块 - 类型定义
 */

export type EnvironmentStatus = 'creating' | 'running' | 'stopped' | 'error' | 'deleting'

export interface Environment {
  id: string
  customer_id: number
  workspace_id?: number
  host_id: string
  name: string
  description?: string
  image: string
  status: EnvironmentStatus | string
  cpu: number
  memory: number
  gpu: number
  storage?: number | null
  ssh_port?: number | null
  rdp_port?: number | null
  jupyter_port?: number | null
  container_id?: string
  pod_name?: string
  created_at: string
  updated_at: string
  started_at?: string | null
  stopped_at?: string | null
}

export interface CreateEnvironmentRequest {
  customer_id: number
  workspace_id?: number
  name: string
  description?: string
  image: string
  cpu: number
  memory: number
  gpu: number
  storage?: number
  command?: string[]
  args?: string[]
  env?: Record<string, string>
}
