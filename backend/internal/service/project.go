package service

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"

    "github.com/Ruturaj-7802/taskflow/internal/dto"
    "github.com/Ruturaj-7802/taskflow/internal/model"
    "github.com/Ruturaj-7802/taskflow/internal/repository"
)

var (
    ErrNotFound  = errors.New("not found")
    ErrForbidden = errors.New("forbidden")
)

type ProjectService struct {
    projectRepo *repository.ProjectRepository
    taskRepo    *repository.TaskRepository
}

func NewProjectService(pr *repository.ProjectRepository, tr *repository.TaskRepository) *ProjectService {
    return &ProjectService{projectRepo: pr, taskRepo: tr}
}

func (s *ProjectService) List(ctx context.Context, userID uuid.UUID) ([]dto.ProjectResponse, error) {
    projects, err := s.projectRepo.ListForUser(ctx, userID)
    if err != nil {
        return nil, err
    }
    resp := make([]dto.ProjectResponse, len(projects))
    for i, p := range projects {
        resp[i] = toProjectResponse(p)
    }
    return resp, nil
}

func (s *ProjectService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
    p := &model.Project{
        ID:          uuid.New(),
        Name:        req.Name,
        Description: req.Description,
        OwnerID:     userID,
        CreatedAt:   time.Now(),
    }
    if err := s.projectRepo.Create(ctx, p); err != nil {
        return nil, err
    }
    r := toProjectResponse(*p)
    return &r, nil
}

func (s *ProjectService) Get(ctx context.Context, id, userID uuid.UUID) (*dto.ProjectDetailResponse, error) {
    p, err := s.projectRepo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, err
    }

    tasks, err := s.taskRepo.ListByProject(ctx, id, repository.TaskFilter{})
    if err != nil {
        return nil, err
    }

    taskDTOs := make([]dto.TaskResponse, len(tasks))
    for i, t := range tasks {
        taskDTOs[i] = toTaskResponse(t)
    }

    return &dto.ProjectDetailResponse{
        ProjectResponse: toProjectResponse(*p),
        Tasks:           taskDTOs,
    }, nil
}

func (s *ProjectService) Update(ctx context.Context, id, userID uuid.UUID, req dto.UpdateProjectRequest) (*dto.ProjectResponse, error) {
    p, err := s.projectRepo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) { return nil, ErrNotFound }
        return nil, err
    }
    if p.OwnerID != userID {
        return nil, ErrForbidden  // ← 403, not 404: ownership violation
    }

    if req.Name != nil { p.Name = *req.Name }
    if req.Description != nil { p.Description = req.Description }

    if err := s.projectRepo.Update(ctx, p); err != nil {
        return nil, err
    }
    r := toProjectResponse(*p)
    return &r, nil
}

func (s *ProjectService) Delete(ctx context.Context, id, userID uuid.UUID) error {
    p, err := s.projectRepo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) { return ErrNotFound }
        return err
    }
    if p.OwnerID != userID {
        return ErrForbidden
    }

    // Delete tasks first (no ON DELETE CASCADE in migrations — explicit is better)
    if err := s.taskRepo.DeleteByProject(ctx, id); err != nil {
        return err
    }
    return s.projectRepo.Delete(ctx, id)
}

func (s *ProjectService) Stats(ctx context.Context, id, userID uuid.UUID) (*dto.ProjectStatsResponse, error) {
    _, err := s.projectRepo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) { return nil, ErrNotFound }
        return nil, err
    }

    byStatus, byAssignee, err := s.taskRepo.StatsByProject(ctx, id)
    if err != nil {
        return nil, err
    }
    return &dto.ProjectStatsResponse{ByStatus: byStatus, ByAssignee: byAssignee}, nil
}

func toProjectResponse(p model.Project) dto.ProjectResponse {
    return dto.ProjectResponse{
        ID: p.ID, Name: p.Name, Description: p.Description,
        OwnerID: p.OwnerID, CreatedAt: p.CreatedAt,
    }
}

func toTaskResponse(t model.Task) dto.TaskResponse {
    return dto.TaskResponse{
        ID:          t.ID,
        Title:       t.Title,
        Description: t.Description,
        Status:      t.Status,
        Priority:    t.Priority,
        ProjectID:   t.ProjectID,
        AssigneeID:  t.AssigneeID,
        DueDate:     t.DueDate,
        CreatedAt:   t.CreatedAt,
        UpdatedAt:   t.UpdatedAt,
    }
}