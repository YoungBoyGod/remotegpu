# Preseed自动化安装配置

## 一、Preseed简介

### 1.1 什么是Preseed

Preseed是Ubuntu/Debian的自动化安装配置系统，允许预先配置安装过程中的所有选项，实现无人值守安装。

**Ubuntu 22.04使用的是Autoinstall（基于cloud-init）：**
- 配置文件格式：YAML
- 配置文件名：user-data、meta-data
- 功能更强大，支持更多自定义

### 1.2 Preseed的作用

在PXE启动中，Preseed负责：
1. 自动分区和格式化磁盘
2. 配置网络和主机名
3. 创建用户和设置密码
4. 选择要安装的软件包
5. 执行后期配置脚本

### 1.3 工作流程

```
PXE启动 → 加载内核和initrd
        → 下载user-data和meta-data
        → 自动安装系统
        → 执行late-commands
        → 重启进入新系统
```

---

## 二、创建Preseed目录

### 2.1 目录结构

```bash
# 创建preseed目录
sudo mkdir -p /var/www/html/preseed

# 创建子目录（可选，用于不同配置）
sudo mkdir -p /var/www/html/preseed/default
sudo mkdir -p /var/www/html/preseed/gpu-server
sudo mkdir -p /var/www/html/preseed/workstation
```

### 2.2 设置权限

```bash
# 设置目录权限
sudo chown -R www-data:www-data /var/www/html/preseed
sudo chmod -R 755 /var/www/html/preseed
```

---

## 三、创建user-data配置文件

### 3.1 基础配置

**创建 /var/www/html/preseed/user-data：**

```yaml
#cloud-config
autoinstall:
  version: 1

  # 语言和键盘
  locale: zh_CN.UTF-8
  keyboard:
    layout: us
    variant: ''

  # 网络配置（使用DHCP）
  network:
    network:
      version: 2
      ethernets:
        any:
          match:
            name: en*
          dhcp4: true
          dhcp6: false

  # 存储配置（自动分区）
  storage:
    layout:
      name: lvm
      match:
        size: largest

  # 用户配置
  identity:
    hostname: remotegpu-server
    username: remotegpu
    password: "$6$rounds=4096$saltsalt$hashed_password"
    realname: RemoteGPU User

  # SSH配置
  ssh:
    install-server: true
    allow-pw: true
    authorized-keys: []

  # 软件包
  packages:
    - ubuntu-desktop
    - openssh-server
    - vim
    - curl
    - wget
    - git
    - htop
    - net-tools

  # 后期配置
  late-commands:
    - echo 'Installation completed' > /target/var/log/autoinstall.log
```

### 3.2 生成加密密码

**使用Python生成密码哈希：**

```bash
# 方法1：使用Python
python3 -c 'import crypt; print(crypt.crypt("remotegpu123", crypt.mksalt(crypt.METHOD_SHA512)))'

# 方法2：使用mkpasswd
sudo apt install whois
mkpasswd -m sha-512 remotegpu123

# 输出示例：
# $6$rounds=4096$saltsalt$hashed_password_string
```

**将生成的哈希值替换到user-data中的password字段。**

### 3.3 完整配置示例

**创建 /var/www/html/preseed/user-data（完整版）：**

```yaml
#cloud-config
autoinstall:
  version: 1

  # 交互模式（设置为false表示完全自动）
  interactive-sections: []

  # 语言和键盘
  locale: zh_CN.UTF-8
  keyboard:
    layout: us
    variant: ''
    toggle: null

  # 时区
  timezone: Asia/Shanghai

  # 网络配置
  network:
    network:
      version: 2
      ethernets:
        any:
          match:
            name: en*
          dhcp4: true
          dhcp6: false
          dhcp-identifier: mac

  # APT镜像源（使用阿里云镜像加速）
  apt:
    primary:
      - arches: [default]
        uri: http://mirrors.aliyun.com/ubuntu/
    security:
      - arches: [default]
        uri: http://mirrors.aliyun.com/ubuntu/

  # 存储配置
  storage:
    layout:
      name: lvm
      match:
        size: largest
        ssd: true
    swap:
      size: 0  # 禁用swap（GPU服务器推荐）

  # 用户配置
  identity:
    hostname: remotegpu-server
    username: remotegpu
    password: "$6$rounds=4096$saltsalt$hashed_password"
    realname: RemoteGPU User

  # SSH配置
  ssh:
    install-server: true
    allow-pw: true
    authorized-keys: []

  # 软件包
  packages:
    # 桌面环境
    - ubuntu-desktop
    - gnome-shell
    - gnome-terminal

    # 开发工具
    - build-essential
    - git
    - vim
    - curl
    - wget

    # 系统工具
    - openssh-server
    - net-tools
    - htop
    - tmux
    - screen

    # Python环境
    - python3
    - python3-pip
    - python3-dev

    # Docker
    - docker.io
    - docker-compose

  # 后期配置脚本
  late-commands:
    # 下载并执行后期配置脚本
    - curtin in-target --target=/target -- wget -O /tmp/post-install.sh http://192.168.1.10/preseed/post-install.sh
    - curtin in-target --target=/target -- chmod +x /tmp/post-install.sh
    - curtin in-target --target=/target -- /tmp/post-install.sh

    # 配置自动登录
    - curtin in-target --target=/target -- mkdir -p /etc/gdm3
    - echo -e "[daemon]\nAutomaticLoginEnable=true\nAutomaticLogin=remotegpu" > /target/etc/gdm3/custom.conf

    # 禁用自动更新
    - curtin in-target --target=/target -- systemctl disable apt-daily.timer
    - curtin in-target --target=/target -- systemctl disable apt-daily-upgrade.timer
```

