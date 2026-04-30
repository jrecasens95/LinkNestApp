package domain

import (
	"time"

	"link-nest/internal/platform/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct{}

func NewJWTService() JWTService {
	return JWTService{}
}

func (JWTService) GenerateToken(subject string) (string, error) {
	claims := jwt.MapClaims{
		"sub": subject,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Current.JWTSecret))
}
