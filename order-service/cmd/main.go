package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Assignment1_AbylayMoldakhmet/order-service/internal/clients"
	orderHttp "Assignment1_AbylayMoldakhmet/order-service/internal/delivery/http"
	"Assignment1_AbylayMoldakhmet/order-service/internal/repository/mongodb"
	"Assignment1_AbylayMoldakhmet/order-service/internal/usecase"
)

func main() {
	mongoURI := "mongodb://localhost:27017"
	databaseName := "ecommerce"
	collectionName := "orders"
	inventoryServiceURL := "http://localhost:8082"

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	db := client.Database(databaseName)
	orderRepo := mongodb.NewOrderRepo(db.Collection(collectionName))
	inventoryClient := clients.NewInventoryClient(inventoryServiceURL) // передаем указатель

	orderUC := usecase.NewOrderUsecase(orderRepo, inventoryClient) // передаем указатель на inventoryClient

	router := gin.Default()
	orderHttp.RegisterOrderRoutes(router, orderUC)

	log.Printf("Order Service is running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
