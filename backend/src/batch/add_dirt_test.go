package batch_test

import (
	"backend/batch"
	"backend/models"
	"testing"
)


// AddDirt() 単体テスト
func TestAddDirt(t *testing.T) {

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


