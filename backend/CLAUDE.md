[根目录](../CLAUDE.md) > **backend**

# Backend 模块文档

后端服务：Go + Gin + PostgreSQL + Redis

## 模块职责

提供博客系统的 RESTful API，包括：
- 文章管理（CRUD、发布/草稿状态）
- 分类与标签管理
- 评论管理（审核流程）
- 管理员认证（JWT）

## 入口与启动

**入口文件**: `cmd/api/main.go`

```bash
# 开发环境
cd backend && make run

# 或直接运行
go run ./cmd/api/main.go
```

**启动流程**:
1. 加载配置 (`config.Load()`)
2. 初始化日志 (`logger.New()`)
3. 连接数据库 (`database.NewPool()`)
4. 初始化 JWT Manager
5. 注入各模块依赖 (Repository -> Service -> Handler)
6. 注册路由 (`app.NewRouter()`)
7. 启动 HTTP 服务器，支持优雅关闭

## 对外接口

### 公开 API (`/api/v1/`)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/health` | 健康检查 |
| GET | `/api/v1/posts` | 文章列表（支持分页、分类/标签过滤、搜索） |
| GET | `/api/v1/posts/:slug` | 文章详情 |
| GET | `/api/v1/categories` | 分类列表 |
| GET | `/api/v1/tags` | 标签列表 |
| GET | `/api/v1/posts/:slug/comments` | 文章评论 |
| POST | `/api/v1/posts/:slug/comments` | 创建评论 |
| POST | `/api/v1/auth/login` | 管理员登录 |
| POST | `/api/v1/auth/refresh` | 刷新 Token |

### 管理 API (`/api/v1/admin/`)

需要 JWT Bearer Token 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/admin/auth/logout` | 登出 |
| GET/POST | `/api/v1/admin/posts` | 文章管理（列表/创建） |
| GET/PUT/DELETE | `/api/v1/admin/posts/:id` | 文章管理（详情/更新/删除） |
| PATCH | `/api/v1/admin/posts/:id/status` | 更新文章状态 |
| CRUD | `/api/v1/admin/categories` | 分类管理 |
| CRUD | `/api/v1/admin/tags` | 标签管理 |
| GET/PATCH/DELETE | `/api/v1/admin/comments` | 评论管理 |

### 统一响应格式

```json
// 成功
{"code": 200, "data": {...}, "message": "success"}

// 分页
{"code": 200, "data": {"items": [...], "total": 100, "page": 1, "size": 10}, "message": "success"}

// 错误
{"code": 4xx, "message": "error description"}
```

## 关键依赖与配置

### 主要依赖

- **Web 框架**: `github.com/gin-gonic/gin`
- **数据库驱动**: `github.com/jackc/pgx/v5`
- **缓存**: `github.com/redis/go-redis/v9`
- **JWT**: `github.com/golang-jwt/jwt/v5`
- **日志**: `go.uber.org/zap`
- **验证**: `github.com/go-playground/validator/v10`

### 环境变量配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `APP_ENV` | 运行环境 | `development` |
| `APP_PORT` | 服务端口 | `8080` |
| `DATABASE_URL` | PostgreSQL 连接串 | 必填 |
| `REDIS_ADDR` | Redis 地址 | `localhost:6379` |
| `REDIS_PASSWORD` | Redis 密码 | 空 |
| `JWT_SECRET` | JWT 密钥 | 必填 |
| `JWT_ACCESS_TTL` | Access Token 有效期 | `15m` |
| `JWT_REFRESH_TTL` | Refresh Token 有效期 | `720h` |
| `PUBLIC_ALLOWED_ORIGINS` | CORS 允许的源 | `http://localhost:5173` |
| `LOG_LEVEL` | 日志级别 | `info` |

## 数据模型

### 核心表结构

- `admin_users` - 管理员账户
- `categories` - 分类
- `tags` - 标签
- `posts` - 文章
- `post_tags` - 文章-标签关联
- `comments` - 评论
- `refresh_tokens` - JWT 刷新令牌
- `audit_logs` - 审计日志

### 迁移文件

位于 `migrations/` 目录，命名格式：`NNNN_description.up.sql` / `NNNN_description.down.sql`

```bash
# 执行迁移
make migrate-up

# 回滚
make migrate-down
```

## 测试与质量

### 当前状态

- **单元测试**: 暂无 `*_test.go` 文件
- **代码检查**: 支持 `golangci-lint`

```bash
cd backend && make lint   # 代码检查
cd backend && make test   # 运行测试（需要添加测试文件）
```

### 质量工具建议

- 添加 `internal/modules/*/repository/*_test.go`
- 添加 `internal/modules/*/core/*_test.go`

## 常见问题 (FAQ)

**Q: 如何创建管理员账户？**

A: 目前需要手动插入数据库或通过迁移脚本创建。密码需使用 bcrypt 哈希。

**Q: 如何调试数据库连接问题？**

A: 检查 `DATABASE_URL` 格式：`postgres://user:password@host:port/dbname?sslmode=disable`

**Q: Redis 连接失败怎么办？**

A: 确认 Redis 服务运行中，检查 `REDIS_ADDR` 和 `REDIS_PASSWORD` 配置。

## 相关文件清单

```
backend/
├── cmd/api/main.go           # 入口
├── internal/
│   ├── app/router.go         # 路由注册
│   ├── config/loader.go      # 配置加载
│   ├── pkg/
│   │   ├── response/         # 统一响应
│   │   ├── pagination/       # 分页工具
│   │   └── middleware/       # 中间件 (CORS, JWT, Logger, Recovery)
│   ├── platform/
│   │   ├── database/         # PostgreSQL 连接池
│   │   ├── cache/            # Redis 客户端
│   │   ├── auth/             # JWT Manager
│   │   └── logger/           # Zap 日志
│   └── modules/
│       ├── auth/             # 认证模块
│       ├── posts/            # 文章模块
│       ├── taxonomy/         # 分类/标签模块
│       ├── comments/         # 评论模块
│       └── shared/errors/    # 共享错误定义
├── migrations/               # 数据库迁移
├── Makefile                  # 构建命令
├── Dockerfile                # Docker 构建
└── go.mod                    # 依赖管理
```

## 变更记录 (Changelog)

| 日期 | 变更 |
|------|------|
| 2026-03-09 | 初始化模块文档 |
