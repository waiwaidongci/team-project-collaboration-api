package model

import "time"

const (
	TaskStatusTodo       = "todo"
	TaskStatusInProgress = "in_progress"
	TaskStatusDone       = "done"

	TaskPriorityLow    = "low"
	TaskPriorityMedium = "medium"
	TaskPriorityHigh   = "high"
	TaskPriorityUrgent = "urgent"

	ProjectRoleOwner  = "owner"
	ProjectRoleMember = "member"

	ActivityTaskCreated     = "task_created"
	ActivityTaskUpdated     = "task_updated"
	ActivityTaskDeleted     = "task_deleted"
	ActivityStatusChanged   = "status_changed"
	ActivityAssigneeChanged = "assignee_changed"
	ActivityCommentAdded    = "comment_added"
)

var ValidTaskStatuses = map[string]bool{
	TaskStatusTodo:       true,
	TaskStatusInProgress: true,
	TaskStatusDone:       true,
}

var ValidTaskPriorities = map[string]bool{
	TaskPriorityLow:    true,
	TaskPriorityMedium: true,
	TaskPriorityHigh:   true,
	TaskPriorityUrgent: true,
}

var ValidProjectRoles = map[string]bool{
	ProjectRoleOwner:  true,
	ProjectRoleMember: true,
}

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Project struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     int64     `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectMember struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"project_id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	ID          int64      `json:"id"`
	ProjectID   int64      `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	Status      string     `json:"status"`
	AssigneeID  *int64     `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
	Labels      []string   `json:"labels"`
	CreatedBy   int64      `json:"created_by"`
	ProjectName string     `json:"project_name,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Comment struct {
	ID         int64     `json:"id"`
	TaskID     int64     `json:"task_id"`
	UserID     int64     `json:"user_id"`
	AuthorName string    `json:"author_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type Activity struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"project_id"`
	TaskID    *int64    `json:"task_id"`
	UserID    int64     `json:"user_id"`
	UserName  string    `json:"user_name"`
	Action    string    `json:"action"`
	OldValue  *string   `json:"old_value,omitempty"`
	NewValue  *string   `json:"new_value,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type PagedTasks struct {
	Items      []Task `json:"items"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

type CreateTaskInput struct {
	Title       string
	Description string
	Priority    string
	Status      string
	AssigneeID  *int64
	DueDate     *time.Time
	Labels      []string
}

type UpdateTaskInput struct {
	Title       *string
	Description *string
	Priority    *string
	Status      *string
	AssigneeID  *int64
	DueDate     *time.Time
	Labels      *[]string
}

type TaskFilter struct {
	Status     string
	AssigneeID int64
	Page       int
	PageSize   int
}
