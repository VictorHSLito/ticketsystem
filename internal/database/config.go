package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConnection *gorm.DB

func CreateDatabaseConnection() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Could not be able to read credential")
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	dsn := "host=localhost user=" + user + " password=" + password + " dbname=" + dbName + " port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error trying to connect to database...")
	}

	DBConnection = db
}
