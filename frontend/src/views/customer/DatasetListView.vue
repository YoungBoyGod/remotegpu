<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { getDatasetList } from '@/api/customer'

const datasets = ref([])
const loading = ref(false)

const loadDatasets = async () => {
  loading.value = true
  try {
    const res = await getDatasetList({ page: 1, pageSize: 20 })
    datasets.value = res.data.list
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDatasets()
})
</script>

<template>
  <div class="dataset-list-view">
    <PageHeader title="数据集管理">
      <template #actions>
        <el-button type="primary" :icon="Plus">上传数据集</el-button>
      </template>
    </PageHeader>

    <el-card>
      <el-table :data="datasets" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="total_size" label="大小" />
        <el-table-column prop="file_count" label="文件数" />
        <el-table-column prop="status" label="状态" />
        <el-table-column prop="created_at" label="上传时间" width="180">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default>
            <el-button size="small">挂载</el-button>
            <el-button size="small" type="danger">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<style scoped>
.dataset-list-view {
  padding: 24px;
}
</style>
