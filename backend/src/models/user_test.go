package models_test

import (
	"backend/models"
	"testing"
	"time"

	"gorm.io/gorm"
)


func TestUser(t *testing.T) {
	models.Init() // TODO: test毎に接続するのはキモい気がする。

	// テーブル初期化
	err := models.DB.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.User{}).Error

	if err != nil {
		t.Fatal(err)
	}

	// INSERT
	users := []models.User{
		{
			UserID:      "user-001",
			UserName:    "syatyo",
			BirthDate:   time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
			Mailadress:  "user1@example.com",
			HealthPoint: 100,
			DirtLevel:   0,
			Combo:       0,
			BgColor:     "#ffb6c1",
		},
		{
			UserID:      "user-002",
			UserName:    "goro",
			BirthDate:   time.Date(2004, 2, 2, 0, 0, 0, 0, time.UTC),
			Mailadress:  "user2@example.com",
			HealthPoint: 90,
			DirtLevel:   10,
			Combo:       1,
			BgColor:     "#add8e6",
		},
		{
			UserID:      "user-003",
			UserName:    "mattu",
			BirthDate:   time.Date(2004, 5, 5, 0, 0, 0, 0, time.UTC),
			Mailadress:  "user5@example.com",
			HealthPoint: 60,
			DirtLevel:   40,
			Combo:       4,
			BgColor:     "#dda0dd",
		},
		{
			UserID:      "user-004",
			UserName:    "saya",
			BirthDate:   time.Date(2004, 3, 3, 0, 0, 0, 0, time.UTC),
			Mailadress:  "user3@example.com",
			HealthPoint: 80,
			DirtLevel:   20,
			Combo:       2,
			BgColor:     "#90ee90",
		},
		{
			UserID:      "user-005",
			UserName:    "yoh",
			BirthDate:   time.Date(2004, 4, 4, 0, 0, 0, 0, time.UTC),
			Mailadress:  "user4@example.com",
			HealthPoint: 70,
			DirtLevel:   30,
			Combo:       3,
			BgColor:     "#ffffe0",
		},

	}

	err = models.DB.Create(&users).Error
	if err != nil {
		t.Fatal(err)
	}

	// SELECT
	var result models.User

	err = models.DB.First(
		&result,
		"user_id = ?",
		"user-001",
	).Error

	if err != nil {
		t.Fatal(err)
	}

	// 検証(insertしたデータと一致してるか確認)
	if result.UserName != "syatyo" {
		t.Fatalf(
			"unexpected UserName: %s",
			result.UserName,
		)
	}

	if result.Mailadress != "user1@example.com" {
		t.Fatalf(
			"unexpected Mailadress: %s",
			result.Mailadress,
		)
	}

	if result.HealthPoint != 100 {
		t.Fatalf(
			"unexpected HealthPoint: %d",
			result.HealthPoint,
		)
	}

	if result.Combo != 0 {
		t.Fatalf(
			"unexpected Combo: %d",
			result.Combo,
		)
	}
}