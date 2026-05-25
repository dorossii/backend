package services_test

import (
	"backend/models"
	"backend/services"
	"testing"
)

func truncateRemindNotices(t *testing.T) {
	t.Helper()

	if err := models.DB.Exec(
		"TRUNCATE TABLE remind_notices",
	).Error; err != nil {
		t.Fatal(err)
	}
}

func seedFriend(t *testing.T) {
	t.Helper()

	friend := models.FriendShips{
		UserID:   "user-001",
		FriendID: "user-002",
		Status:   models.FriendStatusAccepted,
	}

	if err := models.DB.Create(&friend).Error; err != nil {
		t.Fatal(err)
	}
}

func TestPostTaskTauntMessage(t *testing.T) {

	truncateFriendShips(t)
	truncateRemindNotices(t)

	seedFriend(t)

	service := services.TaskService{}
	err := service.PostTaskTauntMessage(
		"user-001",
		"user-002",
		"hello",
	)

	if err != nil {
		t.Fatalf(
			"PostTaskTauntMessage failed: %v",
			err,
		)
	}

	var notice models.RemindNotice

	err = models.DB.First(
		&notice,
		"user_id = ? AND sender_id = ?",
		"user-001",
		"user-002",
	).Error

	if err != nil {
		t.Fatalf("record not found: %v", err)
	}

	if notice.Title != "hello" {
		t.Fatalf(
			"unexpected title: %s",
			notice.Title,
		)
	}
}