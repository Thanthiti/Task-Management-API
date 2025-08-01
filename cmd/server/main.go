package main

import (
	"os"
	"time"

	// "mymodule/config"
	"mymodule/internal/user/handler"
	"mymodule/internal/user/model"
	"mymodule/internal/user/repository"
	"mymodule/internal/user/usecase"
	"mymodule/pkg/auth"
	loger "mymodule/pkg/logger"
	"mymodule/pkg/validator"

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

	db,err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		loger.Log.Error("Failed to connect database : ",err)
	}
	
	if err := db.AutoMigrate(&model.User{});err != nil {
		loger.Log.Error("Failed to connect database : ",err)
	}

	
	jwtKey := os.Getenv("JWT_SECRET") 

	jwtManager := auth.NewJwtManager(jwtKey,time.Hour*2)
	cyptoService := &usecase.DefaultCryptoService{}
	validator := validator.InitValidator()
	
	userRepo := repository.NewGormUserRepository(db)
	useUsecase := usecase.NewUserUsecase(userRepo,cyptoService,jwtManager)
	handler.NewUserHandler(app,useUsecase,jwtManager,validator)
	
	// TestData := model.User{
	// 	Name: "Golang",
	// 	Email: "b@gmil.com",
	// 	Password: "Test1234",
	// }
	
	// result  := userHandler.CreateUser(TestData)
	// logger.Log.Debug(result)	

	app.Listen(":8080")
	

}
