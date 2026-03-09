package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yavon007/blog-dev/backend/internal/modules/media/core"
	sherrors "github.com/yavon007/blog-dev/backend/internal/modules/shared/errors"
	"github.com/yavon007/blog-dev/backend/internal/pkg/middleware"
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

func (h *Handler) RegisterAdmin(rg *gin.RouterGroup) {
	media := rg.Group("/media")
	media.GET("", h.List)
	media.POST("", h.Upload)
	media.GET("/:id", h.GetByID)
	media.DELETE("/:id", h.Delete)
}

func (h *Handler) List(c *gin.Context) {
	filter := core.ListMediaFilter{
		MimeType: c.Query("mime_type"),
		Page:     1,
		PageSize: 20,
	}
	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		filter.Page = page
	}
	if size, err := strconv.Atoi(c.Query("page_size")); err == nil && size > 0 && size <= 100 {
		filter.PageSize = size
	}

	files, total, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		h.log.Error("list media failed", zap.Error(err))
		response.InternalError(c)
		return
	}

	response.Paged(c, files, total, filter.Page, filter.PageSize)
}

func (h *Handler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	altText := c.PostForm("alt_text")

	claims := middleware.GetClaims(c)
	if claims == nil {
		response.Unauthorized(c)
		return
	}

	result, err := h.svc.Upload(c.Request.Context(), file, altText, claims.AdminID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "data": result, "message": "uploaded"})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	file, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sherrors.ErrNotFound) {
			response.NotFound(c, "media not found")
			return
		}
		response.InternalError(c)
		return
	}

	response.OK(c, file)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, sherrors.ErrNotFound) {
			response.NotFound(c, "media not found")
			return
		}
		response.InternalError(c)
		return
	}

	response.OK(c, nil)
}
