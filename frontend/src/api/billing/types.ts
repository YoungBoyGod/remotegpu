/**
 * 计费管理模块 - 类型定义
 */
import { Timestamps } from '../common/types'

// 计费记录
export interface BillingRecord extends Timestamps {
  id: number
  customer_id: number
  env_id?: number
  resource_type: 'cpu' | 'memory' | 'gpu' | 'storage' | 'network'
  quantity: number
  unit_price: number
  amount: number
  billing_period_start: string
  billing_period_end: string
}

// 账户信息
export interface Account {
  customer_id: number
  balance: number
  credit_limit: number
  status: 'active' | 'suspended' | 'overdue'
  currency: string
}

// 账单信息
export interface Invoice extends Timestamps {
  id: number
  customer_id: number
  billing_period: string
  total_amount: number
  status: 'pending' | 'paid' | 'overdue' | 'cancelled'
  due_date: string
  paid_at?: string
}

// 支付记录
export interface Payment extends Timestamps {
  id: number
  customer_id: number
  order_id: string
  amount: number
  payment_method: 'alipay' | 'wechat' | 'bank_transfer'
  status: 'pending' | 'success' | 'failed' | 'refunded'
  paid_at?: string
}
