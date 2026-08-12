package checkout

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.service.Create(
		c.Request.Context(),
		c.GetString("user_id"),
		req,
	)

	if err != nil {
		if errors.Is(err, ErrCartEmpty) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "cart is empty",
			})
			return
		}

		if errors.Is(err, ErrInvalidQuantity) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid cart quantity",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create checkout",
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *Handler) Preview(c *gin.Context) {
	var req Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.service.Preview(
		c.Request.Context(),
		c.GetString("user_id"),
		req,
	)

	if err != nil {
		if errors.Is(err, ErrCartEmpty) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "cart is empty",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
