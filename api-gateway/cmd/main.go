package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Настройка прокси к User Service
	userServiceURL, err := url.Parse("http://user-service:8080")
	if err != nil {
		log.Fatal("Failed to parse user service URL:", err)
	}

	userServiceProxy := httputil.NewSingleHostReverseProxy(userServiceURL)

	// Роуты API Gateway
	r.Any("/auth/*path", gin.WrapH(userServiceProxy))
	r.Any("/users/*path", gin.WrapH(userServiceProxy))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Порт по умолчанию для API Gateway
	}

	log.Printf("API Gateway started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start API Gateway:", err)
	}
}
