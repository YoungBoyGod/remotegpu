# 批量配置生成脚本

## frpc配置生成脚本

**执行位置**：本地开发机

创建脚本 `generate_frpc_config.sh`：

```bash
#!/bin/bash

# 参数
GPU_NUM=$1  # GPU编号（1-200）
SERVER_IP="云服务器IP"
TOKEN="your_secure_token_change_this"

# 本地服务端口（根据实际情况修改）
LOCAL_SSH=22
LOCAL_JUPYTER=8888
LOCAL_TENSORBOARD=6006
LOCAL_SERVICE1=本地端口1
LOCAL_SERVICE2=本地端口2

# 计算远程端口
REMOTE_SSH=$((2200 + GPU_NUM))
REMOTE_JUPYTER=$((8000 + GPU_NUM))
REMOTE_TENSORBOARD=$((9000 + GPU_NUM))
REMOTE_SERVICE1=$((10000 + GPU_NUM))
REMOTE_SERVICE2=$((11000 + GPU_NUM))

# 生成配置
cat > frpc_gpu${GPU_NUM}.toml <<EOF
serverAddr = "${SERVER_IP}"
serverPort = 7000

auth.method = "token"
auth.token = "${TOKEN}"

log.to = "/var/log/frp/frpc.log"
log.level = "info"

[[proxies]]
name = "gpu${GPU_NUM}-ssh"
type = "tcp"
localIP = "127.0.0.1"
localPort = ${LOCAL_SSH}
remotePort = ${REMOTE_SSH}

[[proxies]]
name = "gpu${GPU_NUM}-jupyter"
type = "tcp"
localIP = "127.0.0.1"
localPort = ${LOCAL_JUPYTER}
remotePort = ${REMOTE_JUPYTER}

[[proxies]]
name = "gpu${GPU_NUM}-tensorboard"
type = "tcp"
localIP = "127.0.0.1"
localPort = ${LOCAL_TENSORBOARD}
remotePort = ${REMOTE_TENSORBOARD}

[[proxies]]
name = "gpu${GPU_NUM}-service1"
type = "tcp"
localIP = "127.0.0.1"
localPort = ${LOCAL_SERVICE1}
remotePort = ${REMOTE_SERVICE1}

[[proxies]]
name = "gpu${GPU_NUM}-service2"
type = "tcp"
localIP = "127.0.0.1"
localPort = ${LOCAL_SERVICE2}
remotePort = ${REMOTE_SERVICE2}
EOF

echo "生成配置文件: frpc_gpu${GPU_NUM}.toml"
```

使用方法：
```bash
chmod +x generate_frpc_config.sh

# 生成GPU1的配置
./generate_frpc_config.sh 1

# 批量生成所有配置
for i in {1..200}; do
  ./generate_frpc_config.sh $i
done
```

---

## nginx配置生成脚本

创建脚本 `generate_nginx_config.sh`：

```bash
#!/bin/bash

DOMAIN="gpu.domain.com"
SSL_CERT="/etc/letsencrypt/live/gpu.domain.com/fullchain.pem"
SSL_KEY="/etc/letsencrypt/live/gpu.domain.com/privkey.pem"

for i in {1..200}; do
  JUPYTER_PORT=$((8000 + i))
  TB_PORT=$((9000 + i))
  S1_PORT=$((10000 + i))
  S2_PORT=$((11000 + i))

  cat >> gpu_services.conf <<EOF
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
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
    }
}

EOF
done

echo "生成配置文件: gpu_services.conf"
```
