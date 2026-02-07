<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getSystemConfigs, updateSystemConfigs, getStorageBackends, getStorageStats } from '@/api/admin'
import type { StorageBackend, StorageStats } from '@/api/admin'
import type { SystemConfig } from '@/types/systemConfig'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'

const loading = ref(false)
const saving = ref(false)
const configs = ref<SystemConfig[]>([])
const formValues = ref<Record<string, string>>({})
const activeTab = ref('configs')

// 存储相关
const storageLoading = ref(false)
const storageBackends = ref<StorageBackend[]>([])
const storageStatsData = ref<StorageStats | null>(null)

// 分组中文名映射
const groupLabelMap: Record<string, string> = {
  general: '通用配置',
  system: '系统配置',
  network: '网络配置',
  environment: '环境配置',
}

// 按后端 config_group 字段分组
const groups = computed(() => {
  const map = new Map<string, SystemConfig[]>()
  for (const c of configs.value) {
    const group = c.config_group || 'other'
    if (!map.has(group)) map.set(group, [])
    map.get(group)!.push(c)
  }
  return Array.from(map.entries()).map(([key, items]) => ({
    key,
    label: groupLabelMap[key] || key,
    configs: items,
  }))
})

const loadConfigs = async () => {
  try {
    loading.value = true
    const response = await getSystemConfigs()
    configs.value = response.data
    // 初始化表单值
    const values: Record<string, string> = {}
    for (const c of response.data) {
      values[c.config_key] = c.config_value
    }
    formValues.value = values
  } catch (error) {
    console.error('加载配置失败:', error)
    ElMessage.error('加载配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    saving.value = true
    // 只提交有变化的配置
    const changed: Record<string, string> = {}
    for (const c of configs.value) {
      if (formValues.value[c.config_key] !== c.config_value) {
        changed[c.config_key] = formValues.value[c.config_key] ?? ''
      }
    }
    if (Object.keys(changed).length === 0) {
      ElMessage.info('没有修改')
      return
    }
    await updateSystemConfigs({ configs: changed })
    ElMessage.success('保存成功')
    await loadConfigs()
  } catch (error) {
    console.error('保存配置失败:', error)
    ElMessage.error('保存配置失败')
  } finally {
    saving.value = false
  }
}

const handleReset = () => {
  const values: Record<string, string> = {}
  for (const c of configs.value) {
    values[c.config_key] = c.config_value
  }
  formValues.value = values
}

// 存储数据加载
const loadStorage = async () => {
  storageLoading.value = true
  try {
    const [backendsRes, statsRes] = await Promise.all([
      getStorageBackends(),
      getStorageStats(),
    ])
    storageBackends.value = backendsRes.data?.backends || []
    storageStatsData.value = statsRes.data || null
  } catch (error) {
    console.error('加载存储信息失败:', error)
  } finally {
    storageLoading.value = false
  }
}

// 文件大小格式化
const formatSize = (bytes: number) => {
  if (!bytes || bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(i > 0 ? 1 : 0) + ' ' + units[i]
}

// 存储后端类型标签
const backendTypeTag = (type: string) => {
  const map: Record<string, string> = { minio: '', s3: 'warning', local: 'info' }
  return (map[type] || 'info') as '' | 'success' | 'warning' | 'info' | 'danger'
}

onMounted(() => {
  loadConfigs()
  loadStorage()
})
</script>

<template>
  <div class="platform-settings">
    <div class="page-header">
      <h2 class="page-title">平台配置</h2>
    </div>

    <el-tabs v-model="activeTab">
      <!-- 系统配置标签页 -->
      <el-tab-pane label="系统配置" name="configs">
        <div class="tab-actions">
          <el-button @click="handleReset">重置</el-button>
          <el-button type="primary" :loading="saving" @click="handleSave">保存配置</el-button>
        </div>

        <div v-loading="loading">
          <div v-for="group in groups" :key="group.key" class="config-group">
            <el-card>
              <template #header>
                <span class="group-title">{{ group.label }}</span>
              </template>
              <el-form label-width="200px" label-position="right">
                <el-form-item
                  v-for="item in group.configs"
                  :key="item.config_key"
                  :label="item.description || item.config_key"
                >
                  <!-- boolean 类型 -->
                  <el-switch
                    v-if="item.config_type === 'boolean'"
                    :model-value="formValues[item.config_key] === 'true'"
                    @change="(val: boolean) => (formValues[item.config_key] = String(val))"
                  />
                  <!-- integer 类型 -->
                  <el-input-number
                    v-else-if="item.config_type === 'integer'"
                    :model-value="Number(formValues[item.config_key])"
                    :min="0"
                    :controls="true"
                    style="width: 220px"
                    @change="(val: number | undefined) => (formValues[item.config_key] = String(val ?? 0))"
                  />
                  <!-- string / json 类型 -->
                  <el-input
                    v-else
                    v-model="formValues[item.config_key]"
                    :type="item.config_type === 'json' ? 'textarea' : 'text'"
                    :rows="item.config_type === 'json' ? 4 : undefined"
                    style="max-width: 400px"
                  />
                  <div class="config-key-hint">{{ item.config_key }}</div>
                </el-form-item>
              </el-form>
            </el-card>
          </div>

          <div v-if="groups.length === 0 && !loading" class="empty-state">
            <el-empty description="暂无配置项" />
          </div>
        </div>
      </el-tab-pane>

      <!-- 存储配置标签页 -->
      <el-tab-pane label="存储配置" name="storage">
        <div class="tab-actions">
          <el-button :icon="Refresh" @click="loadStorage">刷新</el-button>
        </div>

        <div v-loading="storageLoading">
          <!-- 存储后端列表 -->
          <el-card class="config-group">
            <template #header>
              <span class="group-title">存储后端</span>
            </template>
            <el-table :data="storageBackends" stripe>
              <el-table-column prop="name" label="名称" min-width="180" />
              <el-table-column prop="type" label="类型" width="120">
                <template #default="{ row }">
                  <el-tag :type="backendTypeTag(row.type)" size="small">{{ row.type }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="默认" width="100">
                <template #default="{ row }">
                  <el-tag v-if="row.is_default" type="success" size="small">默认</el-tag>
                  <span v-else>-</span>
                </template>
              </el-table-column>
            </el-table>
            <el-empty v-if="storageBackends.length === 0 && !storageLoading" description="暂无存储后端" />
          </el-card>

          <!-- 存储统计 -->
          <el-card v-if="storageStatsData" class="config-group">
            <template #header>
              <span class="group-title">存储统计</span>
            </template>
            <el-descriptions :column="3" border>
              <el-descriptions-item label="后端名称">{{ storageStatsData.backend_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="文件数量">{{ storageStatsData.file_count }}</el-descriptions-item>
              <el-descriptions-item label="总大小">{{ formatSize(storageStatsData.total_size) }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped>
.platform-settings {
  padding: 24px;
}

.page-header {
  margin-bottom: 16px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.tab-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-bottom: 16px;
}

.config-group {
  margin-bottom: 20px;
}

.group-title {
  font-size: 16px;
  font-weight: 600;
}

.config-key-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>
