-- 0009_add_posts_search_vector.down.sql

DROP TRIGGER IF EXISTS posts_search_vector_update ON posts;
DROP FUNCTION IF EXISTS update_posts_search_vector();
