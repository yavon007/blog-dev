# 实施计划：博客系统高优先级功能

> 生成时间：2026-03-09

## 任务类型
- [x] 全栈 (→ 并行 Codex + Gemini)

## 功能概述

| 功能 | 优先级 | 预估工时 |
|------|--------|----------|
| SEO 优化 | P0 | 2-3 天 |
| RSS 订阅 | P0 | 1 天 |
| 文章归档 | P1 | 1-2 天 |
| 图片管理 | P1 | 3-4 天 |

## 技术方案

采用 **配置+缓存层** 方案：
- 添加 SEO 字段到 posts 表
- 新增 `media` 模块处理图片
- Redis 缓存 sitemap/RSS
- 前端使用 `@unhead/vue` 管理 Meta

---

## Phase 1: SEO 优化

### 后端任务

#### 1.1 数据库迁移

```sql
-- 0010_add_seo_fields.up.sql
ALTER TABLE posts ADD COLUMN seo_title VARCHAR(255);
ALTER TABLE posts ADD COLUMN seo_description TEXT;
ALTER TABLE posts ADD COLUMN og_image_url VARCHAR(500);
```

```sql
-- 0011_create_site_settings.up.sql
CREATE TABLE site_settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO site_settings (key, value) VALUES
('site_title', 'My Blog'),
('site_description', 'A personal blog'),
('site_url', '');
```

#### 1.2 SEO 模块结构

```
backend/internal/modules/seo/
├── core/
│   └── service.go      # Meta 聚合、Sitemap 生成
├── repository/
│   └── postgres.go     # 查询文章/设置
└── transport/http/
    └── handler.go      # /sitemap.xml, /api/v1/seo/meta
```

#### 1.3 API 设计

```
GET /sitemap.xml           # XML sitemap
GET /api/v1/seo/meta       # 获取全局 SEO 设置
PUT /api/v1/admin/seo/meta # 更新 SEO 设置 (需认证)
```

### 前端任务

#### 1.1 安装依赖

```bash
pnpm add @unhead/vue
```

#### 1.2 Meta 管理配置

```typescript
// main.ts
import { createHead } from '@unhead/vue'
const head = createHead()
app.use(head)
```

#### 1.3 组件更新

- `EditorView.vue`: 添加 SEO 字段表单
- `PostDetailView.vue`: 动态设置 Meta
- `HomeView.vue`: 设置首页 Meta

---

## Phase 2: RSS 订阅

### 后端任务

#### 2.1 Feed 模块

```
backend/internal/modules/feed/
├── core/
│   └── service.go      # RSS 2.0 XML 生成
└── transport/http/
    └── handler.go      # /rss.xml
```

#### 2.2 RSS 输出格式

```xml
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>{{ site_title }}</title>
    <link>{{ site_url }}</link>
    <description>{{ site_description }}</description>
    <item>
      <title>Post Title</title>
      <link>https://example.com/post/slug</link>
      <pubDate>Mon, 09 Mar 2026 00:00:00 +0000</pubDate>
      <description>Post summary...</description>
    </item>
  </channel>
</rss>
```

### 前端任务

#### 2.1 RSS 按钮

```vue
<!-- components/front/RssButton.vue -->
<template>
  <a href="/rss.xml" target="_blank"
     class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
     aria-label="订阅 RSS">
    <span class="i-carbon-rss text-xl text-orange-500" />
  </a>
</template>
```

#### 2.2 布局集成

- `FrontLayout.vue`: 头部添加 RSS 按钮
- `FrontLayout.vue`: 页脚添加 RSS 链接

---

## Phase 3: 文章归档

### 后端任务

#### 3.1 归档查询

```go
// posts/repository/postgres.go
func (r *PostgresRepo) GetArchive(ctx context.Context) ([]ArchiveItem, error) {
    query := `
        SELECT
            EXTRACT(YEAR FROM published_at)::int AS year,
            EXTRACT(MONTH FROM published_at)::int AS month,
            COUNT(*) as count
        FROM posts
        WHERE status = 'published'
        GROUP BY year, month
        ORDER BY year DESC, month DESC
    `
    // ...
}

func (r *PostgresRepo) GetByYearMonth(ctx context.Context, year, month int, page, size int) ([]Post, int64, error) {
    // ...
}
```

#### 3.2 API 设计

```
GET /api/v1/posts/archive              # [{year, month, count}]
GET /api/v1/posts/archive/:year/:month # 分页文章列表
```

### 前端任务

#### 3.1 归档页面

```vue
<!-- views/front/ArchiveView.vue -->
<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">文章归档</h1>
    <div v-for="group in archive" :key="group.year" class="mb-8">
      <h2 class="text-xl font-bold mb-4">{{ group.year }}年</h2>
      <div v-for="month in group.months" :key="month.month" class="mb-4">
        <RouterLink :to="`/archive/${group.year}/${month.month}`"
          class="text-primary-600 dark:text-primary-400 hover:underline">
          {{ month.month }}月 ({{ month.count }}篇)
        </RouterLink>
      </div>
    </div>
  </div>
</template>
```

#### 3.2 路由配置

