package fld_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/fld_dto"
	"github.com/ivan/storage-project-back/internal/utils/validation"
	"github.com/rs/zerolog/log"
)

const (
	SetDelFldDtoKey    = "delDTO"
	SetCreateFldDtoKey = "createDTO"
)

func DelFld(c *gin.Context) {
	fldID := c.Param("fldID")

	var dto fld_dto.DelFld

	parsedID, err := uuid.Parse(fldID)

	if err != nil {
		log.Error().Err(err).Msg("failed to parse fldID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse fldID"})
		c.Abort()
		return
	}

	dto.FldID = parsedID

	if err := validation.Validate.Struct(dto); err != nil {
		log.Error().Err(err).Msg("failed to validate request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to validate request"})
		c.Abort()
		return
	}

	c.Set(SetDelFldDtoKey, dto)

	c.Next()
}

func CreateFld(c *gin.Context) {
	var dto fld_dto.CreateFldReq

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

	c.Set(SetCreateFldDtoKey, dto)

	c.Next()
}
