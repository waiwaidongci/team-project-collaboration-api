package httputil

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"team-project-task-api/internal/pkg/apperror"
)

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

func Error(c *gin.Context, err error) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		c.JSON(appErr.Status, gin.H{
			"error": gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":    "internal_error",
			"message": "internal server error",
		},
	})
}
