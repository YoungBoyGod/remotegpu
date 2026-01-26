# Docs 目录重组方案

## 新的目录结构

```
docs/
├── requirements/              # 需求文档（已存在）
│   ├── README.md
│   └── 01-23 需求文档
│
├── design/                    # 设计文档
│   ├── system_architecture.md      # 系统架构设计
│   ├── module_division.md          # 模块划分
│   ├── database_design.md          # 数据库设计
│   ├── cmdb_design.md              # CMDB 设计
│   ├── ssh_access_design.md        # SSH 访问设计
│   ├── storage_and_image_management.md  # 存储和镜像管理设计
│   ├── windows_remote_access.md    # Windows 远程访问设计
│   └── customer_management.md      # 客户管理设计
│
├── infrastructure/            # 基础设施文档
│   ├── infrastructure.md           # 基础设施配置
│   └── third_party_integration.md  # 第三方集成
│
├── comparison/                # 对比分析文档
│   ├── virtai_cloud_analysis.md    # VirtAI Cloud 分析
│   └── ssh_deployment_comparison.md # SSH 部署方案对比
│
├── reviews/                   # 审查和修复报告
│   ├── review_report.md            # 文档审查报告
│   ├── fix_completed.md            # 修复完成报告
│   └── fix_summary.md              # 修复总结
│
└── REQUIREMENTS.md            # 总需求文档（保留在根目录）
```

## 文件分类说明

### 1. design/ - 设计文档
包含系统架构、模块划分、数据库设计等技术设计文档

### 2. infrastructure/ - 基础设施文档
包含基础设施配置、第三方系统集成等运维相关文档

### 3. comparison/ - 对比分析文档
包含竞品分析、方案对比等调研文档

### 4. reviews/ - 审查报告
包含文档审查、修复报告等质量管理文档

