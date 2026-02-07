#!/bin/bash

# RemoteGPU 服务停止脚本
# 按照反向依赖顺序停止所有服务

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "========================================="
echo "RemoteGPU 服务停止脚本"
echo "========================================="
echo ""

# 获取所有服务目录（排除外部服务）
all_services=("harbor" "guacamole" "uptime-kuma" "grafana" "prometheus" "nginx" "rustfs" "etcd" "redis" "postgresql")

echo "注意: JumpServer 使用外部服务，不在此停止"
echo "停止所有本地服务..."
for service in "${all_services[@]}"; do
    if [ -d "$service" ]; then
        echo "停止 $service..."
        cd "$service"
        docker-compose down
        cd ..
    fi
done
echo ""

echo "========================================="
echo "所有服务已停止！"
echo "========================================="
echo ""
echo "如需删除数据卷，请执行："
echo "  docker volume ls | grep remotegpu"
echo "  docker volume rm <volume-name>"
echo ""
echo "如需删除共享网络，请执行："
echo "  docker network rm remotegpu-network"
echo ""
