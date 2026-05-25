package batch_test

import (
	"backend/models"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	models.InitForTest()
	// テスト用データの準備 (ユーザー2人、ベースタスク3つ)
	if err := createTestUser(); err != nil {
		panic(fmt.Sprintf("failed to create test users: %v", err))
	}

	os.Exit(m.Run())
}

func createTestUser() error {
		// テスト用データの準備 (ユーザー2人)
	users := []models.User{
		{
			UserID:      "user-001",
			UserName:    "syatyo",
			BirthDate:   time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
			Mailadress:  "user1@example.com",
			HealthPoint: 1000,
			DirtLevel:   0,
			Combo:       0,
			BgColor:     "#ffb6c1",
		},
		{
			UserID:      "user-002",
			UserName:    "goro",
			BirthDate:   time.Date(2004, 2, 2, 0, 0, 0, 0, time.UTC),
			Mailadress:  "user2@example.com",
			HealthPoint: 700,
			DirtLevel:   10,
			Combo:       1,
			BgColor:     "#add8e6",
		},
		{
			UserID:      "user-003",
			UserName:    "jiro",
			BirthDate:   time.Date(2004, 3, 3, 0, 0, 0, 0, time.UTC),
			Mailadress:  "user3@example.com",
			HealthPoint: 80,
			DirtLevel:   700,
			Combo:       2,
			BgColor:     "#90ee90",
		},
	}
	if err := models.DB.Create(&users).Error; err != nil {
		return fmt.Errorf("failed to create dummy users: %v", err)
	}
	return nil
}
