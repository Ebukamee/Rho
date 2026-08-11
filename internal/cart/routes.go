package cart

import (
	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
	jwtSecret string,
) {
	cart := router.Group("/api/v1/cart")

	cart.Use(
		middleware.AuthRequired(jwtSecret),
	)

	cart.GET("", handler.Get)

	cart.POST("/items", handler.AddItem)

	cart.PUT("/items/:itemID", handler.UpdateItem)

	cart.DELETE("/items/:itemID", handler.RemoveItem)

	cart.DELETE("", handler.Clear)
}
