# iptables机器隔离配置指南

## 概述

本文档提供可以直接使用的iptables命令,用于实现局域网内机器隔离,同时允许访问公共服务。

## 场景说明

**网络环境**:
- 本地网段: 192.168.1.0/24
- 公共服务: 192.168.1.100-110
- 管理网络: 192.168.200.0/24
- 客户机器: 192.168.1.10-99

**隔离需求**:
- ✅ 客户机器可以访问公共服务
- ✅ 管理网络可以访问所有机器
- ✅ 客户机器可以访问外网
- ❌ 客户机器之间不能互相访问

## 方案一: 基础隔离配置

### 1.1 完整配置脚本

```bash
#!/bin/bash
# isolation_basic.sh - 基础隔离配置脚本

# 定义变量
PUBLIC_SERVICES=(
    "192.168.1.100"  # 存储服务器
    "192.168.1.101"  # 监控服务器
    "192.168.1.102"  # 日志服务器
    "192.168.1.103"  # DNS服务器
    "192.168.1.104"  # NTP服务器
)
MANAGEMENT_NET="192.168.200.0/24"
LOCAL_NET="192.168.1.0/24"

echo "开始配置iptables隔离规则..."

# 1. 清空现有规则
echo "清空现有规则..."
iptables -F INPUT
iptables -F OUTPUT
iptables -F FORWARD

# 2. 设置默认策略(重要:先设置为ACCEPT,避免锁死)
echo "设置默认策略..."
iptables -P INPUT ACCEPT
iptables -P OUTPUT ACCEPT
iptables -P FORWARD DROP

# 3. 允许本机回环接口
echo "允许本机回环..."
iptables -A INPUT -i lo -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

# 4. 允许已建立的连接(非常重要!)
echo "允许已建立的连接..."
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
iptables -A OUTPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# 5. 允许访问公共服务
echo "配置公共服务访问规则..."
for service_ip in "${PUBLIC_SERVICES[@]}"; do
    echo "  允许访问 $service_ip"
    iptables -A OUTPUT -d $service_ip -j ACCEPT
    iptables -A INPUT -s $service_ip -j ACCEPT
done

# 6. 允许管理网络访问
echo "允许管理网络访问..."
iptables -A INPUT -s $MANAGEMENT_NET -j ACCEPT
iptables -A OUTPUT -d $MANAGEMENT_NET -j ACCEPT

# 7. 拒绝来自本地网段的其他访问(核心隔离规则)
echo "配置隔离规则..."
iptables -A INPUT -s $LOCAL_NET -j DROP

# 8. 允许访问外网(默认允许所有出站)
echo "允许访问外网..."
iptables -A OUTPUT -j ACCEPT

# 9. 保存规则
echo "保存iptables规则..."
if command -v iptables-save &> /dev/null; then
    iptables-save > /etc/iptables/rules.v4
    echo "规则已保存到 /etc/iptables/rules.v4"
fi

echo "iptables隔离配置完成!"
echo ""
echo "当前规则:"
iptables -L -n -v
```

### 1.2 使用方法

```bash
# 1. 保存脚本
cat > /usr/local/bin/isolation_basic.sh << 'EOF'
[上面的脚本内容]
EOF

# 2. 添加执行权限
chmod +x /usr/local/bin/isolation_basic.sh

# 3. 执行配置
/usr/local/bin/isolation_basic.sh

# 4. 验证规则
iptables -L -n -v
```

## 方案二: 逐条命令详解

### 2.1 清空规则

```bash
# 清空INPUT链的所有规则
iptables -F INPUT

# 清空OUTPUT链的所有规则
iptables -F OUTPUT

# 清空FORWARD链的所有规则
iptables -F FORWARD
```

**说明**: `-F` (flush) 清空指定链的所有规则,从干净状态开始配置。

### 2.2 设置默认策略

