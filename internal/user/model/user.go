package model

import "time"

type User struct {
	ID        uint       `gorm:"primaryKey"`
	Name      string     `gorm:"not null" validate:"required"`
	Email     string     `gorm:"unique;not null " validate:"required,email"`
	Password  string     `gorm:"not null" validate:"required,main=6"`
	CreatedAt time.Time  
	UpdatedAt time.Time  
	DeletedAt *time.Time `gorm:"index"`  // soft delete
}


