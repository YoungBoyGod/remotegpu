# Harbor 镜像仓库

> ⚠️ **重要提示**: Harbor是企业级应用，配置复杂。**强烈推荐使用官方安装脚本**而不是此docker-compose配置。

## 推荐安装方法（官方脚本）

Harbor官方提供了完整的安装脚本，包含所有必需的配置和初始化步骤。

### 1. 下载Harbor离线安装包

```bash
# 下载最新版本
wget https://github.com/goharbor/harbor/releases/download/v2.9.0/harbor-offline-installer-v2.9.0.tgz

# 解压
tar xzvf harbor-offline-installer-v2.9.0.tgz
cd harbor
```

### 2. 配置Harbor

```bash
# 复制配置模板
cp harbor.yml.tmpl harbor.yml

# 编辑配置文件
vim harbor.yml
```

**关键配置项：**
```yaml
hostname: your-harbor-domain.com

# HTTP配置
http:
  port: 8082

# 数据库配置（使用外部PostgreSQL）
database:
  type: external
  external:
    host: postgresql
    port: 5432
    username: harbor
    password: changeme_db_password
    database: harbor

# Redis配置（使用外部Redis）
redis:
  type: external
  external:
    addr: redis:6379
    password: changeme_redis_password
    db_index: 0

# 管理员密码
harbor_admin_password: Harbor12345
```

### 3. 安装Harbor

```bash
# 执行安装脚本
sudo ./install.sh
```

### 4. 访问Harbor

- 访问地址: http://localhost:8082
- 默认用户名: admin
- 默认密码: Harbor12345（或配置文件中设置的密码）

## Docker客户端配置

```bash
# 配置Docker信任Harbor（HTTP）
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<EOF
{
  "insecure-registries": ["localhost:8082"]
}
EOF

# 重启Docker
sudo systemctl restart docker

# 登录Harbor
docker login localhost:8082
```

## 使用示例

```bash
# 标记镜像
docker tag myimage:latest localhost:8082/library/myimage:latest

# 推送镜像
docker push localhost:8082/library/myimage:latest

# 拉取镜像
docker pull localhost:8082/library/myimage:latest
```

## 为什么推荐官方安装？

1. **完整的初始化**：官方脚本会自动初始化数据库schema
2. **正确的配置**：自动生成所有必需的配置文件
3. **依赖管理**：自动处理组件间的依赖关系
4. **更新支持**：官方提供升级脚本
5. **稳定性**：经过充分测试的部署方案

## 故障排查

### 查看日志
```bash
cd harbor
docker-compose logs -f
```

### 重启服务
```bash
cd harbor
docker-compose down
docker-compose up -d
```

### 卸载Harbor
```bash
cd harbor
docker-compose down -v
```

---

**注意**: 本目录下的docker-compose.yml仅供参考，实际部署请使用官方安装脚本。
