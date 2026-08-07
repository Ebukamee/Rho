package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/rho-commerce/rho/internal/middleware"
)

// RegisterRoutes registers all auth-related routes
func RegisterRoutes(router *gin.Engine, handler *Handler, jwtSecret string) {
	authGroup := router.Group("/api/v1/auth")
	{
		// Public routes
		authGroup.POST("/signup", handler.Signup)
		authGroup.POST("/login", handler.Login)
		authGroup.POST("/refresh", handler.RefreshToken)

		// OAuth routes
		authGroup.GET("/google/login", handler.GoogleLogin)
		authGroup.GET("/google/callback", handler.GoogleCallback)

		// Protected routes (require authentication)
		protected := authGroup.Group("")
		protected.Use(middleware.AuthRequired(jwtSecret))
		{
			protected.POST("/logout", handler.Logout)
			protected.GET("/me", handler.GetProfile)
			protected.PUT("/profile", handler.UpdateProfile)
			protected.PUT("/password", handler.ChangePassword)
		}
	}

	// Admin routes
	usersGroup := router.Group("/api/v1/users")
	usersGroup.Use(middleware.AuthRequired(jwtSecret), middleware.RoleRequired(RoleAdmin, RoleSuperAdmin))
	{
		usersGroup.GET("", handler.ListUsers)
		usersGroup.GET("/:id", handler.GetUser)
		usersGroup.PUT("/:id", handler.UpdateUser)
		usersGroup.DELETE("/:id", handler.DeleteUser)
		usersGroup.PUT("/:id/role", handler.UpdateUserRole)
	}
}
