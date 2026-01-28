# 员工A - 任务分配

## 基本信息
- **负责人**: 员工A
- **主要模块**: Epic 1 - Feature 1.1 工作空间管理
- **预估工时**: 12小时
- **优先级**: P0（最高）
- **开始时间**: 第1周
- **依赖**: Customer模型已完成（已完成）

---

## 任务概述

你负责实现**工作空间管理**功能，这是用户协作的核心模块。包括：
1. Workspace（工作空间）数据模型和DAO
2. WorkspaceMember（工作空间成员）数据模型和DAO

---

## 详细任务列表

### Story 1.1.1: 创建Workspace数据模型（4小时）

#### Task 1.1.1.1: 创建Workspace实体模型（1小时）

**文件路径**: `internal/model/entity/workspace.go`

**Subtask清单**:
- [ ] 1.1.1.1.1: 创建文件和包声明（5分钟）
  - 在`internal/model/entity/`目录创建`workspace.go`
  - 添加package声明和必要的import

- [ ] 1.1.1.1.2: 定义Workspace结构体（15分钟）
  - 参考SQL文件`sql/03_users_and_permissions.sql`
  - 定义所有字段：ID, UUID, OwnerID, Name, Description, Type, MemberCount, Status, CreatedAt, UpdatedAt

- [ ] 1.1.1.1.3: 添加GORM标签（20分钟）
  - 主键字段添加`gorm:"primaryKey"`
  - UUID字段添加`gorm:"type:uuid;unique;not null"`
  - 外键字段添加`gorm:"not null"`
  - 字符串字段添加长度限制

- [ ] 1.1.1.1.4: 实现TableName方法（5分钟）
  - 实现`TableName()`方法返回"workspaces"

- [ ] 1.1.1.1.5: 添加字段注释（10分钟）
  - 为结构体和字段添加Go文档注释

- [ ] 1.1.1.1.6: 代码格式化和检查（5分钟）
  - 运行`go fmt`和`go vet`

**验收标准**:
- [ ] 文件编译无错误
- [ ] 结构体字段与SQL表结构一致
- [ ] GORM标签正确配置
- [ ] 所有导出字段有注释

---

#### Task 1.1.1.2: 创建WorkspaceDao基础结构（25分钟）

**文件路径**: `internal/dao/workspace.go`

**Subtask清单**:
- [ ] 1.1.1.2.1: 创建文件和包声明（5分钟）
- [ ] 1.1.1.2.2: 定义WorkspaceDao结构体（5分钟）
- [ ] 1.1.1.2.3: 实现构造函数NewWorkspaceDao（10分钟）
- [ ] 1.1.1.2.4: 添加结构体注释（5分钟）

**验收标准**:
- [ ] 构造函数正确获取数据库连接
- [ ] 代码符合Go规范

---

#### Task 1.1.1.3: 实现WorkspaceDao的Create方法（45分钟）

**Subtask清单**:
- [ ] 1.1.1.3.1: 定义方法签名（5分钟）
  - `Create(workspace *entity.Workspace) error`

- [ ] 1.1.1.3.2: 实现创建逻辑（15分钟）
  - 使用`db.Create()`创建记录
  - 处理错误返回

- [ ] 1.1.1.3.3: 添加方法注释（5分钟）

- [ ] 1.1.1.3.4: 单元测试（20分钟）
  - 创建`internal/dao/workspace_test.go`
  - 编写TestCreate测试用例

**验收标准**:
- [ ] 方法实现正确
- [ ] 单元测试通过
- [ ] 错误处理完善

---

#### Task 1.1.1.4: 实现WorkspaceDao的查询方法（1小时20分钟）

**Subtask清单**:
- [ ] 1.1.1.4.1: 实现GetByID方法（20分钟）
  - 使用`db.First()`查询
  - 处理记录不存在的情况

- [ ] 1.1.1.4.2: 实现GetByUUID方法（15分钟）
  - 使用`db.Where().First()`查询

- [ ] 1.1.1.4.3: 实现GetByOwnerID方法（15分钟）
  - 使用`db.Where().Find()`查询列表

- [ ] 1.1.1.4.4: 单元测试（30分钟）
  - 编写查询方法的测试用例

**验收标准**:
- [ ] 所有查询方法实现正确
- [ ] 单元测试通过
- [ ] 正确处理记录不存在的情况

---

