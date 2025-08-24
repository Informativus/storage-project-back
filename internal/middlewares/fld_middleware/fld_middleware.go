package fld_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/fld_dto"
	"github.com/ivan/storage-project-back/internal/utils/validation"
	"github.com/rs/zerolog/log"
)

func DelFld(c *gin.Context) {
	fldName := c.Param("fldName")

	var dto fld_dto.DelFld
	dto.Name = fldName

	if err := validation.Validate.Struct(dto); err != nil {
		log.Error().Err(err).Msg("failed to validate request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to validate request"})
		c.Abort()
		return
	}

	c.Set("dltDTO", dto)

	c.Next()
}
