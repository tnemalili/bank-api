package relationaldb

import (
	"os"
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDBClient() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	db, err := gorm.Open(sqlite.Open(dbHost), &gorm.Config{})
	if err != nil {
		log.Println("failed to connect database")
		return nil
	}
	return db
}