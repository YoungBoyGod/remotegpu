#!/bin/bash
# RemoteGPU Agent 升级脚本
# 用法: sudo bash upgrade.sh --download-url <URL或本地路径>

set -euo pipefail

INSTALL_DIR="/opt/remotegpu-agent"
BINARY_NAME="remotegpu-agent"
SERVICE_NAME="remotegpu-agent"
BACKUP_DIR="/opt/remotegpu-agent/backup"

DOWNLOAD_URL=""

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
            --download-url) DOWNLOAD_URL="$2"; shift 2 ;;
            -h|--help)
                echo "用法: sudo bash upgrade.sh --download-url <URL或本地路径>"
                exit 0 ;;
            *) log_error "未知参数: $1"; exit 1 ;;
        esac
    done
}

# ========== 前置检查 ==========
check_prerequisites() {
    if [[ $EUID -ne 0 ]]; then
        log_error "请使用 root 权限运行此脚本"
        exit 1
    fi
    if [[ -z "$DOWNLOAD_URL" ]]; then
        log_error "--download-url 参数必填"
        exit 1
    fi
    if [[ ! -f "$INSTALL_DIR/$BINARY_NAME" ]]; then
        log_error "未找到已安装的 Agent，请先运行 install.sh"
        exit 1
    fi
}

# ========== 备份当前版本 ==========
backup_current() {
    local timestamp
    timestamp=$(date +%Y%m%d_%H%M%S)
    mkdir -p "$BACKUP_DIR"

    log_info "备份当前二进制文件..."
    cp "$INSTALL_DIR/$BINARY_NAME" "$BACKUP_DIR/${BINARY_NAME}.${timestamp}"

    # 只保留最近 5 个备份
    local count
    count=$(ls -1 "$BACKUP_DIR"/${BINARY_NAME}.* 2>/dev/null | wc -l)
    if [[ $count -gt 5 ]]; then
        ls -1t "$BACKUP_DIR"/${BINARY_NAME}.* | tail -n +6 | xargs rm -f
        log_info "清理旧备份，保留最近 5 个"
    fi

    log_info "备份完成: $BACKUP_DIR/${BINARY_NAME}.${timestamp}"
}

# ========== 下载新版本 ==========
download_new_binary() {
    local tmp_file="/tmp/${BINARY_NAME}.new"

    if [[ -f "$DOWNLOAD_URL" ]]; then
        log_info "从本地路径复制: $DOWNLOAD_URL"
        cp "$DOWNLOAD_URL" "$tmp_file"
    else
        log_info "从远程下载: $DOWNLOAD_URL"
        curl -fSL -o "$tmp_file" "$DOWNLOAD_URL"
    fi

    chmod +x "$tmp_file"
    echo "$tmp_file"
}

# ========== 执行升级 ==========
do_upgrade() {
    local tmp_file="$1"

    log_info "停止 Agent 服务..."
    systemctl stop "$SERVICE_NAME"

    log_info "替换二进制文件..."
    mv "$tmp_file" "$INSTALL_DIR/$BINARY_NAME"

    log_info "启动 Agent 服务..."
    systemctl start "$SERVICE_NAME"

    sleep 2
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        log_info "升级成功，Agent 已启动"
    else
        log_error "升级后启动失败，尝试回滚..."
        rollback
    fi
}

# ========== 回滚 ==========
rollback() {
    local latest_backup
    latest_backup=$(ls -1t "$BACKUP_DIR"/${BINARY_NAME}.* 2>/dev/null | head -1)

    if [[ -z "$latest_backup" ]]; then
        log_error "无可用备份，回滚失败"
        exit 1
    fi

    log_warn "回滚到: $latest_backup"
    systemctl stop "$SERVICE_NAME" 2>/dev/null || true
    cp "$latest_backup" "$INSTALL_DIR/$BINARY_NAME"
    systemctl start "$SERVICE_NAME"

    if systemctl is-active --quiet "$SERVICE_NAME"; then
        log_info "回滚成功"
    else
        log_error "回滚后仍无法启动，请手动排查: journalctl -u $SERVICE_NAME"
        exit 1
    fi
}

# ========== 主流程 ==========
main() {
    echo ""
    echo "========================================="
    echo "  RemoteGPU Agent 升级程序"
    echo "========================================="
    echo ""

    parse_args "$@"
    check_prerequisites
    backup_current

    local tmp_file
    tmp_file=$(download_new_binary)
    do_upgrade "$tmp_file"

    echo ""
    log_info "升级完成"
    echo ""
}

main "$@"
