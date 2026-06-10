package services_test

import (
	"backend/batch"
	"backend/models"
	"backend/services"
	"fmt"
	"log"
	"testing"
	"time"
)

// services.GetTasksのテスト
func TestGetTasks(t *testing.T) {
	// テスト用のユーザー作成
	if err := CreateSampleUser(); err != nil {
		t.Fatalf("テストユーザーの作成に失敗: %v", err)
	}

	//　ベースタスクの準備(DueTimeは日数単位)
	baseTasks := []models.BaseTask{
		{
			BaseID:          "base-001",
			TaskName:        "部屋掃除",
			DueTime:         1,
			ImageFlag:       true,
			Description:     "掃除して部屋をきれいにしよう",
			DifficultyLevel: 2,
			Tags:            0,
		},
		{
			BaseID:          "base-002",
			TaskName:        "洗濯物を干す",
			DueTime:         2,
			ImageFlag:       false,
			Description:     "洗濯物を干すのを忘れないようにしよう",
			DifficultyLevel: 4,
			Tags:            3,
		},
		{
			BaseID:          "base-003",
			TaskName:        "夕飯を作る",
			DueTime:         3,
			ImageFlag:       false,
			Description:     "夕飯を作ることを忘れないようにしよう",
			DifficultyLevel: 1,
			Tags:            2,
		},
	}
	if err := models.DB.Create(&baseTasks).Error; err != nil {
		t.Fatalf("failed to create dummy base tasks: %v", err)
	}

	err := batch.CreateTask()
	if err != nil {
		t.Fatalf("タスクの作成に失敗: %v", err)
	}

	// タスクを取得
	tasks, err := services.GetTasks("user-010")
	if err != nil {
		t.Fatalf("タスクの取得に失敗: %v", err)
	}

	// 取得したタスクの数を確認
	if len(tasks) == 0 {
		t.Errorf("タスクが見つかりません")
	}

	//取得したタスクの難易度が1から5の範囲内であることを確認
	for _, task := range tasks {
		if task.DifficultyLevel < 1 || task.DifficultyLevel > 5 {
			t.Errorf("タスクの難易度が不正です: %d", task.DifficultyLevel)
		}
	}

	log.Printf("取得したタスク数: %d", len(tasks))

	// 取得したタスクの内容を確認
	for _, task := range tasks {
		log.Printf("タスクID: %s, タスク名: %s, 期限: %s,難易度: %d", task.TaskID, task.TaskName, task.EndTime.Format("2006-01-02"), task.DifficultyLevel)
	}
}

// テスト用のユーザーを作成する関数
func CreateSampleUser() error {
	users := []models.User{
		{
			UserID:      "user-010",
			UserName:    "syatyo",
			BirthDate:   time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
			Mailadress:  "user1@example.com",
			HealthPoint: 1000,
			DirtLevel:   0,
			Combo:       0,
			BgColor:     "#ffb6c1",
		},
	}
	if err := models.DB.Create(&users).Error; err != nil {
		return fmt.Errorf("failed to create dummy users: %v", err)
	}
	return nil
}

func truncateRemindNotices(t *testing.T) {
	t.Helper()

	if err := models.DB.Exec("TRUNCATE TABLE remind_notices").Error; err != nil {
		t.Fatal(err)
	}
}

func TestPostTaskTauntMessage(t *testing.T) {
	truncateFriendShips(t)
	truncateRemindNotices(t)

	TestRegisterUser(t)
	// seedFriend(t) TODO:呼び出し元がマージされてないので一時的に…
	friend := models.FriendShips{
		UserID:   "user-001",
		FriendID: "user-002",
		Status:   models.FriendStatusAccepted,
	}

	if err := models.DB.Create(&friend).Error; err != nil {
		t.Fatal(err)
	}

	err := services.PostTaskTauntMessage(
		"user-001",
		"user-002",
		"お前の部屋きたなすぎ",
	)

	if err != nil {
		t.Fatalf(
			"PostTaskTauntMessage failed: %v",
			err,
		)
	}

	var notice models.RemindNotice

	err = models.DB.
		First(&notice, "user_id = ?", "user-001").
		Error

	if err != nil {
		t.Fatalf(
			"record not found: %v",
			err,
		)
	}

	if notice.SenderID != "user-002" {
		t.Fatalf(
			"unexpected sender id: %s",
			notice.SenderID,
		)
	}

	if notice.Title != "お前の部屋きたなすぎ" {
		t.Fatalf(
			"unexpected title: %s",
			notice.Title,
		)
	}
}

