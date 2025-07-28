package usecase

import (
	model "mymodule/internal/user/models"
	"mymodule/pkg/logger"
	"mymodule/pkg/auth"
)

type UserRepository interface {
	Save(user model.User) error
	FindByEmail(email string) (*model.User,error)
	FindByID(userID uint) (*model.User,error)
	Update(user model.User) error
	Delete(userID uint) error
}

type UserUsecase interface {
	CreateUser(user model.User) error
	Login(email, password string) (string , error)
	UpdateUser(user model.User) error
	DeleteUser(userID uint) error
}

type UserusecaseImpl struct{
	repo UserRepository
	cypto CryptoService
	token auth.TokenService
}

func NewUserUsecase(repo UserRepository, cypto CryptoService,token auth.TokenService) UserUsecase{
	return  &UserusecaseImpl{
		repo:repo,
		cypto: cypto,
		token: token,
	}
}

func (uc *UserusecaseImpl) CreateUser(user model.User) error{
	hashPassword ,err:= uc.cypto.HashedPassword(user.Password)
	if err != nil {
		logger.Error("Failed to hash password:",err)
		return  err
	}
	user.Password = string(hashPassword)
	if err := uc.repo.Save(user) ;err != nil {
		logger.Log.Warn()
		return  err
	}
	return  nil
}

func (uc *UserusecaseImpl) Login(email ,password string) (string,error){
	
	return  "",nil
}

func (uc *UserusecaseImpl) UpdateUser(user model.User) error{

	return  nil
}

func (uc *UserusecaseImpl) DeleteUser(userID uint) error{
	return nil
}

