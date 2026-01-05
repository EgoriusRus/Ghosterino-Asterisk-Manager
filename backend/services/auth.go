package services

import (
	"os"
	"time"

	"asterisk-manager/domain"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims содержит claims для JWT токена
type JWTClaims struct {
	UserID   uint            `json:"userId"`
	Username string          `json:"username"`
	Role     domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// AuthService сервис авторизации
type AuthService struct {
	secretKey     []byte
	tokenDuration time.Duration
}

// NewAuthService создает новый сервис авторизации
func NewAuthService() *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "asterisk-manager-secret-key-change-in-production"
	}

	return &AuthService{
		secretKey:     []byte(secret),
		tokenDuration: 24 * time.Hour,
	}
}

// GenerateToken генерирует JWT токен для пользователя
func (s *AuthService) GenerateToken(user *domain.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken проверяет JWT токен и возвращает claims
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
