<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, VideoPause } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { getTasks, createTrainingTask, stopTask } from '@/api/customer'

interface Task {
  id: number
  type: string
  status: string
  config: any
  created_at: string
}

const tasks = ref<Task[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const form = ref({
  machine_id: '',
  image: '',
  command: ''
})

const loadTasks = async () => {
  loading.value = true
  try {
    const res = await getTasks({ page: 1, pageSize: 20 })
    tasks.value = res.data.list
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = async () => {
  try {
    await createTrainingTask(form.value)
    ElMessage.success('任务创建成功')
    dialogVisible.value = false
    loadTasks()
  } catch (error) {
    //
  }
}

const handleStop = async (task: Task) => {
  try {
    await ElMessageBox.confirm('确定停止该任务吗?', '提示', { type: 'warning' })
    await stopTask(task.id)
    ElMessage.success('任务已停止')
    loadTasks()
  } catch {
    //
  }
}

onMounted(() => {
  loadTasks()
})
</script>

<template>
  <div class="task-list-view">
    <PageHeader title="任务管理">
      <template #actions>
        <el-button type="primary" :icon="Plus" @click="dialogVisible = true">创建任务</el-button>
      </template>
    </PageHeader>

    <el-card>
      <el-table :data="tasks" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column label="配置">
          <template #default="{ row }">
            <pre class="config-pre">{{ JSON.stringify(row.config, null, 2) }}</pre>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'running' ? 'success' : 'info'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'running' || row.status === 'pending'"
              type="danger"
              size="small"
              :icon="VideoPause"
              @click="handleStop(row)"
            >停止</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="创建训练任务">
      <el-form :model="form" label-width="100px">
        <el-form-item label="机器ID">
          <el-input v-model="form.machine_id" placeholder="Host ID" />
        </el-form-item>
        <el-form-item label="镜像">
          <el-input v-model="form.image" placeholder="docker image" />
        </el-form-item>
        <el-form-item label="启动命令">
          <el-input v-model="form.command" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.task-list-view {
  padding: 24px;
}
.config-pre {
  margin: 0;
  font-size: 12px;
  color: #666;
  max-height: 100px;
  overflow-y: auto;
}
</style>
