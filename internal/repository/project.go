package repository

import (
    "context"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/ruturaj/taskflow/internal/model"
)

type ProjectRepository struct {
    db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
    return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, p *model.Project) error {
    _, err := r.db.Exec(ctx,
        `INSERT INTO projects (id, name, description, owner_id, created_at)
         VALUES ($1, $2, $3, $4, $5)`,
        p.ID, p.Name, p.Description, p.OwnerID, p.CreatedAt,
    )
    return err
}

// ListForUser returns projects where user is owner OR has assigned tasks
func (r *ProjectRepository) ListForUser(ctx context.Context, userID uuid.UUID) ([]model.Project, error) {
    rows, err := r.db.Query(ctx,
        `SELECT DISTINCT p.id, p.name, p.description, p.owner_id, p.created_at
         FROM projects p
         LEFT JOIN tasks t ON t.project_id = p.id
         WHERE p.owner_id = $1 OR t.assignee_id = $1
         ORDER BY p.created_at DESC`,
        userID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var projects []model.Project
    for rows.Next() {
        var p model.Project
        if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt); err != nil {
            return nil, err
        }
        projects = append(projects, p)
    }
    return projects, rows.Err()
}

func (r *ProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
    p := &model.Project{}
    err := r.db.QueryRow(ctx,
        `SELECT id, name, description, owner_id, created_at FROM projects WHERE id = $1`,
        id,
    ).Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt)
    if err != nil {
        return nil, err
    }
    return p, nil
}

func (r *ProjectRepository) Update(ctx context.Context, p *model.Project) error {
    _, err := r.db.Exec(ctx,
        `UPDATE projects SET name = $1, description = $2 WHERE id = $3`,
        p.Name, p.Description, p.ID,
    )
    return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
    _, err := r.db.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
    return err
}