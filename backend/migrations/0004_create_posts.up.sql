-- 0004_create_posts.up.sql

CREATE TYPE post_status AS ENUM ('draft', 'published');

CREATE TABLE IF NOT EXISTS posts (
    id                  BIGSERIAL PRIMARY KEY,
    title               VARCHAR(255) NOT NULL,
    slug                VARCHAR(255) NOT NULL UNIQUE,
    summary             TEXT,
    content_md          TEXT         NOT NULL,
    content_html_cached TEXT,
    cover_url           VARCHAR(512),
    status              post_status  NOT NULL DEFAULT 'draft',
    published_at        TIMESTAMPTZ,
    category_id         BIGINT       REFERENCES categories(id) ON DELETE SET NULL,
    author_id           BIGINT       NOT NULL REFERENCES admin_users(id),
    search_vector       TSVECTOR,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_posts_slug ON posts(slug);
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_published_at ON posts(published_at DESC) WHERE status = 'published';
CREATE INDEX idx_posts_category_id ON posts(category_id);
CREATE INDEX idx_posts_search_vector ON posts USING GIN(search_vector);
