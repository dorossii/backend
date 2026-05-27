package controllers

import (
	"backend/services"

	"github.com/labstack/echo/v4"
)

const ErrTaskFetch = "タスクの取得に失敗しました"

// タスク取得のコントローラー
func GetTask(ctx echo.Context) error {
	//ヘッダーからユーザーIDを取得
	userId := ctx.Get("UserID").(string)

	// サービスからタスクを取得
	tasks, err := services.GetTasks(userId)
	if err != nil {
		return ctx.JSON(500, map[string]string{
			"message": ErrTaskFetch,
		})
	}

	return ctx.JSON(200, map[string]interface{}{
		"tasks": tasks,
	})	
}
