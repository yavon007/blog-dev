# 数据库设计

数据库：PostgreSQL 16

## ER 图

```
admin_users
  ├── 1:N → posts (author_id)
  └── 1:N → refresh_tokens (admin_id)

categories
  └── 1:N → posts (category_id)

tags
  └── M:N → posts (via post_tags)

posts
  ├── 1:N → comments (post_id)
  └── M:N → tags (via post_tags)

comments
  └── 1:N → comments (parent_comment_id, 评论嵌套)
```

## 表结构

### admin_users（管理员）

```sql
CREATE TABLE admin_users (
  id            BIGSERIAL PRIMARY KEY,
  email         VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role          VARCHAR(20)  NOT NULL DEFAULT 'editor',  -- owner | editor
  last_login_at TIMESTAMPTZ,
  created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
```

### categories（分类）

```sql
CREATE TABLE categories (
  id          BIGSERIAL PRIMARY KEY,
  name        VARCHAR(100) NOT NULL UNIQUE,
  slug        VARCHAR(100) NOT NULL UNIQUE,
  description TEXT,
  created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
```

### tags（标签）

```sql
CREATE TABLE tags (
  id         BIGSERIAL PRIMARY KEY,
  name       VARCHAR(100) NOT NULL UNIQUE,
  slug       VARCHAR(100) NOT NULL UNIQUE,
  created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
```

### posts（文章）

```sql
CREATE TABLE posts (
  id                 BIGSERIAL PRIMARY KEY,
  title              VARCHAR(255) NOT NULL,
  slug               VARCHAR(255) NOT NULL UNIQUE,
  summary            TEXT,
  content_md         TEXT         NOT NULL,          -- 原始 Markdown
  content_html_cached TEXT,                          -- 服务端渲染缓存
  cover_url          VARCHAR(512),
  status             VARCHAR(20)  NOT NULL DEFAULT 'draft',  -- draft | published
  published_at       TIMESTAMPTZ,
  category_id        BIGINT       REFERENCES categories(id) ON DELETE SET NULL,
  author_id          BIGINT       NOT NULL REFERENCES admin_users(id),
  search_vector      TSVECTOR,                       -- 全文搜索索引
  created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_slug ON posts(slug);
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_published_at ON posts(published_at DESC);
CREATE INDEX idx_posts_search_vector ON posts USING GIN(search_vector);
```

### post_tags（文章标签关联）

```sql
CREATE TABLE post_tags (
  post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  tag_id  BIGINT NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
  PRIMARY KEY (post_id, tag_id)
);
```

### comments（评论）

```sql
CREATE TABLE comments (
  id                BIGSERIAL PRIMARY KEY,
  post_id           BIGINT       NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  parent_comment_id BIGINT       REFERENCES comments(id) ON DELETE SET NULL,
  author_name       VARCHAR(100) NOT NULL,
  author_email      VARCHAR(255) NOT NULL,
  body              TEXT         NOT NULL,
  status            VARCHAR(20)  NOT NULL DEFAULT 'pending',  -- pending | approved | rejected
  ip_hash           VARCHAR(64),
  user_agent        TEXT,
  created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_status ON comments(status);
```

### refresh_tokens（JWT 刷新令牌）

```sql
CREATE TABLE refresh_tokens (
  id         BIGSERIAL PRIMARY KEY,
  admin_id   BIGINT      NOT NULL REFERENCES admin_users(id) ON DELETE CASCADE,
  token_hash VARCHAR(64) NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### audit_logs（操作日志）

```sql
CREATE TABLE audit_logs (
  id            BIGSERIAL PRIMARY KEY,
  admin_id      BIGINT      REFERENCES admin_users(id) ON DELETE SET NULL,
  action        VARCHAR(100) NOT NULL,
  resource_type VARCHAR(50)  NOT NULL,
  resource_id   BIGINT,
  meta          JSONB,
  created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
```

## 全文搜索

使用 PostgreSQL 原生 `tsvector` 实现文章全文搜索：

```sql
-- search_vector 自动更新触发器
CREATE OR REPLACE FUNCTION update_posts_search_vector()
RETURNS TRIGGER AS $$
BEGIN
  NEW.search_vector :=
    setweight(to_tsvector('simple', COALESCE(NEW.title, '')), 'A') ||
    setweight(to_tsvector('simple', COALESCE(NEW.summary, '')), 'B') ||
    setweight(to_tsvector('simple', COALESCE(NEW.content_md, '')), 'C');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER posts_search_vector_update
  BEFORE INSERT OR UPDATE ON posts
  FOR EACH ROW EXECUTE FUNCTION update_posts_search_vector();
```

查询示例：
```sql
SELECT * FROM posts
WHERE search_vector @@ plainto_tsquery('simple', '关键词')
  AND status = 'published'
ORDER BY ts_rank(search_vector, plainto_tsquery('simple', '关键词')) DESC;
```
