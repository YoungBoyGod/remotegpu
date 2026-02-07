/**
 * 分配相关类型定义
 */

// 分配记录状态
export type AllocationRecordStatus = 'active' | 'expired' | 'reclaimed' | 'pending'

// 分配记录（字段与后端 snake_case 一致）
export interface AllocationRecord {
  id: string
  customer_id: number
  host_id: string
  workspace_id?: number | null
  start_time: string
  end_time: string
  actual_end_time?: string | null
  status: AllocationRecordStatus
  remark?: string
  created_at?: string
  updated_at?: string
  customer?: {
    id: number
    username?: string
    display_name?: string
    company?: string
  }
  host?: {
    id: string
    name?: string
    region?: string
    ip_address?: string
  }
}

// 续期表单
export interface ExtendAllocationForm {
  allocationId: number
  additionalMonths: number
  notes: string
}
