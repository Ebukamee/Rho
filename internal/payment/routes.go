package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
	jwtSecret string,
) {
	payments := router.Group("/api/v1/payments")

	payments.Use(
		middleware.AuthRequired(jwtSecret),
	)

	payments.POST("/initialize", handler.Initialize)
	payments.GET("/:id", handler.Get)
	payments.POST("/:id/verify", handler.Verify)
}
