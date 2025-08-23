package jwt_service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/rs/zerolog/log"
)

type JwtPayload struct {
	ID uuid.UUID
}

type JwtService struct {
	cfg *config.Config
}

func NewJwtService(cfg *config.Config) *JwtService {
	return &JwtService{cfg: cfg}
}

func (s *JwtService) GenerateToken(payload JwtPayload) (string, error) {
	claims := jwt.MapClaims{
		"id":  payload.ID,
		"exp": time.Now().Unix() + s.cfg.ExpiresIn,
		"iat": time.Now().Unix(),
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

func (s *JwtService) ParseToken(tokenString string) (*JwtPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error().Str("alg", fmt.Sprintf("%v", token.Header["alg"])).
				Msg("unexpected signing method")
			return nil, fmt.Errorf("internal server error")
		}
		return []byte(s.cfg.SecretKey), nil
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to parse token")
		return nil, fmt.Errorf("internal server error")
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Error().Msg("invalid token claims type")
		return nil, fmt.Errorf("internal server error")
	}

	idStr, ok := claims["id"].(string)
	if !ok {
		log.Error().Interface("claims", claims).Msg("missing or invalid id claim")
		return nil, fmt.Errorf("internal server error")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("idStr", idStr).Msg("failed to parse UUID from token")
		return nil, fmt.Errorf("internal server error")
	}

	return &JwtPayload{ID: id}, nil
}
