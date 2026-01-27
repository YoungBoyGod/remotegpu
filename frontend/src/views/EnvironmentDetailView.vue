<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useRoleNavigation } from '@/composables/useRoleNavigation'

const route = useRoute()
const router = useRouter()
const { navigateTo } = useRoleNavigation()
const activeTab = ref('overview')

const environment = ref({
  id: '',
  name: '',
  status: 'running',
  image: '',
  cpu: 0,
  memory: 0,
  gpu: '',
  storage: 0,
  sshCommand: '',
  jupyterUrl: '',
  createdAt: '',
  runningTime: '',
})

const loadEnvironment = async () => {
  const id = route.params.id
  // 模拟数据
  environment.value = {
    id: id as string,
    name: 'PyTorch 训练环境',
    status: 'running',
    image: 'pytorch/pytorch:2.0-cuda11.8',
    cpu: 8,
    memory: 32,
    gpu: 'Tesla V100 x2',
    storage: 100,
    sshCommand: 'ssh -p 30001 root@gpu.example.com',
    jupyterUrl: 'http://gpu.example.com:38001',
    createdAt: '2026-01-26 10:30',
    runningTime: '2小时30分',
  }
}

const startEnvironment = async () => {
  ElMessage.success('环境启动中...')
}

const stopEnvironment = async () => {
  ElMessage.success('环境已停止')
}

onMounted(() => {
  loadEnvironment()
})
</script>

<template>
  <div class="environment-detail">
    <div class="page-header">
      <div>
        <h1>{{ environment.name }}</h1>
        <el-tag v-if="environment.status === 'running'" type="success">运行中</el-tag>
        <el-tag v-else type="info">已停止</el-tag>
      </div>
      <div class="actions">
        <el-button v-if="environment.status === 'stopped'" type="success" @click="startEnvironment">
          启动
        </el-button>
        <el-button v-if="environment.status === 'running'" type="warning" @click="stopEnvironment">
          停止
        </el-button>
        <el-button @click="navigateTo('/environments')">返回列表</el-button>
      </div>
    </div>

    <el-tabs v-model="activeTab">
      <el-tab-pane label="概览" name="overview">
        <div class="overview-content">
          <el-card class="info-card">
            <template #header>基本信息</template>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="环境ID">{{ environment.id }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ environment.createdAt }}</el-descriptions-item>
              <el-descriptions-item label="运行时长">{{ environment.runningTime }}</el-descriptions-item>
              <el-descriptions-item label="镜像">{{ environment.image }}</el-descriptions-item>
            </el-descriptions>
          </el-card>

          <el-card class="info-card">
            <template #header>资源配置</template>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="CPU">{{ environment.cpu }} 核</el-descriptions-item>
              <el-descriptions-item label="内存">{{ environment.memory }} GB</el-descriptions-item>
              <el-descriptions-item label="GPU">{{ environment.gpu }}</el-descriptions-item>
              <el-descriptions-item label="存储">{{ environment.storage }} GB</el-descriptions-item>
            </el-descriptions>
          </el-card>

          <el-card class="info-card">
            <template #header>访问信息</template>
            <div class="access-info">
              <div class="access-item">
                <label>SSH 连接：</label>
                <el-input v-model="environment.sshCommand" readonly>
                  <template #append>
                    <el-button>复制</el-button>
                  </template>
                </el-input>
              </div>
              <div class="access-item">
                <label>JupyterLab：</label>
                <el-link :href="environment.jupyterUrl" target="_blank" type="primary">
                  {{ environment.jupyterUrl }}
                </el-link>
              </div>
            </div>
          </el-card>
        </div>
      </el-tab-pane>

      <el-tab-pane label="监控" name="monitoring">
        <div class="monitoring-content">
          <el-empty description="监控数据加载中..." />
        </div>
      </el-tab-pane>

      <el-tab-pane label="日志" name="logs">
        <div class="logs-content">
          <el-input
            type="textarea"
            :rows="20"
            readonly
            placeholder="日志内容..."
          />
        </div>
      </el-tab-pane>

      <el-tab-pane label="设置" name="settings">
        <div class="settings-content">
          <el-form label-width="120px">
            <el-form-item label="环境名称">
              <el-input v-model="environment.name" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary">保存</el-button>
              <el-button type="danger">删除环境</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped>
.environment-detail {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-header h1 {
  font-size: 24px;
  font-weight: 600;
  margin: 0 12px 0 0;
  display: inline-block;
}

.actions {
  display: flex;
  gap: 12px;
}

.overview-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-card {
  margin-bottom: 0;
}

.access-info {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.access-item label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #606266;
}

.monitoring-content,
.logs-content,
.settings-content {
  padding: 20px;
  background: white;
  border-radius: 8px;
}
</style>
