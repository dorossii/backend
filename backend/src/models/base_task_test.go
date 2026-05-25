package models_test

import (
	"backend/models"
	"testing"
)

func TestBaseTask(t *testing.T) {
	// INSERT
	tasks := []models.BaseTask{
		{
			BaseID:          "task-001",
			TaskName:        "部屋掃除",
			Description: "部屋を掃除するタスク",
			DifficultyLevel: 3,
			DueTime:         3,
			ImageFlag:       true,
			Tags:            models.TaskTagCleaning,
		},
		{
			BaseID:          "task-002",
			TaskName:        "洗濯物を干す",
			Description: "洗濯物を干すタスク",
			DifficultyLevel: 2,
			DueTime:         3,
			ImageFlag:       false,
			Tags:            models.TaskTagLaundry,
		},
		{
			BaseID:          "task-003",
			TaskName:        "夕飯を作る",
			Description: "夕飯を作るタスク",
			DifficultyLevel: 4,
			DueTime:         3,
			ImageFlag:       false,
			Tags:            models.TaskTagCooking,
		},
		{
			BaseID:          "task-004",
			TaskName:        "ゴミを出す",
			Description: "ゴミを出すタスク",
			DifficultyLevel: 1,
			DueTime:         3,
			ImageFlag:       false,
			Tags:            models.TaskTagTrash,
		},
	}

	err := models.DB.Create(&tasks).Error
	if err != nil {
		t.Fatal(err)
	}

	// SELECT
	var result models.BaseTask

	err = models.DB.First(
		&result,
		"base_id = ?",
		"task-001",
	).Error

	if err != nil {
		t.Fatal(err)
	}

	// 検証(insertしたデータと一致してるか確認)
	if result.TaskName != "部屋掃除" {
		t.Fatalf(
			"unexpected TaskName: %s",
			result.TaskName,
		)
	}

	if result.Description != "部屋を掃除するタスク" {
		t.Fatalf(
			"unexpected Description: %s",
			result.Description,
		)
	}

	if result.DifficultyLevel != 3 {
		t.Fatalf(
			"unexpected DifficultyLevel: %d",
			result.DifficultyLevel,
		)
	}

	if result.Tags != models.TaskTagCleaning {
		t.Fatal("unexpected tag")
	}

}
