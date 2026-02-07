<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { addMachine } from '@/api/admin'
import type { CreateMachinePayload } from '@/api/admin'
import { ElMessage } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const formData = ref<CreateMachinePayload>({
  name: '',
  hostname: '',
  region: '',
  ip_address: '',
  public_ip: '',
  ssh_host: '',
  ssh_port: 22,
  ssh_username: '',
  ssh_password: '',
  ssh_key: ''
})

const validateConnectionAddress = (_: any, _value: string, callback: (error?: Error) => void) => {
  if (!formData.value.ip_address && !formData.value.hostname) {
    callback(new Error('请输入IP地址或主机名'))
    return
  }
  callback()
}

const rules = {
  name: [{ required: true, message: '请输入机器名称', trigger: 'blur' }],
  region: [{ required: true, message: '请输入区域', trigger: 'blur' }],
  ip_address: [{ validator: validateConnectionAddress, trigger: 'blur' }],
  hostname: [{ validator: validateConnectionAddress, trigger: 'blur' }],
  ssh_username: [{ required: true, message: '请输入用户名', trigger: 'blur' }]
}

const formRef = ref()

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    if (!formData.value.ssh_password && !formData.value.ssh_key) {
      ElMessage.error('SSH 密钥或密码至少填写一个')
      return
    }
    loading.value = true
    const payload = { ...formData.value }
    if (!payload.ip_address && payload.hostname) {
      payload.ip_address = payload.hostname
    }
    const response = await addMachine(payload)
    ElMessage.success('添加机器成功')
    const createdId = response.data?.id
    if (createdId) {
      router.push(`/admin/machines/${createdId}`)
    } else {
      router.push('/admin/machines/list')
    }
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

        <el-form-item label="主机名" prop="hostname">
          <el-input v-model="formData.hostname" placeholder="可选，机器主机名" />
        </el-form-item>

        <el-form-item label="区域" prop="region">
          <el-input v-model="formData.region" placeholder="请输入区域,如:北京/上海" />
        </el-form-item>

        <el-alert
          type="info"
          show-icon
          :closable="false"
          title="硬件配置将由系统自动采集，无需填写。"
          class="machine-tip"
        />

        <el-divider content-position="left">登录信息</el-divider>

        <el-form-item label="连接地址" prop="ip_address">
          <el-input v-model="formData.ip_address" placeholder="IP地址或域名" />
        </el-form-item>

        <el-form-item label="公网IP" prop="public_ip">
          <el-input v-model="formData.public_ip" placeholder="可选，公网访问IP" />
        </el-form-item>

        <el-form-item label="SSH主机" prop="ssh_host">
          <el-input v-model="formData.ssh_host" placeholder="可选，SSH连接主机地址（留空则使用公网IP或连接地址）" />
        </el-form-item>

        <el-form-item label="SSH端口" prop="ssh_port">
          <el-input-number v-model="formData.ssh_port" :min="1" :max="65535" />
        </el-form-item>

        <el-form-item label="用户名" prop="ssh_username">
          <el-input v-model="formData.ssh_username" placeholder="请输入SSH用户名" />
        </el-form-item>

        <el-form-item label="SSH私钥" prop="ssh_key">
          <el-input
            v-model="formData.ssh_key"
            type="textarea"
            placeholder="可选，粘贴SSH私钥"
            :rows="3"
          />
        </el-form-item>

        <el-form-item label="密码" prop="ssh_password">
          <el-input
            v-model="formData.ssh_password"
            type="password"
            placeholder="请输入SSH密码(可选)"
            show-password
          />
        </el-form-item>

        <el-divider content-position="left">外映射配置（可选）</el-divider>

        <el-alert
          type="info"
          show-icon
          :closable="false"
          title="配置 Nginx 反向代理或端口映射，用于外部访问 SSH、Jupyter、VNC 等服务。"
          class="machine-tip"
        />

        <el-form-item label="外部IP/域名">
          <el-input v-model="formData.external_ip" placeholder="可选，对外访问的IP或域名" />
        </el-form-item>

        <el-form-item label="SSH映射端口">
          <el-input-number v-model="formData.external_ssh_port" :min="0" :max="65535" />
          <span class="form-tip">对外暴露的SSH端口，0表示不映射</span>
        </el-form-item>

        <el-form-item label="Jupyter映射端口">
          <el-input-number v-model="formData.external_jupyter_port" :min="0" :max="65535" />
          <span class="form-tip">对外暴露的Jupyter端口，0表示不映射</span>
        </el-form-item>

        <el-form-item label="VNC映射端口">
          <el-input-number v-model="formData.external_vnc_port" :min="0" :max="65535" />
          <span class="form-tip">对外暴露的VNC端口，0表示不映射</span>
        </el-form-item>

        <el-form-item label="Nginx域名">
          <el-input v-model="formData.nginx_domain" placeholder="可选，如 gpu-001.example.com" />
        </el-form-item>

        <el-form-item label="Nginx配置路径">
          <el-input v-model="formData.nginx_config_path" placeholder="可选，如 /etc/nginx/conf.d/gpu-001.conf" />
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

.machine-tip {
  margin-bottom: 16px;
}

.form-tip {
  margin-left: 12px;
  color: #909399;
  font-size: 12px;
}
</style>
