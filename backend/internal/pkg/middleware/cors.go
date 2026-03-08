package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigins string) gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = splitOrigins(allowedOrigins)
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"}
	config.ExposeHeaders = []string{"X-Request-ID"}
	config.AllowCredentials = true
	return cors.New(config)
}

func splitOrigins(s string) []string {
	if s == "" {
		return []string{"*"}
	}
	var origins []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			o := s[start:i]
			if len(o) > 0 {
				origins = append(origins, o)
			}
			start = i + 1
		}
	}
	return origins
}
