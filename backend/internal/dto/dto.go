package dto


import (
    "time"
    "github.com/google/uuid"
)

type RegisterRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthResponse struct {
    Token string   `json:"token"`
    User  UserDTO  `json:"user"`
}

type UserDTO struct {
    ID    uuid.UUID `json:"id"`
    Name  string    `json:"name"`
    Email string    `json:"email"`
}

// Projects
type CreateProjectRequest struct {
    Name        string  `json:"name"`
    Description *string `json:"description"`
}

type UpdateProjectRequest struct {
    Name        *string `json:"name"`
    Description *string `json:"description"`
}

type ProjectResponse struct {
    ID          uuid.UUID  `json:"id"`
    Name        string     `json:"name"`
    Description *string    `json:"description,omitempty"`
    OwnerID     uuid.UUID  `json:"owner_id"`
    CreatedAt   time.Time  `json:"created_at"`
}

type ProjectDetailResponse struct {
    ProjectResponse
    Tasks []TaskResponse `json:"tasks"`
}

type ProjectStatsResponse struct {
    ByStatus   map[string]int `json:"by_status"`
    ByAssignee map[string]int `json:"by_assignee"`
}

// Tasks
type CreateTaskRequest struct {
    Title       string      `json:"title"`
    Description *string     `json:"description,omitempty"`
    Status      *string     `json:"status,omitempty"`
    Priority    *string     `json:"priority,omitempty"`
    AssigneeID  *uuid.UUID  `json:"assignee_id,omitempty"`
    DueDate     *time.Time  `json:"due_date,omitempty"`
}

type UpdateTaskRequest struct {
    Title       *string     `json:"title,omitempty"`
    Description *string     `json:"description,omitempty"`
    Status      *string     `json:"status,omitempty"`
    Priority    *string     `json:"priority,omitempty"`
    AssigneeID  *uuid.UUID  `json:"assignee_id,omitempty"`
    DueDate     *time.Time  `json:"due_date,omitempty"`
}

type TaskResponse struct {
    ID          uuid.UUID  `json:"id"`
    Title       string     `json:"title"`
    Description *string    `json:"description,omitempty"`
    Status      string     `json:"status"`
    Priority    string     `json:"priority"`
    ProjectID   uuid.UUID  `json:"project_id"`
    AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
    DueDate     *time.Time `json:"due_date,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}
