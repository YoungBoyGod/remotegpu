<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, Upload } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import {
  getDatasetList,
  initMultipartUpload,
  completeMultipartUpload,
  mountDataset,
  getMyMachines,
} from '@/api/customer'
import type { Machine } from '@/types/machine'

interface Dataset {
  id: number
  name: string
  description?: string
  total_size: number
  file_count: number
  status: string
  visibility?: string
  created_at: string
  updated_at?: string
}

const datasets = ref<Dataset[]>([])
const total = ref(0)
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)

const loadDatasets = async () => {
  loading.value = true
  try {
    const res = await getDatasetList({ page: page.value, pageSize: pageSize.value })
    datasets.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('加载数据集列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (p: number) => {
  page.value = p
  loadDatasets()
}

const handleSizeChange = (s: number) => {
  pageSize.value = s
  page.value = 1
  loadDatasets()
}

// 文件大小格式化
const formatSize = (bytes: number) => {
  if (!bytes || bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(i > 0 ? 1 : 0) + ' ' + units[i]
}

// 状态标签
const statusTagType = (status: string) => {
  switch (status) {
    case 'ready': return 'success'
    case 'uploading': return 'warning'
    case 'error': return 'danger'
    default: return 'info'
  }
}

const statusLabel = (status: string) => {
  switch (status) {
    case 'ready': return '就绪'
    case 'uploading': return '上传中'
    case 'error': return '错误'
    default: return status
  }
}

const formatDate = (value?: string) => {
  if (!value) return '-'
  return new Date(value).toLocaleString('zh-CN')
}

// 上传数据集
const uploadDialogVisible = ref(false)
const uploadLoading = ref(false)
const uploadFile = ref<File | null>(null)
const uploadProgress = ref(0)
const uploadName = ref('')

const openUploadDialog = () => {
  uploadDialogVisible.value = true
  uploadFile.value = null
  uploadProgress.value = 0
  uploadName.value = ''
}

const handleFileSelect = (file: { raw: File }) => {
  uploadFile.value = file.raw
  if (!uploadName.value) {
    uploadName.value = file.raw.name.replace(/\.[^/.]+$/, '')
  }
  return false
}

const handleUpload = async () => {
  if (!uploadFile.value) {
    ElMessage.warning('请选择文件')
    return
  }
  if (!uploadName.value.trim()) {
    ElMessage.warning('请输入数据集名称')
    return
  }

  uploadLoading.value = true
  uploadProgress.value = 0
  try {
    // 初始化分片上传
    const initRes = await initMultipartUpload({
      filename: uploadFile.value.name,
      size: uploadFile.value.size,
    })
    const { upload_id } = initRes.data

    // 模拟上传进度（实际分片上传需要后端返回预签名 URL）
    uploadProgress.value = 50

    // 完成上传
    await completeMultipartUpload(0, {
      upload_id,
      name: uploadName.value.trim(),
      size: uploadFile.value.size,
    })

    uploadProgress.value = 100
    ElMessage.success('数据集上传成功')
    uploadDialogVisible.value = false
    loadDatasets()
  } catch (error) {
    ElMessage.error('上传失败')
    console.error('上传数据集失败:', error)
  } finally {
    uploadLoading.value = false
  }
}

// 挂载数据集
const mountDialogVisible = ref(false)
const mountLoading = ref(false)
const mountTarget = ref<Dataset | null>(null)
const mountForm = ref({
  machine_id: '',
  mount_path: '/data',
  read_only: true,
})
const machines = ref<Machine[]>([])

const openMountDialog = async (dataset: Dataset) => {
  mountTarget.value = dataset
  mountForm.value = { machine_id: '', mount_path: '/data', read_only: true }
  mountDialogVisible.value = true
  try {
    const res = await getMyMachines({ page: 1, pageSize: 200 })
    machines.value = res.data.list || []
  } catch (error) {
    console.error('加载机器列表失败:', error)
  }
}

const handleMount = async () => {
  if (!mountTarget.value) return
  if (!mountForm.value.machine_id) {
    ElMessage.warning('请选择机器')
    return
  }
  mountLoading.value = true
  try {
    await mountDataset(mountTarget.value.id, {
      machine_id: mountForm.value.machine_id,
      mount_path: mountForm.value.mount_path,
      read_only: mountForm.value.read_only,
    })
    ElMessage.success('挂载成功')
    mountDialogVisible.value = false
  } catch (error) {
    ElMessage.error('挂载失败')
    console.error('挂载数据集失败:', error)
  } finally {
    mountLoading.value = false
  }
}

// 删除数据集（后端暂未实现 DELETE 接口，预留）
const handleDelete = async (dataset: Dataset) => {
  try {
    await ElMessageBox.confirm(
      `确定删除数据集「${dataset.name}」吗？此操作不可恢复。`,
      '确认删除',
      { type: 'warning' }
    )
    ElMessage.info('删除功能暂未开放')
  } catch {
    // 取消
  }
}

onMounted(() => {
  loadDatasets()
})
</script>

<template>
  <div class="dataset-list-view">
    <PageHeader title="数据集管理" subtitle="管理我的训练数据集">
      <template #actions>
        <el-button :icon="Refresh" @click="loadDatasets">刷新</el-button>
        <el-button type="primary" :icon="Upload" @click="openUploadDialog">上传数据集</el-button>
      </template>
    </PageHeader>

    <el-card>
      <el-table :data="datasets" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" min-width="180">
          <template #default="{ row }">
            <span class="dataset-name">{{ row.name }}</span>
            <div v-if="row.description" class="dataset-desc">{{ row.description }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="total_size" label="大小" width="120">
          <template #default="{ row }">
            {{ formatSize(row.total_size) }}
          </template>
        </el-table-column>
        <el-table-column prop="file_count" label="文件数" width="90" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" size="small">
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="上传时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" :disabled="row.status !== 'ready'" @click="openMountDialog(row)">
              挂载
            </el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-if="total > 0"
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

    <!-- 上传对话框 -->
    <el-dialog v-model="uploadDialogVisible" title="上传数据集" width="520px" :close-on-click-modal="false">
      <el-form label-width="100px">
        <el-form-item label="数据集名称">
          <el-input v-model="uploadName" placeholder="请输入数据集名称" />
        </el-form-item>
        <el-form-item label="选择文件">
          <el-upload
            drag
            :auto-upload="false"
            :limit="1"
            :on-change="handleFileSelect"
            :show-file-list="true"
          >
            <el-icon class="el-icon--upload"><Upload /></el-icon>
            <div class="el-upload__text">拖拽文件到此处，或<em>点击上传</em></div>
          </el-upload>
        </el-form-item>
        <el-form-item v-if="uploadProgress > 0" label="上传进度">
          <el-progress :percentage="uploadProgress" :status="uploadProgress === 100 ? 'success' : undefined" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="uploadLoading" @click="handleUpload">上传</el-button>
      </template>
    </el-dialog>

    <!-- 挂载对话框 -->
    <el-dialog v-model="mountDialogVisible" title="挂载数据集" width="480px" :close-on-click-modal="false">
      <el-form label-width="100px">
        <el-form-item label="数据集">
          <span>{{ mountTarget?.name }}</span>
        </el-form-item>
        <el-form-item label="目标机器">
          <el-select v-model="mountForm.machine_id" placeholder="请选择机器" filterable style="width: 100%">
            <el-option
              v-for="m in machines"
              :key="m.id"
              :label="m.name || m.hostname || m.id"
              :value="m.id.toString()"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="挂载路径">
          <el-input v-model="mountForm.mount_path" placeholder="/data" />
        </el-form-item>
        <el-form-item label="只读">
          <el-switch v-model="mountForm.read_only" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="mountDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="mountLoading" @click="handleMount">确认挂载</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.dataset-list-view {
  padding: 24px;
}

.dataset-name {
  font-weight: 500;
  color: #303133;
}

.dataset-desc {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
