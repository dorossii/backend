package models

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB = nil
)

func Init() {
	// データベースを開く
	godotenv.Load()

	// データベースの接続情報
	dsn := os.Getenv("DATABASE_DSN")

	// データベースを開く
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(
		&BaseTask{},
		&FriendShips{},
		&HelpTargets{},
		&HelpedNotice{},
		&RemindNotice{},
		&RescueNotice{},
		&Task{},
		&TrashNotice{},
		&UserRoom{},
		&User{},
	)
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// グローバル変数に格納
	DB = db
}
