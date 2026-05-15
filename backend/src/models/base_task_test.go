package models_test

import (
	"backend/models"
	"testing"

	"gorm.io/gorm"
)

func TestBaseTask(t *testing.T) {
	models.Init() // TODO: test毎に接続するのはキモい気がする。

	// テーブル初期化
	err := models.DB.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.BaseTask{}).Error

	if err != nil {
		t.Fatal(err)
	}

	// INSERT
	tasks := []models.BaseTask{
		{
			TaskID:          "task-001",
			TaskName:        "部屋掃除",
			DifficultyLevel: 3,
			DueTime:         3,
			ImageFlag:       true,
			Tags:            models.TaskTagCleaning,
		},
		{
			TaskID:          "task-002",
			TaskName:        "洗濯物を干す",
			DifficultyLevel: 2,
			DueTime:         3,
			ImageFlag:       false,
			Tags:            models.TaskTagLaundry,
		},
		{
			TaskID:          "task-003",
			TaskName:        "夕飯を作る",
			DifficultyLevel: 4,
			DueTime:         3,
			ImageFlag:       false,
			Tags:            models.TaskTagCooking,
		},
		{
			TaskID:          "task-004",
			TaskName:        "ゴミを出す",
			DifficultyLevel: 1,
			DueTime:         3,
			ImageFlag:       false,
			Tags:            models.TaskTagTrash,
		},
	}

	err = models.DB.Create(&tasks).Error
	if err != nil {
		t.Fatal(err)
	}

	// SELECT
	var result models.BaseTask

	err = models.DB.First(
		&result,
		"task_id = ?",
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
