package main

import (
	"Assignment1_AbylayMoldakhmet/order-service/internal/clients"
	orderGrpc "Assignment1_AbylayMoldakhmet/order-service/internal/delivery/grpc"
	"Assignment1_AbylayMoldakhmet/order-service/internal/repository/mongodb"
	"Assignment1_AbylayMoldakhmet/order-service/internal/usecase"
	"Assignment1_AbylayMoldakhmet/proto/gen"
	"context"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	mongoURI := "mongodb://localhost:27017"
	databaseName := "ecommerce"
	collectionName := "orders"
	inventoryServiceAddress := "localhost:50052"
	grpcPort := ":50053"

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Fatalf("Mongo connect error: %v", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Mongo ping failed: %v", err)
	}

	db := client.Database(databaseName)
	orderRepo := mongodb.NewOrderRepo(db.Collection(collectionName))
	inventoryClient := clients.NewInventoryGRPCClient(inventoryServiceAddress)
	orderUC := usecase.NewOrderUsecase(orderRepo, inventoryClient)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	orderServer := orderGrpc.NewOrderServer(orderUC)
	gen.RegisterOrderServiceServer(grpcServer, orderServer)

	log.Printf("Order gRPC server is running on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
