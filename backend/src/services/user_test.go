package services_test

import (
	"backend/models"
	"backend/services"
	"testing"
)

func TestRegisterUser(t *testing.T) {

	if err := models.DB.Exec("TRUNCATE TABLE users").Error; err != nil {
		t.Fatal(err)
	}

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

	if err := models.DB.Exec("TRUNCATE TABLE users").Error; err != nil {
		t.Fatal(err)
	}

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
}

func TestRegisterUser_InvalidLivingType(t *testing.T) {

	if err := models.DB.Exec("TRUNCATE TABLE users").Error; err != nil {
		t.Fatal(err)
	}

	req := services.RegisterUserRequest{
		BirthDate:  946684800,
		LivingType: "invalid",
	}

	_, err := services.RegisterUser("user-001", "syatyo", "syatyo@example.com", req)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
