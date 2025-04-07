package usecase

import (
	"Assignment1_AbylayMoldakhmet/user-service/internal/domain"
	"Assignment1_AbylayMoldakhmet/user-service/internal/repository"
	"Assignment1_AbylayMoldakhmet/user-service/pkg"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// private реализация интерфейса AuthUsecase
type authUsecase struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthUsecase(repo repository.UserRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (uc *authUsecase) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

func (uc *authUsecase) Register(email, password string) (*domain.User, error) {
	if _, err := uc.repo.GetByEmail(email); err == nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := pkg.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:       primitive.NewObjectID(),
		Email:    email,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := uc.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *authUsecase) Login(email, password string) (string, error) {
	user, err := uc.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials (email)")
	}

	// Добавьте проверку ID
	if user.ID.IsZero() {
		return "", errors.New("user ID is empty")
	}

	if !pkg.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials (password)")
	}

	return uc.GenerateToken(user.ID.Hex())
}

// private реализация интерфейса
type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (uc *userUsecase) GetByID(id string) (*domain.User, error) {
	return uc.repo.GetByID(id)
}

func (uc *userUsecase) Update(user *domain.User) error {
	// Бизнес-логика при необходимости
	return uc.repo.Update(user)
}

func (uc *userUsecase) Delete(id string) error {
	return uc.repo.Delete(id)
}
