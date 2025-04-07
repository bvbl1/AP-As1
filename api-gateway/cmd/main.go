package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Прокидываем регистрацию в user-service
	r.POST("/auth/register", func(c *gin.Context) {
		resp, err := http.Post("http://localhost:8081/auth/register", "application/json", c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "user service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Прокидываем логин в user-service
	r.POST("/auth/login", func(c *gin.Context) {
		resp, err := http.Post("http://localhost:8081/auth/login", "application/json", c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "user service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Проксирование запросов к пользователям
	r.GET("/users/:id", func(c *gin.Context) {
		// 1. Создаем URL для user-service
		url := "http://localhost:8081/users/" + c.Param("id")

		// 2. Создаем новый GET-запрос
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		// 3. Копируем заголовок Authorization
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			req.Header.Set("Authorization", authHeader)
		}

		// 4. Выполняем запрос
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "user service unavailable"})
			return
		}
		defer resp.Body.Close()

		// 5. Возвращаем ответ
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Хелсчек
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Запуск сервера API Gateway
	r.Run(":8000")
}
