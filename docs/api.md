# API 文档

Base URL: `http://localhost:8080/api/v1`

## 认证

管理 API 需在请求头携带 JWT Token：

```
Authorization: Bearer <access_token>
```

## 统一响应格式

**成功（单条/数组）**：
```json
{
  "code": 200,
  "data": {},
  "message": "success"
}
```

**成功（分页）**：
```json
{
  "code": 200,
  "data": {
    "items": [],
    "total": 100,
    "page": 1,
    "size": 10
  },
  "message": "success"
}
```

**失败**：
```json
{
  "code": 400,
  "message": "错误描述"
}
```

---

## 数据结构

### Post（文章）

```json
{
  "id": 1,
  "title": "文章标题",
  "slug": "article-slug",
  "summary": "摘要",
  "content_md": "# Markdown 内容",
  "content_html_cached": "<h1>渲染后的 HTML</h1>",
  "cover_url": "https://example.com/cover.jpg",
  "status": "published",
  "published_at": "2024-01-01T00:00:00Z",
  "category_id": 1,
  "category_name": "技术",
  "author_id": 1,
  "tags": [
    { "id": 1, "name": "Go", "slug": "go" }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int | 文章 ID |
| `title` | string | 标题 |
| `slug` | string | URL 别名 |
| `summary` | string | 摘要 |
| `content_md` | string | Markdown 原文 |
| `content_html_cached` | string | 渲染后 HTML（缓存） |
| `cover_url` | string | 封面图 URL |
| `status` | string | `draft` \| `published` |
| `published_at` | string \| null | 发布时间（ISO 8601） |
| `category_id` | int \| null | 分类 ID |
| `category_name` | string | 分类名称 |
| `author_id` | int | 作者 ID |
| `tags` | Tag[] | 标签列表 |
| `created_at` | string | 创建时间 |
| `updated_at` | string | 更新时间 |

### Tag（标签）

```json
{
  "id": 1,
  "name": "Go",
  "slug": "go",
  "post_count": 5,
  "created_at": "2024-01-01T00:00:00Z"
}
```

> 文章内嵌 `tags` 数组仅含 `id`、`name`、`slug` 三个字段，无 `post_count`。

### Category（分类）

```json
{
  "id": 1,
  "name": "技术",
  "slug": "tech",
  "description": "技术相关文章",
  "post_count": 10,
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Comment（评论）

```json
{
  "id": 1,
  "post_id": 1,
  "parent_comment_id": null,
  "author_name": "张三",
  "author_email": "zhangsan@example.com",
  "body": "评论内容",
  "status": "approved",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int | 评论 ID |
| `post_id` | int | 所属文章 ID |
| `parent_comment_id` | int \| null | 父评论 ID（回复） |
| `author_name` | string | 作者昵称 |
| `author_email` | string | 作者邮箱 |
| `body` | string | 评论内容 |
| `status` | string | `pending` \| `approved` \| `rejected` |
| `created_at` | string | 创建时间 |
| `updated_at` | string | 更新时间 |

### TokenPair（令牌）

```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "expires_in": 900
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `access_token` | string | 访问令牌（JWT） |
| `refresh_token` | string | 刷新令牌 |
| `expires_in` | int | 访问令牌有效期（秒） |

---

## 公开 API

### 健康检查

```
GET /health
```

Response data: `{ "status": "ok" }`

---

### 文章

#### 获取文章列表

```
GET /posts
```

Query 参数：

| 参数 | 类型 | 说明 |
|------|------|------|
| `page` | int | 页码，默认 1 |
| `page_size` | int | 每页数量，默认 10，最大 50 |
| `category` | string | 分类 slug |
| `tag` | string | 标签 slug |
| `q` | string | 搜索关键词（全文搜索） |

Response data: 分页格式，`items` 为 Post[]

#### 获取文章详情

```
GET /posts/:slug
```

Response data: Post

#### 获取文章评论

```
GET /posts/:slug/comments
```

Query 参数：`page`、`page_size`

Response data: 分页格式，`items` 为 Comment[]

#### 提交评论

```
POST /posts/:slug/comments
```

Body：
```json
{
  "author_name": "张三",
  "author_email": "zhangsan@example.com",
  "body": "评论内容",
  "parent_comment_id": null
}
```

Response data: Comment（status 初始为 `pending`）

---

### 分类 & 标签

```
GET /categories
GET /tags
```

Response data: Category[] / Tag[]

---

## 管理 API

> 以下接口（除登录/刷新外）均需携带 `Authorization: Bearer <access_token>`

### 认证

#### 登录

```
POST /auth/login
```

Body：
```json
{
  "email": "admin@example.com",
  "password": "your_password"
}
```

Response data: TokenPair

#### 刷新 Token

```
POST /auth/refresh
```

Body：
```json
{
  "refresh_token": "eyJ..."
}
```

Response data: TokenPair

#### 登出

```
POST /admin/auth/logout    ← 需要 Authorization header
```

Body：
```json
{
  "refresh_token": "eyJ..."
}
```

---

### 文章管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/posts` | 文章列表（含草稿，支持 `status`/`q` 过滤） |
| POST | `/admin/posts` | 创建文章 |
| GET | `/admin/posts/:id` | 文章详情 |
| PUT | `/admin/posts/:id` | 更新文章 |
| DELETE | `/admin/posts/:id` | 删除文章 |
| PATCH | `/admin/posts/:id/status` | 变更发布状态 |

创建/更新文章 Body：
```json
{
  "title": "文章标题",
  "slug": "article-slug",
  "summary": "文章摘要",
  "content_md": "# Markdown 内容",
  "cover_url": "https://example.com/cover.jpg",
  "status": "draft",
  "category_id": 1,
  "tag_ids": [1, 2, 3]
}
```

PATCH status Body：
```json
{
  "status": "published"
}
```

---

### 分类管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/categories` | 列表 |
| POST | `/admin/categories` | 创建 |
| PUT | `/admin/categories/:id` | 更新 |
| DELETE | `/admin/categories/:id` | 删除 |

Body：
```json
{
  "name": "技术",
  "slug": "tech",
  "description": "技术相关文章"
}
```

---

### 标签管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/tags` | 列表 |
| POST | `/admin/tags` | 创建 |
| PUT | `/admin/tags/:id` | 更新 |
| DELETE | `/admin/tags/:id` | 删除 |

Body：
```json
{
  "name": "Go",
  "slug": "go"
}
```

---

### 评论管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/comments` | 评论列表 |
| PATCH | `/admin/comments/:id` | 审核 |
| DELETE | `/admin/comments/:id` | 删除 |

GET 支持 Query 参数：`post_id`、`status`（`pending`/`approved`/`rejected`）、`page`、`page_size`

PATCH Body：
```json
{
  "status": "approved"
}
```
