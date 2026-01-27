<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/common/PageHeader.vue'

interface Device {
  id: number
  name: string
  region: string
  gpuModel: string
  gpuMemory: number
  gpuCount: number
  availableGpu: number
  totalGpu: number
  cpu: string
  memory: number
  disk: number
  systemDisk: number
  cudaVersion: string
  gpuDriver: string
  pricePerHour: number
  expiryDate: string
}

const loading = ref(false)
const devices = ref<Device[]>([])

// 筛选条件
const billingMethod = ref('按需计费')
const selectedRegion = ref('北京B区')
const selectedGpuModels = ref<string[]>([])
const selectedGpuCount = ref(1)

// 地区选项
const regions = [
  { label: '北京B区', value: '北京B区', hot: true },
  { label: '西北B区', value: '西北B区' },
  { label: '重庆A区', value: '重庆A区' },
  { label: '内蒙B区', value: '内蒙B区' },
  { label: '北京A区', value: '北京A区' },
  { label: '佛山区', value: '佛山区' }
]

// GPU型号选项
const gpuModels = [
  { label: '全部', value: 'all' },
  { label: 'RTX 5090', value: 'RTX 5090', count: 1114 },
  { label: 'RTX PRO 6000', value: 'RTX PRO 6000', count: 0 },
  { label: 'vGPU-48GB', value: 'vGPU-48GB', count: 31 },
  { label: 'vGPU-48GB-425W', value: 'vGPU-48GB-425W', count: 89 },
  { label: 'RTX 5090 D', value: 'RTX 5090 D', count: 1 },
  { label: 'RTX 4090D', value: 'RTX 4090D', count: 1 },
  { label: 'RTX 4090', value: 'RTX 4090', count: 87 },
  { label: 'CPU', value: 'CPU', count: 0 }
]

// 过滤后的设备列表
const filteredDevices = computed(() => {
  let result = devices.value

  if (selectedRegion.value) {
    result = result.filter(d => d.region === selectedRegion.value)
  }

  if (selectedGpuModels.value.length > 0 && !selectedGpuModels.value.includes('all')) {
    result = result.filter(d => selectedGpuModels.value.includes(d.gpuModel))
  }

  if (selectedGpuCount.value > 0) {
    result = result.filter(d => d.gpuCount >= selectedGpuCount.value)
  }

  return result
})

// 加载设备列表
const loadDevices = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 500))
    devices.value = [
      {
        id: 1,
        name: '北京B区 / 598机',
        region: '北京B区',
        gpuModel: 'RTX 5090',
        gpuMemory: 32,
        gpuCount: 8,
        availableGpu: 1,
        totalGpu: 8,
        cpu: '16核, Xeon(R) Gold 6459C',
        memory: 90,
        disk: 50,
        systemDisk: 5881,
        cudaVersion: '≤ 13.0',
        gpuDriver: '580.76.05',
        pricePerHour: 3.03,
        expiryDate: '2027-01-01'
      },
      {
        id: 2,
        name: '北京B区 / 353机',
        region: '北京B区',
        gpuModel: 'RTX 5090',
        gpuMemory: 32,
        gpuCount: 8,
        availableGpu: 1,
        totalGpu: 8,
        cpu: '25核, Xeon(R) Platinum 8470Q',
        memory: 90,
        disk: 50,
        systemDisk: 3810,
        cudaVersion: '≤ 13.0',
        gpuDriver: '580.76.05',
        pricePerHour: 3.03,
        expiryDate: '2027-01-01'
      }
    ]
  } catch (error) {
    ElMessage.error('加载设备列表失败')
  } finally {
    loading.value = false
  }
}

// 租用设备
const handleRent = (device: Device) => {
  ElMessage.success(`正在租用设备: ${device.name}`)
  // TODO: 实现租用逻辑
}

onMounted(() => {
  loadDevices()
})
</script>

