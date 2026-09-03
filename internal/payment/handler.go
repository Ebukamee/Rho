package payment

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

func (h *Handler) Initialize(c *gin.Context) {
	var req InitializePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.service.Initialize(
		c.Request.Context(),
		req,
		c.GetString("user_id"),
		c.GetString("email"),
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "order not found",
			})

		case errors.Is(err, ErrProviderNotFound):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "payment provider not supported",
			})

		case errors.Is(err, ErrOrderNotPayable):
			c.JSON(http.StatusConflict, gin.H{
				"error": "order is not payable",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Get(c *gin.Context) {
	payment, err := h.service.GetByID(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "payment not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch payment",
		})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *Handler) Verify(c *gin.Context) {
	payment, err := h.service.Verify(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "payment not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, payment)
}
