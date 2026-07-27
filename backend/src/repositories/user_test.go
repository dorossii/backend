package repositories_test

import (
	"backend/models"
	"backend/repositories"
	"testing"
	"time"
)

func truncateUsers(t *testing.T) {
	t.Helper()
	if err := models.DB.Exec("TRUNCATE TABLE users").Error; err != nil {
		t.Fatal(err)
	}
}

func truncateBaseTasks(t *testing.T) {
	t.Helper()
	if err := models.DB.Exec("TRUNCATE TABLE base_tasks").Error; err != nil {
		t.Fatal(err)
	}
}

func truncateTasks(t *testing.T) {
	t.Helper()
	if err := models.DB.Exec("TRUNCATE TABLE tasks").Error; err != nil {
		t.Fatal(err)
	}
}

func TestCreateUser(t *testing.T) {
	truncateUsers(t)

	user := &models.User{
		UserID:     "user-001",
		UserName:   "syatyo",
		BirthDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Mailadress: "syatyo@example.com",
	}

	err := repositories.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	var result models.User
	err = models.DB.First(&result, "user_id = ?", "user-001").Error
	if err != nil {
		t.Fatalf("failed to find User: %v", err)
	}

	if result.UserName != "syatyo" {
		t.Fatalf("unexpected UserName: %s", result.UserName)
	}
	if result.Mailadress != "syatyo@example.com" {
		t.Fatalf("unexpected Mailadress: %s", result.Mailadress)
	}
	if !result.BirthDate.Equal(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected BirthDate: %v", result.BirthDate)
	}
}

func TestCreateUser_DuplicateID(t *testing.T) {
	truncateUsers(t)

	user := &models.User{
		UserID:     "user-001",
		UserName:   "syatyo",
		BirthDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Mailadress: "syatyo@example.com",
	}

	if err := repositories.CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// 同じIDで再度登録するとエラーになるか確認
	err := repositories.CreateUser(user)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestCreateTaskForUser_Success(t *testing.T) {
	truncateUsers(t)
	truncateBaseTasks(t)
	truncateTasks(t)

	// テスト用ユーザーを作成
	user := &models.User{
		UserID:     "user-001",
		UserName:   "syatyo",
		BirthDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Mailadress: "syatyo@example.com",
	}
	if err := repositories.CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// テスト用BaseTaskを2件作成
	baseTasks := []models.BaseTask{
		{BaseID: "base-001", DueTime: 1, ImageFlag: false},
		{BaseID: "base-002", DueTime: 3, ImageFlag: true},
	}
	if err := models.DB.Create(&baseTasks).Error; err != nil {
		t.Fatalf("failed to create base tasks: %v", err)
	}

	err := repositories.CreateTaskForUser(user.UserID)
	if err != nil {
		t.Fatalf("CreateTaskForUser failed: %v", err)
	}

	var tasks []models.Task
	if err := models.DB.Where("user_id = ?", user.UserID).Find(&tasks).Error; err != nil {
		t.Fatalf("failed to find tasks: %v", err)
	}

	// 1ユーザーにつき2タスク作成されているか
	if len(tasks) != 2 {
		t.Fatalf("unexpected number of tasks: got %d, want 2", len(tasks))
	}

	// BaseIDが重複していないか（同じタスクが2重に割り当たっていないか）
	seen := map[string]bool{}
	for _, task := range tasks {
		if seen[task.BaseID] {
			t.Fatalf("duplicate BaseID assigned to same user: %s", task.BaseID)
		}
		seen[task.BaseID] = true

		// 各フィールドの検証
		if task.TaskID == "" {
			t.Fatal("TaskID should not be empty")
		}
		if task.UserID != user.UserID {
			t.Fatalf("unexpected UserID: %s", task.UserID)
		}
		if task.Status != models.TaskStatusPending {
			t.Fatalf("unexpected Status: %v", task.Status)
		}
		if task.ImageID != "" {
			t.Fatalf("ImageID should be empty on creation, got: %s", task.ImageID)
		}
		if !task.EndTime.After(task.StartTime) {
			t.Fatalf("EndTime should be after StartTime")
		}
	}
}

func TestCreateTaskForUser_RequireImageFalseWhenImageFlagFalse(t *testing.T) {
	truncateUsers(t)
	truncateBaseTasks(t)
	truncateTasks(t)

	user := &models.User{
		UserID:     "user-002",
		UserName:   "tanuki",
		BirthDate:  time.Date(1995, 5, 5, 0, 0, 0, 0, time.UTC),
		Mailadress: "tanuki@example.com",
	}
	if err := repositories.CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// ImageFlagが全てfalseのBaseTaskのみ用意
	baseTasks := []models.BaseTask{
		{BaseID: "base-101", DueTime: 1, ImageFlag: false},
		{BaseID: "base-102", DueTime: 2, ImageFlag: false},
	}
	if err := models.DB.Create(&baseTasks).Error; err != nil {
		t.Fatalf("failed to create base tasks: %v", err)
	}

	if err := repositories.CreateTaskForUser(user.UserID); err != nil {
		t.Fatalf("CreateTaskForUser failed: %v", err)
	}

	var tasks []models.Task
	if err := models.DB.Where("user_id = ?", user.UserID).Find(&tasks).Error; err != nil {
		t.Fatalf("failed to find tasks: %v", err)
	}

	for _, task := range tasks {
		if task.RequireImage {
			t.Fatalf("RequireImage should be false when ImageFlag is false, task: %s", task.TaskID)
		}
	}
}

func TestCreateTaskForUser_InsufficientBaseTasks(t *testing.T) {
	truncateUsers(t)
	truncateBaseTasks(t)
	truncateTasks(t)

	user := &models.User{
		UserID:     "user-003",
		UserName:   "kitsune",
		BirthDate:  time.Date(1998, 3, 3, 0, 0, 0, 0, time.UTC),
		Mailadress: "kitsune@example.com",
	}
	if err := repositories.CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// BaseTaskを1件のみ作成（2件未満）
	baseTask := models.BaseTask{BaseID: "base-201", DueTime: 1, ImageFlag: false}
	if err := models.DB.Create(&baseTask).Error; err != nil {
		t.Fatalf("failed to create base task: %v", err)
	}

	err := repositories.CreateTaskForUser(user.UserID)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestCreateTaskForUser_NoBaseTasks(t *testing.T) {
	truncateUsers(t)
	truncateBaseTasks(t)
	truncateTasks(t)

	user := &models.User{
		UserID:     "user-004",
		UserName:   "usagi",
		BirthDate:  time.Date(1999, 7, 7, 0, 0, 0, 0, time.UTC),
		Mailadress: "usagi@example.com",
	}
	if err := repositories.CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// BaseTaskを1件も作らない
	err := repositories.CreateTaskForUser(user.UserID)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