---

## 四、创建meta-data配置文件

### 4.1 基础meta-data

**创建 /var/www/html/preseed/meta-data：**

```yaml
instance-id: remotegpu-server-001
local-hostname: remotegpu-server
```

### 4.2 高级meta-data配置

**包含网络配置的meta-data：**

```yaml
instance-id: remotegpu-server-001
local-hostname: remotegpu-server

# 网络配置（可选，如果需要固定IP）
network-interfaces: |
  auto eth0
  iface eth0 inet static
    address 192.168.1.101
    netmask 255.255.255.0
    gateway 192.168.1.1
    dns-nameservers 8.8.8.8 8.8.4.4
```

---

## 五、存储配置详解

### 5.1 自动分区（LVM）

**简单LVM配置：**

```yaml
storage:
  layout:
    name: lvm
    match:
      size: largest
```

### 5.2 自定义分区

**手动指定分区方案：**

```yaml
storage:
  config:
    # 选择磁盘
    - type: disk
      id: disk0
      path: /dev/sda
      wipe: superblock

    # 创建分区表
    - type: partition
      id: boot-partition
      device: disk0
      size: 1GB
      flag: boot

    - type: partition
      id: root-partition
      device: disk0
      size: -1  # 使用剩余空间

    # 格式化
    - type: format
      id: boot-fs
      volume: boot-partition
      fstype: ext4

    - type: format
      id: root-fs
      volume: root-partition
      fstype: ext4

    # 挂载点
    - type: mount
      id: boot-mount
      device: boot-fs
      path: /boot

    - type: mount
      id: root-mount
      device: root-fs
      path: /
```

### 5.3 RAID配置

**RAID 1配置示例：**

```yaml
storage:
  config:
    # 磁盘1
    - type: disk
      id: disk0
      path: /dev/sda
      wipe: superblock

    # 磁盘2
    - type: disk
      id: disk1
      path: /dev/sdb
      wipe: superblock

    # 创建RAID
    - type: raid
      id: md0
      name: md0
      raidlevel: 1
      devices:
        - disk0
        - disk1
      spare_devices: []

    # 在RAID上创建分区
    - type: partition
      id: raid-partition
      device: md0
      size: -1

    # 格式化
    - type: format
      id: raid-fs
      volume: raid-partition
      fstype: ext4

    # 挂载
    - type: mount
      id: raid-mount
      device: raid-fs
      path: /
```

---

## 六、软件包配置

### 6.1 基础软件包

```yaml
packages:
  # 最小化安装
  - openssh-server
  - vim
  - curl
  - wget
```

### 6.2 桌面环境

```yaml
packages:
  # Ubuntu桌面
  - ubuntu-desktop

  # 或者使用轻量级桌面
  - xubuntu-desktop  # XFCE
  - lubuntu-desktop  # LXQt
```

### 6.3 开发环境

```yaml
packages:
  # 编译工具
  - build-essential
  - gcc
  - g++
  - make
  - cmake

  # 版本控制
  - git
  - subversion

  # Python
  - python3
  - python3-pip
  - python3-dev
  - python3-venv

  # Node.js
  - nodejs
  - npm
```

### 6.4 GPU相关

```yaml
packages:
  # NVIDIA驱动（需要添加PPA）
  - nvidia-driver-535
  - nvidia-utils-535

  # CUDA工具包（需要额外配置）
  # 通常在late-commands中安装
```

