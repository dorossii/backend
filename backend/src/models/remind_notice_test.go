package models_test

import (
	"backend/models"
	"testing"
	"time"
)

func TestRemindNotice(t *testing.T) {
	// INSERT
	notices := []models.RemindNotice{
		{
			NoticeID:   "notice-002",
			UserID:     "user-001",
			SenderID:   "",
			Title:      "【公式】今日のタスクを確認してください",
			NotifiedAt: time.Now(),
			IsRead:     false,
		},
		{
			NoticeID:   "notice-003",
			UserID:     "user-001",
			SenderID:   "user-002",
			Title:      "【フレンド】user2があなたを応援しています",
			NotifiedAt: time.Now(),
			IsRead:     false,
		},
	}

	err := models.DB.Create(&notices).Error
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

	if results[0].SenderID != "" {
		t.Fatalf(
			"unexpected SenderID: %s",
			results[0].SenderID,
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

	if results[1].SenderID != "user-002" {
		t.Fatalf(
			"unexpected SenderID: %s",
			results[1].SenderID,
		)
	}

	if results[1].Title != "【フレンド】user2があなたを応援しています" {
		t.Fatalf(
			"unexpected Title: %s",
			results[1].Title,
		)
	}
}
