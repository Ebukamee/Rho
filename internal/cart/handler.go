package cart

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

func getUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return "", false
	}

	id, ok := userID.(string)

	if !ok || id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identity",
		})
		return "", false
	}

	return id, true
}

func (h *Handler) Get(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	cart, err := h.service.GetCart(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch cart",
		})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *Handler) AddItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req AddItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cart, err := h.service.AddItem(
		c.Request.Context(),
		userID,
		req,
	)

	if err != nil {
		if errors.Is(err, ErrDuplicateItem) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "product already exists in cart",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to add item to cart",
		})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *Handler) UpdateItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req UpdateItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cart, err := h.service.UpdateItem(
		c.Request.Context(),
		userID,
		c.Param("itemID"),
		req,
	)

	if err != nil {
		if errors.Is(err, ErrCartItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "cart item not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update cart item",
		})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *Handler) RemoveItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	cart, err := h.service.RemoveItem(
		c.Request.Context(),
		userID,
		c.Param("itemID"),
	)

	if err != nil {
		if errors.Is(err, ErrCartItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "cart item not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to remove cart item",
		})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *Handler) Clear(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	if err := h.service.Clear(
		c.Request.Context(),
		userID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to clear cart",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "cart cleared successfully",
	})
}
