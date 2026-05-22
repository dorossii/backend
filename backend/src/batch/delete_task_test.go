package batch_test

import (
	"backend/batch"
	"backend/models"
	"log"
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
	//テスト実行前のタスクの表示
	var beforeTasks []models.Task
	if err := models.DB.Find(&beforeTasks).Error; err != nil {
		t.Fatalf("Failed to find tasks: %v", err)
	}	
	//タスクの表示
	log.Println("タスクの一覧(実行前):")
	for _, task := range beforeTasks {
		log.Printf("- TaskID: %s, UserID: %s, EndTime: %s", task.TaskID, task.UserID, task.EndTime)
	}

	//テスト対象の関数を実行
	err := batch.DeleteTask()
	if err != nil {
		t.Fatalf("Failed to delete tasks: %v", err)
	}

	//削除後のタスクを取得
	var afterTasks []models.Task
	if err := models.DB.Find(&afterTasks).Error; err != nil {
		t.Fatalf("Failed to find tasks: %v", err)
	}
	//削除後のタスクの表示
	log.Println("タスクの一覧(実行後):")
	for _, task := range afterTasks {
		log.Printf("- TaskID: %s, UserID: %s", task.TaskID, task.UserID)
	}

	if err != nil {
		t.Fatalf("Failed to find tasks: %v", err)
	}

	//期限切れのタスクが削除されているかを検証
	if len(afterTasks) != 1 {
		t.Fatalf("Unexpected number of tasks: got %d, want 1", len(afterTasks))
	}

	//残っているタスクがtask2であることを検証
	if afterTasks[0].TaskID != task2.TaskID {
		t.Errorf("Unexpected task remaining: got TaskID %s, want TaskID %s", afterTasks[0].TaskID, task2.TaskID)
	}
}