#### Task 1.1.1.5: 实现WorkspaceDao的更新和删除方法（50分钟）

**Subtask清单**:
- [ ] 1.1.1.5.1: 实现Update方法（15分钟）
  - 使用`db.Save()`更新记录

- [ ] 1.1.1.5.2: 实现Delete方法（15分钟）
  - 使用`db.Delete()`删除记录

- [ ] 1.1.1.5.3: 单元测试（20分钟）
  - 编写更新删除的测试用例

**验收标准**:
- [ ] 更新和删除方法实现正确
- [ ] 单元测试通过

---

### Story 1.1.2: 创建WorkspaceMember数据模型（4小时）

#### Task 1.1.2.1: 创建WorkspaceMember实体模型（45分钟）

**文件路径**: `internal/model/entity/workspace_member.go`

**Subtask清单**:
- [ ] 1.1.2.1.1: 创建文件和结构体（20分钟）
  - 定义所有字段：ID, WorkspaceID, CustomerID, Role, Status, JoinedAt, CreatedAt

- [ ] 1.1.2.1.2: 添加GORM标签和唯一索引（20分钟）
  - 添加唯一索引`gorm:"uniqueIndex:idx_workspace_customer"`

- [ ] 1.1.2.1.3: 实现TableName方法（5分钟）

**验收标准**:
- [ ] 唯一索引正确配置
- [ ] 字段类型正确

---

#### Task 1.1.2.2: 创建WorkspaceMemberDao（3小时15分钟）

**文件路径**: `internal/dao/workspace_member.go`

**Subtask清单**:
- [ ] 1.1.2.2.1: 创建基础结构和构造函数（20分钟）

- [ ] 1.1.2.2.2: 实现Create方法（30分钟）
  - 处理唯一约束冲突

- [ ] 1.1.2.2.3: 实现查询方法（1小时）
  - GetByWorkspaceID
  - GetByCustomerID
  - GetByWorkspaceAndCustomer
  - GetByRole

- [ ] 1.1.2.2.4: 实现Update和Delete方法（30分钟）

- [ ] 1.1.2.2.5: 编写单元测试（1小时）
  - 测试唯一约束是否生效

**验收标准**:
- [ ] 所有方法实现正确
- [ ] 唯一约束测试通过
- [ ] 单元测试覆盖率>80%

---

## 协作和依赖

### 依赖其他员工的工作
- **无直接依赖** - 你的工作只依赖已完成的Customer模型

### 其他员工依赖你的工作
- **员工B**: 需要你完成Workspace模型后才能创建ResourceQuota（工作空间配额）
- **员工E**: 需要你完成Workspace模型后才能创建Environment（环境需要关联工作空间）

### 协作点
- 与**员工B**协作：完成Workspace模型后立即通知，以便其开始ResourceQuota开发
- 与**员工E**协作：WorkspaceID字段的外键关系需要保持一致

---

## 验收检查清单

### 代码质量
- [ ] 所有代码通过`go fmt`格式化
- [ ] 所有代码通过`go vet`检查
- [ ] 所有导出的结构体和方法有注释
- [ ] 错误处理完善

### 测试
- [ ] 所有单元测试通过：`go test ./internal/dao -v`
- [ ] 测试覆盖率>80%
- [ ] 测试数据清理完成

### 数据库
- [ ] 表结构与SQL文件一致
- [ ] GORM标签正确配置
- [ ] 唯一索引生效

### Git提交
- [ ] 提交信息清晰：`feat: 实现Workspace和WorkspaceMember模型及DAO`
- [ ] 代码已推送到远程仓库

---

## 注意事项

1. **UUID字段**: 需要导入`github.com/google/uuid`包
2. **唯一索引**: WorkspaceMember的(workspace_id, customer_id)必须唯一
3. **外键关系**: OwnerID关联customers表的id字段
4. **测试数据**: 测试后必须清理，避免污染数据库
5. **错误处理**: 使用`errors.Is(err, gorm.ErrRecordNotFound)`判断记录不存在

---

## 进度报告

请在完成每个Task后更新`团队协作文档.md`中的进度表，格式如下：

```
| 员工A | Task 1.1.1.1 | ✅ 已完成 | 2024-01-28 10:00 |
```

---

## 联系方式

遇到问题请及时沟通：
- 技术问题：查看`实施步骤指南.md`
- 依赖问题：联系员工B和员工E
- 其他问题：联系项目经理
