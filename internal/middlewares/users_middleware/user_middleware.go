package users_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/user_dto"
	"github.com/ivan/storage-project-back/internal/utils/validation"
	"github.com/rs/zerolog/log"
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
		UrsName: dto.UrsName,
	}

	c.Set("createUserDTO", userDTO)

	c.Next()
}
