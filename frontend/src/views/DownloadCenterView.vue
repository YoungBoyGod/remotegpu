<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Download, Search, Refresh } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'

interface DownloadResource {
  id: number
  name: string
  description: string
  chipType: string
  aiModelType: string
  boardType: string
  productSeries: string
  serverModel: string
  platform: string
  os: string
  version: string
  versionType: 'stable' | 'beta'
  fileSize: number
  md5: string
  publishTime: string
  downloadUrl: string
}

const loading = ref(false)
const resources = ref<DownloadResource[]>([])

// 筛选条件
const selectedChipType = ref<string>('')
const selectedAiModelType = ref<string>('')
const selectedProductSeries = ref<string>('')
const selectedServerModel = ref<string>('')
const selectedPlatform = ref<string>('')
const selectedOs = ref<string>('')
const selectedVersionType = ref<'stable' | 'beta'>('stable')
const searchKeyword = ref('')

// 芯片类型列表
const chipTypes = ref([
  { label: 'AMD', value: 'amd', icon: '🔴' },
  { label: 'Intel', value: 'intel', icon: '🔵' },
  { label: '英伟达', value: 'nvidia', icon: '🟢' },
  { label: '华为昇腾', value: 'huawei-ascend', icon: '🟠' }
])

// AI模型类别
const aiModelTypes = ref([
  { label: '全部', value: '' },
  { label: 'LLM大语言模型', value: 'llm' },
  { label: '计算机视觉', value: 'cv' },
  { label: '语音识别', value: 'asr' },
  { label: '推荐系统', value: 'rec' }
])

// 平台架构
const platforms = ref([
  { label: '全部', value: '' },
  { label: 'x86_64', value: 'x86_64' },
  { label: 'ARM64', value: 'arm64' },
  { label: 'RISC-V', value: 'riscv' }
])

// 操作系统
const osList = ref([
  { label: '全部', value: '' },
  { label: 'Ubuntu 20.04', value: 'ubuntu-20.04' },
  { label: 'Ubuntu 22.04', value: 'ubuntu-22.04' },
  { label: 'CentOS 7', value: 'centos-7' },
  { label: 'CentOS 8', value: 'centos-8' }
])

// 产品系列
const productSeriesList = ref([
  { label: '全部', value: '' },
  { label: 'Radeon RX 7000系列', value: 'rx-7000' },
  { label: 'Radeon RX 6000系列', value: 'rx-6000' },
  { label: 'Xeon可扩展处理器', value: 'xeon-scalable' },
  { label: 'Core系列', value: 'core' },
  { label: 'GeForce RTX 40系列', value: 'rtx-40' },
  { label: 'GeForce RTX 30系列', value: 'rtx-30' }
])

// 已适配服务器型号
const serverModelsList = ref([
  { label: '全部', value: '' },
  { label: 'Dell PowerEdge R750', value: 'dell-r750' },
  { label: 'HP ProLiant DL380', value: 'hp-dl380' },
  { label: 'Lenovo ThinkSystem SR650', value: 'lenovo-sr650' },
  { label: 'Supermicro SYS-420GP', value: 'supermicro-420gp' },
  { label: '浪潮NF5280M6', value: 'inspur-nf5280m6' },
  { label: '华为FusionServer 2288H', value: 'huawei-2288h' }
])

// 过滤后的资源列表
const filteredResources = computed(() => {
  let result = resources.value

  // 芯片类型筛选
  if (selectedChipType.value) {
    result = result.filter(r => r.chipType === selectedChipType.value)
  }

  // AI模型类别筛选
  if (selectedAiModelType.value) {
    result = result.filter(r => r.aiModelType === selectedAiModelType.value)
  }

  // 产品系列筛选
  if (selectedProductSeries.value) {
    result = result.filter(r => r.productSeries === selectedProductSeries.value)
  }

  // 服务器型号筛选
  if (selectedServerModel.value) {
    result = result.filter(r => r.serverModel === selectedServerModel.value)
  }

  // 平台架构筛选
  if (selectedPlatform.value) {
    result = result.filter(r => r.platform === selectedPlatform.value)
  }

  // 操作系统筛选
  if (selectedOs.value) {
    result = result.filter(r => r.os === selectedOs.value)
  }

  // 版本类型筛选
  result = result.filter(r => r.versionType === selectedVersionType.value)

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(r =>
      r.name.toLowerCase().includes(keyword) ||
      r.description.toLowerCase().includes(keyword)
    )
  }

  return result
})

