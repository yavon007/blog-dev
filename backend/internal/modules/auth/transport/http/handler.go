package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/yourblog/backend/internal/modules/auth/core"
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
	auth := rg.Group("/auth")
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh)
}

func (h *Handler) RegisterAdmin(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/logout", h.Logout)
}

func (h *Handler) Login(c *gin.Context) {
	var req core.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	pair, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, sherrors.ErrUnauthorized) {
			response.Unauthorized(c)
			return
		}
		response.InternalError(c)
		return
	}

	response.OK(c, pair)
}

func (h *Handler) Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	pair, err := h.svc.Refresh(c.Request.Context(), body.RefreshToken)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	response.OK(c, pair)
}

func (h *Handler) Logout(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	_ = h.svc.Logout(c.Request.Context(), body.RefreshToken)
	response.OK(c, nil)
}
