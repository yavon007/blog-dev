package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourblog/backend/internal/modules/taxonomy/core"
	sherrors "github.com/yourblog/backend/internal/modules/shared/errors"
	"github.com/yourblog/backend/internal/pkg/response"
)

type Handler struct {
	svc *core.Service
}

func NewHandler(svc *core.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterPublic(rg *gin.RouterGroup) {
	rg.GET("/categories", h.ListCategories)
	rg.GET("/tags", h.ListTags)
}

func (h *Handler) RegisterAdmin(rg *gin.RouterGroup) {
	cats := rg.Group("/categories")
	cats.GET("", h.ListCategories)
	cats.POST("", h.CreateCategory)
	cats.PUT("/:id", h.UpdateCategory)
	cats.DELETE("/:id", h.DeleteCategory)

	tags := rg.Group("/tags")
	tags.GET("", h.ListTags)
	tags.POST("", h.CreateTag)
	tags.PUT("/:id", h.UpdateTag)
	tags.DELETE("/:id", h.DeleteTag)
}

func (h *Handler) ListCategories(c *gin.Context) {
	cats, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		response.InternalError(c)
		return
	}
	response.OK(c, cats)
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var req core.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	cat, err := h.svc.CreateCategory(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, sherrors.ErrConflict) {
			response.Conflict(c, "category name or slug already exists")
			return
		}
		response.InternalError(c)
		return
	}
	response.Created(c, cat)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req core.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	cat, err := h.svc.UpdateCategory(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, sherrors.ErrNotFound) {
			response.NotFound(c, "category not found")
			return
		}
		response.InternalError(c)
		return
	}
	response.OK(c, cat)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.DeleteCategory(c.Request.Context(), id); err != nil {
		if errors.Is(err, sherrors.ErrNotFound) {
			response.NotFound(c, "category not found")
			return
		}
		response.InternalError(c)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) ListTags(c *gin.Context) {
	tags, err := h.svc.ListTags(c.Request.Context())
	if err != nil {
		response.InternalError(c)
		return
	}
	response.OK(c, tags)
}

func (h *Handler) CreateTag(c *gin.Context) {
	var req core.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	tag, err := h.svc.CreateTag(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, sherrors.ErrConflict) {
			response.Conflict(c, "tag already exists")
			return
		}
		response.InternalError(c)
		return
	}
	response.Created(c, tag)
}

func (h *Handler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req core.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	tag, err := h.svc.UpdateTag(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, sherrors.ErrNotFound) {
			response.NotFound(c, "tag not found")
			return
		}
		response.InternalError(c)
		return
	}
	response.OK(c, tag)
}

func (h *Handler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.DeleteTag(c.Request.Context(), id); err != nil {
		if errors.Is(err, sherrors.ErrNotFound) {
			response.NotFound(c, "tag not found")
			return
		}
		response.InternalError(c)
		return
	}
	response.OK(c, nil)
}
