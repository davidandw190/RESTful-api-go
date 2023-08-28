package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database.. \n")
		db.Logger = logger.Default.LogMode(logger.Info)
		log.Println("Running migrations...")
		//TODO: Add migrations

		Database = DbInstance{Db: db}

	}
}
