package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"be-binareversi/model"
)

var DB *gorm.DB

func InitDatabase() {
	const dbPath = "data/rooms.db"

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatalf("failed to create db file: %v", err)
		}
		file.Close()
	}

	// DB接続
	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// マイグレーション
	err = database.AutoMigrate(&model.Room{})
	err = database.AutoMigrate(&model.Player{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	DB = database
	log.Println("Database initialized and migrated.")
}
