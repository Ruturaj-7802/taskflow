package repository

import (
    "context"
    "fmt"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/Ruturaj-7802/taskflow/internal/model"
)

type TaskRepository struct {
    db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
    return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, t *model.Task) error {
    _, err := r.db.Exec(ctx,
        `INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
        t.ID, t.Title, t.Description, t.Status, t.Priority,
        t.ProjectID, t.AssigneeID, t.DueDate, t.CreatedAt, t.UpdatedAt,
    )
    return err
}

type TaskFilter struct {
    Status     *string
    AssigneeID *uuid.UUID
    Page       int
    Limit      int
}

func (r *TaskRepository) ListByProject(ctx context.Context, projectID uuid.UUID, f TaskFilter) ([]model.Task, error) {
    query := `SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
              FROM tasks WHERE project_id = $1`
    args := []any{projectID}
    argIdx := 2

    if f.Status != nil {
        query += fmt.Sprintf(" AND status = $%d", argIdx)
        args = append(args, *f.Status)
        argIdx++
    }
    if f.AssigneeID != nil {
        query += fmt.Sprintf(" AND assignee_id = $%d", argIdx)
        args = append(args, *f.AssigneeID)
        argIdx++
    }

    query += " ORDER BY created_at DESC"

    // Pagination (bonus)
    if f.Limit > 0 {
        query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
        args = append(args, f.Limit, (f.Page-1)*f.Limit)
    }

    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []model.Task
    for rows.Next() {
        var t model.Task
        if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
            &t.ProjectID, &t.AssigneeID, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
            return nil, err
        }
        tasks = append(tasks, t)
    }
    return tasks, rows.Err()
}

func (r *TaskRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
    t := &model.Task{}
    err := r.db.QueryRow(ctx,
        `SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
         FROM tasks WHERE id = $1`, id,
    ).Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
        &t.ProjectID, &t.AssigneeID, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
    if err != nil {
        return nil, err
    }
    return t, nil
}

func (r *TaskRepository) Update(ctx context.Context, t *model.Task) error {
    _, err := r.db.Exec(ctx,
        `UPDATE tasks SET title=$1, description=$2, status=$3, priority=$4,
         assignee_id=$5, due_date=$6, updated_at=$7 WHERE id=$8`,
        t.Title, t.Description, t.Status, t.Priority,
        t.AssigneeID, t.DueDate, t.UpdatedAt, t.ID,
    )
    return err
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
    _, err := r.db.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
    return err
}

func (r *TaskRepository) DeleteByProject(ctx context.Context, projectID uuid.UUID) error {
    _, err := r.db.Exec(ctx, `DELETE FROM tasks WHERE project_id = $1`, projectID)
    return err
}

// Stats (bonus)
func (r *TaskRepository) StatsByProject(ctx context.Context, projectID uuid.UUID) (map[string]int, map[string]int, error) {
    byStatus := map[string]int{"todo": 0, "in_progress": 0, "done": 0}
    rows, err := r.db.Query(ctx,
        `SELECT status, COUNT(*) FROM tasks WHERE project_id = $1 GROUP BY status`,
        projectID)
    if err != nil {
        return nil, nil, err
    }
    defer rows.Close()
    for rows.Next() {
        var s string
        var c int
        rows.Scan(&s, &c)
        byStatus[s] = c
    }

    byAssignee := map[string]int{}
    rows2, err := r.db.Query(ctx,
        `SELECT COALESCE(assignee_id::text, 'unassigned'), COUNT(*) FROM tasks WHERE project_id = $1 GROUP BY assignee_id`,
        projectID)
    if err != nil {
        return nil, nil, err
    }
    defer rows2.Close()
    for rows2.Next() {
        var a string
        var c int
        rows2.Scan(&a, &c)
        byAssignee[a] = c
    }

    return byStatus, byAssignee, nil
}