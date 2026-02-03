/**
 * 分配相关类型定义
 */

// 分配记录状态
export type AllocationRecordStatus = 'active' | 'expired' | 'reclaimed' | 'pending'

// 分配记录
export interface AllocationRecord {
  id: number
  machineId: number
  machineName: string
  customerId: number
  customerName: string
  allocatedAt: string
  duration: number  // 月数
  expiresAt: string
  status: AllocationRecordStatus
  notes?: string
  operator: string
}

// 快速分配表单
export interface QuickAllocateForm {
  customerId: number | null
  machineIds: number[]
  duration: number
  notes: string
}

// 续期表单
export interface ExtendAllocationForm {
  allocationId: number
  additionalMonths: number
  notes: string
}
