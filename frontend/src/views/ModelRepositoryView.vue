<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Download,
  Upload,
  View,
  Edit,
  Delete,
  Search,
  Refresh,
  FolderOpened
} from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'

interface Model {
  id: number
  name: string
  version: string
  framework: string
  size: number
  description: string
  tags: string[]
  author: string
  downloads: number
  createdAt: string
  updatedAt: string
  status: 'active' | 'archived'
}

const loading = ref(false)
const models = ref<Model[]>([])
const searchKeyword = ref('')
const selectedFramework = ref<string>('all')
const selectedStatus = ref<string>('all')

// 对话框状态
const modelDialogVisible = ref(false)
const uploadDialogVisible = ref(false)
const isEditMode = ref(false)
const currentModel = ref<Model | null>(null)

// 模型表单
const modelForm = ref({
  name: '',
  version: '',
  framework: '',
  description: '',
  tags: [] as string[],
  status: 'active' as 'active' | 'archived'
})

// 过滤后的模型列表
const filteredModels = computed(() => {
  let result = models.value

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(m =>
      m.name.toLowerCase().includes(keyword) ||
      m.description.toLowerCase().includes(keyword) ||
      m.author.toLowerCase().includes(keyword) ||
      m.tags.some(tag => tag.toLowerCase().includes(keyword))
    )
  }

  // 框架筛选
  if (selectedFramework.value !== 'all') {
    result = result.filter(m => m.framework === selectedFramework.value)
  }

  // 状态筛选
  if (selectedStatus.value !== 'all') {
    result = result.filter(m => m.status === selectedStatus.value)
  }

  return result
})

// 模型统计
const modelStats = computed(() => {
  return {
    total: models.value.length,
    active: models.value.filter(m => m.status === 'active').length,
    archived: models.value.filter(m => m.status === 'archived').length,
    totalDownloads: models.value.reduce((sum, m) => sum + m.downloads, 0)
  }
})

// 加载模型列表
const loadModels = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 500))
    models.value = [
      {
        id: 1,
        name: 'ResNet-50',
        version: 'v1.0.0',
        framework: 'PyTorch',
        size: 102.5,
        description: '深度残差网络，用于图像分类任务',
        tags: ['图像分类', 'CNN', '预训练'],
        author: '张三',
        downloads: 1250,
        createdAt: '2026-01-15',
        updatedAt: '2026-01-20',
        status: 'active'
      },
      {
        id: 2,
        name: 'BERT-Base',
        version: 'v2.1.0',
        framework: 'TensorFlow',
        size: 438.2,
        description: 'BERT基础模型，用于自然语言处理',
        tags: ['NLP', 'Transformer', '预训练'],
        author: '李四',
        downloads: 3420,
        createdAt: '2026-01-10',
        updatedAt: '2026-01-25',
        status: 'active'
      },
      {
        id: 3,
        name: 'YOLOv5',
        version: 'v1.5.0',
        framework: 'PyTorch',
        size: 89.3,
        description: '目标检测模型，实时性能优秀',
        tags: ['目标检测', 'YOLO', '实时'],
        author: '王五',
        downloads: 2180,
        createdAt: '2026-01-05',
        updatedAt: '2026-01-22',
        status: 'active'
      }
    ]
  } catch (error) {
    ElMessage.error('加载模型列表失败')
  } finally {
    loading.value = false
  }
}

// 格式化文件大小
const formatSize = (size: number) => {
  if (size < 1024) return `${size.toFixed(2)} MB`
  return `${(size / 1024).toFixed(2)} GB`
}

// 获取状态类型
const getStatusType = (status: string) => {
  return status === 'active' ? 'success' : 'info'
}

// 获取状态文本
const getStatusText = (status: string) => {
  return status === 'active' ? '活跃' : '已归档'
}

// 打开新增模型对话框
const handleAddModel = () => {
  isEditMode.value = false
  modelForm.value = {
    name: '',
    version: '',
    framework: '',
    description: '',
    tags: [],
    status: 'active'
  }
  modelDialogVisible.value = true
}

// 打开编辑模型对话框
const handleEditModel = (model: Model) => {
  isEditMode.value = true
  currentModel.value = model
  modelForm.value = {
    name: model.name,
    version: model.version,
    framework: model.framework,
    description: model.description,
    tags: [...model.tags],
    status: model.status
  }
  modelDialogVisible.value = true
}

