package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/internal/auth"
	"github.com/rho-commerce/rho/internal/config"
	"github.com/rho-commerce/rho/internal/database"
	"github.com/rho-commerce/rho/internal/middleware"
)

func main() {
	cfg := config.Load()
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	db := database.Connect(cfg.DatabaseURL)
	defer db.Close()

	authRepository := auth.NewRepository(db)
	authService := auth.NewService(
		authRepository,
		cfg.JWTSecret,
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
	)
	authHandler := auth.NewHandler(authService)

	router := gin.New()
	router.Use(
		middleware.CORS(cfg.CORSOrigins),
		gin.Logger(),
		gin.Recovery(),
	)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	auth.RegisterRoutes(router, authHandler, cfg.JWTSecret)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
