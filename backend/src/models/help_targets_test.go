package models_test

import (
	"backend/models"
	"testing"

	"gorm.io/gorm"
)

func TestHelpTargets(t *testing.T) {
	models.Init() // TODO: test毎に接続するのはキモい気がする。

	// テーブル初期化
	err := models.DB.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.HelpTargets{}).Error

	if err != nil {
		t.Fatal(err)
	}

	// INSERT
	target := models.HelpTargets{
		UserID:   "user-001",
		FriendID: "user-003",
	}

	err = models.DB.Create(&target).Error
	if err != nil {
		t.Fatal(err)
	}

	// SELECT
	var result models.HelpTargets

	err = models.DB.First(
		&result,
		"user_id = ? AND friend_id = ?",
		"user-001",
		"user-003",
	).Error

	if err != nil {
		t.Fatal(err)
	}

	// 検証(insertしたデータと一致してるか確認)
	if result.UserID != "user-001" {
		t.Fatalf(
			"unexpected UserID: %s",
			result.UserID,
		)
	}

	if result.FriendID != "user-003" {
		t.Fatalf(
			"unexpected FriendID: %s",
			result.FriendID,
		)
	}
}