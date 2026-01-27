<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()

const formData = reactive({
  name: '',
  description: '',
  imageId: '',
  command: '',
  workDir: '/workspace',
  gpuCount: 1,
  cpu: 4,
  memory: 16,
  datasetIds: [] as string[],
})

const submitForm = async () => {
  ElMessage.success('训练任务创建成功')
  setTimeout(() => router.push('/training'), 1500)
}
</script>

<template>
  <div class="create-training">
    <div class="page-header">
      <h1>创建训练任务</h1>
    </div>

    <div class="form-container">
      <el-form :model="formData" label-width="120px">
        <el-form-item label="任务名称" required>
          <el-input v-model="formData.name" placeholder="请输入任务名称" />
        </el-form-item>

        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" :rows="3" />
        </el-form-item>

        <el-form-item label="镜像" required>
          <el-select v-model="formData.imageId" placeholder="选择镜像">
            <el-option label="PyTorch 2.0" value="img-1" />
            <el-option label="TensorFlow 2.13" value="img-2" />
          </el-select>
        </el-form-item>

        <el-form-item label="启动命令" required>
          <el-input v-model="formData.command" placeholder="python train.py" />
        </el-form-item>

        <el-form-item label="工作目录">
          <el-input v-model="formData.workDir" />
        </el-form-item>

        <el-form-item label="GPU 数量">
          <el-slider v-model="formData.gpuCount" :min="1" :max="8" />
          <span class="value-label">{{ formData.gpuCount }} 个</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="submitForm">创建任务</el-button>
          <el-button @click="router.push('/training')">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.create-training {
  padding: 24px;
}

.page-header h1 {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 24px 0;
}

.form-container {
  max-width: 800px;
  background: white;
  padding: 32px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.value-label {
  margin-left: 12px;
  color: #606266;
  font-weight: 600;
}
</style>
