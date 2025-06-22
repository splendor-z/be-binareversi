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

	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		log.Fatalf("failed to create data directory: %v", err)
	}

	// DBファイルが存在しない場合は作成
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatalf("failed to create db file: %v", err)
		}
		file.Close()
		log.Println("Database file created.")
	}

	// DB接続
	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// マイグレーション
	if err := database.AutoMigrate(
		&model.Room{},
		&model.Player{},
	); err != nil {
		log.Fatalf("failed to migrate models: %v", err)
	}

	DB = database
	log.Println("Database initialized and migrated.")
}
