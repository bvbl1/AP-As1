package main

import (
	"Assignment1_AbylayMoldakhmet/api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Публичные роуты
	public := r.Group("/auth")
	{
		public.POST("/login", proxy.ToUserService) // Перенаправляем в User Service
		public.POST("/register", proxy.ToUserService)
	}

	// Защищенные роуты
	protected := r.Group("/")
	protected.Use(middleware.JwtAuthMiddleware("your_jwt_secret"))
	{
		protected.GET("/orders", proxy.ToOrderService)
		protected.POST("/orders", proxy.ToOrderService)
		protected.GET("/products", proxy.ToInventoryService)
		// Другие защищенные эндпоинты
	}

	r.Run(":8080")
}
