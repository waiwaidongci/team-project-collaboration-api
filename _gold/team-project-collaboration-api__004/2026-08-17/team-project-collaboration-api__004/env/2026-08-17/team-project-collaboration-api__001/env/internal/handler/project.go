package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"team-project-task-api/internal/middleware"
	"team-project-task-api/internal/pkg/httputil"
	"team-project-task-api/internal/service"
)

type ProjectHandler struct {
	projects *service.ProjectService
}

func NewProjectHandler(projects *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projects: projects}
}

type createProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type inviteMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role"`
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, invalidRequestError(err))
		return
	}
	project, err := h.projects.CreateProject(c.Request.Context(), c.GetInt64(middleware.UserIDKey), req.Name, req.Description)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusCreated, project)
}

func (h *ProjectHandler) List(c *gin.Context) {
	projects, err := h.projects.ListProjects(c.Request.Context(), c.GetInt64(middleware.UserIDKey))
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusOK, gin.H{"items": projects, "total": len(projects)})
}

func (h *ProjectHandler) Get(c *gin.Context) {
	projectID, ok := parseID(c, "projectID")
	if !ok {
		badID(c)
		return
	}
	project, err := h.projects.GetProject(c.Request.Context(), c.GetInt64(middleware.UserIDKey), projectID)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusOK, project)
}

func (h *ProjectHandler) InviteMember(c *gin.Context) {
	projectID, ok := parseID(c, "projectID")
	if !ok {
		badID(c)
		return
	}
	var req inviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, invalidRequestError(err))
		return
	}
	member, err := h.projects.InviteMember(c.Request.Context(), c.GetInt64(middleware.UserIDKey), projectID, req.Email, req.Role)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusCreated, member)
}

func (h *ProjectHandler) ListMembers(c *gin.Context) {
	projectID, ok := parseID(c, "projectID")
	if !ok {
		badID(c)
		return
	}
	members, err := h.projects.ListMembers(c.Request.Context(), c.GetInt64(middleware.UserIDKey), projectID)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusOK, gin.H{"items": members, "total": len(members)})
}

func (h *ProjectHandler) ListActivities(c *gin.Context) {
	projectID, ok := parseID(c, "projectID")
	if !ok {
		badID(c)
		return
	}
	activities, err := h.projects.ListActivities(c.Request.Context(), c.GetInt64(middleware.UserIDKey), projectID)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusOK, gin.H{"items": activities, "total": len(activities)})
}
