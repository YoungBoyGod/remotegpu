#!/bin/bash
set -e

# 配置 SSH
mkdir -p /run/sshd
echo "root:${SSH_PASSWORD:-remotegpu123}" | chpasswd
sed -i 's/#PermitRootLogin.*/PermitRootLogin yes/' /etc/ssh/sshd_config
sed -i 's/#PasswordAuthentication.*/PasswordAuthentication yes/' /etc/ssh/sshd_config

# 启动 SSH 服务
/usr/sbin/sshd

# 启动 Jupyter Lab
exec jupyter lab \
  --ip=0.0.0.0 \
  --port=${JUPYTER_PORT:-8888} \
  --no-browser \
  --allow-root \
  --NotebookApp.token="${JUPYTER_TOKEN:-remotegpu}" \
  --notebook-dir=/workspace
