package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Pagination struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Next   int `json:"next"`
	Prev   int `json:"prev"`
}

type response struct {
	Status     int         `json:"status"`
	Message    string      `json:"message,omitzero"`
	Data       any         `json:"data,omitzero"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Success
func OK(c *gin.Context, data any, pagination *Pagination) {
	if pagination != nil {
		paginatedResponse(c, http.StatusOK, data, pagination)
		return
	}
	successResponse(c, http.StatusOK, data)
}

func Created(c *gin.Context, data any) {
	successResponse(c, http.StatusCreated, data)
}

func NoContent(c *gin.Context) {
	successResponse(c, http.StatusNoContent, nil)
}

// Errors
func Unauthorized(c *gin.Context, error error) {
	errorResponse(c, http.StatusUnauthorized, error)
}

func Forbidden(c *gin.Context, error error) {
	errorResponse(c, http.StatusForbidden, error)
}

func BadRequest(c *gin.Context, error error) {
	errorResponse(c, http.StatusBadRequest, error)
}

func NotFound(c *gin.Context, error error) {
	errorResponse(c, http.StatusNotFound, error)
}

func Conflict(c *gin.Context, error error) {
	errorResponse(c, http.StatusConflict, error)
}

func InternalServerError(c *gin.Context) {
	errorResponse(c, http.StatusInternalServerError, nil)
}

func successResponse(c *gin.Context, code int, data any) {
	var d = data
	if data == nil {
		d = struct{}{}
	}
	c.JSON(code, response{
		Status: code,
		Data:   d,
	})
}

func paginatedResponse(c *gin.Context, code int, data any, pagination *Pagination) {
	var d = data
	if data == nil {
		d = struct{}{}
	}
	c.JSON(code, response{
		Status:     code,
		Data:       d,
		Pagination: pagination,
	})
}

func errorResponse(c *gin.Context, code int, error error) {
	m := http.StatusText(code)
	if error != nil {
		m = error.Error()
	}
	c.JSON(code, response{
		Status:  code,
		Message: m,
	})
}
