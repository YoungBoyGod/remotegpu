#!/bin/bash

# Docker 存储迁移工具
# 用于将 Docker 根目录迁移到更大的分区

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${BLUE}╔═══════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     Docker 存储迁移工具                        ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════╝${NC}"
echo ""

# 检查是否为 root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}[错误] 请使用 root 或 sudo 运行此脚本${NC}"
    exit 1
fi

# 显示当前配置
echo -e "${CYAN}当前 Docker 配置：${NC}"
CURRENT_ROOT=$(docker info 2>/dev/null | grep "Docker Root Dir" | awk '{print $4}')
CURRENT_DRIVER=$(docker info 2>/dev/null | grep "Storage Driver" | awk '{print $3}')

if [ -z "$CURRENT_ROOT" ]; then
    echo -e "${RED}[错误] Docker 未运行或未安装${NC}"
    exit 1
fi

echo "Docker 根目录: ${GREEN}${CURRENT_ROOT}${NC}"
echo "存储驱动: ${GREEN}${CURRENT_DRIVER}${NC}"
echo ""

# 显示磁盘使用情况
echo -e "${CYAN}当前存储使用情况：${NC}"
du -sh ${CURRENT_ROOT} 2>/dev/null || echo "无法读取"
echo ""

# 显示 Docker 资源使用
echo -e "${CYAN}Docker 资源占用：${NC}"
docker system df
echo ""

# 显示可用分区
echo -e "${CYAN}可用磁盘分区：${NC}"
df -h | grep -E "Filesystem|/dev|/mnt|/data"
echo ""

# 选择目标目录
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}请输入新的 Docker 根目录路径：${NC}"
echo "推荐路径："
echo "  /data/docker"
echo "  /mnt/docker"
echo "  /opt/docker"
echo ""
read -p "新路径: " NEW_ROOT

# 验证输入
if [ -z "$NEW_ROOT" ]; then
    echo -e "${RED}[错误] 路径不能为空${NC}"
    exit 1
fi

# 检查路径是否已存在
if [ -d "$NEW_ROOT" ] && [ "$(ls -A $NEW_ROOT)" ]; then
    echo -e "${YELLOW}[警告] 目录已存在且不为空${NC}"
    read -p "是否继续？现有文件可能被覆盖 [y/N]: " confirm
    if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
        echo "已取消"
        exit 0
    fi
fi

# 检查目标分区空间
TARGET_PARTITION=$(df -P "$NEW_ROOT" 2>/dev/null | tail -1 | awk '{print $1}' || df -P "$(dirname "$NEW_ROOT")" | tail -1 | awk '{print $1}')
AVAILABLE_SPACE=$(df -BG "$(dirname "$NEW_ROOT")" | tail -1 | awk '{print $4}' | sed 's/G//')
REQUIRED_SPACE=$(du -sb ${CURRENT_ROOT} 2>/dev/null | awk '{print int($1/1024/1024/1024)+1}')

echo ""
echo -e "${CYAN}空间检查：${NC}"
echo "目标分区: ${TARGET_PARTITION}"
echo "可用空间: ${AVAILABLE_SPACE}GB"
echo "需要空间: ~${REQUIRED_SPACE}GB"
echo ""

if [ "$AVAILABLE_SPACE" -lt "$REQUIRED_SPACE" ]; then
    echo -e "${RED}[错误] 目标分区空间不足！${NC}"
    exit 1
fi

# 确认迁移
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}迁移计划：${NC}"
echo "源目录: ${CURRENT_ROOT}"
echo "目标目录: ${NEW_ROOT}"
echo ""
echo -e "${RED}警告：${NC}"
echo "1. 所有 Docker 容器将被停止"
echo "2. 迁移期间 Docker 服务不可用"
echo "3. 根据数据量，可能需要较长时间"
echo ""
read -p "确认开始迁移? [y/N]: " confirm

if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    echo "已取消迁移"
    exit 0
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}开始迁移...${NC}"
echo ""

# 1. 停止 Docker
echo -e "${CYAN}[1/6] 停止 Docker 服务...${NC}"
systemctl stop docker
systemctl stop docker.socket
sleep 2

# 确认 Docker 已停止
if systemctl is-active --quiet docker; then
    echo -e "${RED}[错误] Docker 服务未能停止${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Docker 已停止${NC}"
echo ""

