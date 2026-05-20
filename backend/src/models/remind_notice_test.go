package models_test

import (
	"backend/models"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestRemindNotice(t *testing.T) {

	// テーブル初期化
	err := models.DB.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.RemindNotice{}).Error

	if err != nil {
		t.Fatal(err)
	}

	// INSERT
	notices := []models.RemindNotice{
		{
			NoticeID:   "notice-002",
			UserID:     "user-001",
			SenderType: models.SenderTypeOfficial,
			Title:      "【公式】今日のタスクを確認してください",
			NotifiedAt: time.Now(),
			IsRead:     false,
		},
		{
			NoticeID:   "notice-003",
			UserID:     "user-001",
			SenderType: models.SenderTypeFriend,
			Title:      "【フレンド】user2があなたを応援しています",
			NotifiedAt: time.Now(),
			IsRead:     false,
		},
	}

	err = models.DB.Create(&notices).Error
	if err != nil {
		t.Fatal(err)
	}

	// SELECT
	var results []models.RemindNotice

	err = models.DB.
		Where("user_id = ?", "user-001").
		Order("notice_id asc").
		Find(&results).Error

	if err != nil {
		t.Fatal(err)
	}

	// 件数確認
	if len(results) != 2 {
		t.Fatalf(
			"unexpected notice count: %d",
			len(results),
		)
	}

	// 検証(insertしたデータと一致してるか確認)
	if results[0].NoticeID != "notice-002" {
		t.Fatalf(
			"unexpected NoticeID: %s",
			results[0].NoticeID,
		)
	}

	if results[0].SenderType != models.SenderTypeOfficial {
		t.Fatalf(
			"unexpected SenderType: %s",
			results[0].SenderType,
		)
	}

	if results[0].Title != "【公式】今日のタスクを確認してください" {
		t.Fatalf(
			"unexpected Title: %s",
			results[0].Title,
		)
	}

	if results[1].NoticeID != "notice-003" {
		t.Fatalf(
			"unexpected NoticeID: %s",
			results[1].NoticeID,
		)
	}

	if results[1].SenderType != models.SenderTypeFriend {
		t.Fatalf(
			"unexpected SenderType: %s",
			results[1].SenderType,
		)
	}

}