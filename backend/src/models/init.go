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

func InitForTest() {
	godotenv.Load()

	dsn := os.Getenv("DATABASE_DSN")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 外部キー制約を無効にして全テーブルを削除
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	db.Migrator().DropTable(
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
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	if err := db.AutoMigrate(
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
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	DB = db
}
