package models_test

import (
	"backend/models"
	"testing"
	"time"

	"gorm.io/gorm"
) 

func TestTask(t *testing.T) {

	// テーブル初期化
	err := models.DB.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.Task{}).Error

	if err != nil {
		t.Fatal(err)
	}

	// INSERT
	tasks := []models.Task{
		{
			TaskID:    "task-001",
			UserID:    "user-001",
			Status:    models.TaskStatusImcomplete,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(24 * time.Hour),
		},
		{
			TaskID:    "task-002",
			UserID:    "user-001",
			Status:    models.TaskStatusPending,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(24 * time.Hour),
		},
		{
			TaskID:    "task-003",
			UserID:    "user-001",
			Status:    models.TaskStatusCompleted,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(24 * time.Hour),
		},
	}

	err = models.DB.Create(&tasks).Error
	if err != nil {
		t.Fatal(err)
	}

	// SELECT
	var result models.Task

	err = models.DB.First(
		&result,
		"task_id = ?",
		"task-001",
	).Error

	if err != nil {
		t.Fatal(err)
	}

	// 検証(insertしたデータと一致してるか確認)
	if result.UserID != "user-001" {
		t.Fatalf(
			"unexpected UserID: %s",
			result.UserID,
		)
	}

	if result.Status != models.TaskStatusImcomplete {
		t.Fatalf(
			"unexpected Status: %d",
			result.Status,
		)
	}

	if result.RequireImage != false {
		t.Fatalf(
			"unexpected RequireImage: %v",
			result.RequireImage,
		)
	}
}