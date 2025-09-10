package model

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusInProcess TaskStatus = "in_process"
	StatusCompleted TaskStatus = "completed"
)

type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null" json:"user_id"` // Associate task with a user
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Status      TaskStatus     `gorm:"type:varchar(20);default:'pending'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
