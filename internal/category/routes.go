package category

import (
	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
	jwtSecret string,
) {
	categories := router.Group("/api/v1/categories")

	// Public endpoints.
	categories.GET("", handler.List)
	categories.GET("/:id", handler.Get)

	// Admin endpoints.
	admin := categories.Group("")
	admin.Use(
		middleware.AuthRequired(jwtSecret),
		middleware.RoleRequired("admin", "super_admin"),
	)

	admin.POST("", handler.Create)
	admin.GET("/admin", handler.AdminList)
	admin.PUT("/:id", handler.Update)
	admin.DELETE("/:id", handler.Delete)
}
