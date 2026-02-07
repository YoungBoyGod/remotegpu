#!/bin/bash
set -e

# 配置 SSH
mkdir -p /run/sshd
echo "root:${SSH_PASSWORD:-remotegpu123}" | chpasswd
sed -i 's/#PermitRootLogin.*/PermitRootLogin yes/' /etc/ssh/sshd_config
sed -i 's/#PasswordAuthentication.*/PasswordAuthentication yes/' /etc/ssh/sshd_config

# 启动 SSH 服务
/usr/sbin/sshd

# 模拟 nvidia-smi（用于测试 GPU 监控）
if [ "${FAKE_GPU:-false}" = "true" ]; then
    cat > /usr/local/bin/nvidia-smi <<'FAKEGPU'
#!/bin/bash
cat <<EOF
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 535.129.03   Driver Version: 535.129.03   CUDA Version: 12.2     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|===============================+======================+======================|
|   0  Tesla V100-SXM2...  On   | 00000000:00:1E.0 Off |                    0 |
| N/A   32C    P0    25W / 300W |      0MiB / 16384MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
EOF
FAKEGPU
    chmod +x /usr/local/bin/nvidia-smi
fi

# 启动 Jupyter（如果启用）
if [ "${JUPYTER_ENABLED:-true}" = "true" ]; then
    jupyter lab \
        --ip=0.0.0.0 \
        --port=${JUPYTER_PORT:-8888} \
        --no-browser \
        --allow-root \
        --NotebookApp.token="${JUPYTER_TOKEN:-remotegpu}" \
        --notebook-dir=/workspace &
fi

# 保持容器运行
exec tail -f /dev/null
