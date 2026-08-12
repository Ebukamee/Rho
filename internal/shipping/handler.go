package shipping

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
	var req CreateShipmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	shipment, err := h.service.Create(
		c.Request.Context(),
		req,
		c.GetString("user_id"),
	)

	if err != nil {
		if errors.Is(err, ErrConflict) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "shipment already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create shipment",
		})
		return
	}

	c.JSON(http.StatusCreated, shipment)
}

func (h *Handler) Get(c *gin.Context) {
	shipment, err := h.service.GetByID(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "shipment not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch shipment",
		})
		return
	}

	c.JSON(http.StatusOK, shipment)
}

func (h *Handler) GetByOrder(c *gin.Context) {
	shipment, err := h.service.GetByOrderID(
		c.Request.Context(),
		c.Param("order_id"),
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "shipment not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch shipment",
		})
		return
	}

	c.JSON(http.StatusOK, shipment)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateShipmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	shipment, err := h.service.Update(
		c.Request.Context(),
		c.Param("id"),
		req,
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "shipment not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update shipment",
		})
		return
	}

	c.JSON(http.StatusOK, shipment)
}

func (h *Handler) Delete(c *gin.Context) {
	err := h.service.Delete(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "shipment not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete shipment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "shipment deleted successfully",
	})
}
