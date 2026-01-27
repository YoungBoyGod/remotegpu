<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, Edit, Delete, Search } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'

interface ResourceConfig {
  cpu: number
  memory: number
  gpu: number
  disk: number
  network: string
  os: string
  [key: string]: any
}

interface Resource {
  id: number
  name: string
  type: string
  status: string
  region: string
  config: ResourceConfig
  createdAt: string
}

const resources = ref<Resource[]>([])
const loading = ref(false)
const searchKeyword = ref('')
const typeFilter = ref('')
const statusFilter = ref('')

// 编辑对话框
const dialogVisible = ref(false)
const dialogTitle = ref('编辑资源')
const formData = ref<Resource>({
  id: 0,
  name: '',
  type: '',
  status: '',
  region: '',
  config: {
    cpu: 0,
    memory: 0,
    gpu: 0,
    disk: 0,
    network: '',
    os: ''
  },
  createdAt: ''
})

// 配置JSON编辑器
const configJson = ref('')

// 分页
const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0
})

// 过滤后的资源列表
const filteredResources = computed(() => {
  let result = resources.value

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(r =>
      r.name.toLowerCase().includes(keyword) ||
      r.region.toLowerCase().includes(keyword)
    )
  }

  if (typeFilter.value) {
    result = result.filter(r => r.type === typeFilter.value)
  }

  if (statusFilter.value) {
    result = result.filter(r => r.status === statusFilter.value)
  }

  return result
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
        name: '主机-001',
        type: '物理机',
        status: '运行中',
        region: '北京一区',
        config: {
          cpu: 16,
          memory: 64,
          gpu: 2,
          disk: 500,
          network: '10Gbps',
          os: 'Ubuntu 22.04'
        },
        createdAt: '2026-01-20 10:00'
      },
      {
        id: 2,
        name: '主机-002',
        type: '虚拟机',
        status: '已停止',
        region: '上海一区',
        config: {
          cpu: 8,
          memory: 32,
          gpu: 1,
          disk: 250,
          network: '1Gbps',
          os: 'CentOS 7'
        },
        createdAt: '2026-01-21 14:30'
      }
    ]
    pagination.value.total = resources.value.length
  } catch (error) {
    ElMessage.error('加载资源列表失败')
  } finally {
    loading.value = false
  }
}

// 编辑资源
const handleEdit = (resource: Resource) => {
  formData.value = { ...resource, config: { ...resource.config } }
  configJson.value = JSON.stringify(resource.config, null, 2)
  dialogTitle.value = '编辑资源'
  dialogVisible.value = true
}

// 添加资源
const handleAdd = () => {
  formData.value = {
    id: 0,
    name: '',
    type: '物理机',
    status: '已停止',
    region: '',
    config: {
      cpu: 0,
      memory: 0,
      gpu: 0,
      disk: 0,
      network: '',
      os: ''
    },
    createdAt: ''
  }
  configJson.value = JSON.stringify(formData.value.config, null, 2)
  dialogTitle.value = '添加资源'
  dialogVisible.value = true
}

// 保存资源
const handleSave = async () => {
  try {
    // 验证JSON格式
    const config = JSON.parse(configJson.value)
    formData.value.config = config

    // TODO: 调用API保存数据
    ElMessage.success('保存成功')
    dialogVisible.value = false
    await loadResources()
  } catch (error) {
    ElMessage.error('配置JSON格式错误，请检查')
  }
}

// 删除资源
const handleDelete = async (resource: Resource) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除资源 "${resource.name}" 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    ElMessage.success('删除成功')
    await loadResources()
  } catch (error) {
    // 用户取消
  }
}

// 清除筛选
const handleClearFilters = () => {
  searchKeyword.value = ''
  typeFilter.value = ''
  statusFilter.value = ''
}

// 格式化配置显示
const formatConfig = (config: ResourceConfig) => {
  return JSON.stringify(config, null, 2)
}

