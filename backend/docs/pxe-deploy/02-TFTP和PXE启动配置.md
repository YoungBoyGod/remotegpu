# TFTP和PXE启动配置

## 一、TFTP服务器作用

TFTP（Trivial File Transfer Protocol）在PXE启动中的作用：
1. 提供PXE启动文件（pxelinux.0）
2. 提供Linux内核（vmlinuz）
3. 提供初始化内存盘（initrd）
4. 提供启动菜单配置

## 二、安装TFTP服务器

```bash
# 安装tftpd-hpa
sudo apt update
sudo apt install -y tftpd-hpa

# 安装PXE启动相关工具
sudo apt install -y pxelinux syslinux-common
```

## 三、配置TFTP服务器

### 3.1 编辑配置文件

**编辑 /etc/default/tftpd-hpa：**

```bash
sudo vim /etc/default/tftpd-hpa
```

**配置内容：**

```conf
# TFTP用户
TFTP_USERNAME="tftp"

# TFTP根目录
TFTP_DIRECTORY="/var/lib/tftpboot"

# 监听地址和端口
TFTP_ADDRESS="0.0.0.0:69"

# TFTP选项
TFTP_OPTIONS="--secure --create"
```

**参数说明：**
- `--secure`: 限制访问TFTP根目录
- `--create`: 允许创建新文件（用于日志等）

### 3.2 创建TFTP目录

```bash
# 创建TFTP根目录
sudo mkdir -p /var/lib/tftpboot

# 设置权限
sudo chown -R tftp:tftp /var/lib/tftpboot
sudo chmod -R 755 /var/lib/tftpboot
```

### 3.3 启动TFTP服务

```bash
# 启动服务
sudo systemctl start tftpd-hpa

# 设置开机自启
sudo systemctl enable tftpd-hpa

# 查看服务状态
sudo systemctl status tftpd-hpa
```

## 四、准备PXE启动文件

### 4.1 复制PXE引导程序

```bash
# 复制pxelinux.0（BIOS启动）
sudo cp /usr/lib/PXELINUX/pxelinux.0 /var/lib/tftpboot/

# 复制UEFI启动文件（如果需要支持UEFI）
sudo cp /usr/lib/syslinux/modules/efi64/syslinux.efi /var/lib/tftpboot/

# 复制syslinux模块
sudo cp /usr/lib/syslinux/modules/bios/*.c32 /var/lib/tftpboot/
```

**重要的.c32模块：**
- `ldlinux.c32`: 核心模块
- `libcom32.c32`: 通用库
- `libutil.c32`: 工具库
- `menu.c32`: 菜单模块
- `vesamenu.c32`: 图形菜单模块

### 4.2 创建配置目录

```bash
# 创建pxelinux配置目录
sudo mkdir -p /var/lib/tftpboot/pxelinux.cfg

# 设置权限
sudo chmod 755 /var/lib/tftpboot/pxelinux.cfg
```

## 五、下载Ubuntu安装文件

### 5.1 下载ISO镜像

```bash
# 创建临时目录
mkdir -p ~/pxe-setup
cd ~/pxe-setup

# 下载Ubuntu 22.04桌面版ISO
wget https://releases.ubuntu.com/22.04/ubuntu-22.04.3-desktop-amd64.iso

# 或使用国内镜像加速
wget https://mirrors.aliyun.com/ubuntu-releases/22.04/ubuntu-22.04.3-desktop-amd64.iso
```

### 5.2 挂载ISO并提取文件

```bash
# 创建挂载点
sudo mkdir -p /mnt/ubuntu-iso

# 挂载ISO
sudo mount -o loop ubuntu-22.04.3-desktop-amd64.iso /mnt/ubuntu-iso

# 提取内核和initrd到TFTP目录
sudo cp /mnt/ubuntu-iso/casper/vmlinuz /var/lib/tftpboot/
sudo cp /mnt/ubuntu-iso/casper/initrd /var/lib/tftpboot/

# 验证文件
ls -lh /var/lib/tftpboot/vmlinuz
ls -lh /var/lib/tftpboot/initrd
```

### 5.3 复制ISO到HTTP目录

```bash
# 创建HTTP服务目录
sudo mkdir -p /var/www/html/ubuntu

# 复制整个ISO内容
sudo cp -r /mnt/ubuntu-iso/* /var/www/html/ubuntu/

# 或者直接复制ISO文件
sudo cp ubuntu-22.04.3-desktop-amd64.iso /var/www/html/ubuntu/

# 卸载ISO
sudo umount /mnt/ubuntu-iso
```

## 六、配置PXE启动菜单

### 6.1 创建默认配置文件

**创建 /var/lib/tftpboot/pxelinux.cfg/default：**

