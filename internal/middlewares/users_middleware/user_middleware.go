package users_middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/user_dto"
	"github.com/ivan/storage-project-back/internal/utils/validation"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
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

	c.Set("userDTO", userDTO)

	c.Next()
}

func DeleteUserMidd(jwt *jwt_service.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(401, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.JSON(401, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)

		jwtPayload, err := jwt.ParseToken(token)

		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		dltUsrDto := user_dto.DeleteUserDto{
			UrsID: jwtPayload.ID,
		}

		c.Set("dltUsrDto", dltUsrDto)

		c.Next()
	}
}
