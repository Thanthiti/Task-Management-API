package repository

import (
	"mymodule/internal/user/usecase"
	model "mymodule/internal/user/models"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) usecase.UserRepository{
	return &GormUserRepository{db: db}
}

func (r * GormUserRepository) Save(user model.User ) error{
	return r.db.Create(&user).Error
}