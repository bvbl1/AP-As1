package main

import (
	"Assignment1_AbylayMoldakhmet/api-gateway/internal/middleware"
	"Assignment1_AbylayMoldakhmet/proto/gen"
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	userServiceAddr      = "localhost:50051"
	inventoryServiceAddr = "localhost:50052"
	orderServiceAddr     = "localhost:50053"
	jwtSecret            = "Vh8yxpK+3AwtcIj0BcX9uz/LmndCrQ7IInYMDXoMLqg="
)

func main() {
	userConn := initGRPCConn(userServiceAddr)
	defer userConn.Close()
	userClient := gen.NewUserServiceClient(userConn)

	inventoryConn := initGRPCConn(inventoryServiceAddr)
	defer inventoryConn.Close()
	inventoryClient := gen.NewInventoryServiceClient(inventoryConn)

	orderConn := initGRPCConn(orderServiceAddr)
	defer orderConn.Close()
	orderClient := gen.NewOrderServiceClient(orderConn)

	r := gin.Default()
	authMiddleware := middleware.JwtAuthMiddleware(jwtSecret)

	r.POST("/auth/register", createAuthHandler(userClient, "Register"))
	r.POST("/auth/login", createAuthHandler(userClient, "Login"))

	protected := r.Group("/")
	protected.Use(authMiddleware)

	protected.GET("/users/:id", createUserHandler(userClient, "GetUserProfile"))
	protected.PUT("/users/:id", createUserHandler(userClient, "UpdateUser"))
	protected.DELETE("/users/:id", createUserHandler(userClient, "DeleteUser"))

	protected.POST("/products", createInventoryHandler(inventoryClient, "CreateProduct"))
	protected.GET("/products/:id", createInventoryHandler(inventoryClient, "GetProduct"))
	protected.PATCH("/products/:id", createInventoryHandler(inventoryClient, "UpdateProduct"))
	protected.DELETE("/products/:id", createInventoryHandler(inventoryClient, "DeleteProduct"))
	protected.GET("/products", createInventoryHandler(inventoryClient, "ListProducts"))

	protected.POST("/orders", createOrderHandler(orderClient, "CreateOrder"))
	protected.GET("/orders/:id", createOrderHandler(orderClient, "GetOrder"))
	protected.PATCH("/orders/:id", createOrderHandler(orderClient, "UpdateOrderStatus"))
	protected.GET("/orders", createOrderHandler(orderClient, "ListOrders"))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Println("API Gateway started on :8000")
	r.Run(":8000")
}

func initGRPCConn(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", addr, err)
	}
	return conn
}

