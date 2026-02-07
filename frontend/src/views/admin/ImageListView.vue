<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { getImageList, syncImages } from '@/api/admin'
import request from '@/utils/request'
import PageHeader from '@/components/common/PageHeader.vue'

interface Image {
  id: number
  name: string
  display_name: string
  description?: string
  category: string
  framework: string
  framework_version?: string
  cuda_version: string
  python_version?: string
  registry_url?: string
  is_official: boolean
  size?: number
  status: string
  created_at: string
}

const images = ref<Image[]>([])
const loading = ref(false)
const syncing = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const filters = ref({
  category: '',
  framework: '',
  status: '',
})

const loadImages = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: page.value,
      pageSize: pageSize.value,
    }
    if (filters.value.category) params.category = filters.value.category
    if (filters.value.framework) params.framework = filters.value.framework
    if (filters.value.status) params.status = filters.value.status

    const res = await getImageList(params)
    images.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('加载镜像列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  page.value = 1
  loadImages()
}

const handleReset = () => {
  filters.value = { category: '', framework: '', status: '' }
  page.value = 1
  loadImages()
}

const handleSync = async () => {
  syncing.value = true
  try {
    await syncImages()
    ElMessage.success('同步任务已触发')
    setTimeout(loadImages, 1000)
  } catch {
    // 拦截器已处理
  } finally {
    syncing.value = false
  }
}

const handlePageChange = (val: number) => {
  page.value = val
  loadImages()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  page.value = 1
  loadImages()
}

// 镜像大小格式化
const formatSize = (bytes?: number) => {
  if (!bytes || bytes === 0) return '-'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(i > 0 ? 1 : 0) + ' ' + units[i]
}

const formatDate = (value?: string) => {
  if (!value) return '-'
  return new Date(value).toLocaleString('zh-CN')
}

const statusLabel = (status: string) => {
  return status === 'active' ? '可用' : status === 'deprecated' ? '已弃用' : status
}

// 详情抽屉
const detailVisible = ref(false)
const detailImage = ref<Image | null>(null)

const handleViewDetail = (row: Image) => {
  detailImage.value = row
  detailVisible.value = true
}

// 删除镜像
const handleDelete = async (row: Image) => {
  try {
    await ElMessageBox.confirm(
      `确定删除镜像「${row.display_name || row.name}」吗？`,
      '确认删除',
      { type: 'warning' }
    )
    await request.delete(`/admin/images/${row.id}`)
    ElMessage.success('删除成功')
    loadImages()
  } catch {
    // 取消或失败
  }
}

onMounted(() => {
  loadImages()
})
</script>

<template>
  <div class="image-list">
    <PageHeader title="镜像管理" subtitle="管理平台可用的容器镜像">
      <template #actions>
        <el-button :icon="Refresh" @click="loadImages">刷新</el-button>
        <el-button type="primary" :icon="Refresh" :loading="syncing" @click="handleSync">
          同步镜像
        </el-button>
      </template>
    </PageHeader>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters">
        <el-form-item label="分类">
          <el-select v-model="filters.category" placeholder="全部分类" clearable style="width: 140px">
            <el-option label="训练" value="training" />
            <el-option label="推理" value="inference" />
            <el-option label="开发" value="development" />
            <el-option label="基础" value="base" />
          </el-select>
        </el-form-item>
        <el-form-item label="框架">
          <el-select v-model="filters.framework" placeholder="全部框架" clearable style="width: 140px">
            <el-option label="PyTorch" value="pytorch" />
            <el-option label="TensorFlow" value="tensorflow" />
            <el-option label="JAX" value="jax" />
            <el-option label="PaddlePaddle" value="paddlepaddle" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="可用" value="active" />
            <el-option label="已弃用" value="deprecated" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <el-table :data="images" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="镜像名称" min-width="200">
          <template #default="{ row }">
            <el-link type="primary" @click="handleViewDetail(row)">
              {{ row.display_name || row.name }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ row.category || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="framework" label="框架" width="130" />
        <el-table-column prop="cuda_version" label="CUDA" width="100" />
        <el-table-column prop="size" label="大小" width="100">
          <template #default="{ row }">
            {{ formatSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleViewDetail(row)">详情</el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </el-card>

    <!-- 详情抽屉 -->
    <el-drawer v-model="detailVisible" title="镜像详情" size="420px">
      <el-descriptions v-if="detailImage" :column="1" border>
        <el-descriptions-item label="ID">{{ detailImage.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ detailImage.name }}</el-descriptions-item>
        <el-descriptions-item label="显示名称">{{ detailImage.display_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ detailImage.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="分类">{{ detailImage.category || '-' }}</el-descriptions-item>
        <el-descriptions-item label="框架">{{ detailImage.framework || '-' }}</el-descriptions-item>
        <el-descriptions-item label="CUDA 版本">{{ detailImage.cuda_version || '-' }}</el-descriptions-item>
        <el-descriptions-item label="Python 版本">{{ detailImage.python_version || '-' }}</el-descriptions-item>
        <el-descriptions-item label="仓库地址">
          <span class="registry-url">{{ detailImage.registry_url || '-' }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="大小">{{ formatSize(detailImage.size) }}</el-descriptions-item>
        <el-descriptions-item label="官方镜像">{{ detailImage.is_official ? '是' : '否' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="detailImage.status === 'active' ? 'success' : 'info'" size="small">
            {{ statusLabel(detailImage.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(detailImage.created_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<style scoped>
.image-list {
  padding: 24px;
}

.filter-card {
  margin-bottom: 16px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.registry-url {
  word-break: break-all;
  font-size: 12px;
  color: #606266;
}
</style>
