# 方案B：子域名方案 - 总体架构

## 架构概览

**目标**：200台GPU机器，每台5个服务，使用独立子域名访问

**子域名数量**：1000个（通过泛域名一次性配置）

---

## 访问示例

### GPU1
- SSH: `ssh user@gpu1.gpu.domain.com`
- Jupyter: `https://gpu1-jupyter.gpu.domain.com`
- TensorBoard: `https://gpu1-tensorboard.gpu.domain.com`
- 服务1: `https://gpu1-service1.gpu.domain.com`
- 服务2: `https://gpu1-service2.gpu.domain.com`

### GPU2-200
同样的命名规则

---

## 技术架构

```
用户
  ↓
泛域名 *.gpu.domain.com（DNS）
  ↓
云服务器（nginx + frps + SSL）
  ↓
frp 隧道（1000个端口映射）
  ↓
本地服务器（frpc）
  ↓
200台 GPU 机器（每台5个服务）
```

---

## 端口分配

| GPU编号 | SSH | Jupyter | TensorBoard | 服务1 | 服务2 |
|---------|-----|---------|-------------|-------|-------|
| GPU1 | 2201 | 8001 | 9001 | 10001 | 11001 |
| GPU2 | 2202 | 8002 | 9002 | 10002 | 11002 |
| ... | ... | ... | ... | ... | ... |
| GPU200 | 2400 | 8200 | 9200 | 10200 | 11200 |

---

## 实施步骤

1. DNS配置（泛域名）
2. SSL证书（泛域名证书）
3. frp服务端配置
4. frp客户端配置（批量）
5. nginx配置（子域名匹配）

详见后续文档。
