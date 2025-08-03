package model

import (
	"time"

	"gorm.io/gorm"
)

// Model for DB(GORM)
type User struct {
	ID        uint       `gorm:"primaryKey"`
	Name      string     `gorm:"not null" validate:"required"`
	Email     string     `gorm:"unique;not null " validate:"required,email"`
	Password  string     `gorm:"not null" validate:"required,main=6"`
	CreatedAt time.Time  
	UpdatedAt time.Time  
	DeletedAt gorm.DeletedAt `gorm:"index"`  
}

//  Request for register model 
type RegisterRequest struct {
	Name     string `json:"name" example:"John Doe" validate:"required"`
	Email    string `json:"email" example:"john@example.com" validate:"required,email"`
	Password string `json:"password" example:"12345678" validate:"required,min=6"`
}

//  Request for login model
type LoginRequest struct {
	Email    string `json:"email" example:"john@example.com" validate:"required,email"`
	Password string `json:"password" example:"12345678" validate:"required,min=6"`
}

// Response model
type UserResponse struct {
	ID    uint   `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

// Update model 
type UpdateUserRequest struct {
	ID    uint   `json:"id" example:"1"`
	Name  string `json:"name,omitempty"`  
	Email string `json:"email,omitempty"`
}



