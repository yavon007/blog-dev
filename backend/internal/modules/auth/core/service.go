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

// GuardInterface abstracts login failure tracking
type GuardInterface interface {
	RegisterFailure(ctx context.Context, email, ip string) (*GuardState, error)
	GetState(ctx context.Context, email, ip string) (*GuardState, error)
	RequireCaptcha(state *GuardState) bool
	Reset(ctx context.Context, email, ip string) error
}

// CaptchaInterface abstracts captcha operations
type CaptchaInterface interface {
	Issue(ctx context.Context, email, ip string) (*CaptchaChallenge, error)
	Validate(ctx context.Context, id, code, email, ip string) error
}

// GuardState represents login failure counts
type GuardState struct {
	IPCount    int
	EmailCount int
}

// CaptchaChallenge represents a generated captcha
type CaptchaChallenge struct {
	ID        string `json:"id"`
	Image     string `json:"image"`
	ExpiresIn int    `json:"expires_in"`
}

// LoginResult wraps the possible outcomes of a login attempt
type LoginResult struct {
	Success       bool
	TokenPair     *TokenPair
	CaptchaReq    bool
	GuardState    *GuardState
}

type Service struct {
	repo       Repository
	jwtMgr     *auth.Manager
	guard      GuardInterface
	captchaSvc CaptchaInterface
}

func NewService(repo Repository, jwtMgr *auth.Manager, guard GuardInterface, captchaSvc CaptchaInterface) *Service {
	return &Service{
		repo:      repo,
		jwtMgr:    jwtMgr,
		guard:     guard,
		captchaSvc: captchaSvc,
	}
}

// Login handles authentication with captcha support.
func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResult, error) {
	ip := req.ClientIP

	// Check current guard state
	state, err := s.guard.GetState(ctx, req.Email, ip)
	if err != nil {
		return nil, fmt.Errorf("get guard state: %w", err)
	}

	// If captcha is required, validate it
	if s.guard.RequireCaptcha(state) {
		if req.CaptchaID == "" || req.CaptchaCode == "" {
			// Captcha required but not provided
			return &LoginResult{
				Success:    false,
				CaptchaReq: true,
				GuardState: state,
			}, nil
		}
		// Validate captcha
		if err := s.captchaSvc.Validate(ctx, req.CaptchaID, req.CaptchaCode, req.Email, ip); err != nil {
			return nil, sherrors.ErrCaptchaInvalid
		}
	}

	// Verify credentials
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		// Register failure
		_, _ = s.guard.RegisterFailure(ctx, req.Email, ip)
		return nil, sherrors.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// Register failure
		_, _ = s.guard.RegisterFailure(ctx, req.Email, ip)
		return nil, sherrors.ErrUnauthorized
	}

	// Success - reset failure counters
	_ = s.guard.Reset(ctx, req.Email, ip)

	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Success:   true,
		TokenPair: tokens,
	}, nil
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

// IssueCaptcha generates a new captcha for the given email and IP
func (s *Service) IssueCaptcha(ctx context.Context, email, ip string) (*CaptchaChallenge, error) {
	return s.captchaSvc.Issue(ctx, email, ip)
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
