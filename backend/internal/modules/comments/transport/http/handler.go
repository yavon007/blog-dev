package http

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourblog/backend/internal/modules/comments/core"
	sherrors "github.com/yourblog/backend/internal/modules/shared/errors"
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
	// Mounted under /posts/:slug by post handler — register separately
	rg.GET("/posts/:slug/comments", h.ListPublic)
	rg.POST("/posts/:slug/comments", h.CreateComment)
}

func (h *Handler) RegisterAdmin(rg *gin.RouterGroup) {
	comments := rg.Group("/comments")
	comments.GET("", h.ListAdmin)
	comments.PATCH("/:id", h.UpdateStatus)
	comments.DELETE("/:id", h.Delete)
}

func (h *Handler) ListPublic(c *gin.Context) {
	// Note: in real impl, resolve slug -> postID first; simplified here
	postID, _ := strconv.ParseInt(c.Query("post_id"), 10, 64)
	p := pagination.FromQuery(c)
	comments, total, err := h.svc.ListPublic(c.Request.Context(), postID, p)
	if err != nil {
		response.InternalError(c)
		return
	}
	response.Paged(c, comments, total, p.Page, p.Size)
}

func (h *Handler) CreateComment(c *gin.Context) {
	postID, _ := strconv.ParseInt(c.Query("post_id"), 10, 64)
	var req core.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ipHash := hashIP(c.ClientIP())
	comment, err := h.svc.Create(c.Request.Context(), postID, req, ipHash, c.Request.UserAgent())
	if err != nil {
		response.InternalError(c)
		return
	}
	response.Created(c, comment)
}

func (h *Handler) ListAdmin(c *gin.Context) {
	p := pagination.FromQuery(c)
	postID, _ := strconv.ParseInt(c.Query("post_id"), 10, 64)
	filter := core.ListCommentsFilter{
		PostID: postID,
		Status: core.CommentStatus(c.Query("status")),
	}
	comments, total, err := h.svc.ListAdmin(c.Request.Context(), filter, p)
	if err != nil {
		response.InternalError(c)
		return
	}
	response.Paged(c, comments, total, p.Page, p.Size)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req core.UpdateCommentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateStatus(c.Request.Context(), id, req); err != nil {
		if errors.Is(err, sherrors.ErrNotFound) {
			response.NotFound(c, "comment not found")
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
			response.NotFound(c, "comment not found")
			return
		}
		response.InternalError(c)
		return
	}
	response.OK(c, nil)
}

func hashIP(ip string) string {
	h := sha256.Sum256([]byte(ip))
	return fmt.Sprintf("%x", h)[:16]
}
