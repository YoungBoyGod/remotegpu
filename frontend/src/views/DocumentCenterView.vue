<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Upload, Edit, Delete, Download } from '@element-plus/icons-vue'
import type { FormInstance, UploadFile } from 'element-plus'
import PageHeader from '@/components/common/PageHeader.vue'
import {
  getDocumentList,
  uploadDocument,
  updateDocument,
  deleteDocument,
  getDocumentCategories,
  getDocumentDownloadUrl,
} from '@/api/admin'
import type { DocumentItem } from '@/api/admin'

// 状态
const documents = ref<DocumentItem[]>([])
const categories = ref<string[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const searchKeyword = ref('')
const selectedCategory = ref('')

// 上传对话框
const uploadDialogVisible = ref(false)
const uploading = ref(false)
const uploadFormRef = ref<FormInstance>()
const uploadForm = reactive({
  title: '',
  category: '',
  file: null as File | null,
})
const uploadRules = {
  title: [{ required: true, message: '请输入文档标题', trigger: 'blur' }],
  category: [{ required: true, message: '请选择文档分类', trigger: 'change' }],
}

// 编辑对话框
const editDialogVisible = ref(false)
const editSaving = ref(false)
const editFormRef = ref<FormInstance>()
const editForm = reactive({
  id: 0,
  title: '',
  category: '',
})

// 加载文档列表
const loadDocuments = async () => {
  loading.value = true
  try {
    const res = await getDocumentList({
      page: page.value,
      pageSize: pageSize.value,
      category: selectedCategory.value || undefined,
      keyword: searchKeyword.value || undefined,
    })
    if (res.code === 0) {
      documents.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch {
    ElMessage.error('加载文档列表失败')
  } finally {
    loading.value = false
  }
}

// 加载分类列表
const loadCategories = async () => {
  try {
    const res = await getDocumentCategories()
    if (res.code === 0) {
      categories.value = res.data || []
    }
  } catch {
    // 忽略
  }
}

// 搜索
const handleSearch = () => {
  page.value = 1
  loadDocuments()
}

// 分类筛选
const handleCategoryChange = () => {
  page.value = 1
  loadDocuments()
}

// 分页
const handlePageChange = (p: number) => {
  page.value = p
  loadDocuments()
}

const handleSizeChange = (s: number) => {
  pageSize.value = s
  page.value = 1
  loadDocuments()
}

// 文件大小格式化
const formatFileSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

// 打开上传对话框
const openUploadDialog = () => {
  uploadForm.title = ''
  uploadForm.category = ''
  uploadForm.file = null
  uploadDialogVisible.value = true
}

// 处理文件选择
const handleFileChange = (uploadFile: UploadFile) => {
  uploadForm.file = uploadFile.raw || null
}

// 提交上传
const handleUpload = async () => {
  if (!uploadFormRef.value) return
  await uploadFormRef.value.validate()
  if (!uploadForm.file) {
    ElMessage.warning('请选择要上传的文件')
    return
  }
  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('title', uploadForm.title)
    formData.append('category', uploadForm.category)
    formData.append('file', uploadForm.file)
    const res = await uploadDocument(formData)
    if (res.code === 0) {
      ElMessage.success('文档上传成功')
      uploadDialogVisible.value = false
      loadDocuments()
      loadCategories()
    } else {
      ElMessage.error(res.msg || '上传失败')
    }
  } catch {
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
  }
}

// 打开编辑对话框
const openEditDialog = (doc: DocumentItem) => {
  editForm.id = doc.id
  editForm.title = doc.title
  editForm.category = doc.category
  editDialogVisible.value = true
}

// 提交编辑
const handleEdit = async () => {
  if (!editFormRef.value) return
  await editFormRef.value.validate()
  editSaving.value = true
  try {
    const res = await updateDocument(editForm.id, {
      title: editForm.title,
      category: editForm.category,
    })
    if (res.code === 0) {
      ElMessage.success('文档更新成功')
      editDialogVisible.value = false
      loadDocuments()
    } else {
      ElMessage.error(res.msg || '更新失败')
    }
  } catch {
    ElMessage.error('更新失败')
  } finally {
    editSaving.value = false
  }
}

// 删除文档
const handleDelete = async (doc: DocumentItem) => {
  await ElMessageBox.confirm(`确定要删除文档「${doc.title}」吗？`, '删除确认', {
    type: 'warning',
  })
  try {
    const res = await deleteDocument(doc.id)
    if (res.code === 0) {
      ElMessage.success('文档已删除')
      loadDocuments()
      loadCategories()
    } else {
      ElMessage.error(res.msg || '删除失败')
    }
  } catch {
    // 用户取消
  }
}

// 下载文档
const handleDownload = async (doc: DocumentItem) => {
  try {
    const res = await getDocumentDownloadUrl(doc.id)
    if (res.code === 0 && res.data?.url) {
      window.open(res.data.url, '_blank')
    } else {
      ElMessage.error('获取下载链接失败')
    }
  } catch {
    ElMessage.error('获取下载链接失败')
  }
}

onMounted(() => {
  loadDocuments()
  loadCategories()
})
</script>

<template>
  <div class="document-center">
    <PageHeader title="文档中心" />

    <el-card shadow="never">
      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-select
          v-model="selectedCategory"
          placeholder="全部分类"
          clearable
          style="width: 180px"
          @change="handleCategoryChange"
        >
          <el-option
            v-for="cat in categories"
            :key="cat"
            :label="cat"
            :value="cat"
          />
        </el-select>
        <el-input
          v-model="searchKeyword"
          placeholder="搜索文档标题"
          clearable
          style="width: 280px"
          @keyup.enter="handleSearch"
          @clear="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" :icon="Upload" @click="openUploadDialog">
          上传文档
        </el-button>
      </div>

      <!-- 文档表格 -->
      <el-table :data="documents" v-loading="loading" stripe>
        <el-table-column prop="title" label="文档标题" min-width="200" />
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.category }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="file_name" label="文件名" min-width="160" show-overflow-tooltip />
        <el-table-column prop="file_size" label="大小" width="100">
          <template #default="{ row }">
            {{ formatFileSize(row.file_size) }}
          </template>
        </el-table-column>
        <el-table-column prop="uploader" label="上传者" width="120">
          <template #default="{ row }">
            {{ row.uploader?.display_name || row.uploader?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="上传时间" width="170">
          <template #default="{ row }">
            {{ row.created_at?.replace('T', ' ').slice(0, 19) || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" :icon="Download" @click="handleDownload(row)">
              下载
            </el-button>
            <el-button link type="primary" :icon="Edit" @click="openEditDialog(row)">
              编辑
            </el-button>
            <el-button link type="danger" :icon="Delete" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty description="暂无文档" />
        </template>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper" v-if="total > 0">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
    </el-card>

    <!-- 上传对话框 -->
    <el-dialog v-model="uploadDialogVisible" title="上传文档" width="500px" destroy-on-close>
      <el-form ref="uploadFormRef" :model="uploadForm" :rules="uploadRules" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="uploadForm.title" placeholder="请输入文档标题" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="uploadForm.category" placeholder="请选择分类" filterable allow-create style="width: 100%">
            <el-option v-for="cat in categories" :key="cat" :label="cat" :value="cat" />
          </el-select>
        </el-form-item>
        <el-form-item label="文件">
          <el-upload
            :auto-upload="false"
            :limit="1"
            :on-change="handleFileChange"
            :show-file-list="true"
          >
            <el-button type="primary" plain>选择文件</el-button>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="uploading" @click="handleUpload">上传</el-button>
      </template>
    </el-dialog>

    <!-- 编辑对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑文档" width="500px" destroy-on-close>
      <el-form ref="editFormRef" :model="editForm" :rules="uploadRules" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="editForm.title" placeholder="请输入文档标题" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="editForm.category" placeholder="请选择分类" filterable allow-create style="width: 100%">
            <el-option v-for="cat in categories" :key="cat" :label="cat" :value="cat" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="editSaving" @click="handleEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.document-center {
  padding: 24px;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
