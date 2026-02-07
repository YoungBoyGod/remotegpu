#!/bin/bash
# RemoteGPU Agent 一键安装脚本
# 用法: sudo bash install.sh [选项]
#   --server-url URL      后端服务地址（必填）
#   --agent-id ID         Agent 唯一标识（默认自动生成）
#   --machine-id ID       机器 ID（默认自动生成）
#   --token TOKEN         认证 Token（必填）
#   --port PORT           Agent 监听端口（默认 8090）
#   --download-url URL    二进制下载地址（默认从 GitHub Release）

set -euo pipefail

# ========== 默认配置 ==========
INSTALL_DIR="/opt/remotegpu-agent"
CONFIG_DIR="/etc/remotegpu-agent"
DATA_DIR="/var/lib/remotegpu-agent"
LOG_DIR="/var/log/remotegpu-agent"
SERVICE_NAME="remotegpu-agent"
BINARY_NAME="remotegpu-agent"

# 参数默认值
SERVER_URL=""
AGENT_ID=""
MACHINE_ID=""
AGENT_TOKEN=""
AGENT_PORT=8090
DOWNLOAD_URL=""

# ========== 颜色输出 ==========
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
            --server-url)  SERVER_URL="$2"; shift 2 ;;
            --agent-id)    AGENT_ID="$2"; shift 2 ;;
            --machine-id)  MACHINE_ID="$2"; shift 2 ;;
            --token)       AGENT_TOKEN="$2"; shift 2 ;;
            --port)        AGENT_PORT="$2"; shift 2 ;;
            --download-url) DOWNLOAD_URL="$2"; shift 2 ;;
            -h|--help)     usage; exit 0 ;;
            *)             log_error "未知参数: $1"; usage; exit 1 ;;
        esac
    done
}

usage() {
    cat <<EOF
用法: sudo bash install.sh [选项]

必填参数:
  --server-url URL      后端服务地址（如 http://192.168.10.210:8080）
  --token TOKEN         Agent 认证 Token

可选参数:
  --agent-id ID         Agent 唯一标识（默认: agent-<hostname>）
  --machine-id ID       机器 ID（默认: <hostname>）
  --port PORT           Agent 监听端口（默认: 8090）
  --download-url URL    二进制文件下载地址或本地路径
  -h, --help            显示帮助信息

示例:
  sudo bash install.sh \\
    --server-url http://192.168.10.210:8080 \\
    --token your_agent_token \\
    --agent-id agent-gpu-01 \\
    --machine-id gpu-01
EOF
}

# ========== 前置检查 ==========
check_prerequisites() {
    # 检查 root 权限
    if [[ $EUID -ne 0 ]]; then
        log_error "请使用 root 权限运行此脚本"
        exit 1
    fi

    # 检查必填参数
    if [[ -z "$SERVER_URL" ]]; then
        log_error "--server-url 参数必填"
        usage
        exit 1
    fi
    if [[ -z "$AGENT_TOKEN" ]]; then
        log_error "--token 参数必填"
        usage
        exit 1
    fi

    # 自动生成 ID
    local hostname
    hostname=$(hostname -s 2>/dev/null || echo "unknown")
    if [[ -z "$AGENT_ID" ]]; then
        AGENT_ID="agent-${hostname}"
        log_info "自动生成 Agent ID: $AGENT_ID"
    fi
    if [[ -z "$MACHINE_ID" ]]; then
        MACHINE_ID="${hostname}"
        log_info "自动生成 Machine ID: $MACHINE_ID"
    fi

    # 检查 systemd
    if ! command -v systemctl &>/dev/null; then
        log_error "系统不支持 systemd"
        exit 1
    fi

    log_info "前置检查通过"
}

# ========== 安装二进制 ==========
install_binary() {
    log_info "创建安装目录..."
    mkdir -p "$INSTALL_DIR" "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR"

    if [[ -n "$DOWNLOAD_URL" ]]; then
        if [[ -f "$DOWNLOAD_URL" ]]; then
            log_info "从本地路径复制二进制文件: $DOWNLOAD_URL"
            cp "$DOWNLOAD_URL" "$INSTALL_DIR/$BINARY_NAME"
        else
            log_info "从远程地址下载二进制文件: $DOWNLOAD_URL"
            curl -fSL -o "$INSTALL_DIR/$BINARY_NAME" "$DOWNLOAD_URL"
        fi
    else
        log_error "请通过 --download-url 指定二进制文件路径或下载地址"
        exit 1
    fi

    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    log_info "二进制文件安装到: $INSTALL_DIR/$BINARY_NAME"
}

