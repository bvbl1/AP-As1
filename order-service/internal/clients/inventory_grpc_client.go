package clients

import (
	"Assignment1_AbylayMoldakhmet/proto/gen"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InventoryGRPCClient struct {
	client gen.InventoryServiceClient
}

func NewInventoryGRPCClient(addr string) *InventoryGRPCClient {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	return &InventoryGRPCClient{
		client: gen.NewInventoryServiceClient(conn),
	}
}

func (c *InventoryGRPCClient) CheckStock(productID string, quantity int) (bool, error) {
	resp, err := c.client.CheckStock(context.Background(), &gen.StockCheckRequest{
		ProductId: productID,
		Quantity:  int32(quantity),
	})
	if err != nil {
		return false, err
	}
	return resp.IsAvailable, nil
}
