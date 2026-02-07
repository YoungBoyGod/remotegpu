FROM python:3.11-slim

# 安装 SSH 和基础工具
RUN apt-get update && apt-get install -y --no-install-recommends \
    openssh-server \
    vim \
    curl \
    wget \
    git \
    htop \
    net-tools \
    && rm -rf /var/lib/apt/lists/*

# 安装 Jupyter
RUN pip install --no-cache-dir \
    jupyterlab \
    notebook \
    numpy \
    pandas

# 创建工作目录
RUN mkdir -p /workspace /run/sshd

# 复制启动脚本
COPY start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE 22 8888

CMD ["/start.sh"]
