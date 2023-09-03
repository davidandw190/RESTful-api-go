package db

import (
	"log"

	"github.com/davidandw190/RESTful-api-go/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Initialize and hold the database connection.
type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

// ConnectDb establishes a connection to the database.
func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database.. \n")
	}

	log.Println("Running db migrations...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
	); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	Database = DbInstance{Db: db}
}