func createOrderHandler(client gen.OrderServiceClient, method string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Order handler called for method: %s, path: %s", method, c.FullPath())

		switch method {
		case "CreateOrder":
			var req gen.CreateOrderRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Printf("CreateOrder bad request: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			log.Printf("CreateOrder request: %+v", req)
			resp, err := client.CreateOrder(c.Request.Context(), &req)
			handleGRPCResponse(c, resp, err)

		case "GetOrder":
			orderID := c.Param("id")
			log.Printf("GetOrder request for ID: %s", orderID)
			resp, err := client.GetOrder(c.Request.Context(), &gen.OrderIDRequest{
				Id: orderID,
			})
			handleGRPCResponse(c, resp, err)

		case "UpdateOrderStatus":
			orderID := c.Param("id")
			log.Printf("UpdateOrderStatus request for ID: %s", orderID)

			var jsonBody struct {
				Status string `json:"status"`
			}

			if err := c.ShouldBindJSON(&jsonBody); err != nil {
				log.Printf("UpdateOrderStatus bad request: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			var status gen.OrderStatus
			switch jsonBody.Status {
			case "pending":
				status = gen.OrderStatus_PENDING
			case "paid":
				status = gen.OrderStatus_PAID
			case "cancelled":
				status = gen.OrderStatus_CANCELLED
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order status"})
				return
			}

			req := &gen.OrderStatusUpdateRequest{
				Id:     orderID,
				Status: status,
			}

			resp, err := client.UpdateOrderStatus(c.Request.Context(), req)
			handleGRPCResponse(c, resp, err)

		case "ListOrders":
			resp, err := client.ListOrders(c.Request.Context(), &gen.OrderFilterRequest{})
			handleGRPCResponse(c, resp, err)

		default:
			log.Printf("Unknown order method: %s", method)
			c.JSON(http.StatusNotFound, gin.H{"error": "method not found"})
		}
	}
}

// Inventory handlers
func createInventoryHandler(client gen.InventoryServiceClient, method string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Inventory handler called for method: %s, path: %s", method, c.FullPath())

		switch method {
		case "CreateProduct":
			var req gen.CreateProductRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Printf("CreateProduct bad request: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			log.Printf("CreateProduct request: %+v", req)
			resp, err := client.CreateProduct(c.Request.Context(), &req)
			handleGRPCResponse(c, resp, err)

		case "GetProduct":
			productID := c.Param("id")
			log.Printf("GetProduct request for ID: %s", productID)

			resp, err := client.GetProduct(c.Request.Context(), &gen.ProductIDRequest{
				Id: productID,
			})

			if err != nil {
				log.Printf("GetProduct error: %v", err)
			} else {
				log.Printf("GetProduct response: %+v", resp)
			}
			handleGRPCResponse(c, resp, err)

		case "UpdateProduct":
			productID := c.Param("id")
			log.Printf("UpdateProduct request for ID: %s", productID)

			var req gen.UpdateProductRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Printf("UpdateProduct bad request: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			req.Id = productID

			resp, err := client.UpdateProduct(c.Request.Context(), &req)
			handleGRPCResponse(c, resp, err)

		case "DeleteProduct":
			productID := c.Param("id")
			log.Printf("DeleteProduct request for ID: %s", productID)

			resp, err := client.DeleteProduct(c.Request.Context(), &gen.ProductIDRequest{
				Id: productID,
			})
			handleGRPCResponse(c, resp, err)

		case "ListProducts":
			log.Printf("ListProducts request with query: %v", c.Request.URL.Query())

			filter := make(map[string]string)
			for k, v := range c.Request.URL.Query() {
				filter[k] = v[0]
			}

			resp, err := client.ListProducts(c.Request.Context(), &gen.ListProductsRequest{
				Filter: filter,
			})
			handleGRPCResponse(c, resp, err)

		default:
			log.Printf("Unknown inventory method: %s", method)
			c.JSON(http.StatusNotFound, gin.H{"error": "method not found"})
		}
	}
}

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

func proxyHandler(baseURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := baseURL + c.Request.URL.Path

		body, _ := io.ReadAll(c.Request.Body)
		req, err := http.NewRequest(c.Request.Method, url, bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		for k, v := range c.Request.Header {
			req.Header[k] = v
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
			return
		}
		defer resp.Body.Close()

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	}
}

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
	case *gen.ProductResponse:
		c.JSON(http.StatusOK, gin.H{
			"id":       resp.Id,
			"name":     resp.Name,
			"price":    resp.Price,
			"category": resp.Category,
			"stock":    resp.Stock,
		})
	case *gen.ListProductsResponse:
		c.JSON(http.StatusOK, resp.Products)

	case *gen.OrderResponse:
		c.JSON(http.StatusOK, gin.H{
			"id":     resp.Id,
			"userId": resp.UserId,
			"status": resp.Status.String(),
			"items":  resp.Items,
		})

	case *gen.OrderListResponse:
		orders := make([]map[string]interface{}, 0)
		for _, o := range resp.Orders {
			orders = append(orders, map[string]interface{}{
				"id":     o.Id,
				"userId": o.UserId,
				"status": o.Status.String(),
				"items":  o.Items,
			})
		}
		c.JSON(http.StatusOK, gin.H{"orders": orders})

	case *emptypb.Empty:
		c.Status(http.StatusNoContent)

	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown response type"})
	}
}

func convertGRPCError(err error) (int, interface{}) {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}
	}

	switch st.Code() {
	case codes.NotFound:
		return http.StatusNotFound, gin.H{"error": st.Message()}
	case codes.InvalidArgument:
		return http.StatusBadRequest, gin.H{"error": st.Message()}
	case codes.PermissionDenied:
		return http.StatusForbidden, gin.H{"error": st.Message()}
	case codes.Unauthenticated:
		return http.StatusUnauthorized, gin.H{"error": st.Message()}
	default:
		return http.StatusInternalServerError, gin.H{"error": st.Message()}
	}
}
