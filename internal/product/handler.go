package product

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/pkg/pagination"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	product, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrConflict) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "product slug or SKU already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *Handler) Get(c *gin.Context) {
	product, err := h.service.GetByID(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch product",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) List(c *gin.Context) {
	params := pagination.Parse(c)

	products, err := h.service.List(
		c.Request.Context(),
		params.Page,
		params.Limit,
		true,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch products",
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) AdminList(c *gin.Context) {
	params := pagination.Parse(c)

	activeOnly := true

	if value := c.Query("active"); value != "" {
		activeOnly, _ = strconv.ParseBool(value)
	}

	products, err := h.service.List(
		c.Request.Context(),
		params.Page,
		params.Limit,
		activeOnly,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch products",
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	product, err := h.service.Update(
		c.Request.Context(),
		c.Param("id"),
		req,
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "product not found",
			})
			return
		}

		if errors.Is(err, ErrConflict) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "product slug or SKU already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update product",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) Delete(c *gin.Context) {
	err := h.service.Delete(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product deleted successfully",
	})
}
