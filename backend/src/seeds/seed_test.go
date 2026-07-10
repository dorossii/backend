package seeds_test

import (
	"backend/models"
	"backend/seeds"
	"testing"
)

func TestSeed(t *testing.T) {
	
	err := seeds.Seed()
	if err != nil {
		t.Errorf("Failed to seed data: %v", err)
	}

	result := models.DB.Where("base_id = ?", "158a17f6-949b-4275-8706-3c2142e76a43").First(&models.BaseTask{})

	if result.Error != nil {
		t.Errorf("Failed to find base task: %v", result.Error)
	}

	if result.RowsAffected == 1 {
		var baseTask models.BaseTask
		result.Scan(&baseTask)

		if baseTask.TaskName != "洗濯物を干す" {
			t.Errorf("Base task name mismatch: expected '洗濯物を干す', got '%s'", baseTask.TaskName)
		}
	}
}
