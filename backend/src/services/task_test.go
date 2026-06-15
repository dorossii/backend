package services_test

import (
	"backend/batch"
	"backend/models"
	"backend/services"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

// 画像アップロードテストに利用するヘルパー
func createFileHeader(t *testing.T, path string) *multipart.FileHeader {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(
		"file",
		filepath.Base(path),
	)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(part, file); err != nil {
		t.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/upload",
		body,
	)

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	_, fileHeader, err := req.FormFile("file")
	if err != nil {
		t.Fatal(err)
	}

	return fileHeader
}

// 画像アップロード(正常系)
func TestPostUploadImage_JPEG(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-jpeg",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusImcomplete,
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	fileHeader := createFileHeader(
		t,
		"../assets/test-images/test.jpg",
	)

	err := services.PostUploadImage(
		"user-001",
		"task-jpeg",
		fileHeader,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var updatedTask models.Task

	if err := models.DB.
		First(&updatedTask, "task_id = ?", "task-jpeg").
		Error; err != nil {
		t.Fatal(err)
	}

	if updatedTask.ImageID == "" {
		t.Fatal("expected image id")
	}
}

// 画像アップロード(画像差し替え)
func TestPostUploadImage_ReplaceImage(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-replace-image",
		BaseID:    "base-001",
		UserID:    "user-001",
		ImageID:   "old-image.jpg",
		Status:    models.TaskStatusImcomplete,
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	fileHeader := createFileHeader(t, "../assets/test-images/test.jpg")

	err := services.PostUploadImage(
		"user-001",
		"task-replace-image",
		fileHeader,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var updatedTask models.Task

	if err := models.DB.
		First(&updatedTask, "task_id = ?", "task-replace-image").
		Error; err != nil {
		t.Fatal(err)
	}

	if updatedTask.ImageID == "" {
		t.Fatal("expected image id")
	}

	if updatedTask.ImageID == "old-image.jpg" {
		t.Fatal("image was not replaced")
	}
}

// 画像アップロード(異常系：タスク不在)
func TestPostUploadImage_TaskNotFound(t *testing.T) {
	fileHeader := createFileHeader(t, "../assets/test-images/test.jpg")

	err := services.PostUploadImage(
		"user-001",
		"not-found-task",
		fileHeader,
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrTaskNotFound {
		t.Fatalf("unexpected error: %v", err)
	}
}

// 画像アップロード(異常系：他人のタスク)
func TestPostUploadImage_PermissionDenied(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-other-user",
		BaseID:    "base-001",
		UserID:    "user-999",
		Status:    models.TaskStatusImcomplete,
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now().Add(1 * time.Hour),
	}
	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	fileHeader := createFileHeader(t, "../assets/test-images/test.jpg")

	err := services.PostUploadImage(
		"user-001",
		"task-other-user",
		fileHeader,
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrTaskPermissionDenied {
		t.Fatalf("unexpected error: %v", err)
	}
}

// 画像アップロード(異常系：JPEG/PNG以外)
func TestPostUploadImage_UnsupportedImageType(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-invalid-image",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusImcomplete,
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	fileHeader := createFileHeader(t, "../assets/test-images/test.txt")

	err := services.PostUploadImage(
		"user-001",
		"task-invalid-image",
		fileHeader,
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrUnsupportedImageType {
		t.Fatalf("unexpected error: %v", err)
	}
}
