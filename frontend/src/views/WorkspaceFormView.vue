<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import PageHeader from '@/components/common/PageHeader.vue'
import { createWorkspace, updateWorkspace, getWorkspaceById } from '@/api/workspace'
import type { CreateWorkspaceRequest, UpdateWorkspaceRequest } from '@/api/workspace/types'

const router = useRouter()
const route = useRoute()

const formRef = ref<FormInstance>()
const loading = ref(false)
const submitting = ref(false)

// 判断是编辑模式还是创建模式
const isEditMode = computed(() => !!route.params.id)
const workspaceId = computed(() => Number(route.params.id))

const formData = reactive<CreateWorkspaceRequest>({
  name: '',
  description: '',
})

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: '请输入工作空间名称', trigger: 'blur' },
    { min: 2, max: 100, message: '名称长度在 2 到 100 个字符', trigger: 'blur' },
  ],
  description: [
    { max: 500, message: '描述不能超过 500 个字符', trigger: 'blur' },
  ],
}

// 加载工作空间详情(编辑模式)
const loadWorkspace = async () => {
  if (!isEditMode.value) return

  loading.value = true
  try {
    const response = await getWorkspaceById(workspaceId.value)
    formData.name = response.data.name
    formData.description = response.data.description
  } catch (error: any) {
    ElMessage.error(error.message || '加载工作空间信息失败')
    router.back()
  } finally {
    loading.value = false
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    if (isEditMode.value) {
      // 编辑模式
      const updateData: UpdateWorkspaceRequest = {
        name: formData.name,
        description: formData.description,
      }
      await updateWorkspace(workspaceId.value, updateData)
      ElMessage.success('更新成功')
    } else {
      // 创建模式
      await createWorkspace(formData)
      ElMessage.success('创建成功')
    }

    router.push('/portal/workspaces')
  } catch (error: any) {
    if (error.message) {
      ElMessage.error(error.message)
    }
  } finally {
    submitting.value = false
  }
}

// 取消
const handleCancel = () => {
  router.back()
}

onMounted(() => {
  loadWorkspace()
})
</script>

<template>
  <div class="workspace-form">
    <PageHeader :title="isEditMode ? '编辑工作空间' : '创建工作空间'" />

    <el-card v-loading="loading">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="120px"
        style="max-width: 600px"
      >
        <el-form-item label="工作空间名称" prop="name">
          <el-input
            v-model="formData.name"
            placeholder="请输入工作空间名称"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="4"
            placeholder="请输入工作空间描述(可选)"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            :loading="submitting"
            @click="handleSubmit"
          >
            {{ isEditMode ? '保存' : '创建' }}
          </el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.workspace-form {
  padding: 24px;
}
</style>