```bash
sudo vim /var/lib/tftpboot/pxelinux.cfg/default
```

**基础配置（文本菜单）：**

```conf
DEFAULT ubuntu-desktop
TIMEOUT 100
PROMPT 0

LABEL ubuntu-desktop
    MENU LABEL Install Ubuntu 22.04 Desktop (Auto)
    KERNEL vmlinuz
    APPEND initrd=initrd boot=casper url=http://192.168.1.10/ubuntu/ubuntu-22.04.3-desktop-amd64.iso autoinstall ds=nocloud-net;s=http://192.168.1.10/preseed/ ---

LABEL ubuntu-desktop-manual
    MENU LABEL Install Ubuntu 22.04 Desktop (Manual)
    KERNEL vmlinuz
    APPEND initrd=initrd boot=casper url=http://192.168.1.10/ubuntu/ubuntu-22.04.3-desktop-amd64.iso ---

LABEL local
    MENU LABEL Boot from local disk
    LOCALBOOT 0
```

**参数说明：**
- `DEFAULT`: 默认启动项
- `TIMEOUT`: 超时时间（单位：1/10秒，100=10秒）
- `PROMPT`: 是否显示提示符（0=不显示）
- `KERNEL`: 内核文件路径
- `APPEND`: 内核启动参数
- `LOCALBOOT 0`: 从本地硬盘启动

### 6.2 图形化菜单配置（可选）

**使用vesamenu.c32创建图形菜单：**

```conf
DEFAULT vesamenu.c32
TIMEOUT 300
PROMPT 0

MENU TITLE RemoteGPU PXE Boot Menu
MENU BACKGROUND splash.png
MENU COLOR screen 37;40 #80ffffff #00000000 std

LABEL ubuntu-auto
    MENU LABEL ^1) Ubuntu 22.04 Desktop - Auto Install
    MENU DEFAULT
    KERNEL vmlinuz
    APPEND initrd=initrd boot=casper url=http://192.168.1.10/ubuntu/ubuntu-22.04.3-desktop-amd64.iso autoinstall ds=nocloud-net;s=http://192.168.1.10/preseed/ ---

LABEL ubuntu-manual
    MENU LABEL ^2) Ubuntu 22.04 Desktop - Manual Install
    KERNEL vmlinuz
    APPEND initrd=initrd boot=casper url=http://192.168.1.10/ubuntu/ubuntu-22.04.3-desktop-amd64.iso ---

LABEL local
    MENU LABEL ^3) Boot from Local Disk
    LOCALBOOT 0

LABEL memtest
    MENU LABEL ^4) Memory Test
    KERNEL memtest86+
```

### 6.3 根据MAC地址使用不同配置

**为特定服务器创建专用配置：**

```bash
# MAC地址：00:11:22:33:44:55
# 配置文件名：01-00-11-22-33-44-55

sudo vim /var/lib/tftpboot/pxelinux.cfg/01-00-11-22-33-44-55
```

**专用配置示例：**

```conf
DEFAULT ubuntu-gpu
TIMEOUT 50

LABEL ubuntu-gpu
    MENU LABEL GPU Server Auto Install
    KERNEL vmlinuz
    APPEND initrd=initrd boot=casper url=http://192.168.1.10/ubuntu/ubuntu-22.04.3-desktop-amd64.iso autoinstall ds=nocloud-net;s=http://192.168.1.10/preseed/gpu-server/ ---
```

**配置文件查找顺序：**
1. `01-<MAC地址>` (如：01-00-11-22-33-44-55)
2. `<IP地址的十六进制>` (如：C0A80165 = 192.168.1.101)
3. `default`

## 七、测试TFTP服务

### 7.1 本地测试

```bash
# 安装tftp客户端
sudo apt install -y tftp

# 测试连接
tftp 192.168.1.10

# 在tftp提示符下测试下载文件
tftp> get pxelinux.0
tftp> get vmlinuz
tftp> quit

# 检查下载的文件
ls -lh pxelinux.0 vmlinuz
```

### 7.2 查看TFTP日志

```bash
# 实时查看TFTP日志
sudo tail -f /var/log/syslog | grep tftpd

# 应该看到类似输出：
# tftpd[1234]: RRQ from 192.168.1.101 filename pxelinux.0
# tftpd[1234]: RRQ from 192.168.1.101 filename pxelinux.cfg/default
# tftpd[1234]: RRQ from 192.168.1.101 filename vmlinuz
```

### 7.3 防火墙配置

```bash
# 允许TFTP流量（UDP 69端口）
sudo ufw allow 69/udp

# 查看防火墙状态
sudo ufw status
```

## 八、HTTP服务器配置

### 8.1 安装Apache

