<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/common/PageHeader.vue'

interface Machine {
  id: number
  name: string
  region: string
  gpuModel: string
  gpuMemory: number
  cpu: string
  memory: number
  disk: number
  cudaVersion: string
  gpuDriver: string
  status: 'available' | 'allocated'
  loginInfo: {
    ip: string
    port: number
    username: string
    password: string
  }
  allocatedTo?: {
    customerId: number
    customerName: string
    allocatedAt: string
    duration: number // 月数
  }
}

interface Customer {
  id: number
  name: string
  company: string
}

interface GpuCategory {
  model: string
  memory: number
  available: Machine[]
  allocated: Machine[]
}

const loading = ref(false)
const machines = ref<Machine[]>([])
const customers = ref<Customer[]>([])

// 分配对话框
const allocationDialogVisible = ref(false)
const currentMachine = ref<Machine | null>(null)
const allocationForm = ref({
  customerId: null as number | null,
  duration: 1, // 默认1个月
  notes: ''
})

// 筛选条件
const selectedGpuModel = ref<string>('all')
const selectedRegion = ref<string>('all')
const selectedStatus = ref<string>('all')

// 获取所有GPU型号
const gpuModels = computed(() => {
  const models = new Set<string>()
  machines.value.forEach(m => models.add(m.gpuModel))
  return ['all', ...Array.from(models)]
})

// 获取所有地区
const regions = computed(() => {
  const regionSet = new Set<string>()
  machines.value.forEach(m => regionSet.add(m.region))
  return ['all', ...Array.from(regionSet)]
})

// 按GPU型号分组（带筛选）
const gpuCategories = computed<GpuCategory[]>(() => {
  // 先筛选机器
  let filteredMachines = machines.value

  if (selectedGpuModel.value !== 'all') {
    filteredMachines = filteredMachines.filter(m => m.gpuModel === selectedGpuModel.value)
  }

  if (selectedRegion.value !== 'all') {
    filteredMachines = filteredMachines.filter(m => m.region === selectedRegion.value)
  }

  if (selectedStatus.value !== 'all') {
    filteredMachines = filteredMachines.filter(m => m.status === selectedStatus.value)
  }

  // 再按GPU型号分组
  const categoryMap = new Map<string, GpuCategory>()

  filteredMachines.forEach(machine => {
    const key = `${machine.gpuModel}-${machine.gpuMemory}`

    if (!categoryMap.has(key)) {
      categoryMap.set(key, {
        model: machine.gpuModel,
        memory: machine.gpuMemory,
        available: [],
        allocated: []
      })
    }

    const category = categoryMap.get(key)!
    if (machine.status === 'available') {
      category.available.push(machine)
    } else {
      category.allocated.push(machine)
    }
  })

  return Array.from(categoryMap.values())
})

