# 📋 实施计划：Go 模块路径修正 + 文档补充代码更新操作

## 任务类型
- [x] 后端（go.mod + import 路径重命名）
- [x] 文档（deployment.md 补充更新章节）

## 背景

- 当前 module：`github.com/yourblog/backend`（占位符）
- 实际 Git 仓库：`https://github.com/yavon007/blog-dev`
- 目标 module：`github.com/yavon007/blog-dev/backend`
- 受影响 Go 文件：17 个

---

## 实施步骤

### 步骤 1：修改 go.mod 模块名

**文件**：`backend/go.mod:1`

```diff
-module github.com/yourblog/backend
+module github.com/yavon007/blog-dev/backend
```

### 步骤 2：全量替换 import 路径

对以下 17 个文件中所有 `github.com/yourblog/backend` 替换为 `github.com/yavon007/blog-dev/backend`：

| 文件 | 说明 |
|------|------|
| `backend/cmd/api/main.go` | 入口 |
| `backend/internal/app/router.go` | 路由 |
| `backend/internal/platform/database/postgres.go` | 数据库 |
| `backend/internal/platform/cache/redis.go` | 缓存 |
| `backend/internal/platform/auth/jwt_manager.go` | JWT |
| `backend/internal/modules/auth/core/service.go` | 认证服务 |
| `backend/internal/modules/auth/repository/postgres.go` | 认证仓储 |
| `backend/internal/modules/auth/transport/http/handler.go` | 认证 handler |
| `backend/internal/modules/posts/core/service.go` | 文章服务 |
| `backend/internal/modules/posts/repository/postgres.go` | 文章仓储 |
| `backend/internal/modules/posts/transport/http/handler.go` | 文章 handler |
| `backend/internal/modules/taxonomy/core/service.go` | 分类服务 |
| `backend/internal/modules/taxonomy/repository/postgres.go` | 分类仓储 |
| `backend/internal/modules/taxonomy/transport/http/handler.go` | 分类 handler |
| `backend/internal/modules/comments/core/service.go` | 评论服务 |
| `backend/internal/modules/comments/repository/postgres.go` | 评论仓储 |
| `backend/internal/modules/comments/transport/http/handler.go` | 评论 handler |
| `backend/internal/pkg/middleware/jwt.go` | JWT 中间件 |

### 步骤 3：验证编译

```bash
cd backend && go build ./...
```

### 步骤 4：文档补充 - 代码更新章节

**文件**：`docs/deployment.md`

在"备份与恢复"章节前（或文末 FAQ 后）新增 `## 代码更新与重新部署` 章节，内容包括：

1. **拉取最新代码**
   ```bash
   cd /opt/blog
   git pull origin main
   ```

2. **重新构建并热更新服务**
   ```bash
   make prod
   # 等价于：docker compose up -d --build
   ```

3. **仅更新后端**（前端无变更时）
   ```bash
   docker compose up -d --build backend
   ```

4. **执行数据库迁移**（如有新迁移文件）
   ```bash
   make migrate-up
   # 等价于：docker compose run --rm migrate
   ```

5. **验证更新结果**
   ```bash
   make ps
   curl http://localhost:8080/api/v1/health
   ```

6. **回滚**（如更新失败）
   ```bash
   git checkout <上一个 commit hash>
   make prod
   ```

---

## 关键文件

| 文件 | 操作 | 说明 |
|------|------|------|
| `backend/go.mod:1` | 修改 | 更新 module 声明 |
| `backend/**/*.go`（17 个） | 修改 | 替换 import 路径 |
| `docs/deployment.md` | 修改 | 新增代码更新章节 |

## 风险与缓解

| 风险 | 缓解措施 |
|------|----------|
| import 替换遗漏导致编译失败 | 替换后运行 `go build ./...` 验证 |
| 文档更新位置不合理 | 插入在"监控与日志"章节之后，"常见问题"之前 |

## SESSION_ID
- CODEX_SESSION: N/A（本任务由 Claude 直接执行，无需多模型会话）
- GEMINI_SESSION: N/A
