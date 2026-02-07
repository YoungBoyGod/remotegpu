# 企业网络方案 - 完整实施指南

## 方案概述

**无需frp，直接连接，性能最优**

---

## 前置条件检查

- [ ] 有企业固定公网IP
- [ ] 有防火墙/路由器管理权限
- [ ] GPU机器在内网（有固定内网IP）
- [ ] 可以配置端口转发规则

---

## 实施步骤

### 第一步：配置DNS（泛域名）
参考：`gpu-dns-ssl-config.md`

```
类型：A
主机记录：*.gpu
记录值：云服务器IP
```

### 第二步：获取SSL证书
参考：`gpu-dns-ssl-config.md`

```bash
certbot certonly --manual -d "*.gpu.domain.com"
```

### 第三步：配置企业防火墙
参考：`enterprise-firewall-config.md`

- 配置1000条端口转发规则
- 限制只允许云服务器IP访问
- 使用批量脚本：`enterprise-batch-scripts.md`

### 第四步：配置云服务器nginx
参考：`enterprise-nginx-config.md`

- 配置子域名匹配
- 代理到企业公网IP
- 使用批量脚本生成配置

### 第五步：测试验证

```bash
# 测试SSH(直接访问企业公网IP)
ssh -p 2201 user@企业公网IP

# 测试Web服务
curl https://gpu1-jupyter.gpu.domain.com
```

**注意**: SSH和Web服务的访问方式不同,详见 `ssh-vs-web-config.md`

---

## 文档索引

1. **总体架构**：`enterprise-direct-architecture.md`
2. **防火墙配置**：`enterprise-firewall-config.md`
3. **nginx配置**：`enterprise-nginx-config.md`
4. **DNS和SSL**：`gpu-dns-ssl-config.md`
5. **批量脚本**：`enterprise-batch-scripts.md`
6. **SSH vs Web配置差异**：`ssh-vs-web-config.md` ⭐重要

---

## 优势总结

- ✅ 无需frp（配置更简单）
- ✅ 性能最优（无隧道开销）
- ✅ 延迟最低（直接连接）
- ✅ 成本最低（只需云服务器）
