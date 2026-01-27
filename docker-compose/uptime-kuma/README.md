# Uptime Kuma 服务监控

## 启动服务

```bash
docker-compose up -d
```

## 访问控制台

- 访问地址: http://localhost:3001
- 首次访问需要创建管理员账户

## 监控配置建议

### HTTP/HTTPS 监控
- API 服务健康检查
- Web 前端可用性
- JupyterLab 服务状态

### TCP 端口监控
- PostgreSQL (5432)
- Redis (6379)
- Kubernetes API (6443)

### Ping 监控
- GPU 节点存活检查
- 网络连通性检查

## 停止服务

```bash
docker-compose down
```