```bash
# 设置INPUT链默认策略为ACCEPT
iptables -P INPUT ACCEPT

# 设置OUTPUT链默认策略为ACCEPT
iptables -P OUTPUT ACCEPT

# 设置FORWARD链默认策略为DROP
iptables -P FORWARD DROP
```

**重要**:
- 先设置为ACCEPT,避免配置过程中锁死自己
- 配置完成后可以改为更严格的策略
- FORWARD链设置为DROP,因为不需要转发

### 2.3 允许回环接口

```bash
# 允许从回环接口进入的流量
iptables -A INPUT -i lo -j ACCEPT

# 允许从回环接口出去的流量
iptables -A OUTPUT -o lo -j ACCEPT
```

**说明**:
- `-i lo`: 指定输入接口为lo(回环接口)
- `-o lo`: 指定输出接口为lo
- `-j ACCEPT`: 动作为接受
- 回环接口用于本机进程间通信,必须允许

### 2.4 允许已建立的连接

```bash
# 允许已建立和相关的连接进入
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# 允许已建立和相关的连接出去
iptables -A OUTPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
```

**说明**:
- `-m state`: 使用state模块
- `--state ESTABLISHED,RELATED`: 匹配已建立和相关的连接
- **非常重要**: 这条规则允许响应流量返回,否则所有连接都会失败

### 2.5 允许访问公共服务

```bash
# 允许访问公共服务(出站)
iptables -A OUTPUT -d 192.168.1.100 -j ACCEPT

# 允许公共服务的响应(入站)
iptables -A INPUT -s 192.168.1.100 -j ACCEPT
```

**说明**:
- `-d 192.168.1.100`: 目标地址为公共服务IP
- `-s 192.168.1.100`: 源地址为公共服务IP
- 需要同时配置出站和入站规则

### 2.6 允许管理网络访问

```bash
# 允许管理网络的流量进入
iptables -A INPUT -s 192.168.200.0/24 -j ACCEPT

# 允许访问管理网络(出站)
iptables -A OUTPUT -d 192.168.200.0/24 -j ACCEPT
```

**说明**:
- `-s 192.168.200.0/24`: 源地址为管理网络
- `-d 192.168.200.0/24`: 目标地址为管理网络
- 管理网络可以访问所有机器

### 2.7 拒绝本地网段(核心隔离规则)

```bash
# 拒绝来自本地网段的流量
iptables -A INPUT -s 192.168.1.0/24 -j DROP
```

**说明**:
- `-s 192.168.1.0/24`: 源地址为本地网段
- `-j DROP`: 动作为丢弃(静默拒绝)
- **核心规则**: 这条规则实现了机器间的隔离
- 放在最后,因为前面的规则已经允许了公共服务

### 2.8 允许访问外网

```bash
# 允许所有出站流量
iptables -A OUTPUT -j ACCEPT
```

**说明**:
- 允许访问外网和其他网段
- 放在最后作为默认规则

## 方案三: 测试和验证

### 3.1 查看当前规则

```bash
# 查看所有规则(详细信息)
iptables -L -n -v

# 查看INPUT链规则
iptables -L INPUT -n -v --line-numbers

# 查看OUTPUT链规则
iptables -L OUTPUT -n -v --line-numbers
```

**说明**:
- `-L`: 列出规则
- `-n`: 以数字形式显示IP和端口
- `-v`: 显示详细信息(包括数据包计数)
- `--line-numbers`: 显示规则行号

### 3.2 测试隔离效果

#### 3.2.1 测试是否能访问公共服务(应该成功)

```bash
# 测试ping公共服务
ping -c 3 192.168.1.100

# 测试SSH连接公共服务
ssh root@192.168.1.100

# 测试HTTP访问公共服务
curl http://192.168.1.100
```

**预期结果**: 所有测试都应该成功

#### 3.2.2 测试是否能访问其他客户机器(应该失败)

```bash
# 测试ping其他客户机器
ping -c 3 192.168.1.20

# 测试SSH连接其他客户机器
ssh root@192.168.1.20

# 测试端口扫描
nc -zv 192.168.1.20 22
```

