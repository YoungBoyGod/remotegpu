#!/bin/bash

# 容器资源限制配置工具
# 根据预设场景快速配置资源限制

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${BLUE}╔═══════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     容器资源限制配置工具                      ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════╝${NC}"
echo ""

# 检查 docker-compose.yml 是否存在
if [ ! -f "docker-compose.yml" ]; then
    echo -e "${RED}[错误] 找不到 docker-compose.yml 文件${NC}"
    exit 1
fi

# 备份原配置
cp docker-compose.yml docker-compose.yml.backup
echo -e "${GREEN}[完成] 已备份配置到 docker-compose.yml.backup${NC}"
echo ""

# 选择配置场景
echo -e "${CYAN}请选择资源配置场景：${NC}"
echo ""
echo -e "${YELLOW}1.${NC} 小型工作空间 ${GREEN}(适合 1-2 用户)${NC}"
echo "   - 内存: 8GB"
echo "   - CPU: 4 核"
echo "   - 磁盘: 50GB (可写层)"
echo "   - GPU: 1 个"
echo ""
echo -e "${YELLOW}2.${NC} 中型工作空间 ${GREEN}(适合 3-5 用户)${NC}"
echo "   - 内存: 16GB"
echo "   - CPU: 8 核"
echo "   - 磁盘: 100GB (可写层)"
echo "   - GPU: 1 个"
echo ""
echo -e "${YELLOW}3.${NC} 大型工作空间 ${GREEN}(适合 5-10 用户)${NC}"
echo "   - 内存: 32GB"
echo "   - CPU: 16 核"
echo "   - 磁盘: 200GB (可写层)"
echo "   - GPU: 1-2 个"
echo ""
echo -e "${YELLOW}4.${NC} 生产环境 ${GREEN}(适合 10+ 用户)${NC}"
echo "   - 内存: 64GB"
echo "   - CPU: 32 核"
echo "   - 磁盘: 500GB (可写层)"
echo "   - GPU: 2-4 个"
echo ""
echo -e "${YELLOW}5.${NC} 自定义配置"
echo ""
echo -e "${YELLOW}6.${NC} 无限制模式 ${RED}(不推荐)${NC}"
echo ""

read -p "请选择 [1-6]: " choice

case $choice in
    1)
        MEMORY="8g"
        MEMORY_RESERVE="4g"
        CPUS="4"
        STORAGE="50G"
        GPU_COUNT="0"
        PIDS="1000"
        ;;
    2)
        MEMORY="16g"
        MEMORY_RESERVE="8g"
        CPUS="8"
        STORAGE="100G"
        GPU_COUNT="0"
        PIDS="2000"
        ;;
    3)
        MEMORY="32g"
        MEMORY_RESERVE="16g"
        CPUS="16"
        STORAGE="200G"
        GPU_COUNT="0"
        PIDS="4000"
        ;;
    4)
        MEMORY="64g"
        MEMORY_RESERVE="32g"
        CPUS="32"
        STORAGE="500G"
        GPU_COUNT="0,1"
        PIDS="8000"
        ;;
    5)
        echo ""
        echo -e "${CYAN}自定义资源配置：${NC}"
        read -p "内存限制 (例如: 16g): " MEMORY
        read -p "内存保留 (例如: 8g): " MEMORY_RESERVE
        read -p "CPU 核心数 (例如: 8): " CPUS
        read -p "磁盘限制 (例如: 100G): " STORAGE
        read -p "GPU 编号 (例如: 0 或 0,1 或 all): " GPU_COUNT
        read -p "进程数限制 (例如: 2000): " PIDS
        ;;
    6)
        echo -e "${RED}[警告] 将移除所有资源限制！${NC}"
        read -p "确认继续? [y/N]: " confirm
        if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
            echo "已取消"
            exit 0
        fi
        MEMORY="unlimited"
        ;;
    *)
        echo -e "${RED}[错误] 无效选择${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

if [ "$MEMORY" != "unlimited" ]; then
    echo -e "${GREEN}正在配置资源限制...${NC}"
    echo ""
    echo "内存限制: $MEMORY"
    echo "内存保留: $MEMORY_RESERVE"
    echo "CPU 核心: $CPUS"
    echo "磁盘限制: $STORAGE"
    echo "GPU 配置: $GPU_COUNT"
    echo "进程限制: $PIDS"
    echo ""

    # 更新 docker-compose.yml
    sed -i "s/mem_limit:.*/mem_limit: $MEMORY/" docker-compose.yml
    sed -i "s/mem_reservation:.*/mem_reservation: $MEMORY_RESERVE/" docker-compose.yml
    sed -i "s/memswap_limit:.*/memswap_limit: $MEMORY/" docker-compose.yml
    sed -i "s/cpus:.*/cpus: $CPUS/" docker-compose.yml
    sed -i "s/pids_limit:.*/pids_limit: $PIDS/" docker-compose.yml

    # 更新 storage_opt（启用磁盘限制）
    sed -i "s/# storage_opt:/storage_opt:/" docker-compose.yml
    sed -i "s/#   size:.*/  size: '$STORAGE'/" docker-compose.yml

    # 更新 GPU 配置
    sed -i "s/- NVIDIA_VISIBLE_DEVICES=.*/- NVIDIA_VISIBLE_DEVICES=$GPU_COUNT/" docker-compose.yml

else
    echo -e "${YELLOW}正在移除资源限制...${NC}"
    echo ""

    # 注释掉资源限制
    sed -i "s/mem_limit:/# mem_limit:/" docker-compose.yml
    sed -i "s/mem_reservation:/# mem_reservation:/" docker-compose.yml
    sed -i "s/memswap_limit:/# memswap_limit:/" docker-compose.yml
    sed -i "s/^[[:space:]]*cpus:/# cpus:/" docker-compose.yml
    sed -i "s/pids_limit:/# pids_limit:/" docker-compose.yml
    sed -i "s/storage_opt:/# storage_opt:/" docker-compose.yml
    sed -i "s/^[[:space:]]*size:/  # size:/" docker-compose.yml
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${GREEN}✓ 配置已更新！${NC}"
echo ""

# 显示差异
echo -e "${CYAN}配置变更：${NC}"
diff docker-compose.yml.backup docker-compose.yml || true
echo ""

# 询问是否重启容器
echo -e "${YELLOW}需要重启容器才能应用新配置${NC}"
read -p "是否立即重启容器? [y/N]: " restart_now

if [[ "$restart_now" =~ ^[Yy]$ ]]; then
    echo ""
    echo -e "${GREEN}正在重启容器...${NC}"
    docker-compose down
    docker-compose up -d
    echo ""
    echo -e "${GREEN}✓ 容器已重启${NC}"
    echo ""
    echo "查看资源使用情况："
    echo -e "${CYAN}docker stats user001-workspace${NC}"
else
    echo ""
    echo "请手动重启容器："
    echo -e "${CYAN}docker-compose down && docker-compose up -d${NC}"
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}配置完成！${NC}"
echo ""
echo "有用的命令："
echo -e "  ${CYAN}docker stats user001-workspace${NC}          - 实时监控资源使用"
echo -e "  ${CYAN}docker exec -it user001-workspace bash${NC}  - 进入容器"
echo -e "  ${CYAN}df -h${NC}                                   - 查看磁盘使用"
echo -e "  ${CYAN}free -h${NC}                                 - 查看内存使用"
echo -e "  ${CYAN}nvidia-smi${NC}                              - 查看 GPU 使用"
echo ""
echo -e "恢复备份："
echo -e "  ${CYAN}cp docker-compose.yml.backup docker-compose.yml${NC}"
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
