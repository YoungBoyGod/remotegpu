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
  contact_person?: string  // 联系人
  contact_email?: string  // 联系邮箱
  contact_phone?: string  // 联系电话
  tenant_id?: number
  allocated_machines?: number  // 分配的机器数量
  storage_used?: number  // 已用存储（GB）
  storage_quota?: number  // 存储配额（GB）
}

// 客户分配的机器信息（来自后端 GetCustomerDetail）
export interface CustomerAllocation {
  allocation_id: string
  machine_id: string
  machine_name?: string
  allocated_at?: string
  end_time?: string
  status?: string
  // 对内连接
  ip_address?: string
  ssh_host?: string
  ssh_port?: number
  ssh_username?: string
  ssh_password?: string
  jupyter_url?: string
  vnc_url?: string
  // 对外连接
  external_ip?: string
  nginx_domain?: string
  external_ssh_port?: number
  external_jupyter_port?: number
  external_vnc_port?: number
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
