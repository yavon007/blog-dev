package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yavon007/blog-dev/backend/internal/app"
	"github.com/yavon007/blog-dev/backend/internal/config"
	"github.com/yavon007/blog-dev/backend/internal/modules/auth/core"
	authrepo "github.com/yavon007/blog-dev/backend/internal/modules/auth/repository"
	"github.com/yavon007/blog-dev/backend/internal/modules/auth/security"
	authhttp "github.com/yavon007/blog-dev/backend/internal/modules/auth/transport/http"
	commentscore "github.com/yavon007/blog-dev/backend/internal/modules/comments/core"
	commentsrepo "github.com/yavon007/blog-dev/backend/internal/modules/comments/repository"
	commentshttp "github.com/yavon007/blog-dev/backend/internal/modules/comments/transport/http"
	feedcore "github.com/yavon007/blog-dev/backend/internal/modules/feed/core"
	feedhttp "github.com/yavon007/blog-dev/backend/internal/modules/feed/transport/http"
	mediacore "github.com/yavon007/blog-dev/backend/internal/modules/media/core"
	mediarepo "github.com/yavon007/blog-dev/backend/internal/modules/media/repository"
	mediahttp "github.com/yavon007/blog-dev/backend/internal/modules/media/transport/http"
	"github.com/yavon007/blog-dev/backend/internal/modules/media/storage"
	postscore "github.com/yavon007/blog-dev/backend/internal/modules/posts/core"
	postsrepo "github.com/yavon007/blog-dev/backend/internal/modules/posts/repository"
	postshttp "github.com/yavon007/blog-dev/backend/internal/modules/posts/transport/http"
	seocore "github.com/yavon007/blog-dev/backend/internal/modules/seo/core"
	seorepo "github.com/yavon007/blog-dev/backend/internal/modules/seo/repository"
	seohttp "github.com/yavon007/blog-dev/backend/internal/modules/seo/transport/http"
	taxonomycore "github.com/yavon007/blog-dev/backend/internal/modules/taxonomy/core"
	taxonomyrepo "github.com/yavon007/blog-dev/backend/internal/modules/taxonomy/repository"
	taxonomyhttp "github.com/yavon007/blog-dev/backend/internal/modules/taxonomy/transport/http"
	"github.com/yavon007/blog-dev/backend/internal/platform/auth"
	"github.com/yavon007/blog-dev/backend/internal/platform/database"
	"github.com/yavon007/blog-dev/backend/internal/platform/logger"
	"go.uber.org/zap"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	// Logger
	log, err := logger.New(cfg.App.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Database
	ctx := context.Background()
	db, err := database.NewPool(ctx, cfg.Database)
	if err != nil {
		log.Fatal("connect database", zap.Error(err))
	}
	defer db.Close()

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("connect redis", zap.Error(err))
	}
	defer rdb.Close()

	// JWT Manager
	jwtMgr := auth.NewManager(cfg.JWT)

	// Security services
	loginGuard := security.NewLoginGuard(rdb, cfg.Security)
	captchaSvc := security.NewCaptchaService(rdb, cfg.Security)

	// Wire modules: auth
	authRepo := authrepo.NewPostgresRepo(db)
	authSvc := core.NewService(authRepo, jwtMgr, loginGuard, captchaSvc)
	authHandler := authhttp.NewHandler(authSvc, captchaSvc)

	// Wire modules: posts
	postsRepo := postsrepo.NewPostgresRepo(db)
	postsSvc := postscore.NewService(postsRepo)
	postsHandler := postshttp.NewHandler(postsSvc, log)

	// Wire modules: taxonomy
	taxonomyRepo := taxonomyrepo.NewPostgresRepo(db)
	taxonomySvc := taxonomycore.NewService(taxonomyRepo)
	taxonomyHandler := taxonomyhttp.NewHandler(taxonomySvc)

	// Wire modules: comments
	commentsRepo := commentsrepo.NewPostgresRepo(db)
	commentsSvc := commentscore.NewService(commentsRepo)
	commentsHandler := commentshttp.NewHandler(commentsSvc)

	// Wire modules: media
	mediaRepo := mediarepo.NewPostgresRepo(db)
	mediaStorage, err := storage.NewLocalStorage("./uploads", cfg.App.BaseURL)
	if err != nil {
		log.Fatal("init media storage", zap.Error(err))
	}
	mediaSvc := mediacore.NewService(mediaRepo, mediaStorage, cfg.App.BaseURL)
	mediaHandler := mediahttp.NewHandler(mediaSvc, log)

	// Wire modules: seo
	seoRepo := seorepo.NewPostgresRepo(db)
	seoSvc := seocore.NewService(seoRepo, cfg.App.BaseURL)
	seoHandler := seohttp.NewHandler(seoSvc, log)

	// Wire modules: feed (reuse SEO repository for posts/settings)
	feedSvc := feedcore.NewService(seoRepo, cfg.App.BaseURL)
	feedHandler := feedhttp.NewHandler(feedSvc, log)

	// Router
	router := app.NewRouter(log, jwtMgr, cfg.App.AllowedOrigins, app.Handlers{
		Auth:     authHandler,
		Posts:    postsHandler,
		Taxonomy: taxonomyHandler,
		Comments: commentsHandler,
		Media:    mediaHandler,
		SEO:      seoHandler,
		Feed:     feedHandler,
	})

	// HTTP Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Info("starting server", zap.String("addr", srv.Addr), zap.String("env", cfg.App.Env))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", zap.Error(err))
	}
	log.Info("server exited")
}
