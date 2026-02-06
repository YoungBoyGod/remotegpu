<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getSystemConfigs, updateSystemConfigs } from '@/api/admin'
import type { SystemConfig } from '@/types/systemConfig'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const saving = ref(false)
const configs = ref<SystemConfig[]>([])
const formValues = ref<Record<string, string>>({})

// 分组规则
const groupDefs = [
  { key: 'basic', label: '基本信息', prefixes: ['system_'] },
  { key: 'port', label: '端口配置', prefixes: ['ssh_', 'rdp_', 'jupyter_'] },
  { key: 'limit', label: '资源限制', prefixes: ['max_', 'default_'] }
]

// 将配置项按前缀分组
const groups = computed(() => {
  const result: { key: string; label: string; configs: SystemConfig[] }[] = []
  const matched = new Set<string>()

  for (const def of groupDefs) {
    const items = configs.value.filter((c) => {
      return def.prefixes.some((p) => c.config_key.startsWith(p))
    })
    items.forEach((c) => matched.add(c.config_key))
    if (items.length > 0) {
      result.push({ key: def.key, label: def.label, configs: items })
    }
  }

  // 未匹配的归入"其他"
  const others = configs.value.filter((c) => !matched.has(c.config_key))
  if (others.length > 0) {
    result.push({ key: 'other', label: '其他配置', configs: others })
  }

  return result
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
        changed[c.config_key] = formValues.value[c.config_key]
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

onMounted(() => {
  loadConfigs()
})
</script>

<template>
  <div class="platform-settings" v-loading="loading">
    <div class="page-header">
      <h2 class="page-title">平台配置</h2>
      <div class="header-actions">
        <el-button @click="handleReset">重置</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存配置</el-button>
      </div>
    </div>

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
</template>

<style scoped>
.platform-settings {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
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
