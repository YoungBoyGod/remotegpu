# 文档合理性审查记录

> 审查日期：2026-01-26
> 覆盖范围：docs 目录内本次指定的设计/对比/集成/审查文档

---

## 发现的问题

1. `docs/design/system_architecture.md` 模块列表仍为 12 个，缺少通知/ Webhook 等模块（与修复总结中“15 个核心模块”的表述不一致）。
   - 位置：`docs/design/system_architecture.md:112`
   - 关联冲突：`docs/reviews/fix_summary.md:24`

2. Windows Server 支持 Docker Desktop 的表述不准确。Docker Desktop 不支持 Windows Server，生产应使用 Windows Server 的容器运行时（如 Mirantis Container Runtime）。
   - 位置：`docs/design/system_architecture.md:511`

3. 数据库为 PostgreSQL，但数据集表的 SQL 使用了 MySQL 风格的 `INDEX` 语法，PostgreSQL 不支持在 `CREATE TABLE` 中这样声明索引。
   - 位置：`docs/design/storage_and_image_management.md:311`

4. 创建数据集 API 的请求结构体缺少 `WorkspaceID` 字段，但后续写库时使用了 `req.WorkspaceID`，会导致必填字段为空或默认值。
   - 位置：`docs/design/storage_and_image_management.md:381`

5. 预签名上传 URL 示例把 `*gin.Context` 直接作为 `context.Context` 使用，类型不匹配（应使用 `c.Request.Context()` 或 `context.Background()`）。
   - 位置：`docs/design/storage_and_image_management.md:452`

6. “负载均衡器转发（推荐）”方案需要为每个端口写死 `stream` 配置，难以支撑动态端口数量和弹性扩缩容，作为推荐方案不合理。
   - 位置：`docs/comparison/ssh_deployment_comparison.md:386`

7. iptables 删除规则示例与添加规则不匹配，缺少 `--to-destination`，实际会删除失败或匹配不到。
   - 位置：`docs/comparison/ssh_deployment_comparison.md:434`

8. GitLab CI 示例把 `curl` 参数拆成多个 `script` 行，`-H` 和 `-d` 会被当作独立命令执行，示例不可用。
   - 位置：`docs/infrastructure/third_party_integration.md:306`

9. 修复总结中的“待修复的问题”小节为空，未记录剩余问题，导致修复状态不可追踪。
   - 位置：`docs/reviews/fix_summary.md:30`
