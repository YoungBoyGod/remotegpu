---
paths:
  - "frontend/src/**/*.{ts,vue}"
---

# 前端代码风格

## Vue 组件

- 使用 `<script setup lang="ts">` 组合式 API
- UI 框架：Element Plus（el-* 组件）
- 组件顺序：`<script>` → `<template>` → `<style scoped>`

## API 请求

- 请求函数放在 `src/api/` 下，按模块分文件
- 返回类型统一为 `Promise<ApiResponse<T>>`
- 使用 `src/api/` 下的 `request` 实例（已封装 baseURL 和拦截器）

## 类型定义

- 放在 `src/types/` 下，按模块分文件
- 字段命名使用 snake_case（与后端 JSON 一致）
- 使用 `interface` 而非 `type`

## 页面风格

- 页面容器 padding: 24px
- 页面标题用 `.page-header` + `.page-title`（font-size: 24px）
- 筛选栏用 `el-card`，数据表格用 `DataTable` 组件
- 消息提示用 `ElMessage`，确认框用 `ElMessageBox`

## 路由

- 管理员路由前缀 `/admin`，客户路由前缀 `/customer`
- name 格式：`角色-模块-操作`（如 `admin-settings-platform`）
- 新路由必须放在通配路由 `:pathMatch(.*)*` 之前
