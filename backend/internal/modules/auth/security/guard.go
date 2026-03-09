package security

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yavon007/blog-dev/backend/internal/config"
	"github.com/yavon007/blog-dev/backend/internal/modules/auth/core"
)

// LoginGuard tracks login failures per IP and email dimensions.
type LoginGuard struct {
	rdb       *redis.Client
	threshold int
	ttl       int64 // seconds
}

func NewLoginGuard(rdb *redis.Client, cfg config.SecurityConfig) *LoginGuard {
	return &LoginGuard{
		rdb:       rdb,
		threshold: cfg.LoginFailThreshold,
		ttl:       int64(cfg.LoginFailTTL.Seconds()),
	}
}

// key patterns
func ipKey(ip string) string    { return "auth:login:fail:ip:" + ip }
func emailKey(email string) string { return "auth:login:fail:email:" + email }

// RegisterFailure increments failure counters for both IP and email.
// Returns the updated state after incrementing.
func (g *LoginGuard) RegisterFailure(ctx context.Context, email, ip string) (*core.GuardState, error) {
	ipCount, err := g.increment(ctx, ipKey(ip))
	if err != nil {
		return nil, fmt.Errorf("increment ip counter: %w", err)
	}
	emailCount, err := g.increment(ctx, emailKey(email))
	if err != nil {
		return nil, fmt.Errorf("increment email counter: %w", err)
	}
	return &core.GuardState{IPCount: ipCount, EmailCount: emailCount}, nil
}

// GetState returns current failure counts without incrementing.
func (g *LoginGuard) GetState(ctx context.Context, email, ip string) (*core.GuardState, error) {
	ipCount, err := g.getCount(ctx, ipKey(ip))
	if err != nil {
		return nil, fmt.Errorf("get ip counter: %w", err)
	}
	emailCount, err := g.getCount(ctx, emailKey(email))
	if err != nil {
		return nil, fmt.Errorf("get email counter: %w", err)
	}
	return &core.GuardState{IPCount: ipCount, EmailCount: emailCount}, nil
}

// RequireCaptcha returns true if either dimension has reached the threshold.
func (g *LoginGuard) RequireCaptcha(state *core.GuardState) bool {
	return state.IPCount >= g.threshold || state.EmailCount >= g.threshold
}

// Reset clears all failure counters for the given email and IP.
func (g *LoginGuard) Reset(ctx context.Context, email, ip string) error {
	if err := g.rdb.Del(ctx, ipKey(ip)).Err(); err != nil {
		return fmt.Errorf("delete ip counter: %w", err)
	}
	if err := g.rdb.Del(ctx, emailKey(email)).Err(); err != nil {
		return fmt.Errorf("delete email counter: %w", err)
	}
	return nil
}

func (g *LoginGuard) increment(ctx context.Context, key string) (int, error) {
	val, err := g.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	// Set TTL on first increment
	if val == 1 {
		if err := g.rdb.Expire(ctx, key, time.Duration(g.ttl)*time.Second).Err(); err != nil {
			return 0, err
		}
	}
	return int(val), nil
}

func (g *LoginGuard) getCount(ctx context.Context, key string) (int, error) {
	val, err := g.rdb.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return val, nil
}
