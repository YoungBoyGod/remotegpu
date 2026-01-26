/**
 * 训练与推理模块 - 类型定义
 */
import { Timestamps, UUIDField } from '../common/types'

// 训练任务
export interface TrainingJob extends Timestamps, UUIDField {
  id: number
  customer_id: number
  workspace_id?: number
  name: string
  description?: string
  image: string
  script: string
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled'
  resources: {
    cpu: number
    memory: number
    gpu: number
  }
  datasets?: number[]
  started_at?: string
  completed_at?: string
  duration?: number
}

// 推理服务
export interface InferenceService extends Timestamps, UUIDField {
  id: number
  customer_id: number
  model_id: number
  name: string
  image: string
  replicas: number
  status: 'deploying' | 'running' | 'stopped' | 'failed'
  endpoint?: string
  resources: {
    cpu: number
    memory: number
    gpu?: number
  }
}

// 实验
export interface Experiment extends Timestamps {
  id: number
  customer_id: number
  name: string
  description?: string
  hyperparameters?: Record<string, any>
}

// 实验运行记录
export interface ExperimentRun extends Timestamps {
  id: number
  experiment_id: number
  run_number: number
  status: string
  metrics?: Record<string, number>
  artifacts?: string[]
}