// 加载机器列表
const loadMachines = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 500))
    machines.value = [
      // RTX 5090 - 未分配
      {
        id: 1,
        name: '北京B区-598机',
        region: '北京B区',
        gpuModel: 'RTX 5090',
        gpuMemory: 32,
        cpu: '16核',
        memory: 90,
        disk: 5881,
        cudaVersion: '≤ 13.0',
        gpuDriver: '580.76.05',
        status: 'available',
        loginInfo: {
          ip: '192.168.1.101',
          port: 22,
          username: 'root',
          password: 'Abc123456'
        }
      },
      {
        id: 2,
        name: '北京B区-353机',
        region: '北京B区',
        gpuModel: 'RTX 5090',
        gpuMemory: 32,
        cpu: '25核',
        memory: 90,
        disk: 3810,
        cudaVersion: '≤ 13.0',
        gpuDriver: '580.76.05',
        status: 'available',
        loginInfo: {
          ip: '192.168.1.102',
          port: 22,
          username: 'root',
          password: 'Abc123456'
        }
      },
      // RTX 5090 - 已分配
      {
        id: 3,
        name: '北京A区-201机',
        region: '北京A区',
        gpuModel: 'RTX 5090',
        gpuMemory: 32,
        cpu: '32核',
        memory: 128,
        disk: 7200,
        cudaVersion: '≤ 13.0',
        gpuDriver: '580.76.05',
        status: 'allocated',
        loginInfo: {
          ip: '192.168.2.201',
          port: 22,
          username: 'root',
          password: 'Abc123456'
        },
        allocatedTo: {
          customerId: 1,
          customerName: '张三',
          allocatedAt: '2026-01-15',
          duration: 3
        }
      },
      // RTX 4090 - 未分配
      {
        id: 4,
        name: '西北B区-102机',
        region: '西北B区',
        gpuModel: 'RTX 4090',
        gpuMemory: 24,
        cpu: '16核',
        memory: 64,
        disk: 4000,
        cudaVersion: '≤ 12.0',
        gpuDriver: '535.54.03',
        status: 'available',
        loginInfo: {
          ip: '192.168.3.102',
          port: 22,
          username: 'root',
          password: 'Abc123456'
        }
      },
      {
        id: 5,
        name: '重庆A区-301机',
        region: '重庆A区',
        gpuModel: 'RTX 4090',
        gpuMemory: 24,
        cpu: '16核',
        memory: 64,
        disk: 4000,
        cudaVersion: '≤ 12.0',
        gpuDriver: '535.54.03',
        status: 'available',
        loginInfo: {
          ip: '192.168.4.301',
          port: 22,
          username: 'root',
          password: 'Abc123456'
        }
      },
      // RTX 4090 - 已分配
      {
        id: 6,
        name: '重庆A区-302机',
        region: '重庆A区',
        gpuModel: 'RTX 4090',
        gpuMemory: 24,
        cpu: '16核',
        memory: 64,
        disk: 4000,
        cudaVersion: '≤ 12.0',
        gpuDriver: '535.54.03',
        status: 'allocated',
        loginInfo: {
          ip: '192.168.4.302',
          port: 22,
          username: 'root',
          password: 'Abc123456'
        },
        allocatedTo: {
          customerId: 2,
          customerName: '李四',
          allocatedAt: '2026-01-20',
          duration: 1
        }
      },
      // vGPU-48GB - 未分配
      {
        id: 7,
        name: '内蒙B区-401机',
        region: '内蒙B区',
        gpuModel: 'vGPU-48GB',
        gpuMemory: 48,
        cpu: '24核',
        memory: 96,
        disk: 6000,
        cudaVersion: '≤ 12.0',
        gpuDriver: '535.54.03',
        status: 'available',
        loginInfo: {
          ip: '192.168.5.401',
          port: 22,
          username: 'root',
          password: 'Abc123456'
        }
      }
    ]
  } catch (error) {
    ElMessage.error('加载机器列表失败')
  } finally {
    loading.value = false
  }
}

// 加载客户列表
const loadCustomers = async () => {
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 300))
    customers.value = [
      { id: 1, name: '张三', company: '科技有限公司' },
      { id: 2, name: '李四', company: '数据科技公司' },
      { id: 3, name: '王五', company: '智能科技公司' },
      { id: 4, name: '赵六', company: '云计算公司' }
    ]
  } catch (error) {
    ElMessage.error('加载客户列表失败')
  }
}

// 打开分配对话框
const handleAllocate = (machine: Machine) => {
  currentMachine.value = machine
  allocationForm.value = {
    customerId: null,
    duration: 1,
    notes: ''
  }
  allocationDialogVisible.value = true
}

