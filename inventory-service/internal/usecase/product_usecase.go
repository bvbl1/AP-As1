package usecase

import (
	"Assignment1_AbylayMoldakhmet/inventory-service/internal/domain"
	"errors"
)

type productUsecase struct {
	repo domain.ProductRepository
}

func NewProductUsecase(repo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{repo: repo}
}

func (uc *productUsecase) Create(product *domain.Product) error {
	if product.Name == "" || product.Price <= 0 || product.Stock < 0 {
		return errors.New("invalid product data")
	}
	return uc.repo.Create(product)
}

func (uc *productUsecase) GetByID(id string) (*domain.Product, error) {
	return uc.repo.GetByID(id)
}

func (uc *productUsecase) Update(product *domain.Product) error {
	if product.ID.IsZero() {
		return errors.New("invalid product ID")
	}
	return uc.repo.Update(product)
}

func (uc *productUsecase) Delete(id string) error {
	return uc.repo.Delete(id)
}

func (uc *productUsecase) List(filter map[string]interface{}) ([]*domain.Product, error) {
	return uc.repo.List(filter)
}

func (uc *productUsecase) CheckStock(productID string, quantity int) (bool, error) {
	product, err := uc.repo.GetByID(productID) // Используем метод репозитория
	if err != nil {
		return false, err
	}
	return product.Stock >= quantity, nil // Проверяем хватает ли товара
}
