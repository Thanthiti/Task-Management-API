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
func (r * GormUserRepository) FindByEmail(email string ) (*model.User,error){
	return &model.User{},nil
}
func (r * GormUserRepository) FindByID(userID uint ) (*model.User,error){
	return &model.User{},nil
}
func (r * GormUserRepository) Update(user model.User ) error{
	return nil
}
func (r * GormUserRepository) Delete(userID uint ) error{
	return nil
}