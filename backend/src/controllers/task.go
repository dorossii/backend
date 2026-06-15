package controllers

import (
	"backend/services"
	"errors"
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

type PostTaskTauntMessageRequest struct {
	FriendID string `json:"friendId"`
	Message  string `json:"message"`
}

// 煽りメッセージの登録
func PostTauntMessageHandler(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)
	
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

// タスク写真アップロード
func PostUploadImageHandler(ctx echo.Context) error {

	userId := ctx.Get("UserID").(string)
	taskID := ctx.Param("id")

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": "画像が選択されていません",
		})
	}

	const maxSize = 10 * 1024 * 1024 // 10MB

	if fileHeader.Size > maxSize {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": "画像サイズが大きすぎます",
		})
	}

	err = services.PostUploadImage(userId, taskID, fileHeader)
	if err != nil {

		switch {

		case errors.Is(err, services.ErrTaskPermissionDenied):
			return ctx.JSON(http.StatusForbidden, echo.Map{
				"error": err.Error(),
			})

		case errors.Is(err, services.ErrUnsupportedImageType):
			return ctx.JSON(http.StatusBadRequest, echo.Map{
				"error": err.Error(),
			})

		case errors.Is(err, services.ErrEmptyImageFile):
			return ctx.JSON(http.StatusBadRequest, echo.Map{
				"error": err.Error(),
			})

		default:
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "internal server error",
			})
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}