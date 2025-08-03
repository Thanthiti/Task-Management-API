package model

import (
	"time"

	"gorm.io/gorm"
)

// Task is the DB model for a task
type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id" example:"1"`
	Title       string     `gorm:"type:text;not null" json:"title" example:"Write blog post" validate:"required"`
	Description string     `gorm:"type:text" json:"description" example:"Write about Clean Architecture"`
	DueDate     *time.Time `gorm:"default:null" json:"due_date,omitempty" example:"2025-08-10T15:00:00Z"`
	Status      string     `gorm:"type:varchar(20);default:'pending'" json:"status" example:"pending" validate:"oneof=pending in_progress completed"`
	UserID      uint       `gorm:"not null" json:"user_id" example:"1"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateTaskRequest is the request model for creating a task
type CreateTaskRequest struct {
	Title       string     `json:"title" example:"Write blog post" validate:"required"`
	Description string     `json:"description,omitempty" example:"Write about Clean Architecture"`
	DueDate     *time.Time `json:"due_date,omitempty" example:"2025-08-10T15:00:00Z"`
}

// UpdateTaskRequest is the request model for updating a task
type UpdateTaskInput struct {
    Title       *string     `json:"title,omitempty"`
    Description *string     `json:"description,omitempty"`
    DueDate     *time.Time  `json:"due_date,omitempty"`
    Status      *string     `json:"status,omitempty" validate:"omitempty,oneof=pending in_progress completed"`
}

// TaskResponse is the response model for a task
type TaskResponse struct {
	ID          uint       `json:"id" example:"1"`
	Title       string     `json:"title" example:"Write blog post"`
	Description string     `json:"description" example:"Write about Clean Architecture"`
	DueDate     *time.Time `json:"due_date,omitempty" example:"2025-08-10T15:00:00Z"`
	Status      string     `json:"status" example:"pending"`
	UserID      uint       `json:"user_id" example:"1"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

