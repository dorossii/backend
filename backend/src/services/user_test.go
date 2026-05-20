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
		UserName:    "syatyo",
		BirthDate:   946684800, // 2000-01-01 00:00:00 UTC
		MailAddress: "syatyo@example.com",
		LivingType:  "alone",
	}

	res, err := services.RegisterUser("user-001", req)
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
	if res.MailAddress != "syatyo@example.com" {
		t.Fatalf("unexpected MailAddress: %s", res.MailAddress)
	}
	if res.LivingType != "alone" {
		t.Fatalf("unexpected LivingType: %s", res.LivingType)
	}
}

func TestRegisterUser_Together(t *testing.T) {

	if err := models.DB.Exec("TRUNCATE TABLE users").Error; err != nil {
		t.Fatal(err)
	}

	req := services.RegisterUserRequest{
		UserName:    "goro",
		BirthDate:   926726400, // 1999-05-15 00:00:00 UTC
		MailAddress: "goro@example.com",
		LivingType:  "together",
	}

	res, err := services.RegisterUser("user-002", req)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	if res.LivingType != "together" {
		t.Fatalf("unexpected LivingType: %s", res.LivingType)
	}
}

func TestRegisterUser_InvalidLivingType(t *testing.T) {

	if err := models.DB.Exec("TRUNCATE TABLE users").Error; err != nil {
		t.Fatal(err)
	}

	req := services.RegisterUserRequest{
		UserName:    "syatyo",
		BirthDate:   946684800,
		MailAddress: "syatyo@example.com",
		LivingType:  "invalid",
	}

	_, err := services.RegisterUser("user-001", req)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

