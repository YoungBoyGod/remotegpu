#!/bin/bash

# Docker 存储监控和清理工具

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# 配置
DOCKER_ROOT=$(docker info 2>/dev/null | grep "Docker Root Dir" | awk '{print $4}')
DATA_DIR="/home/luo/code/remotegpu/docs/dockerfile/allinone/data"

echo -e "${BLUE}╔═══════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     Docker 存储监控和清理工具                  ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════╝${NC}"
echo ""

# 1. Docker 存储概览
echo -e "${CYAN}━━━ Docker 存储使用情况 ━━━${NC}"
docker system df
echo ""

# 2. 详细分类统计
echo -e "${CYAN}━━━ 详细存储统计 ━━━${NC}"

# 镜像
IMAGE_COUNT=$(docker images -q | wc -l)
IMAGE_SIZE=$(docker system df | grep "Images" | awk '{print $3, $4}')
echo -e "镜像: ${GREEN}${IMAGE_COUNT}${NC} 个, 占用 ${YELLOW}${IMAGE_SIZE}${NC}"

# 容器
CONTAINER_COUNT=$(docker ps -a -q | wc -l)
CONTAINER_RUNNING=$(docker ps -q | wc -l)
CONTAINER_SIZE=$(docker system df | grep "Containers" | awk '{print $3, $4}')
echo -e "容器: ${GREEN}${CONTAINER_COUNT}${NC} 个 (运行中: ${GREEN}${CONTAINER_RUNNING}${NC}), 占用 ${YELLOW}${CONTAINER_SIZE}${NC}"

# 卷
VOLUME_COUNT=$(docker volume ls -q | wc -l)
VOLUME_SIZE=$(docker system df | grep "Local Volumes" | awk '{print $3, $4}')
echo -e "卷: ${GREEN}${VOLUME_COUNT}${NC} 个, 占用 ${YELLOW}${VOLUME_SIZE}${NC}"

# 构建缓存
BUILD_CACHE=$(docker system df | grep "Build Cache" | awk '{print $3, $4}')
echo -e "构建缓存: ${YELLOW}${BUILD_CACHE}${NC}"
echo ""

# 3. 宿主机磁盘使用
echo -e "${CYAN}━━━ 宿主机磁盘使用情况 ━━━${NC}"
df -h | grep -E "Filesystem|/dev|/mnt|/data" | grep -v "tmpfs"
echo ""

