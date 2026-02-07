/**
 * 客户相关类型定义
 */

// 客户状态
export type CustomerStatus = 'active' | 'disabled' | 'pending' | 'suspended' | 'deleted'

// 客户信息
export interface Customer {
  id: number
  username?: string
  display_name?: string
  full_name?: string
  company_code?: string
  company?: string
  email?: string
  phone?: string
  status: CustomerStatus
  created_at?: string
  last_login_at?: string

  name?: string
  contactPerson?: string  // 联系人
  contactEmail?: string  // 联系邮箱
  contactPhone?: string  // 联系电话
  tenantId?: number
  allocatedMachines?: number  // 分配的机器数量
  storageUsed?: number  // 已用存储（GB）
  storageQuota?: number  // 存储配额（GB）
  createdAt?: string
  lastLoginAt?: string
}

// 客户分配的机器信息（来自后端 GetCustomerDetail）
export interface CustomerAllocation {
  allocation_id: string
  machine_id: string
  machine_name?: string
  allocated_at?: string
  end_time?: string
  ssh_host?: string
  ssh_port?: number
  jupyter_url?: string
  vnc_url?: string
  status?: string
}

// 客户详情（后端返回 { customer, allocations }）
export interface CustomerDetailResponse {
  customer: Customer
  allocations: CustomerAllocation[]
}

// 添加客户表单
export interface AddCustomerForm {
  username: string
  company_code: string
  company: string
  email: string
  phone?: string
  password: string
}

// 客户表单类型别名
export type CustomerForm = AddCustomerForm
