package guard

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/models/roles_model"
	"github.com/ivan/storage-project-back/internal/repository/user_repo"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
	"github.com/rs/zerolog/log"
)

const SetUsrDtoKey = "usrDTO"

func UsrGuard(usrRepo user_repo.IUserRepo, accessRoles []roles_model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtPayload := c.MustGet(SetJwtDtoKey).(*jwt_service.JwtPayload)

		usr, err := usrRepo.GetUserById(jwtPayload.ID)

		if err != nil || usr == nil {
			log.Error().Err(err).Msg("failed to get user")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid data"})
			c.Abort()
			return
		}

		if !slices.Contains(accessRoles, usr.RoleID) || usr.Blocked {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		tokenModel, err := usrRepo.GetUserAccessByToken(jwtPayload.Token)

		if err != nil || tokenModel == nil || tokenModel.Revoked {
			log.Error().Err(err).Msg("invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set(SetUsrDtoKey, usr)

		c.Next()
	}
}
