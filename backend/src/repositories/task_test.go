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
