package repository

import (
	"errors"
	"fmt"
	"mymodule/internal/user/model"
	"mymodule/internal/user/usecase"
	"mymodule/pkg/logger"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) usecase.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Save(user model.User) error {
	err := r.db.Create(&user).Error
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"email": user.Email,
		}).Error("Failed to create user")
		return err
	}
	logger.Log.WithFields(map[string]interface{}{
		"email": user.Email,
	}).Info("User created successfully")
	return nil
}

func (r *GormUserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil 
	}
	
	if result.Error != nil {
		fmt.Println("7")
		logger.Log.WithFields(map[string]interface{}{
			"email": email,
			"error": result.Error,
			}).Error("Database error finding user by email")
			return nil, result.Error
		}
		
	fmt.Println("8")
	logger.Log.WithFields(map[string]interface{}{"email": email}).Info("User found by email")
	return &user, nil
}


func (r *GormUserRepository) FindByID(userID uint) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, userID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Log.WithFields(map[string]interface{}{
				"userID": userID,
			}).Info("User not found by ID")
		} else {
			logger.Log.WithFields(map[string]interface{}{
				"userID": userID,
				"error": result.Error,
			}).Error("Database error finding user by ID")
		}
		return nil, result.Error
	}

	logger.Log.WithFields(map[string]interface{}{
		"userID": userID,
	}).Info("User found by ID")
	return &user, nil
}

func (r *GormUserRepository) Update(user model.User) error {
	err := r.db.Save(&user).Error
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"userID": user.ID,
		}).Error("Failed to update user")
		return err
	}
	logger.Log.WithFields(map[string]interface{}{
		"userID": user.ID,
	}).Info("User updated successfully")
	return nil
}

func (r *GormUserRepository) Delete(userID uint) error {
	err := r.db.Delete(&model.User{}, userID).Error
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"userID": userID,
		}).Error("Failed to delete user")
		return err
	}
	logger.Log.WithFields(map[string]interface{}{
		"userID": userID,
	}).Info("User deleted successfully")
	return nil
}
