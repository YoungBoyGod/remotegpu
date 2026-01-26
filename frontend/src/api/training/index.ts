/**
 * 训练与推理模块 - API 接口
 */
import request from '../common/request'
import type { TrainingJob, InferenceService, Experiment, ExperimentRun } from './types'
import type { PaginationParams, PaginationResponse, StatusResponse, IdResponse } from '../common/types'

// ==================== 训练任务 ====================

/**
 * 获取训练任务列表
 */
export function getTrainingJobList(params?: PaginationParams & {
  status?: string
  keyword?: string
}) {
  return request.get<PaginationResponse<TrainingJob>>('/training/jobs', { params })
}

/**
 * 创建训练任务
 */
export function createTrainingJob(data: {
  name: string
  image: string
  script: string
  datasets?: number[]
  resources: {
    cpu: number
    memory: number
    gpu: number
  }
  env_vars?: Record<string, string>
}) {
  return request.post<IdResponse & { status: string }>('/training/jobs', data)
}

/**
 * 获取训练任务详情
 */
export function getTrainingJobDetail(id: number) {
  return request.get<TrainingJob>(`/training/jobs/${id}`)
}

/**
 * 停止训练任务
 */
export function stopTrainingJob(id: number) {
  return request.post<StatusResponse>(`/training/jobs/${id}/stop`)
}

/**
 * 获取训练日志
 */
export function getTrainingLogs(id: number, params?: { tail?: number }) {
  return request.get<{ logs: string }>(`/training/jobs/${id}/logs`, { params })
}

// ==================== 推理服务 ====================

/**
 * 获取推理服务列表
 */
export function getInferenceServiceList(params?: PaginationParams) {
  return request.get<PaginationResponse<InferenceService>>('/inference/services', { params })
}

/**
 * 部署推理服务
 */
export function deployInferenceService(data: {
  name: string
  model_id: number
  replicas: number
  resources: {
    cpu: number
    memory: number
    gpu?: number
  }
}) {
  return request.post<IdResponse & { endpoint: string }>('/inference/services', data)
}

/**
 * 获取推理服务详情
 */
export function getInferenceServiceDetail(id: number) {
  return request.get<InferenceService>(`/inference/services/${id}`)
}

/**
 * 停止推理服务
 */
export function stopInferenceService(id: number) {
  return request.post<StatusResponse>(`/inference/services/${id}/stop`)
}

/**
 * 调用推理接口
 */
export function predict(serviceId: number, data: { input_data: any }) {
  return request.post<{ predictions: any }>(`/inference/services/${serviceId}/predict`, data)
}
