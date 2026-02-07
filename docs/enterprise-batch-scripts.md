# 企业方案批量配置脚本

## iptables端口转发批量配置

**执行位置**：企业网关服务器

创建脚本 `setup_port_forwarding.sh`：

```bash
#!/bin/bash

# 企业公网IP（自动获取或手动指定）
PUBLIC_IP=$(curl -s ifconfig.me)

# GPU内网IP范围（根据实际情况修改）
GPU_IP_PREFIX="192.168.10"
GPU_START=101  # GPU1的内网IP最后一位

for i in {1..200}; do
    GPU_NUM=$i
    GPU_IP="${GPU_IP_PREFIX}.$((GPU_START + i - 1))"

    # SSH (22)
    SSH_PORT=$((2200 + GPU_NUM))
    iptables -t nat -A PREROUTING -p tcp --dport $SSH_PORT -j DNAT --to-destination ${GPU_IP}:22
    iptables -t nat -A POSTROUTING -d ${GPU_IP} -p tcp --dport 22 -j MASQUERADE

    # Jupyter (8888)
    JUPYTER_PORT=$((8000 + GPU_NUM))
    iptables -t nat -A PREROUTING -p tcp --dport $JUPYTER_PORT -j DNAT --to-destination ${GPU_IP}:8888
    iptables -t nat -A POSTROUTING -d ${GPU_IP} -p tcp --dport 8888 -j MASQUERADE

    # TensorBoard (6006)
    TB_PORT=$((9000 + GPU_NUM))
    iptables -t nat -A PREROUTING -p tcp --dport $TB_PORT -j DNAT --to-destination ${GPU_IP}:6006
    iptables -t nat -A POSTROUTING -d ${GPU_IP} -p tcp --dport 6006 -j MASQUERADE

    echo "配置完成: GPU${GPU_NUM} (${GPU_IP})"
done

# 保存规则
iptables-save > /etc/iptables/rules.v4
echo "所有端口转发规则已配置并保存"
```

使用方法：
```bash
chmod +x setup_port_forwarding.sh
sudo ./setup_port_forwarding.sh
```

---

## nginx配置批量生成

**执行位置**：云服务器

创建脚本 `generate_nginx_enterprise.sh`：

```bash
#!/bin/bash

ENTERPRISE_IP="企业公网IP"
DOMAIN="gpu.domain.com"
SSL_CERT="/etc/letsencrypt/live/gpu.domain.com/fullchain.pem"
SSL_KEY="/etc/letsencrypt/live/gpu.domain.com/privkey.pem"

# 生成Web服务配置
for i in {1..200}; do
    JUPYTER_PORT=$((8000 + i))
    TB_PORT=$((9000 + i))

    cat >> gpu_enterprise.conf <<EOF
# GPU${i} Jupyter
server {
    listen 443 ssl http2;
    server_name gpu${i}-jupyter.${DOMAIN};
    ssl_certificate ${SSL_CERT};
    ssl_certificate_key ${SSL_KEY};
    location / {
        proxy_pass http://${ENTERPRISE_IP}:${JUPYTER_PORT};
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
    }
}

# GPU${i} TensorBoard
server {
    listen 443 ssl http2;
    server_name gpu${i}-tensorboard.${DOMAIN};
    ssl_certificate ${SSL_CERT};
    ssl_certificate_key ${SSL_KEY};
    location / {
        proxy_pass http://${ENTERPRISE_IP}:${TB_PORT};
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
    }
}

EOF
done

echo "配置文件已生成: gpu_enterprise.conf"
```
