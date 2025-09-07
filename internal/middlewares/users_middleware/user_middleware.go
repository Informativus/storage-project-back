package users_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/user_dto"
	"github.com/ivan/storage-project-back/internal/utils/validation"
	"github.com/rs/zerolog/log"
)

const (
	SetTokenDtoKey      = "tokenDTO"
	SetCreateUserDtoKey = "createUserDTO"
	SetBlockUserDtoKey  = "blockUserDTO"
)

func CreateUserMidd(c *gin.Context) {
	var dto user_dto.CreateUserDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request should be JSON"})
		c.Abort()
		return
	}

	if err := validation.Validate.Struct(dto); err != nil {
		log.Error().Err(err).Msg("failed to validate request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to validate request"})
		c.Abort()
		return
	}

	userDTO := user_dto.CreateUserDto{
		UsrName: dto.UsrName,
	}

	c.Set(SetCreateUserDtoKey, userDTO)

	c.Next()
}

func GetTokenMidd(c *gin.Context) {
	var dto user_dto.GenTokenReq

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request should be JSON"})
		c.Abort()
		return
	}

	if err := validation.Validate.Struct(dto); err != nil {
		log.Error().Err(err).Msg("failed to validate request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to validate request"})
		c.Abort()
		return
	}

	tokenDTO := user_dto.GenTokenReq{
		UsrName: dto.UsrName,
	}

	c.Set(SetTokenDtoKey, tokenDTO)

	c.Next()
}

func BlockUserMiddleware(c *gin.Context) {
	var dto user_dto.BlockUserReq
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		c.Abort()
		return
	}

	if err := validation.Validate.Struct(dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed"})
		c.Abort()
		return
	}

	c.Set(SetBlockUserDtoKey, dto)
	c.Next()
}
