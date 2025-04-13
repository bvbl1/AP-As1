package main

import (
	"Assignment1_AbylayMoldakhmet/user-service/internal/delivery/http"
	"Assignment1_AbylayMoldakhmet/user-service/internal/infrastructure/mongodb"
	"Assignment1_AbylayMoldakhmet/user-service/internal/usecase"
	"Assignment1_AbylayMoldakhmet/user-service/middleware"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	db := client.Database("ecommerce")

	userRepo := mongodb.NewUserRepo(db)

	userUsecase := usecase.NewUserUsecase(userRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo, "Vh8yxpK+3AwtcIj0BcX9uz/LmndCrQ7IInYMDXoMLqg=")

	userHandler := http.NewUserHandler(userUsecase)
	authHandler := http.NewAuthHandler(authUsecase)

	r := gin.Default()

	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	userRoutes := r.Group("/users")
	userRoutes.Use(middleware.JwtAuthMiddleware("Vh8yxpK+3AwtcIj0BcX9uz/LmndCrQ7IInYMDXoMLqg="))
	{
		userRoutes.GET("/:id", userHandler.GetByID)
		userRoutes.PUT("/:id", userHandler.Update)
		userRoutes.DELETE("/:id", userHandler.Delete)
	}

	port := "8081"
	log.Printf("User service started on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
