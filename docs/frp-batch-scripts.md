# frp方案 - 批量配置脚本

## 概述

为200台GPU机器生成配置文件的批量脚本。

---

## 1. 生成frpc配置脚本

**执行位置**: 本地电脑或管理机器

**脚本**: `generate_frpc_configs.sh`

```bash
#!/bin/bash

SERVER_IP="云服务器IP"
TOKEN="your_secure_token_here"

for i in {1..200}; do
    GPU_NUM=$i
    SSH_PORT=$((10000 + GPU_NUM))
    JUPYTER_PORT=$((11000 + GPU_NUM))
    TB_PORT=$((12000 + GPU_NUM))
    SERVICE1_PORT=$((13000 + GPU_NUM))
    SERVICE2_PORT=$((14000 + GPU_NUM))

    cat > frpc_gpu${GPU_NUM}.ini <<EOF
[common]
server_addr = ${SERVER_IP}
server_port = 7000
authentication_method = token
token = ${TOKEN}

[gpu${GPU_NUM}-ssh]
type = tcp
local_ip = 127.0.0.1
local_port = 22
remote_port = ${SSH_PORT}

[gpu${GPU_NUM}-jupyter]
type = tcp
local_ip = 127.0.0.1
local_port = 8888
remote_port = ${JUPYTER_PORT}

[gpu${GPU_NUM}-tensorboard]
type = tcp
local_ip = 127.0.0.1
local_port = 6006
remote_port = ${TB_PORT}

[gpu${GPU_NUM}-service1]
type = tcp
local_ip = 127.0.0.1
local_port = 本地端口1
remote_port = ${SERVICE1_PORT}

[gpu${GPU_NUM}-service2]
type = tcp
local_ip = 127.0.0.1
local_port = 本地端口2
remote_port = ${SERVICE2_PORT}
EOF

    echo "生成: frpc_gpu${GPU_NUM}.ini"
done

echo "所有frpc配置文件已生成"
```

**使用方法**:
```bash
chmod +x generate_frpc_configs.sh
./generate_frpc_configs.sh
```

生成的配置文件需要分发到对应的GPU机器。

---

## 2. 生成nginx配置脚本

**执行位置**: 云服务器

**脚本**: `generate_nginx_frp.sh`

```bash
#!/bin/bash

DOMAIN="gpu.domain.com"
SSL_CERT="/etc/letsencrypt/live/gpu.domain.com/fullchain.pem"
SSL_KEY="/etc/letsencrypt/live/gpu.domain.com/privkey.pem"

cat > /tmp/gpu-frp-web.conf <<EOF
# GPU Web服务配置
EOF

for i in {1..200}; do
    JUPYTER_PORT=$((11000 + i))
    TB_PORT=$((12000 + i))

    cat >> /tmp/gpu-frp-web.conf <<EOF

# GPU${i} Jupyter
server {
    listen 443 ssl http2;
    server_name gpu${i}-jupyter.${DOMAIN};
    ssl_certificate ${SSL_CERT};
    ssl_certificate_key ${SSL_KEY};
    location / {
        proxy_pass http://127.0.0.1:${JUPYTER_PORT};
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
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
        proxy_pass http://127.0.0.1:${TB_PORT};
        proxy_set_header Host \$host;
    }
}
EOF

    echo "生成: GPU${i} 配置"
done

echo "配置文件已生成: /tmp/gpu-frp-web.conf"
echo "请复制到: /etc/nginx/sites-available/gpu-frp"
```

**使用方法**:
```bash
chmod +x generate_nginx_frp.sh
./generate_nginx_frp.sh

# 复制配置
sudo cp /tmp/gpu-frp-web.conf /etc/nginx/sites-available/gpu-frp
sudo ln -s /etc/nginx/sites-available/gpu-frp /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 3. 生成SSH配置脚本

**执行位置**: 用户本地电脑

**脚本**: `generate_ssh_config.sh`

```bash
#!/bin/bash

SERVER_IP="云服务器IP"
USERNAME="your_username"

for i in {1..200}; do
    SSH_PORT=$((10000 + i))
    cat >> ~/.ssh/config <<EOF
Host gpu${i}
    HostName ${SERVER_IP}
    Port ${SSH_PORT}
    User ${USERNAME}
    IdentityFile ~/.ssh/id_rsa

EOF
done

echo "SSH配置已添加到 ~/.ssh/config"
```

**使用方法**:
```bash
chmod +x generate_ssh_config.sh
./generate_ssh_config.sh

# 测试
ssh gpu1
ssh gpu2
```

---

## 4. 批量部署frpc脚本

**前提**: 已配置SSH密钥认证

**脚本**: `deploy_frpc.sh`

```bash
#!/bin/bash

for i in {1..200}; do
    GPU_HOST="gpu${i}"
    CONFIG_FILE="frpc_gpu${i}.ini"

    echo "部署到 ${GPU_HOST}..."

    # 复制配置文件
    scp ${CONFIG_FILE} ${GPU_HOST}:/tmp/frpc.ini

    # 安装和配置frpc
    ssh ${GPU_HOST} << 'EOF'
sudo mv /tmp/frpc.ini /etc/frp/frpc.ini
sudo systemctl restart frpc
sudo systemctl enable frpc
EOF

    echo "${GPU_HOST} 部署完成"
done

echo "所有GPU机器部署完成"
```

---

## 使用流程

1. **生成frpc配置**: `./generate_frpc_configs.sh`
2. **生成nginx配置**: `./generate_nginx_frp.sh`
3. **生成SSH配置**: `./generate_ssh_config.sh`
4. **批量部署frpc**: `./deploy_frpc.sh`

---

## 注意事项

- 修改脚本中的IP、域名、token等参数
- 确保有足够的权限执行脚本
- 批量部署前先在1-2台机器测试
- 保存好生成的配置文件备份
