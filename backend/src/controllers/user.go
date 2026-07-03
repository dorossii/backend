package controllers

import (
	"backend/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterUser(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)
	name := ctx.Get("Name").(string)
	email := ctx.Get("Email").(string)

	var req services.RegisterUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	res, err := services.RegisterUser(userId, name, email, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, res)
}

func CreateUserLifestyle(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)

	var req services.LifestyleRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	res, err := services.CreateUserLifestyle(userId, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, res)
}

func UpdateUserLifestyle(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)

	var req services.LifestyleRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := services.UpdateUserLifestyle(userId, req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{})
}

// UpdateUserSetting は PUT /app/user/setting のハンドラ
func UpdateUserSetting(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string) // 認証ミドルウェアが格納したユーザーIDを取得

	var req services.UserSettingRequest
	if err := ctx.Bind(&req); err != nil { // リクエストボディをバインド
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := services.UpdateUserSetting(userId, req); err != nil { // サービス層に処理を委譲
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{}) // Empty スキーマを返す
}
