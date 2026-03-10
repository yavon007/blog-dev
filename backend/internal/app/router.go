package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authhttp "github.com/yavon007/blog-dev/backend/internal/modules/auth/transport/http"
	commentshttp "github.com/yavon007/blog-dev/backend/internal/modules/comments/transport/http"
	feedhttp "github.com/yavon007/blog-dev/backend/internal/modules/feed/transport/http"
	mediahttp "github.com/yavon007/blog-dev/backend/internal/modules/media/transport/http"
	postshttp "github.com/yavon007/blog-dev/backend/internal/modules/posts/transport/http"
	seohttp "github.com/yavon007/blog-dev/backend/internal/modules/seo/transport/http"
	taxonomyhttp "github.com/yavon007/blog-dev/backend/internal/modules/taxonomy/transport/http"
	"github.com/yavon007/blog-dev/backend/internal/pkg/middleware"
	"github.com/yavon007/blog-dev/backend/internal/platform/auth"
	"go.uber.org/zap"
)

type Handlers struct {
	Auth     *authhttp.Handler
	Posts    *postshttp.Handler
	Taxonomy *taxonomyhttp.Handler
	Comments *commentshttp.Handler
	Media    *mediahttp.Handler
	SEO      *seohttp.Handler
	Feed     *feedhttp.Handler
}

func NewRouter(
	log *zap.Logger,
	jwtMgr *auth.Manager,
	allowedOrigins string,
	h Handlers,
) *gin.Engine {
	r := gin.New()

	// 配置可信代理，正确获取客户端 IP
	// 信任所有私有网络代理（Docker 内部网络）
	r.SetTrustedProxies([]string{
		"10.0.0.0/8",     // Docker 默认网络
		"172.16.0.0/12",  // Docker 默认网络
		"192.168.0.0/16", // 私有网络
		"127.0.0.1",      // 本地
		"::1",            // 本地 IPv6
	})

	// Global middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS(allowedOrigins))

	// Health check
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Static files for uploads
	r.Static("/uploads", "./uploads")
	h.SEO.RegisterSystem(r)
	h.Feed.RegisterSystem(r)

	// Public API
	public := r.Group("/api/v1")
	{
		h.Posts.RegisterPublic(public)
		h.Taxonomy.RegisterPublic(public)
		h.Comments.RegisterPublic(public)
		h.Auth.RegisterPublic(public) // login, refresh
		h.SEO.RegisterPublic(public)
	}

	// Admin API (JWT protected)
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.JWT(jwtMgr))
	{
		h.Auth.RegisterAdmin(admin) // logout
		h.Posts.RegisterAdmin(admin)
		h.Taxonomy.RegisterAdmin(admin)
		h.Comments.RegisterAdmin(admin)
		h.Media.RegisterAdmin(admin)
		h.SEO.RegisterAdmin(admin)
	}

	return r
}
