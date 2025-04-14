package grpc

import (
	"Assignment1_AbylayMoldakhmet/order-service/internal/domain"
	"Assignment1_AbylayMoldakhmet/order-service/internal/usecase"
	"Assignment1_AbylayMoldakhmet/proto/gen"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServer struct {
	gen.UnimplementedOrderServiceServer
	uc usecase.OrderUsecaseInterface // Исправленный тип
}

func NewOrderServer(uc usecase.OrderUsecaseInterface) *OrderServer {
	return &OrderServer{uc: uc}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *gen.CreateOrderRequest) (*gen.OrderResponse, error) {
	order := &domain.Order{
		UserID: req.UserId,
		Items:  convertProtoItems(req.Items),
		Status: domain.StatusPending, // Статус задаем явно
	}

	if err := s.uc.Create(order); err != nil {
		if errors.Is(err, domain.ErrNotEnoughStock) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertOrderToResponse(order), nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *gen.OrderIDRequest) (*gen.OrderResponse, error) {
	order, err := s.uc.GetByID(req.Id)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertOrderToResponse(order), nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *gen.OrderStatusUpdateRequest) (*emptypb.Empty, error) {
	newStatus := convertProtoStatus(req.Status) // Исправленное имя

	if err := s.uc.UpdateStatus(req.Id, newStatus); err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *OrderServer) ListOrders(ctx context.Context, req *gen.OrderFilterRequest) (*gen.OrderListResponse, error) {
	var orders []*domain.Order
	var err error

	if req.UserId != nil && *req.UserId != "" {
		orders, err = s.uc.GetByUserID(*req.UserId)
	} else {
		orders, err = s.uc.GetAll()
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.OrderListResponse{
		Orders: convertOrdersToProto(orders),
	}, nil
}

func convertOrderToResponse(o *domain.Order) *gen.OrderResponse {
	return &gen.OrderResponse{
		Id:     o.ID.Hex(),
		UserId: o.UserID,
		Items:  convertDomainItems(o.Items),
		Status: convertDomainStatus(o.Status),
	}
}

func convertDomainStatus(s domain.OrderStatus) gen.OrderStatus {
	switch s {
	case domain.StatusPending:
		return gen.OrderStatus_PENDING
	case domain.StatusPaid:
		return gen.OrderStatus_PAID
	case domain.StatusCancelled:
		return gen.OrderStatus_CANCELLED
	default:
		return gen.OrderStatus_PENDING
	}
}

func convertProtoStatus(s gen.OrderStatus) domain.OrderStatus {
	switch s {
	case gen.OrderStatus_PENDING:
		return domain.StatusPending
	case gen.OrderStatus_PAID:
		return domain.StatusPaid
	case gen.OrderStatus_CANCELLED:
		return domain.StatusCancelled
	default:
		return domain.StatusPending
	}
}

func convertProtoItems(items []*gen.OrderItem) []domain.OrderItem {
	result := make([]domain.OrderItem, 0, len(items))
	for _, item := range items {
		result = append(result, domain.OrderItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
			Price:     item.Price,
		})
	}
	return result
}

func convertDomainItems(items []domain.OrderItem) []*gen.OrderItem {
	result := make([]*gen.OrderItem, 0, len(items))
	for _, item := range items {
		result = append(result, &gen.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		})
	}
	return result
}

func convertOrdersToProto(orders []*domain.Order) []*gen.OrderResponse {
	result := make([]*gen.OrderResponse, 0, len(orders))
	for _, o := range orders {
		result = append(result, convertOrderToResponse(o))
	}
	return result
}
