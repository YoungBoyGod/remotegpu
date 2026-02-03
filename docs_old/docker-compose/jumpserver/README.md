# JumpServer 堡垒机

> ⚠️ **注意**: 本项目使用外部JumpServer服务，此配置仅供参考。
>
> 如果你需要在本地部署JumpServer进行测试或开发，可以使用此配置。

## 环境变量配置

创建 `.env` 文件：

```bash
SECRET_KEY=your_secret_key_at_least_50_chars_long_random_string
BOOTSTRAP_TOKEN=your_bootstrap_token
DB_HOST=postgresql
DB_PORT=5432
DB_USER=jumpserver
DB_PASSWORD=your_db_password
DB_NAME=jumpserver
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password
```

## 生成密钥

```bash
# 生成 SECRET_KEY
cat /dev/urandom | tr -dc A-Za-z0-9 | head -c 50

# 生成 BOOTSTRAP_TOKEN
cat /dev/urandom | tr -dc A-Za-z0-9 | head -c 24
```

## 启动服务

```bash
docker-compose up -d
```

## 访问控制台

- Web 控制台: http://localhost:8080
- SSH 端口: 2222
- 默认用户名: admin
- 默认密码: admin

**首次登录后请立即修改密码！**

## 停止服务

```bash
docker-compose down
```
