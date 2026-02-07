# 前端开发 2 (frontend-2) 任务清单

## 职责
- 客户端（Customer）页面开发
- 公共组件与布局
- 登录与认证流程

## 任务列表

### T-F2-01 [P0] 登录页面与认证流程
- 登录表单（用户名/密码）
- Token 存储与刷新
- 登出与 Token 清理
- 首次登录强制改密
- 位置: frontend/src/views/LoginView.vue

### T-F2-02 [P1] 客户端机器列表
- 展示用户已分配的机器
- 连接信息展示（SSH/Jupyter/VNC）
- 位置: frontend/src/views/customer/MachineListView.vue

### T-F2-03 [P1] 任务管理页面
- 任务列表与状态展示
- 创建训练任务表单
- 停止任务操作
- 位置: frontend/src/views/customer/TaskListView.vue

### T-F2-04 [P1] 数据集管理页面
- 数据集列表
- 分片上传功能
- 数据集挂载操作
- 位置: frontend/src/views/customer/DatasetListView.vue

### T-F2-05 [P2] SSH 密钥管理页面
- 公钥列表与添加
- 密钥删除
- 位置: frontend/src/views/customer/SSHKeyView.vue
