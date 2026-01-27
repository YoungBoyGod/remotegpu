<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useRoleNavigation } from '@/composables/useRoleNavigation'

const router = useRouter()
const { navigateTo } = useRoleNavigation()

const formData = reactive({
  name: '',
  baseImage: '',
  dockerfile: `FROM pytorch/pytorch:2.0-cuda11.8

RUN pip install numpy pandas scikit-learn

WORKDIR /workspace`,
})

const buildLogs = ref('')
const building = ref(false)

const startBuild = async () => {
  building.value = true
  buildLogs.value = '开始构建镜像...\n'

  // 模拟构建过程
  const logs = [
    'Step 1/3 : FROM pytorch/pytorch:2.0-cuda11.8',
    'Step 2/3 : RUN pip install numpy pandas scikit-learn',
    'Step 3/3 : WORKDIR /workspace',
    '构建成功！',
  ]

  for (const log of logs) {
    await new Promise(resolve => setTimeout(resolve, 1000))
    buildLogs.value += log + '\n'
  }

  building.value = false
  ElMessage.success('镜像构建成功')
  setTimeout(() => navigateTo('/images'), 1500)
}
</script>

<template>
  <div class="build-image">
    <div class="page-header">
      <h1>构建自定义镜像</h1>
    </div>

    <div class="form-container">
      <el-form :model="formData" label-width="120px">
        <el-form-item label="镜像名称" required>
          <el-input v-model="formData.name" placeholder="请输入镜像名称" />
        </el-form-item>

        <el-form-item label="基础镜像">
          <el-select v-model="formData.baseImage" placeholder="选择基础镜像">
            <el-option label="PyTorch 2.0" value="pytorch/pytorch:2.0-cuda11.8" />
            <el-option label="TensorFlow 2.13" value="tensorflow/tensorflow:2.13-gpu" />
            <el-option label="CUDA 12.0" value="nvidia/cuda:12.0-base" />
          </el-select>
        </el-form-item>

        <el-form-item label="Dockerfile">
          <el-input
            v-model="formData.dockerfile"
            type="textarea"
            :rows="12"
            placeholder="输入 Dockerfile 内容"
          />
        </el-form-item>

        <el-form-item v-if="building">
          <el-card>
            <template #header>构建日志</template>
            <pre class="build-logs">{{ buildLogs }}</pre>
          </el-card>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="building" @click="startBuild">
            开始构建
          </el-button>
          <el-button @click="navigateTo('/images')">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.build-image {
  padding: 24px;
}

.page-header h1 {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 24px 0;
}

.form-container {
  max-width: 900px;
  background: white;
  padding: 32px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.build-logs {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 16px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  max-height: 300px;
  overflow-y: auto;
}
</style>
