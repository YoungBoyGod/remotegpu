<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Edit, Delete } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import {
  getAlertList,
  acknowledgeAlert,
  getAlertRuleList,
  createAlertRule,
  updateAlertRule,
  deleteAlertRule,
  getMachineList,
} from '@/api/admin'
import type { AlertRule, AlertRuleForm } from '@/api/admin'
import type { Machine } from '@/types/machine'

interface Alert {
  id: number
  rule_id: number
  host_id: string
  value: number
  message: string
  triggered_at: string
  acknowledged: boolean
  rule: {
    name: string
    severity: string
    metric_type: string
  }
}

// Tab 切换
const activeTab = ref('alerts')

const loading = ref(false)
const alerts = ref<Alert[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const selectedSeverity = ref<string>('')
const selectedAcknowledged = ref<string>('false') // Default show active

// 加载告警列表
const loadAlerts = async () => {
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      pageSize: pageSize.value,
    }
    if (selectedSeverity.value) params.severity = selectedSeverity.value
    if (selectedAcknowledged.value !== '') {
      params.acknowledged = selectedAcknowledged.value === 'true'
    }

    const res = await getAlertList(params)
    alerts.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error(error)
    ElMessage.error('加载告警列表失败')
  } finally {
    loading.value = false
  }
}

// 确认告警
const handleAcknowledge = async (alert: Alert) => {
  try {
    await ElMessageBox.confirm('确认已知晓此告警？', '确认告警', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await acknowledgeAlert(alert.id)
    ElMessage.success('告警已确认')
    loadAlerts()
  } catch {
    // 用户取消
  }
}

const handlePageChange = (val: number) => {
  page.value = val
  loadAlerts()
}

// 获取级别类型
const getSeverityType = (level: string) => {
  const map: Record<string, string> = {
    critical: 'danger',
    warning: 'warning',
    info: 'info'
  }
  return map[level] || 'info'
}

const severityLabel = (level: string) => {
  const map: Record<string, string> = { critical: '严重', warning: '警告', info: '信息' }
  return map[level] || level
}

const metricTypeLabel = (type: string) => {
  const map: Record<string, string> = {
    gpu_usage: 'GPU 使用率',
    gpu_temperature: 'GPU 温度',
    cpu_usage: 'CPU 使用率',
    memory_usage: '内存使用率',
    disk_usage: '磁盘使用率',
  }
  return map[type] || type
}

const conditionLabel = (cond: string) => {
  const map: Record<string, string> = { '>': '大于', '>=': '大于等于', '<': '小于', '<=': '小于等于', '==': '等于' }
  return map[cond] || cond
}

// ==================== 告警规则管理 ====================

const rules = ref<AlertRule[]>([])
const rulesLoading = ref(false)
const rulesTotal = ref(0)
const rulesPage = ref(1)
const rulesPageSize = ref(10)

const loadRules = async () => {
  rulesLoading.value = true
  try {
    const res = await getAlertRuleList({ page: rulesPage.value, pageSize: rulesPageSize.value })
    rules.value = res.data.list || []
    rulesTotal.value = res.data.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    rulesLoading.value = false
  }
}

const handleRulesPageChange = (val: number) => {
  rulesPage.value = val
  loadRules()
}

// 规则表单对话框
const ruleDialogVisible = ref(false)
const ruleDialogTitle = ref('添加告警规则')
const ruleSubmitLoading = ref(false)
const editingRuleId = ref<number | null>(null)
const ruleFormRef = ref()
const ruleForm = ref<AlertRuleForm>({
  name: '',
  metric_type: 'gpu_usage',
  condition: '>',
  threshold: 90,
  severity: 'warning',
  duration: 300,
  enabled: true,
  description: '',
})

const ruleFormRules = {
  name: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
  metric_type: [{ required: true, message: '请选择监控指标', trigger: 'change' }],
  threshold: [{ required: true, message: '请输入阈值', trigger: 'blur' }],
  severity: [{ required: true, message: '请选择告警级别', trigger: 'change' }],
}

const openAddRuleDialog = () => {
  editingRuleId.value = null
  ruleDialogTitle.value = '添加告警规则'
  ruleForm.value = {
    name: '',
    metric_type: 'gpu_usage',
    condition: '>',
    threshold: 90,
    severity: 'warning',
    duration: 300,
    enabled: true,
    description: '',
  }
  ruleDialogVisible.value = true
  ruleFormRef.value?.resetFields()
}

const openEditRuleDialog = (rule: AlertRule) => {
  editingRuleId.value = rule.id
  ruleDialogTitle.value = '编辑告警规则'
  ruleForm.value = {
    name: rule.name,
    metric_type: rule.metric_type,
    condition: rule.condition,
    threshold: rule.threshold,
    severity: rule.severity,
    duration: rule.duration_seconds,
    enabled: rule.enabled,
    description: rule.description || '',
  }
  ruleDialogVisible.value = true
}

const handleRuleSubmit = async () => {
  if (!ruleFormRef.value) return
  await ruleFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return
    ruleSubmitLoading.value = true
    try {
      if (editingRuleId.value) {
        await updateAlertRule(editingRuleId.value, ruleForm.value)
        ElMessage.success('规则已更新')
      } else {
        await createAlertRule(ruleForm.value)
        ElMessage.success('规则已创建')
      }
      ruleDialogVisible.value = false
      loadRules()
    } catch (error) {
      console.error(error)
    } finally {
      ruleSubmitLoading.value = false
    }
  })
}

