# API 使用示例

本文档提供各个模块 API 的使用示例。

## 1. 用户与权限模块

### 用户登录
```typescript
import { authApi } from '@/api'

// 登录
const handleLogin = async () => {
  try {
    const response = await authApi.login({
      username: 'admin',
      password: '123456'
    })

    // 保存 token
    localStorage.setItem('token', response.token)
    console.log('登录成功', response.user)
  } catch (error) {
    console.error('登录失败', error)
  }
}
```

### 获取当前用户信息
```typescript
import { authApi } from '@/api'

const getUserInfo = async () => {
  const user = await authApi.getCurrentUser()
  console.log('用户信息', user)
}
```

### 工作空间管理
```typescript
import { authApi } from '@/api'

// 创建工作空间
const createWorkspace = async () => {
  const result = await authApi.createWorkspace({
    name: '我的工作空间',
    description: '这是一个测试工作空间',
    type: 'team'
  })
  console.log('工作空间 ID', result.id)
}

// 获取工作空间列表
const getWorkspaces = async () => {
  const response = await authApi.getWorkspaceList({
    page: 1,
    page_size: 10
  })
  console.log('工作空间列表', response.items)
}
```

## 2. 环境管理模块

### 创建开发环境
```typescript
import { environmentApi } from '@/api'

const createEnv = async () => {
  const result = await environmentApi.createEnvironment({
    name: '开发环境-1',
    description: 'PyTorch 开发环境',
    image: 'pytorch:2.0-cuda11.8',
    resources: {
      cpu: 4,
      memory: 8192,
      gpu: 1,
      storage: 50
    },
    datasets: [1, 2], // 挂载数据集 ID
    env_vars: {
      PYTHONPATH: '/workspace'
    }
  })
  console.log('环境创建成功', result)
}
```

### 启动/停止环境
```typescript
import { environmentApi } from '@/api'

// 启动环境
const startEnv = async (envId: number) => {
  await environmentApi.startEnvironment(envId)
  console.log('环境启动中...')
}

// 停止环境
const stopEnv = async (envId: number) => {
  await environmentApi.stopEnvironment(envId)
  console.log('环境停止中...')
}
```

### 获取访问信息
```typescript
import { environmentApi } from '@/api'

const getAccess = async (envId: number) => {
  const access = await environmentApi.getEnvironmentAccess(envId)
  console.log('SSH 访问:', `ssh ${access.ssh_username}@${access.ssh_host} -p ${access.ssh_port}`)
  console.log('密码:', access.ssh_password)
}
```

## 3. 数据与存储模块

### 创建数据集
```typescript
import { storageApi } from '@/api'

const createDataset = async () => {
  const result = await storageApi.createDataset({
    name: 'ImageNet-1K',
    description: 'ImageNet 数据集',
    visibility: 'workspace',
    tags: ['image', 'classification']
  })
  console.log('数据集 ID', result.id)
  console.log('存储路径', result.storage_path)
}
```

### 上传文件
```typescript
import { storageApi } from '@/api'

const uploadFile = async (datasetId: number, file: File) => {
  // 1. 获取上传凭证
  const { upload_url } = await storageApi.getUploadUrl(datasetId, {
    file_name: file.name,
    file_size: file.size
  })

  // 2. 上传文件到对象存储
  await fetch(upload_url, {
    method: 'PUT',
    body: file
  })

  // 3. 完成上传
  await storageApi.completeUpload(datasetId, {
    files: [{
      file_name: file.name,
      file_size: file.size
    }]
  })

  console.log('文件上传成功')
}
```

## 4. 训练与推理模块

### 创建训练任务
```typescript
import { trainingApi } from '@/api'

const createTrainingJob = async () => {
  const result = await trainingApi.createTrainingJob({
    name: '图像分类训练',
    image: 'pytorch:2.0-cuda11.8',
    script: 'python train.py --epochs 100',
    datasets: [1], // 数据集 ID
    resources: {
      cpu: 8,
      memory: 16384,
      gpu: 2
    },
    env_vars: {
      BATCH_SIZE: '32',
      LEARNING_RATE: '0.001'
    }
  })
  console.log('训练任务创建成功', result)
}
```

### 查看训练日志
```typescript
import { trainingApi } from '@/api'

const viewLogs = async (jobId: number) => {
  const { logs } = await trainingApi.getTrainingLogs(jobId, {
    tail: 100 // 最后 100 行
  })
  console.log(logs)
}
```

