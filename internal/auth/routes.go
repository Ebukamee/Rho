package auth

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all auth-related routes
func RegisterRoutes(router *gin.Engine, handler *Handler) {
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
		protected.Use(AuthMiddleware())
		{
			protected.POST("/logout", handler.Logout)
			protected.GET("/me", handler.GetProfile)
			protected.PUT("/profile", handler.UpdateProfile)
			protected.PUT("/password", handler.ChangePassword)
		}
	}

	// Admin routes
	usersGroup := router.Group("/api/v1/users")
	usersGroup.Use(AuthMiddleware(), AdminMiddleware())
	{
		usersGroup.GET("", handler.ListUsers)
		usersGroup.GET("/:id", handler.GetUser)
		usersGroup.PUT("/:id", handler.UpdateUser)
		usersGroup.DELETE("/:id", handler.DeleteUser)
		usersGroup.PUT("/:id/role", handler.UpdateUserRole)
	}
}

// AuthMiddleware validates JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement JWT validation
		c.Next()
	}
}

// AdminMiddleware checks if user has admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement admin role check
		c.Next()
	}
}
