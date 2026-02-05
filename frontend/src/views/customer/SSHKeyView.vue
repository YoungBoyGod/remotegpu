<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { getSshKeys, addSshKey, deleteSshKey } from '@/api/customer'

interface SSHKey {
  id: number
  name: string
  fingerprint: string
  created_at: string
}

const keys = ref<SSHKey[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const formRef = ref()
const form = ref({
  name: '',
  publicKey: ''
})
const rules = {
  name: [{ required: true, message: '请输入密钥名称', trigger: 'blur' }],
  publicKey: [{ required: true, message: '请输入公钥内容', trigger: 'blur' }]
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
    if (valid) {
      try {
        await addSshKey(form.value)
        ElMessage.success('添加成功')
        dialogVisible.value = false
        loadKeys()
        form.value = { name: '', publicKey: '' }
      } catch (error) {
        // Error handled by interceptor
      }
    }
  })
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
        <el-button type="primary" :icon="Plus" @click="dialogVisible = true">添加密钥</el-button>
      </template>
    </PageHeader>

    <el-card>
      <el-table :data="keys" v-loading="loading" style="width: 100%">
        <el-table-column prop="name" label="名称" width="200" />
        <el-table-column prop="fingerprint" label="指纹" />
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
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="添加 SSH 密钥" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="例如: My Laptop" />
        </el-form-item>
        <el-form-item label="公钥" prop="publicKey">
          <el-input
            v-model="form.publicKey"
            type="textarea"
            :rows="4"
            placeholder="ssh-rsa AAAAB3NzaC1yc2E..."
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.ssh-key-view {
  padding: 24px;
}
</style>
