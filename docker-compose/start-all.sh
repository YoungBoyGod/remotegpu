#!/bin/bash

# RemoteGPU 服务启动脚本
# 按照依赖顺序启动所有服务

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "========================================="
echo "RemoteGPU 服务启动脚本"
echo "========================================="
echo ""

# 创建共享网络（如果不存在）
echo ">>> 检查共享网络..."
if ! docker network inspect remotegpu-network >/dev/null 2>&1; then
    echo "创建 remotegpu-network 网络..."
    docker network create remotegpu-network
else
    echo "remotegpu-network 网络已存在"
fi
echo ""

# 第一层：基础服务
echo ">>> 第一层：启动基础服务..."
services_layer1=("postgresql" "redis" "etcd")
for service in "${services_layer1[@]}"; do
    echo "启动 $service..."
    cd "$service"
    docker-compose up -d
    cd ..
    sleep 3
done
echo ""

# 第二层：存储和网络
echo ">>> 第二层：启动存储和网络服务..."
services_layer2=("rustfs" "nginx")
for service in "${services_layer2[@]}"; do
    echo "启动 $service..."
    cd "$service"
    docker-compose up -d
    cd ..
    sleep 3
done
echo ""

# 第三层：监控和管理
echo ">>> 第三层：启动监控服务..."
services_layer3=("prometheus" "grafana" "uptime-kuma")
for service in "${services_layer3[@]}"; do
    echo "启动 $service..."
    cd "$service"
    docker-compose up -d
    cd ..
    sleep 2
done
echo ""

# 第四层：测试环境
echo ">>> 第四层：启动测试环境..."
read -p "是否启动测试环境 (test-env)? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "启动 test-env（含 Agent 构建，首次可能较慢）..."
    cd "test-env"
    docker-compose up -d --build
    cd ..
    sleep 5
fi
echo ""

# 第五层：可选服务
echo ">>> 第五层：启动可选服务..."
echo "注意: JumpServer 使用外部服务，不在此启动"
read -p "是否启动 Guacamole, Harbor? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    services_layer4=("guacamole" "harbor")
    for service in "${services_layer4[@]}"; do
        echo "启动 $service..."
        cd "$service"
        docker-compose up -d
        cd ..
        sleep 3
    done
fi
echo ""

echo "========================================="
echo "所有服务启动完成！"
echo "========================================="
echo ""
echo "查看服务状态："
echo "  docker ps"
echo ""
echo "查看服务日志："
echo "  cd <service-dir> && docker-compose logs -f"
echo ""
