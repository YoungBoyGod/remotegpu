<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import FilterBar from '@/components/common/FilterBar.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import ConfigurableTable from '@/components/common/ConfigurableTable.vue'
import { trainingColumns } from '@/config/tableColumns'

const router = useRouter()

interface TrainingJob {
  id: string
  name: string
  status: string
  progress: number
  gpuCount: number
  startTime: string
  runningTime: string
}

const jobs = ref<TrainingJob[]>([])
const loading = ref(false)
const searchText = ref('')

// 过滤后的任务列表
const filteredJobs = computed(() => {
  let result = jobs.value

  // 搜索过滤
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(job =>
      job.name.toLowerCase().includes(search)
    )
  }

  return result
})

const loadJobs = async () => {
  loading.value = true
  try {
    jobs.value = [
      {
        id: 'job-001',
        name: 'ResNet50 训练',
        status: 'running',
        progress: 65,
        gpuCount: 2,
        startTime: '2026-01-26 14:30',
        runningTime: '1小时20分',
      },
      {
        id: 'job-002',
        name: 'BERT 微调',
        status: 'completed',
        progress: 100,
        gpuCount: 4,
        startTime: '2026-01-26 10:00',
        runningTime: '3小时45分',
      },
    ]
  } catch (error) {
    ElMessage.error('加载任务列表失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadJobs()
})
</script>

<template>
  <div class="training-list">
    <PageHeader title="训练任务">
      <template #actions>
        <el-button type="primary" :icon="Plus" @click="router.push('/training/create')">
          创建任务
        </el-button>
      </template>
    </PageHeader>

    <FilterBar
      v-model:search-value="searchText"
      search-placeholder="搜索任务名称"
    />

    <ConfigurableTable
      :columns="trainingColumns"
      :data="filteredJobs"
      :loading="loading"
    >
      <!-- 任务名称列 -->
      <template #name="{ row }">
        <el-link type="primary" @click="router.push(`/training/${row.id}`)">
          {{ row.name }}
        </el-link>
      </template>

      <!-- 状态列 -->
      <template #status="{ row }">
        <StatusTag :status="row.status === 'running' ? '运行中' : row.status === 'completed' ? '已完成' : '失败'" />
      </template>

      <!-- 进度列 -->
      <template #progress="{ row }">
        <el-progress :percentage="row.progress" />
      </template>

      <!-- 操作列 -->
      <template #actions="{ row }">
        <el-button type="primary" size="small" @click="router.push(`/training/${row.id}`)">
          查看
        </el-button>
        <el-button v-if="row.status === 'running'" type="danger" size="small">
          停止
        </el-button>
      </template>
    </ConfigurableTable>
  </div>
</template>

<style scoped>
.training-list {
  padding: 24px;
}
</style>
