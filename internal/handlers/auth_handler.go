package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wanafiq/feed-api/internal/constants"
	"github.com/wanafiq/feed-api/internal/response"
	"github.com/wanafiq/feed-api/internal/services"
	"github.com/wanafiq/feed-api/internal/types"
	"go.uber.org/zap"
)

type AuthHandler struct {
	logger      *zap.SugaredLogger
	authService *services.AuthService
}

func NewAuthHandler(logger *zap.SugaredLogger, authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req types.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	createdUser, err := h.authService.Register(context.Background(), &req)
	if err != nil {
		response.InternalServerError(c)
		return
	}

	response.Created(c, createdUser)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	token, err := h.authService.Login(context.Background(), &req)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrUnauthorized):
			response.Unauthorized(c, nil)
		default:
			response.InternalServerError(c)
		}
		return
	}

	response.OK(c, token, nil)
}
