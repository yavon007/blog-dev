package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yavon007/blog-dev/backend/internal/modules/seo/core"
	"github.com/yavon007/blog-dev/backend/internal/pkg/response"
	"go.uber.org/zap"
)

type Handler struct {
	svc *core.Service
	log *zap.Logger
}

func NewHandler(svc *core.Service, log *zap.Logger) *Handler {
	return &Handler{
		svc: svc,
		log: log,
	}
}

func (h *Handler) RegisterSystem(r *gin.Engine) {
	r.GET("/sitemap.xml", h.Sitemap)
}

func (h *Handler) RegisterPublic(rg *gin.RouterGroup) {
	rg.GET("/seo/meta", h.GetMeta)
}

func (h *Handler) RegisterAdmin(rg *gin.RouterGroup) {
	rg.PUT("/seo/meta", h.UpdateMeta)
}

func (h *Handler) Sitemap(c *gin.Context) {
	xmlBody, err := h.svc.GenerateSitemap(c.Request.Context())
	if err != nil {
		h.log.Error("generate sitemap failed", zap.Error(err))
		response.InternalError(c)
		return
	}
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, xmlBody)
}

func (h *Handler) GetMeta(c *gin.Context) {
	settings, err := h.svc.GetMeta(c.Request.Context())
	if err != nil {
		h.log.Error("get site meta failed", zap.Error(err))
		response.InternalError(c)
		return
	}
	response.OK(c, settings)
}

func (h *Handler) UpdateMeta(c *gin.Context) {
	var req core.UpdateMetaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	settings, err := h.svc.UpdateMeta(c.Request.Context(), req)
	if err != nil {
		h.log.Error("update site meta failed", zap.Error(err))
		response.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    settings,
		"message": "success",
	})
}
