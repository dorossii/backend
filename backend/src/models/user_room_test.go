package models_test

import (
	"backend/models"
	"testing"

	"gorm.io/gorm"
)

func TestUserRoom(t *testing.T) {
	models.Init() // TODO: test毎に接続するのはキモい気がする。

	// テーブル初期化
	err := models.DB.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.UserRoom{}).Error

	if err != nil {
		t.Fatal(err)
	}

	// INSERT
	rooms := []models.UserRoom{
		{
			UserID:    "user-001",
			IsAlone:   true,
			HasWasher: true,
			HasVacuum: true,
		},
		{
			UserID:    "user-002",
			IsAlone:   true,
			HasWasher: true,
		},
		{
			UserID: "user-003",
		},
	}

	err = models.DB.Create(&rooms).Error
	if err != nil {
		t.Fatal(err)
	}

	// SELECT
	var result models.UserRoom

	err = models.DB.First(
		&result,
		"user_id = ?",
		"user-003",
	).Error

	if err != nil {
		t.Fatal(err)
	}

	// 検証(insertしたデータと一致してるか確認)
	if result.UserID != "user-003" {
		t.Fatalf(
			"unexpected UserID: %s",
			result.UserID,
		)
	}

	if result.IsAlone != false {
		t.Fatalf(
			"unexpected IsAlone: %v",
			result.IsAlone,
		)
	}

	if result.HasWasher != true {
		t.Fatalf(
			"unexpected HasWasher: %v",
			result.HasWasher,
		)
	}

	if result.HasVacuum != true {
		t.Fatalf(
			"unexpected HasVacuum: %v",
			result.HasVacuum,
		)
	}
}