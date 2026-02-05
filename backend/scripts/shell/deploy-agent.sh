#!/bin/bash
# Agent 部署脚本
# 使用方法: ./deploy-agent.sh <target_ip> [ssh_user]

TARGET_IP="${1:-}"
SSH_USER="${2:-root}"
AGENT_BIN="/home/luo/code/remotegpu/agent/remotegpu-agent"

if [ -z "$TARGET_IP" ]; then
  echo "用法: $0 <target_ip> [ssh_user]"
  echo "示例: $0 192.168.1.100 root"
  exit 1
fi

echo "=== 部署 Agent 到 $TARGET_IP ==="

# 1. 复制二进制文件
echo "1. 复制二进制文件..."
scp "$AGENT_BIN" "${SSH_USER}@${TARGET_IP}:/usr/local/bin/"

# 2. 创建 systemd 服务
echo "2. 创建 systemd 服务..."
ssh "${SSH_USER}@${TARGET_IP}" 'cat > /etc/systemd/system/remotegpu-agent.service << EOF
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
EOF'

# 3. 启动服务
echo "3. 启动服务..."
ssh "${SSH_USER}@${TARGET_IP}" 'systemctl daemon-reload && systemctl enable remotegpu-agent && systemctl restart remotegpu-agent'

# 4. 验证
echo "4. 验证..."
sleep 2
curl -s "http://${TARGET_IP}:8090/api/v1/ping"
echo ""

echo "=== 部署完成 ==="
