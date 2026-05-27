package services_test

import (
	"backend/models"
	"backend/services"
	"testing"
	"time"
)

func truncateFriendShips(t *testing.T) {
	t.Helper()
	if err := models.DB.Exec("TRUNCATE TABLE friend_ships").Error; err != nil {
		t.Fatal(err)
	}
}

func truncateUsers(t *testing.T) {
	t.Helper()
	if err := models.DB.Exec("TRUNCATE TABLE users").Error; err != nil {
		t.Fatal(err)
	}
}

func createUser(t *testing.T, userID, name, icon, bgColor string) {
	t.Helper()
	u := &models.User{UserID: userID, UserName: name, Icon: icon, BgColor: bgColor, BirthDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	if err := models.DB.Create(u).Error; err != nil {
		t.Fatalf("createUser failed: %v", err)
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

// 自分宛の pending リクエスト一覧が取得できる
func TestGetFriendRequests(t *testing.T) {
	truncateFriendShips(t)

	// user-002, user-003 から user-001 へ申請
	if err := services.SendFriendRequest("user-002", "user-001"); err != nil {
		t.Fatalf("SendFriendRequest failed: %v", err)
	}
	if err := services.SendFriendRequest("user-003", "user-001"); err != nil {
		t.Fatalf("SendFriendRequest failed: %v", err)
	}

	reqs, err := services.GetFriendRequests("user-001")
	if err != nil {
		t.Fatalf("GetFriendRequests failed: %v", err)
	}
	if len(reqs) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(reqs))
	}

	requesters := map[string]bool{}
	for _, r := range reqs {
		requesters[r.RequestUser] = true
	}
	if !requesters["user-002"] || !requesters["user-003"] {
		t.Fatalf("unexpected requesters: %v", requesters)
	}
}

// 自分が送った申請は一覧に含まれない
func TestGetFriendRequests_ExcludesSentRequests(t *testing.T) {
	truncateFriendShips(t)

	// user-001 から user-002 へ申請（自分が送った側）
	if err := services.SendFriendRequest("user-001", "user-002"); err != nil {
		t.Fatalf("SendFriendRequest failed: %v", err)
	}

	reqs, err := services.GetFriendRequests("user-001")
	if err != nil {
		t.Fatalf("GetFriendRequests failed: %v", err)
	}
	if len(reqs) != 0 {
		t.Fatalf("expected 0 requests, got %d", len(reqs))
	}
}

// 承認済みのリクエストは一覧に含まれない
func TestGetFriendRequests_ExcludesAccepted(t *testing.T) {
	truncateFriendShips(t)

	if err := services.SendFriendRequest("user-002", "user-001"); err != nil {
		t.Fatalf("SendFriendRequest failed: %v", err)
	}
	if err := services.AcceptFriendRequest("user-001", "user-002"); err != nil {
		t.Fatalf("AcceptFriendRequest failed: %v", err)
	}

	reqs, err := services.GetFriendRequests("user-001")
	if err != nil {
		t.Fatalf("GetFriendRequests failed: %v", err)
	}
	if len(reqs) != 0 {
		t.Fatalf("expected 0 requests after accept, got %d", len(reqs))
	}
}

// リクエストが0件の場合は空スライスを返す
func TestGetFriendRequests_Empty(t *testing.T) {
	truncateFriendShips(t)

	reqs, err := services.GetFriendRequests("user-001")
	if err != nil {
		t.Fatalf("GetFriendRequests failed: %v", err)
	}
	if len(reqs) != 0 {
		t.Fatalf("expected 0 requests, got %d", len(reqs))
	}
}

// 承認済みフレンドが正しく取得できる（自分が申請した場合）
func TestGetFriends_AsSender(t *testing.T) {
	truncateFriendShips(t)
	truncateUsers(t)

	createUser(t, "user-001", "Alice", "cat", "#ff0000")
	createUser(t, "user-002", "Bob", "dog", "#00ff00")

	if err := services.SendFriendRequest("user-001", "user-002"); err != nil {
		t.Fatalf("SendFriendRequest failed: %v", err)
	}
	if err := services.AcceptFriendRequest("user-002", "user-001"); err != nil {
		t.Fatalf("AcceptFriendRequest failed: %v", err)
	}

	friends, err := services.GetFriends("user-001")
	if err != nil {
		t.Fatalf("GetFriends failed: %v", err)
	}
	if len(friends) != 1 {
		t.Fatalf("expected 1 friend, got %d", len(friends))
	}
	f := friends[0]
	if f.UserID != "user-002" {
		t.Errorf("unexpected user_id: %s", f.UserID)
	}
	if f.Name != "Bob" {
		t.Errorf("unexpected name: %s", f.Name)
	}
	if f.Icon != "dog" {
		t.Errorf("unexpected icon: %s", f.Icon)
	}
	if f.Background != "#00ff00" {
		t.Errorf("unexpected background: %s", f.Background)
	}
}

// 承認済みフレンドが正しく取得できる（相手が申請した場合）
func TestGetFriends_AsReceiver(t *testing.T) {
	truncateFriendShips(t)
	truncateUsers(t)

	createUser(t, "user-001", "Alice", "cat", "#ff0000")
	createUser(t, "user-002", "Bob", "dog", "#00ff00")

	if err := services.SendFriendRequest("user-002", "user-001"); err != nil {
		t.Fatalf("SendFriendRequest failed: %v", err)
	}
	if err := services.AcceptFriendRequest("user-001", "user-002"); err != nil {
		t.Fatalf("AcceptFriendRequest failed: %v", err)
	}

	friends, err := services.GetFriends("user-001")
	if err != nil {
		t.Fatalf("GetFriends failed: %v", err)
	}
	if len(friends) != 1 {
		t.Fatalf("expected 1 friend, got %d", len(friends))
	}
	if friends[0].UserID != "user-002" {
		t.Errorf("unexpected user_id: %s", friends[0].UserID)
	}
}

// pending 状態のフレンドは一覧に含まれない
func TestGetFriends_ExcludesPending(t *testing.T) {
	truncateFriendShips(t)
	truncateUsers(t)

	createUser(t, "user-001", "Alice", "cat", "#ff0000")
	createUser(t, "user-002", "Bob", "dog", "#00ff00")

	if err := services.SendFriendRequest("user-001", "user-002"); err != nil {
		t.Fatalf("SendFriendRequest failed: %v", err)
	}

	friends, err := services.GetFriends("user-001")
	if err != nil {
		t.Fatalf("GetFriends failed: %v", err)
	}
	if len(friends) != 0 {
		t.Fatalf("expected 0 friends, got %d", len(friends))
	}
}

// フレンドが複数いる場合に全員返る
func TestGetFriends_Multiple(t *testing.T) {
	truncateFriendShips(t)
	truncateUsers(t)

	createUser(t, "user-001", "Alice", "cat", "#ff0000")
	createUser(t, "user-002", "Bob", "dog", "#00ff00")
	createUser(t, "user-003", "Carol", "bird", "#0000ff")

	if err := services.SendFriendRequest("user-001", "user-002"); err != nil {
		t.Fatalf("SendFriendRequest failed: %v", err)
	}
	if err := services.AcceptFriendRequest("user-002", "user-001"); err != nil {
		t.Fatalf("AcceptFriendRequest failed: %v", err)
	}
	if err := services.SendFriendRequest("user-003", "user-001"); err != nil {
		t.Fatalf("SendFriendRequest failed: %v", err)
	}
	if err := services.AcceptFriendRequest("user-001", "user-003"); err != nil {
		t.Fatalf("AcceptFriendRequest failed: %v", err)
	}

	friends, err := services.GetFriends("user-001")
	if err != nil {
		t.Fatalf("GetFriends failed: %v", err)
	}
	if len(friends) != 2 {
		t.Fatalf("expected 2 friends, got %d", len(friends))
	}

	ids := map[string]bool{}
	for _, f := range friends {
		ids[f.UserID] = true
	}
	if !ids["user-002"] || !ids["user-003"] {
		t.Errorf("unexpected friend ids: %v", ids)
	}
}

// フレンドが0件の場合は空スライスを返す
func TestGetFriends_Empty(t *testing.T) {
	truncateFriendShips(t)
	truncateUsers(t)

	friends, err := services.GetFriends("user-001")
	if err != nil {
		t.Fatalf("GetFriends failed: %v", err)
	}
	if len(friends) != 0 {
		t.Fatalf("expected 0 friends, got %d", len(friends))
	}
}
