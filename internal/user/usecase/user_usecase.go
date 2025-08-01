package usecase

import (
	"errors"
	"fmt"
	"mymodule/internal/user/model"
	"mymodule/pkg/auth"
	"mymodule/pkg/logger"

	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(userID uint) (*model.User, error)
	Update(user model.User) error
	Delete(userID uint) error
}

type UserUsecase interface {
	Register(user model.User) error
	Login(email, password string) (string, error)
	UpdateUser(user model.User) error
	DeleteUser(userID uint) error
}

type UserusecaseImpl struct {
	repo  UserRepository
	cypto CryptoService
	token auth.TokenService
}

func NewUserUsecase(repo UserRepository, cypto CryptoService, token auth.TokenService) UserUsecase {
	return &UserusecaseImpl{
		repo:  repo,
		cypto: cypto,
		token: token,
	}
}

func (uc *UserusecaseImpl) Register(user model.User) error {
	// Check email already exists
	exitUser, err := uc.repo.FindByEmail(user.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Error("DB error while checking email: ", err)
			return err
		}
	} else if exitUser != nil {
		logger.Log.Warn("Email already exists: ", user.Email)
		return fmt.Errorf("email already registered")
	}

	// Hash password
	hashPassword, err := uc.cypto.HashedPassword(user.Password)
	if err != nil {
		logger.Log.Error("Hash Failed :", err)
		return err
	}
	user.Password = string(hashPassword)

	// Save to DB
	if err := uc.repo.Save(user); err != nil {
		logger.Log.Error("Enable save user : ", err)
		return fmt.Errorf("enable save user : %v", err)
	}

	logger.Log.Info("User created : ", user.Email)
	return nil
}

func (uc *UserusecaseImpl) Login(email, password string) (string, error) {
	user, err := uc.repo.FindByEmail(email)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warn("Login failed : email not found")
			return "", fmt.Errorf("invalid credentials")
		}

		logger.Log.Error("DB error during login: ", err)
		return "", fmt.Errorf("internal server error")
	}

	if user == nil {
		logger.Log.Warn("Login failed : user is nil")
		return "", fmt.Errorf("invalid credentials")
	}

	// Decode password in database and compare
	if !uc.cypto.ComparePassword(user.Password, password) {
		logger.Log.Warn("Login failed : incorrect password")
		return "", fmt.Errorf("invalid credentials")
	}

	// Generate Token
	token, err := uc.token.GenerateToken(user.ID)

	if err != nil {
		logger.Log.Error("Token generation failed : ", err)
		return "", fmt.Errorf("internal server error")
	}

	logger.Log.Info("User logged in : ", email)
	return token, nil
}

func (uc *UserusecaseImpl) UpdateUser(user model.User) error {
	exitUser, err := uc.repo.FindByID(user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("1")
			logger.Log.Warn("Update failed: user not found")
			return fmt.Errorf("user not found")
		}
		fmt.Println("2")
		logger.Log.Error("DB error when finding user by ID: ", err)
		return err
	}
	if exitUser == nil {
		fmt.Println("3")
		logger.Log.Warn("Update failed: user is nil")
		return fmt.Errorf("user not found")
	}
	userWithEmail, err := uc.repo.FindByEmail(user.Email)
	if err != nil {
		fmt.Println("4")
		logger.Log.Error("DB error when checking email uniqueness: ", err)
		return err
	}
	if userWithEmail != nil && userWithEmail.ID != user.ID {
		fmt.Println("5")
		logger.Log.Warn("Update failed: email already registered", user.Email)
		return fmt.Errorf("email already registered")
	}

	exitUser.Name = user.Name
	exitUser.Email = user.Email

	if err := uc.repo.Update(*exitUser); err != nil {
		logger.Log.Error("Update failed : ", err)
		return err
	}

	logger.Log.Info("User update : ", user.ID)
	return nil
}

func (uc *UserusecaseImpl) DeleteUser(userID uint) error {
	user, err := uc.repo.FindByID(userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warn("Delete failed: user not found")
			return fmt.Errorf("user not found")
		}
		logger.Log.Error("DB error when finding user by ID: ", err)
		return err
	}

	if user == nil {
		logger.Log.Warn("Delete failed: user is nil")
		return fmt.Errorf("user not found")
	}

	if err := uc.repo.Delete(userID); err != nil {
		logger.Log.Error("Delete failed : ", err)
		return err
	}

	logger.Log.Info("User deleted : ", userID)
	return nil
}
