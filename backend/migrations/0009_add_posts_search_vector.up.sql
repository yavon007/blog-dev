-- 0009_add_posts_search_vector.up.sql

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

-- Backfill existing rows
UPDATE posts SET search_vector =
    setweight(to_tsvector('simple', COALESCE(title, '')), 'A') ||
    setweight(to_tsvector('simple', COALESCE(summary, '')), 'B') ||
    setweight(to_tsvector('simple', COALESCE(content_md, '')), 'C');
