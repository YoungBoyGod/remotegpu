# 方案5：独立公网服务器 - 第二部分：内网服务器配置

## 第二部分：内网服务器配置

### 步骤7：配置内网服务器防火墙

**执行位置**：内网服务器（192.168.10.210）

```bash
# 如果使用ufw
sudo ufw allow from 公网服务器IP to any port 8080
sudo ufw allow from 公网服务器IP to any port 9980

# 如果使用iptables
sudo iptables -A INPUT -s 公网服务器IP -p tcp --dport 8080 -j ACCEPT
sudo iptables -A INPUT -s 公网服务器IP -p tcp --dport 9980 -j ACCEPT
```

**说明**：限制只允许公网服务器访问内网服务，提高安全性。

### 步骤8：测试网络连通性

**执行位置**：公网服务器

```bash
# 测试能否访问内网backend
curl http://192.168.10.210:8080/api/v1/health

# 测试能否访问内网frontend
curl http://192.168.10.210:9980

# 如果无法访问，检查：
# 1. 内网服务器防火墙
# 2. 网络路由
# 3. VPN/专线是否正常
```

### 步骤9：验证访问

**执行位置**：任意设备浏览器

访问 `https://remotegpu.yourdomain.com`

应该看到：
- ✅ HTTPS安全锁
- ✅ 前端页面正常加载
- ✅ API正常工作

---

## 第三部分：网络连通性方案

### 方案A：通过公网IP访问（最简单）

**前提**：内网服务器有公网IP（45.78.48.169）

**配置**：在公网nginx中使用公网IP
```nginx
proxy_pass http://45.78.48.169:8080;
proxy_pass http://45.78.48.169:9980;
```

**优点**：配置简单
**缺点**：内网服务暴露到公网，安全性较低

### 方案B：通过VPN连接（推荐）

**前提**：公网服务器和内网服务器建立VPN隧道

**常用VPN方案**：
- WireGuard（推荐，性能好）
- OpenVPN（成熟稳定）
- IPsec（企业级）

**配置示例（WireGuard）**：

1. 在两台服务器上安装WireGuard
2. 配置VPN隧道，分配虚拟IP（如10.0.0.1和10.0.0.2）
3. 在公网nginx中使用VPN IP：
```nginx
proxy_pass http://10.0.0.2:8080;
proxy_pass http://10.0.0.2:9980;
```

**优点**：安全，内网服务不暴露
**缺点**：配置稍复杂

### 方案C：通过专线/内网互通

**前提**：两台服务器在同一内网或有专线连接

**配置**：直接使用内网IP
```nginx
proxy_pass http://192.168.10.210:8080;
proxy_pass http://192.168.10.210:9980;
```

**优点**：性能最好，最安全
**缺点**：需要网络基础设施支持

---

## 维护命令

**公网服务器**：
```bash
# 重新加载nginx
sudo systemctl reload nginx

# 查看日志
sudo tail -f /var/log/nginx/remotegpu_access.log

# 续期SSL证书
sudo certbot renew
```

**内网服务器**：
```bash
# 检查服务状态
ps aux | grep "go run"
ps aux | grep "vite"

# 查看端口监听
netstat -tlnp | grep -E "8080|9980"
```

---

## 成本估算

- 公网服务器：约￥50-100/月
- 域名：约￥50-100/年
- VPN（可选）：免费（WireGuard）
- 总计：约￥600-1200/年
