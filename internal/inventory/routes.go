package inventory

import (
	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
	jwtSecret string,
) {
	inventory := router.Group("/api/v1/inventory")

	inventory.Use(
		middleware.AuthRequired(jwtSecret),
		middleware.RoleRequired("admin", "super_admin"),
	)

	inventory.POST("", handler.Create)

	inventory.GET("/:id", handler.Get)

	inventory.GET(
		"/product/:productID",
		handler.GetByProduct,
	)

	inventory.PUT("/:id", handler.Update)

	inventory.POST(
		"/product/:productID/adjust",
		handler.Adjust,
	)

	inventory.DELETE("/:id", handler.Delete)
}
