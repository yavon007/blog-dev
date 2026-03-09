package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yavon007/blog-dev/backend/internal/modules/posts/core"
	sherrors "github.com/yavon007/blog-dev/backend/internal/modules/shared/errors"
	"github.com/yavon007/blog-dev/backend/internal/pkg/middleware"
	"github.com/yavon007/blog-dev/backend/internal/pkg/pagination"
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

func (h *Handler) RegisterPublic(rg *gin.RouterGroup) {
	rg.GET("/posts", h.ListPublic)
	rg.GET("/posts/:slug", h.GetBySlug)
	rg.GET("/posts/archive", h.GetArchive)
	rg.GET("/posts/archive/:year/:month", h.ListByYearMonth)
}

func (h *Handler) RegisterAdmin(rg *gin.RouterGroup) {
	posts := rg.Group("/posts")
	posts.GET("", h.ListAdmin)
	posts.POST("", h.Create)
	posts.GET("/:id", h.GetByID)
	posts.PUT("/:id", h.Update)
	posts.PATCH("/:id/status", h.UpdateStatus)
	posts.DELETE("/:id", h.Delete)
}

func (h *Handler) ListPublic(c *gin.Context) {
	p := pagination.FromQuery(c)
	filter := core.ListFilter{
		Category: c.Query("category"),
		Tag:      c.Query("tag"),
		Query:    c.Query("q"),
	}
	posts, total, err := h.svc.ListPublic(c.Request.Context(), filter, p)
	if err != nil {
		h.log.Error("list public posts failed", zap.Error(err))
		response.InternalError(c)
		return
	}
	response.Paged(c, posts, total, p.Page, p.Size)
}

func (h *Handler) ListAdmin(c *gin.Context) {
	p := pagination.FromQuery(c)
	filter := core.ListFilter{
		Status: core.PostStatus(c.Query("status")),
		Query:  c.Query("q"),
	}
	posts, total, err := h.svc.ListAdmin(c.Request.Context(), filter, p)
	if err != nil {
		h.log.Error("list admin posts failed", zap.Error(err))
		response.InternalError(c)
		return
	}
	response.Paged(c, posts, total, p.Page, p.Size)
}

func (h *Handler) GetBySlug(c *gin.Context) {
	post, err := h.svc.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		if errors.Is(err, sherrors.ErrNotFound) {
			response.NotFound(c, "post not found")
			return
		}
		response.InternalError(c)
		return
	}
	response.OK(c, post)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	post, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sherrors.ErrNotFound) {
			response.NotFound(c, "post not found")
			return
		}
		response.InternalError(c)
		return
	}
	response.OK(c, post)
}

func (h *Handler) Create(c *gin.Context) {
	var req core.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	claims := middleware.GetClaims(c)
	if claims == nil {
		response.Unauthorized(c)
		return
	}
	post, err := h.svc.Create(c.Request.Context(), req, claims.AdminID)
	if err != nil {
		if errors.Is(err, sherrors.ErrConflict) {
			response.Conflict(c, "slug already exists")
			return
		}
		h.log.Error("create post failed", zap.Error(err))
		response.InternalError(c)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "data": post, "message": "created"})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req core.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	post, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, sherrors.ErrNotFound) {
			response.NotFound(c, "post not found")
			return
		}
		h.log.Error("update post failed", zap.Error(err))
		response.InternalError(c)
		return
	}
	response.OK(c, post)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req core.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateStatus(c.Request.Context(), id, req); err != nil {
		if errors.Is(err, sherrors.ErrNotFound) {
			response.NotFound(c, "post not found")
			return
		}
		response.InternalError(c)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, sherrors.ErrNotFound) {
			response.NotFound(c, "post not found")
			return
		}
		response.InternalError(c)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) GetArchive(c *gin.Context) {
	items, err := h.svc.GetArchive(c.Request.Context())
	if err != nil {
		h.log.Error("get archive failed", zap.Error(err))
		response.InternalError(c)
		return
	}
	response.OK(c, items)
}

func (h *Handler) ListByYearMonth(c *gin.Context) {
	year, err1 := strconv.Atoi(c.Param("year"))
	month, err2 := strconv.Atoi(c.Param("month"))
	if err1 != nil || err2 != nil || year < 2000 || month < 1 || month > 12 {
		response.BadRequest(c, "invalid year or month")
		return
	}

	p := pagination.FromQuery(c)
	posts, total, err := h.svc.ListByYearMonth(c.Request.Context(), year, month, p)
	if err != nil {
		h.log.Error("list posts by year/month failed", zap.Error(err))
		response.InternalError(c)
		return
	}
	response.Paged(c, posts, total, p.Page, p.Size)
}
