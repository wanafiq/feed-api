package handlers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wanafiq/feed-api/internal/middleware"
	"github.com/wanafiq/feed-api/internal/models"
	"github.com/wanafiq/feed-api/internal/response"
	"github.com/wanafiq/feed-api/internal/services"
	"github.com/wanafiq/feed-api/internal/types"
	"github.com/wanafiq/feed-api/internal/utils"
	"go.uber.org/zap"
	"strings"
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

	var req types.PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	post, err := h.postService.Save(context.Background(), userCtx.ID, &req)
	if err != nil {
		response.InternalServerError(c)
		return
	}

	response.Created(c, post)
}

func (h *PostHandler) GetAll(c *gin.Context) {
	offset := utils.ParseQueryInt(c, "offset", 0)

	limit := utils.ParseQueryInt(c, "limit", 10)
	if limit > 100 {
		limit = 100
	}

	search := c.DefaultQuery("search", "")

	sort := strings.ToLower(c.DefaultQuery("sort", "desc"))
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	tags := c.QueryArray("tags")
	dateFrom := utils.ParseQueryTime(c, "from")
	dateTo := utils.ParseQueryTime(c, "to")

	filter := models.PostFilter{
		Offset:   offset,
		Limit:    limit,
		Search:   search,
		Sort:     sort,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Tags:     tags,
	}

	posts, count, err := h.postService.GetAll(context.Background(), filter)
	if err != nil {
		response.InternalServerError(c)
		return
	}

	pagination := response.Pagination{
		Total:  count,
		Limit:  limit,
		Offset: offset,
		Next:   utils.Min(offset+limit, count),
		Prev:   utils.Max(offset-limit, 0),
	}

	response.OK(c, posts, &pagination)
}

func (h *PostHandler) GetByID(c *gin.Context) {
	postID := c.Param("postID")
	if postID == "" {
		response.BadRequest(c, errors.New("postID is required"))
		return
	}

	post, err := h.postService.GetPostByID(context.Background(), postID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			response.NotFound(c, nil)
		default:
			response.InternalServerError(c)
		}
		return
	}

	response.OK(c, post, nil)
}

func (h *PostHandler) Update(c *gin.Context) {
	userCtx, exists := middleware.GetUserContext(c)
	if !exists {
		response.Unauthorized(c, nil)
		return
	}

	postID := c.Param("postID")
	if postID == "" {
		response.BadRequest(c, errors.New("postID is required"))
		return
	}

	var req types.PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	post, err := h.postService.Update(context.Background(), userCtx, postID, &req)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			response.NotFound(c, nil)
		default:
			response.InternalServerError(c)
		}
	}

	response.OK(c, post, nil)
}

func (h *PostHandler) Delete(c *gin.Context) {
	postID := c.Param("postID")
	if postID == "" {
		response.BadRequest(c, errors.New("postID is required"))
		return
	}
}
