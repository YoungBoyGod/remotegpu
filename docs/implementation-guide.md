# RemoteGPU 云服务器反向代理 - 完整实施指南

## 文档索引

### 1. 场景选择
- **场景对比**：`network-scenarios-comparison.md`
- 先阅读此文档，确定您的网络场景

### 2. 场景详细文档
- **场景1**：`scenario-1-direct-connection.md`（有固定公网IP）
- **场景2A**：`scenario-2a-firewall-rules.md`（有防火墙权限）
- **场景2B/3**：`scenario-2b3-frp-part1-server.md`、`part2-client.md`、`part3-nginx.md`（frp内网穿透）

### 3. 前端部署
- **前端CDN**：`frontend-cdn-deployment.md`

---

## 快速开始（推荐流程）

### 第一步：确定网络场景

**测试方法**：
```bash
# 在本地服务器执行
curl ifconfig.me

# 记录返回的IP
```

**判断**：
- 如果是固定IP且可以开放端口 → 场景1
- 如果有防火墙且有管理权限 → 场景2A
- 如果是NAT或无防火墙权限 → 场景2B/3（推荐）

---

### 第二步：准备工作

**需要准备**：
1. 云服务器（2核4G，约￥100-150/月）
2. 域名（约￥50-100/年）
3. 确认本地Backend运行正常

**云服务器推荐**：
- 阿里云：华东、华北区域
- 腾讯云：上海、北京区域

---

### 第三步：实施部署

**根据场景选择对应文档**：

#### 场景2B/3（最常见）：
1. 阅读 `scenario-2b3-frp-part1-server.md`（云服务器配置）
2. 阅读 `scenario-2b3-frp-part2-client.md`（本地服务器配置）
3. 阅读 `scenario-2b3-frp-part3-nginx.md`（nginx配置）
4. 阅读 `frontend-cdn-deployment.md`（前端部署）

#### 场景1：
1. 阅读 `scenario-1-direct-connection.md`
2. 阅读 `frontend-cdn-deployment.md`

#### 场景2A：
1. 阅读 `scenario-2a-firewall-rules.md`
2. 阅读 `frontend-cdn-deployment.md`

---

### 第四步：验证测试

**测试清单**：
- [ ] frp隧道连接正常（场景2B/3）
- [ ] API访问正常：`https://remotegpu.yourdomain.com/api/v1/health`
- [ ] 前端访问正常：`https://frontend.yourdomain.com`
- [ ] 前端可以调用API
- [ ] GPU机器通信正常

---

## 常见问题

### Q1：如何判断我的网络场景？
A：参考 `network-scenarios-comparison.md` 中的决策树。

### Q2：frp和直接连接哪个更好？
A：如果有固定公网IP，直接连接性能更好。如果是NAT或有防火墙限制，必须使用frp。

### Q3：前端必须部署到CDN吗？
A：不是必须，但强烈推荐。CDN可以显著提升全国用户的访问速度。

### Q4：成本大概多少？
A：
- 云服务器：￥100-150/月
- CDN：￥10-30/月
- 域名：￥50-100/年
- 总计：约￥110-180/月

---

## 下一步

1. 确定您的网络场景
2. 阅读对应的详细文档
3. 准备云服务器和域名
4. 按照文档步骤实施
5. 如需协助，可以让团队帮忙实施
