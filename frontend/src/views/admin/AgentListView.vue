<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/common/PageHeader.vue'
import { getAgentList } from '@/api/admin'
import type { AgentInfo } from '@/api/admin'

const router = useRouter()
const loading = ref(false)
const agents = ref<AgentInfo[]>([])
const totalCount = ref(0)
const onlineCount = ref(0)
const offlineCount = ref(0)
const statusFilter = ref('')

const filteredAgents = computed(() => {
  if (!statusFilter.value) return agents.value
  return agents.value.filter(a => a.status === statusFilter.value)
})

const loadAgents = async () => {
  loading.value = true
  try {
    const res = await getAgentList()
    agents.value = res.data.agents || []
    totalCount.value = res.data.total
    onlineCount.value = res.data.online
    offlineCount.value = res.data.offline
  } catch (error) {
    console.error('加载 Agent 列表失败:', error)
    ElMessage.error('加载 Agent 列表失败')
  } finally {
    loading.value = false
  }
}

const formatDateTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN')
}

const goToMachine = (machineId: string) => {
  router.push(`/admin/machines/${machineId}`)
}

onMounted(() => {
  loadAgents()
})
</script>

<template>
  <div class="agent-list-view">
    <PageHeader title="Agent 管理" subtitle="查看部署在 GPU 机器上的 Agent 状态">
      <template #actions>
        <el-button :icon="Refresh" @click="loadAgents">刷新</el-button>
      </template>
    </PageHeader>

    <!-- 统计卡片 -->
    <div class="stat-cards">
      <el-card class="stat-card" shadow="hover" @click="statusFilter = ''">
        <div class="stat-number">{{ totalCount }}</div>
        <div class="stat-label">Agent 总数</div>
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
        <span class="filter-tip">共 {{ filteredAgents.length }} 条记录</span>
      </div>
    </el-card>

    <!-- Agent 列表 -->
    <el-table :data="filteredAgents" v-loading="loading" border style="width: 100%">
      <el-table-column label="机器名称" min-width="140">
        <template #default="{ row }">
          <el-button link type="primary" @click="goToMachine(row.machine_id)">
            {{ row.machine_name || row.machine_id }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column prop="ip_address" label="IP 地址" width="140" />
      <el-table-column prop="region" label="区域" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'online' ? 'success' : 'danger'" size="small">
            {{ row.status === 'online' ? '在线' : '离线' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Agent 端口" width="110">
        <template #default="{ row }">
          {{ row.agent_port || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="GPU" min-width="180">
        <template #default="{ row }">
          <template v-if="row.gpu_count > 0">
            <span>{{ row.gpu_count }} 块</span>
            <span class="gpu-model">{{ [...new Set(row.gpu_models)].join(', ') }}</span>
          </template>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="最后心跳" width="180">
        <template #default="{ row }">
          {{ formatDateTime(row.last_heartbeat) }}
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<style scoped>
.agent-list-view {
  padding: 24px;
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

.gpu-model {
  margin-left: 8px;
  color: #909399;
  font-size: 12px;
}
</style>
