package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/wanafiq/feed-api/internal/middleware"
	"github.com/wanafiq/feed-api/internal/response"
	"github.com/wanafiq/feed-api/internal/services"
	"github.com/wanafiq/feed-api/internal/types"
	"go.uber.org/zap"
)

type PostHandler struct {
	logger      *zap.SugaredLogger
	postService *services.PostService
}

func NewPostHandler(logger *zap.SugaredLogger, postService *services.PostService) *PostHandler {
	return &PostHandler{
		logger:      logger,
		postService: postService,
	}
}

func (h *PostHandler) Save(c *gin.Context) {
	userCtx, exists := middleware.GetUserContext(c)
	if !exists {
		response.Unauthorized(c, nil)
		return
	}

	authorID := userCtx.ID

	var req types.SavePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	post, err := h.postService.Save(context.Background(), authorID, &req)
	if err != nil {
		response.InternalServerError(c)
		return
	}

	response.Created(c, post)
}

func (h *PostHandler) GetAll(c *gin.Context) {
	posts, err := h.postService.GetAll(context.Background())
	if err != nil {
		response.InternalServerError(c)
		return
	}

	response.OK(c, posts)
}

func (h *PostHandler) GetAllByUserID(c *gin.Context) {

}

func (h *PostHandler) GetByID(c *gin.Context) {

}

func (h *PostHandler) Update(c *gin.Context) {

}

func (h *PostHandler) Delete(c *gin.Context) {

}
