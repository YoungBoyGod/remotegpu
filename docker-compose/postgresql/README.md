# PostgreSQL 数据库

## 启动服务

```bash
docker-compose up -d
```

## 验证服务

```bash
# 检查健康状态
docker exec remotegpu-postgresql pg_isready -U remotegpu_user -d remotegpu

# 连接数据库
docker exec -it remotegpu-postgresql psql -U remotegpu_user -d remotegpu

# 查看数据库列表
docker exec remotegpu-postgresql psql -U remotegpu_user -c "\l"
```

## 配置说明

### 环境变量

在 `.env` 文件中配置：

```bash
POSTGRES_PASSWORD=your_secure_password
```

### 性能优化

`postgresql.conf` 文件包含了针对生产环境的优化配置：
- 连接数：200
- 共享缓冲区：4GB
- 有效缓存：12GB
- WAL 配置：支持主从复制

根据实际硬件资源调整这些参数。

## 数据备份

```bash
# 备份数据库
docker exec remotegpu-postgresql pg_dump -U remotegpu_user remotegpu > backup.sql

# 恢复数据库
docker exec -i remotegpu-postgresql psql -U remotegpu_user remotegpu < backup.sql
```

## 停止服务

```bash
docker-compose down
```

## 注意事项

- 生产环境必须修改默认密码
- 定期备份数据库
- 监控磁盘空间使用
- 定期清理日志文件
