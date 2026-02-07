#!/bin/bash
# RemoteGPU GPU Exporter 部署脚本
# 在 GPU 机器上部署 DCGM Exporter 或 NVIDIA GPU Exporter
# 用法: sudo bash install-gpu-exporter.sh [选项]
#   --type dcgm|nvidia   Exporter 类型（默认: dcgm）
#   --port PORT          监听端口（默认: 9400）
#   --prometheus URL     Prometheus 地址（用于验证连通性）

set -euo pipefail

# 默认配置
EXPORTER_TYPE="dcgm"
EXPORTER_PORT=9400
PROMETHEUS_URL=""

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

# ========== 参数解析 ==========
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --type)       EXPORTER_TYPE="$2"; shift 2 ;;
            --port)       EXPORTER_PORT="$2"; shift 2 ;;
            --prometheus) PROMETHEUS_URL="$2"; shift 2 ;;
            -h|--help)    usage; exit 0 ;;
            *)            log_error "未知参数: $1"; usage; exit 1 ;;
        esac
    done
}

usage() {
    cat <<EOF
用法: sudo bash install-gpu-exporter.sh [选项]

选项:
  --type dcgm|nvidia    Exporter 类型（默认: dcgm）
    dcgm   — NVIDIA DCGM Exporter（推荐，需要 Docker）
    nvidia — nvidia_gpu_exporter（独立二进制，无需 Docker）
  --port PORT           监听端口（默认: 9400）
  --prometheus URL      Prometheus 地址（可选，用于验证连通性）
  -h, --help            显示帮助信息

示例:
  sudo bash install-gpu-exporter.sh --type dcgm
  sudo bash install-gpu-exporter.sh --type nvidia --port 9835
EOF
}

# ========== 前置检查 ==========
check_prerequisites() {
    if [[ $EUID -ne 0 ]]; then
        log_error "请使用 root 权限运行此脚本"
        exit 1
    fi

    # 检查 NVIDIA 驱动
    if ! command -v nvidia-smi &>/dev/null; then
        log_error "未检测到 nvidia-smi，请先安装 NVIDIA 驱动"
        exit 1
    fi

    log_info "检测到 NVIDIA GPU:"
    nvidia-smi --query-gpu=name,driver_version --format=csv,noheader | while read -r line; do
        log_info "  $line"
    done

    if [[ "$EXPORTER_TYPE" == "dcgm" ]]; then
        if ! command -v docker &>/dev/null; then
            log_error "DCGM Exporter 需要 Docker，请先安装或使用 --type nvidia"
            exit 1
        fi
        # 检查 NVIDIA Container Toolkit
        if ! docker info 2>/dev/null | grep -q "nvidia"; then
            log_warn "未检测到 NVIDIA Container Toolkit，DCGM 可能无法访问 GPU"
            log_warn "安装方法: https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/"
        fi
    fi
}

# ========== 安装 DCGM Exporter（Docker 方式） ==========
install_dcgm_exporter() {
    local container_name="remotegpu-dcgm-exporter"

    # 停止已有容器
    if docker ps -a --format '{{.Names}}' | grep -q "^${container_name}$"; then
        log_warn "检测到已有容器，先移除..."
        docker rm -f "$container_name" 2>/dev/null || true
    fi

    log_info "拉取 DCGM Exporter 镜像..."
    docker pull nvcr.io/nvidia/k8s/dcgm-exporter:3.3.5-3.4.1-ubuntu22.04

    log_info "启动 DCGM Exporter 容器..."
    docker run -d \
        --name "$container_name" \
        --restart always \
        --runtime=nvidia \
        --gpus all \
        -p "${EXPORTER_PORT}:9400" \
        nvcr.io/nvidia/k8s/dcgm-exporter:3.3.5-3.4.1-ubuntu22.04

    sleep 3
    if docker ps --format '{{.Names}}' | grep -q "^${container_name}$"; then
        log_info "DCGM Exporter 启动成功，端口: $EXPORTER_PORT"
    else
        log_error "DCGM Exporter 启动失败"
        docker logs "$container_name" 2>&1 | tail -20
        exit 1
    fi
}

