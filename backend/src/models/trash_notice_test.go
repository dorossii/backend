package models_test

import (
	"backend/models"
	"testing"

	"gorm.io/gorm"
)

func TestTrashNotice(t *testing.T) {

	// テーブル初期化
	err := models.DB.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.TrashNotice{}).Error

	if err != nil {
		t.Fatal(err)
	}

	// INSERT
	notice := models.TrashNotice{
		NoticeID:   "notice-005",
		SenderID:   "user-001",
		ReceiverID: "user-003",
		Count:      1,
	}

	err = models.DB.Create(&notice).Error
	if err != nil {
		t.Fatal(err)
	}

	// SELECT
	var result models.TrashNotice

	err = models.DB.First(
		&result,
		"notice_id = ?",
		"notice-005",
	).Error

	if err != nil {
		t.Fatal(err)
	}

	// 検証(insertしたデータと一致してるか確認)
	if result.NoticeID != "notice-005" {
		t.Fatalf(
			"unexpected NoticeID: %s",
			result.NoticeID,
		)
	}

	if result.SenderID != "user-001" {
		t.Fatalf(
			"unexpected SenderID: %s",
			result.SenderID,
		)
	}

	if result.ReceiverID != "user-003" {
		t.Fatalf(
			"unexpected ReceiverID: %s",
			result.ReceiverID,
		)
	}

	if result.Count != 1 {
		t.Fatalf(
			"unexpected Count: %d",
			result.Count,
		)
	}
}