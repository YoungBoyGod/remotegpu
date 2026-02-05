# RemoteGPU Agent 部署指南

## 概述

Agent 是运行在 GPU 机器上的服务，负责与后端通信。

## 文件位置

- 二进制文件: `/home/luo/code/remotegpu/agent/remotegpu-agent`
- 源代码: `/home/luo/code/remotegpu/agent/cmd/`

## 部署步骤

### 1. 复制二进制文件到目标机器

```bash
scp /home/luo/code/remotegpu/agent/remotegpu-agent root@192.168.1.100:/usr/local/bin/
```

### 2. 在目标机器上启动

```bash
# 设置端口 (默认 8090)
export AGENT_PORT=8090

# 启动
/usr/local/bin/remotegpu-agent
```

### 3. 使用 systemd 管理

```bash
# 创建 service 文件
cat > /etc/systemd/system/remotegpu-agent.service << 'EOF'
[Unit]
Description=RemoteGPU Agent
After=network.target

[Service]
Type=simple
Environment=AGENT_PORT=8090
ExecStart=/usr/local/bin/remotegpu-agent
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
systemctl daemon-reload
systemctl enable remotegpu-agent
systemctl start remotegpu-agent
```

## API 端点

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/v1/ping | 健康检查 |
| GET | /api/v1/system/info | 系统信息 |
| POST | /api/v1/process/stop | 停止进程 |
| POST | /api/v1/ssh/reset | 重置 SSH |
| POST | /api/v1/machine/cleanup | 清理机器 |

## 验证

```bash
curl http://192.168.1.100:8090/api/v1/ping
# 返回: {"ok":true}
```
