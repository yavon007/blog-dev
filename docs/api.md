# API 文档

Base URL: `http://localhost:8080/api/v1`

## 认证

管理 API 需在请求头携带 JWT Token：

```
Authorization: Bearer <access_token>
```

## 统一响应格式

**成功**：
```json
{
  "code": 200,
  "data": {},
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

## 公开 API

### 文章

#### 获取文章列表

```
GET /posts
```

Query 参数：

| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 10，最大 50 |
| category | string | 分类 slug |
| tag | string | 标签 slug |
| q | string | 搜索关键词（全文搜索） |

#### 获取文章详情

```
GET /posts/:slug
```

#### 获取文章评论

```
GET /posts/:slug/comments?page=1&page_size=20
```

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

### 分类 & 标签

```
GET /categories     # 所有分类
GET /tags           # 所有标签
```

### 健康检查

```
GET /health
```

---

## 管理 API

### 认证

#### 登录

```
POST /admin/auth/login
```

Body：
```json
{
  "email": "admin@example.com",
  "password": "your_password"
}
```

Response data：
```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "expires_in": 900
}
```

#### 刷新 Token

```
POST /admin/auth/refresh
```

Body：
```json
{
  "refresh_token": "eyJ..."
}
```

#### 登出

```
POST /admin/auth/logout
```

---

### 文章管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /admin/posts | 文章列表（含草稿） |
| POST | /admin/posts | 创建文章 |
| GET | /admin/posts/:id | 文章详情 |
| PUT | /admin/posts/:id | 更新文章 |
| DELETE | /admin/posts/:id | 删除文章 |
| PATCH | /admin/posts/:id/status | 发布/下线 |

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

### 分类管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /admin/categories | 列表 |
| POST | /admin/categories | 创建 |
| PUT | /admin/categories/:id | 更新 |
| DELETE | /admin/categories/:id | 删除 |

### 标签管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /admin/tags | 列表 |
| POST | /admin/tags | 创建 |
| PUT | /admin/tags/:id | 更新 |
| DELETE | /admin/tags/:id | 删除 |

### 评论管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /admin/comments | 评论列表（支持 status 过滤） |
| PATCH | /admin/comments/:id | 审核（approve/reject） |
| DELETE | /admin/comments/:id | 删除 |

PATCH Body：
```json
{
  "status": "approved"
}
```
