<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useRoleNavigation } from '@/composables/useRoleNavigation'
import StatusTag from '@/components/common/StatusTag.vue'

const { navigateTo } = useRoleNavigation()

interface TrainingJob {
  id: string
  name: string
  status: string
  progress: number
  startTime: string
}

const jobs = ref<TrainingJob[]>([])
const loading = ref(false)

// 加载最近任务
const loadRecentJobs = async () => {
  loading.value = true
  try {
    // TODO: 调用 API
    // const response = await getTrainingJobs({ limit: 5 })
    // jobs.value = response.data.items

    // 模拟数据
    jobs.value = [
      {
        id: 'job-001',
        name: 'ResNet50 训练',
        status: 'running',
        progress: 65,
        startTime: '2026-01-26 14:30',
      },
      {
        id: 'job-002',
        name: 'BERT 微调',
        status: 'completed',
        progress: 100,
        startTime: '2026-01-26 10:00',
      },
      {
        id: 'job-003',
        name: 'YOLOv8 训练',
        status: 'failed',
        progress: 45,
        startTime: '2026-01-25 16:20',
      },
    ]
  } catch (error) {
    ElMessage.error('加载任务列表失败')
  } finally {
    loading.value = false
  }
}

// 查看任务详情
const viewJob = (id: string) => {
  navigateTo(`/training/${id}`)
}

onMounted(() => {
  loadRecentJobs()
})
</script>

<template>
  <div class="recent-jobs">
    <div class="section-header">
      <h3>最近的训练任务</h3>
      <el-link type="primary" :underline="false" @click="navigateTo('/training')">
        查看全部
      </el-link>
    </div>

    <el-table :data="jobs" :loading="loading" style="width: 100%">
      <el-table-column prop="name" label="任务名称" min-width="150">
        <template #default="{ row }">
          <el-link type="primary" @click="viewJob(row.id)">
            {{ row.name }}
          </el-link>
        </template>
      </el-table-column>

      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <StatusTag :status="row.status === 'running' ? '运行中' : row.status === 'completed' ? '已完成' : '失败'" />
        </template>
      </el-table-column>

      <el-table-column prop="progress" label="进度" width="150">
        <template #default="{ row }">
          <el-progress :percentage="row.progress" :status="row.status === 'failed' ? 'exception' : undefined" />
        </template>
      </el-table-column>

      <el-table-column prop="startTime" label="开始时间" width="150" />
    </el-table>
  </div>
</template>

<style scoped>
.recent-jobs {
  background: white;
  padding: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}
</style>
