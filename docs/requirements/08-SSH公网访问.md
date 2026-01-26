# SSH 公网访问

> 所属模块：模块 5 - 网络与访问模块
>
> 功能编号：5.1
>
> 优先级：P0（必须）

---

## 1. 功能概述

### 1.1 功能描述

SSH 公网访问功能允许用户通过公网 IP 和动态分配的端口访问开发环境，支持密码认证和密钥认证，提供安全的远程终端访问能力。

### 1.2 业务价值

- ✅ 用户可从任何地方访问开发环境
- ✅ 支持多种 SSH 客户端（PuTTY、Xshell、Terminal）
- ✅ 安全的认证机制
- ✅ 会话管理和审计

---

## 2. 核心功能

### 2.1 SSH 访问架构

```
用户 → 公网网关 (gateway.example.com:30001)
     → 内网主机 (192.168.1.10)
     → Docker 容器 (172.17.0.2:22)
```

### 2.2 认证方式

**支持的认证方式：**
1. **密码认证**：系统生成随机密码
2. **密钥认证**：用户上传公钥

---

## 3. 实现方案

### 3.1 端口转发

```bash
# 使用 iptables 实现端口转发
iptables -t nat -A PREROUTING -p tcp --dport 30001 \
  -j DNAT --to-destination 172.17.0.2:22
```

### 3.2 密码生成

```go
func GenerateSSHPassword() string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%"
    password := make([]byte, 16)
    for i := range password {
        password[i] = charset[rand.Intn(len(charset))]
    }
    return string(password)
}
```

---

## 4. API 接口

```go
// 获取 SSH 连接信息
GET /api/environments/:id/ssh

Response: {
  "ssh_host": "gateway.example.com",
  "ssh_port": 30001,
  "ssh_username": "developer",
  "ssh_password": "Abc123!@#xyz",
  "ssh_command": "ssh developer@gateway.example.com -p 30001"
}
```

---

## 5. 前端界面

```vue
<template>
  <el-card>
    <h3>SSH 访问信息</h3>
    <el-descriptions :column="1" border>
      <el-descriptions-item label="主机">
        {{ sshInfo.ssh_host }}
      </el-descriptions-item>
      <el-descriptions-item label="端口">
        {{ sshInfo.ssh_port }}
      </el-descriptions-item>
      <el-descriptions-item label="用户名">
        {{ sshInfo.ssh_username }}
      </el-descriptions-item>
      <el-descriptions-item label="密码">
        {{ sshInfo.ssh_password }}
        <el-button @click="copyPassword" size="small">复制</el-button>
      </el-descriptions-item>
    </el-descriptions>

    <el-divider />

    <h4>连接命令</h4>
    <el-input
      v-model="sshInfo.ssh_command"
      readonly
    >
      <template #append>
        <el-button @click="copyCommand">复制</el-button>
      </template>
    </el-input>
  </el-card>
</template>
```

---

## 6. 测试用例

| 用例 | 场景 | 预期结果 |
|------|------|---------|
| TC-01 | 使用密码连接 | 连接成功 |
| TC-02 | 使用密钥连接 | 连接成功 |
| TC-03 | 错误密码 | 连接失败 |

---

**文档版本：** v1.0
**创建日期：** 2026-01-26
