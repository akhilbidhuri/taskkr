package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/akhilbidhuri/taskkr/internal/config"
	"github.com/akhilbidhuri/taskkr/internal/model"
	"github.com/akhilbidhuri/taskkr/internal/repository"
	"github.com/akhilbidhuri/taskkr/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Run AutoMigrate
	if err := db.AutoMigrate(&model.Task{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) repository.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *taskRepository) GetByID(ctx context.Context, id string) (*model.Task, error) {
	var task model.Task
	err := r.db.WithContext(ctx).First(&task, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) List(ctx context.Context, filter *model.TaskFilter) ([]*model.Task, int, error) {
	var tasks []*model.Task
	query := r.db.WithContext(ctx).Model(&model.Task{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Title != "" {
		query = query.Where("title ILIKE ?", "%"+filter.Title+"%")
	}

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	offset := (filter.Page - 1) * filter.PageSize
	err = query.Offset(int(offset)).Limit(int(filter.PageSize)).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, int(total), nil
}

func (r *taskRepository) Update(ctx context.Context, id string, task *model.UpdateTask) (*model.Task, error) {
	result := r.db.WithContext(ctx).
		Model(&model.Task{}).
		Where("id = ?", id).
		Updates(task)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, utils.NoEntryError
	}

	var updatedTask model.Task
	if err := r.db.WithContext(ctx).First(&updatedTask, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &updatedTask, nil
}

func (r *taskRepository) Delete(ctx context.Context, id string) error {
	rowsAffected := r.db.WithContext(ctx).Delete(&model.Task{}, "id = ?", id).RowsAffected
	if rowsAffected == 0 {
		return utils.NoEntryError
	}
	return nil
}
