<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Upload, Delete, Plus } from '@element-plus/icons-vue'
import { batchImportMachines } from '@/api/admin'
import type { ImportMachineItem } from '@/api/admin'

const router = useRouter()
const loading = ref(false)

const emptyRow = (): ImportMachineItem => ({
  host_ip: '',
  ssh_port: 22,
  region: '',
  gpu_model: '',
  gpu_count: 1,
  cpu_cores: 0,
  ram_size: 0,
  disk_size: 0,
  price_hourly: 0,
})

const machines = ref<ImportMachineItem[]>([emptyRow()])

const addRow = () => {
  machines.value.push(emptyRow())
}

const removeRow = (index: number) => {
  if (machines.value.length <= 1) return
  machines.value.splice(index, 1)
}

const validate = (): boolean => {
  for (let i = 0; i < machines.value.length; i++) {
    const m = machines.value[i]!
    if (!m.host_ip) {
      ElMessage.error(`第 ${i + 1} 行：请填写 IP 地址`)
      return false
    }
    if (!m.ssh_port || m.ssh_port < 1 || m.ssh_port > 65535) {
      ElMessage.error(`第 ${i + 1} 行：SSH 端口无效`)
      return false
    }
    if (!m.region) {
      ElMessage.error(`第 ${i + 1} 行：请填写区域`)
      return false
    }
    if (!m.gpu_model) {
      ElMessage.error(`第 ${i + 1} 行：请填写 GPU 型号`)
      return false
    }
  }
  return true
}

const handleSubmit = async () => {
  if (!validate()) return
  loading.value = true
  try {
    const res = await batchImportMachines({ machines: machines.value })
    ElMessage.success(`导入成功，共 ${res.data.count} 台机器`)
    router.push('/admin/machines/list')
  } catch (error: any) {
    ElMessage.error(error?.msg || '导入失败')
  } finally {
    loading.value = false
  }
}

const handleFileUpload = (file: File) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const text = e.target?.result as string
      const parsed = parseCSV(text)
      if (parsed.length === 0) {
        ElMessage.error('CSV 文件为空或格式不正确')
        return
      }
      machines.value = parsed
      ElMessage.success(`已解析 ${parsed.length} 条记录`)
    } catch {
      ElMessage.error('CSV 解析失败，请检查文件格式')
    }
  }
  reader.readAsText(file)
  return false
}

const parseCSV = (text: string): ImportMachineItem[] => {
  const lines = text.trim().split('\n')
  if (lines.length < 2) return []
  const result: ImportMachineItem[] = []
  for (let i = 1; i < lines.length; i++) {
    const line = lines[i]
    if (!line) continue
    const cols = line.split(',').map(s => s.trim())
    if (cols.length < 9) continue
    result.push({
      host_ip: cols[0] ?? '',
      ssh_port: Number(cols[1]) || 22,
      region: cols[2] ?? '',
      gpu_model: cols[3] ?? '',
      gpu_count: Number(cols[4]) || 1,
      cpu_cores: Number(cols[5]) || 0,
      ram_size: Number(cols[6]) || 0,
      disk_size: Number(cols[7]) || 0,
      price_hourly: Number(cols[8]) || 0,
    })
  }
  return result
}

const handleBack = () => {
  router.back()
}
</script>

<template>
  <div class="machine-import">
    <div class="page-header">
      <div>
        <el-button @click="handleBack">返回</el-button>
        <h2 class="page-title">批量导入机器</h2>
      </div>
    </div>

    <el-card class="upload-card">
      <el-alert
        type="info"
        show-icon
        :closable="false"
        title="支持 CSV 文件导入，格式：host_ip, ssh_port, region, gpu_model, gpu_count, cpu_cores, ram_size(GB), disk_size(GB), price_hourly(分)"
        style="margin-bottom: 16px"
      />
      <el-upload
        accept=".csv"
        :show-file-list="false"
        :before-upload="handleFileUpload"
      >
        <el-button :icon="Upload">上传 CSV 文件</el-button>
      </el-upload>
    </el-card>

    <el-card>
      <div class="table-header">
        <span class="table-title">机器列表（{{ machines.length }} 台）</span>
        <el-button type="primary" :icon="Plus" size="small" @click="addRow">添加一行</el-button>
      </div>

      <el-table :data="machines" border stripe style="width: 100%">
        <el-table-column label="序号" width="60" type="index" />
        <el-table-column label="IP 地址" min-width="150">
          <template #default="{ row }">
            <el-input v-model="row.host_ip" placeholder="192.168.1.100" size="small" />
          </template>
        </el-table-column>
        <el-table-column label="SSH 端口" width="100">
          <template #default="{ row }">
            <el-input-number v-model="row.ssh_port" :min="1" :max="65535" size="small" controls-position="right" />
          </template>
        </el-table-column>
        <el-table-column label="区域" width="120">
          <template #default="{ row }">
            <el-input v-model="row.region" placeholder="北京" size="small" />
          </template>
        </el-table-column>
        <el-table-column label="GPU 型号" min-width="140">
          <template #default="{ row }">
            <el-input v-model="row.gpu_model" placeholder="A100" size="small" />
          </template>
        </el-table-column>
        <el-table-column label="GPU 数量" width="100">
          <template #default="{ row }">
            <el-input-number v-model="row.gpu_count" :min="0" size="small" controls-position="right" />
          </template>
        </el-table-column>
        <el-table-column label="CPU 核数" width="100">
          <template #default="{ row }">
            <el-input-number v-model="row.cpu_cores" :min="0" size="small" controls-position="right" />
          </template>
        </el-table-column>
        <el-table-column label="内存(GB)" width="100">
          <template #default="{ row }">
            <el-input-number v-model="row.ram_size" :min="0" size="small" controls-position="right" />
          </template>
        </el-table-column>
        <el-table-column label="磁盘(GB)" width="100">
          <template #default="{ row }">
            <el-input-number v-model="row.disk_size" :min="0" size="small" controls-position="right" />
          </template>
        </el-table-column>
        <el-table-column label="时价(分)" width="100">
          <template #default="{ row }">
            <el-input-number v-model="row.price_hourly" :min="0" size="small" controls-position="right" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="70" fixed="right">
          <template #default="{ $index }">
            <el-button
              link
              type="danger"
              size="small"
              :icon="Delete"
              :disabled="machines.length <= 1"
              @click="removeRow($index)"
            />
          </template>
        </el-table-column>
      </el-table>

      <div class="submit-bar">
        <el-button @click="handleBack">取消</el-button>
        <el-button type="primary" :loading="loading" :icon="Upload" @click="handleSubmit">
          确认导入（{{ machines.length }} 台）
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.machine-import {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 8px 0 0 0;
}

.upload-card {
  margin-bottom: 16px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.table-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.submit-bar {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #ebeef5;
}
</style>