onMounted(() => {
  loadResources()
})
</script>

<template>
  <div class="resource-list">
    <PageHeader title="资源管理">
      <template #actions>
        <el-button type="primary" :icon="Plus" @click="handleAdd">
          添加资源
        </el-button>
      </template>
    </PageHeader>

    <!-- 搜索和筛选栏 -->
    <div class="filter-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索资源名称或地区"
        clearable
        style="width: 300px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>

      <el-select
        v-model="typeFilter"
        placeholder="资源类型"
        clearable
        style="width: 150px"
      >
        <el-option label="物理机" value="物理机" />
        <el-option label="虚拟机" value="虚拟机" />
        <el-option label="云主机" value="云主机" />
      </el-select>

      <el-select
        v-model="statusFilter"
        placeholder="状态"
        clearable
        style="width: 150px"
      >
        <el-option label="运行中" value="运行中" />
        <el-option label="已停止" value="已停止" />
        <el-option label="维护中" value="维护中" />
      </el-select>

      <el-button @click="handleClearFilters">清除筛选</el-button>
      <el-button :icon="Refresh" @click="loadResources">刷新</el-button>
    </div>

    <!-- 数据表格 -->
    <el-table
      :data="filteredResources"
      :loading="loading"
      style="width: 100%"
      stripe
    >
      <el-table-column prop="name" label="资源名称" width="180" />
      <el-table-column prop="type" label="类型" width="120" />
      <el-table-column prop="status" label="状态" width="120">
        <template #default="{ row }">
          <el-tag
            :type="row.status === '运行中' ? 'success' : 'info'"
          >
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="region" label="地区" width="150" />
      <el-table-column label="配置 (JSON)" min-width="300">
        <template #default="{ row }">
          <el-popover
            placement="left"
            :width="400"
            trigger="hover"
          >
            <template #reference>
              <div class="config-preview">
                {{ JSON.stringify(row.config) }}
              </div>
            </template>
            <pre class="config-json">{{ formatConfig(row.config) }}</pre>
          </el-popover>
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="创建时间" width="180" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button
            type="primary"
            size="small"
            :icon="Edit"
            @click="handleEdit(row)"
          >
            编辑
          </el-button>
          <el-button
            type="danger"
            size="small"
            :icon="Delete"
            @click="handleDelete(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-container">
      <el-pagination
        v-model:current-page="pagination.page"
        :page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next, jumper"
      />
    </div>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="700px"
    >
      <el-form :model="formData" label-width="100px">
        <el-form-item label="资源名称">
          <el-input v-model="formData.name" placeholder="请输入资源名称" />
        </el-form-item>
        <el-form-item label="资源类型">
          <el-select v-model="formData.type" placeholder="请选择资源类型">
            <el-option label="物理机" value="物理机" />
            <el-option label="虚拟机" value="虚拟机" />
            <el-option label="云主机" value="云主机" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="formData.status" placeholder="请选择状态">
            <el-option label="运行中" value="运行中" />
            <el-option label="已停止" value="已停止" />
            <el-option label="维护中" value="维护中" />
          </el-select>
        </el-form-item>
        <el-form-item label="地区">
          <el-input v-model="formData.region" placeholder="请输入地区" />
        </el-form-item>
        <el-form-item label="配置 (JSON)">
          <el-input
            v-model="configJson"
            type="textarea"
            :rows="10"
            placeholder="请输入JSON格式的配置"
          />
          <div class="form-tip">
            支持JSON格式，例如: {"cpu": 16, "memory": 64, "gpu": 2}
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.resource-list {
  padding: 24px;
}

.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  padding: 16px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.pagination-container {
  display: flex;
  justify-content: center;
  padding: 24px 0;
}

.config-preview {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
  color: #409EFF;
}

.config-json {
  margin: 0;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 12px;
  line-height: 1.5;
  overflow-x: auto;
}

.form-tip {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
}
</style>
