/**
 * Host Selection Module - Type Definitions
 */
import { Server } from '../cmdb/types'
import { PaginationParams } from '../common/types'

// Host availability status
export type HostAvailabilityStatus = 'available' | 'limited' | 'unavailable'

// Host information (extends Server with host-specific fields)
export interface Host extends Server {
  region: string
  gpu_model: string
  gpu_memory: number
  cuda_version: string
  price_per_hour: number
  availability_status: HostAvailabilityStatus
  available_gpu_count: number
}

// Host filter parameters
export interface HostFilterParams {
  region?: string
  gpu_count?: number | string
  gpu_model?: string
  keyword?: string
}

// Host pricing information
export interface HostPricing {
  base_price: number
  gpu_price: number
  total_price_per_hour: number
  currency: string
}

// Host query parameters (combines filters with pagination)
export interface HostQueryParams extends PaginationParams, HostFilterParams {}
