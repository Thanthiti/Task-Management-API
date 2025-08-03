package main

import (
	"os"
	"time"

	// "mymodule/config"
	"mymodule/pkg/auth"
	loger "mymodule/pkg/logger"
	"mymodule/pkg/validator"

	// User module
	userHandler "mymodule/internal/user/handler"
	userModel "mymodule/internal/user/model"
	userRepo "mymodule/internal/user/repository"
	userUsecase "mymodule/internal/user/usecase"

	// Task module
	taskHandler "mymodule/internal/task/handler"
	taskModel "mymodule/internal/task/model"
	taskRepo "mymodule/internal/task/repository"
	taskUsecase "mymodule/internal/task/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	loger.InitLogger()

	err := godotenv.Load()
	if err != nil {
		loger.Log.Fatal("Error loading .env file")
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		// AllowOrigins:     "https://your-frontend.com", //  Domain frontend
		AllowCredentials: false,
		AllowHeaders:     "Content-Type",
	}))

	// Postgres
	// db := config.InitDB()

	db, err := gorm.Open(sqlite.Open("test_task_management.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		loger.Log.Error("Failed to connect database : ", err)
	}

	if err := db.AutoMigrate(&userModel.User{},&taskModel.Task{}); err != nil {
		loger.Log.Error("Failed to connect database : ", err)
	}

	jwtKey := os.Getenv("JWT_SECRET")
	// === Initialize Core Services ===
	jwtManager := auth.NewJwtManager(jwtKey, time.Hour*2)
	cyptoService := &userUsecase.DefaultCryptoService{}
	validator := validator.InitValidator()
	
	
	// === Setup User Module ===
	userRepo := userRepo.NewGormUserRepository(db)
	useUsecase := userUsecase.NewUserUsecase(userRepo, cyptoService, jwtManager)
	userHandler.NewUserHandler(app, useUsecase, jwtManager, validator)
	
	// === Setup Task Module ===
	taskRepo := taskRepo.NewGormTaskRepository(db)
	taskUsecase := taskUsecase.NewTaskUsecase(taskRepo)
	taskHandler.NewTaskHandler(app,taskUsecase,jwtManager,validator,)
	app.Listen(":8080")

}
