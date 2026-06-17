package controllers

import (
	"backend/logger"
	"backend/services"

	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SendFriendRequest(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)

	friendID := ctx.Request().Header.Get("InviteID")
	if friendID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "InviteID header is required"})
	}

	err := services.SendFriendRequest(userId, friendID)
	if err != nil {
		if errors.Is(err, services.ErrAlreadySent) || errors.Is(err, services.ErrAlreadyReceived) {
			return ctx.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}

func AcceptFriendRequest(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)

	// 相手のIDを取得
	friendID := ctx.Request().Header.Get("FriendId")
	if friendID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "FriendId header is required"})
	}

	err := services.AcceptFriendRequest(userId, friendID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}

// 送られてきたフレンドリクエストを取得する関数
func GetReceivedFriendRequests(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)

	// 相手のIDを取得
	res, err := services.GetFriendRequests(userId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, res)
}


// GetFriends は承認済みフレンドの一覧を返すハンドラ
func GetFriends(ctx echo.Context) error {
	// JWTミドルウェアで検証済みのユーザーIDを取得
	userId := ctx.Get("UserID").(string)

	// フレンド一覧を取得
	friends, err := services.GetFriends(userId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"friends": friends})
}

// DeleteFriend はフレンド関係を削除するハンドラ
func DeleteFriend(ctx echo.Context) error {
	// JWTミドルウェアで検証済みのユーザーIDを取得
	userId := ctx.Get("UserID").(string)

	// パスパラメータから削除対象のフレンドIDを取得
	friendID := ctx.Param("id")
	if friendID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	// フレンド関係を削除
	err := services.DeleteFriend(userId, friendID)
	if err != nil {
		if errors.Is(err, services.ErrFriendShipNotFound) || errors.Is(err, services.ErrFriendShipNotAccepted) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{})
}

func GetInviteURL(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)

	url := fmt.Sprintf("%s/app/friend/invite?inviteid=%s", inviteBaseURL, userId)

	return ctx.JSON(http.StatusOK, map[string]string{
		"URL": url,
	})
}

type PostAttackerSettingsRequest struct {
	TargetUser string `json:"TargetUser"`
}

// 嫌がらせする人の設定
func PostAttackerSettingsHandler(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)

	var req PostAttackerSettingsRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request",
		})
	}

	err := services.PostAttackerSettings(userId, req.TargetUser)
	if err != nil {
		logger.PrintErr("PostAttackerSettingsHandler", err)

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

// レスキューする人の設定
func PostRescuerSettingsHandler(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)

	var req struct {
		TargetUsers []string `json:"TargetUsers"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request",
		})
	}
	
	// ターゲットユーザーの設定
	err := services.PostRescuerSettings(userId, req.TargetUsers)

	if err != nil {
		logger.PrintErr("PostRescuerSettingsHandler", err)
		
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
