package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/internal/auth"
	"github.com/rho-commerce/rho/internal/config"
	"github.com/rho-commerce/rho/internal/database"
	"github.com/rho-commerce/rho/internal/middleware"
	"github.com/rho-commerce/rho/internal/product"
	applogger "github.com/rho-commerce/rho/pkg/logger"
)

func main() {
	cfg := config.Load()
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	appLogger := applogger.New(cfg.Environment)

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
		middleware.RequestLogger(appLogger),
		middleware.CORS(cfg.CORSOrigins),
		gin.Recovery(),
	)
	productRepository := product.NewRepository(db)
	productService := product.NewService(productRepository)
	productHandler := product.NewHandler(productService)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	auth.RegisterRoutes(router, authHandler, cfg.JWTSecret)
	product.RegisterRoutes(router, productHandler, cfg.JWTSecret)

	appLogger.Info(
		"starting Rho API",
		"environment", cfg.Environment,
		"port", cfg.Port,
	)

	if err := router.Run(":" + cfg.Port); err != nil {
		appLogger.Error("server stopped", "error", err)
	}
}
