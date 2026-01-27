#!/bin/bash

# RemoteGPU 服务状态检查脚本

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "========================================="
echo "RemoteGPU 服务状态检查"
echo "========================================="
echo ""
echo "注意: JumpServer 和 Kubernetes 使用外部服务"
echo ""

# 定义所有本地服务
declare -A services=(
    ["postgresql"]="remotegpu-postgresql"
    ["redis"]="remotegpu-redis"
    ["etcd"]="remotegpu-etcd"
    ["rustfs"]="remotegpu-rustfs"
    ["nginx"]="remotegpu-nginx"
    ["prometheus"]="remotegpu-prometheus"
    ["grafana"]="remotegpu-grafana"
    ["uptime-kuma"]="remotegpu-uptime-kuma"
    ["guacamole"]="remotegpu-guacamole"
    ["harbor-nginx"]="remotegpu-harbor-nginx"
)

# 检查每个服务的状态
for service_name in "${!services[@]}"; do
    container_name="${services[$service_name]}"

    # 检查容器是否存在
    if docker ps -a --format '{{.Names}}' | grep -q "^${container_name}$"; then
        # 获取容器状态
        status=$(docker inspect --format='{{.State.Status}}' "$container_name" 2>/dev/null)
        health=$(docker inspect --format='{{.State.Health.Status}}' "$container_name" 2>/dev/null)

        # 显示状态
        printf "%-20s " "$service_name"

        if [ "$status" = "running" ]; then
            echo -n "✓ Running"
            if [ "$health" != "<no value>" ] && [ -n "$health" ]; then
                echo " (Health: $health)"
            else
                echo ""
            fi
        else
            echo "✗ $status"
        fi
    else
        printf "%-20s ✗ Not found\n" "$service_name"
    fi
done

echo ""
echo "========================================="
echo "端口使用情况"
echo "========================================="
echo ""

# 检查端口
declare -A ports=(
    ["5432"]="PostgreSQL"
    ["6379"]="Redis"
    ["2379"]="Etcd"
    ["9000"]="RustFS API"
    ["9001"]="RustFS Console"
    ["80"]="Nginx HTTP"
    ["443"]="Nginx HTTPS"
    ["9090"]="Prometheus"
    ["3000"]="Grafana"
    ["3001"]="Uptime Kuma"
    ["8080"]="JumpServer"
    ["8081"]="Guacamole"
    ["8082"]="Harbor"
)

for port in "${!ports[@]}"; do
    service="${ports[$port]}"
    if netstat -tuln 2>/dev/null | grep -q ":$port "; then
        printf "%-6s ✓ %s\n" "$port" "$service"
    else
        printf "%-6s ✗ %s (未监听)\n" "$port" "$service"
    fi
done

echo ""
