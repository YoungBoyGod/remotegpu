<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Document, Reading } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'

interface DocCategory {
  id: number
  name: string
  icon: string
  count: number
}

interface DocVersion {
  value: string
  label: string
  releaseDate: string
}

interface DocItem {
  id: number
  title: string
  category: string
  version: string
  description: string
  updateDate: string
  views: number
}

const categories = ref<DocCategory[]>([
  { id: 1, name: '快速入门', icon: 'Reading', count: 8 },
  { id: 2, name: '环境部署', icon: 'Setting', count: 12 },
  { id: 3, name: '训练任务', icon: 'TrendCharts', count: 15 },
  { id: 4, name: 'API文档', icon: 'Document', count: 25 },
  { id: 5, name: '常见问题', icon: 'QuestionFilled', count: 20 }
])

// 版本列表
const versions = ref<DocVersion[]>([
  { value: 'v2.0', label: 'v2.0 (最新)', releaseDate: '2026-01-27' },
  { value: 'v1.9', label: 'v1.9', releaseDate: '2026-01-15' },
  { value: 'v1.8', label: 'v1.8', releaseDate: '2025-12-20' },
  { value: 'v1.7', label: 'v1.7', releaseDate: '2025-11-10' }
])

const documents = ref<DocItem[]>([])
const loading = ref(false)
const searchKeyword = ref('')
const selectedCategory = ref('')
const selectedVersion = ref('v2.0') // 默认选择最新版本

// 过滤后的文档列表
const filteredDocuments = computed(() => {
  let result = documents.value

  // 版本筛选
  if (selectedVersion.value) {
    result = result.filter(d => d.version === selectedVersion.value)
  }

  // 分类筛选
  if (selectedCategory.value) {
    result = result.filter(d => d.category === selectedCategory.value)
  }

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(d =>
      d.title.toLowerCase().includes(keyword) ||
      d.description.toLowerCase().includes(keyword)
    )
  }

  return result
})

// 加载文档列表
const loadDocuments = async () => {
  loading.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 500))
    documents.value = [
      {
        id: 1,
        title: '快速开始使用RemoteGPU',
        category: '快速入门',
        version: 'v2.0',
        description: '了解如何快速开始使用RemoteGPU平台进行GPU计算',
        updateDate: '2026-01-27',
        views: 1250
      },
      {
        id: 2,
        title: '创建和管理开发环境',
        category: '环境部署',
        version: 'v2.0',
        description: '学习如何创建、配置和管理您的开发环境',
        updateDate: '2026-01-26',
        views: 890
      },
      {
        id: 3,
        title: '提交训练任务指南',
        category: '训练任务',
        version: 'v2.0',
        description: '详细介绍如何提交和监控训练任务',
        updateDate: '2026-01-25',
        views: 1100
      },
      {
        id: 4,
        title: 'API接口文档',
        category: 'API文档',
        version: 'v2.0',
        description: 'RemoteGPU平台完整的API接口文档',
        updateDate: '2026-01-24',
        views: 2300
      },
      {
        id: 5,
        title: '快速开始使用RemoteGPU',
        category: '快速入门',
        version: 'v1.9',
        description: '了解如何快速开始使用RemoteGPU平台进行GPU计算（v1.9版本）',
        updateDate: '2026-01-15',
        views: 850
      },
      {
        id: 6,
        title: '环境部署指南',
        category: '环境部署',
        version: 'v1.9',
        description: '学习如何在v1.9版本中部署开发环境',
        updateDate: '2026-01-14',
        views: 720
      }
    ]
  } catch (error) {
    ElMessage.error('加载文档列表失败')
  } finally {
    loading.value = false
  }
}

// 查看文档
const handleViewDoc = (doc: DocItem) => {
  ElMessage.info(`查看文档: ${doc.title}`)
  // TODO: 跳转到文档详情页
}

// 选择分类
const handleSelectCategory = (categoryName: string) => {
  selectedCategory.value = selectedCategory.value === categoryName ? '' : categoryName
}

onMounted(() => {
  loadDocuments()
})
</script>

<template>
  <div class="document-center">
    <PageHeader title="文档中心" />

    <div class="doc-container">
      <!-- 版本选择 -->
      <div class="version-section">
        <div class="version-label">文档版本：</div>
        <el-select
          v-model="selectedVersion"
          placeholder="请选择版本"
          style="width: 200px"
        >
          <el-option
            v-for="version in versions"
            :key="version.value"
            :label="version.label"
            :value="version.value"
          >
            <span>{{ version.label }}</span>
            <span style="float: right; color: #909399; font-size: 13px">
              {{ version.releaseDate }}
            </span>
          </el-option>
        </el-select>
      </div>

      <!-- 分类卡片 -->
      <div class="category-section">
        <h3>文档分类</h3>
        <div class="category-grid">
          <div
            v-for="category in categories"
            :key="category.id"
            class="category-card"
            :class="{ active: selectedCategory === category.name }"
            @click="handleSelectCategory(category.name)"
          >
            <el-icon class="category-icon" :size="32">
              <component :is="category.icon" />
            </el-icon>
            <div class="category-name">{{ category.name }}</div>
            <div class="category-count">{{ category.count }} 篇</div>
          </div>
        </div>
      </div>

      <!-- 搜索栏 -->
      <div class="search-section">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索文档标题或内容"
          clearable
          size="large"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>

      <!-- 文档列表 -->
      <div v-loading="loading" class="document-list">
        <div
          v-for="doc in filteredDocuments"
          :key="doc.id"
          class="doc-item"
          @click="handleViewDoc(doc)"
        >
          <div class="doc-icon">
            <el-icon :size="24"><Document /></el-icon>
          </div>
          <div class="doc-content">
            <h4 class="doc-title">{{ doc.title }}</h4>
            <p class="doc-description">{{ doc.description }}</p>
            <div class="doc-meta">
              <span class="doc-category">
                <el-tag size="small">{{ doc.category }}</el-tag>
              </span>
              <span class="doc-date">更新于 {{ doc.updateDate }}</span>
              <span class="doc-views">{{ doc.views }} 次浏览</span>
            </div>
          </div>
          <div class="doc-action">
            <el-icon :size="20"><Reading /></el-icon>
          </div>
        </div>

        <el-empty
          v-if="!loading && filteredDocuments.length === 0"
          description="暂无相关文档"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.document-center {
  padding: 24px;
}

.doc-container {
  background: white;
  border-radius: 8px;
  padding: 24px;
}

.version-section {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.version-label {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.category-section {
  margin-bottom: 32px;
}

.category-section h3 {
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 16px;
}

.category-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 24px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
}

.category-card:hover {
  background: #e8f4ff;
  transform: translateY(-4px);
}

.category-card.active {
  background: #e8f4ff;
  border-color: #409EFF;
}

.category-icon {
  color: #409EFF;
  margin-bottom: 12px;
}

.category-name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 8px;
}

.category-count {
  font-size: 13px;
  color: #909399;
}

.search-section {
  margin-bottom: 24px;
}

.document-list {
  min-height: 400px;
}

.doc-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 12px;
  background: #f8f9fa;
  cursor: pointer;
  transition: all 0.3s;
}

.doc-item:hover {
  background: #e8f4ff;
  transform: translateX(4px);
}

.doc-icon {
  flex-shrink: 0;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border-radius: 8px;
  color: #409EFF;
}

.doc-content {
  flex: 1;
}

.doc-title {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.doc-description {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}

.doc-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 13px;
  color: #909399;
}

.doc-action {
  flex-shrink: 0;
  color: #409EFF;
}

@media (max-width: 1200px) {
  .category-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .category-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
