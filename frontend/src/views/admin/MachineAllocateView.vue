<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getMachineList, getCustomerList, assignMachine } from '@/api/admin'
import type { Machine } from '@/types/machine'
import type { Customer } from '@/types/customer'
import { ElMessage } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const formRef = ref()

const formData = ref({
  machineIds: [] as string[],
  customerId: null as number | null,
  dateRange: [] as string[],
  contactPerson: '',
  notifyMethods: [] as string[],
  remark: ''
})

const machines = ref<Machine[]>([])
const customers = ref<Customer[]>([])

const rules = {
  machineIds: [{ required: true, message: '请选择机器', trigger: 'change', type: 'array' as const }],
  customerId: [{ required: true, message: '请选择客户', trigger: 'change' }],
  dateRange: [{ required: true, message: '请选择占用时间', trigger: 'change' }]
}

const notifyOptions = [
  { label: '邮件', value: 'email' },
  { label: '短信', value: 'sms' },
  { label: '企业微信', value: 'wechat_work' }
]

const loadMachines = async () => {
  try {
    const response = await getMachineList({ page: 1, pageSize: 1000 })
    machines.value = response.data.list
  } catch (error) {
    console.error('加载机器列表失败:', error)
  }
}

const loadCustomers = async () => {
  try {
    const response = await getCustomerList({ page: 1, pageSize: 1000 })
    customers.value = response.data.list.filter((c: Customer) => c.status === 'active')
  } catch (error) {
    console.error('加载客户列表失败:', error)
  }
}

const getMachineLabel = (machine: Machine) => {
  const parts = [machine.name]
  if (machine.region) parts.push(machine.region)
  if (machine.gpus && machine.gpus.length > 0 && machine.gpus[0]) {
    parts.push(`${machine.gpus.length}x ${machine.gpus[0].name}`)
  }
  return parts.join(' - ')
}

const getCustomerLabel = (customer: Customer) => {
  const name = customer.company || customer.display_name || customer.username || ''
  if (customer.company && customer.username) {
    return `${customer.company} (${customer.username})`
  }
  return name
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const startDate = new Date(formData.value.dateRange[0] || '')
    const endDate = new Date(formData.value.dateRange[1] || '')
    const diffMonths = (endDate.getFullYear() - startDate.getFullYear()) * 12
      + (endDate.getMonth() - startDate.getMonth())
    const durationMonths = Math.max(diffMonths, 1)

    const payload = {
      customer_id: formData.value.customerId!,
      duration_months: durationMonths,
      start_time: formData.value.dateRange[0],
      end_time: formData.value.dateRange[1],
      contact_person: formData.value.contactPerson || undefined,
      notify_methods: formData.value.notifyMethods.length > 0 ? formData.value.notifyMethods : undefined,
      remark: formData.value.remark || undefined
    }
    const results = await Promise.allSettled(
      formData.value.machineIds.map(id => assignMachine(id, payload))
    )
    const failed = results.filter(r => r.status === 'rejected')
    if (failed.length === 0) {
      ElMessage.success(`成功分配 ${results.length} 台机器`)
    } else if (failed.length < results.length) {
      ElMessage.warning(`${results.length - failed.length} 台成功，${failed.length} 台失败`)
    } else {
      ElMessage.error('全部分配失败')
    }
    if (failed.length < results.length) {
      router.push('/admin/allocations/list')
    }
  } catch (error: any) {
    if (error !== false) {
      console.error('机器分配失败:', error)
    }
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  router.back()
}

onMounted(() => {
  loadMachines()
  loadCustomers()
})
</script>

<template>
  <div class="machine-allocate">
    <div class="page-header">
      <h2 class="page-title">机器分配</h2>
    </div>

    <el-card>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="120px"
      >
        <el-form-item label="选择机器" prop="machineIds">
          <el-select
            v-model="formData.machineIds"
            placeholder="请选择机器（可多选）"
            filterable
            multiple
            style="width: 100%"
          >
            <el-option
              v-for="machine in machines"
              :key="machine.id"
              :label="getMachineLabel(machine)"
              :value="String(machine.id)"
            />
          </el-select>
        </el-form-item>

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
              :label="getCustomerLabel(customer)"
              :value="customer.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="占用时间" prop="dateRange">
          <el-date-picker
            v-model="formData.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="内部对接人">
          <el-input
            v-model="formData.contactPerson"
            placeholder="请输入内部对接人"
          />
        </el-form-item>

        <el-form-item label="通知方式">
          <el-checkbox-group v-model="formData.notifyMethods">
            <el-checkbox
              v-for="opt in notifyOptions"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="备注">
          <el-input
            v-model="formData.remark"
            type="textarea"
            :rows="3"
            placeholder="可选备注信息"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            确认分配
          </el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.machine-allocate {
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
