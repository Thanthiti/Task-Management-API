package usecase

import (
	"fmt"
	"mymodule/internal/user/model"
	"mymodule/pkg/auth"
	"mymodule/pkg/logger"
)

type UserRepository interface {
	Save(user model.User) error
	FindByEmail(email string) (*model.User,error)
	FindByID(userID uint) (*model.User,error)
	Update(user model.User) error
	Delete(userID uint) error
}

type UserUsecase interface {
	Register(user model.User) error
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

func (uc *UserusecaseImpl) Register(user model.User) error{
	// Check email already exits 
	exitUser ,err := uc.repo.FindByEmail(user.Email)
	if err != nil && exitUser != nil{
		logger.Log.Warn("Email already exits : ", user.Email)
		return  fmt.Errorf("email already registed")
	}
	
	// Hash password  
	hashPassword ,err:= uc.cypto.HashedPassword(user.Password)
	if err != nil {
		logger.Log.Error("Hash Failed :",err)
		return  err
	}
	user.Password = string(hashPassword)
	
	// Save to DB 
	if err := uc.repo.Save(user) ;err != nil {
		logger.Log.Error("Enable save user : ",err)
		return  fmt.Errorf("enable save user : %v",err)
	}

	logger.Log.Info("User created : ", user.Email)
	return  nil
}

func (uc *UserusecaseImpl) Login(email ,password string) (string,error){
	
	user , err := uc.repo.FindByEmail(email)
	if err != nil || user == nil {
		logger.Log.Warn("Login failed : email not found")
		return "",fmt.Errorf("invalid credentials")
	}

	// Decode password in database and compare 
	if !uc.cypto.ComparePassword(user.Password,password){
		logger.Log.Warn("Login failed : incorect password")
		return "",fmt.Errorf("invalid credentials")
	}
	
	// Generate Token
	token,err := uc.token.GenerateToken(user.ID)
	if err != nil {
		logger.Log.Error("Token generation failed : ",err)
		return "",fmt.Errorf("invalid credentials")
		
	}
	
	logger.Log.Info("User logged in : ",email)
	return  token,nil
}

func (uc *UserusecaseImpl) UpdateUser(user model.User) error{

	exitUser,err := uc.repo.FindByID(user.ID)
	if err != nil || exitUser == nil {
		logger.Log.Warn("Update failed: user not found")
		return fmt.Errorf("user not found")
	}

	exitUser.Name = user.Name
	exitUser.Email = user.Email

	if err := uc.repo.Update(*exitUser);err != nil {
		logger.Log.Error("Update failed : ", err)
		return err
	}

	logger.Log.Info("User update : ",user.ID)
	return  nil
}

func (uc *UserusecaseImpl) DeleteUser(userID uint) error{
	user, err := uc.repo.FindByID(userID)
	if err != nil || user == nil {
		logger.Log.Warn("Delete failed: user not found")
		return fmt.Errorf("user not found")
	}
	
	if err := uc.repo.Delete(userID);err != nil {
		logger.Log.Error("Delete failed : ",err)
		return err
	}
	
	logger.Log.Info("User deleted : ",userID)
	return nil
}

