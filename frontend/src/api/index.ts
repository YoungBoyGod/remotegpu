/**
 * API 统一导出
 *
 * 使用方式：
 * import { authApi, adminApi, customerApi } from '@/api'
 */

// 导出业务模块
export * as authApi from './auth'
export * as adminApi from './admin'
export * as customerApi from './customer'

// 导出公共类型
export * from './common/types'
export { default as request } from '@/utils/request'
