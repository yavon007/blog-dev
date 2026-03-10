ALTER TABLE posts
  DROP COLUMN IF EXISTS seo_title,
  DROP COLUMN IF EXISTS seo_description,
  DROP COLUMN IF EXISTS og_image_url;
