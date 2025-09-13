package guard

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
	"github.com/rs/zerolog/log"
)

const SetJwtDtoKey = "jwtDTO"

func JwtGuard(jwt *jwt_service.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			log.Error().Msg("authorization header required")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			log.Error().Msg("invalid authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)

		jwtPayload, err := jwt.ParseToken(token)

		if err != nil {
			log.Error().Err(err).Msg("failed to parse token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set(SetJwtDtoKey, jwtPayload)

		c.Next()
	}
}
