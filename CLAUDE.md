# CLAUDE.md — AI 开发规范

本文档为 AI 助手（Claude Code 等）在此项目中工作的规范约束。

## 项目概述

前后端分离博客系统：
- **前端**：`frontend/` — Vue 3 + TypeScript + Vite + UnoCSS
- **后端**：`backend/` — Go + Gin + PostgreSQL + Redis

## 架构约束

### 全局规则

1. **外部服务写入**：不得直接操作数据库或 Redis，必须通过后端 API
2. **环境变量**：敏感信息（密钥、密码）只写入 `.env`，不进代码
3. **提交规范**：使用 Conventional Commits（`feat:`, `fix:`, `docs:`, `refactor:` 等）

### 后端约束（backend/）

**目录规范**：
- 领域逻辑 → `internal/modules/<module>/core/`
- 数据库操作 → `internal/modules/<module>/repository/`
- HTTP 处理 → `internal/modules/<module>/transport/http/`
- 平台基础设施 → `internal/platform/`
- 跨模块共享 → `internal/modules/shared/`

**代码规范**：
- 使用 `pgx/v5`，禁止使用 `database/sql` 原生接口
- 错误处理：所有 error 必须 wrap 后返回（`fmt.Errorf("...: %w", err)`）
- 日志：使用 `zap`，禁止使用 `fmt.Println` 或标准 `log`
- 验证：使用 `go-playground/validator/v10`，DTO 必须打标签
- 测试：新增 repository 和 service 必须附带单元测试

**API 规范**：
- 版本前缀：`/api/v1/`
- 公开 API：`/api/v1/` 下无需鉴权
- 管理 API：`/api/v1/admin/` 下需要 JWT Bearer Token
- 统一响应格式：`{ "code": 200, "data": {}, "message": "success" }`
- 错误响应格式：`{ "code": 4xx/5xx, "message": "error description" }`

**禁止事项**：
- 禁止在 handler 中直接写 SQL
- 禁止跨模块直接调用 repository
- 禁止在 transport 层出现业务逻辑

### 前端约束（frontend/）

**TypeScript 规范**：
- 严格模式：`tsconfig.json` 中 `strict: true`
- 所有 `.vue` 文件使用 `<script setup lang="ts">`
- 禁止使用 `any`，未知类型用 `unknown`
- API 响应必须定义 TypeScript 接口

**目录规范**：
- 前台页面 → `src/views/front/`
- 后台页面 → `src/views/admin/`
- 前台组件 → `src/components/front/`
- 后台组件 → `src/components/admin/`
- 通用组件 → `src/components/common/`
- API 调用 → `src/api/`（调用 `src/utils/request.ts`）
- 状态管理 → `src/store/`（Pinia，Setup 语法）

**样式规范**：
- 样式使用 UnoCSS 原子类，禁止在 `<style>` 中写大量 CSS
- 主题色通过 `uno.config.ts` 的 `theme.colors.primary` 定义
- 文章内容渲染使用 `prose` class（UnoCSS Typography preset）

**性能规范**：
- `md-editor-v3` 只能在 `/admin/` 路由下使用，且必须 `defineAsyncComponent` 懒加载
- 所有路由组件使用 `() => import(...)` 懒加载
- 禁止在前台（FrontLayout）中引入后台专用依赖

**安全规范**：
- 渲染 Markdown HTML 前必须使用 DOMPurify 净化，防止 XSS
- 禁止将 JWT Token 存入 `<script>` 标签或 URL 参数

## 数据库迁移规范

- 迁移文件在 `backend/migrations/` 目录
- 命名格式：`NNNN_description.up.sql` / `NNNN_description.down.sql`
- 每个 `.up.sql` 必须有对应的 `.down.sql`（回滚操作）
- 禁止修改已执行的迁移文件，需要变更则新增迁移

## 常用命令

```bash
# 后端
cd backend && make run          # 启动后端
cd backend && make test         # 运行测试
cd backend && make migrate-up   # 执行迁移
cd backend && make lint         # 代码检查

# 前端
cd frontend && pnpm dev         # 启动前端
cd frontend && pnpm build       # 构建生产包
cd frontend && pnpm type-check  # TypeScript 类型检查

# Docker
docker compose up -d            # 启动所有服务
docker compose down             # 停止所有服务
docker compose logs -f backend  # 查看后端日志
```
