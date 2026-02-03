# JupyterLab 集成

> 所属模块：模块 3 - 环境管理模块
>
> 功能编号：3.3
>
> 优先级：P1（重要）

---

## 1. 功能概述

### 1.1 功能描述

JupyterLab 集成功能为用户提供基于 Web 的交互式开发环境，支持 Notebook 编辑、代码执行、可视化展示，并与环境管理模块深度集成，实现一键启动和访问。

### 1.2 业务价值

- ✅ 提供友好的交互式开发体验
- ✅ 支持多种编程语言（Python、R、Julia）
- ✅ 内置可视化和数据分析工具
- ✅ 无需本地安装，浏览器即可访问

---

## 2. 核心功能

### 2.1 JupyterLab 启动

**启动方式：**
1. **环境创建时自动启动**：创建环境时选择启动 JupyterLab
2. **手动启动**：在已有环境中启动 JupyterLab
3. **自定义配置**：配置端口、密码、插件等

**启动流程：**
```yaml
启动流程:
  1. 检查环境状态（必须为 running）
  2. 分配访问端口（从端口池）
  3. 生成访问 Token
  4. 在容器中启动 JupyterLab 服务
  5. 配置反向代理（Nginx/Traefik）
  6. 返回访问 URL
```

### 2.2 访问控制

**认证方式：**
- Token 认证（默认）
- 密码认证
- OAuth 集成（企业版）

**访问 URL 格式：**
```
https://jupyter.example.com/env-{environment_id}/?token={access_token}
```

---

## 3. 数据模型

```sql
CREATE TABLE jupyter_instances (
    id BIGSERIAL PRIMARY KEY,
    environment_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,

    -- 访问信息
    access_url VARCHAR(512),
    access_token VARCHAR(128),
    external_port INT,
    internal_port INT DEFAULT 8888,

    -- 配置
    config JSONB,

    -- 状态
    status VARCHAR(20) DEFAULT 'starting',

    started_at TIMESTAMP,
    stopped_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (environment_id) REFERENCES environments(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);
```

---

## 4. 启动实现

### 4.1 JupyterLab 启动器

```go
// JupyterLab 启动器
type JupyterLabLauncher struct {
    dockerClient *docker.Client
    portManager  *PortManager
}

// 启动 JupyterLab
func (l *JupyterLabLauncher) Start(envID string, config JupyterConfig) (*JupyterInstance, error) {
    // 1. 获取环境信息
    env, err := l.getEnvironment(envID)
    if err != nil {
        return nil, err
    }

    if env.Status != "running" {
        return nil, errors.New("环境未运行")
    }

    // 2. 分配端口
    externalPort, err := l.portManager.AllocatePort("jupyter")
    if err != nil {
        return nil, err
    }

    // 3. 生成访问 Token
    accessToken := generateSecureToken(32)

    // 4. 在容器中启动 JupyterLab
    cmd := fmt.Sprintf(
        "jupyter lab --ip=0.0.0.0 --port=8888 --no-browser --allow-root "+
        "--NotebookApp.token='%s' --NotebookApp.base_url='/env-%s'",
        accessToken, envID,
    )

    execConfig := docker.CreateExecOptions{
        Container: env.ContainerID,
        Cmd:       []string{"/bin/bash", "-c", cmd},
        AttachStdout: true,
        AttachStderr: true,
    }

    exec, err := l.dockerClient.CreateExec(execConfig)
    if err != nil {
        return nil, err
    }

    if err := l.dockerClient.StartExec(exec.ID, docker.StartExecOptions{}); err != nil {
        return nil, err
    }

    // 5. 配置端口映射（通过 iptables 或反向代理）
    if err := l.setupPortForwarding(env.HostID, externalPort, env.ContainerID, 8888); err != nil {
        return nil, err
    }

    // 6. 创建实例记录
    instance := &JupyterInstance{
        EnvironmentID: envID,
        CustomerID:    env.CustomerID,
        AccessURL:     fmt.Sprintf("https://jupyter.example.com/env-%s/?token=%s", envID, accessToken),
        AccessToken:   accessToken,
        ExternalPort:  externalPort,
        InternalPort:  8888,
        Status:        "running",
        StartedAt:     time.Now(),
    }

    if err := l.db.Create(instance).Error; err != nil {
        return nil, err
    }

    return instance, nil
}
```

### 4.2 健康检查

