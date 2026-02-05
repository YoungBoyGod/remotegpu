<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { getImageList, syncImages } from '@/api/admin'
import PageHeader from '@/components/common/PageHeader.vue'
import FilterBar from '@/components/common/FilterBar.vue'

interface Image {
  id: number
  name: string
  display_name: string
  category: string
  framework: string
  cuda_version: string
  status: string
  created_at: string
}

const images = ref<Image[]>([])
const loading = ref(false)
const syncing = ref(false)
const searchText = ref('')
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const loadImages = async () => {
  loading.value = true
  try {
    const res = await getImageList({
      page: page.value,
      pageSize: pageSize.value,
      // 这里的搜索目前只能过滤 category/framework，后续后端支持 fuzzy search 更好
      // 暂时前端过滤或忽略 search
    })
    images.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSync = async () => {
  syncing.value = true
  try {
    await syncImages()
    ElMessage.success('同步任务已触发')
    // 延迟一下刷新列表
    setTimeout(loadImages, 1000)
  } catch (error) {
    // Error handled by interceptor
  } finally {
    syncing.value = false
  }
}

const handlePageChange = (val: number) => {
  page.value = val
  loadImages()
}

onMounted(() => {
  loadImages()
})
</script>

<template>
  <div class="image-list">
    <PageHeader title="镜像管理">
      <template #actions>
        <el-button type="primary" :icon="Refresh" :loading="syncing" @click="handleSync">
          同步镜像
        </el-button>
      </template>
    </PageHeader>

    <div class="image-grid">
      <el-table :data="images" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag>{{ row.category }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="framework" label="框架" width="150" />
        <el-table-column prop="cuda_version" label="CUDA" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.image-list {
  padding: 24px;
}

.image-grid {
  margin-top: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
