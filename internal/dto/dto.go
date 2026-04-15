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

// Tasks
type CreateTaskRequest struct {
    Title       string     `json:"title"`
    Description *string    `json:"description"`
    Priority    string     `json:"priority"`
    AssigneeID  *uuid.UUID `json:"assignee_id"`
    DueDate     *string    `json:"due_date"` // "2026-04-15" parsed in handler
}

type UpdateTaskRequest struct {
    Title       *string    `json:"title"`
    Description *string    `json:"description"`
    Status      *string    `json:"status"`
    Priority    *string    `json:"priority"`
    AssigneeID  *uuid.UUID `json:"assignee_id"`
    DueDate     *string    `json:"due_date"`
}

type TaskResponse struct {
    ID          uuid.UUID  `json:"id"`
    Title       string     `json:"title"`
    Description *string    `json:"description,omitempty"`
    Status      string     `json:"status"`
    Priority    string     `json:"priority"`
    ProjectID   uuid.UUID  `json:"project_id"`
    AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
    DueDate     *string    `json:"due_date,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}