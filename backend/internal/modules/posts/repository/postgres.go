package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yavon007/blog-dev/backend/internal/modules/posts/core"
	"github.com/yavon007/blog-dev/backend/internal/pkg/pagination"
	sherrors "github.com/yavon007/blog-dev/backend/internal/modules/shared/errors"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) FindBySlug(ctx context.Context, slug string) (*core.Post, error) {
	const q = `
		SELECT p.id, p.title, p.slug, p.summary, p.content_md, p.content_html_cached,
		       p.cover_url, p.status, p.published_at, p.category_id,
		       COALESCE(c.name, ''), p.author_id, p.created_at, p.updated_at
		FROM posts p
		LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.slug = $1`
	row := r.db.QueryRow(ctx, q, slug)
	post, err := scanPost(row)
	if err != nil {
		return nil, err
	}
	post.Tags, _ = r.findTagsByPostID(ctx, post.ID)
	return post, nil
}

func (r *PostgresRepo) FindByID(ctx context.Context, id int64) (*core.Post, error) {
	const q = `
		SELECT p.id, p.title, p.slug, p.summary, p.content_md, p.content_html_cached,
		       p.cover_url, p.status, p.published_at, p.category_id,
		       COALESCE(c.name, ''), p.author_id, p.created_at, p.updated_at
		FROM posts p
		LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1`
	row := r.db.QueryRow(ctx, q, id)
	post, err := scanPost(row)
	if err != nil {
		return nil, err
	}
	post.Tags, _ = r.findTagsByPostID(ctx, post.ID)
	return post, nil
}

func (r *PostgresRepo) List(ctx context.Context, filter core.ListFilter, p pagination.Params) ([]*core.Post, int64, error) {
	var (
		where  []string
		args   []interface{}
		argIdx = 1
	)

	if filter.Status != "" {
		where = append(where, fmt.Sprintf("p.status = $%d", argIdx))
		args = append(args, string(filter.Status))
		argIdx++
	}
	if filter.Category != "" {
		where = append(where, fmt.Sprintf("c.slug = $%d", argIdx))
		args = append(args, filter.Category)
		argIdx++
	}
	if filter.Tag != "" {
		where = append(where, fmt.Sprintf("EXISTS (SELECT 1 FROM post_tags pt JOIN tags t ON t.id = pt.tag_id WHERE pt.post_id = p.id AND t.slug = $%d)", argIdx))
		args = append(args, filter.Tag)
		argIdx++
	}
	if filter.Query != "" {
		where = append(where, fmt.Sprintf("p.search_vector @@ plainto_tsquery('simple', $%d)", argIdx))
		args = append(args, filter.Query)
		argIdx++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM posts p LEFT JOIN categories c ON c.id = p.category_id %s`, whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count posts: %w", err)
	}

	args = append(args, p.Size, p.Offset())
	listQ := fmt.Sprintf(`
		SELECT p.id, p.title, p.slug, p.summary, p.content_md, p.content_html_cached,
		       p.cover_url, p.status, p.published_at, p.category_id,
		       COALESCE(c.name, ''), p.author_id, p.created_at, p.updated_at
		FROM posts p
		LEFT JOIN categories c ON c.id = p.category_id
		%s ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	rows, err := r.db.Query(ctx, listQ, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()

	var posts []*core.Post
	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			continue
		}
		posts = append(posts, p)
	}

	for _, post := range posts {
		post.Tags, _ = r.findTagsByPostID(ctx, post.ID)
	}

	return posts, total, nil
}

func (r *PostgresRepo) Create(ctx context.Context, req core.CreatePostRequest, authorID int64) (*core.Post, error) {
	var publishedAt interface{} = nil
	if req.Status == core.StatusPublished {
		publishedAt = time.Now()
	}

	const q = `
		INSERT INTO posts (title, slug, summary, content_md, cover_url, status, published_at, category_id, author_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	var id int64
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(ctx, q,
		req.Title, req.Slug, req.Summary, req.ContentMD,
		req.CoverURL, string(req.Status), publishedAt, req.CategoryID, authorID,
	).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, sherrors.ErrConflict
		}
		return nil, fmt.Errorf("create post: %w", err)
	}

	_ = r.syncTags(ctx, id, req.TagIDs)

	return r.FindByID(ctx, id)
}

func (r *PostgresRepo) Update(ctx context.Context, id int64, req core.UpdatePostRequest) (*core.Post, error) {
	var publishedAt interface{} = nil
	if req.Status == core.StatusPublished {
		publishedAt = time.Now()
	}

	const q = `
		UPDATE posts SET title=$1, slug=$2, summary=$3, content_md=$4, cover_url=$5,
		                 status=$6, published_at=COALESCE($7, published_at), category_id=$8, updated_at=NOW()
		WHERE id=$9`
	_, err := r.db.Exec(ctx, q,
		req.Title, req.Slug, req.Summary, req.ContentMD,
		req.CoverURL, string(req.Status), publishedAt, req.CategoryID, id,
	)
	if err != nil {
		return nil, fmt.Errorf("update post: %w", err)
	}
	_ = r.syncTags(ctx, id, req.TagIDs)
	return r.FindByID(ctx, id)
}

func (r *PostgresRepo) UpdateStatus(ctx context.Context, id int64, status core.PostStatus) error {
	var publishedAt interface{} = nil
	if status == core.StatusPublished {
		publishedAt = time.Now()
	}
	_, err := r.db.Exec(ctx,
		`UPDATE posts SET status=$1, published_at=$2, updated_at=NOW() WHERE id=$3`,
		string(status), publishedAt, id,
	)
	return err
}

func (r *PostgresRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM posts WHERE id=$1`, id)
	return err
}

func (r *PostgresRepo) syncTags(ctx context.Context, postID int64, tagIDs []int64) error {
	_, _ = r.db.Exec(ctx, `DELETE FROM post_tags WHERE post_id=$1`, postID)
	for _, tid := range tagIDs {
		_, _ = r.db.Exec(ctx, `INSERT INTO post_tags (post_id, tag_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, postID, tid)
	}
	return nil
}

func (r *PostgresRepo) findTagsByPostID(ctx context.Context, postID int64) ([]core.Tag, error) {
	const q = `SELECT t.id, t.name, t.slug FROM tags t
	            JOIN post_tags pt ON pt.tag_id = t.id WHERE pt.post_id = $1 ORDER BY t.name`
	rows, err := r.db.Query(ctx, q, postID)
	if err != nil {
		return []core.Tag{}, err
	}
	defer rows.Close()
	tags := []core.Tag{}
	for rows.Next() {
		var t core.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err == nil {
			tags = append(tags, t)
		}
	}
	return tags, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanPost(row scannable) (*core.Post, error) {
	var p core.Post
	err := row.Scan(
		&p.ID, &p.Title, &p.Slug, &p.Summary, &p.ContentMD, &p.ContentHTMLCached,
		&p.CoverURL, &p.Status, &p.PublishedAt, &p.CategoryID, &p.CategoryName,
		&p.AuthorID, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, sherrors.ErrNotFound
	}
	p.Tags = []core.Tag{}
	return &p, nil
}
