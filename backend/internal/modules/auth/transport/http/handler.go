package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yavon007/blog-dev/backend/internal/modules/auth/core"
	sherrors "github.com/yavon007/blog-dev/backend/internal/modules/shared/errors"
	"github.com/yavon007/blog-dev/backend/internal/pkg/response"
)

type Handler struct {
	svc        *core.Service
	captchaSvc core.CaptchaInterface
}

func NewHandler(svc *core.Service, captchaSvc core.CaptchaInterface) *Handler {
	return &Handler{svc: svc, captchaSvc: captchaSvc}
}

func (h *Handler) RegisterPublic(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh)
	auth.GET("/captcha", h.GetCaptcha)
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

	// Extract client IP
	req.ClientIP = c.ClientIP()

	result, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, sherrors.ErrUnauthorized) {
			response.Error(c, http.StatusUnauthorized, "invalid credentials")
			return
		}
		if errors.Is(err, sherrors.ErrCaptchaInvalid) {
			response.Error(c, http.StatusBadRequest, "invalid or expired captcha")
			return
		}
		response.InternalError(c)
		return
	}

	// Check if captcha is required
	if !result.Success && result.CaptchaReq {
		response.ErrorWithData(c, http.StatusPreconditionRequired, "captcha required", gin.H{
			"captcha_required": true,
			"failures": gin.H{
				"ip":    result.GuardState.IPCount,
				"email": result.GuardState.EmailCount,
			},
		})
		return
	}

	response.OK(c, result.TokenPair)
}

func (h *Handler) GetCaptcha(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		response.BadRequest(c, "email is required")
		return
	}

	ip := c.ClientIP()
	challenge, err := h.svc.IssueCaptcha(c.Request.Context(), email, ip)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Created(c, challenge)
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

	if err := h.svc.Logout(c.Request.Context(), body.RefreshToken); err != nil {
		response.InternalError(c)
		return
	}

	response.OK(c, nil)
}

