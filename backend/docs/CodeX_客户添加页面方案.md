# CodeX 客户添加页面方案

- 作者: CodeX
- 日期: 2026-02-05

## 目标

完善管理端“添加客户”页面，字段与后端创建接口一致，并确保创建后在客户列表可看到关键信息（公司/联系人/邮箱/电话/状态）。

## 表单字段

必填：
- username
- company
- email
- phone
- password（默认 ChangeME_123）

可选：
- display_name
- full_name

## 接口对齐

- 前端调用：`POST /admin/customers`
- 后端请求体新增可选字段：display_name/full_name/company/phone，password 为空时默认 ChangeME_123
- 后端创建时写入对应 Customer 字段

## 路由与页面

- 新增页面：`/admin/customers/add`
- 客户列表页增加“添加客户”按钮

## 校验规则

- username 必填，长度 >= 3
- email 必填，邮箱格式
- password 必填，长度 >= 6
- phone 非必填（仅格式校验）
