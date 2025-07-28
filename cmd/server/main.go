package main

import (
	"fmt"
	"os"
	"time"

	// "mymodule/config"
	model "mymodule/internal/user/models"
	"mymodule/internal/user/repository"
	"mymodule/internal/user/usecase"
	"mymodule/pkg/auth"
	"mymodule/pkg/logger"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Hello world")
	
	err := godotenv.Load()
	if err != nil {
		logger.Log.Fatal("Error loading .env file")
	}

	logger.InitLogger()
	
	// Postgres
	// db := config.InitDB()

	db,err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect database : ",err)
	}
	
	if err := db.AutoMigrate(&model.User{});err != nil {
		logger.Error("Failed to connect database : ",err)
	}

	
	jwtKey := os.Getenv("JWT_SECRET") 

	jwtManager := auth.NewJwtManager(jwtKey,time.Hour*2)
	cyptoService := &usecase.DefaultCryptoService{}
	
	userRepo := repository.NewGormUserRepository(db)
	userHandler := usecase.NewUserUsecase(userRepo,cyptoService,jwtManager)
	
	TestData := model.User{
		Name: "Golang",
		Email: "a@gmail.com",
		Password: "Test1234",
	}
	
	result  := userHandler.CreateUser(TestData)
	logger.Debug(result)	
	

}
