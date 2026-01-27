<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Download, View, Edit, Delete } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'

interface Release {
  id: number
  version: string
  title: string
  description: string
  releaseDate: string
  status: 'draft' | 'published' | 'archived'
  downloadUrl?: string
  changelog: string[]
}

const releases = ref<Release[]>([])
const loading = ref(false)

// 加载发布版本列表
const loadReleases = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 500))
    releases.value = [
      {
        id: 1,
        version: 'v2.1.0',
        title: 'RemoteGPU 2.1.0 正式版',
        description: '新增算力市场功能，优化资源管理界面',
        releaseDate: '2026-01-27',
        status: 'published',
        downloadUrl: '/downloads/v2.1.0',
        changelog: [
          '新增算力市场功能',
          '优化资源管理界面',
          '修复环境部署bug',
          '提升系统性能'
        ]
      },
      {
        id: 2,
        version: 'v2.0.5',
        title: 'RemoteGPU 2.0.5 稳定版',
        description: '修复已知问题，提升系统稳定性',
        releaseDate: '2026-01-15',
        status: 'published',
        downloadUrl: '/downloads/v2.0.5',
        changelog: [
          '修复训练任务异常退出问题',
          '优化数据集上传速度',
          '改进镜像构建流程'
        ]
      }
    ]
  } catch (error) {
    ElMessage.error('加载版本列表失败')
  } finally {
    loading.value = false
  }
}

// 下载版本
const handleDownload = (release: Release) => {
  ElMessage.success(`开始下载 ${release.version}`)
  // TODO: 实现下载逻辑
}

// 查看详情
const handleView = (release: Release) => {
  ElMessage.info(`查看版本详情: ${release.version}`)
  // TODO: 打开详情对话框
}

onMounted(() => {
  loadReleases()
})
</script>

<template>
  <div class="release-version">
    <PageHeader title="发布版本">
      <template #actions>
        <el-button type="primary" :icon="Plus">
          发布新版本
        </el-button>
      </template>
    </PageHeader>

    <div v-loading="loading" class="release-list">
      <div
        v-for="release in releases"
        :key="release.id"
        class="release-card"
      >
        <div class="release-header">
          <div class="release-info">
            <h2 class="release-version">{{ release.version }}</h2>
            <h3 class="release-title">{{ release.title }}</h3>
          </div>
          <el-tag
            :type="release.status === 'published' ? 'success' : 'info'"
            size="large"
          >
            {{ release.status === 'published' ? '已发布' : '草稿' }}
          </el-tag>
        </div>

        <div class="release-body">
          <p class="release-description">{{ release.description }}</p>
          <div class="release-date">
            发布时间: {{ release.releaseDate }}
          </div>

          <div class="changelog">
            <h4>更新内容:</h4>
            <ul>
              <li v-for="(item, index) in release.changelog" :key="index">
                {{ item }}
              </li>
            </ul>
          </div>
        </div>

        <div class="release-actions">
          <el-button :icon="View" @click="handleView(release)">
            查看详情
          </el-button>
          <el-button
            type="primary"
            :icon="Download"
            @click="handleDownload(release)"
          >
            下载
          </el-button>
        </div>
      </div>

      <el-empty
        v-if="!loading && releases.length === 0"
        description="暂无发布版本"
      />
    </div>
  </div>
</template>

<style scoped>
.release-version {
  padding: 24px;
}

.release-list {
  min-height: 400px;
}

.release-card {
  background: white;
  border-radius: 8px;
  padding: 24px;
  margin-bottom: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.release-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.release-info {
  flex: 1;
}

.release-version {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
  color: #409EFF;
}

.release-title {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
  color: #303133;
}

.release-body {
  margin-bottom: 16px;
}

.release-description {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #606266;
}

.release-date {
  font-size: 13px;
  color: #909399;
  margin-bottom: 16px;
}

.changelog {
  margin-top: 16px;
}

.changelog h4 {
  margin: 0 0 8px 0;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.changelog ul {
  margin: 0;
  padding-left: 20px;
}

.changelog li {
  margin-bottom: 4px;
  font-size: 14px;
  color: #606266;
}

.release-actions {
  display: flex;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}
</style>
