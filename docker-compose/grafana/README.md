# Grafana 可视化

## 启动服务

```bash
docker-compose up -d
```

## 访问控制台

- 访问地址: http://localhost:13000
- 默认用户名: admin
- 默认密码: changeme_grafana_password

**首次登录后请立即修改密码！**

## 配置数据源

### 添加 Prometheus 数据源

1. 登录 Grafana
2. 进入 Configuration > Data Sources
3. 点击 "Add data source"
4. 选择 "Prometheus"
5. 配置 URL: http://prometheus:9090
6. 点击 "Save & Test"

## 导入仪表板

Grafana 提供了大量预制仪表板：

1. 进入 Dashboards > Import
2. 输入仪表板 ID（例如：1860 - Node Exporter Full）
3. 选择数据源
4. 点击 "Import"

### 推荐仪表板

- Node Exporter Full: 1860
- Docker Monitoring: 893
- Kubernetes Cluster: 7249
- NVIDIA GPU: 12239

## 环境变量配置

在 `.env` 文件中配置：

```bash
GF_ADMIN_USER=admin
GF_ADMIN_PASSWORD=your_secure_password
```

## 停止服务

```bash
docker-compose down
```
