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
		url := "http://localhost:8081/users/" + c.Param("id")
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		// Копируем заголовок Authorization
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			req.Header.Set("Authorization", authHeader)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "user service unavailable"})
			return
		}
		defer resp.Body.Close()

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Проксирование запросов к inventory-service
	r.POST("/products", func(c *gin.Context) {
		resp, err := http.Post("http://localhost:8082/products", "application/json", c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "inventory service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	r.GET("/products/:id", func(c *gin.Context) {
		resp, err := http.Get("http://localhost:8082/products/" + c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "inventory service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	r.PATCH("/products/:id", func(c *gin.Context) {
		req, err := http.NewRequest("PATCH", "http://localhost:8082/products/"+c.Param("id"), c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "inventory service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	r.DELETE("/products/:id", func(c *gin.Context) {
		req, err := http.NewRequest("DELETE", "http://localhost:8082/products/"+c.Param("id"), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "inventory service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	r.GET("/products", func(c *gin.Context) {
		resp, err := http.Get("http://localhost:8082/products")
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "inventory service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Проксирование запросов к order-service
	r.POST("/orders", func(c *gin.Context) {
		resp, err := http.Post("http://localhost:8083/orders", "application/json", c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "order service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	r.GET("/orders/:id", func(c *gin.Context) {
		resp, err := http.Get("http://localhost:8083/orders/" + c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "order service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	r.PATCH("/orders/:id", func(c *gin.Context) {
		req, err := http.NewRequest("PATCH", "http://localhost:8083/orders/"+c.Param("id"), c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "order service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	r.GET("/orders", func(c *gin.Context) {
		resp, err := http.Get("http://localhost:8083/orders" + c.Request.URL.String())
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "order service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Хелсчек
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Запуск сервера API Gateway
	r.Run(":8000")
}
