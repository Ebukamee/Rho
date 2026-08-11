package order

import (
	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
	jwtSecret string,
) {

	orders := router.Group("/api/v1/orders")

	orders.Use(
		middleware.AuthRequired(jwtSecret),
	)

	orders.POST("", handler.Create)
	orders.GET("/:id", handler.Get)

	admin := orders.Group("")

	admin.Use(
		middleware.RoleRequired("admin", "super_admin"),
	)

	admin.PUT("/:id/status", handler.UpdateStatus)
}
