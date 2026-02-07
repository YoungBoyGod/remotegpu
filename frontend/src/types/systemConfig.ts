/** 系统配置项 */
export interface SystemConfig {
  id: number
  config_key: string
  config_value: string
  config_type: 'string' | 'integer' | 'boolean' | 'json'
  config_group: string
  description: string
  is_public: boolean
  created_at: string
  updated_at: string
}

/** 配置分组 */
export interface SystemConfigGroup {
  label: string
  key: string
  configs: SystemConfig[]
}

/** 批量更新配置请求 */
export interface UpdateSystemConfigsPayload {
  configs: Record<string, string>
}
