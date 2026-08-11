package order

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {

	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	order, err := h.service.Create(
		c.Request.Context(),
		req,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create order",
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *Handler) Get(c *gin.Context) {

	order, err := h.service.GetByID(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {

		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "order not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch order",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *Handler) UpdateStatus(c *gin.Context) {

	var req UpdateStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.UpdateStatus(
		c.Request.Context(),
		c.Param("id"),
		req.Status,
	)

	if err != nil {

		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "order not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update order",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order updated successfully",
	})
}
