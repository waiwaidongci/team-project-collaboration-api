package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"team-project-task-api/internal/pkg/apperror"
	"team-project-task-api/internal/service"
)

type Handler struct {
	Auth     *AuthHandler
	Projects *ProjectHandler
	Tasks    *TaskHandler
}

func New(authService *service.AuthService, projectService *service.ProjectService, taskService *service.TaskService) *Handler {
	return &Handler{
		Auth:     NewAuthHandler(authService),
		Projects: NewProjectHandler(projectService),
		Tasks:    NewTaskHandler(taskService),
	}
}

func parseID(c *gin.Context, name string) (int64, bool) {
	raw := c.Param(name)
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

func badID(c *gin.Context) {
	c.JSON(400, gin.H{"error": gin.H{"code": "bad_request", "message": "invalid id"}})
}

func parseDate(raw string) (*time.Time, error) {
	if raw == "" {
		return nil, nil
	}
	value, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return nil, apperror.BadRequest("due_date must be in YYYY-MM-DD format", nil)
	}
	return &value, nil
}
