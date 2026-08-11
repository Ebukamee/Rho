package discount

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
	var req CreateDiscountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	d, err := h.service.Create(
		c.Request.Context(),
		req,
	)

	if err != nil {
		if errors.Is(err, ErrConflict) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "discount code already exists",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, d)
}

func (h *Handler) Get(c *gin.Context) {
	d, err := h.service.GetByID(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "discount not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch discount",
		})
		return
	}

	c.JSON(http.StatusOK, d)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateDiscountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	d, err := h.service.Update(
		c.Request.Context(),
		c.Param("id"),
		req,
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "discount not found",
			})
			return
		}

		if errors.Is(err, ErrConflict) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "discount code already exists",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, d)
}

func (h *Handler) Apply(c *gin.Context) {
	var req ApplyDiscountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.service.Apply(
		c.Request.Context(),
		req,
	)

	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, ErrNotFound) {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(
		c.Request.Context(),
		c.Param("id"),
	); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "discount not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete discount",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "discount deleted successfully",
	})
}
