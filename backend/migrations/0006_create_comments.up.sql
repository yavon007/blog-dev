-- 0006_create_comments.up.sql

CREATE TYPE comment_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE IF NOT EXISTS comments (
    id                BIGSERIAL PRIMARY KEY,
    post_id           BIGINT         NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_comment_id BIGINT         REFERENCES comments(id) ON DELETE SET NULL,
    author_name       VARCHAR(100)   NOT NULL,
    author_email      VARCHAR(255)   NOT NULL,
    body              TEXT           NOT NULL,
    status            comment_status NOT NULL DEFAULT 'pending',
    ip_hash           VARCHAR(64),
    user_agent        TEXT,
    created_at        TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_status  ON comments(status);
