package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"team-project-task-api/internal/middleware"
	"team-project-task-api/internal/model"
	"team-project-task-api/internal/pkg/httputil"
	"team-project-task-api/internal/service"
)

type TaskHandler struct {
	tasks *service.TaskService
}

func NewTaskHandler(tasks *service.TaskService) *TaskHandler {
	return &TaskHandler{tasks: tasks}
}

type createTaskRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"`
	Status      string   `json:"status"`
	AssigneeID  *int64   `json:"assignee_id"`
	DueDate     string   `json:"due_date"`
	Labels      []string `json:"labels"`
}

type updateTaskRequest struct {
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Priority    *string   `json:"priority"`
	Status      *string   `json:"status"`
	AssigneeID  *int64    `json:"assignee_id"`
	DueDate     *string   `json:"due_date"`
	Labels      *[]string `json:"labels"`
}

type updateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type updateAssigneeRequest struct {
	AssigneeID int64 `json:"assignee_id" binding:"required"`
}

type addCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *TaskHandler) Create(c *gin.Context) {
	projectID, ok := parseID(c, "projectID")
	if !ok {
		badID(c)
		return
	}
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, invalidRequestError(err))
		return
	}
	dueDate, err := parseDate(req.DueDate)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	if req.Status == "" {
		req.Status = model.TaskStatusTodo
	}
	if req.Priority == "" {
		req.Priority = model.TaskPriorityMedium
	}
	if req.Labels == nil {
		req.Labels = []string{}
	}
	task, err := h.tasks.CreateTask(c.Request.Context(), c.GetInt64(middleware.UserIDKey), projectID, model.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		AssigneeID:  req.AssigneeID,
		DueDate:     dueDate,
		Labels:      req.Labels,
	})
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusCreated, task)
}

func (h *TaskHandler) List(c *gin.Context) {
	projectID, ok := parseID(c, "projectID")
	if !ok {
		badID(c)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	assigneeID, _ := strconv.ParseInt(c.Query("assignee_id"), 10, 64)
	status := strings.TrimSpace(c.Query("status"))
	tasks, total, err := h.tasks.ListTasks(c.Request.Context(), c.GetInt64(middleware.UserIDKey), projectID, model.TaskFilter{
		Status:     status,
		AssigneeID: assigneeID,
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusOK, gin.H{
		"items":       tasks,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages(total, pageSize),
	})
}

func (h *TaskHandler) Get(c *gin.Context) {
	projectID, ok := parseID(c, "projectID")
	if !ok {
		badID(c)
		return
	}
	taskID, ok := parseID(c, "taskID")
	if !ok {
		badID(c)
		return
	}
	task, err := h.tasks.GetTask(c.Request.Context(), c.GetInt64(middleware.UserIDKey), projectID, taskID)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	projectID, ok := parseID(c, "projectID")
	if !ok {
		badID(c)
		return
	}
	taskID, ok := parseID(c, "taskID")
	if !ok {
		badID(c)
		return
	}
	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, invalidRequestError(err))
		return
	}
	input := model.UpdateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		AssigneeID:  req.AssigneeID,
		Labels:      req.Labels,
	}
	if req.DueDate != nil {
		dueDate, err := parseDate(*req.DueDate)
		if err != nil {
			httputil.Error(c, err)
			return
		}
		input.DueDate = dueDate
	}
	task, err := h.tasks.UpdateTask(c.Request.Context(), c.GetInt64(middleware.UserIDKey), projectID, taskID, input)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	projectID, ok := parseID(c, "projectID")
	if !ok {
		badID(c)
		return
	}
	taskID, ok := parseID(c, "taskID")
	if !ok {
		badID(c)
		return
	}
	if err := h.tasks.DeleteTask(c.Request.Context(), c.GetInt64(middleware.UserIDKey), projectID, taskID); err != nil {
		httputil.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	projectID, ok := parseID(c, "projectID")
	if !ok {
		badID(c)
		return
	}
	taskID, ok := parseID(c, "taskID")
	if !ok {
		badID(c)
		return
	}
	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, invalidRequestError(err))
		return
	}
	task, err := h.tasks.UpdateStatus(c.Request.Context(), c.GetInt64(middleware.UserIDKey), projectID, taskID, req.Status)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusOK, task)
}

func (h *TaskHandler) UpdateAssignee(c *gin.Context) {
	projectID, ok := parseID(c, "projectID")
	if !ok {
		badID(c)
		return
	}
	taskID, ok := parseID(c, "taskID")
	if !ok {
		badID(c)
		return
	}
	var req updateAssigneeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, invalidRequestError(err))
		return
	}
	task, err := h.tasks.UpdateAssignee(c.Request.Context(), c.GetInt64(middleware.UserIDKey), projectID, taskID, req.AssigneeID)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusOK, task)
}

func (h *TaskHandler) AddComment(c *gin.Context) {
	projectID, ok := parseID(c, "projectID")
	if !ok {
		badID(c)
		return
	}
	taskID, ok := parseID(c, "taskID")
	if !ok {
		badID(c)
		return
	}
	var req addCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, invalidRequestError(err))
		return
	}
	comment, err := h.tasks.AddComment(c.Request.Context(), c.GetInt64(middleware.UserIDKey), projectID, taskID, req.Content)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusCreated, comment)
}

func (h *TaskHandler) ListComments(c *gin.Context) {
	projectID, ok := parseID(c, "projectID")
	if !ok {
		badID(c)
		return
	}
	taskID, ok := parseID(c, "taskID")
	if !ok {
		badID(c)
		return
	}
	comments, err := h.tasks.ListComments(c.Request.Context(), c.GetInt64(middleware.UserIDKey), projectID, taskID)
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusOK, gin.H{"items": comments, "total": len(comments)})
}

func (h *TaskHandler) ListOverdue(c *gin.Context) {
	tasks, err := h.tasks.ListOverdueTasks(c.Request.Context(), c.GetInt64(middleware.UserIDKey))
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.JSON(c, http.StatusOK, gin.H{"items": tasks, "total": len(tasks)})
}

func totalPages(total int64, pageSize int) int {
	if pageSize < 1 {
		pageSize = 20
	}
	if total <= 0 {
		return 0
	}
	pages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		pages++
	}
	return pages
}
