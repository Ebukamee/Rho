package inventory

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
	var req CreateInventoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	inventory, err := h.service.Create(
		c.Request.Context(),
		req,
	)

	if err != nil {
		if errors.Is(err, ErrConflict) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "inventory already exists for this product",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create inventory",
		})
		return
	}

	c.JSON(http.StatusCreated, inventory)
}

func (h *Handler) GetByProduct(c *gin.Context) {
	inventory, err := h.service.GetByProductID(
		c.Request.Context(),
		c.Param("productID"),
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "inventory not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch inventory",
		})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (h *Handler) Get(c *gin.Context) {
	inventory, err := h.service.GetByID(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "inventory not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch inventory",
		})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateInventoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	inventory, err := h.service.Update(
		c.Request.Context(),
		c.Param("id"),
		req,
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "inventory not found",
			})
			return
		}

		if errors.Is(err, ErrInvalidQuantity) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update inventory",
		})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (h *Handler) Adjust(c *gin.Context) {
	var req AdjustInventoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	inventory, err := h.service.Adjust(
		c.Request.Context(),
		c.Param("productID"),
		req.Quantity,
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "inventory not found",
			})
			return
		}

		if errors.Is(err, ErrInvalidQuantity) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to adjust inventory",
		})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(
		c.Request.Context(),
		c.Param("id"),
	); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "inventory not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete inventory",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "inventory deleted successfully",
	})
}
