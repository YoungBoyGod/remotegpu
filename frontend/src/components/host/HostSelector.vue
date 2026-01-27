<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { Host, HostFilterParams } from '@/api/host/types'
import { getAvailableHosts } from '@/api/host'
import { useHostStore } from '@/stores/host'
import HostFilterBar from './HostFilterBar.vue'
import HostCard from './HostCard.vue'
import EmptyState from '@/components/common/EmptyState.vue'

interface Props {
  selectionMode?: 'single' | 'multiple'
}

const props = withDefaults(defineProps<Props>(), {
  selectionMode: 'single'
})

const emit = defineEmits<{
  select: [host: Host]
}>()

const hostStore = useHostStore()

// 状态
const hosts = ref<Host[]>([])
const loading = ref(false)
const filters = ref<HostFilterParams>({
  region: '',
  gpu_count: '',
  gpu_model: '',
  keyword: ''
})
const selectedHostId = ref<number | null>(null)

// 分页
const pagination = ref({
  page: 1,
  pageSize: 12,
  total: 0
})

// 加载主机列表
const loadHosts = async () => {
  loading.value = true
  try {
    const params = {
      ...filters.value,
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    }

    const response = await getAvailableHosts(params)
    hosts.value = response.data.items
    pagination.value.total = response.data.total
  } catch (error) {
    ElMessage.error('加载主机列表失败')
    console.error('Failed to load hosts:', error)
  } finally {
    loading.value = false
  }
}

// 处理主机选择
const handleSelectHost = (host: Host) => {
  selectedHostId.value = host.id
  hostStore.selectHost(host)
  emit('select', host)
}

// 处理刷新
const handleRefresh = () => {
  loadHosts()
}

// 清除筛选
const handleClearFilters = () => {
  filters.value = {
    region: '',
    gpu_count: '',
    gpu_model: '',
    keyword: ''
  }
}

// 分页变化
const handlePageChange = (page: number) => {
  pagination.value.page = page
  loadHosts()
}

// 监听过滤器变化
watch(filters, () => {
  pagination.value.page = 1
  loadHosts()
}, { deep: true })

// 组件挂载时加载数据
onMounted(() => {
  loadHosts()
})

// 是否显示空状态
const showEmpty = computed(() => {
  return !loading.value && hosts.value.length === 0
})
</script>

<template>
  <div class="host-selector">
    <HostFilterBar
      v-model="filters"
      @refresh="handleRefresh"
    />

    <div v-loading="loading" class="hosts-container">
      <div v-if="!showEmpty" class="hosts-grid">
        <HostCard
          v-for="host in hosts"
          :key="host.id"
          :host="host"
          :selected="selectedHostId === host.id"
          :selection-mode="selectionMode"
          @select="handleSelectHost"
        />
      </div>

      <EmptyState
        v-if="showEmpty"
        title="未找到符合条件的主机"
        description="请尝试调整筛选条件"
        :show-action="true"
        action-text="清除筛选"
        @action="handleClearFilters"
      />
    </div>

    <div v-if="pagination.total > pagination.pageSize" class="pagination-container">
      <el-pagination
        v-model:current-page="pagination.page"
        :page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<style scoped>
.host-selector {
  width: 100%;
}

.hosts-container {
  min-height: 400px;
}

.hosts-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  padding: 24px 0;
}

@media (max-width: 1200px) {
  .hosts-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .hosts-grid {
    grid-template-columns: 1fr;
  }
}
</style>
