package grpc

import (
	"Assignment1_AbylayMoldakhmet/inventory-service/internal/domain"
	"Assignment1_AbylayMoldakhmet/proto/gen"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type InventoryServer struct {
	gen.UnimplementedInventoryServiceServer
	uc domain.ProductUsecase
}

func NewInventoryServer(uc domain.ProductUsecase) *InventoryServer {
	return &InventoryServer{uc: uc}
}

func (s *InventoryServer) CreateProduct(ctx context.Context, req *gen.CreateProductRequest) (*gen.ProductResponse, error) {
	product := &domain.Product{
		Name:     req.Name,
		Price:    float64(req.Price),
		Category: req.Category,
		Stock:    int(req.Stock),
	}

	if err := s.uc.Create(product); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return s.convertProductToResponse(product), nil
}

func (s *InventoryServer) GetProduct(ctx context.Context, req *gen.ProductIDRequest) (*gen.ProductResponse, error) {
	product, err := s.uc.GetByID(req.Id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return s.convertProductToResponse(product), nil
}

func (s *InventoryServer) UpdateProduct(ctx context.Context, req *gen.UpdateProductRequest) (*gen.ProductResponse, error) {
	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	product := &domain.Product{
		ID:       objID,
		Name:     req.Name,
		Price:    float64(req.Price),
		Category: req.Category,
		Stock:    int(req.Stock),
	}

	if err := s.uc.Update(product); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return s.convertProductToResponse(product), nil
}

func (s *InventoryServer) DeleteProduct(ctx context.Context, req *gen.ProductIDRequest) (*emptypb.Empty, error) {
	if err := s.uc.Delete(req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *InventoryServer) ListProducts(ctx context.Context, req *gen.ListProductsRequest) (*gen.ListProductsResponse, error) {
	filter := make(map[string]interface{})
	for k, v := range req.Filter {
		filter[k] = v
	}

	products, err := s.uc.List(filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &gen.ListProductsResponse{
		Products: make([]*gen.ProductResponse, 0, len(products)),
	}

	for _, p := range products {
		response.Products = append(response.Products, s.convertProductToResponse(p))
	}

	return response, nil
}

func (s *InventoryServer) CheckStock(ctx context.Context, req *gen.StockCheckRequest) (*gen.StockCheckResponse, error) {
	available, err := s.uc.CheckStock(req.ProductId, int(req.Quantity))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	product, err := s.uc.GetByID(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &gen.StockCheckResponse{
		IsAvailable:  available,
		CurrentStock: int32(product.Stock),
	}, nil
}

func (s *InventoryServer) convertProductToResponse(p *domain.Product) *gen.ProductResponse {
	return &gen.ProductResponse{
		Id:       p.ID.Hex(),
		Name:     p.Name,
		Price:    float32(p.Price),
		Category: p.Category,
		Stock:    int32(p.Stock),
	}
}
