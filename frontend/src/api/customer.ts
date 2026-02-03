import request from '@/utils/request'
import type { ApiResponse, PageRequest, PageResponse } from '@/types/common'
import type { Machine, MonitoringData } from '@/types/machine'

// ==================== Dashboard ====================

/**
 * 获取客户Dashboard概览数据
 */
export function getDashboardOverview(): Promise<ApiResponse<any>> {
  return request.get('/api/customer/dashboard/stats')
}

// ==================== 我的机器 ====================

/**
 * 获取我的机器列表
 */
export function getMyMachines(params?: PageRequest): Promise<ApiResponse<PageResponse<Machine>>> {
  return request.get('/api/customer/machines', { params })
}

/**
 * 获取机器详情
 */
export function getMachineDetail(id: number): Promise<ApiResponse<Machine>> {
  return request.get(`/api/customer/machines/${id}`)
}

/**
 * 获取机器监控数据
 */
export function getMachineMonitoring(id: number): Promise<ApiResponse<MonitoringData>> {
  return request.get(`/api/customer/machines/${id}/monitoring`)
}

/**
 * 获取连接历史
 */
export function getConnectHistory(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/customer/machines/connect-history', { params })
}

/**
 * 快速连接机器
 */
export function quickConnect(machineId: number, protocol: 'ssh' | 'rdp' | 'vnc'): Promise<ApiResponse<any>> {
  return request.post('/api/customer/machines/quick-connect', { machineId, protocol })
}

// ==================== 任务管理 ====================

/**
 * 获取训练任务列表
 */
export function getTrainingTasks(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/customer/tasks/training', { params })
}

/**
 * 创建训练任务
 */
export function createTrainingTask(data: any): Promise<ApiResponse<any>> {
  return request.post('/api/customer/tasks/training', data)
}

/**
 * 获取推理任务列表
 */
export function getInferenceTasks(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/customer/tasks/inference', { params })
}

/**
 * 创建推理任务
 */
export function createInferenceTask(data: any): Promise<ApiResponse<any>> {
  return request.post('/api/customer/tasks/inference', data)
}

/**
 * 获取任务队列
 */
export function getTaskQueue(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/customer/tasks/queue', { params })
}

/**
 * 获取任务历史
 */
export function getTaskHistory(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/customer/tasks/history', { params })
}

/**
 * 获取任务详情
 */
export function getTaskDetail(id: number): Promise<ApiResponse<any>> {
  return request.get(`/api/customer/tasks/${id}`)
}

/**
 * 取消任务
 */
export function cancelTask(id: number): Promise<ApiResponse<void>> {
  return request.post(`/api/customer/tasks/${id}/stop`)
}

/**
 * 获取任务日志
 */
export function getTaskLogs(id: number): Promise<ApiResponse<string>> {
  return request.get(`/api/customer/tasks/${id}/logs`)
}

// ==================== 镜像市场 ====================

/**
 * 获取镜像市场列表
 */
export function getImageMarket(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/customer/images/market', { params })
}

/**
 * 获取我的镜像列表
 */
export function getMyImages(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/customer/images/my', { params })
}

/**
 * 构建镜像
 */
export function buildImage(data: any): Promise<ApiResponse<any>> {
  return request.post('/api/customer/images/build', data)
}

/**
 * 删除镜像
 */
export function deleteImage(id: number): Promise<ApiResponse<void>> {
  return request.delete(`/api/customer/images/${id}`)
}

// ==================== 数据集管理 ====================

/**
 * 获取数据集列表
 */
export function getDatasetList(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/customer/datasets', { params })
}

/**
 * 上传数据集
 */
export function uploadDataset(data: FormData, onProgress?: (progress: number) => void): Promise<ApiResponse<any>> {
  return request.post('/api/customer/datasets/upload', data, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: (progressEvent) => {
      if (onProgress && progressEvent.total) {
        const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
        onProgress(progress)
      }
    }
  })
}

/**
 * 删除数据集
 */
export function deleteDataset(id: number): Promise<ApiResponse<void>> {
  return request.delete(`/api/customer/datasets/${id}`)
}

/**
 * 获取挂载列表
 */
export function getMountList(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/customer/datasets/mounts', { params })
}

/**
 * 挂载数据集
 */
export function mountDataset(data: { datasetId: number; machineId: number; mountPath: string }): Promise<ApiResponse<any>> {
  return request.post('/api/customer/datasets/mount', data)
}

