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

# 复制启动脚本
COPY start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE 22 8888 8090

CMD ["/start.sh"]