// 是否显示空状态提示
const showEmptyHint = computed(() => {
  return !selectedChipType.value
})

// 加载资源列表
const loadResources = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 500))
    resources.value = [
      {
        id: 1,
        name: 'AMD驱动程序',
        description: 'AMD芯片驱动程序 v2.5.0',
        chipType: 'amd',
        aiModelType: 'llm',
        boardType: 'AMD-RX7900',
        productSeries: 'rx-7000',
        serverModel: 'dell-r750',
        platform: 'x86_64',
        os: 'ubuntu-20.04',
        version: 'v2.5.0',
        versionType: 'stable',
        fileSize: 256.5,
        md5: 'a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6',
        publishTime: '2026-01-20',
        downloadUrl: '#'
      },
      {
        id: 2,
        name: 'Intel SDK',
        description: 'Intel开发工具包 v1.8.3',
        chipType: 'intel',
        aiModelType: 'cv',
        boardType: 'Intel-Xeon',
        productSeries: 'xeon-scalable',
        serverModel: 'hp-dl380',
        platform: 'x86_64',
        os: 'ubuntu-22.04',
        version: 'v1.8.3',
        versionType: 'stable',
        fileSize: 512.8,
        md5: 'b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7',
        publishTime: '2026-01-18',
        downloadUrl: '#'
      },
      {
        id: 3,
        name: 'AMD测试版驱动',
        description: 'AMD芯片驱动程序 v2.6.0-beta',
        chipType: 'amd',
        aiModelType: 'llm',
        boardType: 'AMD-RX7900',
        productSeries: 'rx-7000',
        serverModel: 'supermicro-420gp',
        platform: 'x86_64',
        os: 'ubuntu-22.04',
        version: 'v2.6.0-beta',
        versionType: 'beta',
        fileSize: 268.3,
        md5: 'c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8',
        publishTime: '2026-01-25',
        downloadUrl: '#'
      },
      {
        id: 4,
        name: 'NVIDIA CUDA工具包',
        description: 'NVIDIA CUDA开发工具包 v12.3',
        chipType: 'nvidia',
        aiModelType: 'llm',
        boardType: 'RTX-4090',
        productSeries: 'rtx-40',
        serverModel: 'lenovo-sr650',
        platform: 'x86_64',
        os: 'ubuntu-22.04',
        version: 'v12.3',
        versionType: 'stable',
        fileSize: 3584.2,
        md5: 'd4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9',
        publishTime: '2026-01-22',
        downloadUrl: '#'
      },
      {
        id: 5,
        name: '华为昇腾驱动',
        description: '华为昇腾910B驱动程序 v3.1.0',
        chipType: 'huawei-ascend',
        aiModelType: 'cv',
        boardType: 'Ascend-910B',
        productSeries: 'rx-7000',
        serverModel: 'huawei-2288h',
        platform: 'arm64',
        os: 'centos-8',
        version: 'v3.1.0',
        versionType: 'stable',
        fileSize: 892.6,
        md5: 'e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0',
        publishTime: '2026-01-19',
        downloadUrl: '#'
      },
      {
        id: 6,
        name: 'Intel语音识别SDK',
        description: 'Intel语音识别开发工具包 v2.0.1',
        chipType: 'intel',
        aiModelType: 'asr',
        boardType: 'Intel-Core-i9',
        productSeries: 'core',
        serverModel: 'dell-r750',
        platform: 'x86_64',
        os: 'ubuntu-20.04',
        version: 'v2.0.1',
        versionType: 'stable',
        fileSize: 428.9,
        md5: 'f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1',
        publishTime: '2026-01-17',
        downloadUrl: '#'
      },
      {
        id: 7,
        name: 'NVIDIA推荐系统框架',
        description: 'NVIDIA推荐系统优化框架 v1.5.2',
        chipType: 'nvidia',
        aiModelType: 'rec',
        boardType: 'RTX-3090',
        productSeries: 'rtx-30',
        serverModel: 'inspur-nf5280m6',
        platform: 'x86_64',
        os: 'centos-7',
        version: 'v1.5.2',
        versionType: 'stable',
        fileSize: 1256.4,
        md5: 'g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2',
        publishTime: '2026-01-16',
        downloadUrl: '#'
      },
      {
        id: 8,
        name: 'AMD RX 6000驱动',
        description: 'AMD Radeon RX 6000系列驱动程序 v2.3.1',
        chipType: 'amd',
        aiModelType: 'cv',
        boardType: 'AMD-RX6900XT',
        productSeries: 'rx-6000',
        serverModel: 'hp-dl380',
        platform: 'x86_64',
        os: 'ubuntu-22.04',
        version: 'v2.3.1',
        versionType: 'stable',
        fileSize: 234.7,
        md5: 'h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3',
        publishTime: '2026-01-15',
        downloadUrl: '#'
      },
      {
        id: 9,
        name: '华为昇腾测试版SDK',
        description: '华为昇腾AI开发工具包 v3.2.0-beta',
        chipType: 'huawei-ascend',
        aiModelType: 'llm',
        boardType: 'Ascend-910B',
        productSeries: 'rx-7000',
        serverModel: 'huawei-2288h',
        platform: 'arm64',
        os: 'centos-8',
        version: 'v3.2.0-beta',
        versionType: 'beta',
        fileSize: 1024.5,
        md5: 'i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4',
        publishTime: '2026-01-26',
        downloadUrl: '#'
      },
      {
        id: 10,
        name: 'Intel Core优化库',
        description: 'Intel Core系列性能优化库 v4.2.0',
        chipType: 'intel',
        aiModelType: 'cv',
        boardType: 'Intel-Core-i7',
        productSeries: 'core',
        serverModel: 'lenovo-sr650',
        platform: 'x86_64',
        os: 'ubuntu-20.04',
        version: 'v4.2.0',
        versionType: 'stable',
        fileSize: 356.8,
        md5: 'j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5',
        publishTime: '2026-01-14',
        downloadUrl: '#'
      },
      {
        id: 11,
        name: 'NVIDIA RTX 4090驱动',
        description: 'NVIDIA GeForce RTX 4090显卡驱动 v545.29',
        chipType: 'nvidia',
        aiModelType: 'llm',
        boardType: 'RTX-4090',
        productSeries: 'rtx-40',
        serverModel: 'supermicro-420gp',
        platform: 'x86_64',
        os: 'ubuntu-22.04',
        version: 'v545.29',
        versionType: 'stable',
        fileSize: 678.3,
        md5: 'k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6',
        publishTime: '2026-01-23',
        downloadUrl: '#'
      },
      {
        id: 12,
        name: 'AMD语音识别加速包',
        description: 'AMD语音识别硬件加速包 v1.9.0',
        chipType: 'amd',
        aiModelType: 'asr',
        boardType: 'AMD-RX7900',
        productSeries: 'rx-7000',
        serverModel: 'dell-r750',
        platform: 'x86_64',
        os: 'centos-7',
        version: 'v1.9.0',
        versionType: 'stable',
        fileSize: 445.2,
        md5: 'l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7',
        publishTime: '2026-01-13',
        downloadUrl: '#'
      },
      {
        id: 13,
        name: 'NVIDIA测试版CUDA',
        description: 'NVIDIA CUDA工具包 v12.4-beta',
        chipType: 'nvidia',
        aiModelType: 'cv',
        boardType: 'RTX-4090',
        productSeries: 'rtx-40',
        serverModel: 'hp-dl380',
        platform: 'x86_64',
        os: 'ubuntu-22.04',
        version: 'v12.4-beta',
        versionType: 'beta',
        fileSize: 3698.5,
        md5: 'm3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8',
        publishTime: '2026-01-27',
        downloadUrl: '#'
      }
    ]
  } catch (error) {
    ElMessage.error('加载资源列表失败')
  } finally {
    loading.value = false
  }
}

