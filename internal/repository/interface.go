package repository

import (
	"context"

	"github.com/akhilbidhuri/taskkr/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id string) (*model.Task, error)
	Update(ctx context.Context, id string, task *model.UpdateTask) (*model.Task, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *model.TaskFilter) ([]*model.Task, int, error)
}
