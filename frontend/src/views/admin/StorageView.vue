<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Download, Delete, Search, FolderOpened, Refresh } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/common/PageHeader.vue'
import {
  getStorageBackends,
  getStorageStats,
  getStorageFiles,
  deleteStorageFile,
  getStorageDownloadUrl,
} from '@/api/admin'
import type { StorageBackend, StorageStats, StorageFileInfo } from '@/api/admin'

const loading = ref(false)
const files = ref<StorageFileInfo[]>([])
const fileTotal = ref(0)
const backends = ref<StorageBackend[]>([])
const selectedBackend = ref('')
const prefix = ref('')
const stats = ref<StorageStats | null>(null)

const formatSize = (bytes: number) => {
  if (!bytes || bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`
}

const formatDateTime = (value?: string) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN')
}

// 加载存储后端列表
const loadBackends = async () => {
  try {
    const res = await getStorageBackends()
    if (res.code === 0) {
      backends.value = res.data?.backends || []
      if (backends.value.length > 0 && !selectedBackend.value) {
        const def = backends.value.find(b => b.is_default)
        selectedBackend.value = def ? def.name : (backends.value[0]?.name ?? '')
      }
    }
  } catch {
    // 忽略
  }
}

// 加载存储统计
const loadStats = async () => {
  try {
    const res = await getStorageStats(selectedBackend.value || undefined)
    if (res.code === 0) {
      stats.value = res.data
    }
  } catch {
    // 忽略
  }
}

// 加载文件列表
const loadFiles = async () => {
  loading.value = true
  try {
    const res = await getStorageFiles(
      selectedBackend.value || undefined,
      prefix.value || undefined,
    )
    if (res.code === 0) {
      files.value = res.data?.files || []
      fileTotal.value = res.data?.total || 0
    }
  } catch {
    ElMessage.error('加载文件列表失败')
  } finally {
    loading.value = false
  }
}

// 切换存储后端
const handleBackendChange = () => {
  prefix.value = ''
  loadStats()
  loadFiles()
}

// 搜索（按前缀过滤）
const handleSearch = () => {
  loadFiles()
}

// 进入目录
const handleEnterDir = (dir: StorageFileInfo) => {
  prefix.value = prefix.value ? `${prefix.value}${dir.name}` : dir.name
  loadFiles()
}

// 返回上级目录
const handleGoUp = () => {
  const parts = prefix.value.replace(/\/$/, '').split('/')
  parts.pop()
  prefix.value = parts.length > 0 ? parts.join('/') + '/' : ''
  loadFiles()
}

// 下载文件
const handleDownload = async (file: StorageFileInfo) => {
  try {
    const filePath = prefix.value + file.name
    const res = await getStorageDownloadUrl(selectedBackend.value, filePath)
    if (res.code === 0 && res.data?.url) {
      window.open(res.data.url, '_blank')
    } else {
      ElMessage.error('获取下载链接失败')
    }
  } catch {
    ElMessage.error('获取下载链接失败')
  }
}

// 删除文件
const handleDelete = async (file: StorageFileInfo) => {
  await ElMessageBox.confirm(`确定删除文件「${file.name}」吗？`, '删除确认', { type: 'warning' })
  try {
    const filePath = prefix.value + file.name
    const res = await deleteStorageFile({
      backend: selectedBackend.value || undefined,
      path: filePath,
    })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      loadFiles()
      loadStats()
    } else {
      ElMessage.error(res.msg || '删除失败')
    }
  } catch {
    // 用户取消
  }
}

// 刷新
const handleRefresh = () => {
  loadStats()
  loadFiles()
}

onMounted(async () => {
  await loadBackends()
  loadStats()
  loadFiles()
})
</script>

<template>
  <div class="storage-view">
    <PageHeader title="存储管理" subtitle="管理存储后端与文件">
      <template #actions>
        <el-button :icon="Refresh" @click="handleRefresh">刷新</el-button>
      </template>
    </PageHeader>

    <!-- 存储概览 -->
    <div class="storage-overview">
      <el-card class="overview-card" shadow="hover">
        <div class="overview-content">
          <el-icon class="overview-icon" :size="32"><FolderOpened /></el-icon>
          <div class="overview-info">
            <div class="overview-label">已用空间</div>
            <div class="overview-value">{{ stats ? formatSize(stats.total_size) : '-' }}</div>
          </div>
        </div>
      </el-card>
      <el-card class="overview-card" shadow="hover">
        <div class="overview-content">
          <div class="overview-info">
            <div class="overview-label">文件总数</div>
            <div class="overview-number">{{ stats?.file_count ?? '-' }}</div>
          </div>
        </div>
      </el-card>
      <el-card class="overview-card" shadow="hover">
        <div class="overview-content">
          <div class="overview-info">
            <div class="overview-label">存储后端</div>
            <div class="overview-number">{{ backends.length }}</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 筛选 -->
    <el-card class="filter-card">
      <div class="filter-container">
        <el-select
          v-model="selectedBackend"
          placeholder="选择存储后端"
          style="width: 200px"
          @change="handleBackendChange"
        >
          <el-option
            v-for="b in backends"
            :key="b.name"
            :label="b.name + (b.is_default ? ' (默认)' : '')"
            :value="b.name"
          />
        </el-select>
        <el-input
          v-model="prefix"
          placeholder="按路径前缀过滤"
          :prefix-icon="Search"
          clearable
          style="width: 280px"
          @clear="handleSearch"
          @keyup.enter="handleSearch"
        />
        <el-button type="primary" @click="handleSearch">搜索</el-button>
      </div>
    </el-card>

    <!-- 路径导航 -->
    <div class="breadcrumb-bar" v-if="prefix">
      <el-button link type="primary" @click="handleGoUp">返回上级</el-button>
      <span class="breadcrumb-path">当前路径：/{{ prefix }}</span>
    </div>

    <!-- 文件列表 -->
    <el-table :data="files" v-loading="loading" border style="width: 100%">
      <el-table-column prop="name" label="文件名" min-width="280" show-overflow-tooltip>
        <template #default="{ row }">
          <el-button v-if="row.is_dir" link type="primary" @click="handleEnterDir(row)">
            {{ row.name }}
          </el-button>
          <span v-else>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="类型" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.is_dir" size="small" type="warning">目录</el-tag>
          <span v-else>{{ row.content_type || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="大小" width="120">
        <template #default="{ row }">
          {{ row.is_dir ? '-' : formatSize(row.size) }}
        </template>
      </el-table-column>
      <el-table-column label="修改时间" width="180">
        <template #default="{ row }">
          {{ formatDateTime(row.last_modified) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <template v-if="!row.is_dir">
            <el-button link type="primary" :icon="Download" @click="handleDownload(row)">下载</el-button>
            <el-button link type="danger" :icon="Delete" @click="handleDelete(row)">删除</el-button>
          </template>
        </template>
      </el-table-column>
      <template #empty>
        <el-empty description="暂无文件" />
      </template>
    </el-table>

  </div>
</template>

<style scoped>
.storage-view {
  padding: 24px;
  background: #f5f7fa;
  min-height: 100%;
}

.storage-overview {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.overview-card {
  flex: 1;
}

.overview-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.overview-icon {
  color: #409eff;
}

.overview-label {
  font-size: 13px;
  color: #909399;
}

.overview-value {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.overview-number {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-container {
  display: flex;
  gap: 10px;
}

.breadcrumb-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.breadcrumb-path {
  font-size: 13px;
  color: #606266;
}
</style>
