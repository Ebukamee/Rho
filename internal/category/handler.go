package category

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
	var req CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	category, err := h.service.Create(
		c.Request.Context(),
		req,
	)

	if err != nil {
		if errors.Is(err, ErrConflict) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "category slug already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create category",
		})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *Handler) List(c *gin.Context) {
	params := pagination.Parse(c)

	categories, err := h.service.List(
		c.Request.Context(),
		params.Page,
		params.Limit,
		true,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch categories",
		})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *Handler) AdminList(c *gin.Context) {
	params := pagination.Parse(c)

	activeOnly := true

	if value := c.Query("active"); value != "" {
		activeOnly, _ = strconv.ParseBool(value)
	}

	categories, err := h.service.List(
		c.Request.Context(),
		params.Page,
		params.Limit,
		activeOnly,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch categories",
		})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *Handler) Get(c *gin.Context) {
	category, err := h.service.GetByID(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "category not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch category",
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	category, err := h.service.Update(
		c.Request.Context(),
		c.Param("id"),
		req,
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "category not found",
			})
			return
		}

		if errors.Is(err, ErrConflict) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "category slug already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update category",
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *Handler) Delete(c *gin.Context) {
	err := h.service.Delete(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "category not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "category deleted successfully",
	})
}
