# RustFS K8s 部署配置

## 概述

本目录包含在 Kubernetes 中部署 RustFS 对象存储服务的配置文件。

RustFS 是一个 S3 兼容的对象存储服务，提供：
- S3 API 接口（端口 9000）
- Web 控制台（端口 9001）
- 持久化存储

## 配置文件

| 文件 | 说明 |
|------|------|
| `secret.yaml` | 存储访问密钥（ACCESS_KEY 和 SECRET_KEY） |
| `pvc.yaml` | 持久化存储卷声明（50Gi） |
| `deployment.yaml` | RustFS 部署配置 |
| `service.yaml` | Service 服务暴露 |

## 部署步骤

### 1. 创建 Secret

```bash
kubectl apply -f secret.yaml
```

### 2. 创建 PVC

```bash
kubectl apply -f pvc.yaml
```

### 3. 部署 RustFS

```bash
kubectl apply -f deployment.yaml
```

### 4. 创建 Service

```bash
kubectl apply -f service.yaml
```

### 一键部署

```bash
kubectl apply -f .
```

## 验证部署

### 检查 Pod 状态

```bash
kubectl get pods -l app=rustfs
```

### 检查 Service

```bash
kubectl get svc rustfs
```

### 查看日志

```bash
kubectl logs -l app=rustfs -f
```

## 访问 RustFS

### 集群内访问

- **S3 API**: `http://rustfs.default.svc.cluster.local:9000`
- **控制台**: `http://rustfs.default.svc.cluster.local:9001`

### 集群外访问

如需从集群外访问，可以使用以下方式：

#### 方式1: Port Forward

```bash
# S3 API
kubectl port-forward svc/rustfs 9000:9000

# 控制台
kubectl port-forward svc/rustfs 9001:9001
```

#### 方式2: NodePort

修改 `service.yaml`，将 `type: ClusterIP` 改为 `type: NodePort`：

```yaml
spec:
  type: NodePort
  ports:
  - port: 9000
    targetPort: 9000
    nodePort: 30900  # 可选，指定端口
    name: api
  - port: 9001
    targetPort: 9001
    nodePort: 30901  # 可选，指定端口
    name: console
```

#### 方式3: Ingress

创建 Ingress 资源（需要 Ingress Controller）。

## 配置说明

### 访问密钥

默认访问密钥配置在 `secret.yaml` 中：
- **Access Key**: `rustfsadmin`
- **Secret Key**: `rustfsadmin`

**⚠️ 生产环境请务必修改为强密码！**

修改方式：
```bash
# 编辑 secret.yaml
kubectl edit secret rustfs-secret
```

### 存储容量

默认 PVC 容量为 50Gi，可在 `pvc.yaml` 中修改：

```yaml
resources:
  requests:
    storage: 50Gi  # 修改为所需容量
```

### 资源限制

默认资源配置：
- **Requests**: 512Mi 内存, 500m CPU
- **Limits**: 2Gi 内存, 2000m CPU

可在 `deployment.yaml` 中修改。

## 使用示例

### 使用 AWS CLI

```bash
# 配置 AWS CLI
aws configure set aws_access_key_id rustfsadmin
aws configure set aws_secret_access_key rustfsadmin
aws configure set default.region us-east-1

# 创建 bucket
aws --endpoint-url http://localhost:9000 s3 mb s3://my-bucket

# 上传文件
aws --endpoint-url http://localhost:9000 s3 cp file.txt s3://my-bucket/

# 列出文件
aws --endpoint-url http://localhost:9000 s3 ls s3://my-bucket/
```

### 使用 Go SDK

```go
import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

// 创建 S3 客户端
sess, _ := session.NewSession(&aws.Config{
    Endpoint:         aws.String("http://rustfs.default.svc.cluster.local:9000"),
    Region:           aws.String("us-east-1"),
    Credentials:      credentials.NewStaticCredentials("rustfsadmin", "rustfsadmin", ""),
    S3ForcePathStyle: aws.Bool(true),
})

svc := s3.New(sess)
```

## 故障排查

### Pod 无法启动

```bash
# 查看 Pod 详情
kubectl describe pod -l app=rustfs

# 查看日志
kubectl logs -l app=rustfs
```

### PVC 无法绑定

```bash
# 检查 PVC 状态
kubectl get pvc rustfs-data

# 查看 PVC 详情
kubectl describe pvc rustfs-data
```

常见原因：
- 没有可用的 StorageClass
- 存储容量不足
- 权限问题

### 健康检查失败

检查 RustFS 是否正常响应：

```bash
kubectl exec -it <rustfs-pod-name> -- curl http://localhost:9000/health
```

## 卸载

```bash
# 删除所有资源
kubectl delete -f .

# 或逐个删除
kubectl delete deployment rustfs
kubectl delete service rustfs
kubectl delete pvc rustfs-data
kubectl delete secret rustfs-secret
```

**⚠️ 注意**: 删除 PVC 会导致数据丢失！

## 生产环境建议

1. **修改默认密钥**: 使用强密码
2. **配置备份**: 定期备份 PVC 数据
3. **资源监控**: 监控 CPU、内存、存储使用情况
4. **高可用**: 考虑使用 StatefulSet 和多副本
5. **网络安全**: 配置 NetworkPolicy 限制访问
6. **TLS 加密**: 配置 HTTPS 访问

## 参考资料

- [RustFS 官方文档](https://github.com/rustfs/rustfs)
- [S3 API 文档](https://docs.aws.amazon.com/s3/)
