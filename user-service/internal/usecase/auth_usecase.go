package usecase

import (
	"Assignment1_AbylayMoldakhmet/user-service/internal/repository"
	"Assignment1_AbylayMoldakhmet/user-service/pkg"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUsecase struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthUsecase(repo repository.UserRepository, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// Генерация JWT токена
func (uc *AuthUsecase) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Токен на 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

// Проверка учетных данных и выдача токена
func (uc *AuthUsecase) Login(email, password string) (string, error) {
	// Нам все равно нужен GetByEmail для логина!
	user, err := uc.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !pkg.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return uc.GenerateToken(user.ID.Hex())
}
