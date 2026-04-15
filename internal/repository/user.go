package repository

import (
    "context"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/ruturaj/taskflow/internal/model"
)

type UserRepository struct {
    db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
    _, err := r.db.Exec(ctx,
        `INSERT INTO users (id, name, email, password, created_at)
         VALUES ($1, $2, $3, $4, $5)`,
        user.ID, user.Name, user.Email, user.Password, user.CreatedAt,
    )
    return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    user := &model.User{}
    err := r.db.QueryRow(ctx,
        `SELECT id, name, email, password, created_at FROM users WHERE email = $1`,
        email,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
    user := &model.User{}
    err := r.db.QueryRow(ctx,
        `SELECT id, name, email, created_at FROM users WHERE id = $1`,
        id,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return user, nil
}