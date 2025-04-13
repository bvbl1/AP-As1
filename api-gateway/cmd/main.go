package main

import (
	"Assignment1_AbylayMoldakhmet/api-gateway/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	authMiddleware := middleware.JwtAuthMiddleware(`Vh8yxpK+3AwtcIj0BcX9uz/LmndCrQ7IInYMDXoMLqg=`)

	r.POST("/auth/register", func(c *gin.Context) {
		resp, err := http.Post("http://localhost:8081/auth/register", "application/json", c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "user service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	r.POST("/auth/login", func(c *gin.Context) {
		resp, err := http.Post("http://localhost:8081/auth/login", "application/json", c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "user service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	protected := r.Group("/")
	protected.Use(authMiddleware)

	r.GET("/users/:id", func(c *gin.Context) {
		url := "http://localhost:8081/users/" + c.Param("id")
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

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
	protected.POST("/products", func(c *gin.Context) {
		resp, err := http.Post("http://localhost:8082/products", "application/json", c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "inventory service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	protected.GET("/products/:id", func(c *gin.Context) {
		resp, err := http.Get("http://localhost:8082/products/" + c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "inventory service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	protected.PATCH("/products/:id", func(c *gin.Context) {
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

	protected.DELETE("/products/:id", func(c *gin.Context) {
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

	protected.GET("/products", func(c *gin.Context) {
		resp, err := http.Get("http://localhost:8082/products")
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "inventory service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Проксирование запросов к order-service
	protected.POST("/orders", func(c *gin.Context) {
		resp, err := http.Post("http://localhost:8083/orders", "application/json", c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "order service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	protected.GET("/orders/:id", func(c *gin.Context) {
		resp, err := http.Get("http://localhost:8083/orders/" + c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "order service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	protected.PATCH("/orders/:id", func(c *gin.Context) {
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

	protected.GET("/orders", func(c *gin.Context) {
		req, err := http.NewRequest("GET", "http://localhost:8083/orders"+c.Request.URL.RawQuery, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		if auth := c.GetHeader("Authorization"); auth != "" {
			req.Header.Set("Authorization", auth)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "order service is down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(":8000")
}
