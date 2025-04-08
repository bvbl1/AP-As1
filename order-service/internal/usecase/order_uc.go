package usecase

import (
	"Assignment1_AbylayMoldakhmet/order-service/internal/clients"
	"Assignment1_AbylayMoldakhmet/order-service/internal/domain"
	"Assignment1_AbylayMoldakhmet/order-service/internal/repository/mongodb"
)

type OrderUsecase struct {
	repo            mongodb.OrderRepository
	inventoryClient clients.InventoryClientInterface
}

func NewOrderUsecase(repo mongodb.OrderRepository, inventoryClient clients.InventoryClientInterface) *OrderUsecase {
	return &OrderUsecase{
		repo:            repo,
		inventoryClient: inventoryClient,
	}
}

func (uc *OrderUsecase) Create(order *domain.Order) error {
	for _, item := range order.Items {
		available, err := uc.inventoryClient.CheckStock(item.ProductID, item.Quantity)
		if err != nil || !available {
			return domain.ErrNotEnoughStock
		}
	}

	order.Status = domain.StatusPending

	return uc.repo.Create(order)
}

func (uc *OrderUsecase) GetByID(id string) (*domain.Order, error) {
	return uc.repo.GetByID(id)
}

func (uc *OrderUsecase) UpdateStatus(id string, status domain.OrderStatus) error {
	return uc.repo.UpdateStatus(id, status)
}

func (uc *OrderUsecase) GetByUserID(userID string) ([]*domain.Order, error) {
	return uc.repo.GetByUserID(userID)
}
func (uc *OrderUsecase) GetAll() ([]*domain.Order, error) {
	return uc.repo.GetAll()
}
