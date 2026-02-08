# 第一阶段：编译 Agent（build context 为项目根目录）
FROM golang:1.23-alpine AS agent-builder

WORKDIR /build

COPY agent/go.mod agent/go.sum ./
RUN go mod download

COPY agent/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o remotegpu-agent ./cmd/

# 第二阶段：测试环境
FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

# 安装基础工具和 SSH
RUN apt-get update && apt-get install -y --no-install-recommends \
    openssh-server \
    python3 \
    python3-pip \
    vim \
    curl \
    wget \
    git \
    htop \
    net-tools \
    iproute2 \
    procps \
    sudo \
    && rm -rf /var/lib/apt/lists/*

# 安装 Jupyter
RUN pip3 install --no-cache-dir \
    jupyterlab \
    notebook \
    numpy \
    pandas

# 创建目录
RUN mkdir -p /workspace /run/sshd /opt/remotegpu-agent

# 从构建阶段复制 Agent 二进制
COPY --from=agent-builder /build/remotegpu-agent /opt/remotegpu-agent/remotegpu-agent

# 复制启动脚本（build context 为项目根目录）
COPY docker-compose/test-env/start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE 22 8888 8090

CMD ["/start.sh"]