# 4. Docker 根目录详情
if [ -n "$DOCKER_ROOT" ] && [ -d "$DOCKER_ROOT" ]; then
    echo -e "${CYAN}━━━ Docker 根目录占用 ━━━${NC}"
    echo "位置: ${GREEN}${DOCKER_ROOT}${NC}"

    # 各子目录大小
    echo ""
    echo "详细分布:"
    du -sh ${DOCKER_ROOT}/* 2>/dev/null | sort -h | while read size path; do
        name=$(basename "$path")
        echo "  ${name}: ${YELLOW}${size}${NC}"
    done
    echo ""
fi

# 5. 用户数据目录占用
if [ -d "$DATA_DIR" ]; then
    echo -e "${CYAN}━━━ 用户数据占用 ━━━${NC}"
    du -sh ${DATA_DIR}/* 2>/dev/null | sort -h | while read size path; do
        name=$(basename "$path")
        echo "  ${name}: ${YELLOW}${size}${NC}"
    done
    echo ""
fi

# 6. 检查可回收空间
echo -e "${CYAN}━━━ 可回收空间分析 ━━━${NC}"

# 未使用的镜像
UNUSED_IMAGES=$(docker images -f "dangling=true" -q | wc -l)
if [ "$UNUSED_IMAGES" -gt 0 ]; then
    UNUSED_IMAGE_SIZE=$(docker images -f "dangling=true" --format "{{.Size}}" | sed 's/[^0-9.]//g' | awk '{sum+=$1} END {print sum}')
    echo -e "悬空镜像: ${YELLOW}${UNUSED_IMAGES}${NC} 个 (可清理 ~${UNUSED_IMAGE_SIZE}MB)"
fi

# 停止的容器
STOPPED_CONTAINERS=$(docker ps -a -f "status=exited" -q | wc -l)
if [ "$STOPPED_CONTAINERS" -gt 0 ]; then
    echo -e "已停止容器: ${YELLOW}${STOPPED_CONTAINERS}${NC} 个 (可清理)"
fi

# 未使用的卷
UNUSED_VOLUMES=$(docker volume ls -f "dangling=true" -q | wc -l)
if [ "$UNUSED_VOLUMES" -gt 0 ]; then
    echo -e "未使用的卷: ${YELLOW}${UNUSED_VOLUMES}${NC} 个 (可清理)"
fi

# 构建缓存
BUILD_CACHE_SIZE=$(docker system df | grep "Build Cache" | awk '{print $3}')
if [ -n "$BUILD_CACHE_SIZE" ] && [ "$BUILD_CACHE_SIZE" != "0B" ]; then
    echo -e "构建缓存: ${YELLOW}${BUILD_CACHE_SIZE}${NC} (可清理)"
fi

echo ""

# 7. 存储健康检查
echo -e "${CYAN}━━━ 存储健康检查 ━━━${NC}"

# 检查 Docker 分区使用率
DOCKER_PARTITION=$(df -P "$DOCKER_ROOT" | tail -1)
DOCKER_USAGE=$(echo "$DOCKER_PARTITION" | awk '{print $5}' | sed 's/%//')
DOCKER_AVAILABLE=$(echo "$DOCKER_PARTITION" | awk '{print $4}')

echo -n "Docker 分区使用率: "
if [ "$DOCKER_USAGE" -lt 70 ]; then
    echo -e "${GREEN}${DOCKER_USAGE}%${NC} (良好)"
elif [ "$DOCKER_USAGE" -lt 85 ]; then
    echo -e "${YELLOW}${DOCKER_USAGE}%${NC} (建议清理)"
else
    echo -e "${RED}${DOCKER_USAGE}%${NC} (警告：空间不足！)"
fi

# 检查用户数据分区
if [ -d "$DATA_DIR" ]; then
    DATA_PARTITION=$(df -P "$DATA_DIR" | tail -1)
    DATA_USAGE=$(echo "$DATA_PARTITION" | awk '{print $5}' | sed 's/%//')

    echo -n "数据分区使用率: "
    if [ "$DATA_USAGE" -lt 70 ]; then
        echo -e "${GREEN}${DATA_USAGE}%${NC} (良好)"
    elif [ "$DATA_USAGE" -lt 85 ]; then
        echo -e "${YELLOW}${DATA_USAGE}%${NC} (注意)"
    else
        echo -e "${RED}${DATA_USAGE}%${NC} (警告：空间不足！)"
    fi
fi

echo ""

# 8. 清理建议
if [ "$DOCKER_USAGE" -gt 70 ] || [ "$UNUSED_IMAGES" -gt 0 ] || [ "$STOPPED_CONTAINERS" -gt 0 ]; then
    echo -e "${YELLOW}━━━ 清理建议 ━━━${NC}"
    echo ""
    echo "可执行以下命令清理空间："
    echo ""

    if [ "$UNUSED_IMAGES" -gt 0 ]; then
        echo -e "${CYAN}1. 清理悬空镜像：${NC}"
        echo "   docker image prune"
        echo ""
    fi

    if [ "$STOPPED_CONTAINERS" -gt 0 ]; then
        echo -e "${CYAN}2. 清理停止的容器：${NC}"
        echo "   docker container prune"
        echo ""
    fi

    if [ "$UNUSED_VOLUMES" -gt 0 ]; then
        echo -e "${CYAN}3. 清理未使用的卷：${NC}"
        echo "   docker volume prune"
        echo ""
    fi

    echo -e "${CYAN}4. 清理构建缓存：${NC}"
    echo "   docker builder prune"
    echo ""

    echo -e "${CYAN}5. 完全清理（谨慎！）：${NC}"
    echo "   docker system prune -a --volumes"
    echo ""
fi

# 9. 交互式清理
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
read -p "是否现在执行清理操作? [y/N]: " do_cleanup

if [[ "$do_cleanup" =~ ^[Yy]$ ]]; then
    echo ""
    echo -e "${CYAN}清理选项：${NC}"
    echo "1. 安全清理（推荐）- 仅清理未使用的资源"
    echo "2. 深度清理 - 清理所有未使用的镜像、容器、卷"
    echo "3. 自定义清理"
    echo "4. 取消"
    echo ""
    read -p "请选择 [1-4]: " cleanup_choice

    case $cleanup_choice in
        1)
            echo ""
            echo -e "${GREEN}执行安全清理...${NC}"

            echo "清理悬空镜像..."
            docker image prune -f

            echo "清理停止的容器..."
            docker container prune -f

            echo "清理未使用的网络..."
            docker network prune -f

            echo -e "${GREEN}✓ 安全清理完成${NC}"
            ;;
        2)
            echo ""
            echo -e "${YELLOW}警告：深度清理会删除所有未使用的镜像！${NC}"
            read -p "确认继续? [y/N]: " confirm

            if [[ "$confirm" =~ ^[Yy]$ ]]; then
                echo ""
                echo -e "${GREEN}执行深度清理...${NC}"
                docker system prune -a --volumes -f
                echo -e "${GREEN}✓ 深度清理完成${NC}"
            else
                echo "已取消"
            fi
            ;;
        3)
            echo ""
            echo -e "${CYAN}自定义清理选项：${NC}"

            read -p "清理悬空镜像? [y/N]: " clean_images
            if [[ "$clean_images" =~ ^[Yy]$ ]]; then
                docker image prune -f
            fi

            read -p "清理所有未使用的镜像? [y/N]: " clean_all_images
            if [[ "$clean_all_images" =~ ^[Yy]$ ]]; then
                docker image prune -a -f
            fi

            read -p "清理停止的容器? [y/N]: " clean_containers
            if [[ "$clean_containers" =~ ^[Yy]$ ]]; then
                docker container prune -f
            fi

            read -p "清理未使用的卷? [y/N]: " clean_volumes
            if [[ "$clean_volumes" =~ ^[Yy]$ ]]; then
                docker volume prune -f
            fi

            read -p "清理构建缓存? [y/N]: " clean_build
            if [[ "$clean_build" =~ ^[Yy]$ ]]; then
                docker builder prune -f
            fi

            echo -e "${GREEN}✓ 自定义清理完成${NC}"
            ;;
        4)
            echo "已取消清理"
            ;;
        *)
            echo "无效选择"
            ;;
    esac

    # 显示清理后的状态
    echo ""
    echo -e "${CYAN}清理后状态：${NC}"
    docker system df
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}存储监控完成！${NC}"
echo ""
echo "定期监控建议："
echo "  • 每周运行一次此脚本"
echo "  • 设置 cron 定时任务自动清理"
echo "  • 监控磁盘使用率，保持低于 80%"
echo ""
echo "添加 cron 任务示例："
echo -e "  ${CYAN}0 2 * * 0 $(realpath $0) >> /var/log/docker-storage.log 2>&1${NC}"
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
