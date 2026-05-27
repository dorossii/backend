package batch_test

import (
	"backend/batch"
	"backend/models"
	"testing"
	"gorm.io/gorm"
)

func TestCreateTask(t *testing.T) {
	//　ベースタスクの準備(DueTimeは日数単位)
	baseTasks := []models.BaseTask{
		{BaseID: "base-001", TaskName: "部屋掃除", DueTime: 1, ImageFlag: true},
		{BaseID: "base-002", TaskName: "洗濯物を干す", DueTime: 2, ImageFlag: false},
		{BaseID: "base-003", TaskName: "夕飯を作る", DueTime: 3, ImageFlag: false},
	}
	if err := models.DB.Create(&baseTasks).Error; err != nil {
		t.Fatalf("failed to create dummy base tasks: %v", err)
	}

	// DueTime（日数）をキーにして BaseTask を逆引きできるマップを作成
	baseTaskMap := make(map[int]models.BaseTask)
	for _, bt := range baseTasks {
		baseTaskMap[bt.DueTime] = bt
	}

	// テスト対象の関数を実行
	err := batch.CreateTask()
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	// 結果の検証 (SELECT)
	var createdTasks []models.Task
	if err := models.DB.Find(&createdTasks).Error; err != nil {
		t.Fatal(err)
	}

	// 検証①: 生成された総タスク数の確認 (ユーザー2人 × 2タスク = 4つ)
	expectedTotalTasks := 6
	if len(createdTasks) != expectedTotalTasks {
		t.Fatalf("unexpected total tasks: got %d, want %d", len(createdTasks), expectedTotalTasks)
	}

	// ユーザーごとの生成状況をマップに分類
	tasksByUser := make(map[string][]models.Task)
	for _, task := range createdTasks {
		tasksByUser[task.UserID] = append(tasksByUser[task.UserID], task)
	}

	// 結果の検証 (SELECT)
	var users []models.User
	if err := models.DB.Find(&users).Error; err != nil {
		t.Fatal(err)
	}

	for _, user := range users {
		userTasks := tasksByUser[user.UserID]

		// 検証②: 各ユーザーに過不足なく2つずつ割り振られているか
		if len(userTasks) != 2 {
			t.Fatalf("unexpected task count for user %s: got %d, want 2", user.UserID, len(userTasks))
		}

		for _, task := range userTasks {
			// 初期ステータスがPending（保留中）になっているか
			if task.Status != models.TaskStatusPending {
				t.Errorf("unexpected task status for task %s: got %v, want TaskStatusPending", task.TaskID, task.Status)
			}

			// EndTime から正確な「日数（Days）」を逆算する（24時間で割って四捨五入）
			duration := task.EndTime.Sub(task.StartTime)
			calculatedDueDay := int((duration.Hours() / 24.0) + 0.5)

			// 計算された日数から元になった BaseTask を特定
			matchedBaseTask, exists := baseTaskMap[calculatedDueDay]
			if !exists {
				t.Errorf("failed to map generated task to any BaseTask by DueTime. (TaskID: %s, Calculated DueTime: %d days)",
					task.TaskID, calculatedDueDay)
				continue
			}

			t.Logf("[生成確認] User: %s のタスク (ID: %s) は BaseTask: %s (%s) から正しく生成されました。(設定期間: %d日間)画像要求: %v",
				task.UserID, task.TaskID, matchedBaseTask.BaseID, matchedBaseTask.TaskName, calculatedDueDay, task.RequireImage)
		}

		// 同じユーザー内で重複タスク（EndTimeが同時刻）が割り当たっていないか
		if userTasks[0].EndTime.Equal(userTasks[1].EndTime) {
			t.Errorf("user %s has duplicate tasks assigned (same EndTime: %v)", user.UserID, userTasks[0].EndTime)
		}
	}
}

// ベースタスクが2つ未満の状態で CreateTask を実行したときにエラーになるかのテスト
func TestCreateTask_InsufficientBaseTasks(t *testing.T) {

	models.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Task{})
	models.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.BaseTask{})

	baseTask := models.BaseTask{BaseID: "base-001", TaskName: "部屋掃除"}
	if err := models.DB.Create(&baseTask).Error; err != nil {
		t.Fatalf("failed to create dummy base task: %v", err)
	}

	err := batch.CreateTask()
	if err == nil {
		t.Fatal("expected error due to insufficient base tasks, but got nil")
	}

	expectedErr := "insufficient base tasks available (minimum 2 required)"
	if err.Error() != expectedErr {
		t.Fatalf("unexpected error message: got %q, want %q", err.Error(), expectedErr)
	}
}
