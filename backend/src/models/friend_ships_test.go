package models_test

import (
	"backend/models"
	"testing"
)

func TestFriendShips(t *testing.T) {
	models.Init()

	// テーブル初期化
	err := models.DB.Exec(
		"DELETE FROM friend_ships",
	).Error

	if err != nil {
		t.Fatal(err)
	}

	friendships := []models.FriendShips{
		{
			UserID:   "user-001",
			FriendID: "user-002",
			Status:   models.FriendStatusAccepted,
		},
		{
			UserID:   "user-002",
			FriendID: "user-001",
			Status:   models.FriendStatusAccepted,
		},
		{
			UserID:   "user-001",
			FriendID: "user-003",
			Status:   models.FriendStatusAccepted,
		},
		{
			UserID:   "user-003",
			FriendID: "user-001",
			Status:   models.FriendStatusAccepted,
		},
		{
			UserID:   "user-001",
			FriendID: "user-004",
			Status:   models.FriendStatusPending,
		},
		{
			UserID:   "user-004",
			FriendID: "user-001",
			Status:   models.FriendStatusPending,
		},
		{
			UserID:   "user-002",
			FriendID: "user-003",
			Status:   models.FriendStatusAccepted,
		},
		{
			UserID:   "user-003",
			FriendID: "user-002",
			Status:   models.FriendStatusAccepted,
		},
	}

	err = models.DB.Create(&friendships).Error
	if err != nil {
		t.Fatal(err)
	}

	// SELECT
	var result models.FriendShips

	err = models.DB.First(
		&result,
		"user_id = ? AND friend_id = ?",
		"user-001",
		"user-002",
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

	if result.FriendID != "user-002" {
		t.Fatalf(
			"unexpected FriendID: %s",
			result.FriendID,
		)
	}

	if result.Status != models.FriendStatusAccepted {
		t.Fatalf(
			"unexpected Status: %d",
			result.Status,
		)
	}

}
