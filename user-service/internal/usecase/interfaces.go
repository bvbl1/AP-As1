package usecase

import (
	"Assignment1_AbylayMoldakhmet/user-service/internal/domain"
)

type UserUsecase interface {
	GetByID(id string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
}

type AuthUsecase interface {
	Register(email, password string) (*domain.User, error)
	Login(email, password string) (string, error)
	GenerateToken(userID string) (string, error)
}