**预期结果**:
- ping应该超时或无响应
- SSH连接应该失败
- 端口扫描应该显示连接被拒绝

#### 3.2.3 测试是否能访问外网(应该成功)

```bash
# 测试ping外网
ping -c 3 8.8.8.8

# 测试DNS解析
nslookup www.baidu.com

# 测试HTTP访问
curl https://www.baidu.com
```

**预期结果**: 所有测试都应该成功

#### 3.2.4 测试管理网络访问(应该成功)

```bash
# 从管理网络(192.168.200.x)SSH到本机
# 在管理网络的机器上执行:
ssh root@192.168.1.10
```

**预期结果**: SSH连接应该成功

### 3.3 检查规则计数器

```bash
# 查看规则匹配次数
iptables -L -n -v

# 重置计数器
iptables -Z
```

**说明**:
- 每条规则前面的数字显示匹配的数据包数量
- 可以通过计数器判断规则是否生效
- `-Z`参数可以重置计数器

## 方案四: 故障排查

### 4.1 常见问题

#### 4.1.1 无法访问公共服务

**症状**: ping或SSH连接公共服务失败

**排查步骤**:
```bash
# 1. 检查是否有允许公共服务的规则
iptables -L OUTPUT -n -v | grep 192.168.1.100

# 2. 检查规则顺序(DROP规则是否在ACCEPT规则之前)
iptables -L INPUT -n -v --line-numbers

# 3. 临时删除DROP规则测试
iptables -D INPUT -s 192.168.1.0/24 -j DROP

# 4. 测试网络连通性
ping -c 3 192.168.1.100
```

**解决方法**:
- 确保公共服务的ACCEPT规则在DROP规则之前
- 检查公共服务IP地址是否正确

#### 4.1.2 无法访问外网

**症状**: 无法ping外网或访问互联网

**排查步骤**:
```bash
# 1. 检查OUTPUT链是否有阻止规则
iptables -L OUTPUT -n -v

# 2. 检查DNS是否正常
nslookup www.baidu.com

# 3. 检查路由表
ip route show

# 4. 检查网关是否可达
ping -c 3 192.168.1.1
```

**解决方法**:
- 确保OUTPUT链最后有`-j ACCEPT`规则
- 检查DNS配置(`/etc/resolv.conf`)
- 检查默认网关配置

#### 4.1.3 规则不生效

**症状**: 配置了规则但仍然可以访问其他客户机器

**排查步骤**:
```bash
# 1. 确认规则已加载
iptables -L -n -v

# 2. 检查规则计数器是否增加
iptables -L INPUT -n -v | grep DROP

# 3. 检查是否有其他规则冲突
iptables -S

# 4. 测试从其他机器ping本机
# 在其他机器上执行: ping 本机IP
```

**解决方法**:
- 确认规则已正确添加
- 检查规则顺序,ACCEPT规则不应该在DROP规则之后
- 重启iptables服务

#### 4.1.4 SSH连接断开

**症状**: 配置规则后SSH连接立即断开

**原因**: 没有允许已建立的连接

**解决方法**:
```bash
# 必须添加这条规则(非常重要!)
iptables -I INPUT 1 -m state --state ESTABLISHED,RELATED -j ACCEPT
```

### 4.2 调试技巧

#### 4.2.1 启用日志记录

```bash
# 记录被DROP的数据包
iptables -I INPUT -s 192.168.1.0/24 -j LOG --log-prefix "DROPPED: " --log-level 4

# 查看日志
tail -f /var/log/syslog | grep DROPPED
# 或
tail -f /var/log/kern.log | grep DROPPED
```

#### 4.2.2 临时禁用规则

```bash
# 临时清空所有规则(用于测试)
iptables -F INPUT
iptables -F OUTPUT

# 恢复规则
/usr/local/bin/isolation_basic.sh
```

#### 4.2.3 逐条测试规则

