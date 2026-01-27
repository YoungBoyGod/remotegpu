<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useRoleNavigation } from '@/composables/useRoleNavigation'
import StatusTag from '@/components/common/StatusTag.vue'

const router = useRouter()
const { navigateTo } = useRoleNavigation()

interface Environment {
  id: string
  name: string
  status: string
  gpu: string
  createdAt: string
}

const environments = ref<Environment[]>([])
const loading = ref(false)

// 加载最近环境
const loadRecentEnvironments = async () => {
  loading.value = true
  try {
    // TODO: 调用 API
    // const response = await getEnvironments({ limit: 5 })
    // environments.value = response.data.items

    // 模拟数据
    environments.value = [
      {
        id: 'env-001',
        name: 'PyTorch 训练环境',
        status: 'running',
        gpu: 'Tesla V100 x2',
        createdAt: '2026-01-26 10:30',
      },
      {
        id: 'env-002',
        name: 'TensorFlow 开发',
        status: 'stopped',
        gpu: 'RTX 4090 x1',
        createdAt: '2026-01-25 15:20',
      },
      {
        id: 'env-003',
        name: 'CUDA 测试环境',
        status: 'running',
        gpu: 'A100 x1',
        createdAt: '2026-01-24 09:15',
      },
    ]
  } catch (error) {
    ElMessage.error('加载环境列表失败')
  } finally {
    loading.value = false
  }
}

// 查看环境详情
const viewEnvironment = (id: string) => {
  router.push(`/environments/${id}`)
}

onMounted(() => {
  loadRecentEnvironments()
})
</script>

<template>
  <div class="recent-environments">
    <div class="section-header">
      <h3>最近使用的环境</h3>
      <el-link type="primary" :underline="false" @click="navigateTo('/environments')">
        查看全部
      </el-link>
    </div>

    <el-table :data="environments" :loading="loading" style="width: 100%">
      <el-table-column prop="name" label="环境名称" min-width="150">
        <template #default="{ row }">
          <el-link type="primary" @click="viewEnvironment(row.id)">
            {{ row.name }}
          </el-link>
        </template>
      </el-table-column>

      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <StatusTag :status="row.status === 'running' ? '运行中' : '已停止'" />
        </template>
      </el-table-column>

      <el-table-column prop="gpu" label="GPU" width="150" />

      <el-table-column prop="createdAt" label="创建时间" width="150" />
    </el-table>
  </div>
</template>

<style scoped>
.recent-environments {
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
