package controllers

import (
	"backend/models"
	"backend/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestMain(m *testing.M) {
	models.InitForTest()
	os.Exit(m.Run())
}

func TestGetUserStatus(t *testing.T) {
	models.DB.Exec("TRUNCATE TABLE user_rooms")
	models.DB.Exec("TRUNCATE TABLE users")

	req := services.RegisterUserRequest{
		BirthDate:  946684800,
		LivingType: "alone",
	}
	if _, err := services.RegisterUser("user-001", "syatyo", "syatyo@example.com", req); err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	e := echo.New()
	httpReq := httptest.NewRequest(http.MethodGet, "/user/status", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(httpReq, rec)
	ctx.Set("UserID", "user-001")

	if err := GetUserStatus(ctx); err != nil {
		t.Fatalf("GetUserStatus returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res services.UserStatusResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if res.UserID != "user-001" {
		t.Fatalf("unexpected UserID: %s", res.UserID)
	}
	if res.UserName != "syatyo" {
		t.Fatalf("unexpected UserName: %s", res.UserName)
	}
}

func TestGetUserStatus_NotFound(t *testing.T) {
	models.DB.Exec("TRUNCATE TABLE user_rooms")
	models.DB.Exec("TRUNCATE TABLE users")

	e := echo.New()
	httpReq := httptest.NewRequest(http.MethodGet, "/user/status", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(httpReq, rec)
	ctx.Set("UserID", "not-exist-user")

	if err := GetUserStatus(ctx); err != nil {
		t.Fatalf("GetUserStatus returned error: %v", err)
	}

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}
