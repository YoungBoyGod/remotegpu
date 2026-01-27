<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import { useRoleNavigation } from '@/composables/useRoleNavigation'
import PageHeader from '@/components/common/PageHeader.vue'
import FilterBar from '@/components/common/FilterBar.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import ConfigurableTable from '@/components/common/ConfigurableTable.vue'
import { datasetColumns } from '@/config/tableColumns'

const router = useRouter()
const { navigateTo } = useRoleNavigation()

interface Dataset {
  id: string
  name: string
  version: string
  size: string
  fileCount: number
  visibility: string
  createdAt: string
}

const datasets = ref<Dataset[]>([])
const loading = ref(false)
const searchText = ref('')

// 过滤后的数据集列表
const filteredDatasets = computed(() => {
  let result = datasets.value

  // 搜索过滤
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(ds =>
      ds.name.toLowerCase().includes(search) ||
      ds.version.toLowerCase().includes(search)
    )
  }

  return result
})

const loadDatasets = async () => {
  loading.value = true
  try {
    datasets.value = [
      {
        id: 'ds-001',
        name: 'ImageNet 2012',
        version: 'v1.0',
        size: '150 GB',
        fileCount: 1281167,
        visibility: 'public',
        createdAt: '2026-01-20',
      },
      {
        id: 'ds-002',
        name: 'COCO 2017',
        version: 'v2.1',
        size: '25 GB',
        fileCount: 123287,
        visibility: 'private',
        createdAt: '2026-01-22',
      },
    ]
  } catch (error) {
    ElMessage.error('加载数据集列表失败')
  } finally {
    loading.value = false
  }
}

const deleteDataset = async (id: string) => {
  try {
    await ElMessageBox.confirm('确定要删除这个数据集吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    ElMessage.success('数据集已删除')
    await loadDatasets()
  } catch (error) {
    // 用户取消
  }
}

onMounted(() => {
  loadDatasets()
})
</script>

<template>
  <div class="dataset-list">
    <PageHeader title="数据集管理">
      <template #actions>
        <el-button type="primary" :icon="Upload" @click="navigateTo('/datasets/upload')">
          上传数据集
        </el-button>
      </template>
    </PageHeader>

    <FilterBar
      v-model:search-value="searchText"
      search-placeholder="搜索数据集名称"
    />

    <ConfigurableTable
      :columns="datasetColumns"
      :data="filteredDatasets"
      :loading="loading"
    >
      <!-- 数据集名称列 -->
      <template #name="{ row }">
        <el-link type="primary" @click="router.push(`/datasets/${row.id}`)">
          {{ row.name }}
        </el-link>
      </template>

      <!-- 可见性列 -->
      <template #visibility="{ row }">
        <StatusTag :status="row.visibility === 'public' ? '公开' : '私有'" :type="row.visibility === 'public' ? 'success' : 'info'" />
      </template>

      <!-- 操作列 -->
      <template #actions="{ row }">
        <el-button type="primary" size="small" @click="router.push(`/datasets/${row.id}`)">
          查看
        </el-button>
        <el-button type="danger" size="small" @click="deleteDataset(row.id)">
          删除
        </el-button>
      </template>
    </ConfigurableTable>
  </div>
</template>

<style scoped>
.dataset-list {
  padding: 24px;
}
</style>