```bash
# 安装Apache2
sudo apt install -y apache2

# 启动服务
sudo systemctl start apache2
sudo systemctl enable apache2
```

### 8.2 配置HTTP目录

```bash
# 创建目录结构
sudo mkdir -p /var/www/html/ubuntu
sudo mkdir -p /var/www/html/preseed

# 设置权限
sudo chown -R www-data:www-data /var/www/html
sudo chmod -R 755 /var/www/html
```

### 8.3 测试HTTP服务

```bash
# 测试访问
curl http://192.168.1.10/ubuntu/

# 应该看到目录列表或ISO文件
```

### 8.4 启用目录浏览（可选）

**编辑 /etc/apache2/apache2.conf：**

```apache
<Directory /var/www/html>
    Options Indexes FollowSymLinks
    AllowOverride None
    Require all granted
</Directory>
```

**重启Apache：**

```bash
sudo systemctl restart apache2
```

## 九、完整的目录结构

```
/var/lib/tftpboot/
├── pxelinux.0              # PXE引导程序
├── ldlinux.c32             # 核心模块
├── libcom32.c32            # 通用库
├── libutil.c32             # 工具库
├── menu.c32                # 菜单模块
├── vesamenu.c32            # 图形菜单
├── vmlinuz                 # Linux内核
├── initrd                  # 初始化内存盘
└── pxelinux.cfg/
    ├── default             # 默认配置
    └── 01-00-11-22-33-44-55  # MAC地址专用配置

/var/www/html/
├── ubuntu/
│   ├── ubuntu-22.04.3-desktop-amd64.iso
│   └── (或ISO解压后的文件)
└── preseed/
    ├── user-data           # 自动安装配置
    ├── meta-data           # 元数据
    └── post-install.sh     # 后期配置脚本
```

## 十、故障排查

### 10.1 PXE启动失败

**问题：客户端显示"PXE-E32: TFTP open timeout"**

```bash
# 检查TFTP服务状态
sudo systemctl status tftpd-hpa

# 检查防火墙
sudo ufw status

# 检查文件权限
ls -l /var/lib/tftpboot/pxelinux.0

# 测试TFTP连接
tftp 192.168.1.10
> get pxelinux.0
```

**问题：找不到配置文件**

```bash
# 查看TFTP日志
sudo tail -f /var/log/syslog | grep tftpd

# 检查配置文件是否存在
ls -l /var/lib/tftpboot/pxelinux.cfg/default

# 检查文件权限
sudo chmod 644 /var/lib/tftpboot/pxelinux.cfg/default
```

### 10.2 内核加载失败

**问题：显示"Could not find kernel image"**

```bash
# 检查内核文件
ls -lh /var/lib/tftpboot/vmlinuz
ls -lh /var/lib/tftpboot/initrd

# 重新提取内核
sudo mount -o loop ubuntu-22.04.3-desktop-amd64.iso /mnt/ubuntu-iso
sudo cp /mnt/ubuntu-iso/casper/vmlinuz /var/lib/tftpboot/
sudo cp /mnt/ubuntu-iso/casper/initrd /var/lib/tftpboot/
sudo umount /mnt/ubuntu-iso
```

### 10.3 HTTP下载失败

**问题：无法下载ISO或preseed文件**

```bash
# 检查Apache状态
sudo systemctl status apache2

# 测试HTTP访问
curl -I http://192.168.1.10/ubuntu/ubuntu-22.04.3-desktop-amd64.iso

# 检查文件权限
ls -l /var/www/html/ubuntu/
```

## 十一、性能优化

### 11.1 使用NFS代替HTTP（可选）

**安装NFS服务器：**

```bash
sudo apt install -y nfs-kernel-server

# 配置NFS共享
sudo vim /etc/exports

# 添加：
/var/www/html/ubuntu 192.168.1.0/24(ro,sync,no_subtree_check)

# 重启NFS服务
sudo exportfs -ra
sudo systemctl restart nfs-kernel-server
```

**修改PXE配置使用NFS：**

```conf
APPEND initrd=initrd boot=casper netboot=nfs nfsroot=192.168.1.10:/var/www/html/ubuntu autoinstall ---
```

### 11.2 使用本地镜像源

**在preseed配置中指定本地镜像：**

```yaml
apt:
  primary:
    - arches: [default]
      uri: http://192.168.1.10/ubuntu-mirror/
```

## 十二、总结

**TFTP和PXE配置要点：**
1. ✅ 正确配置TFTP服务器和目录权限
2. ✅ 复制所有必要的PXE启动文件
3. ✅ 提取Ubuntu内核和initrd
4. ✅ 配置PXE启动菜单
5. ✅ 配置HTTP服务器提供安装文件

**下一步：** 配置Preseed自动化安装