# 2. 创建目标目录
echo -e "${CYAN}[2/6] 创建目标目录...${NC}"
mkdir -p "$NEW_ROOT"
echo -e "${GREEN}✓ 目录已创建: ${NEW_ROOT}${NC}"
echo ""

# 3. 同步数据
echo -e "${CYAN}[3/6] 同步数据（这可能需要一些时间）...${NC}"
echo "使用 rsync 同步数据，保持进度显示..."
rsync -aP "${CURRENT_ROOT}/" "${NEW_ROOT}/"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 数据同步完成${NC}"
else
    echo -e "${RED}[错误] 数据同步失败${NC}"
    echo "正在恢复 Docker 服务..."
    systemctl start docker
    exit 1
fi
echo ""

# 4. 备份并更新配置
echo -e "${CYAN}[4/6] 更新 Docker 配置...${NC}"

# 备份原配置
if [ -f /etc/docker/daemon.json ]; then
    cp /etc/docker/daemon.json /etc/docker/daemon.json.backup
    echo "原配置已备份到: /etc/docker/daemon.json.backup"
fi

# 创建新配置
cat > /etc/docker/daemon.json << EOF
{
  "data-root": "${NEW_ROOT}",
  "storage-driver": "overlay2",
  "storage-opts": [
    "overlay2.override_kernel_check=true"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF

echo -e "${GREEN}✓ 配置已更新${NC}"
echo ""

# 5. 启动 Docker
echo -e "${CYAN}[5/6] 启动 Docker 服务...${NC}"
systemctl daemon-reload
systemctl start docker

# 等待 Docker 启动
sleep 3

# 验证 Docker 启动
if ! systemctl is-active --quiet docker; then
    echo -e "${RED}[错误] Docker 启动失败！${NC}"
    echo "请检查日志: journalctl -xeu docker"
    echo ""
    echo "如需回滚，执行："
    echo "  cp /etc/docker/daemon.json.backup /etc/docker/daemon.json"
    echo "  systemctl restart docker"
    exit 1
fi
echo -e "${GREEN}✓ Docker 已启动${NC}"
echo ""

# 6. 验证迁移
echo -e "${CYAN}[6/6] 验证迁移结果...${NC}"

NEW_ROOT_CHECK=$(docker info 2>/dev/null | grep "Docker Root Dir" | awk '{print $4}')

if [ "$NEW_ROOT_CHECK" = "$NEW_ROOT" ]; then
    echo -e "${GREEN}✓ 验证成功！Docker 现在使用新路径${NC}"
    echo ""

    # 显示迁移结果
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}迁移完成！${NC}"
    echo ""
    echo "新 Docker 根目录: ${GREEN}${NEW_ROOT_CHECK}${NC}"
    echo ""

    # 验证镜像和容器
    echo "验证数据完整性..."
    IMAGE_COUNT=$(docker images -q | wc -l)
    CONTAINER_COUNT=$(docker ps -a -q | wc -l)
    echo "镜像数量: ${IMAGE_COUNT}"
    echo "容器数量: ${CONTAINER_COUNT}"
    echo ""

    # 询问是否删除旧数据
    echo -e "${YELLOW}旧数据仍保留在: ${CURRENT_ROOT}${NC}"
    echo -e "${YELLOW}占用空间: $(du -sh ${CURRENT_ROOT} 2>/dev/null | awk '{print $1}')${NC}"
    echo ""
    read -p "确认一切正常后，是否删除旧数据? [y/N]: " delete_old

    if [[ "$delete_old" =~ ^[Yy]$ ]]; then
        echo ""
        echo "删除旧数据..."
        rm -rf "${CURRENT_ROOT}"
        echo -e "${GREEN}✓ 旧数据已删除${NC}"
    else
        echo ""
        echo "旧数据已保留，可稍后手动删除："
        echo -e "  ${CYAN}sudo rm -rf ${CURRENT_ROOT}${NC}"
    fi

else
    echo -e "${RED}[错误] 验证失败！Docker 仍使用旧路径${NC}"
    echo "当前路径: ${NEW_ROOT_CHECK}"
    echo "期望路径: ${NEW_ROOT}"
    exit 1
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}🎉 迁移成功完成！${NC}"
echo ""
echo "后续步骤："
echo "1. 重启相关容器："
echo -e "   ${CYAN}docker-compose up -d${NC}"
echo ""
echo "2. 验证容器运行正常："
echo -e "   ${CYAN}docker ps${NC}"
echo ""
echo "3. 监控存储使用："
echo -e "   ${CYAN}docker system df${NC}"
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
