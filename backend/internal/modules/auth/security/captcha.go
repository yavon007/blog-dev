package security

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yavon007/blog-dev/backend/internal/config"
	"github.com/yavon007/blog-dev/backend/internal/modules/auth/core"
)

const captchaCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// CaptchaService handles captcha generation and validation.
type CaptchaService struct {
	rdb *redis.Client
	ttl int64 // seconds
}

func NewCaptchaService(rdb *redis.Client, cfg config.SecurityConfig) *CaptchaService {
	return &CaptchaService{
		rdb: rdb,
		ttl: int64(cfg.CaptchaTTL.Seconds()),
	}
}

// Issue generates a new captcha challenge for the given email and IP.
func (s *CaptchaService) Issue(ctx context.Context, email, ip string) (*core.CaptchaChallenge, error) {
	// Generate random 6-character code
	code := generateCode(6)

	// Generate ID
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return nil, fmt.Errorf("generate captcha id: %w", err)
	}
	id := hex.EncodeToString(idBytes)

	// Render image
	imageData, err := renderCaptchaImage(code)
	if err != nil {
		return nil, fmt.Errorf("render captcha image: %w", err)
	}

	// Store hashed answer in Redis
	hash := hashAnswer(code)
	meta := captchaMeta{
		Hash:  hash,
		Email: email,
		IP:    ip,
	}
	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return nil, fmt.Errorf("marshal captcha meta: %w", err)
	}

	key := captchaKey(id)
	if err := s.rdb.Set(ctx, key, metaBytes, time.Duration(s.ttl)*time.Second).Err(); err != nil {
		return nil, fmt.Errorf("store captcha: %w", err)
	}

	return &core.CaptchaChallenge{
		ID:        id,
		Image:     imageData,
		ExpiresIn: int(s.ttl),
	}, nil
}

// Validate checks if the provided code matches the stored captcha.
// Returns nil on success, or an error if invalid/expired.
// On any validation attempt (success or fail), the captcha is consumed.
func (s *CaptchaService) Validate(ctx context.Context, id, code, email, ip string) error {
	key := captchaKey(id)

	// Get stored meta
	metaBytes, err := s.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return ErrCaptchaInvalid
	}
	if err != nil {
		return fmt.Errorf("get captcha: %w", err)
	}

	// Always delete after retrieval (single-use)
	defer s.rdb.Del(ctx, key)

	var meta captchaMeta
	if err := json.Unmarshal(metaBytes, &meta); err != nil {
		return fmt.Errorf("unmarshal captcha meta: %w", err)
	}

	// Verify the captcha belongs to this login attempt
	if meta.Email != email || meta.IP != ip {
		return ErrCaptchaInvalid
	}

	// Verify answer
	if hashAnswer(code) != meta.Hash {
		return ErrCaptchaInvalid
	}

	return nil
}

type captchaMeta struct {
	Hash  string `json:"hash"`
	Email string `json:"email"`
	IP    string `json:"ip"`
}

func captchaKey(id string) string {
	return "auth:captcha:" + id
}

func generateCode(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = captchaCharset[int(b[i])%len(captchaCharset)]
	}
	return string(b)
}

func hashAnswer(code string) string {
	h := sha256.Sum256([]byte(code))
	return hex.EncodeToString(h[:])
}

// Error definitions
var ErrCaptchaInvalid = fmt.Errorf("captcha invalid or expired")
