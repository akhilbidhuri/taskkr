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
	// You can add validation logic here if needed
	if task.Title == "" {
		return errors.New("title cannot be empty")
	}
	return s.repo.Create(ctx, task)
}

func (s *TaskService) GetByID(ctx context.Context, id string) (*model.Task, error) {
	return s.repo.GetByID(ctx, id)
}
