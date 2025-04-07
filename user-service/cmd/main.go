package main

import (
	"Assignment1_AbylayMoldakhmet/user-service/internal/infrastructure/mongodb"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Подключение к MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("ecommerce")
	userRepo := mongodb.NewUserRepo(db)

	// Далее инициализация Usecase, Delivery и сервера...
}
