package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

func success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"data":    data,
	})
}

func successWithPagination(c *gin.Context, status int, data interface{}, pagination Pagination) {
	c.JSON(status, gin.H{
		"data":       data,
		"pagination": pagination,
	})
}

func Failure(c *gin.Context, status int, code string, message string) {
	c.JSON(status, gin.H{
		"error": Error{
			Code:    code,
			Message: message,
		},
	})
}

func BadRequest(c *gin.Context, code string, message string) {
	Failure(c, http.StatusBadRequest, code, "bad request: "+message)
}

func Unauthorized(c *gin.Context, code string, message string) {
	Failure(c, http.StatusUnauthorized, code, "unauthorized: "+message)
}

func Forbidden(c *gin.Context, code string, message string) {
	Failure(c, http.StatusForbidden, code, "forbidden: "+message)
}

func NotFound(c *gin.Context, code string, message string) {
	Failure(c, http.StatusNotFound, code, "not found: "+message)
}

func InternalServerError(c *gin.Context, code string, message string) {
	Failure(c, http.StatusInternalServerError, code, "internal server error: "+message)
}
