# Prometheus 监控

## 启动服务

```bash
docker-compose up -d
```

## 访问控制台

- 访问地址: http://localhost:19090
- 无需登录

## 配置说明

### 添加监控目标

编辑 `prometheus.yml` 文件，添加新的监控目标：

```yaml
scrape_configs:
  - job_name: 'my-service'
    static_configs:
      - targets: ['service-host:port']
        labels:
          environment: 'production'
```

### 常用查询示例

```promql
# CPU 使用率
100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)

# 内存使用率
(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100

# 磁盘使用率
(1 - (node_filesystem_avail_bytes / node_filesystem_size_bytes)) * 100
```

## 数据保留

默认保留30天数据，可在 docker-compose.yml 中修改：

```yaml
command:
  - '--storage.tsdb.retention.time=30d'
```

## 停止服务

```bash
docker-compose down
```
