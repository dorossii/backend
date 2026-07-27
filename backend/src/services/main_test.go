package services_test

import (
	"backend/models"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	models.InitForTest()
	seedBaseTasksForRegisterUser()
	os.Exit(m.Run())
}

// seedBaseTasksForRegisterUser は RegisterUser がユーザー登録時に
// タスクを2つ自動生成する仕様(repositories.CreateTaskForUser)を満たすため、
// テストDB全体で共有する最低限のBaseTaskを用意する。
func seedBaseTasksForRegisterUser() {
	baseTasks := []models.BaseTask{
		{
			BaseID:          "test-fixture-base-task-1",
			TaskName:        "テスト用タスク1",
			Description:     "TestMainが用意する固定BaseTask",
			DifficultyLevel: 1,
			DueTime:         7,
			ImageFlag:       false,
			Tags:            models.TaskTagCleaning,
		},
		{
			BaseID:          "test-fixture-base-task-2",
			TaskName:        "テスト用タスク2",
			Description:     "TestMainが用意する固定BaseTask",
			DifficultyLevel: 1,
			DueTime:         7,
			ImageFlag:       false,
			Tags:            models.TaskTagLaundry,
		},
	}

	if err := models.DB.Create(&baseTasks).Error; err != nil {
		panic("failed to seed base tasks for TestMain: " + err.Error())
	}
}
