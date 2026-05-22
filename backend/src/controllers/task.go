package controllers

import "github.com/labstack/echo/v4"


func GetTask(ctx echo.Context) error {
	//ヘッダーからユーザーIDを取得
	userId := ctx.Get("UserID").(string)

	

	return ctx.JSON(200, map[string]string{
		"message": "get task",
		"userId": userId,
	})	
}