---

## 七、后期配置脚本

### 7.1 创建post-install.sh

**创建 /var/www/html/preseed/post-install.sh：**

```bash
#!/bin/bash
# Ubuntu 22.04 Desktop 后期配置脚本

set -e

echo "=== 开始后期配置 ==="

# 1. 配置Docker
echo "配置Docker..."
usermod -aG docker remotegpu
systemctl enable docker
systemctl start docker

# 2. 安装Jupyter Lab
echo "安装Jupyter Lab..."
pip3 install --upgrade pip
pip3 install jupyterlab notebook ipywidgets

# 3. 安装Code Server
echo "安装Code Server..."
curl -fsSL https://code-server.dev/install.sh | sh
systemctl enable --now code-server@remotegpu

# 4. 配置SSH
echo "配置SSH..."
sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config
sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/' /etc/ssh/sshd_config
systemctl restart sshd

# 5. 配置防火墙
echo "配置防火墙..."
ufw allow 22/tcp
ufw allow 8888/tcp
ufw allow 8080/tcp
ufw --force enable

# 6. 禁用屏幕锁定
echo "禁用屏幕锁定..."
sudo -u remotegpu dbus-launch gsettings set org.gnome.desktop.screensaver lock-enabled false
sudo -u remotegpu dbus-launch gsettings set org.gnome.desktop.session idle-delay 0

# 7. 创建工作目录
echo "创建工作目录..."
mkdir -p /home/remotegpu/workspace
chown -R remotegpu:remotegpu /home/remotegpu/workspace

# 8. 配置Jupyter自动启动
echo "配置Jupyter自动启动..."
cat > /etc/systemd/system/jupyter.service << 'EOF'
[Unit]
Description=Jupyter Lab
After=network.target

[Service]
Type=simple
User=remotegpu
WorkingDirectory=/home/remotegpu/workspace
ExecStart=/usr/local/bin/jupyter lab --ip=0.0.0.0 --port=8888 --no-browser --allow-root
Restart=always

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable jupyter
systemctl start jupyter

# 9. 清理
echo "清理安装文件..."
apt clean
rm -rf /tmp/*

echo "=== 后期配置完成 ==="
```

**设置脚本权限：**

```bash
sudo chmod +x /var/www/html/preseed/post-install.sh
```

### 7.2 高级后期配置

**包含NVIDIA驱动安装：**

```bash
#!/bin/bash
set -e

echo "=== 安装NVIDIA驱动 ==="

# 添加NVIDIA PPA
add-apt-repository -y ppa:graphics-drivers/ppa
apt update

# 安装驱动
apt install -y nvidia-driver-535 nvidia-utils-535

# 安装CUDA（可选）
wget https://developer.download.nvidia.com/compute/cuda/repos/ubuntu2204/x86_64/cuda-ubuntu2204.pin
mv cuda-ubuntu2204.pin /etc/apt/preferences.d/cuda-repository-pin-600
wget https://developer.download.nvidia.com/compute/cuda/12.0.0/local_installers/cuda-repo-ubuntu2204-12-0-local_12.0.0-525.60.13-1_amd64.deb
dpkg -i cuda-repo-ubuntu2204-12-0-local_12.0.0-525.60.13-1_amd64.deb
cp /var/cuda-repo-ubuntu2204-12-0-local/cuda-*-keyring.gpg /usr/share/keyrings/
apt update
apt install -y cuda

echo "=== NVIDIA驱动安装完成 ==="
```

---

## 八、测试和验证

### 8.1 验证配置文件

**检查YAML语法：**

```bash
# 安装yamllint
sudo apt install yamllint

# 验证user-data
yamllint /var/www/html/preseed/user-data

# 验证meta-data
yamllint /var/www/html/preseed/meta-data
```

### 8.2 测试HTTP访问

```bash
# 测试user-data
curl http://192.168.1.10/preseed/user-data

# 测试meta-data
curl http://192.168.1.10/preseed/meta-data

# 测试post-install.sh
curl http://192.168.1.10/preseed/post-install.sh
```

### 8.3 验证PXE配置

**检查PXE启动菜单是否正确引用preseed：**

```bash
cat /var/lib/tftpboot/pxelinux.cfg/default

# 应该包含：
# APPEND initrd=initrd boot=casper ... ds=nocloud-net;s=http://192.168.1.10/preseed/
```

---

## 九、常见配置场景

### 9.1 GPU服务器配置

**user-data示例：**

