package middlewares

import (
	"backend/logger"
	"os"

	"github.com/gorilla/sessions"
)

const AdminSessionName = "admin-session"

var Store *sessions.CookieStore

func InitSession() {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		logger.PrintErr("SESSION_SECRETが設定されていません")
		return
	}

	Store = sessions.NewCookieStore([]byte(secret))
}
