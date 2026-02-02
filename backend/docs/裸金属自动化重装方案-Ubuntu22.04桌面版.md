# Ubuntu 22.04 桌面版裸金属自动化重装方案

## 一、方案架构

### 1.1 整体架构

```
PXE服务器 ──┬── DHCP服务（分配IP、指定PXE服务器）
           ├── TFTP服务（提供启动文件）
           ├── HTTP服务（提供安装镜像和配置文件）
           └── NFS服务（可选，共享安装文件）

裸金属服务器 ── 网络启动 → 自动安装 → 自动配置 → 就绪
```

### 1.2 核心组件

- **DHCP服务器**：分配IP，指定PXE启动文件
- **TFTP服务器**：提供PXE启动文件（pxelinux.0、内核、initrd）
- **HTTP服务器**：提供Ubuntu安装镜像和preseed配置文件
- **Preseed文件**：自动化安装配置（无人值守安装）
- **后期配置脚本**：安装完成后的自动化配置

### 1.3 时间估算

- **首次部署PXE服务器**：2-4小时
- **制作preseed配置**：1-2小时
- **单台机器重装时间**：15-25分钟（取决于网络速度）

---

## 二、PXE服务器搭建

### 2.1 服务器要求

**硬件要求：**
- CPU：2核心
- 内存：4GB
- 磁盘：100GB（存储安装镜像）
- 网络：千兆网卡

**软件要求：**
- Ubuntu 22.04 Server
- 与裸金属服务器在同一局域网

### 2.2 安装必要软件

```bash
#!/bin/bash
# 安装PXE服务器所需软件

# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装DHCP服务器
sudo apt install -y isc-dhcp-server

# 安装TFTP服务器
sudo apt install -y tftpd-hpa

# 安装HTTP服务器
sudo apt install -y apache2

# 安装PXE启动相关工具
sudo apt install -y pxelinux syslinux-common

# 安装NFS服务器（可选）
sudo apt install -y nfs-kernel-server
```

### 2.3 配置DHCP服务器

**编辑 /etc/dhcp/dhcpd.conf：**

```bash
# 全局配置
default-lease-time 600;
max-lease-time 7200;
authoritative;

# 子网配置
subnet 192.168.1.0 netmask 255.255.255.0 {
    range 192.168.1.100 192.168.1.200;
    option routers 192.168.1.1;
    option domain-name-servers 8.8.8.8, 8.8.4.4;
    option subnet-mask 255.255.255.0;

    # PXE启动配置
    next-server 192.168.1.10;  # PXE服务器IP
    filename "pxelinux.0";      # PXE启动文件
}

# 为特定MAC地址分配固定IP（可选）
host gpu-server-01 {
    hardware ethernet 00:11:22:33:44:55;
    fixed-address 192.168.1.101;
}
```

**配置网络接口：**

编辑 /etc/default/isc-dhcp-server：
```bash
INTERFACESv4="ens33"  # 替换为你的网卡名称
```

**启动DHCP服务：**
```bash
sudo systemctl restart isc-dhcp-server
sudo systemctl enable isc-dhcp-server
sudo systemctl status isc-dhcp-server
```

### 2.4 配置TFTP服务器

**编辑 /etc/default/tftpd-hpa：**

```bash
TFTP_USERNAME="tftp"
TFTP_DIRECTORY="/var/lib/tftpboot"
TFTP_ADDRESS="0.0.0.0:69"
TFTP_OPTIONS="--secure"
```

**创建TFTP目录并设置权限：**

```bash
sudo mkdir -p /var/lib/tftpboot
sudo chmod 777 /var/lib/tftpboot
```

**复制PXE启动文件：**

```bash
# 复制pxelinux.0
sudo cp /usr/lib/PXELINUX/pxelinux.0 /var/lib/tftpboot/

# 复制syslinux模块
sudo cp /usr/lib/syslinux/modules/bios/*.c32 /var/lib/tftpboot/

# 创建pxelinux配置目录
sudo mkdir -p /var/lib/tftpboot/pxelinux.cfg
```

**启动TFTP服务：**

```bash
sudo systemctl restart tftpd-hpa
sudo systemctl enable tftpd-hpa
sudo systemctl status tftpd-hpa
```

### 2.5 下载Ubuntu安装镜像

