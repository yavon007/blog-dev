package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourblog/backend/internal/modules/posts/core"
	sherrors "github.com/yourblog/backend/internal/modules/shared/errors"
	"github.com/yourblog/backend/internal/pkg/middleware"
	"github.com/yourblog/backend/internal/pkg/pagination"
	"github.com/yourblog/backend/internal/pkg/response"
)

type Handler struct {
	svc *core.Service
}

func NewHandler(svc *core.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterPublic(rg *gin.RouterGroup) {
	rg.GET("/posts", h.ListPublic)
	rg.GET("/posts/:slug", h.GetBySlug)
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
	post, err := h.svc.Create(c.Request.Context(), req, claims.AdminID)
	if err != nil {
		if errors.Is(err, sherrors.ErrConflict) {
			response.Conflict(c, "slug already exists")
			return
		}
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
