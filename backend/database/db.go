package database

import (
	"log"

	"github.com/CBYeuler/automation-backend/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("../data/automation.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connection established")

	MigrateModels()
}

func MigrateModels() {
	err := DB.AutoMigrate(&models.Machine{})
	if err != nil {
		log.Fatal("Failed to migrate database models:", err)
	}
	log.Println("Database models migrated successfully")
}

func GetDB() *gorm.DB {
	return DB
}
