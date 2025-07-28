package main

import (
	"fmt"
	// "mymodule/config"
	model "mymodule/internal/user/models"
	"mymodule/internal/user/repository"
	"mymodule/internal/user/usecase"
	"mymodule/pkg/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Hello world")
	logger.InitLogger()
	// Postgres
	// db := config.InitDB()

	db,err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect database : %v",err)
	}
	
	if err := db.AutoMigrate(&model.User{});err != nil {
		logger.Error("Failed to connect database : %v",err)
	}
	
	userRepo := repository.NewGormUserRepository(db)
	cyptoService := &usecase.DefaultCryptoService{}

	userHandler := usecase.NewUserUsecase(userRepo,cyptoService)
	
	TestData := model.User{
		Name: "Golang",
		Email: "a@gmail.com",
		Password: "Test1234",
	}
	
	result  := userHandler.CreateUser(TestData)
	logger.Debug(result)	
	

}
