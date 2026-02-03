# Redis 缓存服务

## 启动服务

```bash
docker-compose up -d
```

## 验证服务

```bash
# 检查健康状态
docker exec remotegpu-redis redis-cli ping

# 连接 Redis（需要密码）
docker exec -it remotegpu-redis redis-cli -a changeme_redis_password

# 查看信息
docker exec remotegpu-redis redis-cli -a changeme_redis_password INFO
```

## 配置说明

### 修改密码

编辑 `redis.conf` 文件：

```ini
requirepass your_new_password
```

### 内存配置

默认最大内存为 4GB，可根据实际需求调整：

```ini
maxmemory 4gb
maxmemory-policy allkeys-lru
```

## 使用场景

- 用户会话存储
- JWT Token 黑名单
- 端口映射缓存
- 分布式锁
- 消息队列

## 数据备份

```bash
# 触发 RDB 快照
docker exec remotegpu-redis redis-cli -a changeme_redis_password BGSAVE

# 复制备份文件
docker cp remotegpu-redis:/data/dump.rdb ./redis-backup.rdb
```

## 停止服务

```bash
docker-compose down
```
