package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourblog/backend/internal/modules/auth/core"
	sherrors "github.com/yourblog/backend/internal/modules/shared/errors"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) FindByEmail(ctx context.Context, email string) (*core.AdminUser, error) {
	const q = `SELECT id, email, password_hash, role, last_login_at, created_at, updated_at
               FROM admin_users WHERE email = $1`
	row := r.db.QueryRow(ctx, q, email)
	return scanAdminUser(row)
}

func (r *PostgresRepo) FindByID(ctx context.Context, id int64) (*core.AdminUser, error) {
	const q = `SELECT id, email, password_hash, role, last_login_at, created_at, updated_at
               FROM admin_users WHERE id = $1`
	row := r.db.QueryRow(ctx, q, id)
	return scanAdminUser(row)
}

func (r *PostgresRepo) UpdateLastLogin(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `UPDATE admin_users SET last_login_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *PostgresRepo) SaveRefreshToken(ctx context.Context, adminID int64, tokenHash string, expiresAt time.Time) error {
	const q = `INSERT INTO refresh_tokens (admin_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, q, adminID, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}
	return nil
}

func (r *PostgresRepo) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1 AND revoked_at IS NULL`,
		tokenHash,
	)
	return err
}

func (r *PostgresRepo) FindRefreshToken(ctx context.Context, tokenHash string) (int64, error) {
	var adminID int64
	err := r.db.QueryRow(ctx,
		`SELECT admin_id FROM refresh_tokens
         WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()`,
		tokenHash,
	).Scan(&adminID)
	if err != nil {
		return 0, sherrors.ErrUnauthorized
	}
	return adminID, nil
}

func scanAdminUser(row pgx.Row) (*core.AdminUser, error) {
	var u core.AdminUser
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, sherrors.ErrNotFound
	}
	return &u, nil
}
