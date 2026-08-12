package checkout

import (
	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
	jwtSecret string,
) {
	checkout := router.Group("/api/v1/checkout")

	checkout.Use(
		middleware.AuthRequired(jwtSecret),
	)

	checkout.POST(
		"/preview",
		handler.Preview,
	)

	checkout.POST(
		"",
		handler.Create,
	)
}