# ========== 安装 nvidia_gpu_exporter（二进制方式） ==========
install_nvidia_exporter() {
    local install_dir="/opt/nvidia-gpu-exporter"
    local binary_name="nvidia_gpu_exporter"
    local service_name="nvidia-gpu-exporter"
    local version="1.2.1"
    local arch
    arch=$(uname -m)

    case "$arch" in
        x86_64)  arch="amd64" ;;
        aarch64) arch="arm64" ;;
        *) log_error "不支持的架构: $arch"; exit 1 ;;
    esac

    mkdir -p "$install_dir"

    local download_url="https://github.com/utkuozdemir/nvidia_gpu_exporter/releases/download/v${version}/nvidia_gpu_exporter_${version}_linux_${arch}.tar.gz"

    log_info "下载 nvidia_gpu_exporter v${version} (${arch})..."
    curl -fSL -o /tmp/nvidia_gpu_exporter.tar.gz "$download_url"

    log_info "解压安装..."
    tar -xzf /tmp/nvidia_gpu_exporter.tar.gz -C "$install_dir" "$binary_name"
    chmod +x "$install_dir/$binary_name"
    rm -f /tmp/nvidia_gpu_exporter.tar.gz

    # 创建 systemd 服务
    log_info "创建 systemd 服务..."
    cat > "/etc/systemd/system/${service_name}.service" <<EOF
[Unit]
Description=NVIDIA GPU Exporter
After=network-online.target

[Service]
Type=simple
ExecStart=${install_dir}/${binary_name} --web.listen-address=:${EXPORTER_PORT}
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable "$service_name"
    systemctl start "$service_name"

    sleep 2
    if systemctl is-active --quiet "$service_name"; then
        log_info "nvidia_gpu_exporter 启动成功，端口: $EXPORTER_PORT"
    else
        log_error "nvidia_gpu_exporter 启动失败"
        journalctl -u "$service_name" --no-pager -n 20
        exit 1
    fi
}

# ========== 验证 Exporter ==========
verify_exporter() {
    log_info "验证 Exporter metrics 端点..."
    local url="http://localhost:${EXPORTER_PORT}/metrics"

    if curl -sf "$url" | grep -q "gpu"; then
        log_info "Metrics 端点正常，已检测到 GPU 指标"
    else
        log_warn "Metrics 端点可能异常，请手动检查: curl $url"
    fi

    if [[ -n "$PROMETHEUS_URL" ]]; then
        log_info "提示: 请在 Prometheus 配置中添加以下 scrape target:"
        local host_ip
        host_ip=$(hostname -I | awk '{print $1}')
        echo ""
        echo "  - job_name: 'nvidia-gpu'"
        echo "    static_configs:"
        echo "      - targets: ['${host_ip}:${EXPORTER_PORT}']"
        echo "        labels:"
        echo "          group: 'gpu-nodes'"
        echo ""
    fi
}

# ========== 主流程 ==========
main() {
    echo ""
    echo "========================================="
    echo "  RemoteGPU GPU Exporter 部署"
    echo "========================================="
    echo ""

    parse_args "$@"
    check_prerequisites

    case "$EXPORTER_TYPE" in
        dcgm)
            log_info "安装 DCGM Exporter (Docker)..."
            install_dcgm_exporter
            ;;
        nvidia)
            log_info "安装 nvidia_gpu_exporter (二进制)..."
            install_nvidia_exporter
            ;;
        *)
            log_error "不支持的类型: $EXPORTER_TYPE（可选: dcgm, nvidia）"
            exit 1
            ;;
    esac

    verify_exporter

    echo ""
    log_info "GPU Exporter 部署完成"
    echo ""
}

main "$@"
