# 企业防火墙端口转发配置

## 配置原理

将企业公网IP的端口，转发到内网GPU机器的对应端口。

---

## 端口转发规则

### GPU1 示例（内网IP: 192.168.10.101）

| 服务 | 公网端口 | 内网IP | 内网端口 | 协议 |
|------|---------|--------|---------|------|
| SSH | 2201 | 192.168.10.101 | 22 | TCP |
| Jupyter | 8001 | 192.168.10.101 | 8888 | TCP |
| TensorBoard | 9001 | 192.168.10.101 | 6006 | TCP |
| 服务1 | 10001 | 192.168.10.101 | 本地端口1 | TCP |
| 服务2 | 11001 | 192.168.10.101 | 本地端口2 | TCP |

### GPU2-200 同理

---

## 配置方法（根据设备类型）

### 方法1：Linux iptables

**执行位置**：企业网关服务器

```bash
# GPU1 SSH
iptables -t nat -A PREROUTING -p tcp --dport 2201 -j DNAT --to-destination 192.168.10.101:22
iptables -t nat -A POSTROUTING -d 192.168.10.101 -p tcp --dport 22 -j MASQUERADE

# GPU1 Jupyter
iptables -t nat -A PREROUTING -p tcp --dport 8001 -j DNAT --to-destination 192.168.10.101:8888
iptables -t nat -A POSTROUTING -d 192.168.10.101 -p tcp --dport 8888 -j MASQUERADE

# 保存规则
iptables-save > /etc/iptables/rules.v4
```

### 方法2：企业路由器Web界面

**执行位置**：路由器管理界面

配置示例：
- 外部端口：2201
- 内部IP：192.168.10.101
- 内部端口：22
- 协议：TCP

### 方法3：批量配置脚本

见后续文档。

---

## 安全建议

1. **限制来源IP**：
   - 只允许云服务器IP访问
   - 防止其他IP直接访问

2. **防火墙规则**：
```bash
# 只允许云服务器IP
iptables -A INPUT -s 云服务器IP -p tcp --dport 2201:11200 -j ACCEPT
iptables -A INPUT -p tcp --dport 2201:11200 -j DROP
```