<template>
  <div class="computing-market">
    <PageHeader title="算力市场" />

    <div class="market-container">
      <!-- 筛选区域 -->
      <div class="filter-section">
        <!-- 计费方式 -->
        <div class="filter-row">
          <span class="filter-label">计费方式</span>
          <el-radio-group v-model="billingMethod">
            <el-radio-button label="按需计费" />
            <el-radio-button label="包日" />
            <el-radio-button label="包周" />
            <el-radio-button label="包月" />
          </el-radio-group>
        </div>

        <!-- 选择地区 -->
        <div class="filter-row">
          <span class="filter-label">选择地区</span>
          <div class="region-buttons">
            <el-button
              v-for="region in regions"
              :key="region.value"
              :type="selectedRegion === region.value ? 'primary' : 'default'"
              size="small"
              @click="selectedRegion = region.value"
            >
              {{ region.label }}
              <el-tag v-if="region.hot" type="danger" size="small" style="margin-left: 4px">
                PRO9000
              </el-tag>
            </el-button>
          </div>
        </div>

        <!-- GPU型号 -->
        <div class="filter-row">
          <span class="filter-label">GPU型号</span>
          <el-checkbox-group v-model="selectedGpuModels">
            <el-checkbox
              v-for="model in gpuModels"
              :key="model.value"
              :label="model.value"
            >
              {{ model.label }}
              <span v-if="model.count !== undefined" class="model-count">
                ({{ model.count }}/{{ model.count }})
              </span>
            </el-checkbox>
          </el-checkbox-group>
        </div>

        <!-- GPU数量 -->
        <div class="filter-row">
          <span class="filter-label">GPU数量</span>
          <el-radio-group v-model="selectedGpuCount">
            <el-radio-button :label="1">1</el-radio-button>
            <el-radio-button :label="2">2</el-radio-button>
            <el-radio-button :label="3">3</el-radio-button>
            <el-radio-button :label="4">4</el-radio-button>
            <el-radio-button :label="5">5</el-radio-button>
            <el-radio-button :label="6">6</el-radio-button>
            <el-radio-button :label="7">7</el-radio-button>
            <el-radio-button :label="8">8</el-radio-button>
            <el-radio-button :label="10">10</el-radio-button>
            <el-radio-button :label="12">12</el-radio-button>
          </el-radio-group>
        </div>
      </div>

      <!-- 设备列表 -->
      <div v-loading="loading" class="device-list">
        <div
          v-for="device in filteredDevices"
          :key="device.id"
          class="device-card"
        >
          <div class="device-header">
            <div class="device-title">
              <span class="device-region">{{ device.region }}</span>
              <span class="device-name">/ {{ device.name }}</span>
              <span class="device-expiry">可用期限：{{ device.expiryDate }}</span>
            </div>
            <div class="device-availability">
              <span class="availability-label">空闲/总量</span>
              <span class="availability-value">{{ device.availableGpu }}/{{ device.totalGpu }}</span>
            </div>
          </div>

          <div class="device-title-main">
            {{ device.gpuModel }} / {{ device.gpuMemory }} GB
          </div>

          <div class="device-specs">
            <div class="spec-column">
              <div class="spec-title">每GPU分配</div>
              <div class="spec-item">
                <span class="spec-label">CPU:</span>
                <span class="spec-value">{{ device.cpu }}</span>
              </div>
              <div class="spec-item">
                <span class="spec-label">内存:</span>
                <span class="spec-value">{{ device.memory }} GB</span>
              </div>
            </div>

            <div class="spec-column">
              <div class="spec-title">硬盘</div>
              <div class="spec-item">
                <span class="spec-label">系统盘:</span>
                <span class="spec-value">{{ device.disk }} GB</span>
              </div>
              <div class="spec-item">
                <span class="spec-label">数据盘:</span>
                <span class="spec-value">{{ device.systemDisk }} GB</span>
              </div>
            </div>

            <div class="spec-column">
              <div class="spec-title">其它</div>
              <div class="spec-item">
                <span class="spec-label">GPU驱动:</span>
                <span class="spec-value">{{ device.gpuDriver }}</span>
              </div>
              <div class="spec-item">
                <span class="spec-label">CUDA版本:</span>
                <span class="spec-value">{{ device.cudaVersion }}</span>
              </div>
            </div>

            <div class="spec-column price-column">
              <div class="price">
                <span class="price-symbol">¥</span>
                <span class="price-value">{{ device.pricePerHour.toFixed(2) }}</span>
                <span class="price-unit">/时</span>
              </div>
              <div class="price-tip">会员低至7.9折 ¥2.39/时</div>
              <el-button type="primary" size="large" @click="handleRent(device)">
                1卡即租
              </el-button>
            </div>
          </div>
        </div>

        <el-empty
          v-if="!loading && filteredDevices.length === 0"
          description="暂无符合条件的设备"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.computing-market {
  padding: 24px;
}

.market-container {
  background: white;
  border-radius: 8px;
  padding: 24px;
}

.filter-section {
  margin-bottom: 24px;
}

.filter-row {
  display: flex;
  align-items: flex-start;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.filter-row:last-child {
  border-bottom: none;
}

.filter-label {
  min-width: 100px;
  font-weight: 500;
  color: #303133;
  padding-top: 8px;
}

.region-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.model-count {
  color: #909399;
  font-size: 12px;
}

.device-list {
  min-height: 400px;
}

.device-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 16px;
}

.device-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.device-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #606266;
}

.device-region {
  color: #409EFF;
}

.device-expiry {
  color: #909399;
}

.device-availability {
  display: flex;
  align-items: center;
  gap: 8px;
}

.availability-label {
  font-size: 13px;
  color: #909399;
}

.availability-value {
  font-size: 18px;
  font-weight: 600;
  color: #409EFF;
}

.device-title-main {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 16px;
}

.device-specs {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
}

.spec-column {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.spec-title {
  font-size: 13px;
  color: #909399;
  margin-bottom: 4px;
}

.spec-item {
  display: flex;
  gap: 8px;
  font-size: 13px;
}

.spec-label {
  color: #606266;
}

.spec-value {
  color: #303133;
}

.price-column {
  align-items: flex-end;
  text-align: right;
}

.price {
  display: flex;
  align-items: baseline;
  justify-content: flex-end;
  margin-bottom: 4px;
}

.price-symbol {
  font-size: 16px;
  color: #F56C6C;
}

.price-value {
  font-size: 28px;
  font-weight: 600;
  color: #F56C6C;
}

.price-unit {
  font-size: 14px;
  color: #F56C6C;
}

.price-tip {
  font-size: 12px;
  color: #909399;
  margin-bottom: 12px;
}

@media (max-width: 1200px) {
  .device-specs {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
