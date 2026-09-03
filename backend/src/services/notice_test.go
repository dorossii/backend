package services_test

import (
	"backend/models"
	"backend/services"
	"testing"
	"time"
)

func truncateNotices(t *testing.T) {
	t.Helper()
	if err := models.DB.Exec("TRUNCATE TABLE trash_notices").Error; err != nil {
		t.Fatal(err)
	}
	if err := models.DB.Exec("TRUNCATE TABLE rescue_notices").Error; err != nil {
		t.Fatal(err)
	}
	if err := models.DB.Exec("TRUNCATE TABLE helped_notices").Error; err != nil {
		t.Fatal(err)
	}
}

// TestGetNotices はTrash/Rescue/Helpedの3種の通知が時刻降順でまとめて返ることを確認する
func TestGetNotices(t *testing.T) {
	truncateUsers(t)
	truncateNotices(t)

	createUser(t, "user-001", "太郎", "pineTree", "icon1")
	createUser(t, "user-002", "花子", "pineTree", "icon1")

	now := time.Now()

	if err := models.DB.Create(&models.TrashNotice{
		NoticeID:   "notice-trash-1",
		SenderID:   "user-002",
		ReceiverID: "user-001",
		Count:      3,
		CreatedAt:  now.Add(-2 * time.Minute),
	}).Error; err != nil {
		t.Fatalf("failed to create TrashNotice: %v", err)
	}

	if err := models.DB.Create(&models.RescueNotice{
		NoticeID:  "notice-rescue-1",
		TargetID:  "user-002",
		HelperID:  "user-001",
		CreatedAt: now.Add(-1 * time.Minute),
	}).Error; err != nil {
		t.Fatalf("failed to create RescueNotice: %v", err)
	}

	if err := models.DB.Create(&models.HelpedNotice{
		NoticeID:  "notice-helped-1",
		TargetID:  "user-001",
		HelperID:  "user-002",
		CreatedAt: now,
	}).Error; err != nil {
		t.Fatalf("failed to create HelpedNotice: %v", err)
	}

	notices, err := services.GetNotices("user-001")
	if err != nil {
		t.Fatalf("GetNotices failed: %v", err)
	}

	// user-001 は Trash の受信者、Rescue の Helper(HelperID)、Helped の対象(TargetID)を
	// それぞれ兼ねるため、3件すべてが該当する
	if len(notices) != 3 {
		t.Fatalf("expected 3 notices for user-001, got %d: %+v", len(notices), notices)
	}

	// 新しい順（HelpedNotice -> RescueNotice -> TrashNotice）で返る
	if notices[0].SenderType != "user-002" {
		t.Fatalf("unexpected first notice SenderType: %s", notices[0].SenderType)
	}
	if notices[0].Title != "花子さんにレスキューされました" {
		t.Fatalf("unexpected first notice Title: %s", notices[0].Title)
	}

	if notices[1].SenderType != "user-002" {
		t.Fatalf("unexpected second notice SenderType: %s", notices[1].SenderType)
	}
	if notices[1].Title != "花子さんをレスキューしました" {
		t.Fatalf("unexpected second notice Title: %s", notices[1].Title)
	}

	if notices[2].SenderType != "user-002" {
		t.Fatalf("unexpected third notice SenderType: %s", notices[2].SenderType)
	}
	if notices[2].Title != "花子さんから汚さ3分の攻撃が届きました" {
		t.Fatalf("unexpected third notice Title: %s", notices[2].Title)
	}
}

// TestGetNotices_Empty は通知が存在しないユーザーに対して空配列を返すことを確認する
func TestGetNotices_Empty(t *testing.T) {
	truncateUsers(t)
	truncateNotices(t)

	createUser(t, "user-001", "太郎", "pineTree", "icon1")

	notices, err := services.GetNotices("user-001")
	if err != nil {
		t.Fatalf("GetNotices failed: %v", err)
	}
	if len(notices) != 0 {
		t.Fatalf("expected 0 notices, got %d: %+v", len(notices), notices)
	}
}
