# PVC 挂载示例

## 概述

本示例演示如何创建一个 1Gi 的 PVC 并挂载到 Pod 中。

## 配置文件

| 文件 | 说明 |
|------|------|
| `pvc-1g.yaml` | 1Gi 持久化存储卷声明 |
| `pod-with-pvc.yaml` | 挂载 PVC 的 Nginx Pod |

## 部署步骤

### 1. 创建 PVC

```bash
kubectl apply -f pvc-1g.yaml
```

### 2. 验证 PVC 状态

```bash
kubectl get pvc test-pvc-1g
```

等待状态变为 `Bound`。

### 3. 创建 Pod

```bash
kubectl apply -f pod-with-pvc.yaml
```

### 4. 验证 Pod 状态

```bash
kubectl get pod test-pod-with-pvc
```

## 验证挂载

### 进入 Pod 查看挂载点

```bash
kubectl exec -it test-pod-with-pvc -- /bin/bash
```

在 Pod 内执行：

```bash
# 查看挂载点
df -h | grep /data

# 创建测试文件
echo "Hello PVC" > /data/test.txt

# 查看文件
cat /data/test.txt
```

### 测试数据持久化

```bash
# 删除 Pod
kubectl delete pod test-pod-with-pvc

# 重新创建 Pod
kubectl apply -f pod-with-pvc.yaml

# 进入新 Pod 验证数据
kubectl exec -it test-pod-with-pvc -- cat /data/test.txt
```

数据应该仍然存在，证明 PVC 持久化成功。

## 清理资源

```bash
# 删除 Pod
kubectl delete -f pod-with-pvc.yaml

# 删除 PVC（会删除数据）
kubectl delete -f pvc-1g.yaml
```

## 配置说明

### PVC 配置

- **名称**: test-pvc-1g
- **容量**: 1Gi
- **访问模式**: ReadWriteOnce（单节点读写）
- **存储类**: standard

### Pod 配置

- **镜像**: nginx:latest
- **挂载路径**: /data
- **PVC**: test-pvc-1g

## 常见问题

### PVC 一直处于 Pending 状态

原因：
- 没有可用的 StorageClass
- 存储容量不足

解决：
```bash
# 查看 StorageClass
kubectl get storageclass

# 查看 PVC 详情
kubectl describe pvc test-pvc-1g
```

### Pod 无法启动

```bash
# 查看 Pod 详情
kubectl describe pod test-pod-with-pvc

# 查看日志
kubectl logs test-pod-with-pvc
```
