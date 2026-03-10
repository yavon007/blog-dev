package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yavon007/blog-dev/backend/internal/modules/feed/core"
	"github.com/yavon007/blog-dev/backend/internal/pkg/response"
	"go.uber.org/zap"
)

type Handler struct {
	svc *core.Service
	log *zap.Logger
}

func NewHandler(svc *core.Service, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

func (h *Handler) RegisterSystem(r *gin.Engine) {
	r.GET("/rss.xml", h.Feed)
}

func (h *Handler) Feed(c *gin.Context) {
	xmlBody, err := h.svc.GenerateFeed(c.Request.Context())
	if err != nil {
		h.log.Error("generate rss failed", zap.Error(err))
		response.InternalError(c)
		return
	}
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, xmlBody)
}
