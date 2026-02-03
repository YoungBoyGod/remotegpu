<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { addMachine } from '@/api/admin'
import type { Machine } from '@/types/machine'
import { ElMessage } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const formData = ref<Partial<Machine>>({
  name: '',
  region: '',
  status: 'online',
  gpuModel: '',
  gpuMemory: 0,
  gpuCount: 1,
  cpu: '',
  memory: 0,
  disk: 0,
  cudaVersion: '',
  gpuDriver: '',
  loginInfo: {
    sshHost: '',
    sshPort: 22,
    username: '',
    password: ''
  }
})

const rules = {
  name: [{ required: true, message: '请输入机器名称', trigger: 'blur' }],
  region: [{ required: true, message: '请输入区域', trigger: 'blur' }],
  gpuModel: [{ required: true, message: '请输入GPU型号', trigger: 'blur' }],
  gpuMemory: [{ required: true, message: '请输入GPU显存', trigger: 'blur' }],
  gpuCount: [{ required: true, message: '请输入GPU数量', trigger: 'blur' }],
  cpu: [{ required: true, message: '请输入CPU信息', trigger: 'blur' }],
  memory: [{ required: true, message: '请输入内存大小', trigger: 'blur' }],
  disk: [{ required: true, message: '请输入磁盘大小', trigger: 'blur' }],
  'loginInfo.sshHost': [{ required: true, message: '请输入IP地址', trigger: 'blur' }],
  'loginInfo.username': [{ required: true, message: '请输入用户名', trigger: 'blur' }]
}

const formRef = ref()

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    await addMachine(formData.value)
    ElMessage.success('添加机器成功')
    router.push('/admin/machines/list')
  } catch (error: any) {
    if (error !== false) {
      console.error('添加机器失败:', error)
    }
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  router.back()
}
</script>

<template>
  <div class="machine-add">
    <div class="page-header">
      <h2 class="page-title">添加机器</h2>
    </div>

    <el-card>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="120px"
      >
        <el-divider content-position="left">基本信息</el-divider>

        <el-form-item label="机器名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入机器名称" />
        </el-form-item>

        <el-form-item label="区域" prop="region">
          <el-input v-model="formData.region" placeholder="请输入区域,如:北京/上海" />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio label="online">在线</el-radio>
            <el-radio label="offline">离线</el-radio>
            <el-radio label="maintenance">维护中</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-divider content-position="left">硬件配置</el-divider>

        <el-form-item label="GPU型号" prop="gpuModel">
          <el-input v-model="formData.gpuModel" placeholder="如: NVIDIA A100" />
        </el-form-item>

        <el-form-item label="GPU显存(GB)" prop="gpuMemory">
          <el-input-number v-model="formData.gpuMemory" :min="1" :max="1000" />
        </el-form-item>

        <el-form-item label="GPU数量" prop="gpuCount">
          <el-input-number v-model="formData.gpuCount" :min="1" :max="16" />
        </el-form-item>

        <el-form-item label="CPU" prop="cpu">
          <el-input v-model="formData.cpu" placeholder="如: Intel Xeon Gold 6248R" />
        </el-form-item>

        <el-form-item label="内存(GB)" prop="memory">
          <el-input-number v-model="formData.memory" :min="1" :max="10000" />
        </el-form-item>

        <el-form-item label="磁盘(GB)" prop="disk">
          <el-input-number v-model="formData.disk" :min="1" :max="100000" />
        </el-form-item>

        <el-form-item label="CUDA版本" prop="cudaVersion">
          <el-input v-model="formData.cudaVersion" placeholder="如: 12.1" />
        </el-form-item>

        <el-form-item label="GPU驱动" prop="gpuDriver">
          <el-input v-model="formData.gpuDriver" placeholder="如: 535.104.05" />
        </el-form-item>

        <el-divider content-position="left">登录信息</el-divider>

        <el-form-item label="IP地址" prop="loginInfo.sshHost">
          <el-input v-model="formData.loginInfo!.sshHost" placeholder="请输入IP地址" />
        </el-form-item>

        <el-form-item label="SSH端口" prop="loginInfo.sshPort">
          <el-input-number v-model="formData.loginInfo!.sshPort" :min="1" :max="65535" />
        </el-form-item>

        <el-form-item label="用户名" prop="loginInfo.username">
          <el-input v-model="formData.loginInfo!.username" placeholder="请输入SSH用户名" />
        </el-form-item>

        <el-form-item label="密码" prop="loginInfo.password">
          <el-input
            v-model="formData.loginInfo!.password"
            type="password"
            placeholder="请输入SSH密码(可选)"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            提交
          </el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.machine-add {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}
</style>