```bash
# 创建HTTP服务目录
sudo mkdir -p /var/www/html/ubuntu

# 下载Ubuntu 22.04桌面版ISO
cd /tmp
wget https://releases.ubuntu.com/22.04/ubuntu-22.04.3-desktop-amd64.iso

# 挂载ISO
sudo mkdir -p /mnt/ubuntu-iso
sudo mount -o loop ubuntu-22.04.3-desktop-amd64.iso /mnt/ubuntu-iso

# 复制安装文件到HTTP目录
sudo cp -r /mnt/ubuntu-iso/* /var/www/html/ubuntu/

# 提取内核和initrd到TFTP目录
sudo cp /mnt/ubuntu-iso/casper/vmlinuz /var/lib/tftpboot/
sudo cp /mnt/ubuntu-iso/casper/initrd /var/lib/tftpboot/

# 卸载ISO
sudo umount /mnt/ubuntu-iso
```

### 2.6 配置PXE启动菜单

**创建 /var/lib/tftpboot/pxelinux.cfg/default：**

```bash
DEFAULT ubuntu-desktop
TIMEOUT 100
PROMPT 0

LABEL ubuntu-desktop
    MENU LABEL Install Ubuntu 22.04 Desktop (Auto)
    KERNEL vmlinuz
    APPEND initrd=initrd boot=casper url=http://192.168.1.10/ubuntu/ubuntu-22.04.3-desktop-amd64.iso autoinstall ds=nocloud-net;s=http://192.168.1.10/preseed/ ---
```

---

## 三、Preseed自动化配置

### 3.1 创建Preseed目录

```bash
sudo mkdir -p /var/www/html/preseed
```

### 3.2 创建user-data文件

**创建 /var/www/html/preseed/user-data：**

```yaml
#cloud-config
autoinstall:
  version: 1

  # 语言和键盘
  locale: zh_CN.UTF-8
  keyboard:
    layout: us

  # 网络配置（使用DHCP）
  network:
    network:
      version: 2
      ethernets:
        any:
          match:
            name: en*
          dhcp4: true

  # 存储配置（自动分区）
  storage:
    layout:
      name: lvm

  # 用户配置
  identity:
    hostname: remotegpu-server
    username: remotegpu
    password: "$6$rounds=4096$saltsalt$hashed_password"  # 需要替换为实际的加密密码
    realname: RemoteGPU User

  # SSH配置
  ssh:
    install-server: true
    allow-pw: true

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
    - nvidia-driver-535  # NVIDIA驱动
    - docker.io
    - python3-pip

  # 后期配置脚本
  late-commands:
    - curtin in-target --target=/target -- wget -O /tmp/post-install.sh http://192.168.1.10/preseed/post-install.sh
    - curtin in-target --target=/target -- chmod +x /tmp/post-install.sh
    - curtin in-target --target=/target -- /tmp/post-install.sh
```

**生成加密密码：**

```bash
# 生成密码哈希（密码：remotegpu123）
python3 -c 'import crypt; print(crypt.crypt("remotegpu123", crypt.mksalt(crypt.METHOD_SHA512)))'

# 将输出的哈希值替换到user-data中的password字段
```

### 3.3 创建meta-data文件

**创建 /var/www/html/preseed/meta-data：**

```yaml
instance-id: remotegpu-server
local-hostname: remotegpu-server
```

### 3.4 创建后期配置脚本

**创建 /var/www/html/preseed/post-install.sh：**

```bash
#!/bin/bash
# Ubuntu 22.04桌面版后期配置脚本

set -e

echo "=== 开始后期配置 ==="

# 1. 配置自动登录
mkdir -p /etc/gdm3
cat > /etc/gdm3/custom.conf << 'EOF'
[daemon]
AutomaticLoginEnable=true
AutomaticLogin=remotegpu
EOF

# 2. 禁用屏幕锁定
sudo -u remotegpu dbus-launch gsettings set org.gnome.desktop.screensaver lock-enabled false
sudo -u remotegpu dbus-launch gsettings set org.gnome.desktop.session idle-delay 0

# 3. 配置Docker
usermod -aG docker remotegpu
systemctl enable docker
systemctl start docker

# 4. 安装Jupyter
pip3 install jupyterlab notebook

# 5. 安装Code Server
curl -fsSL https://code-server.dev/install.sh | sh
systemctl enable --now code-server@remotegpu

# 6. 配置SSH
sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config
sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/' /etc/ssh/sshd_config
systemctl restart sshd

# 7. 配置防火墙
ufw allow 22/tcp
ufw allow 8888/tcp
ufw allow 8080/tcp
ufw --force enable

# 8. 创建工作目录
mkdir -p /home/remotegpu/workspace
chown -R remotegpu:remotegpu /home/remotegpu/workspace

# 9. 配置Jupyter自动启动
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

# 10. 清理安装文件
apt clean
rm -rf /tmp/*

echo "=== 后期配置完成 ==="
```

