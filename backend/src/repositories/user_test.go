package repositories_test

import (
	"backend/models"
	"backend/repositories"
	"testing"
	"time"
)

func truncateUsers(t *testing.T) {
	t.Helper()
	if err := models.DB.Exec("TRUNCATE TABLE users").Error; err != nil {
		t.Fatal(err)
	}
}

func TestCreateUser(t *testing.T) {
	truncateUsers(t)

	user := &models.User{
		UserID:     "user-001",
		UserName:   "syatyo",
		BirthDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Mailadress: "syatyo@example.com",
	}

	err := repositories.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	var result models.User
	err = models.DB.First(&result, "user_id = ?", "user-001").Error
	if err != nil {
		t.Fatalf("failed to find User: %v", err)
	}

	if result.UserName != "syatyo" {
		t.Fatalf("unexpected UserName: %s", result.UserName)
	}
	if result.Mailadress != "syatyo@example.com" {
		t.Fatalf("unexpected Mailadress: %s", result.Mailadress)
	}
	if !result.BirthDate.Equal(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected BirthDate: %v", result.BirthDate)
	}
}

func TestCreateUser_DuplicateID(t *testing.T) {
	truncateUsers(t)

	user := &models.User{
		UserID:     "user-001",
		UserName:   "syatyo",
		BirthDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Mailadress: "syatyo@example.com",
	}

	if err := repositories.CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// 同じIDで再度登録するとエラーになるか確認
	err := repositories.CreateUser(user)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