// 保存模型
const handleSaveModel = async () => {
  try {
    // TODO: 调用API保存数据
    await new Promise(resolve => setTimeout(resolve, 300))

    if (isEditMode.value && currentModel.value) {
      // 更新模型
      Object.assign(currentModel.value, {
        ...modelForm.value,
        updatedAt: new Date().toLocaleDateString('zh-CN')
      })
      ElMessage.success('模型信息已更新')
    } else {
      // 新增模型
      const newModel: Model = {
        id: models.value.length + 1,
        ...modelForm.value,
        size: 0,
        author: '当前用户',
        downloads: 0,
        createdAt: new Date().toLocaleDateString('zh-CN'),
        updatedAt: new Date().toLocaleDateString('zh-CN')
      }
      models.value.push(newModel)
      ElMessage.success('模型已添加')
    }

    modelDialogVisible.value = false
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

// 删除模型
const handleDeleteModel = async (model: Model) => {
  try {
    await ElMessageBox.confirm('确认删除此模型？', '删除模型', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'error'
    })

    // TODO: 调用API删除
    const index = models.value.findIndex(m => m.id === model.id)
    if (index > -1) {
      models.value.splice(index, 1)
    }
    ElMessage.success('模型已删除')
  } catch {
    // 用户取消
  }
}

// 下载模型
const handleDownloadModel = (model: Model) => {
  ElMessage.success(`开始下载模型: ${model.name}`)
  // TODO: 实现下载逻辑
  model.downloads++
}

// 打开上传对话框
const handleUploadModel = () => {
  uploadDialogVisible.value = true
}

onMounted(() => {
  loadModels()
})
</script>

<template>
  <div class="model-repository">
    <PageHeader title="模型仓库" />

    <!-- 统计卡片 -->
    <div class="stats-container">
      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #409EFF">
            <el-icon :size="32"><FolderOpened /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">模型总数</div>
            <div class="stat-value">{{ modelStats.total }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #67C23A">
            <el-icon :size="32"><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">活跃模型</div>
            <div class="stat-value">{{ modelStats.active }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #E6A23C">
            <el-icon :size="32"><Download /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">总下载量</div>
            <div class="stat-value">{{ modelStats.totalDownloads }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon" style="background: #909399">
            <el-icon :size="32"><Upload /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">已归档</div>
            <div class="stat-value">{{ modelStats.archived }}</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 筛选和操作区域 -->
    <el-card class="filter-card">
      <div class="filter-container">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索模型名称、描述或标签"
          :prefix-icon="Search"
          clearable
          style="width: 300px"
        />

        <el-select v-model="selectedFramework" placeholder="框架" clearable style="width: 150px">
          <el-option label="全部框架" value="all" />
          <el-option label="PyTorch" value="PyTorch" />
          <el-option label="TensorFlow" value="TensorFlow" />
          <el-option label="ONNX" value="ONNX" />
          <el-option label="Keras" value="Keras" />
        </el-select>

        <el-select v-model="selectedStatus" placeholder="状态" clearable style="width: 120px">
          <el-option label="全部状态" value="all" />
          <el-option label="活跃" value="active" />
          <el-option label="已归档" value="archived" />
        </el-select>

        <el-button :icon="Refresh" @click="loadModels">刷新</el-button>
        <el-button type="primary" :icon="Upload" @click="handleUploadModel">
          上传模型
        </el-button>
        <el-button type="success" :icon="Plus" @click="handleAddModel">
          新增模型
        </el-button>
      </div>
    </el-card>

    <!-- 模型列表 -->
    <div v-loading="loading" class="model-grid">
      <el-card
        v-for="model in filteredModels"
        :key="model.id"
        class="model-card"
        shadow="hover"
      >
        <template #header>
          <div class="model-header">
            <div class="model-title">
              <h3>{{ model.name }}</h3>
              <el-tag size="small">{{ model.version }}</el-tag>
            </div>
            <el-tag :type="getStatusType(model.status)" size="small">
              {{ getStatusText(model.status) }}
            </el-tag>
          </div>
        </template>

        <div class="model-content">
          <p class="model-description">{{ model.description }}</p>

          <div class="model-info">
            <div class="info-item">
              <span class="info-label">框架:</span>
              <span class="info-value">{{ model.framework }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">大小:</span>
              <span class="info-value">{{ formatSize(model.size) }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">作者:</span>
              <span class="info-value">{{ model.author }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">下载量:</span>
              <span class="info-value">{{ model.downloads }}</span>
            </div>
          </div>

          <div class="model-tags">
            <el-tag
              v-for="tag in model.tags"
              :key="tag"
              size="small"
              type="info"
              style="margin-right: 8px"
            >
              {{ tag }}
            </el-tag>
          </div>

          <div class="model-meta">
            <span>创建: {{ model.createdAt }}</span>
            <span>更新: {{ model.updatedAt }}</span>
          </div>
        </div>

        <template #footer>
          <div class="model-actions">
            <el-button size="small" :icon="Download" @click="handleDownloadModel(model)">
              下载
            </el-button>
            <el-button size="small" :icon="Edit" @click="handleEditModel(model)">
              编辑
            </el-button>
            <el-button
              size="small"
              type="danger"
              :icon="Delete"
              @click="handleDeleteModel(model)"
            >
              删除
            </el-button>
          </div>
        </template>
      </el-card>

      <el-empty
        v-if="!loading && filteredModels.length === 0"
        description="暂无模型"
      />
    </div>

    <!-- 模型编辑对话框 -->
    <el-dialog
      v-model="modelDialogVisible"
      :title="isEditMode ? '编辑模型' : '新增模型'"
      width="600px"
    >
      <el-form :model="modelForm" label-width="100px">
        <el-form-item label="模型名称">
          <el-input v-model="modelForm.name" placeholder="请输入模型名称" />
        </el-form-item>
        <el-form-item label="版本号">
          <el-input v-model="modelForm.version" placeholder="例如: v1.0.0" />
        </el-form-item>
        <el-form-item label="框架">
          <el-select v-model="modelForm.framework" placeholder="请选择框架">
            <el-option label="PyTorch" value="PyTorch" />
            <el-option label="TensorFlow" value="TensorFlow" />
            <el-option label="ONNX" value="ONNX" />
            <el-option label="Keras" value="Keras" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="modelForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入模型描述"
          />
        </el-form-item>
        <el-form-item label="标签">
          <el-select
            v-model="modelForm.tags"
            multiple
            filterable
            allow-create
            placeholder="请选择或输入标签"
          >
            <el-option label="图像分类" value="图像分类" />
            <el-option label="目标检测" value="目标检测" />
            <el-option label="NLP" value="NLP" />
            <el-option label="预训练" value="预训练" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="modelForm.status" placeholder="请选择状态">
            <el-option label="活跃" value="active" />
            <el-option label="已归档" value="archived" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="modelDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveModel">保存</el-button>
      </template>
    </el-dialog>

    <!-- 上传对话框 -->
    <el-dialog v-model="uploadDialogVisible" title="上传模型" width="600px">
      <el-upload
        drag
        action="#"
        :auto-upload="false"
      >
        <el-icon class="el-icon--upload"><Upload /></el-icon>
        <div class="el-upload__text">
          将模型文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持 .pth, .h5, .onnx, .pb 等格式，单个文件不超过 5GB
          </div>
        </template>
      </el-upload>
      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary">开始上传</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.model-repository {
  padding: 24px;
}

.stats-container {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  cursor: pointer;
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.filter-card {
  margin-bottom: 24px;
}

.filter-container {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.model-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  min-height: 400px;
}

.model-card {
  transition: all 0.3s;
}

.model-card:hover {
  transform: translateY(-4px);
}

.model-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.model-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.model-title h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.model-content {
  min-height: 200px;
}

.model-description {
  margin: 0 0 16px 0;
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}

.model-info {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
  margin-bottom: 16px;
}

.info-item {
  font-size: 13px;
}

.info-label {
  color: #909399;
  margin-right: 4px;
}

.info-value {
  color: #303133;
  font-weight: 500;
}

.model-tags {
  margin-bottom: 16px;
}

.model-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #909399;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.model-actions {
  display: flex;
  gap: 8px;
}

@media (max-width: 1200px) {
  .stats-container {
    grid-template-columns: repeat(2, 1fr);
  }

  .model-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-container {
    grid-template-columns: 1fr;
  }

  .model-grid {
    grid-template-columns: 1fr;
  }

  .filter-container {
    flex-direction: column;
  }

  .filter-container > * {
    width: 100% !important;
  }
}
</style>
