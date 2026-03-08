package core

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/yavon007/blog-dev/backend/internal/platform/auth"
	sherrors "github.com/yavon007/blog-dev/backend/internal/modules/shared/errors"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*AdminUser, error)
	UpdateLastLogin(ctx context.Context, id int64) error
	SaveRefreshToken(ctx context.Context, adminID int64, tokenHash string, expiresAt time.Time) error
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	FindRefreshToken(ctx context.Context, tokenHash string) (int64, error) // returns adminID
	FindByID(ctx context.Context, id int64) (*AdminUser, error)
}

type Service struct {
	repo   Repository
	jwtMgr *auth.Manager
}

func NewService(repo Repository, jwtMgr *auth.Manager) *Service {
	return &Service{repo: repo, jwtMgr: jwtMgr}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*TokenPair, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, sherrors.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, sherrors.ErrUnauthorized
	}

	return s.generateTokenPair(ctx, user)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	hash := hashToken(refreshToken)
	adminID, err := s.repo.FindRefreshToken(ctx, hash)
	if err != nil {
		return nil, sherrors.ErrUnauthorized
	}

	user, err := s.repo.FindByID(ctx, adminID)
	if err != nil {
		return nil, sherrors.ErrUnauthorized
	}

	// Rotate: revoke old, issue new
	_ = s.repo.RevokeRefreshToken(ctx, hash)
	return s.generateTokenPair(ctx, user)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	hash := hashToken(refreshToken)
	return s.repo.RevokeRefreshToken(ctx, hash)
}

func (s *Service) generateTokenPair(ctx context.Context, user *AdminUser) (*TokenPair, error) {
	accessToken, err := s.jwtMgr.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.jwtMgr.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	expiresAt := time.Now().Add(s.jwtMgr.RefreshTTL())
	if err := s.repo.SaveRefreshToken(ctx, user.ID, hashToken(refreshToken), expiresAt); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	_ = s.repo.UpdateLastLogin(ctx, user.ID)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes in seconds
	}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}
