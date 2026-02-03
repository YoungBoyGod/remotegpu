# RustFS 对象存储

RustFS 是一个高性能的对象存储服务，提供 S3 兼容的 API。

## 启动服务

```bash
docker-compose up -d
```

## 访问控制台

- API 端点: http://localhost:9000
- 管理控制台: http://localhost:9001
- 默认用户名: admin
- 默认密码: changeme_rustfs_password

**首次登录后请立即修改密码！**

## 环境变量配置

在 `.env` 文件中配置：

```bash
RUSTFS_ROOT_USER=admin
RUSTFS_ROOT_PASSWORD=your_secure_password
```

## 使用场景

- 数据集存储
- 模型文件存储
- 训练日志存储
- 镜像备份

## 验证服务

```bash
# 检查健康状态
curl http://localhost:9000/health

# 查看容器日志
docker-compose logs -f
```

## 停止服务

```bash
docker-compose down
```

## 注意事项

- 生产环境必须修改默认密码
- 定期备份数据
- 监控磁盘空间使用