```bash
# 先只添加基础规则
iptables -A INPUT -i lo -j ACCEPT
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# 测试是否正常
ping -c 3 8.8.8.8

# 再逐条添加其他规则
iptables -A INPUT -s 192.168.1.100 -j ACCEPT
# 测试...
```

## 方案五: 规则持久化

### 5.1 Ubuntu/Debian系统

#### 5.1.1 使用iptables-persistent

```bash
# 1. 安装iptables-persistent
apt-get update
apt-get install -y iptables-persistent

# 2. 保存当前规则
iptables-save > /etc/iptables/rules.v4

# 3. 重启后自动加载
systemctl enable netfilter-persistent
systemctl start netfilter-persistent

# 4. 手动重新加载规则
iptables-restore < /etc/iptables/rules.v4
```

#### 5.1.2 使用systemd服务

创建systemd服务文件:

```bash
cat > /etc/systemd/system/iptables-restore.service <<'EOF'
[Unit]
Description=Restore iptables rules
Before=network-pre.target
Wants=network-pre.target

[Service]
Type=oneshot
ExecStart=/sbin/iptables-restore /etc/iptables/rules.v4
ExecReload=/sbin/iptables-restore /etc/iptables/rules.v4
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
EOF

# 启用服务
systemctl daemon-reload
systemctl enable iptables-restore.service
systemctl start iptables-restore.service
```

### 5.2 CentOS/RHEL系统

```bash
# 1. 安装iptables-services
yum install -y iptables-services

# 2. 保存规则
service iptables save

# 3. 启用开机自启
systemctl enable iptables
systemctl start iptables

# 4. 查看保存的规则
cat /etc/sysconfig/iptables
```

### 5.3 验证持久化

```bash
# 1. 保存当前规则
iptables-save > /etc/iptables/rules.v4

# 2. 重启系统
reboot

# 3. 重启后检查规则是否存在
iptables -L -n -v

# 4. 如果规则不存在,手动加载
iptables-restore < /etc/iptables/rules.v4
```

## 方案六: 与RemoteGPU平台集成

### 6.1 自动化部署方案

#### 6.1.1 使用Ansible批量部署

创建Ansible playbook:

```yaml
# isolation_deploy.yml
---
- name: Deploy iptables isolation rules
  hosts: all_machines
  become: yes
  vars:
    public_services:
      - "192.168.1.100"  # 存储服务器
      - "192.168.1.101"  # 监控服务器
      - "192.168.1.102"  # 日志服务器
    management_net: "192.168.200.0/24"
    local_net: "192.168.1.0/24"

  tasks:
    - name: Install iptables-persistent
      apt:
        name: iptables-persistent
        state: present
        update_cache: yes

    - name: Copy isolation script
      copy:
        src: isolation_basic.sh
        dest: /usr/local/bin/isolation_basic.sh
        mode: '0755'

    - name: Execute isolation script
      shell: /usr/local/bin/isolation_basic.sh

    - name: Verify rules are applied
      shell: iptables -L -n -v
      register: iptables_output

    - name: Display iptables rules
      debug:
        var: iptables_output.stdout_lines
```

执行部署:

```bash
# 部署到所有机器
ansible-playbook -i inventory.ini isolation_deploy.yml

# 部署到特定组
ansible-playbook -i inventory.ini isolation_deploy.yml --limit customer_machines
```

#### 6.1.2 在机器分配时自动配置

在后端服务中集成防火墙配置:

```go
// backend/internal/service/firewall.go
package service

import (
    "fmt"
    "os/exec"
    "strings"
)

type FirewallService struct {
    publicServices []string
    managementNet  string
    localNet       string
}

// ConfigureIsolation 配置机器隔离规则
func (s *FirewallService) ConfigureIsolation(machineIP string) error {
    script := s.generateIptablesScript()

    // 通过SSH执行脚本
    cmd := exec.Command("ssh", fmt.Sprintf("root@%s", machineIP), script)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to configure isolation: %v, output: %s", err, output)
    }

    return nil
}

func (s *FirewallService) generateIptablesScript() string {
    var sb strings.Builder

    sb.WriteString("#!/bin/bash\n")
    sb.WriteString("iptables -F INPUT\n")
    sb.WriteString("iptables -F OUTPUT\n")
    sb.WriteString("iptables -A INPUT -i lo -j ACCEPT\n")
    sb.WriteString("iptables -A OUTPUT -o lo -j ACCEPT\n")
    sb.WriteString("iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT\n")
    sb.WriteString("iptables -A OUTPUT -m state --state ESTABLISHED,RELATED -j ACCEPT\n")

    // 添加公共服务规则
    for _, serviceIP := range s.publicServices {
        sb.WriteString(fmt.Sprintf("iptables -A OUTPUT -d %s -j ACCEPT\n", serviceIP))
        sb.WriteString(fmt.Sprintf("iptables -A INPUT -s %s -j ACCEPT\n", serviceIP))
    }

    // 添加管理网络规则
    sb.WriteString(fmt.Sprintf("iptables -A INPUT -s %s -j ACCEPT\n", s.managementNet))
    sb.WriteString(fmt.Sprintf("iptables -A OUTPUT -d %s -j ACCEPT\n", s.managementNet))

    // 拒绝本地网段
    sb.WriteString(fmt.Sprintf("iptables -A INPUT -s %s -j DROP\n", s.localNet))

    // 允许访问外网
    sb.WriteString("iptables -A OUTPUT -j ACCEPT\n")

    // 保存规则
    sb.WriteString("iptables-save > /etc/iptables/rules.v4\n")

    return sb.String()
}

// VerifyIsolation 验证隔离规则是否生效
func (s *FirewallService) VerifyIsolation(machineIP string) (bool, error) {
    cmd := exec.Command("ssh", fmt.Sprintf("root@%s", machineIP),
        "iptables -L INPUT -n | grep -q 'DROP.*192.168.1.0/24'")

    err := cmd.Run()
    if err != nil {
        return false, nil
    }

    return true, nil
}
```

### 6.2 监控和告警

#### 6.2.1 定期检查脚本

```bash
#!/bin/bash
# /usr/local/bin/check_isolation.sh

# 检查iptables规则是否存在
if ! iptables -L INPUT -n | grep -q "DROP.*192.168.1.0/24"; then
    echo "WARNING: Isolation rules not found on $(hostname)!"

    # 发送告警到监控系统
    curl -X POST "http://monitoring-server/alert" \
         -H "Content-Type: application/json" \
         -d "{\"host\": \"$(hostname)\", \"message\": \"Isolation rules missing\"}"

    # 尝试自动恢复
    /usr/local/bin/isolation_basic.sh

    exit 1
fi

echo "Isolation rules OK on $(hostname)"
exit 0
```

添加到crontab:

```bash
# 每5分钟检查一次
*/5 * * * * /usr/local/bin/check_isolation.sh >> /var/log/isolation_check.log 2>&1
```

#### 6.2.2 Prometheus监控指标

```bash
# /usr/local/bin/iptables_exporter.sh
#!/bin/bash

# 导出iptables规则数量
echo "# HELP iptables_rules_count Number of iptables rules"
echo "# TYPE iptables_rules_count gauge"
echo "iptables_rules_count{chain=\"INPUT\"} $(iptables -L INPUT -n | grep -c '^')"
echo "iptables_rules_count{chain=\"OUTPUT\"} $(iptables -L OUTPUT -n | grep -c '^')"

# 导出DROP规则匹配次数
echo "# HELP iptables_drop_packets Packets dropped by iptables"
echo "# TYPE iptables_drop_packets counter"
DROP_COUNT=$(iptables -L INPUT -n -v | grep DROP | awk '{sum+=$1} END {print sum}')
echo "iptables_drop_packets ${DROP_COUNT:-0}"
```

