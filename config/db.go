package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
func InitDB() *gorm.DB{
	dsn := ConnectDB()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect DB: %v",err)
	}
	return db
}

func ConnectDB() string{
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	host := os.Getenv("DB_HOST")
	port :=os.Getenv("DB_PORT") 
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSL")

		if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		panic("database environment variables not set")
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)

}
