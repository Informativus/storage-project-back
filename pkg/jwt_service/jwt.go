package jwt_service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
    "github.com/ivan/storage-project-back/pkg/config"
)

type JwtPayload struct {
	FolderName string
}

type JwtService struct {
	cfg *config.Config
}

func NewJwtService(cfg *config.Config) *JwtService {
	return &JwtService{cfg: cfg}
}

func (s *JwtService) GenerateToken(payload JwtPayload) (string, error) {
	claims := jwt.MapClaims{
		"folderName": payload.FolderName,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.cfg.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *JwtService) VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.SecretKey), nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}
