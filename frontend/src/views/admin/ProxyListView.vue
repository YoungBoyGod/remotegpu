<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Refresh, Delete, View } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/common/PageHeader.vue'
import { getProxyNodes, deleteProxyNode } from '@/api/proxy'
import type { ProxyNode } from '@/api/proxy'

const loading = ref(false)
const nodes = ref<ProxyNode[]>([])
const statusFilter = ref('')
let refreshTimer: ReturnType<typeof setInterval> | null = null

const filteredNodes = computed(() => {
  if (!statusFilter.value) return nodes.value
  return nodes.value.filter(n => n.status === statusFilter.value)
})

const onlineCount = computed(() => nodes.value.filter(n => n.status === 'online').length)
const offlineCount = computed(() => nodes.value.filter(n => n.status === 'offline').length)

const loadNodes = async () => {
  loading.value = true
  try {
    const res = await getProxyNodes()
    nodes.value = res.data || []
  } catch (error) {
    console.error('加载 Proxy 节点列表失败:', error)
    ElMessage.error('加载 Proxy 节点列表失败')
  } finally {
    loading.value = false
  }
}

const handleDelete = (row: ProxyNode) => {
  ElMessageBox.confirm(
    `确定要删除 Proxy 节点「${row.name}」吗？此操作不可恢复。`,
    '删除确认',
    { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
  ).then(async () => {
    try {
      await deleteProxyNode(row.id)
      ElMessage.success('删除成功')
      await loadNodes()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

const portUsagePercent = (row: ProxyNode) => {
  const total = row.range_end - row.range_start + 1
  if (total <= 0) return 0
  return Math.round((row.used_ports / total) * 100)
}

const formatDateTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN')
}

// 详情弹窗
const detailVisible = ref(false)
const detailNode = ref<ProxyNode | null>(null)

const showDetail = (row: ProxyNode) => {
  detailNode.value = row
  detailVisible.value = true
}

onMounted(() => {
  loadNodes()
  refreshTimer = setInterval(loadNodes, 30000)
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
})
</script>

<template>
  <div class="proxy-list-view">
    <PageHeader title="Proxy 管理" subtitle="管理 Proxy 节点和端口映射">
      <template #actions>
        <el-button :icon="Refresh" @click="loadNodes">刷新</el-button>
      </template>
    </PageHeader>

    <!-- 统计卡片 -->
    <div class="stat-cards">
      <el-card class="stat-card" shadow="hover" @click="statusFilter = ''">
        <div class="stat-number">{{ nodes.length }}</div>
        <div class="stat-label">节点总数</div>
      </el-card>
      <el-card class="stat-card stat-online" shadow="hover" @click="statusFilter = 'online'">
        <div class="stat-number">{{ onlineCount }}</div>
        <div class="stat-label">在线</div>
      </el-card>
      <el-card class="stat-card stat-offline" shadow="hover" @click="statusFilter = 'offline'">
        <div class="stat-number">{{ offlineCount }}</div>
        <div class="stat-label">离线</div>
      </el-card>
    </div>

    <!-- 筛选 -->
    <el-card class="filter-card">
      <div class="filter-container">
        <el-select v-model="statusFilter" placeholder="状态筛选" clearable style="width: 150px">
          <el-option label="在线" value="online" />
          <el-option label="离线" value="offline" />
        </el-select>
        <span class="filter-tip">共 {{ filteredNodes.length }} 条记录</span>
      </div>
    </el-card>

    <!-- 节点列表 -->
    <el-table :data="filteredNodes" v-loading="loading" border style="width: 100%">
      <el-table-column prop="id" label="ID" width="220" show-overflow-tooltip />
      <el-table-column prop="name" label="名称" min-width="120" />
      <el-table-column label="地址" min-width="160">
        <template #default="{ row }">
          {{ row.host }}:{{ row.http_port }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.status === 'online' ? 'success' : 'danger'" size="small">
            {{ row.status === 'online' ? '在线' : '离线' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="端口范围" width="140">
        <template #default="{ row }">
          {{ row.range_start }} - {{ row.range_end }}
        </template>
      </el-table-column>
      <el-table-column label="活跃映射" width="90" align="center">
        <template #default="{ row }">
          {{ row.active_mappings }}
        </template>
      </el-table-column>
      <el-table-column label="端口使用率" width="150">
        <template #default="{ row }">
          <el-progress :percentage="portUsagePercent(row)" :stroke-width="14" :text-inside="true" />
        </template>
      </el-table-column>
      <el-table-column label="最后心跳" width="180">
        <template #default="{ row }">
          {{ formatDateTime(row.last_heartbeat) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" :icon="View" @click="showDetail(row)">详情</el-button>
          <el-button link type="danger" :icon="Delete" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" title="Proxy 节点详情" width="560px">
      <template v-if="detailNode">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID" :span="2">{{ detailNode.id }}</el-descriptions-item>
          <el-descriptions-item label="名称">{{ detailNode.name }}</el-descriptions-item>
          <el-descriptions-item label="版本">{{ detailNode.version }}</el-descriptions-item>
          <el-descriptions-item label="主机">{{ detailNode.host }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="detailNode.status === 'online' ? 'success' : 'danger'" size="small">
              {{ detailNode.status === 'online' ? '在线' : '离线' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="API 端口">{{ detailNode.api_port }}</el-descriptions-item>
          <el-descriptions-item label="HTTP 端口">{{ detailNode.http_port }}</el-descriptions-item>
          <el-descriptions-item label="端口范围">{{ detailNode.range_start }} - {{ detailNode.range_end }}</el-descriptions-item>
          <el-descriptions-item label="已用端口">{{ detailNode.used_ports }}</el-descriptions-item>
          <el-descriptions-item label="活跃映射">{{ detailNode.active_mappings }}</el-descriptions-item>
          <el-descriptions-item label="最后心跳">{{ formatDateTime(detailNode.last_heartbeat) }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDateTime(detailNode.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ formatDateTime(detailNode.updated_at) }}</el-descriptions-item>
        </el-descriptions>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.proxy-list-view {
  padding: 24px;
  background: #f5f7fa;
  min-height: 100%;
}

.stat-cards {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  flex: 1;
  cursor: pointer;
  text-align: center;
}

.stat-number {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
}

.stat-online .stat-number {
  color: #67c23a;
}

.stat-offline .stat-number {
  color: #f56c6c;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-tip {
  font-size: 13px;
  color: #909399;
}
</style>
