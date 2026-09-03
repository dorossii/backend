package controllers

import (
	"backend/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetNotices はアプリ内通知一覧を返すハンドラ
func GetNotices(ctx echo.Context) error {
	// JWTミドルウェアで検証済みのユーザーIDを取得
	userID := ctx.Get("UserID").(string)

	notices, err := services.GetNotices(userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// トップレベル配列として返す
	return ctx.JSON(http.StatusOK, notices)
}
