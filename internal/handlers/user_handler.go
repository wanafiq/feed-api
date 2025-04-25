package handlers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wanafiq/feed-api/internal/middleware"
	"github.com/wanafiq/feed-api/internal/response"
	"github.com/wanafiq/feed-api/internal/services"
	"go.uber.org/zap"
)

type UserHandler struct {
	logger      *zap.SugaredLogger
	userService *services.UserService
}

func NewUserHandler(logger *zap.SugaredLogger, userService *services.UserService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

func (h *UserHandler) GetByID(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		response.BadRequest(c, errors.New("userID is required"))
		return
	}

	user, err := h.userService.GetByID(c, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			response.NotFound(c, nil)
		default:
			response.InternalServerError(c)
		}
		return
	}

	response.OK(c, user, nil)
}

func (h *UserHandler) Deactivate(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		response.BadRequest(c, errors.New("userID is required"))
		return
	}

	user, err := h.userService.GetByID(c, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			response.NotFound(c, nil)
		default:
			response.InternalServerError(c)
		}
		return
	}

	response.OK(c, user, nil)
}

func (h *UserHandler) Follow(c *gin.Context) {
	userCtx, exists := middleware.GetUserContext(c)
	if !exists {
		response.Unauthorized(c, nil)
		return
	}

	followerID := userCtx.ID

	followeeID := c.Param("userID")
	if followeeID == "" {
		response.BadRequest(c, errors.New("userID is required"))
		return
	}

	err := h.userService.Follow(context.Background(), followerID, followeeID)
	if err != nil {
		response.InternalServerError(c)
		return
	}

	response.NoContent(c)
}

func (h *UserHandler) Unfollow(c *gin.Context) {
	userCtx, exists := middleware.GetUserContext(c)
	if !exists {
		response.Unauthorized(c, nil)
		return
	}

	followerID := userCtx.ID

	followeeID := c.Param("userID")
	if followeeID == "" {
		response.BadRequest(c, errors.New("userID is required"))
		return
	}

	err := h.userService.Unfollow(context.Background(), followerID, followeeID)
	if err != nil {
		response.InternalServerError(c)
		return
	}

	response.NoContent(c)
}
