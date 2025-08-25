package guard

import (
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/models/roles_model"
	"github.com/ivan/storage-project-back/internal/repository/user_repo"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
	"github.com/rs/zerolog/log"
)

func AuthGuard(jwt *jwt_service.JwtService, usrRepo *user_repo.UserRepo, accessRoles []roles_model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			log.Error().Msg("authorization header required")
			c.JSON(401, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			log.Error().Msg("invalid authorization header")
			c.JSON(401, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)

		jwtPayload, err := jwt.ParseToken(token)

		if err != nil {
			log.Error().Err(err).Msg("failed to parse token")
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		usr, err := usrRepo.GetUserById(jwtPayload.ID)

		if err != nil || usr == nil {
			log.Error().Err(err).Msg("failed to get user")
			c.JSON(401, gin.H{"error": "invalid data"})
			c.Abort()
			return
		}

		if !slices.Contains(accessRoles, usr.RoleID) || usr.Blocked {
			c.JSON(403, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		usrAccess, err := usrRepo.GetUserAccessById(usr.ID)

		if err != nil || usrAccess == nil {
			log.Error().Err(err).Msg("failed to get user access")
			c.JSON(401, gin.H{"error": "invalid data"})
			c.Abort()
			return
		}

		if usrAccess.Revoked {
			c.JSON(403, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		c.Set("usrDTO", usr)

		c.Next()
	}
}
