<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoleNavigation } from '@/composables/useRoleNavigation'
import type { Environment } from '@/api/environment/types'
import {
  getEnvironmentDetail,
  startEnvironment as startEnv,
  stopEnvironment as stopEnv,
  deleteEnvironment as deleteEnv,
} from '@/api/environment'

const route = useRoute()
const { navigateTo } = useRoleNavigation()
const activeTab = ref('overview')
const environment = ref<Environment | null>(null)

const formatDateTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

const formatMemoryToGB = (value?: number | null) => {
  if (!value) return '-'
  const gb = value / 1024
  return `${gb < 1 ? 1 : Math.round(gb)} GB`
}

const formatStorage = (value?: number | null) => {
  if (!value) return '-'
  const gb = value / 1024
  return `${gb < 1 ? 1 : Math.round(gb)} GB`
}

const statusLabel = computed(() => {
  if (!environment.value) return ''
  return environment.value.status === 'running' ? '运行中' :
    environment.value.status === 'stopped' ? '已停止' : environment.value.status
})

const loadEnvironment = async () => {
  const id = route.params.id as string
  try {
    const response = await getEnvironmentDetail(id)
    environment.value = response.data
  } catch (error) {
    ElMessage.error('加载环境详情失败')
  }
}

const startEnvironment = async () => {
  if (!environment.value) return
  await startEnv(environment.value.id)
  ElMessage.success('环境启动中...')
  await loadEnvironment()
}

const stopEnvironment = async () => {
  if (!environment.value) return
  await stopEnv(environment.value.id)
  ElMessage.success('环境已停止')
  await loadEnvironment()
}

const deleteEnvironment = async () => {
  if (!environment.value) return
  try {
    await ElMessageBox.confirm('确定要删除这个环境吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await deleteEnv(environment.value.id)
    ElMessage.success('环境已删除')
    navigateTo('/environments')
  } catch (error) {
    // 用户取消
  }
}

onMounted(() => {
  loadEnvironment()
})
</script>

<template>
  <div v-if="environment" class="environment-detail">
    <div class="page-header">
      <div>
        <h1>{{ environment.name }}</h1>
        <el-tag v-if="environment.status === 'running'" type="success">{{ statusLabel }}</el-tag>
        <el-tag v-else type="info">{{ statusLabel }}</el-tag>
      </div>
      <div class="actions">
        <el-button v-if="environment.status === 'stopped'" type="success" @click="startEnvironment">
          启动
        </el-button>
        <el-button v-if="environment.status === 'running'" type="warning" @click="stopEnvironment">
          停止
        </el-button>
        <el-button type="danger" @click="deleteEnvironment">删除</el-button>
        <el-button @click="navigateTo('/environments')">返回列表</el-button>
      </div>
    </div>

    <el-tabs v-model="activeTab">
      <el-tab-pane label="概览" name="overview">
        <div class="overview-content">
          <el-card class="info-card">
            <template #header>基本信息</template>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="环境ID">{{ environment.id }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ formatDateTime(environment.created_at) }}</el-descriptions-item>
              <el-descriptions-item label="主机ID">{{ environment.host_id }}</el-descriptions-item>
              <el-descriptions-item label="镜像">{{ environment.image }}</el-descriptions-item>
            </el-descriptions>
          </el-card>

          <el-card class="info-card">
            <template #header>资源配置</template>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="CPU">{{ environment.cpu }} 核</el-descriptions-item>
              <el-descriptions-item label="内存">{{ formatMemoryToGB(environment.memory) }}</el-descriptions-item>
              <el-descriptions-item label="GPU">{{ environment.gpu }} 张</el-descriptions-item>
              <el-descriptions-item label="存储">{{ formatStorage(environment.storage) }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
  <el-empty v-else description="正在加载环境信息..." />
</template>

<style scoped>
.environment-detail {
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

.actions {
  display: flex;
  gap: 12px;
}

.overview-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-card {
  margin-bottom: 0;
}

</style>
