package pagination

import (
	"strconv"

	"github.com/gin-conic/gin"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

func Parse(c *gin.Context) Params {
	page := positiveInt(c.Query("page"), DefaultPage)
	limit := positiveInt(c.Query("limit"), DefaultLimit)

	if limit > MaxLimit {
		limit = MaxLimit
	}

	return Params{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

func TotalPages(total, limit int) int {
	if total == 0 || limit <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}

func positiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}

	return parsed
}