// 提交分配
const handleSubmitAllocation = async () => {
  if (!allocationForm.value.customerId) {
    ElMessage.warning('请选择客户')
    return
  }

  try {
    // TODO: 调用API提交分配
    await new Promise(resolve => setTimeout(resolve, 500))

    const customer = customers.value.find(c => c.id === allocationForm.value.customerId)

    // 更新机器状态
    if (currentMachine.value) {
      currentMachine.value.status = 'allocated'
      currentMachine.value.allocatedTo = {
        customerId: allocationForm.value.customerId!,
        customerName: customer?.name || '',
        allocatedAt: new Date().toLocaleDateString('zh-CN'),
        duration: allocationForm.value.duration
      }
    }

    ElMessage.success(`已成功将 ${currentMachine.value?.name} 分配给 ${customer?.name}，时长 ${allocationForm.value.duration} 个月`)
    allocationDialogVisible.value = false
  } catch (error) {
    ElMessage.error('分配失败')
  }
}

// 回收机器
const handleReclaim = async (machine: Machine) => {
  try {
    await ElMessageBox.confirm(
      `确认回收机器 ${machine.name}？该机器当前分配给 ${machine.allocatedTo?.customerName}`,
      '回收确认',
      {
        confirmButtonText: '确认回收',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    // TODO: 调用API回收机器
    await new Promise(resolve => setTimeout(resolve, 500))

    // 更新机器状态
    machine.status = 'available'
    delete machine.allocatedTo

    ElMessage.success(`已成功回收机器 ${machine.name}`)
  } catch {
    // 用户取消
  }
}

// 一键登录
const handleQuickLogin = async (machine: Machine) => {
  // 构造WebSSH URL (假设WebSSH服务部署在 /webssh 路径)
  // 格式: /webssh?hostname=ip&port=port&username=username&password=password
  const websshUrl = `/webssh?hostname=${machine.loginInfo.ip}&port=${machine.loginInfo.port}&username=${machine.loginInfo.username}&password=${encodeURIComponent(machine.loginInfo.password)}`

  // 构造Jumpserver URL (假设Jumpserver部署在配置的地址)
  // 格式: http://jumpserver/terminal/connect/?asset_id=xxx
  const jumpserverUrl = `http://192.168.10.249:8080/terminal/connect/?hostname=${machine.loginInfo.ip}&port=${machine.loginInfo.port}&username=${machine.loginInfo.username}`

  // SSH命令(作为备选方案)
  const sshCommand = `ssh ${machine.loginInfo.username}@${machine.loginInfo.ip} -p ${machine.loginInfo.port}`

  // 弹出选择对话框
  try {
    await ElMessageBox.confirm(
      '请选择登录方式',
      '一键登录',
      {
        distinguishCancelAndClose: true,
        confirmButtonText: 'WebSSH (浏览器)',
        cancelButtonText: 'Jumpserver',
        type: 'info'
      }
    )
    // 用户选择WebSSH
    window.open(websshUrl, '_blank')
  } catch (action) {
    if (action === 'cancel') {
      // 用户选择Jumpserver
      window.open(jumpserverUrl, '_blank')
    } else {
      // 用户关闭对话框,复制SSH命令作为备选
      try {
        await navigator.clipboard.writeText(sshCommand)
        ElMessage.info({
          message: `SSH命令已复制到剪贴板\n密码: ${machine.loginInfo.password}`,
          duration: 5000,
          showClose: true
        })
      } catch (error) {
        ElMessageBox.alert(
          `SSH命令: ${sshCommand}\n密码: ${machine.loginInfo.password}`,
          '登录信息',
          {
            confirmButtonText: '确定',
            type: 'info'
          }
        )
      }
    }
  }
}

onMounted(() => {
  loadMachines()
  loadCustomers()
})
</script>

<template>
  <div class="computing-market">
    <PageHeader title="算力市场" />

    <div class="market-container" v-loading="loading">
      <!-- 筛选标签区域 -->
      <div class="filter-section">
        <!-- GPU型号筛选 -->
        <div class="filter-group">
          <span class="filter-label">显卡分类：</span>
          <div class="filter-tags">
            <el-tag
              v-for="model in gpuModels"
              :key="model"
              :type="selectedGpuModel === model ? 'primary' : 'info'"
              :effect="selectedGpuModel === model ? 'dark' : 'plain'"
              class="filter-tag"
              @click="selectedGpuModel = model"
            >
              {{ model === 'all' ? '全部' : model }}
            </el-tag>
          </div>
        </div>

        <!-- 地区筛选 -->
        <div class="filter-group">
          <span class="filter-label">地区分类：</span>
          <div class="filter-tags">
            <el-tag
              v-for="region in regions"
              :key="region"
              :type="selectedRegion === region ? 'success' : 'info'"
              :effect="selectedRegion === region ? 'dark' : 'plain'"
              class="filter-tag"
              @click="selectedRegion = region"
            >
              {{ region === 'all' ? '全部' : region }}
            </el-tag>
          </div>
        </div>

        <!-- 状态筛选 -->
        <div class="filter-group">
          <span class="filter-label">状态分类：</span>
          <div class="filter-tags">
            <el-tag
              :type="selectedStatus === 'all' ? 'warning' : 'info'"
              :effect="selectedStatus === 'all' ? 'dark' : 'plain'"
              class="filter-tag"
              @click="selectedStatus = 'all'"
            >
              全部
            </el-tag>
            <el-tag
              :type="selectedStatus === 'available' ? 'warning' : 'info'"
              :effect="selectedStatus === 'available' ? 'dark' : 'plain'"
              class="filter-tag"
              @click="selectedStatus = 'available'"
            >
              未分配
            </el-tag>
            <el-tag
              :type="selectedStatus === 'allocated' ? 'warning' : 'info'"
              :effect="selectedStatus === 'allocated' ? 'dark' : 'plain'"
              class="filter-tag"
              @click="selectedStatus = 'allocated'"
            >
              已分配
            </el-tag>
          </div>
        </div>
      </div>

      <!-- GPU型号分类列表 -->
      <div v-for="category in gpuCategories" :key="`${category.model}-${category.memory}`" class="gpu-category">
        <div class="category-header">
          <h3>{{ category.model }} / {{ category.memory }}GB</h3>
          <div class="category-stats">
            <el-tag type="success">未分配: {{ category.available.length }}</el-tag>
            <el-tag type="warning" style="margin-left: 8px">已分配: {{ category.allocated.length }}</el-tag>
          </div>
        </div>

        <div class="category-content">
          <!-- 左栏：未分配 -->
          <div class="machines-column">
            <div class="column-header">
              <h4>未分配机器</h4>
            </div>
            <div class="machines-list">
              <div
                v-for="machine in category.available"
                :key="machine.id"
                class="machine-card"
              >
                <div class="machine-info">
                  <div class="machine-name">{{ machine.name }}</div>
                  <div class="machine-specs">
                    <span>{{ machine.region }}</span>
                    <span>CPU: {{ machine.cpu }}</span>
                    <span>内存: {{ machine.memory }}GB</span>
                    <span>硬盘: {{ machine.disk }}GB</span>
                  </div>
                  <div class="machine-driver">
                    <span>CUDA: {{ machine.cudaVersion }}</span>
                    <span>驱动: {{ machine.gpuDriver }}</span>
                  </div>
                  <div class="login-info">
                    <div class="login-item">
                      <span class="login-label">IP地址：</span>
                      <span class="login-value">{{ machine.loginInfo.ip }}:{{ machine.loginInfo.port }}</span>
                    </div>
                    <div class="login-item">
                      <span class="login-label">用户名：</span>
                      <span class="login-value">{{ machine.loginInfo.username }}</span>
                    </div>
                    <div class="login-item">
                      <span class="login-label">密码：</span>
                      <span class="login-value">{{ machine.loginInfo.password }}</span>
                    </div>
                  </div>
                </div>
                <div class="machine-actions">
                  <el-button type="primary" @click="handleAllocate(machine)">
                    立即分配
                  </el-button>
                  <el-button type="success" @click="handleQuickLogin(machine)">
                    一键登录
                  </el-button>
                </div>
              </div>
              <el-empty
                v-if="category.available.length === 0"
                description="暂无未分配机器"
                :image-size="80"
              />
            </div>
          </div>

          <!-- 右栏：已分配 -->
          <div class="machines-column">
            <div class="column-header">
              <h4>已分配机器</h4>
            </div>
            <div class="machines-list">
              <div
                v-for="machine in category.allocated"
                :key="machine.id"
                class="machine-card allocated"
              >
                <div class="machine-info">
                  <div class="machine-name">{{ machine.name }}</div>
                  <div class="machine-specs">
                    <span>{{ machine.region }}</span>
                    <span>CPU: {{ machine.cpu }}</span>
                    <span>内存: {{ machine.memory }}GB</span>
                    <span>硬盘: {{ machine.disk }}GB</span>
                  </div>
                  <div class="machine-driver">
                    <span>CUDA: {{ machine.cudaVersion }}</span>
                    <span>驱动: {{ machine.gpuDriver }}</span>
                  </div>
                  <div class="login-info">
                    <div class="login-item">
                      <span class="login-label">IP地址：</span>
                      <span class="login-value">{{ machine.loginInfo.ip }}:{{ machine.loginInfo.port }}</span>
                    </div>
                    <div class="login-item">
                      <span class="login-label">用户名：</span>
                      <span class="login-value">{{ machine.loginInfo.username }}</span>
                    </div>
                    <div class="login-item">
                      <span class="login-label">密码：</span>
                      <span class="login-value">{{ machine.loginInfo.password }}</span>
                    </div>
                  </div>
                  <div class="allocation-info">
                    <el-tag type="info" size="small">
                      客户: {{ machine.allocatedTo?.customerName }}
                    </el-tag>
                    <el-tag type="warning" size="small" style="margin-left: 4px">
                      时长: {{ machine.allocatedTo?.duration }}个月
                    </el-tag>
                    <span class="allocation-date">
                      分配时间: {{ machine.allocatedTo?.allocatedAt }}
                    </span>
                  </div>
                </div>
                <div class="machine-actions">
                  <el-button type="danger" @click="handleReclaim(machine)">
                    立即回收
                  </el-button>
                  <el-button type="success" @click="handleQuickLogin(machine)">
                    一键登录
                  </el-button>
                </div>
              </div>
              <el-empty
                v-if="category.allocated.length === 0"
                description="暂无已分配机器"
                :image-size="80"
              />
            </div>
          </div>
        </div>
      </div>

      <el-empty
        v-if="gpuCategories.length === 0"
        description="暂无机器数据"
      />
    </div>

    <!-- 分配对话框 -->
    <el-dialog
      v-model="allocationDialogVisible"
      title="分配主机资源"
      width="600px"
    >
      <div v-if="currentMachine" class="allocation-dialog">
        <!-- 机器信息 -->
        <div class="device-info-section">
          <h4>机器信息</h4>
          <div class="info-row">
            <span class="info-label">机器名称：</span>
            <span class="info-value">{{ currentMachine.name }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">GPU型号：</span>
            <span class="info-value">{{ currentMachine.gpuModel }} / {{ currentMachine.gpuMemory }}GB</span>
          </div>
          <div class="info-row">
            <span class="info-label">配置：</span>
            <span class="info-value">{{ currentMachine.cpu }} / {{ currentMachine.memory }}GB内存</span>
          </div>
        </div>

        <!-- 分配表单 -->
        <el-form :model="allocationForm" label-width="100px" style="margin-top: 20px">
          <el-form-item label="选择客户" required>
            <el-select
              v-model="allocationForm.customerId"
              placeholder="请选择客户"
              style="width: 100%"
            >
              <el-option
                v-for="customer in customers"
                :key="customer.id"
                :label="`${customer.name} (${customer.company})`"
                :value="customer.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="分配时长" required>
            <el-input-number
              v-model="allocationForm.duration"
              :min="1"
              :max="12"
              style="width: 100%"
            />
            <span style="margin-left: 8px; color: #909399">个月</span>
          </el-form-item>

          <el-form-item label="备注">
            <el-input
              v-model="allocationForm.notes"
              type="textarea"
              :rows="3"
              placeholder="请输入备注信息"
            />
          </el-form-item>
        </el-form>
      </div>

      <template #footer>
        <el-button @click="allocationDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitAllocation">确认分配</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.computing-market {
  padding: 24px;
}

.market-container {
  background: white;
  border-radius: 8px;
  padding: 24px;
  min-height: 400px;
}

/* 筛选区域 */
.filter-section {
  margin-bottom: 24px;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;
}

.filter-group {
  display: flex;
  align-items: flex-start;
  margin-bottom: 16px;
}

.filter-group:last-child {
  margin-bottom: 0;
}

.filter-label {
  min-width: 100px;
  font-size: 14px;
  font-weight: 500;
  color: #606266;
  padding-top: 6px;
  flex-shrink: 0;
}

.filter-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  flex: 1;
}

.filter-tag {
  cursor: pointer;
  transition: all 0.3s;
  user-select: none;
}

.filter-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

/* GPU分类 */
.gpu-category {
  margin-bottom: 32px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
}

.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: #f5f7fa;
  border-bottom: 1px solid #e4e7ed;
}

.category-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.category-stats {
  display: flex;
  gap: 8px;
}

/* 左右两栏布局 */
.category-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1px;
  background: #e4e7ed;
}

