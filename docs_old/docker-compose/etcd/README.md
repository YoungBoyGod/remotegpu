# Etcd 配置中心

## 启动服务

```bash
docker-compose up -d
```

## 验证服务

```bash
# 检查健康状态
docker exec remotegpu-etcd etcdctl endpoint health

# 查看成员列表
docker exec remotegpu-etcd etcdctl member list

# 测试写入
docker exec remotegpu-etcd etcdctl put test "hello"
docker exec remotegpu-etcd etcdctl get test
```

## 停止服务

```bash
docker-compose down
```
