package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/Ruturaj-7802/taskflow/internal/dto"
	"github.com/Ruturaj-7802/taskflow/internal/model"
	"github.com/Ruturaj-7802/taskflow/internal/repository"
)

type TaskService struct {
	taskRepo    *repository.TaskRepository
	projectRepo *repository.ProjectRepository
}

func NewTaskService(taskRepo *repository.TaskRepository, projectRepo *repository.ProjectRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo, projectRepo: projectRepo}
}

func (s *TaskService) Create(ctx context.Context, projectID string, req dto.CreateTaskRequest, ownerID uuid.UUID) (model.Task, error) {
	pid, err := uuid.Parse(projectID)
	if err != nil {
		return model.Task{}, errors.New("invalid project id")
	}

	p, err := s.projectRepo.GetByID(ctx, pid)
	if err != nil {
		return model.Task{}, errors.New("project not found")
	}
	if p.OwnerID != ownerID {
		return model.Task{}, errors.New("forbidden")
	}

	// description := ""
	// if req.Description != nil {
	// 	description = *req.Description
	// }

	status := "todo"
	if req.Status != nil {
		status = *req.Status
	}

	priority := "medium"
	if req.Priority != nil {
		priority = *req.Priority
	}

	task := &model.Task{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    priority,
		ProjectID:   pid,
		AssigneeID:  req.AssigneeID,
		DueDate:     req.DueDate,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return model.Task{}, err
	}

	return *task, nil
}

func (s *TaskService) List(ctx context.Context, projectID string, ownerID uuid.UUID) ([]model.Task, error) {
	pid, err := uuid.Parse(projectID)
	if err != nil {
		return nil, errors.New("invalid project id")
	}

	p, err := s.projectRepo.GetByID(ctx, pid)
	if err != nil {
		return nil, errors.New("project not found")
	}
	if p.OwnerID != ownerID {
		return nil, errors.New("forbidden")
	}

	return s.taskRepo.ListByProject(ctx, pid, repository.TaskFilter{})
}

func (s *TaskService) Get(ctx context.Context, id string, ownerID uuid.UUID) (model.Task, error) {
	tid, err := uuid.Parse(id)
	if err != nil {
		return model.Task{}, errors.New("invalid id")
	}

	t, err := s.taskRepo.FindByID(ctx, tid)
	if err != nil {
		return model.Task{}, errors.New("not found")
	}

	p, err := s.projectRepo.GetByID(ctx, t.ProjectID)
	if err != nil || p.OwnerID != ownerID {
		return model.Task{}, errors.New("forbidden")
	}

	return *t, nil
}

func (s *TaskService) Update(ctx context.Context, id string, req dto.UpdateTaskRequest, ownerID uuid.UUID) (model.Task, error) {
	tid, err := uuid.Parse(id)
	if err != nil {
		return model.Task{}, errors.New("invalid id")
	}

	t, err := s.taskRepo.FindByID(ctx, tid)
	if err != nil {
		return model.Task{}, errors.New("not found")
	}

	p, err := s.projectRepo.GetByID(ctx, t.ProjectID)
	if err != nil || p.OwnerID != ownerID {
		return model.Task{}, errors.New("forbidden")
	}

	if req.Title != nil {
		t.Title = *req.Title
	}
	if req.Description != nil {
		t.Description = req.Description
	}
	if req.Status != nil {
		t.Status = *req.Status
	}
	if req.Priority != nil {
		t.Priority = *req.Priority
	}
	if req.AssigneeID != nil {
		t.AssigneeID = req.AssigneeID
	}
	if req.DueDate != nil {
		t.DueDate = req.DueDate
	}

	if err := s.taskRepo.Update(ctx, t); err != nil {
		return model.Task{}, err
	}

	return *t, nil
}

func (s *TaskService) Delete(ctx context.Context, id string, ownerID uuid.UUID) error {
	tid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid id")
	}

	t, err := s.taskRepo.FindByID(ctx, tid)
	if err != nil {
		return errors.New("not found")
	}

	p, err := s.projectRepo.GetByID(ctx, t.ProjectID)
	if err != nil || p.OwnerID != ownerID {
		return errors.New("forbidden")
	}

	return s.taskRepo.Delete(ctx, tid)
}