func TestPostTaskTauntMessage_FriendNotFound(t *testing.T) {
	truncateFriendShips(t)
	truncateRemindNotices(t)

	TestRegisterUser(t)

	err := services.PostTaskTauntMessage(
		"user-001",
		"user-999",
		"test message",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrFriendNotFound {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

// タスクステータス更新(完了: 正常系)
func TestPutTaskStatus_Complete(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-Complete",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusImcomplete,
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	resp, err := services.PutTaskStatus(
		"user-001",
		"task-Complete",
		services.TaskStatusComplete,
		"",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.IsChanged {
		t.Fatal("expected IsChanged=true")
	}

	var updatedTask models.Task

	if err := models.DB.
		First(&updatedTask, "task_id = ?", "task-Complete").
		Error; err != nil {
		t.Fatal(err)
	}

	if updatedTask.Status != models.TaskStatusCompleted {
		t.Fatalf(
			"unexpected status: %v",
			updatedTask.Status,
		)
	}
}

// タスクステータス更新(完了: タスク不存在)
func TestPutTaskStatus_TaskNotFound(t *testing.T) {
	TestRegisterUser(t)

	_, err := services.PutTaskStatus(
		"user-001",
		"task-TaskNotFound",
		services.TaskStatusComplete,
		"",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrTaskNotFound {
		t.Fatalf("unexpected error: %v", err)
	}
}

// タスクステータス更新(完了: 有効期間外)
func TestPutTaskStatus_Expired(t *testing.T) {
	CreateSampleUser()

	task := models.Task{
		TaskID:    "task-Expired",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusImcomplete,
		StartTime: time.Now().Add(-2 * time.Hour),
		EndTime:   time.Now().Add(-1 * time.Hour), // 既に終了
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	_, err := services.PutTaskStatus(
		"user-001",
		"task-Expired",
		services.TaskStatusComplete,
		"",
	)

	if err != services.ErrTaskExpired {
		t.Fatalf(
			"expected ErrTaskExpired, got %v",
			err,
		)
	}
}

// タスクステータス更新(認証待ち: 正常系)
func TestPutTaskStatus_Pending(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:       "task-pending",
		BaseID:       "base-001",
		UserID:       "user-001",
		Status:       models.TaskStatusImcomplete,
		StartTime:    time.Now().Add(-1 * time.Hour),
		EndTime:      time.Now().Add(1 * time.Hour),
		ImageID:      "image-001",
		RequireImage: false,
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	resp, err := services.PutTaskStatus(
		"user-001",
		"task-pending",
		services.TaskStatusPending,
		"",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.IsChanged {
		t.Fatal("expected IsChanged=true")
	}

	if resp.RequireImage {
		t.Fatal("expected RequireImage=false")
	}

	var updatedTask models.Task

	if err := models.DB.
		First(&updatedTask, "task_id = ?", "task-pending").
		Error; err != nil {
		t.Fatal(err)
	}

	if updatedTask.Status != models.TaskStatusPending {
		t.Fatalf(
			"unexpected status: %v",
			updatedTask.Status,
		)
	}
}

// タスクステータス更新(承認待ち: 写真必須なのに画像なし)
func TestPutTaskStatus_Pending_RequireImageButNoImageID(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:       "task-pending-no-image",
		BaseID:       "base-001",
		UserID:       "user-001",
		Status:       models.TaskStatusImcomplete,
		StartTime:    time.Now().Add(-1 * time.Hour),
		EndTime:      time.Now().Add(1 * time.Hour),
		ImageID:      "",     // 画像なし
		RequireImage: true,   // 画像必須
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	_, err := services.PutTaskStatus(
		"user-001",
		"task-pending-no-image",
		services.TaskStatusPending,
		"",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrInvalidRequest {
		t.Fatalf("unexpected error: %v", err)
	}
}

// タスクステータス更新(未完了: 正常系)
func TestPutTaskStatus_Incomplete(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-incomplete",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusPending, // Pendingから差し戻し
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	message := "普通に汚い"

	resp, err := services.PutTaskStatus(
		"user-001",
		"task-incomplete",
		services.TaskStatusIncomplete,
		message,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.IsChanged {
		t.Fatal("expected IsChanged=true")
	}

	var updatedTask models.Task

	if err := models.DB.
		First(&updatedTask, "task_id = ?", "task-incomplete").
		Error; err != nil {
		t.Fatal(err)
	}

	if updatedTask.Status != models.TaskStatusImcomplete {
		t.Fatalf(
			"unexpected status: %v",
			updatedTask.Status,
		)
	}

	if updatedTask.Message != message {
		t.Fatalf(
			"unexpected message: %v",
			updatedTask.Message,
		)
	}
}

// タスクステータス更新(未完了: 認証待ち以外から戻そうとする)
func TestPutTaskStatus_Incomplete_NotPending(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-incomplete-not-pending",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusCompleted,
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	_, err := services.PutTaskStatus(
		"user-001",
		"task-incomplete-not-pending",
		services.TaskStatusIncomplete,
		"差し戻し理由",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrTaskStatusAlreadyUpdated {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

// タスクステータス更新(未完了: 拒否理由なし)
func TestPutTaskStatus_Incomplete_EmptyMessage(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-incomplete-empty-message",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusPending,
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	_, err := services.PutTaskStatus(
		"user-001",
		"task-incomplete-empty-message",
		services.TaskStatusIncomplete,
		"",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrInvalidRequest {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

// タスクステータス更新(ステータス不正)
func TestPutTaskStatus_InvalidStatus(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-invalid-status",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusImcomplete,
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	_, err := services.PutTaskStatus(
		"user-001",
		"task-invalid-status",
		"invalid-status",
		"",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrInvalidTaskStatus {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}


// タスクステータス更新(同じステータスへの更新完了)
func TestPutTaskStatus_AlreadyUpdated(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-already-updated",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusPending,
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	_, err := services.PutTaskStatus(
		"user-001",
		"task-already-updated",
		services.TaskStatusPending,
		"",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrTaskStatusAlreadyUpdated {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}