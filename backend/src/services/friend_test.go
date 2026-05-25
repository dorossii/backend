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

// user-002 -> user-001 の申請を user-001 が承認する正常系
func TestAcceptFriendRequest(t *testing.T) {
	truncateFriendShips(t)

	if err := services.SendFriendRequest("user-002", "user-001"); err != nil {
		t.Fatalf("SendFriendRequest failed: %v", err)
	}

	if err := services.AcceptFriendRequest("user-001", "user-002"); err != nil {
		t.Fatalf("AcceptFriendRequest failed: %v", err)
	}

	var fs models.FriendShips
	if err := models.DB.First(&fs, "user_id = ? AND friend_id = ?", "user-002", "user-001").Error; err != nil {
		t.Fatalf("record not found: %v", err)
	}
	if fs.Status != models.FriendStatusAccepted {
		t.Fatalf("unexpected status: %d", fs.Status)
	}
}

// リクエストが存在しない場合はエラー
func TestAcceptFriendRequest_NotFound(t *testing.T) {
	truncateFriendShips(t)

	err := services.AcceptFriendRequest("user-001", "user-002")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

// 自分が申請者側の場合は承認不可
func TestAcceptFriendRequest_CannotAcceptOwnRequest(t *testing.T) {
	truncateFriendShips(t)

	// user-001 -> user-002 の申請
	if err := services.SendFriendRequest("user-001", "user-002"); err != nil {
		t.Fatalf("SendFriendRequest failed: %v", err)
	}

	// 申請者 user-001 が自分の申請を承認しようとする（不正）
	err := services.AcceptFriendRequest("user-001", "user-002")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