# ========== 生成配置文件 ==========
generate_config() {
    local config_file="$CONFIG_DIR/agent.yaml"

    if [[ -f "$config_file" ]]; then
        log_warn "配置文件已存在，备份为 ${config_file}.bak"
        cp "$config_file" "${config_file}.bak"
    fi

    log_info "生成配置文件: $config_file"
    cat > "$config_file" <<EOF
# RemoteGPU Agent 配置文件（由安装脚本自动生成）

port: ${AGENT_PORT}
db_path: "${DATA_DIR}/tasks.db"
max_workers: 4

server:
  url: "${SERVER_URL}"
  agent_id: "${AGENT_ID}"
  machine_id: "${MACHINE_ID}"
  token: "${AGENT_TOKEN}"
  timeout: 30s

poll:
  interval: 5s
  batch_size: 10

limits:
  max_output_size: 1048576

security:
  allowed_commands:
    - "python"
    - "python3"
    - "bash"
    - "sh"
    - "nvidia-smi"
    - "docker"
  blocked_patterns:
    - "rm -rf /"
    - "mkfs"
    - "dd if="
    - "> /dev/sd"
EOF

    chmod 600 "$config_file"
    log_info "配置文件生成完成"
}

# ========== 创建 systemd 服务 ==========
create_service() {
    local service_file="/etc/systemd/system/${SERVICE_NAME}.service"

    log_info "创建 systemd 服务: $service_file"
    cat > "$service_file" <<EOF
[Unit]
Description=RemoteGPU Agent
Documentation=https://github.com/remotegpu/remotegpu
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=${INSTALL_DIR}/${BINARY_NAME}
WorkingDirectory=${INSTALL_DIR}
Restart=always
RestartSec=10
LimitNOFILE=65536

Environment=AGENT_PORT=${AGENT_PORT}
Environment=AGENT_DB_PATH=${DATA_DIR}/tasks.db
EnvironmentFile=-${CONFIG_DIR}/agent.env

StandardOutput=append:${LOG_DIR}/agent.log
StandardError=append:${LOG_DIR}/agent-error.log

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable "$SERVICE_NAME"
    log_info "systemd 服务创建完成"
}

# ========== 配置日志轮转 ==========
setup_logrotate() {
    log_info "配置日志轮转..."
    cat > /etc/logrotate.d/remotegpu-agent <<EOF
${LOG_DIR}/*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    copytruncate
}
EOF
    log_info "日志轮转配置完成"
}

# ========== 主流程 ==========
main() {
    echo ""
    echo "========================================="
    echo "  RemoteGPU Agent 安装程序"
    echo "========================================="
    echo ""

    parse_args "$@"
    check_prerequisites

    # 如果服务已运行，先停止
    if systemctl is-active --quiet "$SERVICE_NAME" 2>/dev/null; then
        log_warn "检测到 Agent 正在运行，先停止..."
        systemctl stop "$SERVICE_NAME"
    fi

    install_binary
    generate_config
    create_service
    setup_logrotate

    # 启动服务
    log_info "启动 Agent 服务..."
    systemctl start "$SERVICE_NAME"

    # 验证
    sleep 2
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        log_info "Agent 服务启动成功"
    else
        log_error "Agent 服务启动失败，请检查日志: journalctl -u $SERVICE_NAME"
        exit 1
    fi

    echo ""
    echo "========================================="
    echo "  安装完成"
    echo "========================================="
    echo ""
    echo "  安装目录:   $INSTALL_DIR"
    echo "  配置文件:   $CONFIG_DIR/agent.yaml"
    echo "  数据目录:   $DATA_DIR"
    echo "  日志目录:   $LOG_DIR"
    echo "  服务名称:   $SERVICE_NAME"
    echo ""
    echo "  常用命令:"
    echo "    systemctl status $SERVICE_NAME"
    echo "    systemctl restart $SERVICE_NAME"
    echo "    journalctl -u $SERVICE_NAME -f"
    echo ""
}

main "$@"
