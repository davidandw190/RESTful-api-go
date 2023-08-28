package db

import (
	"log"

	"github.com/davidandw190/RESTful-api-go/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DbInstance holds the database connection.
type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

// ConnectDb establishes a connection to the database.
func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database.. \n")
	}

	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running migrations...")
	db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{})

	Database = DbInstance{Db: db}

}
