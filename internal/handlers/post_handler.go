package handlers

import (
	"context"
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

	h.logger.Debugw("filter", "filter", filter)

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

func (h *PostHandler) GetAllByUserID(c *gin.Context) {

}

func (h *PostHandler) GetByID(c *gin.Context) {

}

func (h *PostHandler) Update(c *gin.Context) {

}

func (h *PostHandler) Delete(c *gin.Context) {

}