```yaml
autoinstall:
  version: 1

  identity:
    hostname: gpu-server-01
    username: remotegpu
    password: "$6$..."

  storage:
    layout:
      name: lvm
    swap:
      size: 0  # GPU服务器不需要swap

  packages:
    - ubuntu-desktop
    - openssh-server
    - nvidia-driver-535
    - docker.io

  late-commands:
    - curtin in-target -- wget -O /tmp/gpu-setup.sh http://192.168.1.10/preseed/gpu-setup.sh
    - curtin in-target -- bash /tmp/gpu-setup.sh
```

### 9.2 工作站配置

**user-data示例：**

```yaml
autoinstall:
  version: 1

  identity:
    hostname: workstation-01
    username: remotegpu
    password: "$6$..."

  packages:
    - ubuntu-desktop
    - openssh-server
    - build-essential
    - git
    - vim
    - code  # VS Code

  late-commands:
    - curtin in-target -- wget -O /tmp/workstation-setup.sh http://192.168.1.10/preseed/workstation-setup.sh
    - curtin in-target -- bash /tmp/workstation-setup.sh
```

### 9.3 最小化服务器配置

**user-data示例：**

```yaml
autoinstall:
  version: 1

  identity:
    hostname: minimal-server
    username: remotegpu
    password: "$6$..."

  packages:
    - openssh-server
    - vim
    - curl

  # 不安装桌面环境
  # 最小化安装
```

---

## 十、故障排查

### 10.1 安装卡住

**问题：安装过程中卡在某个步骤**

```bash
# 在安装过程中按Alt+F2切换到控制台
# 查看安装日志
tail -f /var/log/installer/autoinstall-user-data

# 查看cloud-init日志
tail -f /var/log/cloud-init.log
tail -f /var/log/cloud-init-output.log
```

### 10.2 配置文件错误

**问题：user-data配置错误导致安装失败**

```bash
# 检查HTTP服务器日志
sudo tail -f /var/log/apache2/access.log

# 确认文件可访问
curl -I http://192.168.1.10/preseed/user-data

# 验证YAML语法
yamllint /var/www/html/preseed/user-data
```

### 10.3 密码无法登录

**问题：安装完成后无法使用配置的密码登录**

```bash
# 重新生成密码哈希
python3 -c 'import crypt; print(crypt.crypt("your_password", crypt.mksalt(crypt.METHOD_SHA512)))'

# 确认密码哈希格式正确
# 应该是：$6$rounds=4096$...

# 检查user-data中的密码字段
grep password /var/www/html/preseed/user-data
```

### 10.4 后期脚本失败

**问题：late-commands执行失败**

```bash
# 查看安装日志
cat /var/log/installer/autoinstall-user-data

# 查看curtin日志
cat /var/log/installer/curtin-install.log

# 检查脚本是否可下载
curl http://192.168.1.10/preseed/post-install.sh

# 检查脚本权限
ls -l /var/www/html/preseed/post-install.sh
```

---

## 十一、优化建议

### 11.1 加速安装

**使用本地镜像源：**

```yaml
apt:
  primary:
    - arches: [default]
      uri: http://192.168.1.10/ubuntu-mirror/
```

**预下载软件包：**

```bash
# 在PXE服务器上创建本地apt缓存
sudo apt install apt-cacher-ng

# 在user-data中配置
apt:
  http_proxy: http://192.168.1.10:3142
```

### 11.2 减少安装时间

**最小化软件包：**

```yaml
packages:
  # 只安装必要的包
  - openssh-server
  - vim

# 在late-commands中安装其他软件
late-commands:
  - curtin in-target -- apt install -y ubuntu-desktop
```

### 11.3 安全加固

**禁用root登录：**

```yaml
late-commands:
  - curtin in-target -- passwd -l root
  - curtin in-target -- sed -i 's/#PermitRootLogin.*/PermitRootLogin no/' /etc/ssh/sshd_config
```

**配置防火墙：**

```yaml
late-commands:
  - curtin in-target -- ufw default deny incoming
  - curtin in-target -- ufw default allow outgoing
  - curtin in-target -- ufw allow 22/tcp
  - curtin in-target -- ufw --force enable
```

---

## 十二、总结

**Preseed配置要点：**
1. ✅ 正确配置user-data和meta-data
2. ✅ 生成安全的密码哈希
3. ✅ 选择合适的存储布局
4. ✅ 配置必要的软件包
5. ✅ 编写后期配置脚本

**下一步：** 配置后期脚本和故障排查

