<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { addCustomer } from '@/api/admin'
import type { AddCustomerForm } from '@/types/customer'
import { ElMessage } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const formRef = ref()

const formData = ref<AddCustomerForm>({
  username: '',
  company_code: '',
  company: '',
  email: '',
  phone: '',
  password: 'ChangeME_123'
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  company_code: [{ required: true, message: '请输入公司代号', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }
  ],
  phone: [{ message: '请输入联系电话', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    await addCustomer(formData.value)
    ElMessage.success('添加客户成功')
    router.push('/admin/customers/list')
  } catch (error: any) {
    if (error !== false) {
      console.error('添加客户失败:', error)
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
  <div class="customer-add">
    <div class="page-header">
      <h2 class="page-title">添加客户</h2>
    </div>

    <el-card>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="120px"
      >
        <el-divider content-position="left">账号信息</el-divider>

        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="用于登录" />
        </el-form-item>

        <el-form-item label="公司代号" prop="company_code">
          <el-input v-model="formData.company_code" placeholder="请输入公司代号" />
        </el-form-item>

        <el-form-item label="公司名称" prop="company">
          <el-input v-model="formData.company" placeholder="可选" />
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>

        <el-form-item label="电话（可选）" prop="phone">
          <el-input v-model="formData.phone" placeholder="可选" />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="formData.password" type="password" show-password placeholder="默认 ChangeME_123" />
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
.customer-add {
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
