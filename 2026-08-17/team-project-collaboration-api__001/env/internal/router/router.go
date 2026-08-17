package router

import (
	"github.com/gin-gonic/gin"

	"team-project-task-api/internal/handler"
	"team-project-task-api/internal/middleware"
)

func New(h *handler.Handler, jwtSecret string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/register", h.Auth.Register)
	auth.POST("/login", h.Auth.Login)

	authed := api.Group("")
	authed.Use(middleware.Auth(jwtSecret))
	authed.GET("/auth/me", h.Auth.Me)

	projects := authed.Group("/projects")
	projects.POST("", h.Projects.Create)
	projects.GET("", h.Projects.List)
	projects.GET("/:projectID", h.Projects.Get)
	projects.POST("/:projectID/members", h.Projects.InviteMember)
	projects.GET("/:projectID/members", h.Projects.ListMembers)
	projects.GET("/:projectID/activities", h.Projects.ListActivities)

	projects.POST("/:projectID/tasks", h.Tasks.Create)
	projects.GET("/:projectID/tasks", h.Tasks.List)
	projects.GET("/:projectID/tasks/:taskID", h.Tasks.Get)
	projects.PATCH("/:projectID/tasks/:taskID", h.Tasks.Update)
	projects.DELETE("/:projectID/tasks/:taskID", h.Tasks.Delete)
	projects.PATCH("/:projectID/tasks/:taskID/status", h.Tasks.UpdateStatus)
	projects.PATCH("/:projectID/tasks/:taskID/assignee", h.Tasks.UpdateAssignee)
	projects.POST("/:projectID/tasks/:taskID/comments", h.Tasks.AddComment)
	projects.GET("/:projectID/tasks/:taskID/comments", h.Tasks.ListComments)

	authed.GET("/tasks/overdue", h.Tasks.ListOverdue)

	return r
}
