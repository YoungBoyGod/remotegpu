<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete, CopyDocument } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { getSshKeys, addSshKey, deleteSshKey } from '@/api/customer'

interface SSHKey {
  id: number
  name: string
  fingerprint: string
  key_type?: string
  public_key?: string
  created_at: string
}

const SUPPORTED_KEY_TYPES = ['ssh-rsa', 'ssh-ed25519', 'ecdsa-sha2-nistp256', 'ecdsa-sha2-nistp384', 'ecdsa-sha2-nistp521']

// 从公钥内容中解析密钥类型
const parseKeyType = (publicKey: string): string => {
  const trimmed = publicKey.trim()
  const prefix = trimmed.split(/\s+/)[0] || ''
  if (SUPPORTED_KEY_TYPES.includes(prefix)) return prefix
  return ''
}

// 密钥类型显示名称
const keyTypeLabel = (type?: string): string => {
  if (!type) return '-'
  if (type.startsWith('ecdsa-')) return 'ECDSA'
  if (type === 'ssh-ed25519') return 'ED25519'
  if (type === 'ssh-rsa') return 'RSA'
  return type
}

// 公钥格式校验
const validatePublicKey = (_rule: any, value: string, callback: (error?: Error) => void) => {
  if (!value) {
    callback(new Error('请输入公钥内容'))
    return
  }
  const trimmed = value.trim()
  const keyType = parseKeyType(trimmed)
  if (!keyType) {
    callback(new Error('不支持的密钥格式，请使用 ssh-rsa、ssh-ed25519 或 ecdsa 格式'))
    return
  }
  const parts = trimmed.split(/\s+/)
  if (parts.length < 2 || !parts[1]) {
    callback(new Error('公钥格式不完整，缺少密钥数据'))
    return
  }
  callback()
}

const keys = ref<SSHKey[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const submitLoading = ref(false)
const formRef = ref()
const form = ref({
  name: '',
  publicKey: ''
})
const rules = {
  name: [
    { required: true, message: '请输入密钥名称', trigger: 'blur' },
    { max: 64, message: '名称不超过 64 个字符', trigger: 'blur' },
  ],
  publicKey: [
    { required: true, validator: validatePublicKey, trigger: 'blur' },
  ],
}

const loadKeys = async () => {
  loading.value = true
  try {
    const res = await getSshKeys()
    keys.value = res.data
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return
    submitLoading.value = true
    try {
      await addSshKey(form.value)
      ElMessage.success('添加成功')
      dialogVisible.value = false
      loadKeys()
      form.value = { name: '', publicKey: '' }
    } catch (error) {
      // Error handled by interceptor
    } finally {
      submitLoading.value = false
    }
  })
}

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

// 打开添加对话框
const openAddDialog = () => {
  form.value = { name: '', publicKey: '' }
  dialogVisible.value = true
  // 重置表单验证状态
  formRef.value?.resetFields()
}

const handleDelete = async (key: SSHKey) => {
  try {
    await ElMessageBox.confirm(`确认删除密钥 "${key.name}"?`, '提示', {
      type: 'warning'
    })
    await deleteSshKey(key.id)
    ElMessage.success('删除成功')
    loadKeys()
  } catch {
    // Cancelled
  }
}

onMounted(() => {
  loadKeys()
})
</script>

<template>
  <div class="ssh-key-view">
    <PageHeader title="SSH 密钥管理">
      <template #actions>
        <el-button type="primary" :icon="Plus" @click="openAddDialog">添加密钥</el-button>
      </template>
    </PageHeader>

    <el-card>
      <el-table :data="keys" v-loading="loading" style="width: 100%">
        <el-table-column prop="name" label="名称" width="200" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ keyTypeLabel(row.key_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="指纹" min-width="280">
          <template #default="{ row }">
            <div class="fingerprint-cell">
              <code class="fingerprint-text">{{ row.fingerprint }}</code>
              <el-button link :icon="CopyDocument" size="small" @click="copyToClipboard(row.fingerprint)" />
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button type="danger" :icon="Delete" circle size="small" @click="handleDelete(row)" />
          </template>
        </el-table-column>
        <template #empty>
          <el-empty description="暂无 SSH 密钥">
            <el-button type="primary" :icon="Plus" @click="openAddDialog">添加密钥</el-button>
          </el-empty>
        </template>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="添加 SSH 密钥" width="560px" :close-on-click-modal="false">
      <el-alert
        type="info"
        show-icon
        :closable="false"
        title="支持 ssh-rsa、ssh-ed25519、ecdsa 格式的公钥"
        style="margin-bottom: 16px"
      />
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="例如: My Laptop" />
        </el-form-item>
        <el-form-item label="公钥" prop="publicKey">
          <el-input
            v-model="form.publicKey"
            type="textarea"
            :rows="5"
            placeholder="ssh-rsa AAAAB3NzaC1yc2E... user@host"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.ssh-key-view {
  padding: 24px;
}

.fingerprint-cell {
  display: flex;
  align-items: center;
  gap: 4px;
}

.fingerprint-text {
  font-size: 13px;
  color: #606266;
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 3px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
