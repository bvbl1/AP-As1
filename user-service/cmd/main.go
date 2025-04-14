package main

import (
	"Assignment1_AbylayMoldakhmet/proto/gen"
	grpcDelivery "Assignment1_AbylayMoldakhmet/user-service/internal/delivery/grpc"
	"Assignment1_AbylayMoldakhmet/user-service/internal/infrastructure/mongodb"
	"Assignment1_AbylayMoldakhmet/user-service/internal/usecase"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const (
	mongoURI        = "mongodb://localhost:27017"
	databaseName    = "ecommerce"
	jwtSecret       = "Vh8yxpK+3AwtcIj0BcX9uz/LmndCrQ7IInYMDXoMLqg="
	grpcPort        = ":50051"
	shutdownTimeout = 5 * time.Second
)

func main() {
	// Инициализация подключения к MongoDB
	client := connectMongoDB()
	defer client.Disconnect(context.Background())

	// Инициализация репозиториев и use cases
	repo := mongodb.NewUserRepo(client.Database(databaseName))
	authUC := usecase.NewAuthUsecase(repo, jwtSecret)
	userUC := usecase.NewUserUsecase(repo)

	// Создание и запуск gRPC сервера
	grpcServer := createGRPCServer(authUC, userUC)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("Starting graceful shutdown...")

		grpcServer.GracefulStop()
		log.Println("gRPC server stopped")

		client.Disconnect(context.Background())
		log.Println("MongoDB connection closed")
	}()

	log.Printf("gRPC server started on %s", grpcPort)
	if err := grpcServer.Serve(createListener()); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func connectMongoDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	return client
}

func createGRPCServer(authUC usecase.AuthUsecase, userUC usecase.UserUsecase) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	userServer := grpcDelivery.NewUserServer(authUC, userUC)
	gen.RegisterUserServiceServer(server, userServer)

	return server
}

func createListener() net.Listener {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	return lis
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	log.Printf("GRPC call: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	log.Printf("Method: %s, Duration: %s, Error: %v",
		info.FullMethod,
		time.Since(start),
		err)

	return resp, err
}
