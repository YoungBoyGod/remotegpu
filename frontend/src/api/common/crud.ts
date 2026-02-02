import request from '../request'

/**
 * 分页参数
 */
export interface PaginationParams {
  page?: number
  page_size?: number
}

/**
 * 分页响应
 */
export interface PaginationResponse<T> {
  data: T[]
  total: number
  page: number
  page_size: number
}

/**
 * 状态响应
 */
export interface StatusResponse {
  message: string
}

/**
 * ID响应
 */
export interface IdResponse {
  id: number | string
}

/**
 * CRUD API 配置
 */
export interface CrudConfig {
  baseUrl: string
  listEndpoint?: string
  detailEndpoint?: string
  createEndpoint?: string
  updateEndpoint?: string
  deleteEndpoint?: string
}

/**
 * 创建 CRUD API
 * @param config API 配置
 * @returns CRUD API 对象
 */
export function createCrudApi<T, CreateReq = Partial<T>, UpdateReq = Partial<T>>(
  config: CrudConfig
) {
  const {
    baseUrl,
    listEndpoint = baseUrl,
    detailEndpoint = baseUrl,
    createEndpoint = baseUrl,
    updateEndpoint = baseUrl,
    deleteEndpoint = baseUrl,
  } = config

  return {
    /**
     * 获取列表（分页）
     */
    list: (page: number = 1, pageSize: number = 10) =>
      request.get<PaginationResponse<T>>(listEndpoint, {
        params: { page, page_size: pageSize },
      }),

    /**
     * 根据ID获取详情
     */
    getById: (id: number | string) =>
      request.get<T>(`${detailEndpoint}/${id}`),

    /**
     * 创建
     */
    create: (data: CreateReq) =>
      request.post<T>(createEndpoint, data),

    /**
     * 更新
     */
    update: (id: number | string, data: UpdateReq) =>
      request.put<T>(`${updateEndpoint}/${id}`, data),

    /**
     * 删除
     */
    delete: (id: number | string) =>
      request.delete<StatusResponse>(`${deleteEndpoint}/${id}`),
  }
}
