# 数据与镜像管理技术方案

> 对象存储、数据集管理、镜像管理完整技术方案
>
> 创建日期：2026-01-26

---

## 目录

1. [整体架构](#1-整体架构)
2. [对象存储方案](#2-对象存储方案)
3. [数据集管理](#3-数据集管理)
4. [模型管理](#4-模型管理)
5. [镜像管理](#5-镜像管理)
6. [文件系统挂载](#6-文件系统挂载)
7. [传统架构 vs K8s 实现](#7-传统架构-vs-k8s-实现)

---

## 1. 整体架构

### 1.1 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                      用户界面层                               │
│  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐            │
│  │数据集  │  │ 模型库 │  │ 镜像库 │  │文件浏览 │            │
│  │管理    │  │        │  │        │  │        │            │
│  └────────┘  └────────┘  └────────┘  └────────┘            │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                      API 服务层                               │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │数据集 API  │  │ 模型 API   │  │ 镜像 API   │            │
│  └────────────┘  └────────────┘  └────────────┘            │
└───────────────────────┬─────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        │               │               │
        ▼               ▼               ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│  对象存储     │ │  镜像仓库     │ │  数据库       │
│  (MinIO)     │ │  (Registry)  │ │ (Postgres)   │
│              │ │              │ │              │
│ /datasets/   │ │ docker.io/   │ │ metadata     │
│ /models/     │ │ registry:2   │ │              │
│ /artifacts/  │ │              │ │              │
└──────────────┘ └──────────────┘ └──────────────┘
        │               │
        │               │
        └───────┬───────┘
                │
                ▼
┌─────────────────────────────────────────────────────────────┐
│                   开发环境容器                                 │
│  ┌─────────────────────────────────────────────────────────┐│
│  │  挂载点：                                                ││
│  │  /gemini/data-1/     -> 数据集 1                        ││
│  │  /gemini/data-2/     -> 数据集 2                        ││
│  │  /gemini/pretrain1/  -> 模型 1                          ││
│  │  /gemini/code/       -> 用户代码                        ││
│  │  /gemini/output/     -> 训练输出                        ││
│  └─────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

### 1.2 核心组件

| 组件 | 功能 | 技术选型 |
|------|------|---------|
| **对象存储** | 存储数据集、模型、制品 | MinIO / Ceph / S3 |
| **镜像仓库** | 存储 Docker 镜像 | Harbor / Registry 2.0 |
| **元数据库** | 存储文件元信息 | PostgreSQL |
| **文件系统** | 容器内挂载 | FUSE / 直接挂载 |
| **上传服务** | 处理大文件上传 | Tus / Resumable.js |
| **镜像构建** | 构建自定义镜像 | Kaniko / BuildKit |

---

## 2. 对象存储方案

### 2.1 方案对比

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|---------|
| **MinIO** | 开源、S3 兼容、易部署 | 需要自己运维 | 自建机房、私有云 |
| **Ceph** | 分布式、高可用 | 复杂、学习曲线高 | 大规模集群 |
| **AWS S3** | 稳定、免运维 | 成本高、供应商锁定 | 海外业务 |
| **阿里云 OSS** | 稳定、国内快 | 成本较高 | 国内业务、快速上线 |
| **腾讯云 COS** | 稳定、价格适中 | 供应商锁定 | 国内业务 |
| **本地存储** | 简单、成本低 | 不易扩展、单点故障 | MVP 阶段 |

**推荐选择：MinIO**（开源、S3 兼容、后期可迁移到云存储）

### 2.2 MinIO 部署

#### 单机部署（开发/测试）

```bash
# Docker 方式部署
docker run -d \
  --name minio \
  -p 9000:9000 \
  -p 9001:9001 \
  -e MINIO_ROOT_USER=admin \
  -e MINIO_ROOT_PASSWORD=your-password \
  -v /data/minio:/data \
  minio/minio server /data --console-address ":9001"
```

#### 分布式部署（生产环境）

```yaml
# docker-compose.yml
version: '3.8'

services:
  minio1:
    image: minio/minio
    volumes:
      - /mnt/disk1:/data1
      - /mnt/disk2:/data2
    command: server http://minio{1...4}/data{1...2} --console-address ":9001"
    environment:
      MINIO_ROOT_USER: admin
      MINIO_ROOT_PASSWORD: your-password

  minio2:
    image: minio/minio
    volumes:
      - /mnt/disk3:/data1
      - /mnt/disk4:/data2
    command: server http://minio{1...4}/data{1...2} --console-address ":9001"
    environment:
      MINIO_ROOT_USER: admin
      MINIO_ROOT_PASSWORD: your-password

  # ... minio3, minio4
```

#### Kubernetes 部署

```bash
# 使用 Helm 部署
helm repo add minio https://charts.min.io/
helm install minio minio/minio \
  --set rootUser=admin \
  --set rootPassword=your-password \
  --set mode=distributed \
  --set replicas=4 \
  --set persistence.size=1Ti
```

### 2.3 存储桶设计

```
MinIO 存储桶结构：

remotegpu/
├── datasets/                    # 数据集存储
│   ├── customer-123/
│   │   ├── dataset-001/
│   │   │   ├── version-1/
│   │   │   │   ├── train/
│   │   │   │   ├── test/
│   │   │   │   └── metadata.json
│   │   │   └── version-2/
│   │   └── dataset-002/
│   └── customer-456/
│
├── models/                      # 模型存储
│   ├── customer-123/
│   │   ├── model-001/
│   │   │   ├── v1.0/
│   │   │   │   ├── model.pth
│   │   │   │   ├── config.json
│   │   │   │   └── README.md
│   │   │   └── v2.0/
│   │   └── model-002/
│   └── pretrained/              # 预训练模型（公共）
│       ├── bert-base/
│       └── resnet50/
│
├── artifacts/                   # 制品/训练输出
│   ├── customer-123/
│   │   ├── training-001/
│   │   │   ├── checkpoints/
│   │   │   ├── logs/
│   │   │   └── results/
│   │   └── training-002/
│   └── customer-456/
│
└── images/                      # 自定义镜像构建产物（可选）
    ├── customer-123/
    │   └── custom-pytorch/
    └── customer-456/
```

### 2.4 Go SDK 封装

```go
// storage/minio_client.go
package storage

import (
    "context"
    "io"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
    client     *minio.Client
    bucketName string
}

// 初始化客户端
func NewMinioClient(endpoint, accessKey, secretKey, bucketName string) (*MinioClient, error) {
    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: false, // 生产环境设置为 true
    })
    if err != nil {
        return nil, err
    }

    // 确保桶存在
    exists, err := client.BucketExists(context.Background(), bucketName)
    if err != nil {
        return nil, err
    }
    if !exists {
        err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
        if err != nil {
            return nil, err
        }
    }

    return &MinioClient{
        client:     client,
        bucketName: bucketName,
    }, nil
}

// 上传文件
func (m *MinioClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64) error {
    _, err := m.client.PutObject(
        ctx,
        m.bucketName,
        objectName,
        reader,
        size,
        minio.PutObjectOptions{
            ContentType: "application/octet-stream",
        },
    )
    return err
}

// 下载文件
func (m *MinioClient) DownloadFile(ctx context.Context, objectName string) (*minio.Object, error) {
    return m.client.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
}

// 删除文件
func (m *MinioClient) DeleteFile(ctx context.Context, objectName string) error {
    return m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
}

// 列出文件
func (m *MinioClient) ListFiles(ctx context.Context, prefix string) ([]string, error) {
    var files []string
    objectCh := m.client.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{
        Prefix:    prefix,
        Recursive: true,
    })

    for object := range objectCh {
        if object.Err != nil {
            return nil, object.Err
        }
        files = append(files, object.Key)
    }

    return files, nil
}

// 生成预签名 URL（用于前端直传）
func (m *MinioClient) GetPresignedUploadURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
    return m.client.PresignedPutObject(ctx, m.bucketName, objectName, expires)
}

// 生成预签名下载 URL
func (m *MinioClient) GetPresignedDownloadURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
    return m.client.PresignedGetObject(ctx, m.bucketName, objectName, expires, nil)
}
```

---

## 3. 数据集管理

### 3.1 数据库设计

```sql
-- 数据集表
CREATE TABLE datasets (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR(64) UNIQUE NOT NULL,
    customer_id BIGINT NOT NULL,
    workspace_id BIGINT NOT NULL,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    visibility VARCHAR(20) DEFAULT 'private',  -- private, workspace, public
    storage_path VARCHAR(512) NOT NULL,        -- datasets/customer-123/dataset-001/
    total_size BIGINT DEFAULT 0,               -- 总大小（字节）
    file_count INT DEFAULT 0,                  -- 文件数量
    status VARCHAR(20) DEFAULT 'uploading',    -- uploading, ready, error
    tags TEXT[],                               -- 标签数组
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_datasets_customer_id ON datasets (customer_id);
CREATE INDEX idx_datasets_workspace_id ON datasets (workspace_id);
CREATE INDEX idx_datasets_visibility ON datasets (visibility);

-- 数据集版本表
CREATE TABLE dataset_versions (
    id BIGSERIAL PRIMARY KEY,
    dataset_id BIGINT NOT NULL,
    version VARCHAR(64) NOT NULL,              -- v1, v2, latest
    storage_path VARCHAR(512) NOT NULL,
    size BIGINT DEFAULT 0,
    file_count INT DEFAULT 0,
    description TEXT,
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (dataset_id) REFERENCES datasets(id) ON DELETE CASCADE,
    UNIQUE(dataset_id, version)
);

-- 数据集文件表（可选，用于详细跟踪）
CREATE TABLE dataset_files (
    id BIGSERIAL PRIMARY KEY,
    dataset_id BIGINT NOT NULL,
    version_id BIGINT,
    file_path VARCHAR(512) NOT NULL,           -- 相对路径
    file_size BIGINT NOT NULL,
    file_type VARCHAR(64),                     -- image/jpeg, text/csv, etc.
    checksum VARCHAR(128),                     -- MD5/SHA256
    uploaded_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (dataset_id) REFERENCES datasets(id) ON DELETE CASCADE
);

CREATE INDEX idx_dataset_files_dataset_id ON dataset_files (dataset_id);

-- 数据集使用记录
CREATE TABLE dataset_usage (
    id BIGSERIAL PRIMARY KEY,
    dataset_id BIGINT NOT NULL,
    env_id VARCHAR(64),                        -- 开发环境/训练任务 ID
    customer_id BIGINT NOT NULL,
    mount_path VARCHAR(256),                   -- /gemini/data-1/
    mounted_at TIMESTAMP DEFAULT NOW(),
    unmounted_at TIMESTAMP,
    FOREIGN KEY (dataset_id) REFERENCES datasets(id) ON DELETE CASCADE
);
```

### 3.2 API 设计

```go
// api/dataset.go
package api

// 创建数据集
// POST /api/datasets
func CreateDataset(c *gin.Context) {
    var req struct {
        WorkspaceID int64    `json:"workspace_id" binding:"required"`
        Name        string   `json:"name" binding:"required"`
        Description string   `json:"description"`
        Visibility  string   `json:"visibility"`
        Tags        []string `json:"tags"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    customer := GetCurrentCustomer(c)

    // 1. 生成 UUID
    datasetUUID := uuid.New().String()

    // 2. 生成存储路径
    storagePath := fmt.Sprintf("datasets/customer-%d/dataset-%s/", customer.ID, datasetUUID)

    // 3. 创建数据库记录
    dataset := &Dataset{
        UUID:        datasetUUID,
        CustomerID:  customer.ID,
        WorkspaceID: req.WorkspaceID,
        Name:        req.Name,
        Description: req.Description,
        Visibility:  req.Visibility,
        StoragePath: storagePath,
        Status:      "uploading",
        Tags:        req.Tags,
    }

    if err := db.Create(dataset).Error; err != nil {
        c.JSON(500, gin.H{"error": "创建失败"})
        return
    }

    c.JSON(200, gin.H{
        "dataset_id":   dataset.UUID,
        "storage_path": storagePath,
    })
}

// 获取上传凭证（预签名 URL）
// POST /api/datasets/:id/upload-url
func GetUploadURL(c *gin.Context) {
    datasetID := c.Param("id")

    var req struct {
        FileName string `json:"file_name" binding:"required"`
        FileSize int64  `json:"file_size"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 1. 查询数据集
    var dataset Dataset
    if err := db.Where("uuid = ? AND customer_id = ?", datasetID, GetCurrentCustomer(c).ID).First(&dataset).Error; err != nil {
        c.JSON(404, gin.H{"error": "数据集不存在"})
        return
    }

    // 2. 生成对象存储路径
    objectName := filepath.Join(dataset.StoragePath, req.FileName)

    // 3. 生成预签名上传 URL（有效期 1 小时）
    uploadURL, err := minioClient.GetPresignedUploadURL(c.Request.Context(), objectName, 1*time.Hour)
    if err != nil {
        c.JSON(500, gin.H{"error": "生成上传链接失败"})
        return
    }

    c.JSON(200, gin.H{
        "upload_url": uploadURL,
        "object_name": objectName,
        "expires_in": 3600,
    })
}

// 完成上传（更新数据集状态）
// POST /api/datasets/:id/complete
func CompleteUpload(c *gin.Context) {
    datasetID := c.Param("id")

    var req struct {
        Files []struct {
            FileName string `json:"file_name"`
            FileSize int64  `json:"file_size"`
        } `json:"files"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 1. 查询数据集
    var dataset Dataset
    if err := db.Where("uuid = ?", datasetID).First(&dataset).Error; err != nil {
        c.JSON(404, gin.H{"error": "数据集不存在"})
        return
    }

    // 2. 计算总大小
    var totalSize int64
    for _, file := range req.Files {
        totalSize += file.FileSize
    }

    // 3. 更新数据集状态
    db.Model(&dataset).Updates(map[string]interface{}{
        "status":     "ready",
        "total_size": totalSize,
        "file_count": len(req.Files),
    })

    c.JSON(200, gin.H{"status": "success"})
}

// 列出数据集
// GET /api/datasets
func ListDatasets(c *gin.Context) {
    customer := GetCurrentCustomer(c)

    var datasets []Dataset
    query := db.Where("customer_id = ?", customer.ID)

    // 筛选条件
    if visibility := c.Query("visibility"); visibility != "" {
        query = query.Where("visibility = ?", visibility)
    }

    if tag := c.Query("tag"); tag != "" {
        query = query.Where("? = ANY(tags)", tag)
    }

    query.Order("created_at DESC").Find(&datasets)

    c.JSON(200, gin.H{
        "datasets": datasets,
    })
}

// 删除数据集
// DELETE /api/datasets/:id
func DeleteDataset(c *gin.Context) {
    datasetID := c.Param("id")
    customer := GetCurrentCustomer(c)

    var dataset Dataset
    if err := db.Where("uuid = ? AND customer_id = ?", datasetID, customer.ID).First(&dataset).Error; err != nil {
        c.JSON(404, gin.H{"error": "数据集不存在"})
        return
    }

    // 1. 检查是否正在使用
    var usageCount int64
    db.Model(&DatasetUsage{}).Where("dataset_id = ? AND unmounted_at IS NULL", dataset.ID).Count(&usageCount)
    if usageCount > 0 {
        c.JSON(400, gin.H{"error": "数据集正在使用中，无法删除"})
        return
    }

    // 2. 删除对象存储中的文件
    go func() {
        files, _ := minioClient.ListFiles(context.Background(), dataset.StoragePath)
        for _, file := range files {
            minioClient.DeleteFile(context.Background(), file)
        }
    }()

    // 3. 删除数据库记录
    db.Delete(&dataset)

    c.JSON(200, gin.H{"status": "deleted"})
}
```

### 3.3 大文件上传方案

#### 方案 A：前端直传（推荐）

```javascript
// 前端实现（React）
async function uploadDataset(datasetId, files) {
  for (const file of files) {
    // 1. 获取预签名上传 URL
    const { upload_url, object_name } = await fetch(
      `/api/datasets/${datasetId}/upload-url`,
      {
        method: 'POST',
        body: JSON.stringify({
          file_name: file.name,
          file_size: file.size,
        }),
      }
    ).then(r => r.json());

    // 2. 直接上传到对象存储
    await fetch(upload_url, {
      method: 'PUT',
      body: file,
      headers: {
        'Content-Type': file.type,
      },
    });
  }

  // 3. 通知后端上传完成
  await fetch(`/api/datasets/${datasetId}/complete`, {
    method: 'POST',
    body: JSON.stringify({
      files: files.map(f => ({
        file_name: f.name,
        file_size: f.size,
      })),
    }),
  });
}
```

#### 方案 B：分片上传（超大文件）

```go
// 使用 Tus 协议实现断点续传
import (
    "github.com/tus/tusd/pkg/filestore"
    "github.com/tus/tusd/pkg/handler"
)

func SetupTusServer() {
    store := filestore.New("/tmp/tusd")
    composer := handler.NewStoreComposer()
    store.UseIn(composer)

    h, err := handler.NewHandler(handler.Config{
        BasePath:      "/files/",
        StoreComposer: composer,
    })

    http.Handle("/files/", http.StripPrefix("/files/", h))
}
```

### 3.4 文件浏览 API

```go
// 浏览数据集文件
// GET /api/datasets/:id/files
func BrowseDatasetFiles(c *gin.Context) {
    datasetID := c.Param("id")
    prefix := c.Query("prefix") // 子目录

    var dataset Dataset
    if err := db.Where("uuid = ?", datasetID).First(&dataset).Error; err != nil {
        c.JSON(404, gin.H{"error": "数据集不存在"})
        return
    }

    // 列出文件
    fullPrefix := filepath.Join(dataset.StoragePath, prefix)
    objects := minioClient.client.ListObjects(c, minioClient.bucketName, minio.ListObjectsOptions{
        Prefix: fullPrefix,
    })

    var files []map[string]interface{}
    for obj := range objects {
        if obj.Err != nil {
            continue
        }

        files = append(files, map[string]interface{}{
            "name":         filepath.Base(obj.Key),
            "path":         strings.TrimPrefix(obj.Key, dataset.StoragePath+"/"),
            "size":         obj.Size,
            "last_modified": obj.LastModified,
            "is_dir":       strings.HasSuffix(obj.Key, "/"),
        })
    }

    c.JSON(200, gin.H{
        "files": files,
    })
}

// 下载单个文件
// GET /api/datasets/:id/download
func DownloadFile(c *gin.Context) {
    datasetID := c.Param("id")
    filePath := c.Query("file")

    var dataset Dataset
    if err := db.Where("uuid = ?", datasetID).First(&dataset).Error; err != nil {
        c.JSON(404, gin.H{"error": "数据集不存在"})
        return
    }

    // 生成预签名下载 URL
    objectName := filepath.Join(dataset.StoragePath, filePath)
    downloadURL, err := minioClient.GetPresignedDownloadURL(c.Request.Context(), objectName, 1*time.Hour)
    if err != nil {
        c.JSON(500, gin.H{"error": "生成下载链接失败"})
        return
    }

    c.JSON(200, gin.H{
        "download_url": downloadURL,
        "expires_in": 3600,
    })
}
```

---

## 4. 模型管理

### 4.1 数据库设计

```sql
-- 模型表
CREATE TABLE models (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR(64) UNIQUE NOT NULL,
    customer_id BIGINT NOT NULL,
    workspace_id BIGINT NOT NULL,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    framework VARCHAR(64),                     -- pytorch, tensorflow, onnx
    task_type VARCHAR(64),                     -- classification, detection, nlp, etc.
    visibility VARCHAR(20) DEFAULT 'private',
    storage_path VARCHAR(512) NOT NULL,
    total_size BIGINT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'uploading',
    tags TEXT[],
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_models_customer_id ON models (customer_id);
CREATE INDEX idx_models_framework ON models (framework);

-- 模型版本表
CREATE TABLE model_versions (
    id BIGSERIAL PRIMARY KEY,
    model_id BIGINT NOT NULL,
    version VARCHAR(64) NOT NULL,
    storage_path VARCHAR(512) NOT NULL,
    size BIGINT DEFAULT 0,
    accuracy FLOAT,                            -- 模型指标
    precision FLOAT,
    recall FLOAT,
    metrics JSONB,                             -- 其他指标
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (model_id) REFERENCES models(id) ON DELETE CASCADE,
    UNIQUE(model_id, version)
);
```

### 4.2 预训练模型库

```go
// 预训练模型配置
var PretrainedModels = []PretrainedModel{
    {
        Name:        "BERT Base",
        Framework:   "pytorch",
        Description: "BERT base model (110M parameters)",
        URL:         "s3://pretrained/bert-base/",
        Size:        440 * 1024 * 1024, // 440MB
    },
    {
        Name:        "ResNet-50",
        Framework:   "pytorch",
        Description: "ResNet-50 pretrained on ImageNet",
        URL:         "s3://pretrained/resnet50/",
        Size:        98 * 1024 * 1024, // 98MB
    },
    // ... 更多预训练模型
}

// 同步预训练模型到用户空间
// POST /api/models/pretrained/:name/sync
func SyncPretrainedModel(c *gin.Context) {
    modelName := c.Param("name")
    customer := GetCurrentCustomer(c)

    // 查找预训练模型
    var pretrainedModel PretrainedModel
    for _, pm := range PretrainedModels {
        if pm.Name == modelName {
            pretrainedModel = pm
            break
        }
    }

    if pretrainedModel.Name == "" {
        c.JSON(404, gin.H{"error": "预训练模型不存在"})
        return
    }

    // 创建模型记录（指向共享存储）
    model := &Model{
        UUID:        uuid.New().String(),
        CustomerID:  customer.ID,
        Name:        pretrainedModel.Name,
        Description: pretrainedModel.Description,
        Framework:   pretrainedModel.Framework,
        StoragePath: pretrainedModel.URL,
        TotalSize:   pretrainedModel.Size,
        Status:      "ready",
        Visibility:  "private",
    }

    db.Create(model)

    c.JSON(200, gin.H{
        "model_id": model.UUID,
    })
}
```

---

## 5. 镜像管理

### 5.1 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                     镜像仓库架构                               │
└─────────────────────────────────────────────────────────────┘

┌──────────────────┐         ┌──────────────────┐
│   官方镜像库      │         │  用户自定义镜像   │
├──────────────────┤         ├──────────────────┤
│ ubuntu20-py38    │         │ customer-123/    │
│ ubuntu20-pytorch │         │   custom-torch   │
│ ubuntu20-tf      │         ├──────────────────┤
│ ubuntu22-py310   │         │ customer-456/    │
└──────────────────┘         │   my-image       │
                             └──────────────────┘
           │                           │
           └───────────┬───────────────┘
                       │
                       ▼
        ┌──────────────────────────────┐
        │   镜像仓库 (Harbor/Registry)  │
        │   registry.example.com       │
        └──────────────────────────────┘
```

### 5.2 镜像仓库部署

#### 使用 Harbor（推荐）

```yaml
# docker-compose.yml for Harbor
version: '3.8'

services:
  harbor-core:
    image: goharbor/harbor-core:v2.8.0
    # ... 配置

  harbor-registry:
    image: goharbor/registry-photon:v2.8.0
    volumes:
      - /data/registry:/storage

  harbor-db:
    image: goharbor/harbor-db:v2.8.0
    volumes:
      - /data/database:/var/lib/postgresql/data
```

#### 使用原生 Registry

```bash
# 启动 Docker Registry
docker run -d \
  --name registry \
  -p 5000:5000 \
  -v /data/registry:/var/lib/registry \
  registry:2
```

### 5.3 官方镜像管理

#### Dockerfile 模板

```dockerfile
# base-image/ubuntu20-python38.Dockerfile
FROM nvidia/cuda:11.7.0-cudnn8-devel-ubuntu20.04

# 设置环境变量
ENV DEBIAN_FRONTEND=noninteractive
ENV PYTHON_VERSION=3.8

# 安装基础软件
RUN apt-get update && apt-get install -y \
    python${PYTHON_VERSION} \
    python3-pip \
    openssh-server \
    git \
    vim \
    tmux \
    wget \
    curl \
    && rm -rf /var/lib/apt/lists/*

# 配置 SSH
RUN mkdir /var/run/sshd && \
    sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/' /etc/ssh/sshd_config && \
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config

# 创建用户
RUN useradd -m -s /bin/bash developer && \
    echo "developer ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# 安装 Python 包
RUN pip3 install --no-cache-dir \
    jupyterlab \
    numpy \
    pandas \
    matplotlib \
    scikit-learn

# 创建工作目录
RUN mkdir -p /gemini/code /gemini/output && \
    chown -R developer:developer /gemini

# 启动脚本
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 22 8888

ENTRYPOINT ["/entrypoint.sh"]
```

```dockerfile
# pytorch-image/ubuntu20-pytorch.Dockerfile
FROM registry.example.com/base/ubuntu20-python38:latest

# 安装 PyTorch
RUN pip3 install --no-cache-dir \
    torch==2.0.0 \
    torchvision==0.15.0 \
    torchaudio==2.0.0 \
    --index-url https://download.pytorch.org/whl/cu117

# 安装常用库
RUN pip3 install --no-cache-dir \
    tensorboard \
    transformers \
    accelerate
```

#### 镜像构建流水线

```yaml
# .gitlab-ci.yml
stages:
  - build
  - push

build-base-images:
  stage: build
  script:
    - docker build -t registry.example.com/base/ubuntu20-python38:latest -f base-image/ubuntu20-python38.Dockerfile .
    - docker build -t registry.example.com/base/ubuntu20-python310:latest -f base-image/ubuntu20-python310.Dockerfile .

build-framework-images:
  stage: build
  needs: ["build-base-images"]
  script:
    - docker build -t registry.example.com/pytorch/ubuntu20-pytorch:2.0 -f pytorch-image/ubuntu20-pytorch.Dockerfile .
    - docker build -t registry.example.com/tensorflow/ubuntu20-tf:2.12 -f tensorflow-image/ubuntu20-tf.Dockerfile .

push-images:
  stage: push
  script:
    - docker push registry.example.com/base/ubuntu20-python38:latest
    - docker push registry.example.com/pytorch/ubuntu20-pytorch:2.0
```

### 5.4 自定义镜像构建

#### 数据库设计

```sql
-- 自定义镜像表
CREATE TABLE custom_images (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR(64) UNIQUE NOT NULL,
    customer_id BIGINT NOT NULL,
    name VARCHAR(256) NOT NULL,
    base_image VARCHAR(512),                   -- 基础镜像
    dockerfile TEXT,                           -- Dockerfile 内容
    image_tag VARCHAR(512),                    -- registry.example.com/customer-123/my-image:v1
    size BIGINT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'building',     -- building, ready, failed
    build_log TEXT,
    visibility VARCHAR(20) DEFAULT 'private',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_custom_images_customer_id ON custom_images (customer_id);

-- 构建历史表
CREATE TABLE image_builds (
    id BIGSERIAL PRIMARY KEY,
    image_id BIGINT NOT NULL,
    build_number INT NOT NULL,
    status VARCHAR(20),                        -- building, success, failed
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    build_log TEXT,
    FOREIGN KEY (image_id) REFERENCES custom_images(id) ON DELETE CASCADE
);
```

#### API 实现

```go
// 创建自定义镜像
// POST /api/images/custom
func CreateCustomImage(c *gin.Context) {
    var req struct {
        Name       string `json:"name" binding:"required"`
        BaseImage  string `json:"base_image" binding:"required"`
        Dockerfile string `json:"dockerfile" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    customer := GetCurrentCustomer(c)

    // 1. 验证 Dockerfile（禁止某些指令）
    if err := validateDockerfile(req.Dockerfile); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 2. 生成镜像标签
    imageUUID := uuid.New().String()
    imageTag := fmt.Sprintf("registry.example.com/customer-%d/%s:v1", customer.ID, req.Name)

    // 3. 创建数据库记录
    image := &CustomImage{
        UUID:       imageUUID,
        CustomerID: customer.ID,
        Name:       req.Name,
        BaseImage:  req.BaseImage,
        Dockerfile: req.Dockerfile,
        ImageTag:   imageTag,
        Status:     "building",
    }

    db.Create(image)

    // 4. 异步构建镜像
    go buildImage(image)

    c.JSON(200, gin.H{
        "image_id": imageUUID,
        "status":   "building",
    })
}

// Dockerfile 验证
func validateDockerfile(dockerfile string) error {
    // 禁止的指令
    forbidden := []string{"FROM", "EXPOSE", "CMD", "ENTRYPOINT", "USER"}

    lines := strings.Split(dockerfile, "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        for _, cmd := range forbidden {
            if strings.HasPrefix(strings.ToUpper(line), cmd) {
                return fmt.Errorf("禁止使用 %s 指令", cmd)
            }
        }
    }

    return nil
}

// 构建镜像（使用 Kaniko）
func buildImage(image *CustomImage) {
    // 1. 生成完整的 Dockerfile
    fullDockerfile := fmt.Sprintf("FROM %s\n%s", image.BaseImage, image.Dockerfile)

    // 2. 将 Dockerfile 写入临时目录
    buildDir := filepath.Join("/tmp/builds", image.UUID)
    os.MkdirAll(buildDir, 0755)
    ioutil.WriteFile(filepath.Join(buildDir, "Dockerfile"), []byte(fullDockerfile), 0644)

    // 3. 使用 Kaniko 构建
    cmd := exec.Command(
        "docker", "run",
        "-v", fmt.Sprintf("%s:/workspace", buildDir),
        "gcr.io/kaniko-project/executor:latest",
        "--dockerfile=/workspace/Dockerfile",
        "--destination="+image.ImageTag,
        "--cache=true",
    )

    // 4. 捕获日志
    output, err := cmd.CombinedOutput()
    buildLog := string(output)

    // 5. 更新状态
    if err != nil {
        db.Model(image).Updates(map[string]interface{}{
            "status":    "failed",
            "build_log": buildLog,
        })
    } else {
        db.Model(image).Updates(map[string]interface{}{
            "status":    "ready",
            "build_log": buildLog,
        })
    }

    // 6. 清理
    os.RemoveAll(buildDir)
}
```

#### Dockerfile 编辑器（前端）

```typescript
// ImageEditor.tsx
import React, { useState } from 'react';
import { MonacoEditor } from '@monaco-editor/react';

const dockerfileTemplate = `# 安装 Python 包
RUN pip install --no-cache-dir \\
    torch \\
    transformers \\
    numpy

# 安装系统包
RUN apt-get update && apt-get install -y \\
    ffmpeg \\
    libsm6 \\
    libxext6

# 设置环境变量
ENV MY_VAR=value

# 创建目录
RUN mkdir -p /app
`;

export const ImageEditor: React.FC = () => {
  const [dockerfile, setDockerfile] = useState(dockerfileTemplate);
  const [baseImage, setBaseImage] = useState('ubuntu20-pytorch:2.0');

  const handleBuild = async () => {
    const response = await fetch('/api/images/custom', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: 'my-custom-image',
        base_image: baseImage,
        dockerfile: dockerfile,
      }),
    });

    const data = await response.json();
    console.log('Build started:', data);
  };

  return (
    <div>
      <h2>自定义镜像构建</h2>

      <div>
        <label>基础镜像：</label>
        <select value={baseImage} onChange={e => setBaseImage(e.target.value)}>
          <option value="ubuntu20-pytorch:2.0">Ubuntu 20.04 + PyTorch 2.0</option>
          <option value="ubuntu20-tf:2.12">Ubuntu 20.04 + TensorFlow 2.12</option>
        </select>
      </div>

      <MonacoEditor
        height="400px"
        language="dockerfile"
        value={dockerfile}
        onChange={value => setDockerfile(value || '')}
      />

      <button onClick={handleBuild}>开始构建</button>
    </div>
  );
};
```

---

## 6. 文件系统挂载

### 6.1 挂载方案对比

| 方案 | 实现方式 | 优点 | 缺点 |
|------|---------|------|------|
| **直接挂载** | Volume Mount | 性能好、简单 | 需要预先下载 |
| **FUSE 挂载** | s3fs / goofys | 按需加载 | 性能较差 |
| **JuiceFS** | 分布式文件系统 | 性能好、功能强 | 部署复杂 |
| **符号链接** | Symlink | 灵活 | 管理复杂 |

**推荐：直接挂载（MVP）+ JuiceFS（扩展）**

### 6.2 实现方案

#### 方案 A：预下载 + 直接挂载（推荐）

```go
// 容器启动前，将数据集下载到本地
func PrepareDatasets(envID string, datasets []Dataset) error {
    for i, dataset := range datasets {
        // 1. 创建挂载目录
        mountPath := filepath.Join("/data/environments", envID, fmt.Sprintf("data-%d", i+1))
        os.MkdirAll(mountPath, 0755)

        // 2. 从对象存储同步数据
        cmd := exec.Command(
            "rclone", "sync",
            fmt.Sprintf("minio:%s/%s", bucketName, dataset.StoragePath),
            mountPath,
        )
        if err := cmd.Run(); err != nil {
            return err
        }
    }

    return nil
}

// Docker 启动时挂载
docker run \
  -v /data/environments/env-123/data-1:/gemini/data-1:ro \
  -v /data/environments/env-123/data-2:/gemini/data-2:ro \
  ...
```

#### 方案 B：FUSE 挂载（按需加载）

```bash
# 使用 s3fs 挂载 MinIO
s3fs remotegpu /mnt/s3 \
  -o url=http://minio:9000 \
  -o use_path_request_style \
  -o passwd_file=/etc/passwd-s3fs

# 容器挂载
docker run \
  -v /mnt/s3/datasets/customer-123/dataset-001:/gemini/data-1:ro \
  ...
```

#### 方案 C：JuiceFS（生产推荐）

```bash
# 1. 部署 JuiceFS
juicefs format \
  --storage minio \
  --bucket http://minio:9000/remotegpu \
  --access-key admin \
  --secret-key password \
  redis://redis:6379/0 \
  remotegpu-fs

# 2. 挂载文件系统
juicefs mount redis://redis:6379/0 /mnt/jfs

# 3. 容器使用
docker run \
  -v /mnt/jfs/datasets/customer-123/dataset-001:/gemini/data-1:ro \
  ...
```

### 6.3 Kubernetes PV/PVC 方案

```yaml
# 使用 CSI 驱动挂载对象存储
apiVersion: v1
kind: PersistentVolume
metadata:
  name: dataset-pv-123
spec:
  capacity:
    storage: 100Gi
  accessModes:
    - ReadOnlyMany
  csi:
    driver: csi-s3
    volumeHandle: dataset-123
    volumeAttributes:
      bucket: remotegpu
      prefix: datasets/customer-123/dataset-001/

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: dataset-pvc-123
spec:
  accessModes:
    - ReadOnlyMany
  resources:
    requests:
      storage: 100Gi
  volumeName: dataset-pv-123

---
# Pod 使用
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: dev-env
    volumeMounts:
    - name: dataset-1
      mountPath: /gemini/data-1
      readOnly: true
  volumes:
  - name: dataset-1
    persistentVolumeClaim:
      claimName: dataset-pvc-123
```

---

## 7. 传统架构 vs K8s 实现

### 7.1 传统架构实现

```go
// 创建开发环境时准备数据
func CreateEnvironmentWithDatasets(req CreateEnvRequest) error {
    // 1. 选择主机
    host := scheduler.SelectHost(req.Resources)

    // 2. 准备数据集（下载到主机本地）
    for i, datasetID := range req.Datasets {
        dataset := getDataset(datasetID)
        localPath := fmt.Sprintf("/data/environments/%s/data-%d", req.EnvID, i+1)

        // 调用主机 Agent 下载数据
        agent.DownloadDataset(host, dataset.StoragePath, localPath)
    }

    // 3. 创建容器（挂载本地目录）
    volumes := []string{
        fmt.Sprintf("/data/environments/%s/code:/gemini/code", req.EnvID),
    }
    for i := range req.Datasets {
        volumes = append(volumes, fmt.Sprintf(
            "/data/environments/%s/data-%d:/gemini/data-%d:ro",
            req.EnvID, i+1, i+1,
        ))
    }

    agent.CreateContainer(host, CreateContainerRequest{
        EnvID:   req.EnvID,
        Image:   req.Image,
        Volumes: volumes,
    })

    return nil
}
```

### 7.2 Kubernetes 实现

```go
// 创建开发环境时创建 PVC
func CreateEnvironmentWithDatasets(req CreateEnvRequest) error {
    // 1. 为每个数据集创建 PVC
    for i, datasetID := range req.Datasets {
        dataset := getDataset(datasetID)

        pvc := &corev1.PersistentVolumeClaim{
            ObjectMeta: metav1.ObjectMeta{
                Name: fmt.Sprintf("dataset-pvc-%s-%d", req.EnvID, i+1),
            },
            Spec: corev1.PersistentVolumeClaimSpec{
                AccessModes: []corev1.PersistentVolumeAccessMode{
                    corev1.ReadOnlyMany,
                },
                Resources: corev1.ResourceRequirements{
                    Requests: corev1.ResourceList{
                        "storage": resource.MustParse("100Gi"),
                    },
                },
                VolumeName: fmt.Sprintf("dataset-pv-%s", dataset.UUID),
            },
        }

        k8sClient.CoreV1().PersistentVolumeClaims("dev-environments").Create(ctx, pvc, metav1.CreateOptions{})
    }

    // 2. 创建 Pod（引用 PVC）
    volumeMounts := []corev1.VolumeMount{
        {Name: "code", MountPath: "/gemini/code"},
    }
    volumes := []corev1.Volume{
        {
            Name: "code",
            VolumeSource: corev1.VolumeSource{
                PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
                    ClaimName: fmt.Sprintf("code-pvc-%s", req.EnvID),
                },
            },
        },
    }

    for i := range req.Datasets {
        volumeMounts = append(volumeMounts, corev1.VolumeMount{
            Name:      fmt.Sprintf("dataset-%d", i+1),
            MountPath: fmt.Sprintf("/gemini/data-%d", i+1),
            ReadOnly:  true,
        })
        volumes = append(volumes, corev1.Volume{
            Name: fmt.Sprintf("dataset-%d", i+1),
            VolumeSource: corev1.VolumeSource{
                PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
                    ClaimName: fmt.Sprintf("dataset-pvc-%s-%d", req.EnvID, i+1),
                },
            },
        })
    }

    pod := buildPodSpec(req, volumeMounts, volumes)
    k8sClient.CoreV1().Pods("dev-environments").Create(ctx, pod, metav1.CreateOptions{})

    return nil
}
```

---

## 8. 最佳实践

### 8.1 性能优化

1. **数据集缓存**
   - 对于热门数据集，在每台主机上缓存
   - 使用 LRU 策略清理

2. **增量同步**
   - 使用 rsync/rclone 增量同步
   - 避免重复下载

3. **压缩存储**
   - 对大数据集使用压缩
   - 容器内按需解压

4. **并行上传**
   - 使用分片并行上传
   - 提高上传速度

### 8.2 成本优化

1. **存储分层**
   - 热数据：高性能存储
   - 冷数据：归档存储
   - 定期清理过期数据

2. **重复数据删除**
   - 使用哈希检测重复文件
   - 符号链接替代重复存储

3. **镜像层共享**
   - 基础镜像统一管理
   - 充分利用 Docker 层缓存

---

## 9. 总结

### 推荐技术栈

**存储层：**
- MinIO（对象存储）
- Harbor（镜像仓库）
- JuiceFS（文件系统，可选）

**数据管理：**
- 预签名 URL（前端直传）
- 分片上传（大文件）
- 版本管理（数据集/模型版本）

**镜像管理：**
- 官方镜像库（预构建）
- Dockerfile 自定义（Kaniko 构建）
- 镜像缓存（加速拉取）

---

**文档结束**

这套方案可以支持从 MVP 到大规模生产的完整演进路径。
