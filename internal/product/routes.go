package product

import (
	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
	jwtSecret string,
) {
	products := router.Group("/api/v1/products")

	// Public storefront endpoints.
	products.GET("", handler.List)
	products.GET("/:id", handler.Get)

	// Admin management endpoints.
	admin := products.Group("")
	admin.Use(
		middleware.AuthRequired(jwtSecret),
		middleware.RoleRequired("admin", "super_admin"),
	)

	admin.POST("", handler.Create)
	admin.GET("/admin", handler.AdminList)
	admin.PUT("/:id", handler.Update)
	admin.DELETE("/:id", handler.Delete)
}