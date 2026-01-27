# Apache Guacamole 远程桌面网关

## 数据库初始化

首次启动前需要初始化数据库：

```bash
# 1. 生成初始化脚本
docker run --rm guacamole/guacamole:1.5.4 /opt/guacamole/bin/initdb.sh --postgresql > initdb.sql

# 2. 在 PostgreSQL 中创建数据库和用户
docker exec -it remotegpu-postgresql psql -U postgres -c "CREATE DATABASE guacamole;"
docker exec -it remotegpu-postgresql psql -U postgres -c "CREATE USER guacamole WITH PASSWORD 'changeme_db_password';"
docker exec -it remotegpu-postgresql psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE guacamole TO guacamole;"

# 3. 导入初始化脚本
docker exec -i remotegpu-postgresql psql -U guacamole -d guacamole < initdb.sql
```

## 启动服务

```bash
docker-compose up -d
```

## 访问控制台

- 访问地址: http://localhost:8081/guacamole/
- 默认用户名: guacadmin
- 默认密码: guacadmin

**首次登录后请立即修改密码！**

## 支持的协议

- RDP (远程桌面)
- VNC
- SSH
- Telnet

## 停止服务

```bash
docker-compose down
```
