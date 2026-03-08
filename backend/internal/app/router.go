package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authhttp "github.com/yourblog/backend/internal/modules/auth/transport/http"
	commentshttp "github.com/yourblog/backend/internal/modules/comments/transport/http"
	postshttp "github.com/yourblog/backend/internal/modules/posts/transport/http"
	taxonomyhttp "github.com/yourblog/backend/internal/modules/taxonomy/transport/http"
	"github.com/yourblog/backend/internal/pkg/middleware"
	"github.com/yourblog/backend/internal/platform/auth"
	"go.uber.org/zap"
)

type Handlers struct {
	Auth     *authhttp.Handler
	Posts    *postshttp.Handler
	Taxonomy *taxonomyhttp.Handler
	Comments *commentshttp.Handler
}

func NewRouter(
	log *zap.Logger,
	jwtMgr *auth.Manager,
	allowedOrigins string,
	h Handlers,
) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS(allowedOrigins))

	// Health check
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Public API
	public := r.Group("/api/v1")
	{
		h.Posts.RegisterPublic(public)
		h.Taxonomy.RegisterPublic(public)
		h.Comments.RegisterPublic(public)
		h.Auth.RegisterPublic(public) // login, refresh
	}

	// Admin API (JWT protected)
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.JWT(jwtMgr))
	{
		h.Auth.RegisterAdmin(admin) // logout
		h.Posts.RegisterAdmin(admin)
		h.Taxonomy.RegisterAdmin(admin)
		h.Comments.RegisterAdmin(admin)
	}

	return r
}
