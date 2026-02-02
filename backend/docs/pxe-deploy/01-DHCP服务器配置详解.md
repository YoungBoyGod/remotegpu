# DHCP服务器配置详解

## 一、DHCP服务器作用

DHCP服务器在PXE启动中的作用：
1. 为裸金属服务器分配IP地址
2. 告知PXE服务器的位置（next-server）
3. 指定PXE启动文件名（filename）

## 二、安装DHCP服务器

```bash
# 安装isc-dhcp-server
sudo apt update
sudo apt install -y isc-dhcp-server

# 检查安装
dpkg -l | grep isc-dhcp-server
```

## 三、网络规划

### 3.1 确定网络参数

**示例网络配置：**
```
网段：192.168.1.0/24
网关：192.168.1.1
DNS：8.8.8.8, 8.8.4.4
PXE服务器IP：192.168.1.10
DHCP地址池：192.168.1.100 - 192.168.1.200
```

### 3.2 查看网卡信息

```bash
# 查看网卡名称
ip addr show

# 示例输出：
# 2: ens33: <BROADCAST,MULTICAST,UP,LOWER_UP>
#     inet 192.168.1.10/24
```

记录网卡名称（如：ens33），后续配置需要使用。

## 四、配置DHCP服务器

### 4.1 编辑主配置文件

**编辑 /etc/dhcp/dhcpd.conf：**

```bash
sudo vim /etc/dhcp/dhcpd.conf
```

**基础配置：**

```conf
# 全局配置
default-lease-time 600;           # 默认租约时间（秒）
max-lease-time 7200;              # 最大租约时间（秒）
authoritative;                     # 声明为权威DHCP服务器

# 日志级别
log-facility local7;

# 子网配置
subnet 192.168.1.0 netmask 255.255.255.0 {
    # 地址池范围
    range 192.168.1.100 192.168.1.200;

    # 网关
    option routers 192.168.1.1;

    # DNS服务器
    option domain-name-servers 8.8.8.8, 8.8.4.4;

    # 子网掩码
    option subnet-mask 255.255.255.0;

    # 域名
    option domain-name "remotegpu.local";

    # PXE启动配置
    next-server 192.168.1.10;      # TFTP服务器IP
    filename "pxelinux.0";          # PXE启动文件
}
```

### 4.2 固定IP分配（可选）

**为特定MAC地址分配固定IP：**

```conf
# GPU服务器1
host gpu-server-01 {
    hardware ethernet 00:11:22:33:44:55;
    fixed-address 192.168.1.101;
    option host-name "gpu-server-01";
}

# GPU服务器2
host gpu-server-02 {
    hardware ethernet 00:11:22:33:44:56;
    fixed-address 192.168.1.102;
    option host-name "gpu-server-02";
}

# GPU服务器3
host gpu-server-03 {
    hardware ethernet 00:11:22:33:44:57;
    fixed-address 192.168.1.103;
    option host-name "gpu-server-03";
}
```

**获取服务器MAC地址：**

```bash
# 方法1：在服务器上查看
ip link show

# 方法2：从DHCP日志查看
sudo tail -f /var/log/syslog | grep DHCPDISCOVER

# 方法3：从路由器管理界面查看
```

### 4.3 配置网络接口

**编辑 /etc/default/isc-dhcp-server：**

```bash
sudo vim /etc/default/isc-dhcp-server
```

**指定监听的网卡：**

```conf
# IPv4配置
INTERFACESv4="ens33"

# IPv6配置（如果不需要可以留空）
INTERFACESv6=""
```

**注意：** 将 `ens33` 替换为你的实际网卡名称。

## 五、启动和测试

### 5.1 检查配置文件语法

```bash
# 检查配置文件是否有语法错误
sudo dhcpd -t -cf /etc/dhcp/dhcpd.conf

# 如果没有错误，会输出：
# Internet Systems Consortium DHCP Server 4.4.1
# Copyright 2004-2018 Internet Systems Consortium.
# All rights reserved.
```

### 5.2 启动DHCP服务

```bash
# 启动服务
sudo systemctl start isc-dhcp-server

# 设置开机自启
sudo systemctl enable isc-dhcp-server

# 查看服务状态
sudo systemctl status isc-dhcp-server
```

**正常输出示例：**

```
● isc-dhcp-server.service - ISC DHCP IPv4 server
     Loaded: loaded (/lib/systemd/system/isc-dhcp-server.service; enabled)
     Active: active (running) since Sun 2026-02-02 10:00:00 UTC; 5s ago
```

### 5.3 测试DHCP服务

**方法1：使用测试客户端**

```bash
# 在另一台机器上测试
sudo dhclient -v ens33

# 查看获取的IP
ip addr show ens33
```

**方法2：查看DHCP日志**

```bash
# 实时查看DHCP日志
sudo tail -f /var/log/syslog | grep dhcpd

# 应该看到类似输出：
# DHCPDISCOVER from 00:11:22:33:44:55 via ens33
# DHCPOFFER on 192.168.1.101 to 00:11:22:33:44:55 via ens33
# DHCPREQUEST for 192.168.1.101 from 00:11:22:33:44:55 via ens33
# DHCPACK on 192.168.1.101 to 00:11:22:33:44:55 via ens33
```

## 六、高级配置

### 6.1 多网段支持

**配置多个子网：**

```conf
# 子网1：管理网络
subnet 192.168.1.0 netmask 255.255.255.0 {
    range 192.168.1.100 192.168.1.200;
    option routers 192.168.1.1;
    next-server 192.168.1.10;
    filename "pxelinux.0";
}

# 子网2：业务网络
subnet 192.168.2.0 netmask 255.255.255.0 {
    range 192.168.2.100 192.168.2.200;
    option routers 192.168.2.1;
    next-server 192.168.1.10;  # 仍然指向同一个PXE服务器
    filename "pxelinux.0";
}
```