**设置脚本权限：**

```bash
sudo chmod +x /var/www/html/preseed/post-install.sh
```

---

## 四、裸金属服务器配置

### 4.1 BIOS设置

**进入BIOS设置（通常按F2或Del）：**

1. **启动顺序设置：**
   - 将"Network Boot"或"PXE Boot"设置为第一启动项
   - 或者设置为"Hard Disk → Network"

2. **网络启动设置：**
   - 启用"PXE Boot"
   - 启用"Network Stack"
   - 选择正确的网卡（如果有多个）

3. **保存并退出**

### 4.2 网络启动流程

**启动过程：**

```
1. 服务器开机
   ↓
2. BIOS检测到网络启动
   ↓
3. 发送DHCP请求
   ↓
4. DHCP服务器响应（分配IP + PXE服务器地址）
   ↓
5. 从TFTP服务器下载pxelinux.0
   ↓
6. 加载PXE启动菜单
   ↓
7. 自动选择"Install Ubuntu 22.04 Desktop (Auto)"
   ↓
8. 下载内核和initrd
   ↓
9. 启动安装程序
   ↓
10. 从HTTP服务器下载preseed配置
   ↓
11. 自动安装系统（无人值守）
   ↓
12. 执行后期配置脚本
   ↓
13. 重启进入新系统
```

---

## 五、管理和维护

### 5.1 快速重装脚本

**创建 /usr/local/bin/reinstall-server.sh：**

```bash
#!/bin/bash
# 远程触发服务器重装

SERVER_IP=$1
SERVER_MAC=$2

if [ -z "$SERVER_IP" ] || [ -z "$SERVER_MAC" ]; then
    echo "用法: $0 <服务器IP> <服务器MAC地址>"
    exit 1
fi

echo "准备重装服务器: $SERVER_IP ($SERVER_MAC)"

# 1. 通过SSH触发重启到网络启动
ssh remotegpu@$SERVER_IP "sudo reboot" || true

# 2. 等待服务器关机
sleep 10

# 3. 使用Wake-on-LAN唤醒服务器（如果支持）
# wakeonlan $SERVER_MAC

echo "服务器正在重装，预计15-25分钟完成"
echo "可以通过以下命令监控进度："
echo "  ping $SERVER_IP"
echo "  ssh remotegpu@$SERVER_IP"
```

### 5.2 监控安装进度

```bash
# 查看DHCP日志
sudo tail -f /var/log/syslog | grep dhcpd

# 查看TFTP日志
sudo tail -f /var/log/syslog | grep tftpd

# 查看HTTP访问日志
sudo tail -f /var/log/apache2/access.log
```

### 5.3 故障排查

**问题1：服务器无法网络启动**

```bash
# 检查DHCP服务
sudo systemctl status isc-dhcp-server
sudo journalctl -u isc-dhcp-server -f

# 检查TFTP服务
sudo systemctl status tftpd-hpa
sudo journalctl -u tftpd-hpa -f

# 测试TFTP连接
tftp 192.168.1.10
> get pxelinux.0
> quit
```

**问题2：安装过程卡住**

```bash
# 检查HTTP服务
sudo systemctl status apache2

# 检查preseed文件
curl http://192.168.1.10/preseed/user-data
curl http://192.168.1.10/preseed/meta-data

# 检查后期配置脚本
curl http://192.168.1.10/preseed/post-install.sh
```

**问题3：安装完成但配置不正确**

```bash
# 查看安装日志（在目标服务器上）
sudo cat /var/log/installer/autoinstall-user-data
sudo cat /var/log/cloud-init.log
sudo cat /var/log/cloud-init-output.log
```

---

## 六、优化建议

### 6.1 加速安装

**使用本地镜像源：**

在user-data中添加：
```yaml
apt:
  primary:
    - arches: [default]
      uri: http://mirrors.aliyun.com/ubuntu/
```

