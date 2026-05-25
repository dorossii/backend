package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestGetInviteURL(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/friend/invite", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("UserID", "0182624a-a06f-4cd2-811b-1dd75a554e7f")

	if err := GetInviteURL(ctx); err != nil {
		t.Fatalf("GetInviteURL returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	url, ok := res["URL"]
	if !ok {
		t.Fatal("response does not contain URL field")
	}

	expected := inviteBaseURL + "/app/friend/invite?inviteid=0182624a-a06f-4cd2-811b-1dd75a554e7f"
	if url != expected {
		t.Fatalf("unexpected URL: got %s, want %s", url, expected)
	}
}

func TestGetInviteURL_BaseURLFromEnv(t *testing.T) {
	original := inviteBaseURL
	inviteBaseURL = "https://example.com"
	defer func() { inviteBaseURL = original }()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/friend/invite", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("UserID", "test-user-id")

	if err := GetInviteURL(ctx); err != nil {
		t.Fatalf("GetInviteURL returned error: %v", err)
	}

	var res map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	expected := "https://example.com/app/friend/invite?inviteid=test-user-id"
	if res["URL"] != expected {
		t.Fatalf("unexpected URL: got %s, want %s", res["URL"], expected)
	}
}
