package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"team-project-task-api/internal/middleware"
	"team-project-task-api/internal/pkg/apperror"
	"team-project-task-api/internal/pkg/httputil"
	"team-project-task-api/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, invalidRequestError(err))
		return
	}
	result, err := h.auth.Register(c.Request.Context(), req.Email, req.Name, req.Password)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusCreated, result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, invalidRequestError(err))
		return
	}
	result, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusOK, result)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetInt64(middleware.UserIDKey)
	user, err := h.auth.GetUser(c.Request.Context(), userID)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusOK, user)
}

func invalidRequestError(err error) error {
	return apperror.BadRequest("invalid request body", err)
}
