package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/internal/auth"
	"github.com/rho-commerce/rho/internal/cart"
	"github.com/rho-commerce/rho/internal/category"
	"github.com/rho-commerce/rho/internal/config"
	"github.com/rho-commerce/rho/internal/database"
	"github.com/rho-commerce/rho/internal/discount"
	"github.com/rho-commerce/rho/internal/inventory"
	"github.com/rho-commerce/rho/internal/middleware"
	"github.com/rho-commerce/rho/internal/order"
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

	cartRepo := cart.NewRepository(db)
	cartService := cart.NewService(cartRepo)
	cartHandler := cart.NewHandler(cartService)

	cart.RegisterRoutes(router, cartHandler, cfg.JWTSecret)

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	category.RegisterRoutes(router, categoryHandler, cfg.JWTSecret)
	inventoryRepo := inventory.NewRepository(db)
	inventoryService := inventory.NewService(inventoryRepo)
	inventoryHandler := inventory.NewHandler(inventoryService)

	inventory.RegisterRoutes(
		router,
		inventoryHandler,
		cfg.JWTSecret,
	)

	discountRepo := discount.NewRepository(db)
	discountService := discount.NewService(discountRepo)
	discountHandler := discount.NewHandler(discountService)

	discount.RegisterRoutes(
		router,
		discountHandler,
		cfg.JWTSecret,
	)

	orderRepo := order.NewRepository(db)
	orderService := order.NewService(orderRepo)
	orderHandler := order.NewHandler(orderService)

	order.RegisterRoutes(router, orderHandler, cfg.JWTSecret)

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
