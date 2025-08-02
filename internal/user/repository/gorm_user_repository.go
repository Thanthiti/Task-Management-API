package repository

import (
	"errors"
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
	if err := r.db.Create(&user).Error; err != nil {
		logger.LogUser(user).Error("Failed to create user")
		return err
	}
	logger.LogUser(user).Info("User created successfully")
	return nil
}

func (r *GormUserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Log.WithField("email", email).Info("User not found by email")
		return nil, nil
	}
	if result.Error != nil {
		logger.Log.WithFields(map[string]interface{}{
			"email": email,
			"error": result.Error,
		}).Error("Database error finding user by email")
		return nil, result.Error
	}

	logger.LogUser(user).Info("User found by email")
	return &user, nil
}

func (r *GormUserRepository) FindByID(userID uint) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, userID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Log.WithField("userID", userID).Info("User not found by ID")
		return nil, result.Error
	}
	if result.Error != nil {
		logger.Log.WithFields(map[string]interface{}{
			"userID": userID,
			"error":  result.Error,
		}).Error("Database error finding user by ID")
		return nil, result.Error
	}

	logger.LogUser(user).Info("User found by ID")
	return &user, nil
}

func (r *GormUserRepository) Update(user model.User) error {
	if err := r.db.Save(&user).Error; err != nil {
		logger.LogUser(user).Error("Failed to update user")
		return err
	}
	logger.LogUser(user).Info("User updated successfully")
	return nil
}

func (r *GormUserRepository) Delete(userID uint) error {
	if err := r.db.Delete(&model.User{}, userID).Error; err != nil {
		logger.Log.WithField("userID", userID).Error("Failed to delete user")
		return err
	}
	logger.Log.WithField("userID", userID).Info("User deleted successfully")
	return nil
}
