/**
 * 客户相关类型定义
 */

// 客户状态
export type CustomerStatus = 'active' | 'disabled' | 'pending'

// 客户信息
export interface Customer {
  id: number
  name: string
  company: string
  email: string
  phone: string
  contactPerson: string  // 联系人
  contactEmail: string  // 联系邮箱
  contactPhone: string  // 联系电话
  status: CustomerStatus
  tenantId: number
  allocatedMachines: number  // 分配的机器数量
  storageUsed: number  // 已用存储（GB）
  storageQuota: number  // 存储配额（GB）
  createdAt: string
  lastLoginAt?: string
}

// 客户详情
export interface CustomerDetail extends Customer {
  machines: Array<{
    id: number
    name: string
    gpuModel: string
    allocatedAt: string
    expiresAt: string
  }>
  usageStats: {
    allocatedMachines: number  // 分配的机器数
    runningTasks: number  // 运行中任务数
    totalTasks: number  // 总任务数
    storageUsed: number  // 存储使用（GB）
    totalUsageHours: number
    gpuUsageAvg: number
  }
  operationLogs: Array<{
    id: number
    action: string
    operator: string
    timestamp: string
    details: string
  }>
}

// 添加客户表单
export interface AddCustomerForm {
  name: string
  company: string
  email: string
  phone: string
  storageQuota: number
  password: string
}

// 客户表单类型别名
export type CustomerForm = AddCustomerForm
