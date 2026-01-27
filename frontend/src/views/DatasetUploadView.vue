<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import { useRoleNavigation } from '@/composables/useRoleNavigation'
import type { FormInstance, UploadProps } from 'element-plus'

const router = useRouter()
const { navigateTo } = useRoleNavigation()
const formRef = ref<FormInstance>()
const uploadProgress = ref(0)
const uploading = ref(false)

const formData = reactive({
  name: '',
  description: '',
  dataType: 'image',
  visibility: 'private',
  tags: [] as string[],
})

const handleUpload: UploadProps['onChange'] = (file) => {
  uploading.value = true
  // 模拟上传进度
  const interval = setInterval(() => {
    uploadProgress.value += 10
    if (uploadProgress.value >= 100) {
      clearInterval(interval)
      uploading.value = false
      ElMessage.success('上传成功')
      setTimeout(() => navigateTo('/datasets'), 1500)
    }
  }, 500)
}

const submitForm = async () => {
  await formRef.value?.validate()
  ElMessage.success('数据集创建成功')
  navigateTo('/datasets')
}
</script>

<template>
  <div class="upload-dataset">
    <div class="page-header">
      <h1>上传数据集</h1>
    </div>

    <div class="form-container">
      <el-form ref="formRef" :model="formData" label-width="120px">
        <el-form-item label="数据集名称" required>
          <el-input v-model="formData.name" placeholder="请输入数据集名称" />
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入数据集描述"
          />
        </el-form-item>

        <el-form-item label="数据类型">
          <el-select v-model="formData.dataType" placeholder="请选择数据类型">
            <el-option label="图像" value="image" />
            <el-option label="文本" value="text" />
            <el-option label="音频" value="audio" />
            <el-option label="视频" value="video" />
          </el-select>
        </el-form-item>

        <el-form-item label="可见性">
          <el-radio-group v-model="formData.visibility">
            <el-radio value="public">公开</el-radio>
            <el-radio value="private">私有</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="上传文件">
          <el-upload
            drag
            :auto-upload="false"
            :on-change="handleUpload"
            multiple
          >
            <el-icon class="el-icon--upload"><Upload /></el-icon>
            <div class="el-upload__text">
              拖拽文件到此处或<em>点击上传</em>
            </div>
          </el-upload>
        </el-form-item>

        <el-form-item v-if="uploading">
          <el-progress :percentage="uploadProgress" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="submitForm">创建数据集</el-button>
          <el-button @click="navigateTo('/datasets')">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.upload-dataset {
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
</style>