```typescript
// router/index.ts
{
  path: 'archive',
  name: 'Archive',
  component: () => import('@/views/front/ArchiveView.vue'),
},
{
  path: 'archive/:year/:month',
  name: 'ArchiveMonth',
  component: () => import('@/views/front/ArchiveMonthView.vue'),
},
```

---

## Phase 4: 图片管理

### 后端任务

#### 4.1 数据库迁移

```sql
-- 0012_create_media_files.up.sql
CREATE TABLE media_files (
    id BIGSERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    width INT,
    height INT,
    alt_text TEXT,
    storage VARCHAR(50) DEFAULT 'local',
    path VARCHAR(500) NOT NULL,
    url VARCHAR(500) NOT NULL,
    uploaded_by BIGINT REFERENCES admin_users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_media_files_created_at ON media_files(created_at DESC);
```

#### 4.2 Media 模块

```
backend/internal/modules/media/
├── core/
│   ├── entity.go       # MediaFile struct
│   └── service.go      # Upload, Delete, List
├── repository/
│   └── postgres.go     # CRUD
├── storage/
│   └── local.go        # 本地存储实现
└── transport/http/
    └── handler.go      # API handlers
```

#### 4.3 API 设计

```
POST   /api/v1/admin/media              # 上传图片 (multipart)
GET    /api/v1/admin/media              # 图片列表
DELETE /api/v1/admin/media/:id          # 删除图片
```

### 前端任务

#### 4.1 上传组件

```vue
<!-- components/admin/ImageDropzone.vue -->
<template>
  <div
    class="border-2 border-dashed border-gray-300 dark:border-gray-600
           rounded-lg p-8 text-center cursor-pointer
           hover:border-primary-500 hover:bg-gray-50 dark:hover:bg-gray-800"
    @dragover.prevent="isDragging = true"
    @dragleave="isDragging = false"
    @drop.prevent="handleDrop"
    @click="$refs.input.click()">
    <input ref="input" type="file" accept="image/*" class="hidden" @change="handleChange" />
    <span class="i-carbon-cloud-upload text-4xl text-gray-400" />
    <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
      拖拽图片到此处或点击上传
    </p>
  </div>
</template>
```

#### 4.2 媒体库模态框

```vue
<!-- components/admin/MediaLibraryModal.vue -->
<template>
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
    <div class="bg-white dark:bg-gray-800 rounded-lg w-full max-w-4xl max-h-[80vh] overflow-hidden">
      <div class="p-4 border-b border-gray-200 dark:border-gray-700 flex justify-between">
        <h3 class="font-bold">媒体库</h3>
        <button @click="$emit('close')">×</button>
      </div>
      <div class="p-4">
        <ImageDropzone @upload="handleUpload" />
        <ImageGrid :images="images" :selected="selected" @select="selected = $event" />
      </div>
      <div class="p-4 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-2">
        <button class="btn-secondary" @click="$emit('close')">取消</button>
        <button class="btn-primary" @click="confirm">确认选择</button>
      </div>
    </div>
  </div>
</template>
```

#### 4.3 编辑器集成

在 `EditorView.vue` 中：
1. 工具栏添加图片按钮
2. 点击弹出 `MediaLibraryModal`
3. 选择后插入 Markdown 图片语法

---

## 关键文件清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `backend/migrations/0010_add_seo_fields.up.sql` | 新建 | SEO 字段迁移 |
| `backend/migrations/0011_create_site_settings.up.sql` | 新建 | 站点设置表 |
| `backend/migrations/0012_create_media_files.up.sql` | 新建 | 媒体文件表 |
| `backend/internal/modules/seo/` | 新建 | SEO 模块 |
| `backend/internal/modules/feed/` | 新建 | RSS 模块 |
| `backend/internal/modules/media/` | 新建 | 媒体模块 |
| `backend/internal/app/router.go` | 修改 | 注册新路由 |
| `backend/cmd/api/main.go` | 修改 | 注入新模块 |
| `frontend/src/main.ts` | 修改 | 配置 @unhead/vue |
| `frontend/src/router/index.ts` | 修改 | 添加归档路由 |
| `frontend/src/views/front/ArchiveView.vue` | 新建 | 归档页面 |
| `frontend/src/views/admin/MediaView.vue` | 新建 | 媒体管理页 |
| `frontend/src/components/front/RssButton.vue` | 新建 | RSS 按钮 |
| `frontend/src/components/admin/ImageDropzone.vue` | 新建 | 上传组件 |
| `frontend/src/components/admin/ImageGrid.vue` | 新建 | 图片网格 |
| `frontend/src/components/admin/MediaLibraryModal.vue` | 新建 | 图片选择器 |

---

## 风险与缓解

| 风险 | 缓解措施 |
|------|----------|
| 图片存储空间增长 | 配置上传限制(5MB)、支持压缩、软删除机制 |
| Sitemap/RSS 生成性能 | Redis 缓存 + ETag + 文章变更时失效 |
| SEO 字段迁移兼容 | 默认值回退到 title/summary |
| 归档查询性能 | 数据库索引 + 可选物化视图 |
| 上传安全风险 | MIME 校验、文件大小限制、路径遍历检查 |

---

## SESSION_ID

- CODEX_SESSION: `019cd331-0e2e-7f50-bb93-74dbb77997c1`
- GEMINI_SESSION: `59e07c5f-43c8-4b6b-82c8-6183959ece8f`
