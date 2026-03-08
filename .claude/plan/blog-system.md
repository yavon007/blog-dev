# 博客系统实施计划

> 版本: v1.0 | 日期: 2026-03-08

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + Vite + pnpm + UnoCSS + Vue Router 4 + Pinia |
| 后端 | Go + Gin v1.12 |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7.4 |
| 容器化 | docker-compose |
| 编辑器 | md-editor-v3 (Markdown) |

---

## 项目结构

```
blog-dev/
├── frontend/          # Vue 3 前端
├── backend/           # Go 后端
├── docker-compose.yml # 根级编排（可选）
├── README.md
├── CLAUDE.md
└── docs/
```

---

## 后端结构 (backend/)

```
backend/
├── Makefile
├── go.mod
├── go.sum
├── .env.example
├── README.md
├── configs/
│   └── config.yaml
├── deployments/
│   ├── docker-compose.yml
│   └── Dockerfile.api
├── docs/
│   └── openapi/
│       └── blog-api.yaml
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── bootstrap.go
│   │   ├── router.go
│   │   └── shutdown.go
│   ├── config/
│   │   └── loader.go
│   ├── platform/
│   │   ├── auth/jwt_manager.go
│   │   ├── cache/redis.go
│   │   ├── database/postgres.go
│   │   ├── logger/logger.go
│   │   └── observability/metrics.go
│   ├── pkg/
│   │   ├── middleware/
│   │   │   ├── cors.go
│   │   │   ├── jwt.go
│   │   │   ├── rate_limit.go
│   │   │   ├── request_id.go
│   │   │   ├── logging.go
│   │   │   └── recovery.go
│   │   ├── pagination/pagination.go
│   │   ├── response/response.go
│   │   └── validator/validator.go
│   └── modules/
│       ├── auth/
│       │   ├── core/
│       │   ├── repository/
│       │   └── transport/http/
│       ├── posts/
│       │   ├── core/
│       │   ├── repository/
│       │   └── transport/http/
│       ├── taxonomy/
│       │   ├── core/
│       │   ├── repository/
│       │   └── transport/http/
│       ├── comments/
│       │   ├── core/
│       │   ├── repository/
│       │   └── transport/http/
│       └── shared/
│           ├── dto/
│           ├── errors/
│           └── mapper/
├── migrations/
│   ├── 0001_create_admin_users.up.sql
│   ├── 0001_create_admin_users.down.sql
│   ├── 0002_create_categories.up.sql
│   ├── 0002_create_categories.down.sql
│   ├── 0003_create_tags.up.sql
│   ├── 0003_create_tags.down.sql
│   ├── 0004_create_posts.up.sql
│   ├── 0004_create_posts.down.sql
│   ├── 0005_create_post_tags.up.sql
│   ├── 0005_create_post_tags.down.sql
│   ├── 0006_create_comments.up.sql
│   ├── 0006_create_comments.down.sql
│   ├── 0007_create_refresh_tokens.up.sql
│   ├── 0007_create_refresh_tokens.down.sql
│   ├── 0008_create_audit_logs.up.sql
│   ├── 0008_create_audit_logs.down.sql
│   ├── 0009_add_posts_search_vector.up.sql
│   └── 0009_add_posts_search_vector.down.sql
└── scripts/
    ├── wait-for.sh
    └── migrate.sh
```

### 数据库表设计

```sql
admin_users(id, email UNIQUE, password_hash, role ENUM{owner,editor}, last_login_at, created_at)
categories(id, name UNIQUE, slug UNIQUE, description, created_at)
tags(id, name UNIQUE, slug UNIQUE, created_at)
posts(id, title, slug UNIQUE, summary, content_md, content_html_cached, cover_url,
      status ENUM{draft,published}, published_at, category_id FK, author_id FK,
      search_vector tsvector, created_at, updated_at)
post_tags(post_id FK, tag_id FK, PRIMARY KEY(post_id, tag_id))
comments(id, post_id FK, author_name, author_email, body,
         status ENUM{pending,approved,rejected}, parent_comment_id NULLABLE,
         ip_hash, user_agent, created_at, updated_at)
refresh_tokens(id, admin_id FK, token_hash, expires_at, revoked_at)
audit_logs(id, admin_id FK, action, resource_type, resource_id, meta, created_at)
```

### API 路由设计

