<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useRoleNavigation } from '@/composables/useRoleNavigation'

const route = useRoute()
const router = useRouter()
const { navigateTo } = useRoleNavigation()

const dataset = ref({
  id: '',
  name: '',
  version: '',
  size: '',
  fileCount: 0,
  description: '',
  createdAt: '',
})

const loadDataset = async () => {
  const id = route.params.id
  dataset.value = {
    id: id as string,
    name: 'ImageNet 2012',
    version: 'v1.0',
    size: '150 GB',
    fileCount: 1281167,
    description: 'ImageNet 大规模视觉识别挑战数据集',
    createdAt: '2026-01-20',
  }
}

onMounted(() => {
  loadDataset()
})
</script>

<template>
  <div class="dataset-detail">
    <div class="page-header">
      <h1>{{ dataset.name }}</h1>
      <el-button @click="navigateTo('/datasets')">返回列表</el-button>
    </div>

    <el-card>
      <template #header>基本信息</template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="数据集ID">{{ dataset.id }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ dataset.version }}</el-descriptions-item>
        <el-descriptions-item label="大小">{{ dataset.size }}</el-descriptions-item>
        <el-descriptions-item label="文件数量">{{ dataset.fileCount }}</el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ dataset.createdAt }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ dataset.description }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<style scoped>
.dataset-detail {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-header h1 {
  font-size: 24px;
  font-weight: 600;
  margin: 0;
}
</style>
