package service

import (
	"context"
	"errors"
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

func (s *TaskService) Create(ctx context.Context, projectID string, req dto.CreateTaskRequest, ownerID string) (model.Task, error) {
	p, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return model.Task{}, errors.New("project not found")
	}
	if p.OwnerID != ownerID {
		return model.Task{}, errors.New("forbidden")
	}
	status := req.Status
	if status == "" {
		status = "todo"
	}
	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}
	t := model.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    priority,
		ProjectID:   projectID,
		AssigneeID:  req.AssigneeID,
	}
	return s.taskRepo.Create(ctx, t)
}

func (s *TaskService) List(ctx context.Context, projectID, ownerID string) ([]model.Task, error) {
	p, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}
	if p.OwnerID != ownerID {
		return nil, errors.New("forbidden")
	}
	return s.taskRepo.ListByProject(ctx, projectID)
}

func (s *TaskService) Get(ctx context.Context, id, ownerID string) (model.Task, error) {
	t, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return model.Task{}, errors.New("not found")
	}
	p, err := s.projectRepo.GetByID(ctx, t.ProjectID)
	if err != nil || p.OwnerID != ownerID {
		return model.Task{}, errors.New("forbidden")
	}
	return t, nil
}

func (s *TaskService) Update(ctx context.Context, id string, req dto.UpdateTaskRequest, ownerID string) (model.Task, error) {
	t, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return model.Task{}, errors.New("not found")
	}
	p, err := s.projectRepo.GetByID(ctx, t.ProjectID)
	if err != nil || p.OwnerID != ownerID {
		return model.Task{}, errors.New("forbidden")
	}
	t.Title = req.Title
	t.Description = req.Description
	t.Status = req.Status
	t.Priority = req.Priority
	t.AssigneeID = req.AssigneeID
	return s.taskRepo.Update(ctx, t)
}

func (s *TaskService) Delete(ctx context.Context, id, ownerID string) error {
	t, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("not found")
	}
	p, err := s.projectRepo.GetByID(ctx, t.ProjectID)
	if err != nil || p.OwnerID != ownerID {
		return errors.New("forbidden")
	}
	return s.taskRepo.Delete(ctx, id)
}