const handleDeleteRule = async (rule: AlertRule) => {
  try {
    await ElMessageBox.confirm(`确认删除规则「${rule.name}」？`, '确认删除', { type: 'warning' })
    await deleteAlertRule(rule.id)
    ElMessage.success('规则已删除')
    loadRules()
  } catch {
    // 取消
  }
}

const handleToggleRule = async (rule: AlertRule) => {
  try {
    await updateAlertRule(rule.id, { enabled: !rule.enabled })
    ElMessage.success(rule.enabled ? '规则已禁用' : '规则已启用')
    loadRules()
  } catch (error) {
    console.error(error)
  }
}

// ==================== 离线机器 ====================

const offlineMachines = ref<Machine[]>([])
const offlineLoading = ref(false)

const loadOfflineMachines = async () => {
  offlineLoading.value = true
  try {
    const res = await getMachineList({
      page: 1,
      pageSize: 100,
      status: 'offline',
    })
    offlineMachines.value = res.data.list || []
  } catch (error) {
    console.error(error)
  } finally {
    offlineLoading.value = false
  }
}

const getOfflineDuration = (lastHeartbeat?: string | null) => {
  if (!lastHeartbeat) return '未知'
  const last = new Date(lastHeartbeat).getTime()
  if (Number.isNaN(last)) return '未知'
  const diff = Date.now() - last
  const minutes = Math.floor(diff / 60000)
  if (minutes < 60) return `${minutes} 分钟`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours} 小时 ${minutes % 60} 分钟`
  const days = Math.floor(hours / 24)
  return `${days} 天 ${hours % 24} 小时`
}

const formatDateTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN')
}

// Tab 切换时加载对应数据
watch(activeTab, (tab) => {
  if (tab === 'alerts') loadAlerts()
  else if (tab === 'rules') loadRules()
  else if (tab === 'offline') loadOfflineMachines()
})

onMounted(() => {
  loadAlerts()
})
</script>

<template>
  <div class="alert-center">
    <PageHeader title="告警中心" />

    <el-tabs v-model="activeTab">
      <!-- ==================== 告警记录 Tab ==================== -->
      <el-tab-pane label="告警记录" name="alerts">
        <!-- 筛选区域 -->
        <el-card class="filter-card">
          <div class="filter-container">
            <el-select v-model="selectedSeverity" placeholder="告警级别" clearable style="width: 150px" @change="loadAlerts">
              <el-option label="严重" value="critical" />
              <el-option label="警告" value="warning" />
              <el-option label="信息" value="info" />
            </el-select>

            <el-select v-model="selectedAcknowledged" placeholder="告警状态" style="width: 150px" @change="loadAlerts">
              <el-option label="未确认" value="false" />
              <el-option label="已确认" value="true" />
              <el-option label="全部" value="" />
            </el-select>

            <el-button type="primary" :icon="Search" @click="loadAlerts">刷新</el-button>
          </div>
        </el-card>

        <!-- 告警列表 -->
        <el-card class="alert-list-card">
          <div v-loading="loading" class="alert-list">
            <div
              v-for="alert in alerts"
              :key="alert.id"
              class="alert-item"
              :class="`alert-${alert.rule?.severity || 'info'}`"
            >
              <div class="alert-header">
                <div class="alert-title-section">
                  <el-tag :type="getSeverityType(alert.rule?.severity)" size="small">
                    {{ severityLabel(alert.rule?.severity) }}
                  </el-tag>
                  <h4 class="alert-title">{{ alert.rule?.name || '未知规则' }}</h4>
                </div>
                <el-tag v-if="alert.acknowledged" type="success" size="small">已确认</el-tag>
                <el-tag v-else type="danger" size="small">未确认</el-tag>
              </div>

              <p class="alert-description">{{ alert.message }}（当前值: {{ alert.value }}）</p>

              <div class="alert-meta">
                <span class="meta-item">主机: {{ alert.host_id }}</span>
                <span class="meta-item">触发时间: {{ new Date(alert.triggered_at).toLocaleString('zh-CN') }}</span>
              </div>

              <div class="alert-actions" v-if="!alert.acknowledged">
                <el-button size="small" type="warning" @click="handleAcknowledge(alert)">
                  确认告警
                </el-button>
              </div>
            </div>

            <el-empty v-if="!loading && alerts.length === 0" description="暂无告警信息" />
          </div>

          <div class="pagination-container">
            <el-pagination
              v-model:current-page="page"
              v-model:page-size="pageSize"
              :total="total"
              layout="total, prev, pager, next"
              @current-change="handlePageChange"
            />
          </div>
        </el-card>
      </el-tab-pane>

      <!-- ==================== 告警规则 Tab ==================== -->
      <el-tab-pane label="告警规则" name="rules">
        <div class="rules-toolbar">
          <el-button type="primary" :icon="Plus" @click="openAddRuleDialog">添加规则</el-button>
        </div>

        <el-card>
          <el-table :data="rules" v-loading="rulesLoading" stripe>
            <template #empty>
              <el-empty description="暂无告警规则" />
            </template>
            <el-table-column prop="name" label="规则名称" min-width="160" />
            <el-table-column label="监控指标" width="140">
              <template #default="{ row }">
                {{ metricTypeLabel(row.metric_type) }}
              </template>
            </el-table-column>
            <el-table-column label="触发条件" width="140">
              <template #default="{ row }">
                {{ conditionLabel(row.condition) }} {{ row.threshold }}
              </template>
            </el-table-column>
            <el-table-column label="持续时间" width="100">
              <template #default="{ row }">
                {{ row.duration_seconds }}s
              </template>
            </el-table-column>
            <el-table-column label="级别" width="90">
              <template #default="{ row }">
                <el-tag :type="getSeverityType(row.severity)" size="small">
                  {{ severityLabel(row.severity) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
                  {{ row.enabled ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="handleToggleRule(row)">
                  {{ row.enabled ? '禁用' : '启用' }}
                </el-button>
                <el-button link type="primary" size="small" :icon="Edit" @click="openEditRuleDialog(row)">
                  编辑
                </el-button>
                <el-button link type="danger" size="small" :icon="Delete" @click="handleDeleteRule(row)">
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination-container" v-if="rulesTotal > 0">
            <el-pagination
              v-model:current-page="rulesPage"
              v-model:page-size="rulesPageSize"
              :total="rulesTotal"
              layout="total, prev, pager, next"
              @current-change="handleRulesPageChange"
            />
          </div>
        </el-card>
      </el-tab-pane>

      <!-- ==================== 离线机器 Tab ==================== -->
      <el-tab-pane label="离线机器" name="offline">
        <el-card>
          <div class="offline-toolbar">
            <el-button type="primary" :icon="Search" @click="loadOfflineMachines">刷新</el-button>
            <span class="offline-count" v-if="offlineMachines.length > 0">
              共 {{ offlineMachines.length }} 台离线
            </span>
          </div>

          <div v-loading="offlineLoading" class="offline-list">
            <div
              v-for="machine in offlineMachines"
              :key="machine.id"
              class="offline-item"
            >
              <div class="offline-header">
                <div class="offline-title-section">
                  <el-tag type="danger" size="small">离线</el-tag>
                  <h4 class="offline-name">{{ machine.name || machine.hostname || machine.id }}</h4>
                </div>
                <el-tag type="warning" size="small">
                  已离线 {{ getOfflineDuration(machine.last_heartbeat) }}
                </el-tag>
              </div>

              <div class="offline-meta">
                <span class="meta-item">IP: {{ machine.ip_address || '-' }}</span>
                <span class="meta-item">区域: {{ machine.region || '-' }}</span>
                <span class="meta-item">最后心跳: {{ formatDateTime(machine.last_heartbeat) }}</span>
              </div>

              <div class="offline-reason">
                <span class="reason-label">可能原因：</span>
                <span v-if="!machine.last_heartbeat">从未上报心跳，Agent 可能未部署或未启动</span>
                <span v-else>心跳超时，机器可能关机、网络中断或 Agent 进程异常</span>
              </div>
            </div>

            <el-empty v-if="!offlineLoading && offlineMachines.length === 0" description="当前没有离线机器" />
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 规则编辑对话框 -->
    <el-dialog v-model="ruleDialogVisible" :title="ruleDialogTitle" width="560px" :close-on-click-modal="false">
      <el-form ref="ruleFormRef" :model="ruleForm" :rules="ruleFormRules" label-width="100px">
        <el-form-item label="规则名称" prop="name">
          <el-input v-model="ruleForm.name" placeholder="例如: GPU 使用率过高" />
        </el-form-item>
        <el-form-item label="监控指标" prop="metric_type">
          <el-select v-model="ruleForm.metric_type" style="width: 100%">
            <el-option label="GPU 使用率" value="gpu_usage" />
            <el-option label="GPU 温度" value="gpu_temperature" />
            <el-option label="CPU 使用率" value="cpu_usage" />
            <el-option label="内存使用率" value="memory_usage" />
            <el-option label="磁盘使用率" value="disk_usage" />
          </el-select>
        </el-form-item>
        <el-form-item label="触发条件">
          <div style="display: flex; gap: 12px; width: 100%">
            <el-select v-model="ruleForm.condition" style="width: 140px">
              <el-option label="大于" value=">" />
              <el-option label="大于等于" value=">=" />
              <el-option label="小于" value="<" />
              <el-option label="小于等于" value="<=" />
            </el-select>
            <el-input-number v-model="ruleForm.threshold" :min="0" :max="100" style="flex: 1" />
          </div>
        </el-form-item>
        <el-form-item label="持续时间">
          <el-input-number v-model="ruleForm.duration" :min="0" :step="60" />
          <span style="margin-left: 8px; color: #909399">秒</span>
        </el-form-item>
        <el-form-item label="告警级别" prop="severity">
          <el-radio-group v-model="ruleForm.severity">
            <el-radio value="critical">严重</el-radio>
            <el-radio value="warning">警告</el-radio>
            <el-radio value="info">信息</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="ruleForm.enabled" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="ruleForm.description" type="textarea" :rows="2" placeholder="可选" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="ruleDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="ruleSubmitLoading" @click="handleRuleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.alert-center {
  padding: 24px;
  background: #f5f7fa;
  min-height: 100%;
}

.filter-card {
  margin-bottom: 24px;
}

.filter-container {
  display: flex;
  gap: 12px;
}

.alert-list {
  min-height: 400px;
}

.alert-item {
  padding: 20px;
  margin-bottom: 16px;
  border-radius: 8px;
  border-left: 4px solid #dcdfe6;
  background: #f8f9fa;
  transition: all 0.3s;
}

.alert-critical { border-left-color: #f56c6c; background: #fef0f0; }
.alert-warning { border-left-color: #e6a23c; background: #fdf6ec; }
.alert-info { border-left-color: #409eff; background: #ecf5ff; }

.alert-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.alert-title-section {
  display: flex;
  align-items: center;
  gap: 8px;
}

.alert-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.alert-description {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #606266;
}

.alert-meta {
  display: flex;
  gap: 16px;
  color: #909399;
  font-size: 13px;
  margin-bottom: 12px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.rules-toolbar {
  margin-bottom: 16px;
  display: flex;
  justify-content: flex-end;
}

.offline-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.offline-count {
  font-size: 14px;
  color: #f56c6c;
  font-weight: 600;
}

.offline-list {
  min-height: 200px;
}

.offline-item {
  padding: 16px 20px;
  margin-bottom: 12px;
  border-radius: 8px;
  border-left: 4px solid #f56c6c;
  background: #fef0f0;
}

.offline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.offline-title-section {
  display: flex;
  align-items: center;
  gap: 8px;
}

.offline-name {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.offline-meta {
  display: flex;
  gap: 16px;
  color: #909399;
  font-size: 13px;
  margin-bottom: 8px;
}

.offline-reason {
  font-size: 13px;
  color: #606266;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 4px;
}

.offline-reason .reason-label {
  color: #e6a23c;
  font-weight: 600;
}
</style>
