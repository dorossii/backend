package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RequireAdminSession は管理画面APIのセッション認証状態を確認するミドルウェア
func RequireAdminSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		session, err := Store.Get(ctx.Request(), AdminSessionName)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		if authenticated, ok := session.Values["authenticated"].(bool); !ok || !authenticated {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		return next(ctx)
	}
}

// RequireAdminSessionPage は管理画面の画面表示用に、未認証時はログインページへリダイレクトするミドルウェア
func RequireAdminSessionPage(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		session, err := Store.Get(ctx.Request(), AdminSessionName)
		if err != nil {
			return ctx.Redirect(http.StatusFound, "login.html")
		}

		if authenticated, ok := session.Values["authenticated"].(bool); !ok || !authenticated {
			return ctx.Redirect(http.StatusFound, "login.html")
		}

		return next(ctx)
	}
}
