package core

import "time"

type AdminUser struct {
	ID           int64
	Email        string
	PasswordHash string
	Role         string
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type LoginRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	CaptchaID   string `json:"captcha_id,omitempty"`
	CaptchaCode string `json:"captcha_code,omitempty"`
	ClientIP    string `json:"-"` // populated by handler
}

// LoginResponse wraps TokenPair with optional captcha requirement info.
type LoginResponse struct {
	*TokenPair
	CaptchaRequired bool              `json:"captcha_required,omitempty"`
	Failures        map[string]int    `json:"failures,omitempty"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds
}