```
公开 API (无鉴权)
GET    /api/v1/posts                  列表（分页、搜索、分类、标签过滤）
GET    /api/v1/posts/:slug            文章详情
GET    /api/v1/posts/:slug/comments   评论列表
POST   /api/v1/posts/:slug/comments   提交评论（待审核）
GET    /api/v1/categories             分类列表
GET    /api/v1/tags                   标签列表
GET    /api/v1/health                 健康检查

管理 API (JWT 鉴权)
POST   /api/v1/admin/auth/login       登录
POST   /api/v1/admin/auth/refresh     刷新 Token
POST   /api/v1/admin/auth/logout      登出
CRUD   /api/v1/admin/posts            文章管理
PATCH  /api/v1/admin/posts/:id/status 发布/下线
CRUD   /api/v1/admin/categories       分类管理
CRUD   /api/v1/admin/tags             标签管理
GET    /api/v1/admin/comments         评论列表（含审核状态）
PATCH  /api/v1/admin/comments/:id     审核（approve/reject）
DELETE /api/v1/admin/comments/:id     删除
```

### 后端依赖

```
github.com/gin-gonic/gin v1.12.0
github.com/gin-contrib/cors v1.7.6
github.com/golang-jwt/jwt/v5 v5.2.2
github.com/jackc/pgx/v5 v5.6.0
github.com/redis/go-redis/v9 v9.7.3
go.uber.org/zap v1.27.1
github.com/go-playground/validator/v10 v10.22.1
github.com/caarlos0/env/v11 v11.3.1
github.com/stretchr/testify v1.9.0
```

---

## 前端结构 (frontend/)

```
frontend/
├── package.json
├── vite.config.ts
├── uno.config.ts
├── tsconfig.json
├── index.html
└── src/
    ├── main.ts
    ├── App.vue
    ├── api/
    │   ├── article.ts
    │   ├── category.ts
    │   ├── comment.ts
    │   └── user.ts
    ├── assets/
    ├── components/
    │   ├── common/        # 通用组件 (Toast, Modal, Skeleton)
    │   ├── front/         # 前台组件 (PostCard, TagBadge)
    │   └── admin/         # 后台组件 (DataTable, StatusBadge)
    ├── layouts/
    │   ├── FrontLayout.vue
    │   └── AdminLayout.vue
    ├── router/
    │   └── index.ts
    ├── store/
    │   ├── index.ts
    │   ├── user.ts        # Token, 用户信息, 登出
    │   └── app.ts         # 主题, 侧边栏状态
    ├── utils/
    │   └── request.ts     # Axios 封装
    └── views/
        ├── front/
        │   ├── Home.vue
        │   ├── PostDetail.vue
        │   └── Tags.vue
        └── admin/
            ├── Login.vue
            ├── Dashboard.vue
            ├── PostList.vue
            ├── Editor.vue          # 懒加载 md-editor-v3
            ├── CategoryManage.vue
            ├── TagManage.vue
            └── CommentManage.vue
```

### 前端依赖

```json
dependencies:
  vue: ^3.4.21
  vue-router: ^4.3.0
  pinia: ^2.1.7
  axios: ^1.6.8
  md-editor-v3: ^4.11.0
  @vueuse/core: ^10.9.0
  @unocss/reset: ^0.58.0

devDependencies:
  vite: ^5.2.6
  @vitejs/plugin-vue: ^5.0.4
  unocss: ^0.58.0
  @unocss/preset-typography: ^0.58.0
  @unocss/preset-icons: ^0.58.0
  @iconify-json/carbon: ^1.1.31
  typescript: ^5.4.3
  vue-tsc: ^2.0.7
```

---

## 文档输出

- `README.md` - 项目说明、快速开始、环境要求
- `CLAUDE.md` - AI 开发规范（代码风格、架构约束、目录说明）
- `docs/api.md` - API 文档说明
- `docs/database.md` - 数据库设计说明
- `docs/architecture.md` - 系统架构说明
- `backend/docs/openapi/blog-api.yaml` - OpenAPI 规范

---

## 实施步骤

### Phase 1: 项目脚手架（优先）
1. 初始化 Git，创建 .gitignore
2. 创建根级 README.md、CLAUDE.md、docs/ 目录
3. 初始化后端 Go 模块，创建目录结构
4. 初始化前端 Vite 项目，安装依赖

### Phase 2: 后端核心
5. 编写数据库迁移文件 (0001-0009)
6. 实现 platform 层（DB、Cache、JWT、Logger）
7. 实现中间件（CORS、JWT、Rate Limit、Logging、Recovery）
8. 实现 modules（auth → taxonomy → posts → comments）
9. 实现路由注册和应用启动

### Phase 3: 前端核心
10. 配置 Vite + UnoCSS + Router + Pinia
11. 实现布局组件（FrontLayout、AdminLayout）
12. 实现 Axios 请求层和 Pinia stores
13. 实现前台页面（Home、PostDetail、Tags）
14. 实现后台页面（Login、Dashboard、Editor、管理页）

### Phase 4: 容器化与文档
15. 编写 docker-compose.yml（含 postgres、redis、api、migrate）
16. 编写 Dockerfile.api（多阶段构建）
17. 完善 Makefile 命令
18. 编写完整文档（README、CLAUDE.md、docs/）
