package usecase
import (
	model "mymodule/internal/user/models"
)

type UserRepository interface {
	Save(user model.User) error
}

type UserUsecase interface {
	CreateUser(user model.User) error
}

type UserusecaseImpl struct{
	repo UserRepository
}

func NewUserUsecase(repo UserRepository) UserUsecase{
	return  &UserusecaseImpl{repo:repo}
}

func (uc *UserusecaseImpl) CreateUser(user model.User) error{
	// pass
	return  nil
}

