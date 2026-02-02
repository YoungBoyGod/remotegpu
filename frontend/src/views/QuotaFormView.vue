<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import PageHeader from '@/components/common/PageHeader.vue'
import { setQuota, updateQuota, getQuotaById } from '@/api/quota'
import type { SetQuotaRequest, UpdateQuotaRequest, QuotaLevel } from '@/api/quota/types'

const router = useRouter()
const route = useRoute()

const formRef = ref<FormInstance>()
const loading = ref(false)
const submitting = ref(false)

// 判断是编辑模式还是创建模式
const isEditMode = computed(() => !!route.params.id)
const quotaId = computed(() => Number(route.params.id))

const formData = reactive<SetQuotaRequest>({
  customer_id: 0,
  workspace_id: null,
  max_gpu: 0,
  max_cpu: 0,
  max_memory: 0,
  max_storage: 0,
  max_environments: 0,
  quota_level: 'free',
})

// 配额级别选项
const quotaLevelOptions = [
  { label: '免费版', value: 'free' },
  { label: '基础版', value: 'basic' },
  { label: '专业版', value: 'pro' },
  { label: '企业版', value: 'enterprise' },
]

// 表单验证规则
const rules: FormRules = {
  customer_id: [
    { required: true, message: '请输入客户ID', trigger: 'blur' },
    { type: 'number', min: 1, message: '客户ID必须大于0', trigger: 'blur' },
  ],
  max_gpu: [
    { required: true, message: '请输入GPU配额', trigger: 'blur' },
    { type: 'number', min: 0, message: 'GPU配额不能为负数', trigger: 'blur' },
  ],
  max_cpu: [
    { required: true, message: '请输入CPU配额', trigger: 'blur' },
    { type: 'number', min: 0, message: 'CPU配额不能为负数', trigger: 'blur' },
  ],
  max_memory: [
    { required: true, message: '请输入内存配额', trigger: 'blur' },
    { type: 'number', min: 0, message: '内存配额不能为负数', trigger: 'blur' },
  ],
  max_storage: [
    { required: true, message: '请输入存储配额', trigger: 'blur' },
    { type: 'number', min: 0, message: '存储配额不能为负数', trigger: 'blur' },
  ],
  max_environments: [
    { required: true, message: '请输入环境数量配额', trigger: 'blur' },
    { type: 'number', min: 0, message: '环境数量配额不能为负数', trigger: 'blur' },
  ],
}

// 加载配额详情(编辑模式)
const loadQuota = async () => {
  if (!isEditMode.value) return

  loading.value = true
  try {
    const response = await getQuotaById(quotaId.value)
    const quota = response.data
    formData.customer_id = quota.customer_id
    formData.workspace_id = quota.workspace_id
    formData.max_gpu = quota.max_gpu
    formData.max_cpu = quota.max_cpu
    formData.max_memory = quota.max_memory
    formData.max_storage = quota.max_storage
    formData.max_environments = quota.max_environments
    formData.quota_level = quota.quota_level
  } catch (error: any) {
    ElMessage.error(error.message || '加载配额信息失败')
    router.back()
  } finally {
    loading.value = false
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    if (isEditMode.value) {
      // 编辑模式
      const updateData: UpdateQuotaRequest = {
        max_gpu: formData.max_gpu,
        max_cpu: formData.max_cpu,
        max_memory: formData.max_memory,
        max_storage: formData.max_storage,
        max_environments: formData.max_environments,
        quota_level: formData.quota_level,
      }
      await updateQuota(quotaId.value, updateData)
      ElMessage.success('更新成功')
    } else {
      // 创建模式
      await setQuota(formData)
      ElMessage.success('创建成功')
    }

    router.push('/admin/quotas')
  } catch (error: any) {
    if (error.message) {
      ElMessage.error(error.message)
    }
  } finally {
    submitting.value = false
  }
}

// 取消
const handleCancel = () => {
  router.back()
}

onMounted(() => {
  loadQuota()
})
</script>

<template>
  <div class="quota-form">
    <PageHeader :title="isEditMode ? '编辑配额' : '设置配额'" />

    <el-card v-loading="loading">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="140px"
        style="max-width: 600px"
      >
        <el-form-item label="客户ID" prop="customer_id">
          <el-input-number
            v-model="formData.customer_id"
            :min="1"
            :disabled="isEditMode"
            placeholder="请输入客户ID"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="工作空间ID" prop="workspace_id">
          <el-input-number
            v-model="formData.workspace_id"
            :min="1"
            :disabled="isEditMode"
            placeholder="留空表示用户级配额"
            style="width: 100%"
            clearable
          />
          <div class="form-tip">留空表示设置用户级配额,填写则为工作空间级配额</div>
        </el-form-item>

        <el-form-item label="配额级别" prop="quota_level">
          <el-select v-model="formData.quota_level" style="width: 100%">
            <el-option
              v-for="option in quotaLevelOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
        </el-form-item>

        <el-divider content-position="left">资源配额</el-divider>

        <el-form-item label="GPU数量" prop="max_gpu">
          <el-input-number
            v-model="formData.max_gpu"
            :min="0"
            placeholder="请输入GPU配额"
            style="width: 100%"
          />
          <div class="form-tip">单位: 个</div>
        </el-form-item>

        <el-form-item label="CPU核心数" prop="max_cpu">
          <el-input-number
            v-model="formData.max_cpu"
            :min="0"
            placeholder="请输入CPU配额"
            style="width: 100%"
          />
          <div class="form-tip">单位: 核</div>
        </el-form-item>

        <el-form-item label="内存" prop="max_memory">
          <el-input-number
            v-model="formData.max_memory"
            :min="0"
            :step="1024"
            placeholder="请输入内存配额"
            style="width: 100%"
          />
          <div class="form-tip">单位: MB (1024MB = 1GB)</div>
        </el-form-item>

        <el-form-item label="存储" prop="max_storage">
          <el-input-number
            v-model="formData.max_storage"
            :min="0"
            :step="100"
            placeholder="请输入存储配额"
            style="width: 100%"
          />
          <div class="form-tip">单位: GB</div>
        </el-form-item>

        <el-form-item label="环境数量" prop="max_environments">
          <el-input-number
            v-model="formData.max_environments"
            :min="0"
            placeholder="请输入环境数量配额"
            style="width: 100%"
          />
          <div class="form-tip">单位: 个</div>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            :loading="submitting"
            @click="handleSubmit"
          >
            {{ isEditMode ? '保存' : '创建' }}
          </el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.quota-form {
  padding: 24px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>
