package repositories_test

import (
	"backend/models"
	"backend/repositories"
	"testing"
	"time"
)

// GetUserTasksがbase_tasks.tagsをTaskResponse.Tagへ正しくマッピングすることを確認する。
// 修正前は "base_tasks.tags" にエイリアスがなく、構造体フィールド"Tag"(json:"tag")と
// カラム名"tags"が一致しないためGORMがScanできず、常に0(ゼロ値)が返っていた。
func TestGetUserTasks_TagMapping(t *testing.T) {
	user := models.User{
		UserID:     "repo-test-user-001",
		UserName:   "repo-test-user",
		BirthDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Mailadress: "repo-test-user@example.com",
	}
	if err := models.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create dummy user: %v", err)
	}

	baseTask := models.BaseTask{
		BaseID:          "repo-test-base-tag",
		TaskName:        "洗濯物を干す",
		Description:     "洗濯物を干すのを忘れないようにしよう",
		DifficultyLevel: 4,
		DueTime:         2,
		Tags:            models.TaskTagAther, // 3
	}
	if err := models.DB.Create(&baseTask).Error; err != nil {
		t.Fatalf("failed to create dummy base task: %v", err)
	}

	task := models.Task{
		TaskID:    "repo-test-task-tag",
		BaseID:    baseTask.BaseID,
		UserID:    user.UserID,
		Status:    models.TaskStatusIncomplete,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	}
	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatalf("failed to create dummy task: %v", err)
	}

	results, err := repositories.GetUserTasks(user.UserID)
	if err != nil {
		t.Fatalf("GetUserTasks failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 task, got %d", len(results))
	}

	got := results[0]
	if got.TaskID != task.TaskID {
		t.Fatalf("unexpected TaskID: %s", got.TaskID)
	}

	// 修正前はここが常に0になっていた
	if got.Tag != int(models.TaskTagAther) {
		t.Fatalf("unexpected Tag: got %d, want %d", got.Tag, int(models.TaskTagAther))
	}
}

// GetUserTasksのImageIdが、base_tasksのマスタ画像ではなく、
// tasksに保存されたユーザーのアップロード画像(完了証明写真)のIDを返すことを確認する。
func TestGetUserTasks_ImageIdIsTaskUploadedImage(t *testing.T) {
	user := models.User{
		UserID:     "repo-test-user-002",
		UserName:   "repo-test-user-2",
		BirthDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Mailadress: "repo-test-user-2@example.com",
	}
	if err := models.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create dummy user: %v", err)
	}

	baseTask := models.BaseTask{
		BaseID:          "repo-test-base-image",
		TaskName:        "食材の買い出し",
		Description:     "スーパーで翌日以降に必要な食材リストを確認しながら購入する。",
		DifficultyLevel: 3,
		DueTime:         2,
		ImageId:         "kaimono.webp", // base_tasksのマスタ画像(ファイル名)
		Tags:            models.TaskTagAther,
	}
	if err := models.DB.Create(&baseTask).Error; err != nil {
		t.Fatalf("failed to create dummy base task: %v", err)
	}

	task := models.Task{
		TaskID:    "repo-test-task-image",
		BaseID:    baseTask.BaseID,
		UserID:    user.UserID,
		Status:    models.TaskStatusPending,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
		ImageID:   "uploaded-image-uuid", // ユーザーがアップロードした完了証明写真のID
	}
	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatalf("failed to create dummy task: %v", err)
	}

	results, err := repositories.GetUserTasks(user.UserID)
	if err != nil {
		t.Fatalf("GetUserTasks failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 task, got %d", len(results))
	}

	got := results[0]
	if got.ImageID != task.ImageID {
		t.Fatalf("expected ImageID %s (task's uploaded image), got %s", task.ImageID, got.ImageID)
	}
}

func TestGetMessageUserId(t *testing.T) {
	// テスト用ユーザーを作成
	user := models.User{
		UserID:     "test-user-001",
		UserName:   "Test User",
		BirthDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Mailadress: "test-user-001@example.com",
	}
	if err := repositories.CreateUser(&user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// テスト用フレンドを2人作成
	friend1 := models.User{
		UserID:     "test-friend-001",
		UserName:   "Test Friend 1",
		BirthDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Mailadress: "test-friend-001@example.com",
	}
	friend2 := models.User{
		UserID:     "test-friend-002",
		UserName:   "Test Friend 2",
		BirthDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Mailadress: "test-friend-002@example.com",
	}
	if err := repositories.CreateUser(&friend1); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if err := repositories.CreateUser(&friend2); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// user -> friend1 (userが申請した方向), friend2 -> user (friend2が申請した方向) の
	// 双方向の友達関係を作成し、方向に依らず候補として拾えることを確認する
	if err := repositories.CreateFriendShip(&models.FriendShips{
		UserID:   user.UserID,
		FriendID: friend1.UserID,
		Status:   models.FriendStatusAccepted,
	}); err != nil {
		t.Fatalf("CreateFriendShip failed: %v", err)
	}
	if err := repositories.CreateFriendShip(&models.FriendShips{
		UserID:   friend2.UserID,
		FriendID: user.UserID,
		Status:   models.FriendStatusAccepted,
	}); err != nil {
		t.Fatalf("CreateFriendShip failed: %v", err)
	}

	// friend1 にはすでにリマインド通知済みなので候補から除外される
	message := models.RemindNotice{
		NoticeID:   "test-notice-001",
		UserID:     friend1.UserID,
		SenderID:   user.UserID,
		Title:      "Test Notice",
		NotifiedAt: time.Now(),
		IsRead:     false,
	}
	if err := models.DB.Create(&message).Error; err != nil {
		t.Fatalf("failed to create test message: %v", err)
	}

	// メッセージ対象を取得
	result, err := repositories.GetMessageUserId(user.UserID)
	if err != nil {
		t.Fatalf("GetMessageUserId failed: %v", err)
	}

	// friend1 は通知済みなので除外され、friend2 が唯一の候補として返るはず
	if result != friend2.UserID {
		t.Errorf("expected %s, got %s", friend2.UserID, result)
	}
}

func TestGetMessageUserId_NoEligibleFriend(t *testing.T) {
	// フレンドが1人もいないユーザーで呼び出した場合、エラーではなく空文字列が返ることを確認する
	user := models.User{
		UserID:     "test-user-002",
		UserName:   "Test User No Friends",
		BirthDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Mailadress: "test-user-002@example.com",
	}
	if err := repositories.CreateUser(&user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	result, err := repositories.GetMessageUserId(user.UserID)
	if err != nil {
		t.Fatalf("GetMessageUserId failed: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty result, got %s", result)
	}
}
