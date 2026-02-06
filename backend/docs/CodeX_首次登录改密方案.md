# CodeX 首次登录强制改密方案

- 作者: CodeX
- 日期: 2026-02-06

## 目标

客户使用默认密码 `ChangeME_123` 首次登录后必须修改密码，防止默认口令长期存在。

## 方案摘要

- 新增字段：`customers.must_change_password`（默认 false）。
- 创建客户时：若使用默认密码（空密码会自动填充 `ChangeME_123`），则写入 `must_change_password=true`。
- 登录/刷新响应：返回 `must_change_password` 标记。
- 个人资料接口：包含 `must_change_password` 字段。
- 前端：检测 `must_change_password=true` 时强制跳转到改密页面。
- 改密接口：`POST /api/v1/auth/password/change`，修改成功后清除 `must_change_password`。

## 必要接口

- 登录响应新增字段：`must_change_password`。
- 刷新响应新增字段：`must_change_password`。
- 改密请求：
  - 请求体：`old_password`, `new_password`。
  - 成功后返回 `message=ok`。

## 前端规则

- 登录成功后若 `must_change_password=true`，仅允许访问改密页与登出。
- 改密成功后跳转到角色首页。

## 实施清单

- 数据库：新增字段 `customers.must_change_password`（默认 false）。
- 后端：
  - 登录/刷新响应带回 `must_change_password`。
  - 修改密码接口校验旧密码并写入新 hash。
  - 修改成功后清除 `must_change_password`。
- 前端：
  - 登录成功后读取 profile 的 `must_change_password`。
  - 路由守卫拦截并强制跳转到改密页。
  - 改密页提交成功后跳转到角色首页。