```go
// 健康检查
func (l *JupyterLabLauncher) HealthCheck(instanceID int64) error {
    var instance JupyterInstance
    if err := l.db.First(&instance, instanceID).Error; err != nil {
        return err
    }

    // 检查端口是否可访问
    url := fmt.Sprintf("http://localhost:%d/api", instance.ExternalPort)
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return errors.New("JupyterLab 服务异常")
    }

    return nil
}
```

---

## 5. API 接口

### 5.1 启动 JupyterLab

```go
POST /api/environments/:id/jupyter/start
Body: {
  "password": "optional_password",
  "plugins": ["jupyterlab-git", "jupyterlab-lsp"]
}

Response: {
  "instance_id": 123,
  "access_url": "https://jupyter.example.com/env-abc123/?token=xyz789",
  "status": "running"
}
```

### 5.2 停止 JupyterLab

```go
POST /api/environments/:id/jupyter/stop

Response: {
  "status": "stopped"
}
```

### 5.3 查询 JupyterLab 状态

```go
GET /api/environments/:id/jupyter

Response: {
  "instance_id": 123,
  "status": "running",
  "access_url": "https://jupyter.example.com/env-abc123/?token=xyz789",
  "started_at": "2026-01-26T10:00:00Z"
}
```

---

## 6. 前端界面

```vue
<template>
  <el-card>
    <h3>JupyterLab</h3>

    <div v-if="!jupyterInstance">
      <el-button type="primary" @click="startJupyter" :loading="starting">
        启动 JupyterLab
      </el-button>
    </div>

    <div v-else>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(jupyterInstance.status)">
            {{ jupyterInstance.status }}
          </el-tag>
        </el-descriptions-item>

        <el-descriptions-item label="访问地址">
          <el-link :href="jupyterInstance.access_url" target="_blank" type="primary">
            {{ jupyterInstance.access_url }}
          </el-link>
        </el-descriptions-item>

        <el-descriptions-item label="启动时间">
          {{ formatTime(jupyterInstance.started_at) }}
        </el-descriptions-item>
      </el-descriptions>

      <el-button-group style="margin-top: 20px">
        <el-button type="primary" @click="openJupyter">
          打开 JupyterLab
        </el-button>
        <el-button @click="stopJupyter">
          停止
        </el-button>
      </el-button-group>
    </div>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  environmentId: String
})

const jupyterInstance = ref(null)
const starting = ref(false)

const startJupyter = async () => {
  starting.value = true
  try {
    const response = await fetch(`/api/environments/${props.environmentId}/jupyter/start`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' }
    })

    const data = await response.json()
    jupyterInstance.value = data

    ElMessage.success('JupyterLab 启动成功')
  } catch (error) {
    ElMessage.error('启动失败: ' + error.message)
  } finally {
    starting.value = false
  }
}

const openJupyter = () => {
  window.open(jupyterInstance.value.access_url, '_blank')
}

const stopJupyter = async () => {
  try {
    await fetch(`/api/environments/${props.environmentId}/jupyter/stop`, {
      method: 'POST'
    })

    jupyterInstance.value = null
    ElMessage.success('JupyterLab 已停止')
  } catch (error) {
    ElMessage.error('停止失败: ' + error.message)
  }
}

onMounted(async () => {
  // 查询是否已有运行中的 JupyterLab
  try {
    const response = await fetch(`/api/environments/${props.environmentId}/jupyter`)
    const data = await response.json()

    if (data.status === 'running') {
      jupyterInstance.value = data
    }
  } catch (error) {
    // 未启动
  }
})
</script>
```

---

## 7. 反向代理配置

### 7.1 Nginx 配置

```nginx
# JupyterLab 反向代理
location ~ ^/env-([a-zA-Z0-9-]+)/ {
    set $env_id $1;

    # 动态查询后端端口
    set $backend_port 0;
    access_by_lua_block {
        local redis = require "resty.redis"
        local red = redis:new()
        red:connect("127.0.0.1", 6379)

        local port = red:get("jupyter:env:" .. ngx.var.env_id)
        if port then
            ngx.var.backend_port = port
        end
    }

    proxy_pass http://127.0.0.1:$backend_port;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

    # WebSocket 支持
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
}
```

---

## 8. 测试用例

| 用例 | 场景 | 预期结果 |
|------|------|---------|
| TC-01 | 启动 JupyterLab | 启动成功，返回访问 URL |
| TC-02 | 访问 JupyterLab | 浏览器打开，可正常使用 |
| TC-03 | 停止 JupyterLab | 停止成功，端口释放 |
| TC-04 | 环境停止时自动停止 JupyterLab | JupyterLab 自动停止 |

---

**文档版本：** v1.0
**创建日期：** 2026-01-26
