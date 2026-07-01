package controllers

import (
	// "backend/logger"
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

	return ctx.JSON(200, tasks)
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
		case errors.Is(err, services.ErrTaskNotFound):
			return ctx.JSON(http.StatusNotFound, echo.Map{
				"error": err.Error(),
			})

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


type PutTaskStatusRequest struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// タスクステータス更新(完了•未完了)
func PutTaskStatusHandler(ctx echo.Context) error {
	//ヘッダーからユーザーIDを取得
	userId := ctx.Get("UserID").(string)
	taskID := ctx.Param("id")

	var req PutTaskStatusRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request",
		})
	}

	// TODO:ここのエラーハンドリング修正の余地あり（ケースで分けてるのが多すぎる気がするから）
	result, err := services.PutTaskStatus(userId, taskID, req.Status, req.Message)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidTaskStatus):
			return ctx.JSON(http.StatusBadRequest, echo.Map{
				"error": "invalid task status",
			})

		case errors.Is(err, services.ErrInvalidRequest):
			return ctx.JSON(http.StatusBadRequest, echo.Map{
				"error": "invalid request",
			})

		case errors.Is(err, services.ErrTaskNotFound):
			return ctx.JSON(http.StatusNotFound, echo.Map{
				"error": "task not found",
			})

		case errors.Is(err, services.ErrTaskExpired):
			return ctx.JSON(http.StatusConflict, echo.Map{
				"error": "task expired",
			})

		case errors.Is(err, services.ErrInvalidTaskState):
			return ctx.JSON(http.StatusConflict, echo.Map{
				"error": "invalid task state",
			})

		case errors.Is(err, services.ErrTaskPermissionDenied):
			return ctx.JSON(http.StatusForbidden, echo.Map{
				"error": "forbidden",
			})

		default:
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "internal server error",
			})
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"isChanged":    result.IsChanged,
		"requireImage": result.RequireImage,
	})
}

// タスク写真取得
func GetTaskImageHandler(ctx echo.Context) error {
	imageID := ctx.Param("imageId")

	filePath, _, err := services.GetTaskImage(imageID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrTaskNotFound):
			return ctx.JSON(http.StatusNotFound, echo.Map{
				"error": err.Error(),
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "internal server error",
			})
		}
	}

	return ctx.File(filePath)
}
