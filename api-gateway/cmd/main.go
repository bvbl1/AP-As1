package main

import (
	"Assignment1_AbylayMoldakhmet/api-gateway/internal/middleware"
	"Assignment1_AbylayMoldakhmet/proto/gen"
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	userServiceAddr     = "localhost:50051"       // gRPC адрес user-service
	inventoryServiceURL = "http://localhost:8082" // REST адрес inventory-service
	orderServiceURL     = "http://localhost:8083" // REST адрес order-service
	jwtSecret           = "Vh8yxpK+3AwtcIj0BcX9uz/LmndCrQ7IInYMDXoMLqg="
)

func main() {
	// Создаем gRPC соединение для user-service
	userConn, err := grpc.Dial(userServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer userConn.Close()

	userClient := gen.NewUserServiceClient(userConn)

	r := gin.Default()
	authMiddleware := middleware.JwtAuthMiddleware(jwtSecret)

	// Auth endpoints
	r.POST("/auth/register", createAuthHandler(userClient, "Register"))
	r.POST("/auth/login", createAuthHandler(userClient, "Login"))

	protected := r.Group("/")
	protected.Use(authMiddleware)

	// User endpoints
	protected.GET("/users/:id", createUserHandler(userClient, "GetUserProfile"))
	protected.PUT("/users/:id", createUserHandler(userClient, "UpdateUser"))
	protected.DELETE("/users/:id", createUserHandler(userClient, "DeleteUser"))

	// Inventory endpoints (остаются REST)
	protected.POST("/products", proxyHandler(inventoryServiceURL))
	protected.GET("/products/:id", proxyHandler(inventoryServiceURL))
	protected.PATCH("/products/:id", proxyHandler(inventoryServiceURL))
	protected.DELETE("/products/:id", proxyHandler(inventoryServiceURL))
	protected.GET("/products", proxyHandler(inventoryServiceURL))

	// Order endpoints (остаются REST)
	protected.POST("/orders", proxyHandler(orderServiceURL))
	protected.GET("/orders/:id", proxyHandler(orderServiceURL))
	protected.PATCH("/orders/:id", proxyHandler(orderServiceURL))
	protected.GET("/orders", proxyHandler(orderServiceURL))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(":8000")
}

// Обработчик для gRPC вызовов
func createUserHandler(client gen.UserServiceClient, method string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch method {
		case "GetUserProfile":
			resp, err := client.GetUserProfile(c.Request.Context(), &gen.UserIDRequest{
				UserId: c.Param("id"),
			})
			handleGRPCResponse(c, resp, err)

		case "UpdateUser":
			var requestBody gen.UpdateUserRequest
			if err := c.ShouldBindJSON(&requestBody); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			resp, err := client.UpdateUser(c.Request.Context(), &requestBody)
			handleGRPCResponse(c, resp, err)

		case "DeleteUser":
			resp, err := client.DeleteUser(c.Request.Context(), &gen.UserIDRequest{
				UserId: c.Param("id"),
			})
			handleGRPCResponse(c, resp, err)
		}
	}
}

// Обработчик для аутентификации
func createAuthHandler(client gen.UserServiceClient, method string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch method {
		case "Register":
			var req gen.RegisterRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			resp, err := client.Register(c.Request.Context(), &req)
			handleGRPCResponse(c, resp, err)

		case "Login":
			var req gen.LoginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			resp, err := client.Login(c.Request.Context(), &req)
			handleGRPCResponse(c, resp, err)
		}
	}
}

// Обработчик для REST прокси
func proxyHandler(baseURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := baseURL + c.Request.URL.Path

		// Создаем новый запрос
		body, _ := io.ReadAll(c.Request.Body)
		req, err := http.NewRequest(c.Request.Method, url, bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		// Копируем заголовки
		for k, v := range c.Request.Header {
			req.Header[k] = v
		}

		// Выполняем запрос
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
			return
		}
		defer resp.Body.Close()

		// Копируем ответ
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	}
}

// Универсальная обработка gRPC ответов
func handleGRPCResponse(c *gin.Context, response interface{}, err error) {
	if err != nil {
		c.JSON(convertGRPCError(err))
		return
	}

	switch resp := response.(type) {
	case *gen.UserResponse:
		c.JSON(http.StatusOK, gin.H{
			"id":    resp.Id,
			"email": resp.Email,
			"role":  resp.Role,
		})
	case *gen.LoginResponse:
		c.JSON(http.StatusOK, gin.H{
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
		})
	case *emptypb.Empty:
		c.Status(http.StatusNoContent)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown response type"})
	}
}

// Конвертация gRPC ошибок в HTTP статусы
func convertGRPCError(err error) (int, interface{}) {
	// Реализуйте конвертацию на основе gRPC status codes
	return http.StatusInternalServerError, gin.H{"error": err.Error()}
}
