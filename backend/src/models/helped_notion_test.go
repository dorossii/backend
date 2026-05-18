package models_test

import (
	"backend/models"
	"testing"

	"gorm.io/gorm"
)

func TestHelpedNotice(t *testing.T) {

	// テーブル初期化
	err := models.DB.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.HelpedNotice{}).Error

	if err != nil {
		t.Fatal(err)
	}

	// INSERT
	notice := models.HelpedNotice{
		NoticeID: "notice-001",
		TargetID: "user-001",
		HelperID: "user-002",
		IsRead:   false,
	}

	err = models.DB.Create(&notice).Error
	if err != nil {
		t.Fatal(err)
	}

	// SELECT
	var result models.HelpedNotice

	err = models.DB.First(
		&result,
		"notice_id = ?",
		"notice-001",
	).Error

	if err != nil {
		t.Fatal(err)
	}

	// 検証(insertしたデータと一致してるか確認)
	if result.NoticeID != "notice-001" {
		t.Fatalf(
			"unexpected NoticeID: %s",
			result.NoticeID,
		)
	}

	if result.TargetID != "user-001" {
		t.Fatalf(
			"unexpected TargetID: %s",
			result.TargetID,
		)
	}

	if result.HelperID != "user-002" {
		t.Fatalf(
			"unexpected HelperID: %s",
			result.HelperID,
		)
	}

	if result.IsRead != false {
		t.Fatalf(
			"unexpected IsRead: %v",
			result.IsRead,
		)
	}
}