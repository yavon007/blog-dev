package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourblog/backend/internal/modules/comments/core"
	"github.com/yourblog/backend/internal/pkg/pagination"
	sherrors "github.com/yourblog/backend/internal/modules/shared/errors"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) List(ctx context.Context, filter core.ListCommentsFilter, p pagination.Params) ([]*core.Comment, int64, error) {
	args := []interface{}{}
	argIdx := 1
	where := "WHERE 1=1"

	if filter.PostID > 0 {
		where += fmt.Sprintf(" AND post_id=$%d", argIdx)
		args = append(args, filter.PostID)
		argIdx++
	}
	if filter.Status != "" {
		where += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, string(filter.Status))
		argIdx++
	}

	var total int64
	if err := r.db.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM comments %s`, where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, p.Size, p.Offset())
	q := fmt.Sprintf(`
		SELECT id, post_id, parent_comment_id, author_name, author_email, body, status, created_at, updated_at
		FROM comments %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1)

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []*core.Comment
	for rows.Next() {
		var c core.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParentCommentID, &c.AuthorName, &c.AuthorEmail,
			&c.Body, &c.Status, &c.CreatedAt, &c.UpdatedAt); err == nil {
			comments = append(comments, &c)
		}
	}
	return comments, total, nil
}

func (r *PostgresRepo) FindByID(ctx context.Context, id int64) (*core.Comment, error) {
	var c core.Comment
	err := r.db.QueryRow(ctx,
		`SELECT id, post_id, parent_comment_id, author_name, author_email, body, status, created_at, updated_at
		 FROM comments WHERE id=$1`, id,
	).Scan(&c.ID, &c.PostID, &c.ParentCommentID, &c.AuthorName, &c.AuthorEmail,
		&c.Body, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, sherrors.ErrNotFound
	}
	return &c, nil
}

func (r *PostgresRepo) Create(ctx context.Context, postID int64, req core.CreateCommentRequest, ipHash, ua string) (*core.Comment, error) {
	var c core.Comment
	err := r.db.QueryRow(ctx, `
		INSERT INTO comments (post_id, parent_comment_id, author_name, author_email, body, ip_hash, user_agent)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, post_id, parent_comment_id, author_name, author_email, body, status, created_at, updated_at`,
		postID, req.ParentCommentID, req.AuthorName, req.AuthorEmail, req.Body, ipHash, ua,
	).Scan(&c.ID, &c.PostID, &c.ParentCommentID, &c.AuthorName, &c.AuthorEmail,
		&c.Body, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	return &c, nil
}

func (r *PostgresRepo) UpdateStatus(ctx context.Context, id int64, status core.CommentStatus) error {
	_, err := r.db.Exec(ctx,
		`UPDATE comments SET status=$1, updated_at=NOW() WHERE id=$2`, string(status), id)
	return err
}

func (r *PostgresRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM comments WHERE id=$1`, id)
	return err
}
