package batch_test

import (
	"backend/batch"
	"backend/models"
	"testing"
	"time"
)


// AddDirt() 単体テスト
func TestAddDirt(t *testing.T) {

		// テスト用データの準備 (ユーザー2人、ベースタスク3つ)
	users := []models.User{
		{
			UserID:      "user-001",
			UserName:    "syatyo",
			BirthDate:   time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
			Mailadress:  "user1@example.com",
			HealthPoint: 100,
			DirtLevel:   0,
			Combo:       0,
			BgColor:     "#ffb6c1",
		},
		{
			UserID:      "user-002",
			UserName:    "goro",
			BirthDate:   time.Date(2004, 2, 2, 0, 0, 0, 0, time.UTC),
			Mailadress:  "user2@example.com",
			HealthPoint: 90,
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
		t.Fatalf("failed to create dummy users: %v", err)
	}

	//テスト対象の関数を実行
	err := batch.AddDirt()
	if err != nil {
		t.Fatalf("AddDirt() returned an error: %v", err)
	}

	// 結果の検証
	var updatedUsers []models.User
	err = models.DB.Find(&updatedUsers).Error
	if err != nil {
		t.Fatalf("failed to find updated users: %v", err)
	}

	// ユーザーごとにDirtLevelが正しく更新されているかを検証
	for _, user := range updatedUsers {
		expectedDirtLevel := user.DirtLevel
		if user.UserID == "user-001" {
			expectedDirtLevel = 3 // 初期値0 + 3
		} else if user.UserID == "user-002" {
			expectedDirtLevel = 13 // 初期値10 + 3
		}
		if user.DirtLevel != expectedDirtLevel {
			t.Errorf("User %s: expected DirtLevel %d, got %d", user.UserID, expectedDirtLevel, user.DirtLevel)
		}
		//ユーザーデータの表示
		t.Logf("UserID: %s, DirtLevel: %d", user.UserID, user.DirtLevel)
	}
	// ユーザー3はDirtLevelが700以上なので、更新後も700のままであることを検証
	for _, user := range updatedUsers {
		if user.UserID == "user-003" {
			if user.DirtLevel != 700 {
				t.Errorf("User %s: expected DirtLevel to remain 700, got %d", user.UserID, user.DirtLevel)
			}
		}
	}
}


