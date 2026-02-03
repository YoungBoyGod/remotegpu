<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { quickAllocate, getAvailableMachines } from '@/api/admin'
import { getCustomerList } from '@/api/admin'
import type { QuickAllocateForm } from '@/types/allocation'
import type { Machine } from '@/types/machine'
import type { Customer } from '@/types/customer'
import { ElMessage } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const formData = ref<QuickAllocateForm>({
  customerId: null,
  machineIds: [],
  duration: 30,
  notes: ''
})

const customers = ref<Customer[]>([])
const machines = ref<Machine[]>([])

const rules = {
  customerId: [{ required: true, message: '请选择客户', trigger: 'change' }],
  machineIds: [{ required: true, message: '请选择机器', trigger: 'change' }],
  duration: [{ required: true, message: '请输入分配时长', trigger: 'blur' }]
}

const formRef = ref()

const loadCustomers = async () => {
  try {
    const response = await getCustomerList({ page: 1, pageSize: 1000 })
    customers.value = response.data.list.filter(c => c.status === 'active')
  } catch (error) {
    console.error('加载客户列表失败:', error)
  }
}

const loadMachines = async () => {
  try {
    const response = await getAvailableMachines()
    machines.value = response.data
  } catch (error) {
    console.error('加载可用机器失败:', error)
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    await quickAllocate(formData.value)
    ElMessage.success('分配成功')
    router.push('/admin/allocations/list')
  } catch (error: any) {
    if (error !== false) {
      console.error('分配失败:', error)
    }
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  router.back()
}

onMounted(() => {
  loadCustomers()
  loadMachines()
})
</script>

<template>
  <div class="quick-allocate">
    <div class="page-header">
      <h2 class="page-title">快速分配</h2>
    </div>

    <el-card>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="120px"
      >
        <el-form-item label="选择客户" prop="customerId">
          <el-select
            v-model="formData.customerId"
            placeholder="请选择客户"
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="customer in customers"
              :key="customer.id"
              :label="customer.name"
              :value="customer.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="选择机器" prop="machineIds">
          <el-select
            v-model="formData.machineIds"
            placeholder="请选择机器"
            filterable
            multiple
            style="width: 100%"
          >
            <el-option
              v-for="machine in machines"
              :key="machine.id"
              :label="`${machine.name} (${machine.region}) - ${machine.gpuCount}x ${machine.gpuModel}`"
              :value="machine.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="分配时长" prop="duration">
          <el-input-number
            v-model="formData.duration"
            :min="1"
            :max="365"
            style="width: 200px"
          />
          <span style="margin-left: 10px">天</span>
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
.quick-allocate {
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
