package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yavon007/blog-dev/backend/internal/pkg/response"
	"github.com/yavon007/blog-dev/backend/internal/platform/auth"
)

const AdminClaimsKey = "admin_claims"

func JWT(jwtMgr *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtMgr.ParseAccessToken(tokenStr)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set(AdminClaimsKey, claims)
		c.Next()
	}
}

func GetClaims(c *gin.Context) *auth.Claims {
	if v, exists := c.Get(AdminClaimsKey); exists {
		if claims, ok := v.(*auth.Claims); ok {
			return claims
		}
	}
	return nil
}
