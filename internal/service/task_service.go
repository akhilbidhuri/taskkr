package service

import (
	"context"
	"errors"

	"github.com/akhilbidhuri/taskkr/internal/repository"

	"github.com/akhilbidhuri/taskkr/internal/model"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) Create(ctx context.Context, task *model.Task) error {
	if task.Title == "" {
		return errors.New("title cannot be empty")
	}
	return s.repo.Create(ctx, task)
}

func (s *TaskService) GetByID(ctx context.Context, id string) (*model.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) List(ctx context.Context, filter *model.TaskFilter) ([]*model.Task, int, error) {
	return s.repo.List(ctx, filter)
}

func (s *TaskService) Update(ctx context.Context, id string, task *model.UpdateTask) (*model.Task, error) {
	return s.repo.Update(ctx, id, task)
}

func (s *TaskService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
