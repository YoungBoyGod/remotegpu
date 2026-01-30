<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useRoleNavigation } from '@/composables/useRoleNavigation'
import { createEnvironment } from '@/api/environment'
import { useAuthStore } from '@/stores/auth'
import type { FormInstance, FormRules } from 'element-plus'

const { navigateTo } = useRoleNavigation()
const authStore = useAuthStore()
const currentStep = ref(0)
const formRef = ref<FormInstance>()

const formData = reactive({
  name: '',
  description: '',
  image: '',
  cpu: 4,
  memory: 16,
  gpu: 0,
  storage: 50,
  command: '',
  args: '',
  env: '',
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入环境名称', trigger: 'blur' },
    { min: 3, max: 50, message: '长度在 3 到 50 个字符', trigger: 'blur' },
  ],
  image: [
    { required: true, message: '请输入镜像地址', trigger: 'blur' },
  ],
}

const parseArgs = (value: string) => {
  return value
    .split(' ')
    .map(item => item.trim())
    .filter(Boolean)
}

const parseEnv = (value: string) => {
  const result: Record<string, string> = {}
  value
    .split('\n')
    .map(line => line.trim())
    .filter(Boolean)
    .forEach((line) => {
      const [key, ...rest] = line.split('=')
      if (!key) return
      result[key.trim()] = rest.join('=').trim()
    })
  return result
}

const nextStep = async () => {
  if (currentStep.value === 0) {
    await formRef.value?.validateField(['name'])
  }
  if (currentStep.value === 1) {
    await formRef.value?.validateField(['image'])
  }
  currentStep.value++
}

const prevStep = () => {
  currentStep.value--
}

const submitForm = async () => {
  if (!authStore.user) {
    await authStore.fetchProfile()
  }

  if (!authStore.user) {
    ElMessage.error('未获取到用户信息')
    return
  }

  try {
    await createEnvironment({
      customer_id: authStore.user.id,
      name: formData.name,
      description: formData.description || undefined,
      image: formData.image,
      cpu: formData.cpu,
      memory: formData.memory * 1024,
      gpu: formData.gpu,
      storage: formData.storage ? formData.storage * 1024 : undefined,
      command: formData.command ? parseArgs(formData.command) : undefined,
      args: formData.args ? parseArgs(formData.args) : undefined,
      env: formData.env ? parseEnv(formData.env) : undefined,
    })
    ElMessage.success('环境创建成功')
    navigateTo('/environments')
  } catch (error) {
    ElMessage.error('创建失败')
  }
}

onMounted(() => {
  if (!authStore.user && authStore.isAuthenticated) {
    authStore.fetchProfile()
  }
})
</script>

<template>
  <div class="create-environment">
    <div class="page-header">
      <h1>创建开发环境</h1>
    </div>

    <el-steps :active="currentStep" align-center class="steps">
      <el-step title="基本信息" />
      <el-step title="镜像配置" />
      <el-step title="资源配置" />
      <el-step title="确认创建" />
    </el-steps>

    <div class="form-container">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px">
        <!-- 步骤1: 基本信息 -->
        <div v-show="currentStep === 0" class="step-content">
          <el-form-item label="环境名称" prop="name">
            <el-input v-model="formData.name" placeholder="请输入环境名称" />
          </el-form-item>
          <el-form-item label="描述">
            <el-input
              v-model="formData.description"
              type="textarea"
              :rows="3"
              placeholder="请输入环境描述（可选）"
            />
          </el-form-item>
        </div>

        <!-- 步骤2: 镜像配置 -->
        <div v-show="currentStep === 1" class="step-content">
          <el-form-item label="镜像地址" prop="image">
            <el-input v-model="formData.image" placeholder="例如: pytorch/pytorch:2.0-cuda11.8" />
          </el-form-item>
          <el-form-item label="启动命令">
            <el-input v-model="formData.command" placeholder="例如: /bin/bash" />
          </el-form-item>
          <el-form-item label="启动参数">
            <el-input v-model="formData.args" placeholder="例如: -c sleep 3600" />
          </el-form-item>
          <el-form-item label="环境变量">
            <el-input
              v-model="formData.env"
              type="textarea"
              :rows="4"
              placeholder="每行一个 KEY=VALUE，例如: CUDA_VISIBLE_DEVICES=0"
            />
          </el-form-item>
        </div>

        <!-- 步骤3: 资源配置 -->
        <div v-show="currentStep === 2" class="step-content">
          <el-form-item label="CPU 核心数">
            <el-slider v-model="formData.cpu" :min="1" :max="32" :marks="{ 1: '1', 16: '16', 32: '32' }" />
            <span class="value-label">{{ formData.cpu }} 核</span>
          </el-form-item>
          <el-form-item label="内存大小">
            <el-slider v-model="formData.memory" :min="2" :max="128" :marks="{ 2: '2GB', 64: '64GB', 128: '128GB' }" />
            <span class="value-label">{{ formData.memory }} GB</span>
          </el-form-item>
          <el-form-item label="GPU 数量">
            <el-slider v-model="formData.gpu" :min="0" :max="8" :marks="{ 0: '0', 4: '4', 8: '8' }" />
            <span class="value-label">{{ formData.gpu }} 个</span>
          </el-form-item>
          <el-form-item label="存储空间">
            <el-slider v-model="formData.storage" :min="10" :max="500" :marks="{ 10: '10GB', 250: '250GB', 500: '500GB' }" />
            <span class="value-label">{{ formData.storage }} GB</span>
          </el-form-item>
        </div>

        <!-- 步骤4: 确认创建 -->
        <div v-show="currentStep === 3" class="step-content">
          <el-descriptions title="配置摘要" :column="2" border>
            <el-descriptions-item label="环境名称">{{ formData.name }}</el-descriptions-item>
            <el-descriptions-item label="镜像">{{ formData.image || '-' }}</el-descriptions-item>
            <el-descriptions-item label="CPU">{{ formData.cpu }} 核</el-descriptions-item>
            <el-descriptions-item label="内存">{{ formData.memory }} GB</el-descriptions-item>
            <el-descriptions-item label="GPU">{{ formData.gpu }} 个</el-descriptions-item>
            <el-descriptions-item label="存储">{{ formData.storage }} GB</el-descriptions-item>
          </el-descriptions>
        </div>
      </el-form>

      <div class="form-actions">
        <el-button v-if="currentStep > 0" @click="prevStep">上一步</el-button>
        <el-button v-if="currentStep < 3" type="primary" @click="nextStep">下一步</el-button>
        <el-button v-if="currentStep === 3" type="primary" @click="submitForm">创建环境</el-button>
        <el-button @click="navigateTo('/environments')">取消</el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.create-environment {
  padding: 24px;
}

.page-header h1 {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 24px 0;
}

.steps {
  margin-bottom: 40px;
}

.form-container {
  max-width: 800px;
  margin: 0 auto;
  background: white;
  padding: 32px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.step-content {
  min-height: 300px;
}

.value-label {
  margin-left: 12px;
  color: #606266;
  font-weight: 600;
}

.form-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 32px;
}
</style>
