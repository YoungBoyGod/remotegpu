# 前端 CDN 部署指南

## 适用场景

所有网络场景都可以使用前端 CDN 加速。

---

## 方案选择

### 推荐：阿里云 OSS + CDN

**成本**：约￥10-30/月
**优点**：国内访问快，配置简单

### 备选：腾讯云 COS + CDN

**成本**：约￥10-30/月
**优点**：与阿里云类似

---

## 实施步骤（以阿里云为例）

### 步骤1：构建前端

**执行位置**：本地开发机或本地服务器

```bash
cd /home/luo/code/remotegpu/frontend

# 安装依赖（如果还没安装）
npm install

# 构建生产版本
npm run build

# 生成 dist/ 目录
```

### 步骤2：创建 OSS 存储桶

**执行位置**：阿里云控制台

1. 登录阿里云控制台
2. 进入 OSS 服务
3. 创建 Bucket：
   - 名称：remotegpu-frontend
   - 区域：选择离用户最近的区域
   - 读写权限：公共读
   - 其他：默认

### 步骤3：上传前端文件

**方式A：使用阿里云控制台**
- 进入 Bucket
- 上传 dist/ 目录下的所有文件

**方式B：使用 ossutil 工具**

```bash
# 安装 ossutil
wget http://gosspublic.alicdn.com/ossutil/1.7.15/ossutil64
chmod 755 ossutil64

# 配置
./ossutil64 config

# 上传
./ossutil64 cp -r dist/ oss://remotegpu-frontend/ --update
```

### 步骤4：配置 CDN

**执行位置**：阿里云控制台

1. 进入 CDN 服务
2. 添加域名：
   - 加速域名：frontend.yourdomain.com
   - 业务类型：图片小文件
   - 源站类型：OSS 域名
   - 源站域名：选择刚创建的 Bucket

3. 配置 HTTPS：
   - 上传 SSL 证书
   - 或使用免费证书

4. 配置缓存规则：
   - `/index.html`：不缓存
   - `*.js, *.css`：缓存 30 天
   - `*.jpg, *.png`：缓存 30 天

### 步骤5：配置 DNS

**执行位置**：域名服务商控制台

添加 CNAME 记录：
- 类型：CNAME
- 主机记录：frontend
- 记录值：CDN 提供的 CNAME 地址
- TTL：600

---

## 前端配置调整

### 修改 API 地址

**执行位置**：本地开发机

编辑前端配置文件（通常是 `.env.production`）：

```bash
# API 地址指向云服务器
VITE_API_BASE_URL=https://remotegpu.yourdomain.com
```

重新构建并上传。

---

## 验证

访问：`https://frontend.yourdomain.com`

应该看到前端页面正常加载。