.machines-column {
  background: white;
  min-height: 300px;
}

.column-header {
  padding: 12px 16px;
  background: #fafafa;
  border-bottom: 1px solid #e4e7ed;
}

.column-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #606266;
}

.machines-list {
  padding: 16px;
}

/* 机器卡片 */
.machine-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  margin-bottom: 12px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  transition: all 0.3s;
}

.machine-card:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.machine-card:last-child {
  margin-bottom: 0;
}

.machine-card.allocated {
  background: #fff9e6;
  border-color: #ffd666;
}

.machine-info {
  flex: 1;
}

.machine-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.machine-specs {
  display: flex;
  gap: 12px;
  font-size: 13px;
  color: #606266;
  margin-bottom: 6px;
}

.machine-specs span {
  display: inline-block;
}

.machine-driver {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #909399;
}

.login-info {
  margin-top: 8px;
  padding: 8px;
  background: #f0f9ff;
  border-radius: 4px;
  border: 1px solid #d1e7ff;
}

.login-item {
  display: flex;
  align-items: center;
  font-size: 12px;
  margin-bottom: 4px;
}

.login-item:last-child {
  margin-bottom: 0;
}

.login-label {
  color: #606266;
  min-width: 60px;
  font-weight: 500;
}

.login-value {
  color: #303133;
  font-family: 'Courier New', monospace;
  font-weight: 600;
}

.allocation-info {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}

.allocation-date {
  font-size: 12px;
  color: #909399;
  margin-left: 8px;
}

.machine-actions {
  margin-left: 16px;
}

/* 分配对话框 */
.allocation-dialog {
  padding: 10px 0;
}

.device-info-section {
  background: #f5f7fa;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.device-info-section h4 {
  margin: 0 0 12px 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.info-row {
  display: flex;
  margin-bottom: 8px;
  font-size: 14px;
}

.info-row:last-child {
  margin-bottom: 0;
}

.info-label {
  color: #606266;
  min-width: 100px;
}

.info-value {
  color: #303133;
  font-weight: 500;
}

/* 响应式 */
@media (max-width: 1200px) {
  .category-content {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .machine-card {
    flex-direction: column;
    align-items: flex-start;
  }

  .machine-actions {
    margin-left: 0;
    margin-top: 12px;
    width: 100%;
  }

  .machine-actions .el-button {
    width: 100%;
  }
}
</style>
