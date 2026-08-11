package discount

import (
	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
	jwtSecret string,
) {
	discounts := router.Group("/api/v1/discounts")

	// Storefront: customers can validate a coupon.
	discounts.Use(
		middleware.AuthRequired(jwtSecret),
	)

	discounts.POST("/apply", handler.Apply)

	// Admin management.
	admin := discounts.Group("")
	admin.Use(
		middleware.RoleRequired("admin", "super_admin"),
	)

	admin.POST("", handler.Create)
	admin.GET("/:id", handler.Get)
	admin.PUT("/:id", handler.Update)
	admin.DELETE("/:id", handler.Delete)
}
