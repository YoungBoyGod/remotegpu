<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import { useRoleNavigation } from '@/composables/useRoleNavigation'
import PageHeader from '@/components/common/PageHeader.vue'
import FilterBar from '@/components/common/FilterBar.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import ConfigurableTable from '@/components/common/ConfigurableTable.vue'
import { environmentColumns } from '@/config/tableColumns'
import type { Environment } from '@/api/environment/types'
import {
  getEnvironmentList,
  startEnvironment as startEnv,
  stopEnvironment as stopEnv,
  deleteEnvironment as deleteEnv,
} from '@/api/environment'

const { navigateTo } = useRoleNavigation()

const environments = ref<Environment[]>([])
const loading = ref(false)
const searchText = ref('')
const statusFilter = ref('')

const formatDateTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

const formatMemoryToGB = (value: number) => {
  if (!value) return 0
  const gb = value / 1024
  return gb < 1 ? 1 : Math.round(gb)
}

const formatDuration = (startAt?: string | null) => {
  if (!startAt) return '-'
  const start = new Date(startAt)
  if (Number.isNaN(start.getTime())) return '-'
  const diffMs = Date.now() - start.getTime()
  if (diffMs < 0) return '-'
  const diffMinutes = Math.floor(diffMs / 60000)
  const hours = Math.floor(diffMinutes / 60)
  const minutes = diffMinutes % 60
  if (hours > 0) {
    return `${hours}小时${minutes}分`
  }
  return `${minutes}分`
}

const displayEnvironments = computed(() => {
  return environments.value.map(env => ({
    id: env.id,
    name: env.name,
    status: env.status,
    image: env.image,
    gpu: env.gpu,
    cpu: env.cpu,
    memory: formatMemoryToGB(env.memory),
    runningTime: env.status === 'running' ? formatDuration(env.started_at) : '-',
    createdAt: formatDateTime(env.created_at),
  }))
})

const statusTextMap: Record<string, string> = {
  running: '运行中',
  stopped: '已停止',
  creating: '创建中',
  deleting: '删除中',
  error: '错误',
}

// 过滤后的环境列表
const filteredEnvironments = computed(() => {
  let result = displayEnvironments.value

  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(env =>
      env.name.toLowerCase().includes(search) ||
      env.image.toLowerCase().includes(search)
    )
  }

  if (statusFilter.value) {
    result = result.filter(env => env.status === statusFilter.value)
  }

  return result
})

// 加载环境列表
const loadEnvironments = async () => {
  loading.value = true
  try {
    const response = await getEnvironmentList()
    environments.value = response.data
  } catch (error) {
    ElMessage.error('加载环境列表失败')
  } finally {
    loading.value = false
  }
}

// 启动环境
const startEnvironment = async (id: string) => {
  try {
    await startEnv(id)
    ElMessage.success('环境启动中...')
    await loadEnvironments()
  } catch (error) {
    ElMessage.error('启动失败')
  }
}

// 停止环境
const stopEnvironment = async (id: string) => {
  try {
    await stopEnv(id)
    ElMessage.success('环境已停止')
    await loadEnvironments()
  } catch (error) {
    ElMessage.error('停止失败')
  }
}

// 删除环境
const deleteEnvironment = async (id: string) => {
  try {
    await ElMessageBox.confirm('确定要删除这个环境吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await deleteEnv(id)
    ElMessage.success('环境已删除')
    await loadEnvironments()
  } catch (error) {
    // 用户取消
  }
}

onMounted(() => {
  loadEnvironments()
})
</script>

<template>
  <div class="environment-list">
    <PageHeader title="开发环境">
      <template #actions>
        <el-button type="primary" :icon="Plus" @click="navigateTo('/environments/create')">
          创建环境
        </el-button>
      </template>
    </PageHeader>

    <FilterBar
      v-model:search-value="searchText"
      search-placeholder="搜索环境名称"
    >
      <template #filters>
        <el-select v-model="statusFilter" placeholder="状态筛选" style="width: 150px" clearable>
          <el-option label="运行中" value="running" />
          <el-option label="已停止" value="stopped" />
          <el-option label="错误" value="error" />
        </el-select>
      </template>
      <template #actions>
        <el-button :icon="Refresh" @click="loadEnvironments">刷新</el-button>
      </template>
    </FilterBar>

    <ConfigurableTable
      :columns="environmentColumns"
      :data="filteredEnvironments"
      :loading="loading"
    >
      <!-- 环境名称列 -->
      <template #name="{ row }">
        <el-link type="primary" @click="navigateTo(`/environments/${row.id}`)">
          {{ row.name }}
        </el-link>
      </template>

      <!-- 状态列 -->
      <template #status="{ row }">
        <StatusTag :status="statusTextMap[row.status] || row.status" />
      </template>

      <!-- CPU/内存列 -->
      <template #cpu-memory="{ row }">
        {{ row.cpu }}核 / {{ row.memory }}GB
      </template>

      <!-- 操作列 -->
      <template #actions="{ row }">
          <el-button
            v-if="row.status === 'stopped'"
            type="success"
            size="small"
            @click="startEnvironment(row.id)"
          >
            启动
          </el-button>
          <el-button
            v-if="row.status === 'running'"
            type="warning"
            size="small"
            @click="stopEnvironment(row.id)"
          >
            停止
          </el-button>
          <el-button type="danger" size="small" @click="deleteEnvironment(row.id)">
            删除
          </el-button>
      </template>
    </ConfigurableTable>
  </div>
</template>

<style scoped>
.environment-list {
  padding: 24px;
}
</style>
