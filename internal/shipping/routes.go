package shipping

import (
	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
	jwtSecret string,
) {
	shipping := router.Group("/api/v1/shipping")

	shipping.Use(
		middleware.AuthRequired(jwtSecret),
	)

	shipping.POST(
		"",
		handler.Create,
	)

	shipping.GET(
		"/:id",
		handler.Get,
	)

	shipping.GET(
		"/order/:order_id",
		handler.GetByOrder,
	)

	shipping.PUT(
		"/:id",
		handler.Update,
	)

	shipping.DELETE(
		"/:id",
		handler.Delete,
	)
}