### 6.3 配置管理

#### 6.3.1 集中配置文件

```yaml
# /etc/remotegpu/firewall_config.yaml
public_services:
  - ip: "192.168.1.100"
    name: "storage_server"
    description: "NFS/JuiceFS存储服务器"
  - ip: "192.168.1.101"
    name: "monitoring_server"
    description: "Prometheus监控服务器"
  - ip: "192.168.1.102"
    name: "log_server"
    description: "ELK日志服务器"

management_network:
  cidr: "192.168.200.0/24"
  description: "管理网络,可以访问所有机器"

local_network:
  cidr: "192.168.1.0/24"
  description: "客户机器所在网段"

isolation_enabled: true
auto_recovery: true
check_interval: 300  # 秒
```

#### 6.3.2 动态更新公共服务列表

```bash
#!/bin/bash
# update_public_services.sh

# 从配置文件读取公共服务列表
PUBLIC_SERVICES=$(yq eval '.public_services[].ip' /etc/remotegpu/firewall_config.yaml)

# 清空现有的公共服务规则
iptables -D INPUT -s 192.168.1.100 -j ACCEPT 2>/dev/null
iptables -D OUTPUT -d 192.168.1.100 -j ACCEPT 2>/dev/null

# 重新添加规则
for service_ip in $PUBLIC_SERVICES; do
    iptables -I INPUT -s $service_ip -j ACCEPT
    iptables -I OUTPUT -d $service_ip -j ACCEPT
done

# 保存规则
iptables-save > /etc/iptables/rules.v4
```

## 总结

### 核心要点

本文档提供了完整的iptables机器隔离配置方案,包括:

✅ **方案一**: 完整的bash脚本,可直接使用
✅ **方案二**: 逐条命令详解,便于理解和自定义
✅ **方案三**: 测试验证方法,确保配置正确
✅ **方案四**: 故障排查指南,快速解决问题
✅ **方案五**: 规则持久化配置,重启后自动生效
✅ **方案六**: 与RemoteGPU平台集成方案

### 隔离效果

配置完成后,实现以下隔离效果:

- ✅ 客户机器可以访问公共服务(存储、监控等)
- ✅ 管理网络可以访问所有机器
- ✅ 客户机器可以访问外网
- ❌ 客户机器之间不能互相访问

### 适用场景

本方案特别适合以下场景:

1. **大二层网络**: 所有机器在同一网段,无法使用VLAN隔离
2. **普通交换机**: 交换机不支持Private VLAN等高级功能
3. **成本敏感**: 无预算购买高级网络设备
4. **Linux环境**: 所有机器运行Linux系统
5. **RemoteGPU平台**: 多租户GPU云平台的机器隔离需求

### 实施建议

**快速开始**:
1. 复制方案一的完整脚本
2. 修改公共服务IP列表
3. 执行脚本并测试
4. 配置规则持久化

**生产环境**:
1. 使用Ansible批量部署
2. 与平台集成,自动配置新机器
3. 配置监控告警
4. 定期检查规则是否生效

### 注意事项

⚠️ **重要提醒**:

1. **已建立连接规则**: 必须添加`ESTABLISHED,RELATED`规则,否则SSH会断开
2. **规则顺序**: ACCEPT规则必须在DROP规则之前
3. **测试验证**: 配置后务必测试各种访问场景
4. **规则持久化**: 确保重启后规则自动生效
5. **备份配置**: 保存原有规则,便于回滚

### 相关文档

- [局域网机器隔离方案](./局域网机器隔离方案.md) - 多种隔离方案对比
- [物理机和容器访问架构说明](./物理机和容器访问架构说明.md) - 访问管理架构
- [访问管理架构实施清单](./访问管理架构实施清单.md) - 完整实施指南

---

**文档版本**: v1.0
**创建日期**: 2026-02-02
**最后更新**: 2026-02-02
**维护者**: RemoteGPU开发团队
**状态**: 可直接使用