// 格式化文件大小
const formatSize = (size: number) => {
  if (size < 1024) return `${size.toFixed(2)} MB`
  return `${(size / 1024).toFixed(2)} GB`
}

// 选择芯片类型
const selectChipType = (chipType: string) => {
  selectedChipType.value = chipType
}

// 下载资源
const handleDownload = (resource: DownloadResource) => {
  ElMessage.success(`开始下载: ${resource.name}`)
  // TODO: 实现下载逻辑
}

// 清除筛选
const clearFilters = () => {
  selectedChipType.value = ''
  selectedAiModelType.value = ''
  selectedProductSeries.value = ''
  selectedServerModel.value = ''
  selectedPlatform.value = ''
  selectedOs.value = ''
  searchKeyword.value = ''
}

onMounted(() => {
  loadResources()
})
</script>

<template>
  <div class="download-center">
    <PageHeader title="下载中心" />

    <!-- 芯片类型选择 -->
    <el-card class="chip-selector-card">
      <div class="chip-types">
        <div
          v-for="chip in chipTypes"
          :key="chip.value"
          class="chip-type-item"
          :class="{ active: selectedChipType === chip.value }"
          @click="selectChipType(chip.value)"
        >
          <div class="chip-icon">{{ chip.icon }}</div>
          <div class="chip-label">{{ chip.label }}</div>
        </div>
      </div>
    </el-card>

    <!-- 筛选区域 -->
    <el-card class="filter-card">
      <div class="filter-header">
        <span class="filter-title">
          请先选择数级类别 <span class="required">*（必选）</span>
        </span>
      </div>
      <div class="filter-container">
        <div class="filter-group">
          <label class="filter-label">AI模型类别筛选器</label>
          <el-select
            v-model="selectedAiModelType"
            placeholder="请选择AI模型类别"
            clearable
            style="width: 200px"
            :disabled="!selectedChipType"
          >
            <el-option
              v-for="type in aiModelTypes"
              :key="type.value"
              :label="type.label"
              :value="type.value"
            />
          </el-select>
        </div>

        <div class="filter-group">
          <label class="filter-label">产品系列筛选器</label>
          <el-select
            v-model="selectedProductSeries"
            placeholder="请选择产品系列"
            clearable
            style="width: 200px"
            :disabled="!selectedChipType"
          >
            <el-option
              v-for="series in productSeriesList"
              :key="series.value"
              :label="series.label"
              :value="series.value"
            />
          </el-select>
        </div>

        <div class="filter-group">
          <label class="filter-label">服务器型号筛选器</label>
          <el-select
            v-model="selectedServerModel"
            placeholder="请选择服务器型号"
            clearable
            style="width: 200px"
            :disabled="!selectedChipType"
          >
            <el-option
              v-for="model in serverModelsList"
              :key="model.value"
              :label="model.label"
              :value="model.value"
            />
          </el-select>
        </div>

        <div class="filter-group">
          <label class="filter-label">平台架构筛选器</label>
          <el-select
            v-model="selectedPlatform"
            placeholder="请选择平台架构"
            clearable
            style="width: 200px"
            :disabled="!selectedChipType"
          >
            <el-option
              v-for="platform in platforms"
              :key="platform.value"
              :label="platform.label"
              :value="platform.value"
            />
          </el-select>
        </div>

        <div class="filter-group">
          <label class="filter-label">操作系统筛选器</label>
          <el-select
            v-model="selectedOs"
            placeholder="请选择操作系统"
            clearable
            style="width: 200px"
            :disabled="!selectedChipType"
          >
            <el-option
              v-for="os in osList"
              :key="os.value"
              :label="os.label"
              :value="os.value"
            />
          </el-select>
        </div>
      </div>
    </el-card>

    <!-- 版本类型和搜索 -->
    <el-card class="version-card">
      <div class="version-container">
        <div class="version-tabs">
          <el-radio-group v-model="selectedVersionType">
            <el-radio-button label="stable">正式版本</el-radio-button>
            <el-radio-button label="beta">测试版本</el-radio-button>
          </el-radio-group>
          <div v-if="selectedChipType" class="chip-hint">
            请先选择相关芯片
          </div>
        </div>

        <div class="search-actions">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索资源名称或描述"
            :prefix-icon="Search"
            clearable
            style="width: 300px"
          />
          <el-button :icon="Refresh" @click="loadResources">刷新</el-button>
          <el-button @click="clearFilters">清除筛选</el-button>
        </div>
      </div>
    </el-card>

    <!-- 资源列表 -->
    <el-card class="resource-list-card">
      <div v-if="showEmptyHint" class="empty-hint">
        <el-empty description="请先选择相关芯片">
          <el-icon :size="80" color="#909399">
            <Download />
          </el-icon>
        </el-empty>
      </div>

      <el-table
        v-else
        :data="filteredResources"
        :loading="loading"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="资源名" min-width="200" />
        <el-table-column prop="description" label="资源包描述" min-width="250" />
        <el-table-column prop="publishTime" label="发布时间" width="120" />
        <el-table-column label="文件大小" width="120">
          <template #default="{ row }">
            {{ formatSize(row.fileSize) }}
          </template>
        </el-table-column>
        <el-table-column prop="aiModelType" label="AI模型类别" width="130" />
        <el-table-column prop="boardType" label="板卡类型" width="120" />
        <el-table-column prop="platform" label="平台架构" width="120" />
        <el-table-column prop="os" label="操作系统" width="140" />
        <el-table-column prop="md5" label="MD5" width="280" show-overflow-tooltip />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              :icon="Download"
              @click="handleDownload(row)"
            >
              下载
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty
        v-if="!showEmptyHint && filteredResources.length === 0"
        description="暂无匹配的资源"
      />
    </el-card>
  </div>
