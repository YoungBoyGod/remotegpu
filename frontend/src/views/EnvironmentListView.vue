<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import { useRoleNavigation } from '@/composables/useRoleNavigation'
import PageHeader from '@/components/common/PageHeader.vue'
import FilterBar from '@/components/common/FilterBar.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import ConfigurableTable from '@/components/common/ConfigurableTable.vue'
import { environmentColumns } from '@/config/tableColumns'

const router = useRouter()
const { navigateTo } = useRoleNavigation()

interface Environment {
  id: string
  name: string
  status: string
  image: string
  gpu: string
  cpu: number
  memory: number
  runningTime: string
  createdAt: string
}

const environments = ref<Environment[]>([])
const loading = ref(false)
const searchText = ref('')
const statusFilter = ref('')

// 过滤后的环境列表
const filteredEnvironments = computed(() => {
  let result = environments.value

  // 搜索过滤
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(env =>
      env.name.toLowerCase().includes(search) ||
      env.image.toLowerCase().includes(search)
    )
  }

  // 状态过滤
  if (statusFilter.value) {
    result = result.filter(env => env.status === statusFilter.value)
  }

  return result
})

// 加载环境列表
const loadEnvironments = async () => {
  loading.value = true
  try {
    // 模拟数据
    environments.value = [
      {
        id: 'env-001',
        name: 'PyTorch 训练环境',
        status: 'running',
        image: 'pytorch/pytorch:2.0-cuda11.8',
        gpu: 'Tesla V100 x2',
        cpu: 8,
        memory: 32,
        runningTime: '2小时30分',
        createdAt: '2026-01-26 10:30',
      },
      {
        id: 'env-002',
        name: 'TensorFlow 开发',
        status: 'stopped',
        image: 'tensorflow/tensorflow:2.13-gpu',
        gpu: 'RTX 4090 x1',
        cpu: 4,
        memory: 16,
        runningTime: '-',
        createdAt: '2026-01-25 15:20',
      },
    ]
  } catch (error) {
    ElMessage.error('加载环境列表失败')
  } finally {
    loading.value = false
  }
}

// 启动环境
const startEnvironment = async (id: string) => {
  try {
    ElMessage.success('环境启动中...')
    await loadEnvironments()
  } catch (error) {
    ElMessage.error('启动失败')
  }
}

// 停止环境
const stopEnvironment = async (id: string) => {
  try {
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
        <el-link type="primary" @click="router.push(`/environments/${row.id}`)">
          {{ row.name }}
        </el-link>
      </template>

      <!-- 状态列 -->
      <template #status="{ row }">
        <StatusTag :status="row.status === 'running' ? '运行中' : row.status === 'stopped' ? '已停止' : '错误'" />
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