/**
 * 卸载数据集
 */
export function unmountDataset(id: number): Promise<ApiResponse<void>> {
  return request.post(`/api/customer/datasets/unmount/${id}`)
}

// ==================== 模型管理 ====================

/**
 * 获取模型列表
 */
export function getModelList(params: PageRequest): Promise<ApiResponse<PageResponse<any>>> {
  return request.get('/api/customer/models', { params })
}

/**
 * 上传模型
 */
export function uploadModel(data: FormData): Promise<ApiResponse<any>> {
  return request.post('/api/customer/models/upload', data, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

/**
 * 删除模型
 */
export function deleteModel(id: number): Promise<ApiResponse<void>> {
  return request.delete(`/api/customer/models/${id}`)
}

// ==================== 监控与分析 ====================

/**
 * 获取实时监控数据
 */
export function getRealtimeMonitoring(machineId?: number): Promise<ApiResponse<any>> {
  return request.get('/api/customer/monitoring/realtime', { params: { machineId } })
}

/**
 * 获取使用统计数据
 */
export function getUsageStats(params?: { startDate?: string; endDate?: string }): Promise<ApiResponse<any>> {
  return request.get('/api/customer/statistics/usage', { params })
}

// ==================== 文件管理 ====================

/**
 * 浏览文件
 */
export function browseFiles(path: string): Promise<ApiResponse<any>> {
  return request.get('/api/customer/files/browse', { params: { path } })
}

/**
 * 上传文件
 */
export function uploadFile(data: FormData): Promise<ApiResponse<any>> {
  return request.post('/api/customer/files/upload', data, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

/**
 * 下载文件
 */
export function downloadFile(path: string): Promise<Blob> {
  return request.get('/api/customer/files/download', {
    params: { path },
    responseType: 'blob'
  })
}

/**
 * 删除文件
 */
export function deleteFile(path: string): Promise<ApiResponse<void>> {
  return request.delete('/api/customer/files', { params: { path } })
}

// ==================== 开发工具 ====================

/**
 * 启动Jupyter
 */
export function startJupyter(machineId: number): Promise<ApiResponse<{ url: string; token: string }>> {
  return request.post('/api/customer/tools/jupyter/start', { machineId })
}

/**
 * 停止Jupyter
 */
export function stopJupyter(machineId: number): Promise<ApiResponse<void>> {
  return request.post('/api/customer/tools/jupyter/stop', { machineId })
}

/**
 * 启动TensorBoard
 */
export function startTensorBoard(machineId: number, logDir: string): Promise<ApiResponse<{ url: string }>> {
  return request.post('/api/customer/tools/tensorboard/start', { machineId, logDir })
}

/**
 * 停止TensorBoard
 */
export function stopTensorBoard(machineId: number): Promise<ApiResponse<void>> {
  return request.post('/api/customer/tools/tensorboard/stop', { machineId })
}

/**
 * 获取Terminal连接信息
 */
export function getTerminalConnection(machineId: number): Promise<ApiResponse<{ wsUrl: string; token: string }>> {
  return request.get(`/api/customer/tools/terminal/${machineId}/connection`)
}

// ==================== 设置 ====================

/**
 * 获取个人设置
 */
export function getProfile(): Promise<ApiResponse<any>> {
  return request.get('/api/customer/settings/profile')
}

/**
 * 更新个人设置
 */
export function updateProfile(data: any): Promise<ApiResponse<any>> {
  return request.put('/api/customer/settings/profile', data)
}

/**
 * 获取SSH密钥列表
 */
export function getSshKeys(): Promise<ApiResponse<any[]>> {
  return request.get('/api/customer/settings/ssh-keys')
}

/**
 * 添加SSH密钥
 */
export function addSshKey(data: { name: string; publicKey: string }): Promise<ApiResponse<any>> {
  return request.post('/api/customer/settings/ssh-keys', data)
}

/**
 * 删除SSH密钥
 */
export function deleteSshKey(id: number): Promise<ApiResponse<void>> {
  return request.delete(`/api/customer/settings/ssh-keys/${id}`)
}

/**
 * 获取通知设置
 */
export function getNotificationSettings(): Promise<ApiResponse<any>> {
  return request.get('/api/customer/settings/notifications')
}

/**
 * 更新通知设置
 */
export function updateNotificationSettings(data: any): Promise<ApiResponse<any>> {
  return request.put('/api/customer/settings/notifications', data)
}
