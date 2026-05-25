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

	log.Printf("取得したタスク数: %d", len(tasks))

	// 取得したタスクの内容を確認
	for _, task := range tasks {
		log.Printf("タスクID: %s, タスク名: %s, 期限: %s", task.TaskID, task.TaskName, task.EndTime.Format("2006-01-02"))
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
