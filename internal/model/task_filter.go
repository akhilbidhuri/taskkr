package model

type TaskFilter struct {
	Status   *TaskStatus
	Title    *string
	Page     int
	PageSize int
}