## 5. 监控告警模块

### 获取 GPU 监控数据
```typescript
import { monitoringApi } from '@/api'

const getGPUMetrics = async (gpuId: number) => {
  const { metrics } = await monitoringApi.getGPUMetrics(gpuId, {
    start_time: '2024-01-01T00:00:00Z',
    end_time: '2024-01-02T00:00:00Z'
  })

  // 绘制图表
  metrics.forEach(m => {
    console.log(`GPU 使用率: ${m.gpu_usage}%, 温度: ${m.temperature}°C`)
  })
}
```

### 创建告警规则
```typescript
import { monitoringApi } from '@/api'

const createAlert = async () => {
  const result = await monitoringApi.createAlertRule({
    name: 'GPU 温度过高',
    metric: 'gpu_temperature',
    threshold: 85,
    comparison: 'gt',
    severity: 'warning',
    notification_channels: ['email', 'webhook']
  })
  console.log('告警规则创建成功', result)
}
```

## 6. 计费管理模块

### 查询账户余额
```typescript
import { billingApi } from '@/api'

const checkBalance = async () => {
  const account = await billingApi.getAccountBalance()
  console.log(`当前余额: ${account.balance} 元`)
  console.log(`信用额度: ${account.credit_limit} 元`)
}
```

### 查询计费记录
```typescript
import { billingApi } from '@/api'

const getBillingRecords = async () => {
  const response = await billingApi.getBillingRecords({
    page: 1,
    page_size: 20,
    start_date: '2024-01-01',
    end_date: '2024-01-31',
    resource_type: 'gpu'
  })

  console.log('计费记录', response.items)
  console.log('总金额', response.items.reduce((sum, r) => sum + r.amount, 0))
}
```

## 7. 完整示例：创建环境并训练模型

```typescript
import { environmentApi, storageApi, trainingApi } from '@/api'

const fullWorkflow = async () => {
  try {
    // 1. 创建数据集
    console.log('步骤 1: 创建数据集')
    const dataset = await storageApi.createDataset({
      name: 'MNIST',
      description: '手写数字数据集',
      visibility: 'private'
    })

    // 2. 上传数据（省略文件上传代码）
    console.log('步骤 2: 上传数据')
    // ... 上传文件

    // 3. 创建开发环境
    console.log('步骤 3: 创建开发环境')
    const env = await environmentApi.createEnvironment({
      name: 'MNIST 训练环境',
      image: 'pytorch:2.0-cuda11.8',
      resources: {
        cpu: 4,
        memory: 8192,
        gpu: 1
      },
      datasets: [dataset.id]
    })

    // 4. 等待环境启动
    console.log('步骤 4: 等待环境启动')
    await new Promise(resolve => setTimeout(resolve, 30000))

    // 5. 创建训练任务
    console.log('步骤 5: 创建训练任务')
    const job = await trainingApi.createTrainingJob({
      name: 'MNIST 分类训练',
      image: 'pytorch:2.0-cuda11.8',
      script: 'python train.py',
      datasets: [dataset.id],
      resources: {
        cpu: 4,
        memory: 8192,
        gpu: 1
      }
    })

    console.log('训练任务已提交', job)

  } catch (error) {
    console.error('工作流执行失败', error)
  }
}
```

## 8. 错误处理

```typescript
import { environmentApi } from '@/api'
import { ElMessage } from 'element-plus'

const createEnvWithErrorHandling = async () => {
  try {
    const result = await environmentApi.createEnvironment({
      name: '测试环境',
      image: 'ubuntu:20.04',
      resources: {
        cpu: 2,
        memory: 4096
      }
    })

    ElMessage.success('环境创建成功')
    return result

  } catch (error: any) {
    // 错误已经在 request.ts 中统一处理
    // 这里可以做额外的业务逻辑处理
    console.error('创建环境失败', error)
    throw error
  }
}
```

## 9. 使用 TypeScript 类型

```typescript
import type { Environment, CreateEnvironmentRequest } from '@/api/environment/types'
import { environmentApi } from '@/api'

// 使用类型定义
const envData: CreateEnvironmentRequest = {
  name: '类型安全的环境',
  image: 'pytorch:2.0',
  resources: {
    cpu: 4,
    memory: 8192,
    gpu: 1
  }
}

const env: Environment = await environmentApi.createEnvironment(envData)
```
