package main

import (
	"Assignment1_AbylayMoldakhmet/user-service/internal/delivery/http"
	"Assignment1_AbylayMoldakhmet/user-service/internal/infrastructure/mongodb"
	"Assignment1_AbylayMoldakhmet/user-service/internal/usecase"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 1. Подключение к MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	// Проверка подключения
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	db := client.Database("ecommerce")

	// 2. Инициализация репозитория
	userRepo := mongodb.NewUserRepo(db)

	// 3. Инициализация Usecase
	userUsecase := usecase.NewUserUsecase(userRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo, os.Getenv("JWT_SECRET"))

	// 4. Инициализация HTTP обработчиков
	userHandler := http.NewUserHandler(userUsecase)
	authHandler := http.NewAuthHandler(authUsecase)

	// 5. Настройка Gin
	r := gin.Default()

	// 6. Регистрация роутов
	// Auth routes
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	// User CRUD routes (требуют JWT)
	userRoutes := r.Group("/users")
	// userRoutes.Use(middleware.JwtAuthMiddleware(os.Getenv("JWT_SECRET"))) // Раскомментировать после настройки middleware
	{
		userRoutes.GET("/:id", userHandler.GetByID)
		userRoutes.PUT("/:id", userHandler.Update)
		userRoutes.DELETE("/:id", userHandler.Delete)
	}

	// 7. Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Порт по умолчанию
	}
	log.Printf("Server started on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
