# 系统架构说明

## 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                     用户浏览器                           │
│  ┌─────────────────────┐  ┌──────────────────────────┐  │
│  │   前台博客 (/)       │  │   管理后台 (/admin)       │  │
│  │  Vue 3 + UnoCSS     │  │  Vue 3 + md-editor-v3    │  │
│  └──────────┬──────────┘  └────────────┬─────────────┘  │
└─────────────┼───────────────────────────┼────────────────┘
              │ HTTP/REST API             │ HTTP/REST API
              │ /api/v1/                  │ /api/v1/admin/
              ▼                           ▼
┌─────────────────────────────────────────────────────────┐
│              Go + Gin 后端 (:8080)                       │
│                                                         │
│  ┌──────────┐  ┌──────────┐  ┌────────┐  ┌──────────┐  │
│  │  auth    │  │  posts   │  │taxonomy│  │ comments │  │
│  │  module  │  │  module  │  │ module │  │  module  │  │
│  └────┬─────┘  └────┬─────┘  └───┬────┘  └────┬─────┘  │
│       └─────────────┴────────────┴─────────────┘        │
│                          │                              │
│              ┌───────────┴───────────┐                  │
│              │   platform layer      │                  │
│              │ (DB, Cache, Auth, Log)│                  │
│              └───────────┬───────────┘                  │
└──────────────────────────┼──────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
   ┌──────────▼──────────┐  ┌──────────▼──────────┐
   │   PostgreSQL 16      │  │     Redis 7          │
   │   (主数据库)          │  │  (缓存 + 限流)        │
   └─────────────────────┘  └─────────────────────┘
```

## 前端架构

### Layout 隔离策略

```
App.vue
  └── RouterView
       ├── FrontLayout.vue      # 前台布局（导航 + 页脚）
       │    └── RouterView
       │         ├── Home.vue
       │         ├── PostDetail.vue
       │         └── Tags.vue
       ├── AdminLayout.vue      # 后台布局（侧边栏 + 顶栏）
       │    └── RouterView
       │         ├── Dashboard.vue
       │         ├── Editor.vue   ← 异步加载 md-editor-v3
       │         ├── PostList.vue
       │         └── ...
       └── Login.vue
```

**代码分割效果**：
- 前台包（FrontLayout + 所有前台页面）：< 200KB gzip
- 后台包（AdminLayout + 管理页面）：按需加载
- Markdown 编辑器包：仅 `/admin/article/edit` 时加载（~400KB）

### 状态管理

```
Pinia Stores
├── useUserStore    → JWT Token, 用户信息, 登录/登出
└── useAppStore     → 深浅色主题, 侧边栏开关
```

## 后端架构

### 请求生命周期

```
HTTP Request
    │
    ▼
Gin Router
    │
    ├── Middleware Chain
    │     ├── RequestID（生成追踪 ID）
    │     ├── Logger（结构化日志）
    │     ├── Recovery（Panic 恢复）
    │     ├── CORS（跨域策略）
    │     ├── RateLimit（IP 限流，Redis 支撑）
    │     └── JWT（管理路由鉴权）
    │
    ▼
Handler（transport/http）
    │  ← DTO 绑定 + 验证
    ▼
Service（core）
    │  ← 业务逻辑
    ▼
Repository（repository）
    │  ← SQL 查询
    ▼
PostgreSQL
```

### 模块职责

| 模块 | 职责 |
|------|------|
| auth | 管理员登录、JWT 签发/验证/刷新/撤销 |
| posts | 文章 CRUD、发布状态管理、全文搜索 |
| taxonomy | 分类 + 标签 CRUD，文章关联 |
| comments | 游客评论提交、审核流程 |
| shared | 通用 DTO、错误码、数据映射 |

## 数据流

### 文章发布流程

```
管理员编辑 → POST /api/v1/admin/posts
    → 验证（标题非空、slug 唯一）
    → 保存 content_md（原始 Markdown）
    → 生成 content_html_cached（服务端渲染）
    → 更新 search_vector（tsvector，全文索引）
    → status = draft

管理员发布 → PATCH /api/v1/admin/posts/:id/status
    → status = published, published_at = now()
```

### 评论审核流程

```
游客提交 → POST /api/v1/posts/:slug/comments
    → status = pending

管理员审核 → PATCH /api/v1/admin/comments/:id
    → status = approved | rejected

公开 API 只返回 status = approved 的评论
```
