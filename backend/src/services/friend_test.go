package services_test

import (
	"backend/models"
	"backend/services"
	"testing"
)

func truncateFriendShips(t *testing.T) {
	t.Helper()
	if err := models.DB.Exec("TRUNCATE TABLE friend_ships").Error; err != nil {
		t.Fatal(err)
	}
}

func TestSendFriendRequest(t *testing.T) {
	truncateFriendShips(t)

	err := services.SendFriendRequest("user-001", "user-002")
	if err != nil {
		t.Fatalf("SendFriendRequest failed: %v", err)
	}

	var fs models.FriendShips
	if err := models.DB.First(&fs, "user_id = ? AND friend_id = ?", "user-001", "user-002").Error; err != nil {
		t.Fatalf("record not found: %v", err)
	}
	if fs.Status != models.FriendStatusPending {
		t.Fatalf("unexpected status: %d", fs.Status)
	}
}

func TestSendFriendRequest_AlreadySent(t *testing.T) {
	truncateFriendShips(t)

	if err := services.SendFriendRequest("user-001", "user-002"); err != nil {
		t.Fatalf("first SendFriendRequest failed: %v", err)
	}

	err := services.SendFriendRequest("user-001", "user-002")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err != services.ErrAlreadySent {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendFriendRequest_AlreadyReceived(t *testing.T) {
	truncateFriendShips(t)

	// 相手から先に申請済み
	if err := services.SendFriendRequest("user-002", "user-001"); err != nil {
		t.Fatalf("reverse SendFriendRequest failed: %v", err)
	}

	err := services.SendFriendRequest("user-001", "user-002")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err != services.ErrAlreadyReceived {
		t.Fatalf("unexpected error: %v", err)
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

// 嫌がらせ設定
func TestPostAttackerSettings(t *testing.T) {
	truncateFriendShips(t)

	TestRegisterUser(t)
	seedFriend(t)

	err := services.PostAttackerSettings("user-001", "user-002")
	if err != nil {
		t.Fatalf("PostAttackerSettings failed: %v", err)
	}

	var setting models.User

	err = models.DB.
		First(&setting, "user_id = ?", "user-001").
		Error

	if err != nil {
		t.Fatalf("record not found: %v", err)
	}

	if setting.TargetUser != "user-002" {
		t.Fatalf(
			"unexpected target user: %s",
			setting.TargetUser,
		)
	}
}

// 嫌がらせ設定(ランダム攻撃設定)
func TestPostAttackerSettings_RandomMode(t *testing.T) {
	truncateFriendShips(t)

	TestRegisterUser(t)

	err := services.PostAttackerSettings("user-001", "")
	if err != nil {
		t.Fatalf("PostAttackerSettings failed: %v", err)
	}

	var setting models.User

	err = models.DB.
		First(&setting, "user_id = ?", "user-001").
		Error

	if err != nil {
		t.Fatalf("record not found: %v", err)
	}

	if setting.TargetUser != "" {
		t.Fatalf(
			"expected empty target user, got: %s",
			setting.TargetUser,
		)
	}
}

// 嫌がらせ設定(エラー系:フレンドがいない)
func TestPostAttackerSettings_FriendNotFound(t *testing.T) {
	truncateFriendShips(t)

	TestRegisterUser(t)

	err := services.PostAttackerSettings("user-001", "user-999")

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrFriendNotFound {
		t.Fatalf("unexpected error: %v", err)
	}
}