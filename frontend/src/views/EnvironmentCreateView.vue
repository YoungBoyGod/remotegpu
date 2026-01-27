<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useRoleNavigation } from '@/composables/useRoleNavigation'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const { navigateTo } = useRoleNavigation()
const currentStep = ref(0)
const formRef = ref<FormInstance>()

const formData = reactive({
  name: '',
  description: '',
  osType: 'linux',
  imageId: '',
  cpu: 4,
  memory: 16,
  gpuCount: 1,
  gpuModel: '',
  storage: 50,
  enableSSH: true,
  enableJupyter: false,
  datasets: [] as string[],
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入环境名称', trigger: 'blur' },
    { min: 3, max: 50, message: '长度在 3 到 50 个字符', trigger: 'blur' },
  ],
}

const images = [
  { id: 'img-1', name: 'PyTorch 2.0 + CUDA 11.8', type: 'pytorch' },
  { id: 'img-2', name: 'TensorFlow 2.13 + CUDA 11.8', type: 'tensorflow' },
  { id: 'img-3', name: 'CUDA 12.0 Base', type: 'cuda' },
]

const nextStep = async () => {
  if (currentStep.value === 0) {
    await formRef.value?.validate()
  }
  currentStep.value++
}

const prevStep = () => {
  currentStep.value--
}

const submitForm = async () => {
  try {
    ElMessage.success('环境创建中，请稍候...')
    setTimeout(() => {
      navigateTo('/environments')
    }, 1500)
  } catch (error) {
    ElMessage.error('创建失败')
  }
}
</script>

<template>
  <div class="create-environment">
    <div class="page-header">
      <h1>创建开发环境</h1>
    </div>

    <el-steps :active="currentStep" align-center class="steps">
      <el-step title="基本信息" />
      <el-step title="选择镜像" />
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
          <el-form-item label="操作系统">
            <el-radio-group v-model="formData.osType">
              <el-radio value="linux">Linux</el-radio>
              <el-radio value="windows">Windows</el-radio>
            </el-radio-group>
          </el-form-item>
        </div>

        <!-- 步骤2: 选择镜像 -->
        <div v-show="currentStep === 1" class="step-content">
          <el-form-item label="选择镜像">
            <el-radio-group v-model="formData.imageId">
              <div class="image-list">
                <div
                  v-for="img in images"
                  :key="img.id"
                  class="image-card"
                  :class="{ selected: formData.imageId === img.id }"
                  @click="formData.imageId = img.id"
                >
                  <el-radio :value="img.id">{{ img.name }}</el-radio>
                </div>
              </div>
            </el-radio-group>
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
            <el-slider v-model="formData.gpuCount" :min="0" :max="8" :marks="{ 0: '0', 4: '4', 8: '8' }" />
            <span class="value-label">{{ formData.gpuCount }} 个</span>
          </el-form-item>
          <el-form-item label="存储空间">
            <el-slider v-model="formData.storage" :min="10" :max="500" :marks="{ 10: '10GB', 250: '250GB', 500: '500GB' }" />
            <span class="value-label">{{ formData.storage }} GB</span>
          </el-form-item>
          <el-form-item label="网络配置">
            <el-checkbox v-model="formData.enableSSH">启用 SSH 访问</el-checkbox>
            <el-checkbox v-model="formData.enableJupyter">启用 JupyterLab</el-checkbox>
          </el-form-item>
        </div>

        <!-- 步骤4: 确认创建 -->
        <div v-show="currentStep === 3" class="step-content">
          <el-descriptions title="配置摘要" :column="2" border>
            <el-descriptions-item label="环境名称">{{ formData.name }}</el-descriptions-item>
            <el-descriptions-item label="操作系统">{{ formData.osType === 'linux' ? 'Linux' : 'Windows' }}</el-descriptions-item>
            <el-descriptions-item label="CPU">{{ formData.cpu }} 核</el-descriptions-item>
            <el-descriptions-item label="内存">{{ formData.memory }} GB</el-descriptions-item>
            <el-descriptions-item label="GPU">{{ formData.gpuCount }} 个</el-descriptions-item>
            <el-descriptions-item label="存储">{{ formData.storage }} GB</el-descriptions-item>
          </el-descriptions>
          <div class="cost-estimate">
            <h3>预估费用</h3>
            <p class="cost">¥ 12.50 / 小时</p>
          </div>
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

.image-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
}

.image-card {
  padding: 16px;
  border: 2px solid #dcdfe6;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.image-card:hover {
  border-color: #409eff;
}

.image-card.selected {
  border-color: #409eff;
  background: #ecf5ff;
}

.value-label {
  margin-left: 12px;
  color: #606266;
  font-weight: 600;
}

.cost-estimate {
  margin-top: 24px;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;
  text-align: center;
}

.cost-estimate h3 {
  margin: 0 0 12px 0;
  font-size: 16px;
  color: #606266;
}

.cost {
  font-size: 32px;
  font-weight: 600;
  color: #409eff;
  margin: 0;
}

.form-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 32px;
}
</style>
