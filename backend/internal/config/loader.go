package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	RateLimit RateLimitConfig
}

type AppConfig struct {
	Env             string
	Name            string
	Port            int
	BaseURL         string
	AllowedOrigins  string
	LogLevel        string
}

type DatabaseConfig struct {
	URL         string
	MaxOpen     int
	MaxIdle     int
	MaxLifetime time.Duration
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type RateLimitConfig struct {
	PublicRPS int
	AdminRPS  int
}

func Load() (*Config, error) {
	// Load .env file if exists (ignore error in production)
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Env:            getEnv("APP_ENV", "development"),
			Name:           getEnv("APP_NAME", "blog-backend"),
			Port:           getEnvInt("APP_PORT", 8080),
			BaseURL:        getEnv("APP_BASE_URL", "http://localhost:8080"),
			AllowedOrigins: getEnv("PUBLIC_ALLOWED_ORIGINS", "http://localhost:5173"),
			LogLevel:       getEnv("LOG_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			URL:         getEnvRequired("DATABASE_URL"),
			MaxOpen:     getEnvInt("DATABASE_MAX_OPEN", 50),
			MaxIdle:     getEnvInt("DATABASE_MAX_IDLE", 10),
			MaxLifetime: getEnvDuration("DATABASE_MAX_LIFETIME", 30*time.Minute),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     getEnvRequired("JWT_SECRET"),
			AccessTTL:  getEnvDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTTL: getEnvDuration("JWT_REFRESH_TTL", 720*time.Hour),
		},
		RateLimit: RateLimitConfig{
			PublicRPS: getEnvInt("RATE_LIMIT_PUBLIC_RPS", 20),
			AdminRPS:  getEnvInt("RATE_LIMIT_ADMIN_RPS", 5),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvRequired(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return v
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultVal
}
