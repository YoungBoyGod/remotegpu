<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const job = ref({
  id: '',
  name: '',
  status: 'running',
  progress: 0,
  startTime: '',
  logs: '',
})

const loadJob = async () => {
  const id = route.params.id
  job.value = {
    id: id as string,
    name: 'ResNet50 训练',
    status: 'running',
    progress: 65,
    startTime: '2026-01-26 14:30',
    logs: '[2026-01-26 14:30:00] 开始训练...\n[2026-01-26 14:31:00] Epoch 1/100, Loss: 2.345\n[2026-01-26 14:32:00] Epoch 2/100, Loss: 2.123\n',
  }
}

onMounted(() => {
  loadJob()
})
</script>

<template>
  <div class="training-detail">
    <div class="page-header">
      <div>
        <h1>{{ job.name }}</h1>
        <el-tag v-if="job.status === 'running'" type="primary">运行中</el-tag>
        <el-tag v-else-if="job.status === 'completed'" type="success">已完成</el-tag>
      </div>
      <el-button @click="router.push('/training')">返回列表</el-button>
    </div>

    <el-card class="info-card">
      <template #header>任务信息</template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="任务ID">{{ job.id }}</el-descriptions-item>
        <el-descriptions-item label="开始时间">{{ job.startTime }}</el-descriptions-item>
        <el-descriptions-item label="进度" :span="2">
          <el-progress :percentage="job.progress" />
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card class="logs-card">
      <template #header>训练日志</template>
      <el-input
        v-model="job.logs"
        type="textarea"
        :rows="20"
        readonly
        class="logs-textarea"
      />
    </el-card>
  </div>
</template>

<style scoped>
.training-detail {
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
  margin: 0 12px 0 0;
  display: inline-block;
}

.info-card,
.logs-card {
  margin-bottom: 20px;
}

.logs-textarea {
  font-family: 'Courier New', monospace;
  font-size: 13px;
}
</style>
