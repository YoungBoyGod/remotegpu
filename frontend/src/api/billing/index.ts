/**
 * 计费管理模块 - API 接口
 */
import request from '@/utils/request'
import type { BillingRecord, Account, Invoice, Payment } from './types'
import type { PaginationParams, PaginationResponse, IdResponse } from '../common/types'

// ==================== 账户管理 ====================

/**
 * 获取账户余额
 */
export function getAccountBalance() {
  return request.get<Account>('/billing/account')
}

// ==================== 计费记录 ====================

/**
 * 获取计费记录
 */
export function getBillingRecords(params?: PaginationParams & {
  start_date?: string
  end_date?: string
  resource_type?: string
}) {
  return request.get<PaginationResponse<BillingRecord>>('/billing/records', { params })
}

// ==================== 账单管理 ====================

/**
 * 获取账单列表
 */
export function getInvoiceList(params?: PaginationParams & {
  status?: string
  start_date?: string
  end_date?: string
}) {
  return request.get<PaginationResponse<Invoice>>('/billing/invoices', { params })
}

/**
 * 获取账单详情
 */
export function getInvoiceDetail(id: number) {
  return request.get<Invoice & { items: BillingRecord[] }>(`/billing/invoices/${id}`)
}

/**
 * 下载账单
 */
export function downloadInvoice(id: number) {
  return request.get<{ download_url: string }>(`/billing/invoices/${id}/download`)
}

// ==================== 充值支付 ====================

/**
 * 创建充值订单
 */
export function createRechargeOrder(data: {
  amount: number
  payment_method: string
}) {
  return request.post<IdResponse & { payment_url: string }>('/billing/recharge', data)
}

/**
 * 查询支付状态
 */
export function getPaymentStatus(orderId: string) {
  return request.get<{ status: string; paid_at?: string }>(`/billing/payments/${orderId}`)
}

/**
 * 获取支付记录
 */
export function getPaymentHistory(params?: PaginationParams) {
  return request.get<PaginationResponse<Payment>>('/billing/payments', { params })
}
