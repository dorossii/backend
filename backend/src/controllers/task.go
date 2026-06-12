package controllers

import (
	"backend/services"
	"net/http"
	"errors"

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

type PostTaskTauntMessageRequest struct {
	FriendID string `json:"friendId"`
	Message  string `json:"message"`
}

// 煽りメッセージの登録
func PostTauntMessageHandler(ctx echo.Context) error {
	userId := ctx.Request().Header.Get("UserID")

	var req PostTaskTauntMessageRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request",
		})
	}

	err := services.PostTaskTauntMessage(userId, req.FriendID, req.Message)
	if err != nil {
		if errors.Is(err, services.ErrFriendNotFound) {
			return ctx.JSON(http.StatusForbidden, echo.Map{
				"error": "friend not found",
			})
		}

		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error": "internal server error",
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}
