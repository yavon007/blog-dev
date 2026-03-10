package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yavon007/blog-dev/backend/internal/modules/seo/core"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) GetSiteSettings(ctx context.Context) (*core.SiteSettings, error) {
	const q = `
		SELECT site_title, site_description, default_meta_title,
		       default_meta_description, COALESCE(og_image_url, ''), updated_at
		FROM site_settings WHERE id = 1`

	var settings core.SiteSettings
	err := r.db.QueryRow(ctx, q).Scan(
		&settings.SiteTitle,
		&settings.SiteDescription,
		&settings.DefaultMetaTitle,
		&settings.DefaultMetaDescription,
		&settings.OgImageURL,
		&settings.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return defaultSettings(), nil
		}
		return nil, fmt.Errorf("get site settings: %w", err)
	}
	return &settings, nil
}

func (r *PostgresRepo) UpsertSiteSettings(ctx context.Context, req core.UpdateMetaRequest) (*core.SiteSettings, error) {
	const q = `
		INSERT INTO site_settings (id, site_title, site_description, default_meta_title, default_meta_description, og_image_url)
		VALUES (1, $1, $2, $3, $4, NULLIF($5, ''))
		ON CONFLICT (id) DO UPDATE SET
			site_title = EXCLUDED.site_title,
			site_description = EXCLUDED.site_description,
			default_meta_title = EXCLUDED.default_meta_title,
			default_meta_description = EXCLUDED.default_meta_description,
			og_image_url = EXCLUDED.og_image_url,
			updated_at = NOW()
		RETURNING site_title, site_description, default_meta_title,
		          default_meta_description, COALESCE(og_image_url, ''), updated_at`

	var settings core.SiteSettings
	err := r.db.QueryRow(ctx, q,
		req.SiteTitle,
		req.SiteDescription,
		req.DefaultMetaTitle,
		req.DefaultMetaDescription,
		req.OgImageURL,
	).Scan(
		&settings.SiteTitle,
		&settings.SiteDescription,
		&settings.DefaultMetaTitle,
		&settings.DefaultMetaDescription,
		&settings.OgImageURL,
		&settings.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert site settings: %w", err)
	}
	return &settings, nil
}

func (r *PostgresRepo) ListPublishedPosts(ctx context.Context, limit int) ([]core.PostMeta, error) {
	baseQuery := `
		SELECT slug,
		       COALESCE(seo_title, title) AS seo_title,
		       COALESCE(seo_description, summary, '') AS seo_description,
		       COALESCE(og_image_url, cover_url, '') AS og_image_url,
		       COALESCE(summary, '') AS summary,
		       published_at,
		       updated_at
		FROM posts
		WHERE status = 'published'
		ORDER BY published_at DESC NULLS LAST, updated_at DESC`

	if limit > 0 {
		baseQuery = fmt.Sprintf("%s LIMIT %d", baseQuery, limit)
	}

	rows, err := r.db.Query(ctx, baseQuery)
	if err != nil {
		return nil, fmt.Errorf("list published posts: %w", err)
	}
	defer rows.Close()

	var items []core.PostMeta
	for rows.Next() {
		var meta core.PostMeta
		err := rows.Scan(
			&meta.Slug,
			&meta.SEOTitle,
			&meta.SEODescription,
			&meta.OGImageURL,
			&meta.Summary,
			&meta.PublishedAt,
			&meta.UpdatedAt,
		)
		if err == nil {
			items = append(items, meta)
		}
	}
	return items, nil
}

func defaultSettings() *core.SiteSettings {
	return &core.SiteSettings{
		SiteTitle:              "My Blog",
		SiteDescription:        "",
		DefaultMetaTitle:       "My Blog",
		DefaultMetaDescription: "",
		OgImageURL:             "",
		UpdatedAt:              time.Now().UTC(),
	}
}
