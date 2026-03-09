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
