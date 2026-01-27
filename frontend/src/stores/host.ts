import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Host } from '@/api/host/types'

export const useHostStore = defineStore('host', () => {
  // 选中的主机
  const selectedHost = ref<Host | null>(null)

  // 可用地区列表
  const availableRegions = ref<string[]>([])

  // 可用GPU型号列表
  const availableGpuModels = ref<string[]>([])

  // 选择主机
  const selectHost = (host: Host) => {
    selectedHost.value = host
  }

  // 清除选择
  const clearSelection = () => {
    selectedHost.value = null
  }

  // 设置可用地区
  const setAvailableRegions = (regions: string[]) => {
    availableRegions.value = regions
  }

  // 设置可用GPU型号
  const setAvailableGpuModels = (models: string[]) => {
    availableGpuModels.value = models
  }

  return {
    selectedHost,
    availableRegions,
    availableGpuModels,
    selectHost,
    clearSelection,
    setAvailableRegions,
    setAvailableGpuModels,
  }
})
