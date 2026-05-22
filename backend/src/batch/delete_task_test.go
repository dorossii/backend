package batch_test

import (
	"backend/batch"
	"backend/models"
	"testing"
	"time"
)

func TestDeleteTask(t *testing.T) {
	//テストデータの作成
	// 期限切れのタスク
	task1 := models.Task{
		TaskID:       "task-001",
		UserID:       "user-001",
		Status:       0,
		StartTime:    time.Now().Add(-10 * time.Hour),
		EndTime:      time.Now().Add(-1 * time.Hour), // 1時間前に期限切れ
		ImageID:      "",
		RequireImage: false,
	} 
	// 期限切れでない
	task2 := models.Task{
		TaskID:       "task-002",
		UserID:       "user-002",
		Status:       0,
		StartTime:    time.Now().Add(-10 * time.Hour),
		EndTime:      time.Now().Add(1 * time.Hour), // 1時間後に期限切れ
		ImageID:      "",
		RequireImage: false,
	} 
	if err := models.DB.Create(&task1).Error; err != nil {
		t.Fatalf("Failed to create task1: %v", err)
	}
	if err := models.DB.Create(&task2).Error; err != nil {
		t.Fatalf("Failed to create task2: %v", err)
	}

	//テスト対象の関数を実行
	err := batch.DeleteTask()
	if err != nil {
		t.Fatalf("Failed to delete tasks: %v", err)
	}

	//結果の検証
	var tasks []models.Task
	err = models.DB.Find(&tasks).Error

	if err != nil {
		t.Fatalf("Failed to find tasks: %v", err)
	}

	//期限切れのタスクが削除されているかを検証
	if len(tasks) != 1 {
		t.Fatalf("Unexpected number of tasks: got %d, want 1", len(tasks))
	}

	//残っているタスクがtask2であることを検証
	if tasks[0].TaskID != task2.TaskID {
		t.Errorf("Unexpected task remaining: got TaskID %s, want TaskID %s", tasks[0].TaskID, task2.TaskID)
	}
}


