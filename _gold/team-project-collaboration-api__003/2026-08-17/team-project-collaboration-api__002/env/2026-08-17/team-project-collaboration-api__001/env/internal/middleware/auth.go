package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"team-project-task-api/internal/pkg/apperror"
	"team-project-task-api/internal/pkg/token"
)

const (
	UserIDKey = "user_id"
	EmailKey  = "email"
	NameKey   = "name"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization"))
		if header == "" {
			abortWithError(c, apperror.Unauthorized("authorization header is required"))
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			abortWithError(c, apperror.Unauthorized("invalid authorization header"))
			return
		}
		claims, err := token.Parse(strings.TrimSpace(parts[1]), secret)
		if err != nil {
			abortWithError(c, apperror.Unauthorized("invalid or expired token"))
			return
		}
		c.Set(UserIDKey, claims.UserID)
		c.Set(EmailKey, claims.Email)
		c.Set(NameKey, claims.Name)
		c.Next()
	}
}

func abortWithError(c *gin.Context, err error) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		c.AbortWithStatusJSON(appErr.Status, gin.H{
			"error": gin.H{"code": appErr.Code, "message": appErr.Message},
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{"code": "internal_error", "message": "internal server error"},
	})
}
