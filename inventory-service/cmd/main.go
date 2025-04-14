package main

import (
	"context"
	"log"
	"net"

	inventorygrpc "Assignment1_AbylayMoldakhmet/inventory-service/internal/delivery/grpc"
	"Assignment1_AbylayMoldakhmet/inventory-service/internal/repository/mongodb"
	"Assignment1_AbylayMoldakhmet/inventory-service/internal/usecase"
	"Assignment1_AbylayMoldakhmet/proto/gen"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	mongoURI := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("ecommerce")
	productRepo := mongodb.NewProductRepo(db)
	productUsecase := usecase.NewProductUsecase(productRepo)

	grpcServer := grpc.NewServer()
	grpcInventoryServer := inventorygrpc.NewInventoryServer(productUsecase)

	gen.RegisterInventoryServiceServer(grpcServer, grpcInventoryServer)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Inventory gRPC server started on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
