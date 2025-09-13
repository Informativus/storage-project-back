package guard

import (
	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/models/protection_group"
	"github.com/ivan/storage-project-back/internal/services/security_service"
)

func FldGuard(securityService *security_service.SecurityService, accessGroups []protection_group.ProtectionGroupsType) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