</template>

<style scoped>
.download-center {
  padding: 24px;
}

.chip-selector-card {
  margin-bottom: 24px;
}

.chip-types {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.chip-type-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 120px;
  height: 120px;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  background: white;
}

.chip-type-item:hover {
  border-color: #409eff;
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
}

.chip-type-item.active {
  border-color: #409eff;
  background: #ecf5ff;
}

.chip-icon {
  font-size: 48px;
  margin-bottom: 8px;
}

.chip-label {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.filter-card {
  margin-bottom: 24px;
}

.filter-header {
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.filter-title {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.required {
  color: #f56c6c;
  font-size: 13px;
}

.filter-container {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-label {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
}

.version-card {
  margin-bottom: 24px;
}

.version-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.version-tabs {
  display: flex;
  align-items: center;
  gap: 16px;
}

.chip-hint {
  font-size: 13px;
  color: #909399;
}

.search-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.resource-list-card {
  margin-bottom: 24px;
}

.empty-hint {
  padding: 60px 0;
  text-align: center;
}

@media (max-width: 768px) {
  .chip-types {
    justify-content: center;
  }

  .filter-container {
    flex-direction: column;
  }

  .version-container {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }

  .search-actions {
    width: 100%;
    flex-direction: column;
  }

  .search-actions > * {
    width: 100%;
  }
}
</style>
