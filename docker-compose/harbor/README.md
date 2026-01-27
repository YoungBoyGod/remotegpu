# Harbor 镜像仓库

## 安装说明

Harbor 推荐使用官方安装脚本部署，这里提供简化的 Docker Compose 配置。

## 环境变量配置

创建 `.env` 文件：

```bash
HARBOR_ADMIN_PASSWORD=changeme_harbor_password
HARBOR_DATABASE_PASSWORD=changeme_db_password
```

## 启动服务

```bash
docker-compose up -d
```

## 访问控制台

- 访问地址: http://localhost:8082
- 默认用户名: admin
- 默认密码: Harbor12345 (或 .env 中配置的密码)

**首次登录后请立即修改密码！**

## Docker 客户端配置

```bash
# 配置 Docker 信任 Harbor（HTTP）
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<EOF
{
  "insecure-registries": ["localhost:8082"]
}
EOF

# 重启 Docker
sudo systemctl restart docker

# 登录 Harbor
docker login localhost:8082
```

## 推送镜像示例

```bash
# 标记镜像
docker tag myimage:latest localhost:8082/library/myimage:latest

# 推送镜像
docker push localhost:8082/library/myimage:latest
```

## 生产环境建议

生产环境建议使用官方安装包：

```bash
wget https://github.com/goharbor/harbor/releases/download/v2.9.0/harbor-offline-installer-v2.9.0.tgz
tar xzvf harbor-offline-installer-v2.9.0.tgz
cd harbor
cp harbor.yml.tmpl harbor.yml
# 编辑 harbor.yml 配置文件
./install.sh
```

## 停止服务

```bash
docker-compose down
```