### 6.2 根据客户端类型分配不同配置

```conf
# 定义客户端类别
class "pxe-clients" {
    match if substring (option vendor-class-identifier, 0, 9) = "PXEClient";
}

subnet 192.168.1.0 netmask 255.255.255.0 {
    pool {
        allow members of "pxe-clients";
        range 192.168.1.100 192.168.1.150;
        next-server 192.168.1.10;
        filename "pxelinux.0";
    }

    pool {
        deny members of "pxe-clients";
        range 192.168.1.151 192.168.1.200;
    }
}
```

### 6.3 DHCP中继（跨网段PXE）

**如果PXE服务器和裸金属服务器不在同一网段：**

```bash
# 在路由器或网关上配置DHCP中继
# 示例（Cisco路由器）：
interface GigabitEthernet0/0
 ip helper-address 192.168.1.10
```

## 七、故障排查

### 7.1 服务无法启动

**检查配置文件语法：**

```bash
sudo dhcpd -t -cf /etc/dhcp/dhcpd.conf
```

**查看详细错误日志：**

```bash
sudo journalctl -u isc-dhcp-server -n 50
```

**常见错误：**

1. **端口被占用**
   ```bash
   # 检查端口67是否被占用
   sudo netstat -tulnp | grep :67

   # 如果被占用，停止其他DHCP服务
   sudo systemctl stop dnsmasq
   ```

2. **网卡配置错误**
   ```bash
   # 确认网卡名称正确
   ip link show

   # 修改 /etc/default/isc-dhcp-server
   ```

3. **权限问题**
   ```bash
   # 检查配置文件权限
   ls -l /etc/dhcp/dhcpd.conf

   # 应该是：-rw-r--r-- root root
   ```

### 7.2 客户端无法获取IP

**检查网络连通性：**

```bash
# 在PXE服务器上抓包
sudo tcpdump -i ens33 port 67 or port 68

# 应该看到DHCP请求
```

**检查防火墙：**

```bash
# 允许DHCP流量
sudo ufw allow 67/udp
sudo ufw allow 68/udp

# 或者临时关闭防火墙测试
sudo ufw disable
```

### 7.3 PXE启动失败

**检查next-server和filename配置：**

```bash
# 确认TFTP服务器可访问
tftp 192.168.1.10
> get pxelinux.0
> quit

# 如果失败，检查TFTP服务
sudo systemctl status tftpd-hpa
```

## 八、维护和监控

### 8.1 查看租约信息

```bash
# 查看当前租约
cat /var/lib/dhcp/dhcpd.leases

# 示例输出：
# lease 192.168.1.101 {
#   starts 4 2026/02/02 10:00:00;
#   ends 4 2026/02/02 10:10:00;
#   hardware ethernet 00:11:22:33:44:55;
#   client-hostname "gpu-server-01";
# }
```

### 8.2 清理过期租约

```bash
# 停止服务
sudo systemctl stop isc-dhcp-server

# 清理租约文件
sudo rm /var/lib/dhcp/dhcpd.leases
sudo touch /var/lib/dhcp/dhcpd.leases

# 重启服务
sudo systemctl start isc-dhcp-server
```

### 8.3 监控脚本

**创建监控脚本：**

```bash
#!/bin/bash
# /usr/local/bin/monitor-dhcp.sh

# 检查DHCP服务状态
if ! systemctl is-active --quiet isc-dhcp-server; then
    echo "DHCP服务已停止，尝试重启..."
    systemctl restart isc-dhcp-server

    # 发送告警（可选）
    # curl -X POST https://api.example.com/alert \
    #   -d "message=DHCP服务异常"
fi

# 检查租约数量
LEASE_COUNT=$(grep "^lease" /var/lib/dhcp/dhcpd.leases | wc -l)
echo "当前租约数量: $LEASE_COUNT"

# 如果租约过多，发送告警
if [ $LEASE_COUNT -gt 90 ]; then
    echo "警告：租约数量接近上限"
fi
```

**添加到crontab：**

```bash
# 每5分钟检查一次
crontab -e

# 添加：
*/5 * * * * /usr/local/bin/monitor-dhcp.sh >> /var/log/dhcp-monitor.log 2>&1
```

## 九、安全建议

### 9.1 限制DHCP服务范围

```conf
# 只允许已知MAC地址获取IP
subnet 192.168.1.0 netmask 255.255.255.0 {
    deny unknown-clients;  # 拒绝未知客户端

    # 只有在host声明中的MAC地址才能获取IP
}

# 然后为每台服务器添加host声明
host gpu-server-01 {
    hardware ethernet 00:11:22:33:44:55;
    fixed-address 192.168.1.101;
}
```

### 9.2 日志审计

```bash
# 配置详细日志
# 在 /etc/dhcp/dhcpd.conf 中添加：
log-facility local7;

# 配置rsyslog
sudo vim /etc/rsyslog.d/dhcpd.conf

# 添加：
local7.*    /var/log/dhcpd.log

# 重启rsyslog
sudo systemctl restart rsyslog
```

### 9.3 网络隔离

- 将PXE网络与生产网络隔离
- 使用VLAN划分管理网络
- 配置防火墙规则限制访问

## 十、总结

**DHCP配置要点：**
1. ✅ 正确配置网段和地址池
2. ✅ 指定next-server和filename
3. ✅ 为服务器分配固定IP（推荐）
4. ✅ 配置正确的网卡接口
5. ✅ 测试验证配置

**下一步：** 配置TFTP服务器和PXE启动文件
