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
	var bodyDto fld_dto.CreateFldBody

	if err := c.ShouldBindJSON(&bodyDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request should be JSON"})
		c.Abort()
		return
	}

	if err := validation.Validate.Struct(bodyDto); err != nil {
		log.Error().Err(err).Msg("failed to validate request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to validate request"})
		c.Abort()
		return
	}

	parsedID, err := uuid.Parse(c.Param("fldID"))

	if err != nil {
		log.Error().Err(err).Msg("failed to parse fldID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse fldID"})
		c.Abort()
		return
	}

	dto := fld_dto.CreateFldDto{
		Name:     bodyDto.Name,
		ParentID: parsedID,
	}
	log.Info().Interface("dto", dto).Msg("dto")

	c.Set(SetCreateFldDtoKey, dto)

	c.Next()
}
