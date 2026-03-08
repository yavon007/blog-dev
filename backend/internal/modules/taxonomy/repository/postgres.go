package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourblog/backend/internal/modules/taxonomy/core"
	sherrors "github.com/yourblog/backend/internal/modules/shared/errors"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// --- Categories ---

func (r *PostgresRepo) ListCategories(ctx context.Context) ([]*core.Category, error) {
	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.name, c.slug, c.description, COUNT(p.id) AS post_count, c.created_at
		FROM categories c
		LEFT JOIN posts p ON p.category_id = c.id AND p.status = 'published'
		GROUP BY c.id ORDER BY c.name`)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var cats []*core.Category
	for rows.Next() {
		var c core.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.PostCount, &c.CreatedAt); err == nil {
			cats = append(cats, &c)
		}
	}
	return cats, nil
}

func (r *PostgresRepo) FindCategoryByID(ctx context.Context, id int64) (*core.Category, error) {
	var c core.Category
	err := r.db.QueryRow(ctx,
		`SELECT id, name, slug, description, 0, created_at FROM categories WHERE id=$1`, id,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.PostCount, &c.CreatedAt)
	if err != nil {
		return nil, sherrors.ErrNotFound
	}
	return &c, nil
}

func (r *PostgresRepo) CreateCategory(ctx context.Context, req core.CreateCategoryRequest) (*core.Category, error) {
	var c core.Category
	err := r.db.QueryRow(ctx,
		`INSERT INTO categories (name, slug, description) VALUES ($1,$2,$3)
		 RETURNING id, name, slug, description, created_at`,
		req.Name, req.Slug, req.Description,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, sherrors.ErrConflict
		}
		return nil, fmt.Errorf("create category: %w", err)
	}
	return &c, nil
}

func (r *PostgresRepo) UpdateCategory(ctx context.Context, id int64, req core.UpdateCategoryRequest) (*core.Category, error) {
	var c core.Category
	err := r.db.QueryRow(ctx,
		`UPDATE categories SET name=$1, slug=$2, description=$3 WHERE id=$4
		 RETURNING id, name, slug, description, created_at`,
		req.Name, req.Slug, req.Description, id,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}
	return &c, nil
}

func (r *PostgresRepo) DeleteCategory(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM categories WHERE id=$1`, id)
	return err
}

// --- Tags ---

func (r *PostgresRepo) ListTags(ctx context.Context) ([]*core.Tag, error) {
	rows, err := r.db.Query(ctx, `
		SELECT t.id, t.name, t.slug, COUNT(pt.post_id) AS post_count, t.created_at
		FROM tags t
		LEFT JOIN post_tags pt ON pt.tag_id = t.id
		GROUP BY t.id ORDER BY t.name`)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer rows.Close()

	var tags []*core.Tag
	for rows.Next() {
		var t core.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.PostCount, &t.CreatedAt); err == nil {
			tags = append(tags, &t)
		}
	}
	return tags, nil
}

func (r *PostgresRepo) FindTagByID(ctx context.Context, id int64) (*core.Tag, error) {
	var t core.Tag
	err := r.db.QueryRow(ctx,
		`SELECT id, name, slug, 0, created_at FROM tags WHERE id=$1`, id,
	).Scan(&t.ID, &t.Name, &t.Slug, &t.PostCount, &t.CreatedAt)
	if err != nil {
		return nil, sherrors.ErrNotFound
	}
	return &t, nil
}

func (r *PostgresRepo) CreateTag(ctx context.Context, req core.CreateTagRequest) (*core.Tag, error) {
	var t core.Tag
	err := r.db.QueryRow(ctx,
		`INSERT INTO tags (name, slug) VALUES ($1,$2) RETURNING id, name, slug, created_at`,
		req.Name, req.Slug,
	).Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, sherrors.ErrConflict
		}
		return nil, fmt.Errorf("create tag: %w", err)
	}
	return &t, nil
}

func (r *PostgresRepo) UpdateTag(ctx context.Context, id int64, req core.UpdateTagRequest) (*core.Tag, error) {
	var t core.Tag
	err := r.db.QueryRow(ctx,
		`UPDATE tags SET name=$1, slug=$2 WHERE id=$3 RETURNING id, name, slug, created_at`,
		req.Name, req.Slug, id,
	).Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("update tag: %w", err)
	}
	return &t, nil
}

func (r *PostgresRepo) DeleteTag(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tags WHERE id=$1`, id)
	return err
}
