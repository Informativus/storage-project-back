package error_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/pkg/errsvc"
	"github.com/rs/zerolog/log"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		err := c.Errors.Last().Err

		if appErr, ok := err.(*errsvc.AppError); ok {
			log.Error().Str("trace", appErr.Trace()).Msg(appErr.Error())

			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Key,
				"message": appErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Internal server error",
		})
	}
}
