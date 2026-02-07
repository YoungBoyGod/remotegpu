export type TaskStatus =
  | 'pending'
  | 'queued'
  | 'assigned'
  | 'running'
  | 'completed'
  | 'failed'
  | 'cancelled'
  | 'stopped'
  | 'preempted'
  | 'suspended'

export interface Task {
  id: string
  name: string
  type: string
  status: TaskStatus | string
  customer_id?: number
  host_id?: string
  image_id?: number
  image?: {
    id?: number
    name?: string
  }
  command?: string
  env_vars?: Record<string, string> | string
  exit_code?: number
  error_msg?: string
  started_at?: string | null
  ended_at?: string | null
  created_at?: string
}

export interface TaskLogResponse {
  stdout?: string
  stderr?: string
  logs?: string
}

export interface TaskResultResponse {
  storage_type?: string
  presigned_url?: string
  url?: string
  expires_in?: number
  filename?: string
}
