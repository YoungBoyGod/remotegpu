/**
 * API 统一导出
 *
 * 使用方式：
 * import { authApi, environmentApi, cmdbApi } from '@/api'
 *
 * 或者：
 * import * as api from '@/api'
 */

// 导出所有模块
export * as authApi from './auth'
export * as cmdbApi from './cmdb'
export * as environmentApi from './environment'
export * as schedulerApi from './scheduler'
export * as storageApi from './storage'
export * as imageApi from './image'
export * as billingApi from './billing'
export * as monitoringApi from './monitoring'
export * as artifactApi from './artifact'
export * as issueApi from './issue'
export * as requirementApi from './requirement'
export * as notificationApi from './notification'
export * as webhookApi from './webhook'

// 导出公共模块
export * from './common/types'
export { default as request } from './common/request'
