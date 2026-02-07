# frp方案 - 完整实施指南

## 方案选择

frp方案提供两种nginx配置方式,请根据实际需求选择:

### 方案A: Nginx Proxy Manager (NPM) - 推荐用于开发环境

**优势**:
- ✅ 图形化界面,操作简单
- ✅ 自动申请和续期SSL证书
- ✅ 无需手动编辑nginx配置
- ✅ 适合不熟悉nginx的用户

**适用场景**: 开发环境、测试环境、快速搭建

**文档**: 详见 **[frp-npm-guide.md](frp-npm-guide.md)** - NPM完整实施指南

### 方案B: 手动配置Nginx - 推荐用于生产环境

**优势**:
- ✅ 性能更好,资源占用低
- ✅ 配置灵活,完全自定义
- ✅ 适合生产环境

**适用场景**: 生产环境、对性能要求高的场景

**文档**: 本文档介绍手动nginx配置方案

---

## 方案概述

**适用场景**: 开发环境,无固定公网IP,或无防火墙管理权限

**核心组件**:
- 云服务器(公网IP) - 运行frps服务端 + nginx + SSL
- GPU机器(内网) - 运行frpc客户端
- 域名 - 泛域名解析

---

## 前置条件检查

- [ ] 有一台云服务器(公网IP)
- [ ] 有域名(可以配置DNS)
- [ ] GPU机器可以访问外网
- [ ] 云服务器和GPU机器可以互相通信

---

## 网络架构

```
用户
  ↓
泛域名 *.gpu.domain.com
  ↓
云服务器(nginx + SSL + frps)
  ↓
frp隧道
  ↓
GPU机器(frpc客户端)
  ↓
本地服务(SSH/Jupyter/TensorBoard等)
```

---

## 实施步骤概览

### 第一步:配置DNS(泛域名)
- 配置 `*.gpu.domain.com` 指向云服务器IP
- 详见: `frp-step1-dns.md`

### 第二步:获取SSL证书
- 使用certbot获取泛域名证书
- 详见: `frp-step2-ssl.md`

### 第三步:安装和配置frps(服务端)
- 在云服务器安装frp
- 配置frps.ini
- 详见: `frp-step3-frps.md`

### 第四步:配置nginx
- 配置子域名代理
- 配置SSL
- 详见: `frp-step4-nginx.md`

### 第五步:安装和配置frpc(客户端)
- 在GPU机器安装frp
- 配置frpc.ini
- 详见: `frp-step5-frpc.md`

### 第六步:测试验证
- 测试SSH连接
- 测试Web服务访问
- 详见: `frp-step6-test.md`

---

## 端口分配规则

### frp内部端口(云服务器本地)
- SSH: 10001-10200 (GPU1-200)
- Jupyter: 11001-11200
- TensorBoard: 12001-12200
- 服务1: 13001-13200
- 服务2: 14001-14200

### 访问方式

**SSH:**
```bash
ssh -p 10001 user@云服务器IP    # GPU1
ssh -p 10002 user@云服务器IP    # GPU2
```

**Web服务:**
```bash
https://gpu1-jupyter.gpu.domain.com
https://gpu2-tensorboard.gpu.domain.com
```

---

## 文档索引

### 手动Nginx方案 (本文档)

1. **第一步 - DNS配置**: `frp-step1-dns.md`
2. **第二步 - SSL证书**: `frp-step2-ssl.md`
3. **第三步 - frps服务端**: `frp-step3-frps.md`
4. **第四步 - nginx配置**: `frp-step4-nginx.md`
5. **第五步 - frpc客户端**: `frp-step5-frpc.md`
6. **第六步 - 测试验证**: `frp-step6-test.md`
7. **批量配置脚本**: `frp-batch-scripts.md`

### NPM方案 (推荐用于开发环境)

1. **NPM完整指南**: `frp-npm-guide.md` - 完整实施步骤和架构说明
2. **NPM安装**: `frp-npm-installation.md` - 安装和初始配置
3. **NPM代理配置**: `frp-npm-proxy-config.md` - 配置反向代理和SSL
4. **NPM批量配置**: `frp-npm-batch-config.md` - 批量添加200台机器配置

---

## 优势和注意事项

### 优势
- ✅ 无需固定公网IP
- ✅ 无需配置防火墙
- ✅ 支持子域名访问
- ✅ 统一SSL证书管理

### 注意事项
- ⚠️ 需要保持frpc客户端运行
- ⚠️ 性能略低于直连方案(有隧道开销)
- ⚠️ 需要云服务器有足够带宽
- ⚠️ SSH和Web服务配置方式不同(详见 `ssh-vs-web-config.md`)

---

## 快速开始

1. 按顺序阅读 `frp-step1-dns.md` 到 `frp-step6-test.md`
2. 每完成一步,进行测试验证
3. 使用批量脚本简化配置(200台机器)
4. 遇到问题查看对应步骤的详细文档
