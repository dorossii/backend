package batch

import (
	"backend/models"
	"testing"
)

func TestChangeHP(t *testing.T) {
	//テスト対象の関数を実行
	err := ChangeHP()
	if err != nil {
		t.Fatalf("Failed to change HP: %v", err)
	}

	//結果の検証
	var users []models.User
	err = models.DB.Find(&users).Error
	if err != nil {
		t.Fatalf("Failed to find users: %v", err)
	}

	//ユーザーごとにHPが正しく更新されているかを検証
	for _, user := range users {
		//HPは1000を上限とし100%なため
		expectedHP := user.HealthPoint
		if user.UserID == "user-001" {
			expectedHP = 1000 // 100% -> 100% (1000 - 0)
		}
		if user.UserID == "user-002" {
			expectedHP = 760 // 70% -> 76% (700 + 60)
		}
		if user.UserID == "user-003" {
			expectedHP = 55 // 8% -> 5.5% (80 - 25)
		}
		if user.HealthPoint != expectedHP {
			t.Errorf("User %s: expected HealthPoint %d, got %d", user.UserID, expectedHP, user.HealthPoint)
		}
		
		//ユーザーデータの表示
		t.Logf("DirtLevel: %d, UserID: %s, HealthPoint: %d", user.DirtLevel, user.UserID, user.HealthPoint)
	}
}
