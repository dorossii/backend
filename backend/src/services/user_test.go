package services_test

import (
	"backend/models"
	"backend/services"
	"testing"
)

func truncateUsersAndRooms(t *testing.T) {
	t.Helper()

	if err := models.DB.Exec("TRUNCATE TABLE user_rooms").Error; err != nil {
		t.Fatal(err)
	}
	if err := models.DB.Exec("TRUNCATE TABLE users").Error; err != nil {
		t.Fatal(err)
	}
}

func TestRegisterUser(t *testing.T) {

	truncateUsersAndRooms(t)

	req := services.RegisterUserRequest{
		BirthDate:  946684800, // 2000-01-01 00:00:00 UTC
		LivingType: "alone",
	}

	res, err := services.RegisterUser("user-001", "syatyo", "syatyo@example.com", req)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	if res.UserID != "user-001" {
		t.Fatalf("unexpected UserID: %s", res.UserID)
	}
	if res.UserName != "syatyo" {
		t.Fatalf("unexpected UserName: %s", res.UserName)
	}
	if res.BirthDate != 946684800 {
		t.Fatalf("unexpected BirthDate: %d", res.BirthDate)
	}
	if res.LivingType != "alone" {
		t.Fatalf("unexpected LivingType: %s", res.LivingType)
	}

	var room models.UserRoom
	if err := models.DB.First(&room, "user_id = ?", "user-001").Error; err != nil {
		t.Fatalf("UserRoom not created: %v", err)
	}
	if !room.IsAlone {
		t.Fatalf("expected IsAlone=true for livingType=alone")
	}
	if !room.HasWasher || !room.HasVacuum || !room.HasRobot || !room.UseTableware || !room.HasDishwasher {
		t.Fatalf("expected default lifestyle fields to be true: %+v", room) // register時点では全てデフォルト値(true)であることを確認する
	}

	// ユーザー二人目を追加
	req = services.RegisterUserRequest{
		BirthDate:  926726400, // 1999-05-15 00:00:00 UTC
		LivingType: "alone",
	}

	res, err = services.RegisterUser("user-002", "goro", "goro@example.com", req)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	if res.UserID != "user-002" {
		t.Fatalf("unexpected UserID: %s", res.UserID)
	}
	if res.UserName != "goro" {
		t.Fatalf("unexpected UserName: %s", res.UserName)
	}
	if res.BirthDate != 926726400 {
		t.Fatalf("unexpected BirthDate: %d", res.BirthDate)
	}
	if res.LivingType != "alone" {
		t.Fatalf("unexpected LivingType: %s", res.LivingType)
	}

	// ユーザー003
	req = services.RegisterUserRequest{
		BirthDate:  946684800, // 2000-01-01 00:00:00 UTC
		LivingType: "alone",
	}

	res, err = services.RegisterUser("user-003", "usr003", "usr003@example.com", req)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	if res.UserID != "user-003" {
		t.Fatalf("unexpected UserID: %s", res.UserID)
	}

	if res.UserName != "usr003" {
		t.Fatalf("unexpected UserName: %s", res.UserName)
	}

	if res.BirthDate != 946684800 {
		t.Fatalf("unexpected BirthDate: %d", res.BirthDate)
	}

	if res.LivingType != "alone" {
		t.Fatalf("unexpected LivingType: %s", res.LivingType)
	}
}

func TestRegisterUser_Family(t *testing.T) {

	truncateUsersAndRooms(t)

	req := services.RegisterUserRequest{
		BirthDate:  926726400, // 1999-05-15 00:00:00 UTC
		LivingType: "family",
	}

	res, err := services.RegisterUser("user-002", "goro", "goro@example.com", req)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	if res.LivingType != "family" {
		t.Fatalf("unexpected LivingType: %s", res.LivingType)
	}

	var room models.UserRoom
	if err := models.DB.First(&room, "user_id = ?", "user-002").Error; err != nil {
		t.Fatalf("UserRoom not created: %v", err)
	}
	if room.IsAlone {
		t.Fatalf("expected IsAlone=false for livingType=family")
	}
}

func TestRegisterUser_InvalidLivingType(t *testing.T) {

	truncateUsersAndRooms(t)

	req := services.RegisterUserRequest{
		BirthDate:  946684800,
		LivingType: "invalid",
	}

	_, err := services.RegisterUser("user-001", "syatyo", "syatyo@example.com", req)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
