CREATE TABLE IF NOT EXISTS site_settings (
  id SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
  site_title VARCHAR(255) NOT NULL DEFAULT 'My Blog',
  site_description TEXT NOT NULL DEFAULT '',
  default_meta_title VARCHAR(255) NOT NULL DEFAULT 'My Blog',
  default_meta_description TEXT NOT NULL DEFAULT '',
  og_image_url VARCHAR(512),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO site_settings (id, site_title, site_description, default_meta_title, default_meta_description, og_image_url)
VALUES (1, 'My Blog', '', 'My Blog', '', NULL)
ON CONFLICT (id) DO NOTHING;
