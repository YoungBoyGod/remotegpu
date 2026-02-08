/**
 * 机器相关类型定义
 */

// 设备在线状态
export type DeviceStatus = 'online' | 'offline'

// 机器分配状态
export type MachineAllocationStatus = 'idle' | 'allocated' | 'maintenance'

export interface MachineGPU {
  id: number
  host_id: string
  index: number
  uuid: string
  name: string
  memory_total_mb: number
  brand?: string
  status?: string
  health_status?: string
  allocated_to?: string
  updated_at?: string
}

export interface MachineAllocation {
  id: string
  customer_id: number
  host_id: string
  start_time: string
  end_time: string
  actual_end_time?: string | null
  status: string
  remark?: string
  created_at?: string
  updated_at?: string
  customer?: {
    id: number
    username?: string
    display_name?: string
    company?: string
  }
}

// 机器信息
export interface Machine {
  id: number | string
  name: string
  region: string
  device_status?: DeviceStatus
  allocation_status?: MachineAllocationStatus
  needs_collect?: boolean
  hostname?: string
  ip_address?: string
  public_ip?: string
  ssh_port?: number
  jupyter_port?: number
  vnc_port?: number
  ssh_host?: string
  ssh_username?: string
  ssh_password?: string
  ssh_key?: string
  ssh_command?: string
  jupyter_url?: string
  jupyter_token?: string
  vnc_url?: string
  vnc_password?: string
  external_ip?: string
  external_ssh_port?: number
  external_jupyter_port?: number
  external_vnc_port?: number
  nginx_domain?: string
  nginx_config_path?: string
  agent_port?: number
  start_time?: string
  end_time?: string
  os_type?: string
  os_version?: string
  cpu_info?: string
  total_cpu?: number
  total_memory_gb?: number
  total_disk_gb?: number
  health_status?: string
  deployment_mode?: string
  last_heartbeat?: string | null
  created_at?: string
  updated_at?: string
  gpus?: MachineGPU[]
  allocations?: MachineAllocation[]
}
