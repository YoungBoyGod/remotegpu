/**
 * 机器相关类型定义
 */

// 机器状态
export type MachineStatus = 'idle' | 'allocated' | 'maintenance' | 'offline'

// 分配状态
export type AllocationStatus = 'allocated' | 'unallocated' | 'expiring'

// SSH 登录信息
export interface LoginInfo {
  sshHost: string
  sshPort: number
  username: string
  password: string
  jupyterUrl?: string
}

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
  status: MachineStatus
  needs_collect?: boolean
  hostname?: string
  ip_address?: string
  public_ip?: string
  ssh_port?: number
  ssh_username?: string
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
  gpuModel: string
  gpuMemory: number
  gpuCount: number
  gpuUsage: number  // GPU使用率 0-100
  cpu: string
  memory: number
  memoryUsage: number  // 内存使用率
  disk: number
  diskUsage: number  // 磁盘使用率
  cudaVersion: string
  gpuDriver: string
  allocatedAt?: string  // 分配时间
  loginInfo: LoginInfo
  allocatedTo?: {
    customerId: number
    customerName: string
    allocatedAt: string
    duration: number  // 月数
    expiresAt: string  // 到期时间
  }
}

// 机器监控数据
export interface MachineMonitoring {
  machineId: number
  timestamp: number
  gpuUsage: number
  gpuMemory: number
  gpuTemperature: number
  cpuUsage: number
  memoryUsage: number
  diskUsage: number
  networkIn: number
  networkOut: number
}

// 监控数据类型别名
export type MonitoringData = MachineMonitoring

// 进程信息
export interface ProcessInfo {
  pid: number
  name: string
  user: string
  gpuUsage: number
  memoryUsage: number
  cpuUsage: number
}

// 添加机器表单
export interface AddMachineForm {
  name: string
  hostname?: string
  region: string
  ipAddress: string
  publicIp?: string
  sshPort: number
  sshUsername: string
  sshPassword?: string
}
