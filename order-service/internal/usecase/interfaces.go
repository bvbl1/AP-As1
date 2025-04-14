package usecase

import "Assignment1_AbylayMoldakhmet/order-service/internal/domain"

type OrderUsecaseInterface interface {
	Create(order *domain.Order) error
	GetByID(id string) (*domain.Order, error)
	UpdateStatus(id string, status domain.OrderStatus) error
	GetByUserID(userID string) ([]*domain.Order, error)
	GetAll() ([]*domain.Order, error)
}
