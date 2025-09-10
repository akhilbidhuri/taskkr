package model

type TaskFilter struct {
	Status   TaskStatus
	Title    string
	Page     uint
	PageSize uint
}
