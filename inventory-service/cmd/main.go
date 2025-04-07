package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Assignment1_AbylayMoldakhmet/inventory-service/internal/delivery"
	"Assignment1_AbylayMoldakhmet/inventory-service/internal/repository/mongodb"
	"Assignment1_AbylayMoldakhmet/inventory-service/internal/usecase"
)

func main() {
	mongoURI := "mongodb://localhost:27017"

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("ecommerce")

	productRepo := mongodb.NewProductRepo(db)
	productUsecase := usecase.NewProductUsecase(productRepo)

	r := gin.Default()
	delivery.NewProductHandler(r, productUsecase)

	port := "8082"
	log.Printf("Inventory Service started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server error:", err)
	}
}
