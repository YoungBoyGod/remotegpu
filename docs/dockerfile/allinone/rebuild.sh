#!/bin/bash

# 重新构建镜像并重启容器
# 用于应用 Dockerfile 更新

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔═══════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     镜像重建和容器重启工具                    ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════╝${NC}"
echo ""

# 检查 Dockerfile 是否存在
if [ ! -f "Dockerfile" ]; then
    echo -e "${RED}[错误] 找不到 Dockerfile${NC}"
    exit 1
fi

# 检查 docker-compose.yml 是否存在
if [ ! -f "docker-compose.yml" ]; then
    echo -e "${RED}[错误] 找不到 docker-compose.yml${NC}"
    exit 1
fi

echo -e "${CYAN}准备重建镜像...${NC}"
echo ""

# 显示当前运行的容器
echo -e "${CYAN}当前运行的容器：${NC}"
docker-compose ps
echo ""

# 询问是否继续
read -p "这将停止现有容器并重建镜像，是否继续? [y/N]: " confirm

if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    echo "已取消"
    exit 0
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}开始重建...${NC}"
echo ""

# 1. 停止容器
echo -e "${CYAN}[1/4] 停止容器...${NC}"
docker-compose down
echo -e "${GREEN}✓ 容器已停止${NC}"
echo ""

# 2. 删除旧镜像（可选）
read -p "是否删除旧镜像以节省空间? [y/N]: " delete_old

if [[ "$delete_old" =~ ^[Yy]$ ]]; then
    echo -e "${CYAN}删除旧镜像...${NC}"
    docker rmi gpu-workspace:latest || echo "旧镜像不存在或已被使用"
    echo ""
fi

# 3. 构建新镜像
echo -e "${CYAN}[2/4] 构建新镜像（这可能需要几分钟）...${NC}"
docker build -t gpu-workspace:latest .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 镜像构建成功${NC}"
else
    echo -e "${RED}[错误] 镜像构建失败${NC}"
    exit 1
fi
echo ""

# 4. 启动容器
echo -e "${CYAN}[3/4] 启动容器...${NC}"
docker-compose up -d

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 容器已启动${NC}"
else
    echo -e "${RED}[错误] 容器启动失败${NC}"
    exit 1
fi
echo ""

# 5. 验证
echo -e "${CYAN}[4/4] 验证...${NC}"
sleep 3

# 检查容器状态
if docker-compose ps | grep -q "Up"; then
    echo -e "${GREEN}✓ 容器运行正常${NC}"
else
    echo -e "${RED}[错误] 容器未正常运行${NC}"
    echo "查看日志："
    docker-compose logs --tail=50
    exit 1
fi

# 检查系统库是否安装
echo ""
echo -e "${CYAN}验证 OpenCV 依赖...${NC}"
if docker exec user001-workspace bash -c "ldconfig -p | grep -q libGL.so.1"; then
    echo -e "${GREEN}✓ libGL.so.1 已安装${NC}"
else
    echo -e "${YELLOW}⚠ libGL.so.1 未找到，但可能不影响使用${NC}"
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}🎉 重建完成！${NC}"
echo ""

echo -e "${CYAN}后续步骤：${NC}"
echo ""
echo "1. 查看容器日志："
echo -e "   ${BLUE}docker-compose logs -f${NC}"
echo ""
echo "2. 进入容器测试："
echo -e "   ${BLUE}docker exec -it user001-workspace bash${NC}"
echo ""
echo "3. 测试 OpenCV："
echo -e "   ${BLUE}docker exec user001-workspace python3 -c 'import cv2; print(cv2.__version__)'${NC}"
echo ""
echo "4. 安装深度学习框架："
echo -e "   ${BLUE}# 在 Jupyter 或 SSH 中运行${NC}"
echo -e "   ${BLUE}pip install torch torchvision ultralytics${NC}"
echo ""
echo "5. 访问服务："
echo -e "   Jupyter Lab: ${BLUE}http://localhost:18888${NC}"
echo -e "   VSCode Web:  ${BLUE}http://localhost:18080${NC}"
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
