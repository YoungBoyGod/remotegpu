/**
 * 机器相关类型定义
 */

// 机器状态
export type MachineStatus = 'online' | 'offline' | 'maintenance'

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

// 机器信息
export interface Machine {
  id: number
  name: string
  region: string
  status: MachineStatus
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
  region: string
  location?: string
  gpuModel: string
  gpuMemory: number
  gpuCount: number
  cudaVersion: string
  gpuDriver: string
  cpu: string
  memory: number
  disk: number
  internalIp: string
  sshHost: string
  sshPort: number
  sshUsername: string
  sshPassword: string
}
