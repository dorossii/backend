package controllers

import (
	"backend/services"
	"net/http"

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

type PutTaskStatusRequest struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
}

// タスクステータス更新(完了•未完了)
func PutTaskStatusHandler(ctx echo.Context) error {
	//ヘッダーからユーザーIDを取得
	id := ctx.Request().Header.Get("UserID")
	taskID := ctx.Param("id")

	var req PutTaskStatusRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request",
		})
	}

	result, err := services.PutTaskStatus(id, taskID, req.Status, req.Message)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"isChanged":    result.IsChanged,
		"requireImage": result.RequireImage,
	})
}
