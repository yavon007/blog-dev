package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yavon007/blog-dev/backend/internal/modules/media/core"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, file *core.MediaFile) error {
	const q = `
		INSERT INTO media_files (filename, original_name, mime_type, size, width, height, alt_text, storage, path, url, uploaded_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`
	return r.db.QueryRow(ctx, q,
		file.Filename, file.OriginalName, file.MimeType, file.Size,
		file.Width, file.Height, file.AltText, file.Storage,
		file.Path, file.URL, file.UploadedBy, file.CreatedAt,
	).Scan(&file.ID)
}

func (r *PostgresRepo) FindByID(ctx context.Context, id int64) (*core.MediaFile, error) {
	const q = `
		SELECT id, filename, original_name, mime_type, size, width, height, alt_text, storage, path, url, uploaded_by, created_at, deleted_at
		FROM media_files WHERE id = $1 AND deleted_at IS NULL
	`
	var f core.MediaFile
	err := r.db.QueryRow(ctx, q, id).Scan(
		&f.ID, &f.Filename, &f.OriginalName, &f.MimeType, &f.Size,
		&f.Width, &f.Height, &f.AltText, &f.Storage, &f.Path, &f.URL,
		&f.UploadedBy, &f.CreatedAt, &f.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("media file not found")
	}
	return &f, nil
}

func (r *PostgresRepo) List(ctx context.Context, filter core.ListMediaFilter) ([]*core.MediaFile, int64, error) {
	var where []string
	var args []interface{}
	argIdx := 1

	if filter.MimeType != "" {
		where = append(where, fmt.Sprintf("mime_type LIKE $%d", argIdx))
		args = append(args, filter.MimeType+"%")
		argIdx++
	}

	where = append(where, "deleted_at IS NULL")
	whereClause := "WHERE " + strings.Join(where, " AND ")

	// Count
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM media_files %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count media: %w", err)
	}

	// List
	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)
	listQ := fmt.Sprintf(`
		SELECT id, filename, original_name, mime_type, size, width, height, alt_text, storage, path, url, uploaded_by, created_at, deleted_at
		FROM media_files %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	rows, err := r.db.Query(ctx, listQ, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list media: %w", err)
	}
	defer rows.Close()

	var files []*core.MediaFile
	for rows.Next() {
		var f core.MediaFile
		if err := rows.Scan(
			&f.ID, &f.Filename, &f.OriginalName, &f.MimeType, &f.Size,
			&f.Width, &f.Height, &f.AltText, &f.Storage, &f.Path, &f.URL,
			&f.UploadedBy, &f.CreatedAt, &f.DeletedAt,
		); err == nil {
			files = append(files, &f)
		}
	}

	return files, total, nil
}

func (r *PostgresRepo) SoftDelete(ctx context.Context, id int64) error {
	const q = `UPDATE media_files SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, q, id)
	return err
}