**使用SSD存储安装文件：**
```bash
# 将/var/www/html挂载到SSD
sudo mkdir -p /mnt/ssd/www
sudo mount /dev/nvme0n1p1 /mnt/ssd/www
sudo ln -s /mnt/ssd/www /var/www/html
```

### 6.2 批量管理

**创建服务器清单：**

```bash
# servers.txt
192.168.1.101,00:11:22:33:44:55,gpu-server-01
192.168.1.102,00:11:22:33:44:56,gpu-server-02
192.168.1.103,00:11:22:33:44:57,gpu-server-03
```

**批量重装脚本：**

```bash
#!/bin/bash
# 批量重装服务器

while IFS=',' read -r ip mac hostname; do
    echo "重装服务器: $hostname ($ip)"
    ./reinstall-server.sh $ip $mac &
    sleep 5
done < servers.txt

wait
echo "所有服务器重装任务已启动"
```

### 6.3 自定义配置

**根据MAC地址使用不同配置：**

在DHCP配置中：
```bash
host gpu-server-01 {
    hardware ethernet 00:11:22:33:44:55;
    fixed-address 192.168.1.101;
    filename "pxelinux.cfg/gpu-server-01";
}
```

创建专用配置：
```bash
# /var/lib/tftpboot/pxelinux.cfg/gpu-server-01
DEFAULT ubuntu-desktop
LABEL ubuntu-desktop
    KERNEL vmlinuz
    APPEND initrd=initrd autoinstall ds=nocloud-net;s=http://192.168.1.10/preseed/gpu-server-01/
```

---

## 七、集成到RemoteGPU平台

### 7.1 API接口设计

```go
// ReinstallServer 重装服务器
func (s *HostService) ReinstallServer(hostID string) error {
    host, err := s.hostDao.GetByID(hostID)
    if err != nil {
        return err
    }

    // 1. 标记服务器状态为"重装中"
    host.Status = "reinstalling"
    s.hostDao.Update(host)

    // 2. 触发远程重启
    err = s.remoteManager.RebootToNetwork(host.IPAddress)
    if err != nil {
        return err
    }

    // 3. 启动监控任务
    go s.monitorReinstall(hostID)

    return nil
}

// monitorReinstall 监控重装进度
func (s *HostService) monitorReinstall(hostID string) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    timeout := time.After(30 * time.Minute)

    for {
        select {
        case <-ticker.C:
            // 检查服务器是否可访问
            if s.checkServerReady(hostID) {
                // 更新状态为"可用"
                s.updateHostStatus(hostID, "active")
                return
            }
        case <-timeout:
            // 超时，标记为失败
            s.updateHostStatus(hostID, "reinstall_failed")
            return
        }
    }
}
```

### 7.2 前端界面

**服务器管理页面添加"重装系统"按钮：**

```vue
<template>
  <el-button
    type="danger"
    @click="handleReinstall(server.id)"
    :loading="reinstalling">
    重装系统
  </el-button>
</template>

<script>
export default {
  methods: {
    async handleReinstall(serverId) {
      const confirm = await this.$confirm(
        '重装系统将清除所有数据，是否继续？',
        '警告',
        { type: 'warning' }
      )

      if (confirm) {
        this.reinstalling = true
        try {
          await api.reinstallServer(serverId)
          this.$message.success('重装任务已启动，预计15-25分钟完成')
        } catch (error) {
          this.$message.error('重装失败：' + error.message)
        } finally {
          this.reinstalling = false
        }
      }
    }
  }
}
</script>
```

---

## 八、总结

### 8.1 方案优势

✅ **完全自动化**：无需人工干预，15-25分钟完成重装
✅ **标准化配置**：每次重装后配置完全一致
✅ **易于维护**：修改preseed文件即可更新配置
✅ **成本低**：使用开源软件，无额外费用

### 8.2 实施步骤

1. **搭建PXE服务器**（2-4小时）
2. **配置preseed文件**（1-2小时）
3. **测试单台服务器**（30分钟）
4. **批量部署**（根据服务器数量）

### 8.3 注意事项

⚠️ **网络要求**：PXE服务器和裸金属服务器必须在同一局域网
⚠️ **BIOS设置**：确保启用网络启动
⚠️ **备份数据**：重装会清除所有数据
⚠️ **测试环境**：先在测试环境验证配置

### 8.4 后续优化

- 添加多版本支持（Ubuntu 20.04、22.04、24.04）
- 集成GPU驱动自动安装
- 添加RAID配置支持
- 实现Web管理界面
